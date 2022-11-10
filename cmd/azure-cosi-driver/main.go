// Copyright 2021 The Ceph-COSI Authors.
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

package main

import (
	"flag"
	"Azure/azure-cosi-driver/pkg/driver"
	identityserver "Azure/azure-cosi-driver/pkg/server/identity"
	provisionerserver "Azure/azure-cosi-driver/pkg/server/provisioner"

	"k8s.io/klog"
)

var (
	endpoint                   = flag.String("endpoint", driver.DefaultEndpoint, "endpoint for the GRPC server")
	kubeconfig                 = flag.String("kubeconfig", "", "Absolute path to the kubeconfig file. Required only when running out of cluster.")
	cloudConfigSecretName      = flag.String("cloud-config-secret-name", "azure-cloud-provider", "cloud config secret name")
	cloudConfigSecretNamespace = flag.String("cloud-config-secret-namespace", "kube-system", "cloud config secret namespace")
)

func init() {
	klog.InitFlags(nil)
}

func main() {
	flag.Parse()
	defer klog.Flush()

	provServer, err := provisionerserver.NewProvisionerServer(*kubeconfig, *cloudConfigSecretName, *cloudConfigSecretNamespace)
	if err != nil {
		klog.Exitf("Error creating ProvisionerServer: %v", err)
	}
	identityServer, err := identityserver.NewIdentityServer(driver.DriverName)
	if err != nil {
		klog.Exitf("Error creating IdentityServer: %v", err)
	}

	err = driver.RunServerWithSignalHandler(*endpoint, identityServer, provServer)
	if err != nil {
		klog.Exitf("Error when running driver: %v", err)
	}
}
