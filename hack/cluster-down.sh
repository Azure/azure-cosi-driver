while getopts "n:r:" flag;do
    case "${flag}" in
        n) 
            CLUSTER_NAME=$OPTARG
            echo "Cluster Name: $CLUSTER_NAME"
            ;;
        r) 
            RESOURCE_GROUP=$OPTARG
            echo "Resource Group: $RESOURCE_GROUP"
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

service_principal=$(az aks show --name $CLUSTER_NAME --resource-group $RESOURCE_GROUP --query servicePrincipalProfile.clientId -o tsv)
echo -e "\nDeleting Service Principal $service_principal"
az ad sp delete --id $service_principal

echo -e "\nDeleting Cluster $CLUSTER_NAME from Resource Group $RESOURCE_GROUP"
az aks delete --name $CLUSTER_NAME --resource-group $RESOURCE_GROUP