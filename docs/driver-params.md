# Driver Parameters

### BucketClass parameters
|Name            | Meaning | Available Value | Mandatory |
|----------------|---------|-----------------|-----------|
| bucketunittype | Decide whether the bucket is a container or a storage account (container by default) | container, storageaccount | yes   |
| createbucket | automatically creates bucket (default yes) | true, false | no   |
| createstorageaccount | automatically creates storage acc | true, false | no |
| subscriptionid | Subscription ID | string | no   |
| storageaccountname | Name of the Storage Account | string | yes   |
| region | Storage Account Region | [availability zones](https://learn.microsoft.com/en-us/azure/reliability/availability-zones-service-support); example format: eastus | yes   |
| accesstier | [manages storage pricing](https://learn.microsoft.com/en-us/azure/storage/blobs/access-tiers-overview) | hot, cool, archive | no   |
| skuname | Stock Keeping Unit | Standard_LRS, Standard_GRS, Standard_RAGRS, Premium_LRS | no   |
| resourcegroup | Resource group of the cluster | string | yes   |
| allowblobaccess |  | true, false | no   |
| allowsharedaccess |  | true, false | no   |
| enableblobversioning |  | true, false | no   |
| enableblobdeleteretention | [adds retention period for deleted blobs](https://learn.microsoft.com/en-us/azure/storage/blobs/soft-delete-blob-enable?tabs=azure-CLI) | true, false | no   |
| blobdeleteretentiondays | days that a blob lasts when deleted | positive int | no   |
| enablecontainerdeleteretention | [adds retention period for deleted containers](https://learn.microsoft.com/en-us/azure/storage/blobs/soft-delete-container-enable?tabs=azure-portal)  | true, false | no   |
| containerdeleteretentiondays | days that a container lasts when deleted  | positive int | no   |

### BucketAccessClass parameters
|Name            | Meaning | Available Value | Mandatory |
|----------------|---------|-----------------|-----------|
| bucketunittype | Decide whether the bucket is a container or a storage account | container, storageaccount | yes   |
| signedversion | Signed storage service version (has default value) | 2015-04-05 or later | no   |
| signedip | specified ip address, or range of ip addresses to accept requests |  | no   |
| validationperiod | how long the token lasts (ms) | uint64(default 7 days) | no   |
| signedprotocol | determins protocol used | "https", "https,http"(default) | no   |
| enablelist | enables list access | true(default), false | no   |
| enableread | enables read access | true(default), false | no   |
| enablewrite | enables write access | true, false | no   |
| enabledelete | enables delete access | true, false | no   |
| enablepermanentdelete | enables permanent delete access | true, false | no   |
| enableadd | enables add operations | true, false | no   |
| enabletags | enables blob tag operations | true, false | no   |
| enablefilter | enables filtering by blob tag | true, false | no   |
| allowservicesignedresourcetype | gives access to service level apis | true, false | no   |
| allowcontainersignedresourcetype | gives access to container level apis | true(default), false | no   |
| allowobjectsignedresourcetype (default)| gives access to object level apis | true(default), false | no   |