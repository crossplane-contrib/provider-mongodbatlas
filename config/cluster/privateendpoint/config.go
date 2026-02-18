package privateendpoint

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	common "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/common"
)

const group = "privateendpoint"

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_private_endpoint_regional_mode", func(r *config.Resource) {
		r.ShortGroup = group
		r.Kind = "RegionalMode"

		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_privatelink_endpoint", func(r *config.Resource) {
		r.ShortGroup = group
		r.Kind = "Resource"

		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 1, "project_id", "provider_name", "region")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromSegment("-", 1)
	})

	p.AddResourceConfigurator("mongodbatlas_privatelink_endpoint_service", func(r *config.Resource) {
		r.ShortGroup = group
		r.Kind = "Service"

		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
			"private_link_id": {
				TerraformName: "mongodbatlas_privatelink_endpoint",
			},
		}
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("--", 2, "project_id", "private_link_id", "provider_name")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromSegment("--", 2)
	})
}
