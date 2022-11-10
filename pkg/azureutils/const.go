// Copyright 2021 The Kubernetes Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package azureutils

import (
	"fmt"
	"strings"
)

const (
	TagsDelimiter        = ","
	TagKeyValueDelimiter = "="

	TrueValue                  = "true"
	FalseValue                 = "false"
	StorageAccountTypeField    = "storageaccounttype"
	LocationField              = "location"
	KindField                  = "kind"
	TagsField                  = "tags"
	VNResourceIdsField         = "virtualnetworkresourceids"
	HTTPSTrafficOnlyField      = "httpstrafficonly"
	CreatePrivateEndpointField = "createprivateendpoint"
	HNSEnabledField            = "hnsenabled"
	EnableNFSV3Field           = "enablenfsv3"
	EnableLargeFileSharesField = "enablelargefileshares"
)

// ConvertTagsToMap convert the tags from string to map
// the valid tags format is "key1=value1,key2=value2", which could be converted to
// {"key1": "value1", "key2": "value2"}
func ConvertTagsToMap(tags string) (map[string]string, error) {
	m := make(map[string]string)
	if tags == "" {
		return m, nil
	}
	s := strings.Split(tags, TagsDelimiter)
	for _, tag := range s {
		kv := strings.Split(tag, TagKeyValueDelimiter)
		if len(kv) != 2 {
			return nil, fmt.Errorf("Tags '%s' are invalid, the format should like: 'key1=value1,key2=value2'", tags)
		}
		key := strings.TrimSpace(kv[0])
		if key == "" {
			return nil, fmt.Errorf("Tags '%s' are invalid, the format should like: 'key1=value1,key2=value2'", tags)
		}
		value := strings.TrimSpace(kv[1])
		m[key] = value
	}

	return m, nil
}

func ConvertMapToMapPointer(origin map[string]string) map[string]*string {
	newly := make(map[string]*string)
	for k, v := range origin {
		value := v
		newly[k] = &value
	}
	return newly
}
