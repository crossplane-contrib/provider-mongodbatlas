package resources

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

func ConfigureProject(p *config.Provider) {
	p.AddResourceConfigurator(refs.TFProject, func(r *config.Resource) {
		r.ExternalName = providerAssignedID()
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_third_party_integration", func(r *config.Resource) {
		r.ShortGroup = ""
		// Terraform Read looks the integration up by project_id and type and
		// sets the ID to the Atlas integration ID, so the external name is that
		// provider-assigned ID (the default identifierFromProvider).
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_project_invitation", func(r *config.Resource) {
		r.ExternalName = importJoinedIDAssigned([]string{refs.ProjectID, "username"}, "invitation_id")
		r.TerraformResource.DeprecationMessage = "This resource is deprecated. Migrate to mongodbatlas_cloud_user_project_assignment for managing project membership."
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"ip_address"},
		}
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_project_ip_access_list", func(r *config.Resource) {
		r.ExternalName = accessListImportJoinedID([]string{refs.ProjectID})
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"ip_address"},
		}
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_project_service_account_access_list_entry", func(r *config.Resource) {
		r.ShortGroup = "project"
		r.Kind = "ServiceAccountAccessListEntry"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			refs.ClientID: {
				TerraformName: refs.TFProjectServiceAccount,
				Extractor:     refs.ExtractParamPath(refs.ClientID, true),
			},
		}
		r.ExternalName.GetIDFn = refs.AccessListGetIDFn(refs.ProjectID, "client_id")
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromAccessListState(refs.ProjectID)
	})

	p.AddResourceConfigurator("mongodbatlas_project_service_account_secret", func(r *config.Resource) {
		r.ShortGroup = "project"
		r.Kind = "ServiceAccountSecret"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			refs.ClientID: {
				TerraformName: refs.TFProjectServiceAccount,
				Extractor:     refs.ExtractParamPath(refs.ClientID, true),
			},
		}
		r.ExternalName.GetIDFn = refs.GetIDFromParamsAndExternalName("/", 2, refs.ProjectID, "client_id")
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromIDOrState("/", 2, 0, "secret_id")
	})

	p.AddResourceConfigurator(refs.TFProjectMCPConfig, func(r *config.Resource) {
		r.ShortGroup = "project"
		r.Kind = "MCPConfig"
		// Atlas returns both ip_address and cidr_block for an entry, but
		// Terraform accepts exactly one of them per entry.
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"ip_access_list"},
		}
		r.ExternalName = computedKeyID(refs.MCPConfigID)
		r.TerraformPluginFrameworkIsStateEmptyFn = refs.StateEmptyWhenAttributeUnset(refs.MCPConfigID)
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_project_mcp_config_secret", func(r *config.Resource) {
		r.ShortGroup = "project"
		r.Kind = "MCPConfigSecret"
		r.ExternalName = computedKeyID("secret_id")
		r.TerraformPluginFrameworkIsStateEmptyFn = refs.StateEmptyWhenAttributeUnset("secret_id")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			refs.MCPConfigID: {
				TerraformName: refs.TFProjectMCPConfig,
				Extractor:     refs.ExtractParamPath(refs.MCPConfigID, true),
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_project_service_account", func(r *config.Resource) {
		r.ShortGroup = "project"
		r.Kind = "ServiceAccount"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
		r.ExternalName.GetIDFn = refs.GetIDFromParamsAndExternalName("/", 1, refs.ProjectID)
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromIDOrState("/", 1, 0, refs.ClientID)
		r.ExternalName.SetIdentifierArgumentFn = refs.SetIdentifierArgument(refs.ClientID)
		r.TerraformPluginFrameworkIsStateEmptyFn = refs.StateEmptyWhenAttributeUnset(refs.ClientID)
	})
}
