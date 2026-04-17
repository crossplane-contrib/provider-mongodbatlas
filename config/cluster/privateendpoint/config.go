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
		// ID format: {project_id}-{endpoint_id}-{provider_name}-{region}
		// project_id is hex (no dashes). endpoint_id is the external name (hex).
		// provider_name and region are suffix segments to skip from the right.
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 1, "project_id", "provider_name", "region")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("-", 1, 2)
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
		// ID format: {project_id}--{private_link_id}--{endpoint_service_id}--{provider_name}
		// All segments are fixed-format (hex/enum), separator is "--" (unambiguous).
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("--", 2, "project_id", "private_link_id", "provider_name")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("--", 2, 1)
	})
}
