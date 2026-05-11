package alert

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

const group = "alert"

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_alert_configuration", func(r *config.Resource) {
		r.ShortGroup = group
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})
}
