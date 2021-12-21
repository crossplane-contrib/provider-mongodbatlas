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
	xpref "github.com/crossplane/crossplane-runtime/pkg/reference"
	xpresource "github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/terrajet/pkg/resource"
	"github.com/pkg/errors"
)

const (
	// ErrFmtNoAttribute is an error string for not-found attributes
	ErrFmtNoAttribute = `"attribute not found: %s`
	// ErrFmtUnexpectedType is an error string for attribute map values of unexpected type
	ErrFmtUnexpectedType = `unexpected type for attribute %s: Expecting a string`

	commonConfigPackagePath = "github.com/crossplane-contrib/provider-jet-mongodbatlas/config/common"
	// ExtractResourceIDFuncPath holds the MongoDBAtlas resource ID extractor func name
	ExtractResourceIDFuncPath = commonConfigPackagePath + ".ExtractResourceID()"
)

// GetAttributeValue reads a string attribute from the specified map
func GetAttributeValue(attrMap map[string]interface{}, attr string) (string, error) {
	v, ok := attrMap[attr]
	if !ok {
		return "", errors.Errorf(ErrFmtNoAttribute, attr)
	}
	vStr, ok := v.(string)
	if !ok {
		return "", errors.Errorf(ErrFmtUnexpectedType, attr)
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
