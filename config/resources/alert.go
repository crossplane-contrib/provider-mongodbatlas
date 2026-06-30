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
		for _, f := range []string{
			"matcher",
			"metric_threshold_config",
			"notification",
			"threshold_config",
		} {
			r.ServerSideApplyMergeStrategies[f] = config.MergeStrategy{
				ListMergeStrategy: config.ListMergeStrategy{
					MergeStrategy: config.ListTypeAtomic,
				},
			}
		}
	})
}
