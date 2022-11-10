package azureutils

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"Azure/azure-cosi-driver/pkg/constant"
	"Azure/azure-cosi-driver/pkg/types"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	azure "sigs.k8s.io/cloud-provider-azure/pkg/provider"
	"sigs.k8s.io/cloud-provider-azure/pkg/retry"
)

func TestCreateBucket(t *testing.T) {
	tests := []struct {
		testName    string
		bucket      string
		expectedURL string
		expectedErr error
		params      map[string]string
	}{
		{
			testName: "Parsing Error (invalid bucket unit type)",
			bucket:   constant.ValidContainerURL,
			params:   map[string]string{constant.BucketUnitTypeField: "invalid type"},
			expectedErr: status.Error(codes.Unknown, fmt.Sprintf("Error parsing parameters : %v",
				status.Error(codes.InvalidArgument, fmt.Sprintf("Invalid BucketUnitType %s", "invalid type")))),
		},
		{
			testName: "Create storage Account Bucket",
			bucket:   constant.ValidAccount,
			params: map[string]string{
				constant.BucketUnitTypeField:     constant.StorageAccount.String(),
				constant.StorageAccountNameField: constant.ValidAccount},
			expectedURL: constant.ValidAccountURL,
			expectedErr: nil,
		},
		{
			testName: "Create Container Bucket (invalid account)",
			bucket:   constant.ValidContainer,
			params: map[string]string{
				constant.BucketUnitTypeField:     constant.Container.String(),
				constant.StorageAccountNameField: constant.InvalidAccount},
			expectedURL: constant.ValidContainerURL,
			expectedErr: status.Error(codes.Internal, fmt.Sprintf("Could not ensure storage account %s exists: %v", constant.InvalidAccount,
				fmt.Errorf("could not get storage key for storage account "+constant.InvalidAccount+": "+retry.GetError(&http.Response{}, fmt.Errorf("Invalid Account")).Error().Error()))),
		},
	}
	ctrl := gomock.NewController(t)
	cloud := azure.GetTestCloud(ctrl)
	keyList := make([]storage.AccountKey, 0)
	keyList = append(keyList, storage.AccountKey{KeyName: to.StringPtr(constant.ValidAccount), Value: to.StringPtr("val")})
	cloud.StorageAccountClient = NewMockSAClient(context.Background(), ctrl, "", "", "", &keyList)

	for _, test := range tests {
		base64ID, err := CreateBucket(context.Background(), constant.ValidAccount, test.params, cloud)

		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("\nTestCase: %s\nExpected Error: %v\nActual Error: %v", test.testName, test.expectedErr, err)
		}

		id, _ := types.DecodeToBucketID(base64ID)
		url := ""
		if id != nil {
			url = id.URL
		}
		if err == nil && !reflect.DeepEqual(url, test.expectedURL) {
			t.Errorf("\nTestCase: %s\nExpected URL: %v\nActual URL: %v", test.testName, test.expectedURL, url)
		}
	}
}

func TestDeleteBucket(t *testing.T) {
	tests := []struct {
		testName    string
		id          *types.BucketID
		expectedErr error
	}{
		{
			testName: "Individual Blob Unit Type Unsupported",
			id: &types.BucketID{
				SubID:         constant.ValidSub,
				ResourceGroup: constant.ValidResourceGroup,
				URL:           constant.ValidBlobURL,
			},
			expectedErr: status.Error(codes.InvalidArgument, "Individual Blobs unsupported. Please use Blob Containers or Storage Accounts instead."),
		},
		{
			testName: "Delete storage Account Bucket",
			id: &types.BucketID{
				SubID:         constant.ValidSub,
				ResourceGroup: constant.ValidResourceGroup,
				URL:           constant.ValidAccountURL,
			},
			expectedErr: nil,
		},
		{
			testName: "Delete Container Bucket (Invalid credentials)",
			id: &types.BucketID{
				SubID:         constant.ValidSub,
				ResourceGroup: constant.ValidResourceGroup,
				URL:           constant.ValidContainerURL,
			},
			expectedErr: fmt.Errorf("Error deleting container %s in storage account %s : %v", constant.ValidContainer, constant.ValidAccount, fmt.Errorf("Delete \"https://validaccount.blob.core.windows.net/validcontainer?restype=container\": dial tcp: lookup validaccount.blob.core.windows.net: no such host")),
		},
	}
	ctrl := gomock.NewController(t)
	cloud := azure.GetTestCloud(ctrl)
	keyList := make([]storage.AccountKey, 0)
	keyList = append(keyList, storage.AccountKey{KeyName: to.StringPtr(constant.ValidAccount), Value: to.StringPtr(base64.StdEncoding.EncodeToString([]byte{1, 2, 3, 4}))})
	cloud.StorageAccountClient = NewMockSAClient(context.Background(), ctrl, "", "", "", &keyList)

	for _, test := range tests {
		base64ID, err := test.id.Encode()
		if err != nil {
			t.Errorf("encoding error: %s", err.Error())
		}

		err = DeleteBucket(context.Background(), base64ID, cloud)
		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("\nTestCase: %s\nExpected Error: %v\nActual Error: %v", test.testName, test.expectedErr, err)
		}
	}
}

