echo -e "\nDeleting Kube Resources"
kubectl delete -k ./. 
echo -e "\n"

echo -e "\nDeleting COSI Driver and Sidecar"
DRIVER_NAME=$(dirname $(dirname "$(realpath ${BASH_SOURCE[0]})"))
cd $DRIVER_NAME
echo $(pwd)

echo -e "\nDeleting COSI Image"
PROPERTY_FILE="./resources/cosi-driver-azure.properties"
export REGISTRY=`cat $PROPERTY_FILE | grep AZURE_DRIVER_IMAGE_ORG | cut -d'=' -f2`
export IMAGE_VERSION=`cat $PROPERTY_FILE | grep AZURE_DRIVER_IMAGE_VERSION | cut -d'=' -f2`
echo "Registry: $REGISTRY"
echo "Version: $IMAGE_VERSION"
docker image rm "$REGISTRY/azure-cosi-driver:$IMAGE_VERSION"

echo "Deleting CRD's for COSI and COSI Controller"
kubectl delete -k github.com/kubernetes-sigs/container-object-storage-interface-controller
kubectl delete -k github.com/kubernetes-sigs/container-object-storage-interface-api