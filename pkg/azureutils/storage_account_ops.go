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
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/Azure/azure-cosi-driver/pkg/types"
	azure "sigs.k8s.io/cloud-provider-azure/pkg/provider"
)

func DeleteStorageAccount(
	ctx context.Context,
	id *types.BucketID,
	cloud *azure.Cloud) error {
	SAClient := cloud.StorageAccountClient
	err := SAClient.Delete(ctx, id.SubID, id.ResourceGroup, getStorageAccountNameFromContainerURL(id.URL))
	if err != nil {
		return err.Error()
	}
	return nil
}

func createStorageAccountBucket(ctx context.Context,
	bucketName string,
	parameters *BucketClassParameters,
	cloud *azure.Cloud) (string, error) {
	accName, _, err := cloud.EnsureStorageAccount(ctx, getAccountOptions(parameters), "")
	if err != nil {
		return "", status.Error(codes.Internal, fmt.Sprintf("Could not create storage account: %v", err))
	}

	accURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accName)

	id := types.BucketID{
		ResourceGroup: parameters.resourceGroup,
		URL:           accURL,
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

// creates SAS and returns service client with sas
func createAccountSASURL(ctx context.Context, bucketID string, parameters *BucketAccessClassParameters, accountKey string) (string, string, error) {
	account := getStorageAccountNameFromContainerURL(bucketID)
	cred, err := azblob.NewSharedKeyCredential(account, accountKey)
	if err != nil {
		return "", "", err
	}

	resources := sas.AccountResourceTypes{}
	if parameters.allowServiceSignedResourceType {
		resources.Service = true
	}
	if parameters.allowContainerSignedResourceType {
		resources.Container = true
	}
	if parameters.allowObjectSignedResourceType {
		resources.Object = true
	}

	permission := sas.AccountPermissions{}
	permission.List = parameters.enableList
	permission.Read = parameters.enableRead
	permission.Write = parameters.enableWrite
	permission.Delete = parameters.enableDelete
	permission.DeletePreviousVersion = parameters.enablePermanentDelete
	permission.Add = parameters.enableAdd
	permission.Tag = parameters.enableTags
	permission.FilterByTags = parameters.enableFilter

	start := time.Now()
	expiry := start.Add(time.Millisecond * time.Duration(parameters.validationPeriod))
	services := &sas.AccountServices{Blob: true}
	sasQueryParams := sas.AccountSignatureValues{
		Protocol:      parameters.signedProtocol,
		StartTime:     start,
		ExpiryTime:    expiry,
		Permissions:   permission.String(),
		ResourceTypes: resources.String(),
		Services:      services.String(),
		IPRange:       parameters.signedIP,
		Version:       parameters.signedversion,
	}

	queryParams, err := sasQueryParams.SignWithSharedKey(cred)
	if err != nil {
		return "", "", err
	}
	sasURL := fmt.Sprintf("%s/?%s", strings.TrimSuffix(bucketID, "/"), queryParams.Encode())
	return sasURL, bucketID, nil
}
