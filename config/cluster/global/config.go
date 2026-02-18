package global

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_global_cluster_config", func(r *config.Resource) {
		r.ShortGroup = "global"
		r.Kind = "ClusterConfig"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

}
