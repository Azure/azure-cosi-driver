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

package driver

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
	spec "sigs.k8s.io/container-object-storage-interface-spec"
)

const (
	DefaultEndpoint = "unix:///var/lib/cosi/cosi.sock"
	DriverName      = "blob.cosi.azure.com"
	waitTime        = 5
)

// Run COSI gRPC Services.
// Set up signal handlers for SIGINT and SIGTERM.
func RunServerWithSignalHandler(
	endpoint string,
	identityServer spec.IdentityServer,
	provisionerServer spec.ProvisionerServer) error {
	server, err := StartServers(endpoint, identityServer, provisionerServer)
	if err != nil {
		return nil
	}

	// Registering signal handlers
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func(c chan os.Signal) {
		sig := <-c
		signal.Stop(sigChan)

		klog.InfoS("Received signal", "signal", sig)

		// Initially try to gracefully stop the server
		server.GracefulStop()

		time.Sleep(waitTime * time.Second)
		// Stopping server
		server.Stop()

		time.Sleep(waitTime * time.Second)
		klog.Info("Signal handler timed out, forcing process exit")
		klog.Flush()
		os.Exit(1)
	}(sigChan)

	// Wait till all the WaitGroup are released
	server.Wait()

	return nil
}

// Run the GRPC services
func StartServers(
	endpoint string,
	identityServer spec.IdentityServer,
	provServer spec.ProvisionerServer) (*COSIServer, error) {
	proto, addr, err := parseEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	grpcServer := newCOSIServer(proto, addr, identityServer, provServer)
	if err := grpcServer.startServer(); err != nil {
		klog.Errorf("Error starting GRPC server %v", err)
		return nil, err
	}

	return grpcServer, nil
}

// Parse the endpoint to extract the protocol and address
func parseEndpoint(endpoint string) (string, string, error) {
	if !strings.HasPrefix(endpoint, "unix://") &&
		!strings.HasPrefix(endpoint, "tcp://") {
		return "", "", status.Error(codes.InvalidArgument, fmt.Sprintf("Unable to parse endpoint protocol %s", endpoint))
	}

	pieces := strings.SplitN(endpoint, "://", 2)
	if pieces[1] == "" {
		return "", "", status.Error(codes.InvalidArgument, fmt.Sprintf("Endpoint %s missing address", endpoint))
	}

	return pieces[0], pieces[1], nil
}
