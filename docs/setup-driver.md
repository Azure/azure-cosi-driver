## Remotely Install with Required CRDs and Controller with Shell Script

### Install the driver

 - Create New Cluster and Install
 
```console
curl -skSL https://raw.githubusercontent.com/Azure/azure-cosi-driver/master/hack/{azure-cluster-up.sh,cosi-install.sh} | bash -s -- -n <cluster-name> -r <resourge-group> -s <subscription-id>
```
 - Install Driver without creating new Cluster
 
```console
curl -skSL https://raw.githubusercontent.com/Azure/azure-cosi-driver/master/hack/cosi-install.sh | bash
```

 - Uninstall Driver and Delete Cluster
 
```console
 curl -skSL https://raw.githubusercontent.com/Azure/azure-cosi-driver/master/hack/{cosi-uninstall.sh,azure-cluster-down.sh} | bash -s -- -n <cluster-name> -r <resourge-group>
```

 - Uninstall Driver
 
```console
 curl -skSL https://raw.githubusercontent.com/Azure/azure-cosi-driver/master/hack/cosi-uninstall.sh | bash
```

## Locally Install with Required CRDs and Controller with Shell Script

### Clone the COSI repo

```console
git@github.com:Azure/azure-cosi-driver.git
```

### Install the driver

 - Create New Cluster and Install
 
 ```console
 ./hack/cosi-azure-cluster-up.sh -n <cluster_name> -r <resource_group> -s <subscription_id>
 ```

 - Install Driver without creating new Cluster

 ```console
 ./hack/cosi-install.sh
 ```