package resources

import (
	"context"

	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

func ConfigureDatabase(p *config.Provider, pwGen func(string, string) config.NewInitializerFn) {
	p.AddResourceConfigurator("mongodbatlas_database_user", func(r *config.Resource) {
		r.Version = refs.VersionV1Alpha3
		r.ExternalName = importJoinedID([]string{refs.ProjectID, "username", "auth_database_name"}, "username")
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"x509_type", "ldap_auth_type", "aws_iam_type"},
		}
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
		r.InitializerFns = append(r.InitializerFns,
			pwGen(
				"spec.forProvider.passwordSecretRef",
				"spec.writeConnectionSecretToRef",
			))
		r.TerraformResource.Schema["password"].Description = "Password for the " +
			"database user. Do not set passwordSecretRef directly, it is wired " +
			"automatically by the initializer. Instead, use writeConnectionSecretToRef: " +
			"pre-populate the Secret with a 'password' key for BYOP, or leave it " +
			"empty for auto-generation. See docs/password-management.md."
	})

	p.AddResourceConfigurator("mongodbatlas_custom_db_role", func(r *config.Resource) {
		r.ShortGroup = "database"
		r.Kind = "CustomRole"
		r.ExternalName = importJoinedID([]string{refs.ProjectID, refs.RoleName}, refs.RoleName)
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_x509_authentication_database_user", func(r *config.Resource) {
		r.ShortGroup = "database"
		r.Kind = "X509UserAuthentication"
		r.ExternalName = importJoinedID([]string{refs.ProjectID}, refs.ProjectID)
		// Terraform Read decodes username from the state ID to read the user
		// certificates. Without it, Read takes the project-level path.
		projectGetIDFn := r.ExternalName.GetIDFn
		r.ExternalName.GetIDFn = func(ctx context.Context, externalName string, parameters, providerConfig map[string]any) (string, error) {
			id, err := projectGetIDFn(ctx, externalName, parameters, providerConfig)
			username, _ := parameters["username"].(string)
			if err != nil || id == "" || username == "" {
				return id, err
			}
			values := decodeAtlasStateID(id)
			values["username"] = username
			return encodeAtlasStateID(values), nil
		}
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})
}
