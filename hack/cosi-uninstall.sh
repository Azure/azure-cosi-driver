#!/bin/bash
while getopts "v:" flag;do
    case "${flag}" in
        v) 
            VERSION=$OPTARG
            echo "Version: $VERSION"
            ;;
        *)
            echo "Unknown argument $OPTARG"
            ;;
    esac
done

if [ -z $VERSION ]; then
    VERSION="remote"
fi

echo -e "\nDeleting Kube Resources"
if [ $VERSION = "local" ] || [ $VERSION = "push" ]; 
then
    kubectl delete -k ./. 
    echo -e "\nDeleting COSI Driver and Sidecar"
    DRIVER_NAME=$(dirname "$(dirname "$(realpath "${BASH_SOURCE[0]}")")")
    cd $DRIVER_NAME || echo "Unable to find directory $DRIVER_NAME"; exit 1
    echo "$(pwd)"
else
    kubectl delete -k github.com/Azure/azure-cosi-driver
fi
echo -e "\n"

if [ $VERSION = "push" ]; then
    echo -e "\nDeleting COSI Image"
    PROPERTY_FILE="./resources/cosi-driver-azure.properties"
    export REGISTRY=`cat "${PROPERTY_FILE}" | grep AZURE_DRIVER_IMAGE_ORG | cut -d'=' -f2`
    export IMAGE_VERSION=`cat "${PROPERTY_FILE}" | grep AZURE_DRIVER_IMAGE_VERSION | cut -d'=' -f2`
    echo "Registry: $REGISTRY"
    echo "Version: $IMAGE_VERSION"
    docker image rm "$REGISTRY/azure-cosi-driver:$IMAGE_VERSION"
fi

echo "Deleting CRD's for COSI and COSI Controller"
kubectl delete -k github.com/kubernetes-sigs/container-object-storage-interface-controller
kubectl delete -k github.com/kubernetes-sigs/container-object-storage-interface-api