package constant

const (
	// BucketClassParameterFields
	BucketUnitTypeField                 = "bucketunittype"
	CreateBucketField                   = "createbucket"
	CreateStorageAccountField           = "createstorageaccount"
	SubscriptionIDField                 = "subscriptionid"
	StorageAccountNameField             = "storageaccountname"
	ContainerNameField                  = "containername"
	RegionField                         = "region"
	AccessTierField                     = "accesstier"
	SKUNameField                        = "skuname"
	ResourceGroupField                  = "resourcegroup"
	AllowBlobAccessField                = "allowblobaccess"
	AllowSharedAccessKeyField           = "allowsharedaccesskey"
	EnableBlobVersioningField           = "enableblobversioning"
	EnableBlobDeleteRetentionField      = "enableblobdeleteretention"
	BlobDeleteRetentionDaysField        = "blobdeleteretentiondays"
	EnableContainerDeleteRetentionField = "enablecontainerdeleteretention"
	ContainerDeleteRetentionDaysField   = "containerdeleteretentiondays"
)

type BucketUnitType int
type AccessTier int
type SKU int
type Kind int

const (
	Container BucketUnitType = iota
	StorageAccount
)

const (
	Hot AccessTier = iota
	Cool
	Archive
)

const (
	StandardLRS SKU = iota
	StandardGRS
	StandardRAGRS
	PremiumLRS
)

const (
	StorageV2 Kind = iota
	Storage
	BlobStorage
	BlockBlobStorage
	FileStorage
)

func (b BucketUnitType) String() string {
	switch b {
	case Container:
		return "container"
	case StorageAccount:
		return "storageaccount"
	}
	return "unknown"
}

func (s SKU) String() string {
	switch s {
	case StandardLRS:
		return "Standard_LRS"
	case StandardGRS:
		return "Standard_GRS"
	case StandardRAGRS:
		return "Standard_RAGRS"
	case PremiumLRS:
		return "Premium_LRS"
	}
	return "unknown"
}

func (a AccessTier) String() string {
	switch a {
	case Hot:
		return "hot"
	case Cool:
		return "cool"
	case Archive:
		return "archive"
	}
	return "unknown"
}

func (a Kind) String() string {
	switch a {
	case StorageV2:
		return "StorageV2"
	case Storage:
		return "Storage"
	case BlobStorage:
		return "BlobStorage"
	case BlockBlobStorage:
		return "BlockBlobStorage"
	case FileStorage:
		return "FileStorage"
	}
	return "unknown"
}
