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

echo "Getting CRD's for COSI and COSI Controller"
kubectl create -k github.com/kubernetes-sigs/container-object-storage-interface-api
kubectl create -k github.com/kubernetes-sigs/container-object-storage-interface-controller

echo -e "\nInstalling COSI Driver and Sidecar"

if [ $VERSION = "local" ] || [ $VERSION = "push" ];  then
    DRIVER_NAME="$(dirname "$(dirname "$(realpath ${BASH_SOURCE[0]})")")"
    cd $DRIVER_NAME || exit
    echo "$(pwd)"
fi

if [ $VERSION = "push" ]; then
    echo -e "\nVerifying Binaries"
    make all

    echo -e "\nBuilding COSI Image"
    IMAGE_ID=$(docker build -q .)
    echo "Image ID: $IMAGE_ID"
    PROPERTY_FILE="./resources/cosi-driver-azure.properties"
    export REGISTRY=`cat $PROPERTY_FILE | grep AZURE_DRIVER_IMAGE_ORG | cut -d'=' -f2`
    export IMAGE_VERSION=`cat $PROPERTY_FILE | grep AZURE_DRIVER_IMAGE_VERSION | cut -d'=' -f2`
    echo "Registry: $REGISTRY"
    echo "Version: $IMAGE_VERSION"
    docker tag $IMAGE_ID "$REGISTRY/azure-cosi-driver:$IMAGE_VERSION"

    echo -e "\nPushing COSI Image $REGISTRY/azure-cosi-driver:$IMAGE_VERSION"
    docker push "$REGISTRY/azure-cosi-driver:$IMAGE_VERSION"
fi

echo -e "\nCreating Kube Resources"
if [ $VERSION = "local" ] || [ $VERSION = "push" ]; 
then
    kubectl create -k ./. 
else
    kubectl create -k github.com/Azure/azure-cosi-driver
fi
echo -e "\n"
