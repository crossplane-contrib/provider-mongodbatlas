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
		r.ExternalName = importJoinedID([]string{refs.ProjectID, refs.Name}, "--", refs.Name)
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_federated_query_limit", func(r *config.Resource) {
		r.ShortGroup = groupFederated
		r.Kind = "QueryLimit"
		r.ExternalName = importJoinedID([]string{refs.ProjectID, "tenant_name", "limit_name"}, "--", "limit_name")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_privatelink_endpoint_service_data_federation_online_archive", func(r *config.Resource) {
		r.ShortGroup = groupFederated
		r.Kind = "PrivateLinkEndpointService"
		r.ExternalName = importJoinedID([]string{refs.ProjectID, "endpoint_id"}, "--", "endpoint_id")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_federated_settings_identity_provider", func(r *config.Resource) {
		r.ShortGroup = groupFederated
		r.Kind = "SettingsIdentityProvider"
		r.ExternalName = importJoinedID([]string{"federation_settings_id"}, "-", "okta_idp_id")
	})

	p.AddResourceConfigurator("mongodbatlas_federated_settings_org_config", func(r *config.Resource) {
		r.ShortGroup = groupFederated
		r.Kind = "OrgConfigSettings"
		r.ExternalName = importJoinedID([]string{"federation_settings_id", refs.OrgID}, "-", refs.OrgID)
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_federated_settings_org_role_mapping", func(r *config.Resource) {
		r.ShortGroup = groupFederated
		r.Kind = "RoleMapping"
		r.ExternalName = importJoinedID([]string{"federation_settings_id", refs.OrgID}, "-", "role_mapping_id")
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})
}
