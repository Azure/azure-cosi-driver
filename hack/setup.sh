#!/bin/bash
while getopts "n:r:l:s:v:" flag;do
    case "${flag}" in
        n) 
            CLUSTER_NAME=$OPTARG
            echo "Cluster Name: $CLUSTER_NAME"
            ;;
        r) 
            RESOURCE_GROUP=$OPTARG
            echo "Resource Group: $RESOURCE_GROUP"
            ;;
        l) 
            LOCATION=$OPTARG
            echo "Location: $LOCATION"
            ;;
        s) 
            SUBSCRIPTION_ID=$OPTARG
            echo "Subscription ID: $SUBSCRIPTION_ID"
            ;;
        v) 
            VERSION=$OPTARG
            echo "version: $VERSION"
            ;;
        *)
            echo "Unknown argument $OPTARG"
            ;;
    esac
done

#check mandatory flags
if [ -z $CLUSTER_NAME ]; then
    echo "Error: Missing argument Cluster Name (flag -n)"
    exit 1
fi
if [ -z $RESOURCE_GROUP ]; then
    echo "Error: Missing argument Resource Group Name (flag -r)"
    exit 1
fi
if [ -z $SUBSCRIPTION_ID ]; then
    echo "Subscription ID (flag -s) not given, getting subID from current context"
    SUBSCRIPTION_ID=$(az account show --query id --output tsv)
fi
if [ -z $VERSION ]; then
    VERSION="local"
fi

DRIVER_NAME=$(dirname "$(realpath ${BASH_SOURCE[0]})")

echo -e "\nRunning azure-cluster-up.sh"
if [ -z $LOCATION ];
then
    source "$DRIVER_NAME/azure-cluster-up.sh" -n $CLUSTER_NAME -r $RESOURCE_GROUP -s $SUBSCRIPTION_ID -v $VERSION
else
    source "$DRIVER_NAME/azure-cluster-up.sh" -n $CLUSTER_NAME -r $RESOURCE_GROUP -l $LOCATION -s $SUBSCRIPTION_ID -v $VERSION
fi

echo -e "\nrunning cosi-install.sh"
source "$DRIVER_NAME/cosi-install.sh" -v $VERSION
