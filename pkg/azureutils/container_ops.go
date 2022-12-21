// Copyright 2021 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package azureutils

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-cosi-driver/pkg/types"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
	azure "sigs.k8s.io/cloud-provider-azure/pkg/provider"
)

const (
	AccessKey = ""
)

var (
	storageAccountRE = regexp.MustCompile(`https://(.+).blob.core.windows.net/([^/]*)/?(.*)`)
)

func createContainerBucket(
	ctx context.Context,
	bucketName string,
	parameters *BucketClassParameters,
	cloud *azure.Cloud) (string, error) {
	accOptions := getAccountOptions(parameters)
	_, key, err := cloud.EnsureStorageAccount(ctx, accOptions, "")
	if err != nil {
		return "", status.Error(codes.Internal, fmt.Sprintf("Could not ensure storage account %s exists: %v", accOptions.Name, err))
	}
	containerParams := make(map[string]string) //NOTE: Container parameters still need to be filled/implemented

	container, err := createAzureContainer(ctx, parameters.storageAccountName, key, bucketName, containerParams)
	if err != nil {
		return "", err
	}

	id := types.BucketID{
		ResourceGroup: parameters.resourceGroup,
		URL:           container,
	}
	if parameters.subscriptionID != "" {
		id.SubID = parameters.subscriptionID
	} else {
		id.SubID = cloud.SubscriptionID
	}
	base64ID, err := id.Encode()
	if err != nil {
		return "", status.Error(codes.InvalidArgument, fmt.Sprintf("could not encode ID: %v", err))
	}

	return base64ID, nil
}

func DeleteContainerBucket(
	ctx context.Context,
	bucketID *types.BucketID,
	cloud *azure.Cloud) error {
	// Get storage account name from bucket url
	storageAccountName := getStorageAccountNameFromContainerURL(bucketID.URL)
	// Get access keys for the storage account
	accessKey, err := cloud.GetStorageAccesskey(ctx, bucketID.SubID, storageAccountName, bucketID.ResourceGroup)
	if err != nil {
		return err
	}

	containerName := getContainerNameFromContainerURL(bucketID.URL)
	err = deleteAzureContainer(ctx, storageAccountName, accessKey, containerName)
	if err != nil {
		return fmt.Errorf("Error deleting container %s in storage account %s : %v", containerName, storageAccountName, err)
	}

	// Now, we check and delete the storage account if its empty
	return nil
}

func getStorageAccountNameFromContainerURL(containerURL string) string {
	storageAccountName, _, _, err := parseContainerURL(containerURL)
	if err != nil {
		return ""
	}

	return storageAccountName
}

func getContainerNameFromContainerURL(containerURL string) string {
	_, containerName, _, err := parseContainerURL(containerURL)
	if err != nil {
		return ""
	}

	return containerName
}

func deleteAzureContainer(
	ctx context.Context,
	storageAccount,
	accessKey,
	containerName string) error {
	containerClient, err := createContainerClient(storageAccount, accessKey, containerName)

	if err != nil {
		return err
	}

	_, err = containerClient.Delete(ctx, nil)
	return err
}

func createContainerClient(
	storageAccount string,
	accessKey string,
	containerName string) (*container.Client, error) {
	// Create credentials
	credential, err := container.NewSharedKeyCredential(storageAccount, accessKey)
	if err != nil {
		return nil, fmt.Errorf("Invalid credentials with error : %v", err)
	}

	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", storageAccount, containerName)

	containerClient, err := container.NewClientWithSharedKeyCredential(containerURL, credential, nil)

	return containerClient, err
}

func parseContainerURL(containerURL string) (string, string, string, error) {
	matches := storageAccountRE.FindStringSubmatch(containerURL)
	if len(matches) < 2 {
		errStr := fmt.Sprintf("Invalid URL has been passed: %s", containerURL)
		klog.Errorf("Error in parseContainerURL :: %s", errStr)
		return "", "", "", errors.New(errStr)
	}

	storageAccount := matches[1]
	containerName := ""
	if len(matches) > 2 {
		containerName = matches[2]
	}

	blobName := ""
	if len(matches) > 3 {
		blobName = matches[3]
	}

	return storageAccount, containerName, blobName, nil
}

func createAzureContainer(
	ctx context.Context,
	storageAccount string,
	accessKey string,
	containerName string,
	parameters map[string]string) (string, error) {
	if len(storageAccount) == 0 || len(accessKey) == 0 {
		return "", fmt.Errorf("Invalid storage account or access key")
	}

	containerClient, err := createContainerClient(storageAccount, accessKey, containerName)
	if err != nil {
		return "", err
	}

	// Lets create a container with the containerClient
	_, err = containerClient.Create(ctx, &container.CreateOptions{
		Metadata: parameters,
		Access:   nil,
	})
	if err != nil {
		if err.Error() == "ResourceExistsError" {
			return containerClient.URL(), nil
		}
		return "", fmt.Errorf("Error creating container from containterURL : %s, Error : %v", containerClient.URL(), err)
	}

	return containerClient.URL(), nil
}

func createContainerSASURL(ctx context.Context, bucketID string, parameters *BucketAccessClassParameters, accountKey string) (string, string, error) {
	account, containerName, _, err := parseContainerURL(bucketID)
	if err != nil {
		return "", "", err
	}

	if containerName == "" {
		return "", "", fmt.Errorf("Error in createContainerSASURL as containerName is empty for bucketID: %s", bucketID)
	}

	cred, err := container.NewSharedKeyCredential(account, accountKey)
	if err != nil {
		return "", "", err
	}

	permission := sas.ContainerPermissions{}
	permission.List = parameters.enableList
	permission.Read = parameters.enableRead
	permission.Write = parameters.enableWrite
	permission.Delete = parameters.enableDelete
	permission.DeletePreviousVersion = parameters.enablePermanentDelete
	permission.Add = parameters.enableAdd
	permission.FilterByTags = parameters.enableTags

	start := time.Now()
	expiry := start.Add(time.Millisecond * time.Duration(parameters.validationPeriod))

	sasQueryParams, err := sas.BlobSignatureValues{
		Protocol:      sas.Protocol(parameters.signedProtocol),
		StartTime:     start,
		ExpiryTime:    expiry,
		Permissions:   permission.String(),
		IPRange:       sas.IPRange(parameters.signedIP),
		Version:       parameters.signedversion,
		ContainerName: containerName,
	}.SignWithSharedKey(cred)

	if err != nil {
		return "", "", err
	}

	queryParams := sasQueryParams.Encode()
	accountID := fmt.Sprintf("https://%s.blob.core.windows.net/", account)
	sasURL := fmt.Sprintf("%s?%s", accountID, queryParams)
	return sasURL, accountID, nil
}
