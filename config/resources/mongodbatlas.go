package resources

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

func ConfigureMongoDBAtlas(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_cluster", func(r *config.Resource) {
		r.Version = refs.VersionV1Alpha2
		r.UseAsync = true
		r.Kind = "Cluster"
		r.ExternalName = importJoinedIDMapped([]string{refs.ProjectID, refs.Name}, map[string]string{refs.ProjectID: refs.ProjectID, refs.Name: refs.ClusterName}, refs.ClusterName)
		r.TerraformResource.DeprecationMessage = "This resource is deprecated and will be removed in the next major version. Please use AdvancedCluster (mongodbatlas_advanced_cluster) instead."
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_advanced_cluster", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Version = refs.VersionV1Alpha3
		r.Kind = "AdvancedCluster"
		r.UseAsync = true
		r.ExternalName = templated("{{ .parameters.project_id }}-{{ .parameters.name }}")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_flex_cluster", func(r *config.Resource) {
		r.ShortGroup = ""
		r.UseAsync = true
		r.Kind = "FlexCluster"
		r.ExternalName = templated("{{ .parameters.project_id }}-{{ .parameters.name }}")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cluster_outage_simulation", func(r *config.Resource) {
		r.ShortGroup = "cluster"
		r.Kind = "OutageSimulation"
		r.ExternalName = config.IdentifierFromProvider
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_mongodb_employee_access_grant", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "EmployeeAccessGrant"
		r.ExternalName = templated("{{ .parameters.project_id }}/{{ .parameters.cluster_name }}")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_auditing", func(r *config.Resource) {
		r.UseAsync = true
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_team", func(r *config.Resource) {
		r.Kind = "Team"
		r.ExternalName = importJoinedID([]string{refs.OrgID}, "-", "id")
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_team_project_assignment", func(r *config.Resource) {
		r.ShortGroup = "team"
		r.Kind = "ProjectAssignment"
		r.ExternalName = templated("{{ .parameters.project_id }}/{{ .parameters.team_id }}")
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_api_key", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "APIKey"
		r.ExternalName = importJoinedID([]string{refs.OrgID}, "-", "api_key_id")
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_api_key_project_assignment", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "APIKeyProjectAssignment"
		r.ExternalName = templated("{{ .parameters.project_id }}/{{ .parameters.api_key_id }}")
		r.References = config.References{
			"api_key_id": {
				TerraformName: "mongodbatlas_api_key",
			},
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_encryption_at_rest", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "EncryptionAtRest"
		r.ExternalName = templated("{{ .parameters.project_id }}")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_encryption_at_rest_private_endpoint", func(r *config.Resource) {
		r.ShortGroup = "encryptionatrest"
		r.Kind = "PrivateEndpoint"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
		r.ExternalName.GetIDFn = refs.GetIDFromParamsAndExternalName("-", 2, refs.ProjectID, "cloud_provider")
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromID("-", 2, 0)
	})

	p.AddResourceConfigurator("mongodbatlas_event_trigger", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "EventTrigger"
		r.ExternalName = importJoinedID([]string{refs.ProjectID, "app_id"}, "--", "trigger_id")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_global_cluster_config", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "GlobalClusterConfig"
		r.ExternalName = importJoinedID([]string{refs.ProjectID, refs.ClusterName}, "-", refs.ClusterName)
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_log_integration", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "LogIntegration"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
		// otel_supplied_headers is a list-of-objects field marked sensitive in the
		// Terraform schema. The code generator only supports string-like sensitive
		// fields, so we clear the flag here. The nested value field is still
		// sensitive and will be handled by users via Kubernetes secrets.
		if sch, ok := r.TerraformResource.Schema["otel_supplied_headers"]; ok {
			sch.Sensitive = false
		}
		r.ExternalName.GetIDFn = refs.GetIDFromParamsAndExternalName("/", 1, refs.ProjectID)
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromIDOrState("/", 1, 0, "type")
	})

	p.AddResourceConfigurator(refs.TFOrganization, func(r *config.Resource) {
		r.Kind = "Organization"
		r.ExternalName = importJoinedID([]string{}, "-", "org_id")
	})

	p.AddResourceConfigurator("mongodbatlas_org_invitation", func(r *config.Resource) {
		r.ShortGroup = "org"
		r.Kind = "Invitation"
		r.ExternalName = importJoinedIDHidden([]string{refs.OrgID, "username"}, "-", "invitation_id")
		r.TerraformResource.DeprecationMessage = "This resource is deprecated. Migrate to mongodbatlas_cloud_user_org_assignment for managing organization membership."
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_custom_dns_configuration_cluster_aws", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "CustomDNSConfigurationClusterAWS"
		r.ExternalName = templated("{{ .parameters.project_id }}")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_maintenance_window", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "MaintenanceWindow"
		r.ExternalName = templated("{{ .parameters.project_id }}")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_service_account_project_assignment", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "ServiceAccountProjectAssignment"
		r.ExternalName = templated("{{ .parameters.project_id }}/{{ .parameters.client_id }}")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_service_account", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "ServiceAccount"
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
		r.ExternalName.GetIDFn = refs.GetIDFromParamsAndExternalName("/", 1, refs.OrgID)
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromIDOrState("/", 1, 0, "client_id")
	})

	p.AddResourceConfigurator("mongodbatlas_service_account_secret", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "ServiceAccountSecret"
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
		r.ExternalName.GetIDFn = refs.GetIDFromParamsAndExternalName("/", 2, refs.OrgID, "client_id")
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromIDOrState("/", 2, 0, "secret_id")
	})

	p.AddResourceConfigurator("mongodbatlas_service_account_access_list_entry", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "ServiceAccountAccessListEntry"
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"ip_address"},
		}
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
		r.ExternalName.GetIDFn = refs.AccessListGetIDFn(refs.OrgID, "client_id")
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromAccessListState(refs.OrgID)
	})

	p.AddResourceConfigurator("mongodbatlas_push_based_log_export", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "PushBasedLogExport"
		r.ExternalName = templated("{{ .parameters.project_id }}")
		r.TerraformResource.DeprecationMessage = "This resource is deprecated and will be removed in the next major version. Please use mongodbatlas_log_integration instead."
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_resource_policy", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "ResourcePolicy"
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
		r.ExternalName.GetIDFn = refs.GetIDFromParamsAndExternalName("-", 1, refs.OrgID)
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromID("-", 1, 0)
	})

	p.AddResourceConfigurator("mongodbatlas_online_archive", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "OnlineArchive"
		r.ExternalName = importJoinedID([]string{refs.ProjectID, refs.ClusterName}, "-", "archive_id")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_project_api_key", func(r *config.Resource) {
		r.ExternalName = importJoinedID([]string{refs.ProjectID}, "-", "api_key_id")
	})

	p.AddResourceConfigurator("mongodbatlas_access_list_api_key", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "AccessListAPIKey"
		r.ExternalName = accessListImportJoinedID([]string{refs.OrgID, "api_key_id"})
		r.ExternalName.GetIDFn = refs.AccessListGetIDFn(refs.OrgID, "api_key_id")
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"ip_address"},
		}
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
			"api_key_id": {
				TerraformName: "mongodbatlas_api_key",
			},
		}
	})
}
