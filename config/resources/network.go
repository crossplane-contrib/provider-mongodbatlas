package resources

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

func ConfigureNetwork(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_network_container", func(r *config.Resource) {
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_network_peering", func(r *config.Resource) {
		r.References = config.References{
			"container_id": {
				TerraformName: "mongodbatlas_network_container",
			},
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})
}