func TestParseBucketClassParameters(t *testing.T) {
	tests := []struct {
		testName       string
		parameters     map[string]string
		expectedErr    error
		expectedParams BucketClassParameters
	}{
		{
			testName:       "BucketUnitType Container",
			parameters:     map[string]string{constant.BucketUnitTypeField: constant.Container.String()},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{bucketUnitType: constant.Container},
		},
		{
			testName:       "BucketUnitType Account",
			parameters:     map[string]string{constant.BucketUnitTypeField: constant.StorageAccount.String()},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{bucketUnitType: constant.StorageAccount, createStorageAccount: to.BoolPtr(true)},
		},
		{
			testName:       "Create Bucket True",
			parameters:     map[string]string{constant.CreateBucketField: TrueValue},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{createBucket: true},
		},
		{
			testName:       "Create StorageAccountField True",
			parameters:     map[string]string{constant.BucketUnitTypeField: constant.Container.String()},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{bucketUnitType: constant.Container},
		},
		{
			testName:       "Valid SubscriptionID",
			parameters:     map[string]string{constant.SubscriptionIDField: constant.ValidSub},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{subscriptionID: constant.ValidSub},
		},
		{
			testName:       "Valid Storage Account",
			parameters:     map[string]string{constant.StorageAccountNameField: constant.ValidAccount},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{storageAccountName: constant.ValidAccount},
		},
		{
			testName:       "Valid Region",
			parameters:     map[string]string{constant.RegionField: constant.ValidRegion},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{region: constant.ValidRegion},
		},
		{
			testName:       "Access Tier Field Hot",
			parameters:     map[string]string{constant.AccessTierField: constant.Hot.String()},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{accessTier: constant.Hot},
		},
		{
			testName:       "Access Tier Field Cool",
			parameters:     map[string]string{constant.AccessTierField: constant.Cool.String()},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{accessTier: constant.Cool},
		},
		{
			testName:       "Access Tier Field Archive",
			parameters:     map[string]string{constant.AccessTierField: constant.Archive.String()},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{accessTier: constant.Archive},
		},
		{
			testName:       "SKU Name Field StandardLRS",
			parameters:     map[string]string{constant.SKUNameField: constant.StandardLRS.String()},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{SKUName: constant.StandardLRS},
		},
		{
			testName:       "SKU Name Field StandardGRS",
			parameters:     map[string]string{constant.SKUNameField: constant.StandardGRS.String()},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{SKUName: constant.StandardGRS},
		},
		{
			testName:       "SKU Name Field StandardRAGRS",
			parameters:     map[string]string{constant.SKUNameField: constant.StandardRAGRS.String()},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{SKUName: constant.StandardRAGRS},
		},
		{
			testName:       "SKU Name Field PremiumLRS",
			parameters:     map[string]string{constant.SKUNameField: constant.PremiumLRS.String()},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{SKUName: constant.PremiumLRS},
		},
		{
			testName:       "SKU Name Field unsupported",
			parameters:     map[string]string{constant.SKUNameField: "foobar"},
			expectedErr:    status.Error(codes.InvalidArgument, fmt.Sprintf("Access Tier %s is unsupported", "foobar")),
			expectedParams: BucketClassParameters{},
		},
		{
			testName:       "Valid Resource Group",
			parameters:     map[string]string{constant.ResourceGroupField: constant.ValidResourceGroup},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{resourceGroup: constant.ValidResourceGroup},
		},
		{
			testName:       "AllowBlobAccess True",
			parameters:     map[string]string{constant.AllowBlobAccessField: TrueValue},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{allowBlobAccess: true},
		},
		{
			testName:       "SharedAccessKey True",
			parameters:     map[string]string{constant.AllowSharedAccessKeyField: TrueValue},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{allowSharedAccessKey: true},
		},
		{
			testName:       "BlobVersioning True",
			parameters:     map[string]string{constant.EnableBlobVersioningField: TrueValue},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{enableBlobVersioning: true},
		},
		{
			testName:       "EnableBlobDeleteRetention True",
			parameters:     map[string]string{constant.EnableBlobDeleteRetentionField: TrueValue},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{enableBlobDeleteRetention: true},
		},
		{
			testName:       "BlobRetentionDays 1",
			parameters:     map[string]string{constant.BlobDeleteRetentionDaysField: "1"},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{blobDeleteRetentionDays: 1},
		},
		{
			testName:       "BlobRetentionDays Not a number",
			parameters:     map[string]string{constant.BlobDeleteRetentionDaysField: "foobar"},
			expectedErr:    status.Error(codes.InvalidArgument, "strconv.Atoi: parsing \"foobar\": invalid syntax"),
			expectedParams: BucketClassParameters{},
		},
		{
			testName:       "EnableContainerDeleteRetention True",
			parameters:     map[string]string{constant.EnableContainerDeleteRetentionField: TrueValue},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{enableContainerDeleteRetention: true},
		},
		{
			testName:       "ContainerRetentionDays 1",
			parameters:     map[string]string{constant.ContainerDeleteRetentionDaysField: "1"},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{containerDeleteRetentionDays: 1},
		},
		{
			testName:       "ContainerRetentionDays 1",
			parameters:     map[string]string{constant.ContainerDeleteRetentionDaysField: "foobar"},
			expectedErr:    status.Error(codes.InvalidArgument, "strconv.Atoi: parsing \"foobar\": invalid syntax"),
			expectedParams: BucketClassParameters{},
		},
		{
			testName:       "storage account type",
			parameters:     map[string]string{StorageAccountTypeField: "unittest"},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{storageAccountType: "unittest"},
		},
		{
			testName:       "account kind StorageV2",
			parameters:     map[string]string{KindField: constant.StorageV2.String()},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{kind: constant.StorageV2},
		},
		{
			testName:       "account kind Storage",
			parameters:     map[string]string{KindField: constant.Storage.String()},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{kind: constant.Storage},
		},
		{
			testName:       "account kind BlobStorage",
			parameters:     map[string]string{KindField: constant.BlobStorage.String()},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{kind: constant.BlobStorage},
		},
		{
			testName:       "account kind BlockBlobStorage",
			parameters:     map[string]string{KindField: constant.BlockBlobStorage.String()},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{kind: constant.BlockBlobStorage},
		},
		{
			testName:       "account kind FileStorage",
			parameters:     map[string]string{KindField: constant.FileStorage.String()},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{kind: constant.FileStorage},
		},
		{
			testName:       "tags",
			parameters:     map[string]string{TagsField: "key1=value1,key2=value2"},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{tags: map[string]string{"key1": "value1", "key2": "value2"}},
		},
		{
			testName:       "VN Resource ID's",
			parameters:     map[string]string{VNResourceIdsField: "foo,bar"},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{virtualNetworkResourceIDs: []string{"foo", "bar"}},
		},
		{
			testName:       "Create Private Endpoint",
			parameters:     map[string]string{CreatePrivateEndpointField: TrueValue},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{createPrivateEndpoint: true},
		},
		{
			testName:       "HNS Enabled",
			parameters:     map[string]string{HNSEnabledField: TrueValue},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{isHnsEnabled: true},
		},
		{
			testName:       "enableNfsV3",
			parameters:     map[string]string{EnableNFSV3Field: TrueValue},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{enableNfsV3: true},
		},
		{
			testName:       "enableLargeFileShare",
			parameters:     map[string]string{EnableLargeFileSharesField: TrueValue},
			expectedErr:    nil,
			expectedParams: BucketClassParameters{enableLargeFileShare: true},
		},
	}
	for _, test := range tests {
		params, err := parseBucketClassParameters(test.parameters)
		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("\nTestCase: %s\nExpected Error: %v\nActual Error: %v", test.testName, test.expectedErr, err)
		}

		if params != nil && test.expectedParams.createStorageAccount != nil && params.createStorageAccount != nil {
			if to.Bool(test.expectedParams.createStorageAccount) == to.Bool(params.createStorageAccount) {
				test.expectedParams.createStorageAccount = params.createStorageAccount
			}
		}

		if err == nil && !reflect.DeepEqual(*params, test.expectedParams) {
			t.Errorf("\nTestCase: %s\nExpected Params: %+v\nActual Params: %+v", test.testName, test.expectedParams, params)
		}
	}
}

