package federated

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	common "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/common"
)

const group = "federation"

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_federated_database_instance", func(r *config.Resource) {
		r.ShortGroup = group
		r.Kind = "DatabaseInstance"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_federated_query_limit", func(r *config.Resource) {
		r.ShortGroup = group
		r.Kind = "QueryLimit"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_privatelink_endpoint_service_data_federation_online_archive", func(r *config.Resource) {
		r.ShortGroup = group
		r.Kind = "PrivateLinkEndpointService"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_federated_settings_identity_provider", func(r *config.Resource) {
		r.ShortGroup = group
		r.Kind = "IdentityProvider"
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 1, "federation_settings_id")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromSegment("-")
	})

	p.AddResourceConfigurator("mongodbatlas_federated_settings_org_config", func(r *config.Resource) {
		r.ShortGroup = group
		r.Kind = "OrgConfigSettings"
		r.References = config.References{
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_federated_settings_org_role_mapping", func(r *config.Resource) {
		r.ShortGroup = group
		r.Kind = "RoleMapping"
		r.References = config.References{
			"federation_settings_id": {
				TerraformName: "mongodbatlas_organization",
			},
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
		}
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 2, "federation_settings_id", "org_id")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromSegment("-", 1)
	})
}
