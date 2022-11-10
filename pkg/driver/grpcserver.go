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
	"context"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"

	klog "k8s.io/klog/v2"

	spec "sigs.k8s.io/container-object-storage-interface-spec"
)

type COSIServer struct {
	waitGroup       sync.WaitGroup
	endpointProto   string
	endpointAddress string
	server          *grpc.Server
}

func newCOSIServer(
	endpointProto string,
	endpointAddr string,
	identityServer spec.IdentityServer,
	provisionerServer spec.ProvisionerServer) *COSIServer {
	serverOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			func(
				ctx context.Context,
				req interface{},
				info *grpc.UnaryServerInfo,
				handler grpc.UnaryHandler) (interface{}, error) {
				klog.V(2).InfoS("GRPC call", "method", info.FullMethod, "request", req)

				resp, err := handler(ctx, req)
				if err != nil {
					klog.Errorf("GRPC error %v", err)
				} else {
					klog.V(2).InfoS("GRPC response", "method", info.FullMethod, "response", resp)
				}

				return resp, err
			}),
	}

	server := grpc.NewServer(serverOpts...)
	spec.RegisterIdentityServer(server, identityServer)
	spec.RegisterProvisionerServer(server, provisionerServer)

	return &COSIServer{
		endpointProto:   endpointProto,
		endpointAddress: endpointAddr,
		server:          server,
	}
}

func (s *COSIServer) startServer() error {
	if s.endpointProto == "unix" {
		// Removing the existing socket file
		if err := os.Remove(s.endpointAddress); err != nil && !os.IsNotExist(err) {
			klog.Errorf("Unable to delete socket file: %s, Error: %v", s.endpointAddress, err)
			return err
		}
	}

	listener, err := net.Listen(s.endpointProto, s.endpointAddress)
	if err != nil {
		return err
	}

	s.waitGroup.Add(1)

	go func() {
		defer s.waitGroup.Done()

		klog.Infof("Starting GRPC server at %s://%s", s.endpointProto, s.endpointAddress)
		if err := s.server.Serve(listener); err != nil {
			klog.Errorf("Error starting GRPC server : %v", err)
		}

		klog.Info("GRPC Server loop finished")
	}()

	return nil
}

func (s *COSIServer) GracefulStop() {
	klog.Info("Requesting graceful GRPC server stop")
	s.server.GracefulStop()
}

func (s *COSIServer) Stop() {
	klog.Info("Requesting GRPC server stop")
	s.server.Stop()
}

func (s *COSIServer) Wait() {
	s.waitGroup.Wait()
	klog.Info("GRPC Server stopped")
}
