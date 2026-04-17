package alert

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	common "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/common"
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

		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 1, "project_id")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("-", 1, 0)
	})
}
