package resources

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

func ConfigurePrivateEndpoint(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_private_endpoint_regional_mode", func(r *config.Resource) {
		r.ShortGroup = "privateendpoint"
		r.Kind = "RegionalMode"
		r.ExternalName = templated("{{ .parameters.project_id }}")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_privatelink_endpoint", func(r *config.Resource) {
		r.ShortGroup = "privateendpoint"
		r.Kind = "Resource"
		r.ExternalName = importJoinedIDOrdered([]string{refs.ProjectID, "private_link_id", refs.ProviderName, refs.Region}, "private_link_id")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_privatelink_endpoint_service", func(r *config.Resource) {
		r.ShortGroup = "privateendpoint"
		r.Kind = "Service"
		r.ExternalName = importJoinedID([]string{refs.ProjectID, "private_link_id", "endpoint_service_id", refs.ProviderName}, "--", "endpoint_service_id")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			"private_link_id": {
				TerraformName: "mongodbatlas_privatelink_endpoint",
			},
		}
	})
}
