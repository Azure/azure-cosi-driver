package provisionerserver

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"github.com/Azure/azure-cosi-driver/pkg/constant"
	"github.com/Azure/azure-cosi-driver/pkg/types"
	"reflect"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"sigs.k8s.io/cloud-provider-azure/pkg/azureclients/storageaccountclient/mockstorageaccountclient"
	azure "sigs.k8s.io/cloud-provider-azure/pkg/provider"
	"sigs.k8s.io/cloud-provider-azure/pkg/retry"
	spec "sigs.k8s.io/container-object-storage-interface-spec"
)

func NewMockSAClient(ctx context.Context, ctrl *gomock.Controller, subsID, rg, accName string, keyList *[]storage.AccountKey) *mockstorageaccountclient.MockInterface {
	cl := mockstorageaccountclient.NewMockInterface(ctrl)

	cl.EXPECT().
		Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(constant.ValidAccount)).
		Return(nil).
		AnyTimes()

	cl.EXPECT().
		Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(constant.ValidAccountURL)).
		Return(nil).
		AnyTimes()

	cl.EXPECT().
		Delete(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(constant.ValidAccount)).
		Return(retry.GetError(&http.Response{}, status.Error(codes.NotFound, "could not find storage account"))).
		AnyTimes()

	cl.EXPECT().
		Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Eq(constant.ValidAccount), gomock.Any()).
		Return(nil).
		AnyTimes()

	cl.EXPECT().
		Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(constant.ValidAccount), gomock.Any()).
		Return(retry.GetError(&http.Response{}, status.Error(codes.NotFound, "could not find storage account"))).
		AnyTimes()

	accountList := []storage.Account{{Name: to.StringPtr(constant.ValidAccount), AccountProperties: &storage.AccountProperties{}}}
	cl.EXPECT().
		ListByResourceGroup(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(accountList, nil).
		AnyTimes()

	cl.EXPECT().
		ListKeys(gomock.Any(), gomock.Any(), gomock.Any(), constant.ValidAccount).
		Return(storage.AccountListKeysResult{Keys: keyList}, nil).
		AnyTimes()

	cl.EXPECT().
		ListKeys(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(constant.ValidAccount)).
		Return(storage.AccountListKeysResult{Keys: keyList}, retry.GetError(&http.Response{}, fmt.Errorf("Invalid Account"))).
		AnyTimes()

	return cl
}

func newFakeProvisioner(ctrl *gomock.Controller) spec.ProvisionerServer {
	cloud := azure.GetTestCloud(ctrl)
	keyList := make([]storage.AccountKey, 0)
	keyList = append(keyList, storage.AccountKey{KeyName: to.StringPtr(constant.ValidAccount), Value: to.StringPtr(base64.StdEncoding.EncodeToString([]byte{1, 2, 3, 4}))})
	cloud.StorageAccountClient = NewMockSAClient(context.Background(), ctrl, "", "", "", &keyList)

	return &provisioner{
		nameToBucketMap:   make(map[string]*bucketDetails),
		bucketsLock:       sync.RWMutex{},
		bucketIDToNameMap: make(map[string]string),
		cloud:             cloud,
	}
}

func TestDriverCreateBucket(t *testing.T) {
	tests := []struct {
		testName    string
		bucketName  string
		params      map[string]string
		expectedErr error
	}{
		{
			testName:    "Missing Parameters",
			bucketName:  constant.ValidAccount,
			expectedErr: status.Error(codes.InvalidArgument, "Parameters missing. Cannot initialize Azure bucket."),
		},
		{
			testName:   "Create New Storage Account Bucket",
			bucketName: constant.ValidAccount,
			params: map[string]string{
				constant.BucketUnitTypeField:     constant.StorageAccount.String(),
				constant.StorageAccountNameField: constant.ValidAccount,
			},
			expectedErr: nil,
		},
		{
			testName:   "Create New Container Bucket (Could not ensure storage account)",
			bucketName: constant.ValidContainer,
			params: map[string]string{
				constant.BucketUnitTypeField:     constant.Container.String(),
				constant.StorageAccountNameField: constant.InvalidAccount,
			},
			expectedErr: status.Error(codes.Internal, fmt.Sprintf("Could not ensure storage account %s exists: %v", constant.InvalidAccount,
				fmt.Errorf("could not get storage key for storage account %s: %w", constant.InvalidAccount,
					retry.GetError(&http.Response{}, fmt.Errorf("Invalid Account")).Error()))),
		},
	}

	ctrl := gomock.NewController(t)
	pr := newFakeProvisioner(ctrl)

	for _, test := range tests {
		resp, err := pr.DriverCreateBucket(context.Background(), &spec.DriverCreateBucketRequest{
			Name:       test.bucketName,
			Parameters: test.params,
		})
		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("\nTestCase: %s\nexpected: %v\nactual: %v", test.testName, test.expectedErr, err)
		}
		if err == nil && reflect.DeepEqual(nil, resp) {
			t.Errorf("\nTestCase: %s\nresponse is nil", test.testName)
		}
	}
}

