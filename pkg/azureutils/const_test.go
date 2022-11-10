package azureutils

import (
	"fmt"
	"Azure/azure-cosi-driver/pkg/constant"
	"Azure/azure-cosi-driver/pkg/types"
	"reflect"
	"testing"

	"github.com/Azure/go-autorest/autorest/to"
)

func TestEncode(t *testing.T) {
	id := &types.BucketID{
		SubID:         "subid",
		ResourceGroup: constant.ValidResourceGroup,
		URL:           constant.ValidContainerURL,
	}
	var id2 *types.BucketID
	data, err := id.Encode()
	if err != nil {
		t.Errorf(err.Error())
	}
	id2, err = types.DecodeToBucketID(data)
	if err != nil {
		t.Errorf(err.Error())
	}
	if id2.URL != id.URL {
		t.Errorf("\nExpected Output: %v\nActual Output: %v", id.URL, id2.URL)
	}
}

func TestConvertTagsToMap(t *testing.T) {
	tests := []struct {
		testName    string
		input       string
		expectedErr error
		expectedOut map[string]string
	}{
		{
			testName:    "correct input/output",
			input:       "key1=value1,key2=value2",
			expectedErr: nil,
			expectedOut: map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			testName:    "too many values in key",
			input:       "key1=value1=value3,key2=value2",
			expectedErr: fmt.Errorf("Tags '%s' are invalid, the format should like: 'key1=value1,key2=value2'", "key1=value1=value3,key2=value2"),
			expectedOut: map[string]string{},
		},
		{
			testName:    "missing key",
			input:       "value1,key2=value2",
			expectedErr: fmt.Errorf("Tags '%s' are invalid, the format should like: 'key1=value1,key2=value2'", "value1,key2=value2"),
			expectedOut: map[string]string{},
		},
	}
	for _, test := range tests {
		output, err := ConvertTagsToMap(test.input)
		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("\nTestCase: %s\nExpected Error: %v\nActual Error: %v", test.testName, test.expectedErr, err)
		}
		if err == nil && !reflect.DeepEqual(output, test.expectedOut) {
			t.Errorf("\nTestCase: %s\nExpected Output: %v\nActual Output: %v", test.testName, test.expectedOut, output)
		}
	}
}
func TestConvertMapToMapPointer(t *testing.T) {
	tests := []struct {
		testName    string
		input       map[string]string
		expectedOut map[string]*string
	}{
		{
			testName:    "correct input/output",
			input:       map[string]string{"key1": "value1", "key2": "value2"},
			expectedOut: map[string]*string{"key1": to.StringPtr("value1"), "key2": to.StringPtr("value2")},
		},
	}
	for _, test := range tests {
		output := ConvertMapToMapPointer(test.input)
		if !reflect.DeepEqual(output, test.expectedOut) {
			t.Errorf("\nTestCase: %s\nExpected Output: %v\nActual Output: %v", test.testName, test.expectedOut, output)
		}
	}
}
