package resources

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

func ConfigureSearch(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_search_deployment", func(r *config.Resource) {
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			refs.ClusterName: {
				TerraformName: refs.TFCluster,
				Extractor:     refs.ExtractParamPath("name", false),
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_search_index", func(r *config.Resource) {
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			refs.ClusterName: {
				TerraformName: refs.TFCluster,
				Extractor:     refs.ExtractParamPath("name", false),
			},
		}
	})
}
