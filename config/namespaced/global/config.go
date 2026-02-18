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

	p.AddResourceConfigurator("mongodbatlas_organization", func(r *config.Resource) {
		r.ShortGroup = "global"
		r.Kind = "Organization"
	})

	p.AddResourceConfigurator("mongodbatlas_org_invitation", func(r *config.Resource) {
		r.ShortGroup = "global"
		r.Kind = "Invitation"
		r.TerraformResource.DeprecationMessage = "This resource is deprecated. Migrate to mongodbatlas_cloud_user_org_assignment for managing organization membership."

		r.References = config.References{
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
		}
	})

}
