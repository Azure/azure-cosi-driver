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
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	clientSet "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	azure "sigs.k8s.io/cloud-provider-azure/pkg/provider"
)

const (
	// DefaultAzureCredentialFileEnv is the default azure credentials file env variable
	DefaultAzureCredentialFileEnv = "AZURE_CREDENTIAL_FILE"
	// DefaultCredFilePathLinux is default creds file for linux machine
	DefaultCredFilePathLinux = "/etc/kubernetes/azure.json"
	// DefaultCredFilePathWindows is default creds file for windows machine
	DefaultCredFilePathWindows = "C:\\k\\azure.json"
)

// GetKubeConfig gets config object from config file
func getKubeConfig(kubeconfig string) (config *rest.Config, err error) {

	if kubeconfig == "" {
		// if kubeconfig path is not passed
		// read the incluster config
		config, err = rest.InClusterConfig()

		// if we couldn't get the in-cluster config
		// get kubeconfig path from environment variable
		if err != nil {
			kubeconfig = os.Getenv("KUBECONFIG")
			if kubeconfig == "" {
				kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
			}
		} else {
			return config, err
		}
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)

	return config, err
}

func GetKubeClient(kubeconfig string) (*clientSet.Clientset, error) {
	config, err := getKubeConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	return clientSet.NewForConfig(config)
}

// GetAzureCloudProvider get Azure Cloud Provider
func GetAzureCloudProvider(
	kubeClient clientSet.Interface,
	secretName,
	secretNamespace string) (*azure.Cloud, error) {
	az := &azure.Cloud{
		InitSecretConfig: azure.InitSecretConfig{
			SecretName:      secretName,
			SecretNamespace: secretNamespace,
			CloudConfigKey:  "cloud-config",
		},
	}
	if kubeClient != nil {
		klog.Infof("reading cloud config from secret")
		az.KubeClient = kubeClient
		if err := az.InitializeCloudFromSecret(); err != nil {
			klog.Infof("InitializeCloudFromSecret failed with error: %v", err)
		}
	}

	if az.TenantID == "" || az.SubscriptionID == "" || az.ResourceGroup == "" {
		klog.Infof("could not read cloud config from secret")
		credFile, ok := os.LookupEnv(DefaultAzureCredentialFileEnv)
		if ok && strings.TrimSpace(credFile) != "" {
			klog.Infof("%s env var set as %v", DefaultAzureCredentialFileEnv, credFile)
		} else {
			if runtime.GOOS == "windows" {
				credFile = DefaultCredFilePathWindows
			} else {
				credFile = DefaultCredFilePathLinux
			}

			klog.Infof("use default %s env var: %v", DefaultAzureCredentialFileEnv, credFile)
		}

		f, err := os.Open(credFile)
		if err != nil {
			klog.Errorf("Failed to load config from file: %s, Error: %+v", credFile, err)
			return nil, fmt.Errorf("Failed to load config from file: %s, cloud not get azure cloud provider", credFile)
		}
		defer f.Close()

		klog.Infof("read cloud config from file: %s successfully", credFile)
		if az, err = azure.NewCloudWithoutFeatureGates(f, false); err != nil {
			return az, err
		}
	}

	// reassign kubeClient
	if kubeClient != nil && az.KubeClient == nil {
		az.KubeClient = kubeClient
	}
	return az, nil
}
