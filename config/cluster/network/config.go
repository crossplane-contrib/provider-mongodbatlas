package network

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	common "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/common"
)

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_network_container", func(r *config.Resource) {
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 1, "project_id")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromSegment("-")
	})

	p.AddResourceConfigurator("mongodbatlas_network_peering", func(r *config.Resource) {
		r.References = config.References{
			"container_id": {
				TerraformName: "mongodbatlas_network_container",
			},
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 1, "project_id", "provider_name")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromSegment("-", 1)
	})

}
