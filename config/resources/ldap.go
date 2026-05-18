package resources

import (
	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/types/comments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

func ConfigureLDAP(p *config.Provider, pwGen func(string, string) config.NewInitializerFn) {
	p.AddResourceConfigurator("mongodbatlas_ldap_configuration", func(r *config.Resource) {
		r.ShortGroup = "ldap"
		r.Kind = "Configuration"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
		desc, _ := comments.New("If true, the bind password will be auto-generated and"+
			" stored in the Secret referenced by the bindPasswordSecretRef field.",
			comments.WithTFTag("-"))
		r.TerraformResource.Schema["auto_generate_bind_password"] = &schema.Schema{
			Type:        schema.TypeBool,
			Optional:    true,
			Description: desc.String(),
		}
		r.InitializerFns = append(r.InitializerFns,
			pwGen(
				"spec.forProvider.bindPasswordSecretRef",
				"spec.forProvider.autoGenerateBindPassword",
			))
		r.TerraformResource.Schema["bind_password"].Description = "Password to authenticate the bind user." +
			" If you set autoGenerateBindPassword to true, the Secret referenced here will be" +
			" created or updated with the generated password if it does not already contain one."
	})

	p.AddResourceConfigurator("mongodbatlas_ldap_verify", func(r *config.Resource) {
		r.ShortGroup = "ldap"
		r.Kind = "Verify"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
		r.ExternalName.GetIDFn = refs.GetIDFromParamsAndExternalName("-", 1, refs.ProjectID)
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromID("-", 1, 0)
	})
}
