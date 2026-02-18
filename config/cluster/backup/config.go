package backup

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

const group = "backup"

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_backup_compliance_policy", func(r *config.Resource) {
		r.ShortGroup = group
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})
}
