package resources

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

func ConfigureAlert(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_alert_configuration", func(r *config.Resource) {
		r.ShortGroup = "alert"
		r.ExternalName = importJoinedID([]string{refs.ProjectID}, "-", "id")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})
}
