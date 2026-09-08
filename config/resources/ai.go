package resources

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

const groupAI = "ai"

func ConfigureAI(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_ai_model_api_key", func(r *config.Resource) {
		r.ShortGroup = groupAI
		r.ExternalName = importJoinedIDAssigned([]string{refs.ProjectID, "api_key_id"}, "/", "api_key_id")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_ai_model_rate_limit", func(r *config.Resource) {
		r.ShortGroup = groupAI
		r.ExternalName = importJoinedID([]string{refs.ProjectID, "cloud", "geography", "model_group_name"}, "/", "model_group_name")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})
}