func TestDriverDeleteBucket(t *testing.T) {
	tests := []struct {
		testName    string
		bucketID    *types.BucketID
		expectedErr error
	}{
		{
			testName: "Delete Storage Account Bucket",
			bucketID: &types.BucketID{
				SubID:         constant.ValidSub,
				ResourceGroup: constant.ValidResourceGroup,
				URL:           constant.ValidAccountURL,
			},
			expectedErr: nil,
		},
		{
			testName: "Delete Container Bucket",
			bucketID: &types.BucketID{
				SubID:         constant.ValidSub,
				ResourceGroup: constant.ValidResourceGroup,
				URL:           constant.ValidContainerURL,
			},
			expectedErr: fmt.Errorf("Error deleting container %s in storage account %s : %v", constant.ValidContainer, constant.ValidAccount, fmt.Errorf("Delete \"https://validaccount.blob.core.windows.net/validcontainer?restype=container\": dial tcp: lookup validaccount.blob.core.windows.net: no such host")),
		},
	}

	ctrl := gomock.NewController(t)
	pr := newFakeProvisioner(ctrl)

	for _, test := range tests {
		data, _ := test.bucketID.Encode()
		resp, err := pr.DriverDeleteBucket(context.Background(), &spec.DriverDeleteBucketRequest{
			BucketId: data,
		})
		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("\nTestCase: %s\nexpected: %v\nactual: %v", test.testName, test.expectedErr, err)
		}
		if err == nil && reflect.DeepEqual(nil, resp) {
			t.Errorf("\nTestCase: %s\nresponse is nil", test.testName)
		}
	}
}

func TestDriverGrantBucketAccess(t *testing.T) {
	tests := []struct {
		testName    string
		url         string
		authType    spec.AuthenticationType
		params      map[string]string
		expectedErr error
	}{
		{
			testName:    "No Parameters",
			expectedErr: status.Error(codes.InvalidArgument, "Parameters missing. Cannot initialize Azure bucket."),
		},
		{
			testName:    "Unknown auth type",
			authType:    spec.AuthenticationType_UnknownAuthenticationType,
			params:      map[string]string{},
			expectedErr: status.Error(codes.InvalidArgument, "AuthenticationType not provided in GrantBucketAccess request."),
		},
		{
			testName:    "IAM Not yet Implemented",
			authType:    spec.AuthenticationType_IAM,
			params:      map[string]string{},
			expectedErr: status.Error(codes.Unimplemented, "AuthenticationType IAM not implemented."),
		},
		{
			testName:    "Key Auth Type",
			authType:    spec.AuthenticationType_Key,
			url:         constant.ValidAccountURL,
			params:      map[string]string{},
			expectedErr: nil,
		},
	}

	ctrl := gomock.NewController(t)
	pr := newFakeProvisioner(ctrl)

	for _, test := range tests {
		bucketID := types.BucketID{
			URL: test.url,
		}
		id, err := bucketID.Encode()
		if err != nil {
			t.Errorf("encoding error: %s", err.Error())
		}

		resp, err := pr.DriverGrantBucketAccess(context.Background(), &spec.DriverGrantBucketAccessRequest{
			BucketId:           id,
			AuthenticationType: test.authType,
			Parameters:         test.params,
		})
		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("\nTestCase: %s\nexpected: %v\nactual: %v", test.testName, test.expectedErr, err)
		}
		if err == nil && reflect.DeepEqual(nil, resp) {
			t.Errorf("\nTestCase: %s\nresponse is nil", test.testName)
		}
	}
}
