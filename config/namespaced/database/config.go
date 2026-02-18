package database

import (
	"context"

	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/common"
)

const group = "database"

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
		r.ShortGroup = group
		r.Kind = "CustomRole"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_x509_authentication_database_user", func(r *config.Resource) {
		r.ShortGroup = group
		r.Kind = "X509UserAuthentication"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})
}
