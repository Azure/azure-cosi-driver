package azureutils

import (
	"os"
	"reflect"
	"testing"
)

func createTestFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func TestGetAzureCloudProvider(t *testing.T) {
	emptyKubeConfig := "empty-kube-config"
	err := createTestFile(emptyKubeConfig)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		if err := os.Remove(emptyKubeConfig); err != nil {
			t.Error(err)
		}
	}()
	tests := []struct {
		testName    string
		kubeconfig  string
		expectedErr error
	}{
		/*{
			testName:    "no kubeconfig, no credential file",
			kubeconfig:  "",
			expectedErr: nil,
		},*/
	}
	for _, test := range tests {
		kubeclient, err := GetKubeClient(test.kubeconfig)
		if err != nil {
			t.Errorf("Testcase: %s\nError getting kubeclient from kubeconfig %s", test.testName, test.kubeconfig)
		}

		cloud, err := GetAzureCloudProvider(kubeclient, "", "")
		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("\nTestCase: %s\nExpected Error: %v\nActual Error: %v", test.testName, test.expectedErr, err)
		}
		if cloud == nil {
			t.Errorf("return value of get cloud provider should not be nil even if there is an error")
		} /* else if !reflect.DeepEqual(cloud.Environment.StorageEndpointSuffix, constant.CloudDefaultURL) {
			t.Errorf("\nTestCase: %s\nExpected Output: %v\nActual Output: %v", test.testName, constant.CloudDefaultURL, cloud.Environment.StorageEndpointSuffix)
		}*/
	}
}
