package ldap

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	common "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/common"
)

const group = "ldap"

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_ldap_configuration", func(r *config.Resource) {
		r.ShortGroup = group
		r.Kind = "Configuration"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_ldap_verify", func(r *config.Resource) {
		r.ShortGroup = group
		r.Kind = "Verify"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 1, "project_id")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("-", 1, 0)
	})

}