func TestParseBucketAccessClassParameters(t *testing.T) {

}

func TestGetAccountOptions(t *testing.T) {
	t.Run("All Variables Filled", func(t *testing.T) {
		input := &BucketClassParameters{
			storageAccountName:        constant.ValidAccount,
			resourceGroup:             constant.ValidResourceGroup,
			region:                    constant.ValidRegion,
			storageAccountType:        constant.ValidAccountType,
			kind:                      constant.StorageV2,
			tags:                      map[string]string{"foo": "bar"},
			virtualNetworkResourceIDs: []string{"id1"},
			enableHTTPSTrafficOnly:    true,
			createPrivateEndpoint:     true,
			isHnsEnabled:              true,
			enableNfsV3:               true,
			enableLargeFileShare:      true,
		}
		expectedOutput := azure.AccountOptions{
			Name:                      constant.ValidAccount,
			ResourceGroup:             constant.ValidResourceGroup,
			Location:                  constant.ValidRegion,
			Type:                      constant.ValidAccountType,
			Kind:                      constant.StorageV2.String(),
			Tags:                      map[string]string{"foo": "bar"},
			VirtualNetworkResourceIDs: []string{"id1"},
			EnableHTTPSTrafficOnly:    true,
			CreatePrivateEndpoint:     true,
			IsHnsEnabled:              to.BoolPtr(true),
			EnableNfsV3:               to.BoolPtr(true),
			EnableLargeFileShare:      true,
		}
		output := getAccountOptions(input)
		if !reflect.DeepEqual(*output, expectedOutput) {
			t.Errorf("\nExpected Options: %+v\nActual Options: %+v", expectedOutput, output)
		}
	})
}
