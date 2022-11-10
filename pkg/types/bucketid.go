package types

import (
	"encoding/base64"
	"encoding/json"
)

// bucketID is returned by the DriverCreateBucket function call as an encoded string with the subID, resource group and the URL of the bucket.
// These details are required by DriverDeleteBucket and DriverGrantBucketAccess.
type BucketID struct {
	SubID         string `json:"subscriptionID"`
	ResourceGroup string `json:"resourceGroup"`
	URL           string `json:"url"`
}

// Marshals bucketID struct into json bytes, then encodes into base64
func (id *BucketID) Encode() (string, error) {
	data, err := json.Marshal(id)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// Decodes base64 string to bucketID pointer struct
func DecodeToBucketID(id string) (*BucketID, error) {
	data, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return nil, err
	}

	bID := &BucketID{}
	err = json.Unmarshal(data, bID)
	if err != nil {
		return nil, err
	}
	return bID, nil
}
