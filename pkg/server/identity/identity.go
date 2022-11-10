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

package identityserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	spec "sigs.k8s.io/container-object-storage-interface-spec"
)

type identity struct {
	spec.UnimplementedIdentityServer
	driverName string
}

// type-check
var _ spec.IdentityServer = &identity{}

func NewIdentityServer(driverName string) (spec.IdentityServer, error) {
	return &identity{
		driverName: driverName,
	}, nil
}

func (id identity) DriverGetInfo(
	ctx context.Context,
	req *spec.DriverGetInfoRequest) (*spec.DriverGetInfoResponse, error) {
	if id.driverName == "" {
		return nil, status.Error(codes.InvalidArgument, "Driver name not found")
	}

	return &spec.DriverGetInfoResponse{
		Name: id.driverName,
	}, nil
}
