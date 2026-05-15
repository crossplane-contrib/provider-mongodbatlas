package resources

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

const groupFederated = "federated"

func ConfigureFederated(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_federated_database_instance", func(r *config.Resource) {
		r.ShortGroup = groupFederated
		r.Kind = "DatabaseInstance"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_federated_query_limit", func(r *config.Resource) {
		r.ShortGroup = groupFederated
		r.Kind = "QueryLimit"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_privatelink_endpoint_service_data_federation_online_archive", func(r *config.Resource) {
		r.ShortGroup = groupFederated
		r.Kind = "PrivateLinkEndpointService"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_federated_settings_identity_provider", func(r *config.Resource) {
		r.ShortGroup = groupFederated
		r.Kind = "SettingsIdentityProvider"
		r.ExternalName.GetIDFn = refs.GetIDFromParamsAndExternalName("-", 1, "federation_settings_id")
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromID("-", 1, 0)
	})

	p.AddResourceConfigurator("mongodbatlas_federated_settings_org_config", func(r *config.Resource) {
		r.ShortGroup = groupFederated
		r.Kind = "OrgConfigSettings"
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_federated_settings_org_role_mapping", func(r *config.Resource) {
		r.ShortGroup = groupFederated
		r.Kind = "RoleMapping"
		r.References = config.References{
			"federation_settings_id": {
				TerraformName: refs.TFOrganization,
			},
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
		// ID format: {federation_settings_id}-{org_id}-{role_mapping_id}
		// All three segments are hex IDs (no dashes).
		r.ExternalName.GetIDFn = refs.GetIDFromParamsAndExternalName("-", 2, "federation_settings_id", refs.OrgID)
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromID("-", 2, 0)
	})
}
