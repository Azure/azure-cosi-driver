while getopts "n:r:l:s:" flag;do
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
    esac
done

echo $PWD

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
    echo "Subscription ID (flag -s) no given, getting subID from current context"
    SUBSCRIPTION_ID=$(az account show --query id --output tsv)
fi

echo -e "\nChecking if Resource Group $RESOURCE_GROUP Exists"
if [ $(az group exists -n $RESOURCE_GROUP) = true ];
then
    echo "$RESOURCE_GROUP exists"
else
    echo "$RESOURCE_GROUP does not exist"
    echo "Creating new Resource Group $RESOURCE_GROUP"
    if [ -z $LOCATION]; then
        echo "Error: Cannot create Resource group without Location (flag -l)"
        exit 1
    fi
    az group create -l $LOCATION -n $RESOURCE_GROUP
fi

echo -e "\nCreating Service Principal"
sp=$(az ad sp create-for-rbac --scopes /subscriptions/$SUBSCRIPTION_ID/resourceGroups/$RESOURCE_GROUP --role Contributor)
username=$(jq -r '.appId' <<< "$sp")
password=$(jq -r '.password' <<< "$sp")

echo -e "\nSpinning up Azure Kubernetes Cluster $CLUSTER_NAME"
az aks create --resource-group $RESOURCE_GROUP --name $CLUSTER_NAME --enable-addons monitoring --generate-ssh-keys --service-principal $username --client-secret $password

echo -e "\nGetting Credentials for Cluster $CLUSTER_NAME"
az aks get-credentials --resource-group $RESOURCE_GROUP --name $CLUSTER_NAME
echo -e "\n"

DRIVER_NAME=$(dirname "$(realpath ${BASH_SOURCE[0]})")
source "$DRIVER_NAME/cosi-install.sh"