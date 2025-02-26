/*
Copyright 2021 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"encoding/base64"
	"fmt"
	"strings"

	xpref "github.com/crossplane/crossplane-runtime/pkg/reference"
	xpresource "github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/upjet/pkg/resource"
	"github.com/pkg/errors"
)

const (
	errUnevenCount = "argument count should be even: expecting key-value pairs"
	// errFmtNoAttribute is an error string for not-found attributes
	errFmtNoAttribute = `"attribute not found: %s`
	// errFmtUnexpectedType is an error string for attribute map values of unexpected type
	errFmtUnexpectedType = `unexpected type for attribute %s: Expecting a string`

	commonConfigPackagePath = "github.com/crossplane-contrib/provider-mongodbatlas/config/common"
	// ExtractResourceIDFuncPath holds the MongoDBAtlas resource ID extractor func name
	ExtractResourceIDFuncPath = commonConfigPackagePath + ".ExtractResourceID()"

	// VersionV1Alpha2 is used as minimum version for all manually configured resources.
	VersionV1Alpha2 = "v1alpha2"
)

const (
	// APISPackagePath is the package path for generated APIs root package
	APISPackagePath = "github.com/crossplane-contrib/provider-mongodbatlas/apis"
)

// GetAttributeValue reads a string attribute from the specified map
func GetAttributeValue(attrMap map[string]interface{}, attr string) (string, error) {
	v, ok := attrMap[attr]
	if !ok {
		return "", errors.Errorf(errFmtNoAttribute, attr)
	}
	vStr, ok := v.(string)
	if !ok {
		return "", errors.Errorf(errFmtUnexpectedType, attr)
	}
	return vStr, nil
}

// ExtractResourceID extracts the value of `spec.atProvider.id`
// from a Terraformed resource. If mr is not a Terraformed
// resource, returns an empty string.
func ExtractResourceID() xpref.ExtractValueFn {
	return func(mr xpresource.Managed) string {
		tr, ok := mr.(resource.Terraformed)
		if !ok {
			return ""
		}
		return tr.GetID()
	}
}

// Base64EncodeTokens base64-encode key-value pairs using a colon
// as a separator between them and concatenate pairs with hyphens
func Base64EncodeTokens(keyVal ...interface{}) (string, error) {
	if len(keyVal)%2 == 1 {
		return "", errors.New(errUnevenCount)
	}
	result := ""
	for i := 0; i < len(keyVal); i += 2 {
		encodedPair := fmt.Sprintf("%s:%s", base64.StdEncoding.EncodeToString([]byte(keyVal[i].(string))), base64.StdEncoding.EncodeToString([]byte(keyVal[i+1].(string))))
		switch result {
		case "":
			result = encodedPair
		default:
			result = fmt.Sprintf("%s-%s", result, encodedPair)
		}
	}
	return result, nil
}

// SetIdentifierFunc sets the identifier attribute `name` from a composite
// external-name where the identifier resides at index 0 of a colon-delimited
// string.
func SetIdentifierFunc(base map[string]interface{}, externalName string) {
	parts := strings.Split(externalName, ":")
	base["name"] = parts[0]
}
