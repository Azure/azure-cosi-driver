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
	"encoding/base64"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"Azure/azure-cosi-driver/pkg/constant"
	"Azure/azure-cosi-driver/pkg/types"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sigs.k8s.io/cloud-provider-azure/pkg/azureclients/storageaccountclient/mockstorageaccountclient"
	azure "sigs.k8s.io/cloud-provider-azure/pkg/provider"
	"sigs.k8s.io/cloud-provider-azure/pkg/retry"
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

func TestDeleteStorageAccount(t *testing.T) {
	tests := []struct {
		testName    string
		id          *types.BucketID
		expectedErr error
	}{
		{
			testName: "Valid Account",
			id: &types.BucketID{
				SubID:         constant.ValidSub,
				ResourceGroup: constant.ValidResourceGroup,
				URL:           constant.ValidAccountURL,
			},
			expectedErr: nil,
		},
		{
			testName: "Invalid Account",
			id: &types.BucketID{
				SubID:         constant.ValidSub,
				ResourceGroup: constant.ValidResourceGroup,
				URL:           constant.InvalidAccount,
			},
			expectedErr: retry.GetError(&http.Response{}, status.Error(codes.NotFound, "could not find storage account")).Error(),
		},
	}

	ctrl := gomock.NewController(t)
	cloud := azure.GetTestCloud(ctrl)
	keyList := make([]storage.AccountKey, 0)
	cloud.StorageAccountClient = NewMockSAClient(context.Background(), ctrl, "", "", "", &keyList)

	for _, test := range tests {
		err := DeleteStorageAccount(context.Background(), test.id, cloud)
		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("\nTestCase: %s\nExpected: %v\nActual: %v", test.testName, test.expectedErr, err)
		}
	}
}

func TestCreateStorageAccountBucket(t *testing.T) {
	tests := []struct {
		testName    string
		account     string
		accountURL  string
		expectedErr error
	}{
		{
			testName:    "Valid Account",
			account:     constant.ValidAccount,
			accountURL:  constant.ValidAccountURL,
			expectedErr: nil,
		},
		{
			testName: "Invalid Account",
			account:  constant.InvalidAccount,
			expectedErr: status.Error(codes.Internal, fmt.Sprintf("Could not create storage account: %v",
				fmt.Errorf("could not get storage key for storage account "+constant.InvalidAccount+": "+retry.GetError(&http.Response{}, fmt.Errorf("Invalid Account")).Error().Error()))),
		},
	}
	ctrl := gomock.NewController(t)
	cloud := azure.GetTestCloud(ctrl)
	keyList := make([]storage.AccountKey, 0)
	keyList = append(keyList, storage.AccountKey{KeyName: to.StringPtr(constant.ValidAccount), Value: to.StringPtr("val")})
	cloud.StorageAccountClient = NewMockSAClient(context.Background(), ctrl, "", "", "", &keyList)

	for _, test := range tests {
		base64ID, err := createStorageAccountBucket(context.Background(), test.account, &BucketClassParameters{storageAccountName: test.account}, cloud)
		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("\nTestCase: %s\nexpected: %v\nactual: %v", test.testName, test.expectedErr, err)
		}

		id, _ := types.DecodeToBucketID(base64ID)
		accName := ""
		if id != nil {
			accName = id.URL
		}
		if err == nil && !reflect.DeepEqual(accName, test.accountURL) {
			t.Errorf("\nTestCase: %s\nexpected account: %s\nactual account: %s", test.testName, test.account, accName)
		}
	}
}

func TestCreateAccountSASURL(t *testing.T) {
	tests := []struct {
		testName    string
		bucketID    string
		params      *BucketAccessClassParameters
		key         string
		urlIsEmpty  bool
		expectedID  string
		expectedErr error
	}{
		{
			testName:    "Key is illegal base64 data",
			bucketID:    constant.ValidAccountURL,
			params:      &BucketAccessClassParameters{},
			key:         "badkey",
			expectedID:  constant.ValidAccountURL,
			expectedErr: fmt.Errorf("decode account key: %w", base64.CorruptInputError(4)),
		},
		{
			testName:    "Missing permissions, expiry, etc.",
			bucketID:    constant.ValidAccountURL,
			params:      &BucketAccessClassParameters{},
			key:         "",
			expectedID:  constant.ValidAccountURL,
			expectedErr: fmt.Errorf("account SAS is missing at least one of these: ExpiryTime, Permissions, Service, or ResourceType"),
		},
		{
			testName: "Correct Inputs",
			bucketID: constant.ValidAccountURL,
			params: &BucketAccessClassParameters{
				allowServiceSignedResourceType:   true,
				allowContainerSignedResourceType: true,
				allowObjectSignedResourceType:    true,
				enableRead:                       true,
				enableList:                       true,
				validationPeriod:                 1,
			},
			key:         "",
			expectedID:  constant.ValidAccountURL,
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		_, bucketID, err := createAccountSASURL(context.Background(), test.bucketID, test.params, test.key)
		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("\nTestCase: %s\nexpected:\t%v\nactual: \t%v", test.testName, test.expectedErr, err)
		}
		if err == nil && !reflect.DeepEqual(bucketID, test.expectedID) {
			t.Errorf("\nTestCase: %s\nexpected account: %s\nactual account: %s", test.testName, test.expectedID, bucketID)
		}
	}
}
