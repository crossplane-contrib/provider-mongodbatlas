package database

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/common"
)

const group = "database"

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_database_user", func(r *config.Resource) {
		r.Version = common.VersionV1Alpha3
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"x509_type", "ldap_auth_type", "aws_iam_type"},
		}
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
		r.InitializerFns = append(r.InitializerFns,
			common.PasswordGenerator(
				"spec.forProvider.passwordSecretRef",
				"spec.writeConnectionSecretToRef",
			))
		r.TerraformResource.Schema["password"].Description = "Password for the " +
			"database user. When passwordSecretRef is set, that secret is used " +
			"as the password (Bring Your Own Password). Otherwise, a password is " +
			"auto-generated and written to the secret in writeConnectionSecretToRef."
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
