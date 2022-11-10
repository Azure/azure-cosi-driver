echo "Getting CRD's for COSI and COSI Controller"
kubectl create -k github.com/kubernetes-sigs/container-object-storage-interface-api
kubectl create -k github.com/kubernetes-sigs/container-object-storage-interface-controller

echo -e "\nInstalling COSI Driver and Sidecar"
DRIVER_NAME=$(dirname $(dirname "$(realpath ${BASH_SOURCE[0]})"))
cd $DRIVER_NAME
echo $(pwd)

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

echo -e "\nCreating Kube Resources"
kubectl create -k ./. 
echo -e "\n"