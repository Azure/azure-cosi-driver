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

package provisionerserver

import (
	"context"
	"fmt"
	"Azure/azure-cosi-driver/pkg/azureutils"
	"Azure/azure-cosi-driver/pkg/constant"
	"reflect"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
	azure "sigs.k8s.io/cloud-provider-azure/pkg/provider"
	spec "sigs.k8s.io/container-object-storage-interface-spec"
)

type bucketDetails struct {
	bucketID   string
	parameters map[string]string
}

type provisioner struct {
	spec.UnimplementedProvisionerServer

	bucketsLock       sync.RWMutex
	nameToBucketMap   map[string]*bucketDetails
	bucketIDToNameMap map[string]string
	cloud             *azure.Cloud
}

var _ spec.ProvisionerServer = &provisioner{}

func NewProvisionerServer(
	kubeconfig,
	cloudConfigSecretName,
	cloudConfigSecretNamespace string) (spec.ProvisionerServer, error) {
	kubeClient, err := azureutils.GetKubeClient(kubeconfig)
	if err != nil {
		return nil, err
	}
	klog.Infof("Kubeclient : %+v", kubeClient)

	azCloud, err := azureutils.GetAzureCloudProvider(kubeClient, cloudConfigSecretName, cloudConfigSecretNamespace)
	if err != nil {
		return nil, err
	}

	return &provisioner{
		nameToBucketMap:   make(map[string]*bucketDetails),
		bucketsLock:       sync.RWMutex{},
		bucketIDToNameMap: make(map[string]string),
		cloud:             azCloud,
	}, nil
}

func (pr *provisioner) DriverCreateBucket(
	ctx context.Context,
	req *spec.DriverCreateBucketRequest) (*spec.DriverCreateBucketResponse, error) {

	bucketName := req.GetName()
	parameters := req.GetParameters()
	if parameters == nil {
		return nil, status.Error(codes.InvalidArgument, "Parameters missing. Cannot initialize Azure bucket.")
	}

	// Check if a bucket with these set of values exist in the namesToBucketMap
	pr.bucketsLock.RLock()
	currBucket, exists := pr.nameToBucketMap[bucketName]
	pr.bucketsLock.RUnlock()

	if exists {
		bucketParams := currBucket.parameters
		if bucketParams == nil {
			bucketParams = make(map[string]string)
		}
		if reflect.DeepEqual(bucketParams, parameters) {
			return &spec.DriverCreateBucketResponse{
				BucketId: currBucket.bucketID,
			}, nil
		}

		return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("Bucket %s exists with different parameters", bucketName))
	}

	bucketID, err := azureutils.CreateBucket(ctx, bucketName, parameters, pr.cloud)
	if err != nil {
		return nil, err
	}

	// Insert the bucket into the namesToBucketMap
	pr.bucketsLock.RLock()
	pr.nameToBucketMap[bucketName] = &bucketDetails{
		bucketID:   bucketID,
		parameters: parameters,
	}
	pr.bucketIDToNameMap[bucketID] = bucketName
	pr.bucketsLock.RUnlock()

	klog.Infof("DriverCreateBucket :: Bucket id :: %s", bucketID)

	return &spec.DriverCreateBucketResponse{
		BucketId: bucketID,
	}, nil
}

func (pr *provisioner) DriverDeleteBucket(
	ctx context.Context,
	req *spec.DriverDeleteBucketRequest) (*spec.DriverDeleteBucketResponse, error) {
	//determine if the bucket is an account or a blob container
	bucketID := req.BucketId
	err := azureutils.DeleteBucket(ctx, bucketID, pr.cloud)
	if err != nil {
		return nil, err
	}

	klog.Infof("DriverDeleteBucket :: Bucket id :: %s", bucketID)
	if bucketName, ok := pr.bucketIDToNameMap[bucketID]; ok {
		// Remove from the namesToBucketMap
		pr.bucketsLock.RLock()
		delete(pr.nameToBucketMap, bucketName)
		delete(pr.bucketIDToNameMap, bucketID)
		pr.bucketsLock.RUnlock()
	}

	return &spec.DriverDeleteBucketResponse{}, nil
}

func (pr *provisioner) DriverGrantBucketAccess(
	ctx context.Context,
	req *spec.DriverGrantBucketAccessRequest) (*spec.DriverGrantBucketAccessResponse, error) {
	bucketID := req.GetBucketId()
	parameters := req.GetParameters()
	if parameters == nil {
		return nil, status.Error(codes.InvalidArgument, "Parameters missing. Cannot initialize Azure bucket.")
	}

	if req.AuthenticationType == spec.AuthenticationType_UnknownAuthenticationType {
		return nil, status.Error(codes.InvalidArgument, "AuthenticationType not provided in GrantBucketAccess request.")
	}

	var token string
	var err error

	klog.Infof("DriverGrantBucketAccess :: Bucket id :: %s", bucketID)
	if req.AuthenticationType == spec.AuthenticationType_IAM {
		return nil, status.Error(codes.Unimplemented, "AuthenticationType IAM not implemented.")
	} else if req.AuthenticationType == spec.AuthenticationType_Key {
		token, _, err = azureutils.CreateBucketSASURL(ctx, bucketID, parameters, pr.cloud)
		if err != nil {
			return nil, err
		}
	}

	return &spec.DriverGrantBucketAccessResponse{
		AccountId: req.GetName(),
		Credentials: map[string]*spec.CredentialDetails{constant.CredentialType: {
			Secrets: map[string]string{constant.AccessToken: token},
		}},
	}, nil
}

func (pr *provisioner) DriverRevokeBucketAccess(
	ctx context.Context,
	req *spec.DriverRevokeBucketAccessRequest) (*spec.DriverRevokeBucketAccessResponse, error) {
	return &spec.DriverRevokeBucketAccessResponse{}, nil
}
