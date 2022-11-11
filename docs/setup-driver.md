## Clone the COSI repo

```console
git@github.com:Azure/azure-cosi-driver.git
```

## Install with Required CRDs and Controller with Shell Script

 - Create New Cluster and Install
 
 ```console
 ./hack/cosi-azure-cluster-up.sh -n <cluster_name> -r <resource_group> -s <subscription_id>
 ```

 - Install Driver without creating new Cluster

 ```console
 ./hack/cosi-install.sh
 ```

 ## Build the Cosi Driver Image

Build the Image
 ```console
 make all
 docker build .
 ```

Setup values for REGISTRY and IMAGE_VERSION
```console
export REGISTRY=<docker_id> #set the same value to AZURE_DRIVER_IMAGE_ORG in resources/ cosi-driver-azure.properties file in the repo
export IMAGE_VERSION=<version> #Defaults to latest. Set the same value to AZURE_DRIVER_IMAGE_VERSION in resources/cosi-driver-azure.properties
```

Tag and push image to repository
```console
docker tag <image_id> $REGISTRY/azure-cosi-driver:$IMAGE_VERSION
docker push $REGISTRY/azure-cosi-driver:$IMAGE_VERSION
```

Run kustomization file
```console
kubectl create -k ./.
```
