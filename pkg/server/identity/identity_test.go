package identityserver

import (
	"context"
	"Azure/azure-cosi-driver/pkg/constant"
	"reflect"
	"testing"

	spec "sigs.k8s.io/container-object-storage-interface-spec"
)

func TestNewIdentityServer(t *testing.T) {
	tests := []struct {
		testName    string
		driverName  string
		expectedErr error
		expectedID  *identity
	}{
		{
			testName:   "Valid I/O",
			driverName: constant.ValidDriver,
			expectedID: &identity{
				driverName: constant.ValidDriver,
			},
			expectedErr: nil,
		},
	}
	for _, test := range tests {
		id, err := NewIdentityServer(test.driverName)
		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("\nTestCase: %s\nexpected: %v\nactual: %v", test.testName, test.expectedErr, err)
		}
		if err == nil && !reflect.DeepEqual(id, test.expectedID) {
			t.Errorf("\nTestCase: %s\nexpected: %v\nactual: %v", test.testName, test.expectedID, id)
		}
	}
}

func TestDriverGetInfo(t *testing.T) {
	tests := []struct {
		testName    string
		id          *identity
		req         *spec.DriverGetInfoRequest
		expectedRes *spec.DriverGetInfoResponse
		expectedErr error
	}{
		{
			testName: "Valid I/O",
			id: &identity{
				driverName: constant.ValidDriver,
			},
			req: &spec.DriverGetInfoRequest{},
			expectedRes: &spec.DriverGetInfoResponse{
				Name: constant.ValidDriver,
			},
			expectedErr: nil,
		},
	}
	for _, test := range tests {
		res, err := test.id.DriverGetInfo(context.Background(), test.req)
		if !reflect.DeepEqual(err, test.expectedErr) {
			t.Errorf("\nTestCase: %s\nexpected: %v\nactual: %v", test.testName, test.expectedErr, err)
		}
		if err == nil && !reflect.DeepEqual(res, test.expectedRes) {
			t.Errorf("\nTestCase: %s\nexpected: %v\nactual: %v", test.testName, test.expectedRes, res)
		}
	}
}
