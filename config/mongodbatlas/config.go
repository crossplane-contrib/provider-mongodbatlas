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

package mongodbatlas

import (
	"context"
	"fmt"
	"strings"

	"github.com/upbound/upjet/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/common"
)

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_cluster", func(r *config.Resource) {
		r.ExternalName = config.NameAsIdentifier
		r.ExternalName.SetIdentifierArgumentFn = common.SetIdentifierFunc
		r.ExternalName.GetExternalNameFn = getExternalNameFunc
		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]interface{}, providerConfig map[string]interface{}) (string, error) {
			parts := strings.Split(externalName, ":")
			if len(parts) != 2 {
				return "", nil
			}
			return common.Base64EncodeTokens("cluster_id", parts[1], "cluster_name", parameters["name"], "project_id", parameters["project_id"], "provider_name", parameters["provider_name"])
		}
		r.UseAsync = true
	})
	p.AddResourceConfigurator("mongodbatlas_advanced_cluster", func(r *config.Resource) {
		r.ExternalName = config.NameAsIdentifier
		r.ExternalName.SetIdentifierArgumentFn = common.SetIdentifierFunc
		r.ExternalName.GetExternalNameFn = getExternalNameFunc
		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]interface{}, providerConfig map[string]interface{}) (string, error) {
			parts := strings.Split(externalName, ":")
			if len(parts) != 2 {
				return "", nil
			}
			return common.Base64EncodeTokens("cluster_id", parts[1], "cluster_name", parameters["name"], "project_id", parameters["project_id"])
		}
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"bi_connector"},
		}
		r.UseAsync = true
	})
}

func getExternalNameFunc(tfstate map[string]interface{}) (string, error) {
	return fmt.Sprintf("%s:%s", tfstate["name"], tfstate["cluster_id"]), nil
}
