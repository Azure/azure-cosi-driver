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

echo -e "\nrunning cosi-uninstall.sh"
DRIVER_NAME=$(dirname "$(realpath ${BASH_SOURCE[0]})")
source "$DRIVER_NAME/cosi-uninstall.sh"

echo -e "\nRunning cluster-down.sh"
DRIVER_NAME=$(dirname "$(realpath ${BASH_SOURCE[0]})")
source "$DRIVER_NAME/cluster-down.sh" -n $CLUSTER_NAME -r $RESOURCE_GROUP