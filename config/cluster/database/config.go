/*
Copyright 2022 The Crossplane Authors.

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

package database

import (
	"context"

	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/common"
)

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_database_user", func(r *config.Resource) {
		r.Version = common.VersionV1Alpha2
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"x509_type", "ldap_auth_type", "aws_iam_type"},
		}
		r.ExternalName.OmittedFields = []string{"username"}
		r.ExternalName.SetIdentifierArgumentFn = func(base map[string]interface{}, externalName string) {
			base["username"] = externalName
		}
		r.ExternalName.GetExternalNameFn = func(tfstate map[string]interface{}) (string, error) {
			return tfstate["username"].(string), nil
		}
		r.ExternalName.GetIDFn = func(_ context.Context, _ string, parameters map[string]interface{}, providerConfig map[string]interface{}) (string, error) {
			return common.Base64EncodeTokens("auth_database_name", parameters["auth_database_name"], "project_id", parameters["project_id"], "username", parameters["username"])
		}
	})

	p.AddResourceConfigurator("mongodbatlas_custom_db_role", func(r *config.Resource) {
		r.ShortGroup = "database"
		r.Kind = "CustomRole"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_x509_authentication_database_user", func(r *config.Resource) {
		r.ShortGroup = "database"
		r.Kind = "X509UserAuthentication"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})
}
