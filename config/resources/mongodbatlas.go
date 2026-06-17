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
		r.TerraformResource.DeprecationMessage = "This resource is deprecated and will be removed in the next major version. Please use AdvancedCluster (mongodbatlas_advanced_cluster) instead."
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_advanced_cluster", func(r *config.Resource) {
		r.Version = refs.VersionV1Alpha3
		r.Kind = "AdvancedCluster"
		r.UseAsync = true
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_flex_cluster", func(r *config.Resource) {
		r.UseAsync = true
		r.Kind = "FlexCluster"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cluster_outage_simulation", func(r *config.Resource) {
		r.ShortGroup = "cluster"
		r.Kind = "OutageSimulation"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_mongodb_employee_access_grant", func(r *config.Resource) {
		r.Kind = "EmployeeAccessGrant"
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
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
		r.ExternalName.GetIDFn = refs.GetIDFromParamsAndExternalName("-", 1, refs.OrgID)
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromID("-", 1, 0)
	})

	p.AddResourceConfigurator("mongodbatlas_team_project_assignment", func(r *config.Resource) {
		r.ShortGroup = "team"
		r.Kind = "ProjectAssignment"
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_api_key", func(r *config.Resource) {
		r.Kind = "APIKey"
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
		r.ExternalName.GetIDFn = refs.GetIDFromParamsAndExternalName("-", 1, refs.OrgID)
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromID("-", 1, 0)
	})

	p.AddResourceConfigurator("mongodbatlas_api_key_project_assignment", func(r *config.Resource) {
		r.Kind = "APIKeyProjectAssignment"
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
		r.Kind = "EncryptionAtRest"
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
		r.Kind = "EventTrigger"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
		r.ExternalName.GetIDFn = refs.GetIDFromParamsAndExternalName("-", 2, refs.ProjectID, "app_id")
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromID("-", 2, 0)
	})

	p.AddResourceConfigurator("mongodbatlas_global_cluster_config", func(r *config.Resource) {
		r.Kind = "GlobalClusterConfig"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_log_integration", func(r *config.Resource) {
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
	})

	p.AddResourceConfigurator("mongodbatlas_org_invitation", func(r *config.Resource) {
		r.ShortGroup = "org"
		r.Kind = "Invitation"
		r.TerraformResource.DeprecationMessage = "This resource is deprecated. Migrate to mongodbatlas_cloud_user_org_assignment for managing organization membership."
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_custom_dns_configuration_cluster_aws", func(r *config.Resource) {
		r.Kind = "CustomDNSConfigurationClusterAWS"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_maintenance_window", func(r *config.Resource) {
		r.Kind = "MaintenanceWindow"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_service_account_project_assignment", func(r *config.Resource) {
		r.Kind = "ServiceAccountProjectAssignment"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_service_account", func(r *config.Resource) {
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
		r.Kind = "PushBasedLogExport"
		r.TerraformResource.DeprecationMessage = "This resource is deprecated and will be removed in the next major version. Please use mongodbatlas_log_integration instead."
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_resource_policy", func(r *config.Resource) {
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
		r.Kind = "OnlineArchive"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
		// ID format: {project_id}-{cluster_name}-{archive_id}
		// cluster_name may contain dashes, but the extracted archive_id is a hex
		// ID that never contains dashes, so Split+last is safe here.
		r.ExternalName.GetIDFn = refs.GetIDFromParamsAndExternalName("-", 2, refs.ProjectID, refs.ClusterName)
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromSegment("-")
	})

	p.AddResourceConfigurator("mongodbatlas_access_list_api_key", func(r *config.Resource) {
		r.Kind = "AccessListAPIKey"
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
		r.ExternalName.GetIDFn = refs.AccessListGetIDFn(refs.OrgID, "api_key_id")
	})
}
