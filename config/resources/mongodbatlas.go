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
		r.ExternalName = importJoinedIDMapped([]string{refs.ProjectID, refs.Name, refs.ProviderName}, map[string]string{refs.ProjectID: refs.ProjectID, refs.Name: refs.ClusterName, refs.ProviderName: refs.ProviderName}, refs.ClusterName)
		r.TerraformResource.DeprecationMessage = "This resource is deprecated and will be removed in the next major version. Please use AdvancedCluster (mongodbatlas_advanced_cluster) instead."
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
		for _, f := range []string{
			"advanced_configuration",
			"bi_connector_config",
			"labels",
			"pinned_fcv",
			"replication_specs",
			"tags",
		} {
			r.ServerSideApplyMergeStrategies[f] = config.MergeStrategy{
				ListMergeStrategy: config.ListMergeStrategy{
					MergeStrategy: config.ListTypeAtomic,
				},
			}
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
		r.ServerSideApplyMergeStrategies["replication_specs"] = config.MergeStrategy{
			ListMergeStrategy: config.ListMergeStrategy{
				MergeStrategy: config.ListTypeAtomic,
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
		r.ExternalName = importJoinedID([]string{refs.ProjectID, refs.ClusterName}, refs.ClusterName)
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			// Read the name from the cluster's observation: it is only set once
			// the cluster is created and IDLE. A simulation started earlier keeps
			// the cluster UPDATING, so its create would never finish.
			refs.ClusterName: {
				TerraformName: refs.TFCluster,
				Extractor:     refs.ExtractParamPath(refs.Name, true),
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cluster_overload_simulation", func(r *config.Resource) {
		r.ShortGroup = "cluster"
		r.Kind = "OverloadSimulation"
		// Create and Delete wait up to 5 minutes for the simulation state.
		r.UseAsync = true
		r.ExternalName = computedKeyID("simulation_id")
		r.TerraformPluginFrameworkIsStateEmptyFn = refs.StateEmptyWhenAttributeUnset("simulation_id")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			// Same as OutageSimulation: wait until the cluster is created and IDLE.
			refs.ClusterName: {
				TerraformName: refs.TFCluster,
				Extractor:     refs.ExtractParamPath(refs.Name, true),
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cluster_adaptive_settings", func(r *config.Resource) {
		r.ShortGroup = "cluster"
		r.Kind = "AdaptiveSettings"
		r.ExternalName = templated("{{ .parameters.project_id }}/{{ .parameters.cluster_name }}")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			refs.ClusterName: {
				TerraformName: refs.TFCluster,
				Extractor:     refs.ExtractParamPath(refs.Name, true),
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
			refs.ClusterName: {
				TerraformName: refs.TFCluster,
				Extractor:     refs.ExtractParamPath(refs.Name, false),
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
		r.ExternalName = importJoinedIDAssigned([]string{refs.OrgID, "id"}, "id")
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
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			refs.TeamID: {
				TerraformName: refs.TFTeam,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_api_key", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "APIKey"
		r.ExternalName = importJoinedIDAssigned([]string{refs.OrgID, refs.APIKeyID}, refs.APIKeyID)
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
		// Terraform Read lists the project's API keys and reports a key
		// that is not assigned to the project as this error.
		r.ExternalName.IsNotFoundDiagnosticFn = refs.NotFoundWhenErrorSummary("API key not found")
		r.References = config.References{
			refs.APIKeyID: {
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
			"aws_kms_config.role_id": {
				TerraformName: refs.TFProviderAccessAuthorization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_encryption_at_rest_private_endpoint", func(r *config.Resource) {
		r.ShortGroup = "encryptionatrest"
		r.Kind = "PrivateEndpoint"
		r.ExternalName = providerAssignedID()
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_event_trigger", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "EventTrigger"
		r.ExternalName = importJoinedIDAssigned([]string{refs.ProjectID, "app_id", "trigger_id"}, "trigger_id")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_global_cluster_config", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "GlobalClusterConfig"
		r.ExternalName = importJoinedID([]string{refs.ProjectID, refs.ClusterName}, refs.ClusterName)
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			refs.ClusterName: {
				TerraformName: refs.TFCluster,
				Extractor:     refs.ExtractParamPath(refs.Name, false),
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
			refs.IAMRoleID: {
				TerraformName: refs.TFProviderAccessAuthorization,
			},
		}
		// otel_supplied_headers is a list-of-objects field marked sensitive in the
		// Terraform schema. The code generator only supports string-like sensitive
		// fields, so we clear the flag here. The nested value field is still
		// sensitive and will be handled by users via Kubernetes secrets.
		if sch, ok := r.TerraformResource.Schema["otel_supplied_headers"]; ok {
			sch.Sensitive = false
		}
		r.ExternalName = computedKeyID("integration_id")
		r.TerraformPluginFrameworkIsStateEmptyFn = refs.StateEmptyWhenAttributeUnset("integration_id")
	})

	p.AddResourceConfigurator("mongodbatlas_metric_integration", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "MetricIntegration"
		r.ExternalName = computedKeyID("metric_integration_id")
		r.TerraformPluginFrameworkIsStateEmptyFn = refs.StateEmptyWhenAttributeUnset("metric_integration_id")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator(refs.TFOrganization, func(r *config.Resource) {
		r.Kind = "Organization"
		r.ExternalName = importJoinedIDAssigned([]string{"org_id"}, "org_id")
	})

	p.AddResourceConfigurator("mongodbatlas_org_invitation", func(r *config.Resource) {
		r.ShortGroup = "org"
		r.Kind = "Invitation"
		r.ExternalName = importJoinedIDAssigned([]string{refs.OrgID, "username"}, "invitation_id")
		r.TerraformResource.DeprecationMessage = "This resource is deprecated. Migrate to mongodbatlas_cloud_user_org_assignment for managing organization membership."
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_org_delegation_settings", func(r *config.Resource) {
		r.ShortGroup = "org"
		r.Kind = "DelegationSettings"
		r.ExternalName = templated("{{ .parameters.org_id }}")
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_org_maintenance_settings", func(r *config.Resource) {
		r.ShortGroup = "org"
		r.Kind = "MaintenanceSettings"
		r.ExternalName = templated("{{ .parameters.org_id }}")
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_org_log_integration", func(r *config.Resource) {
		r.ShortGroup = "org"
		r.Kind = "LogIntegration"
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
		// Same as mongodbatlas_log_integration: the generator only supports
		// string-like sensitive fields. The nested value stays sensitive.
		if sch, ok := r.TerraformResource.Schema["otel_supplied_headers"]; ok {
			sch.Sensitive = false
		}
		r.ExternalName = computedKeyID("integration_id")
		r.TerraformPluginFrameworkIsStateEmptyFn = refs.StateEmptyWhenAttributeUnset("integration_id")
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
			refs.ClientID: {
				TerraformName: refs.TFServiceAccount,
				Extractor:     refs.ExtractParamPath(refs.ClientID, true),
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
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromIDOrState("/", 1, 0, refs.ClientID)
		r.ExternalName.SetIdentifierArgumentFn = refs.SetIdentifierArgument(refs.ClientID)
		r.TerraformPluginFrameworkIsStateEmptyFn = refs.StateEmptyWhenAttributeUnset(refs.ClientID)
	})

	p.AddResourceConfigurator("mongodbatlas_service_account_secret", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "ServiceAccountSecret"
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
			refs.ClientID: {
				TerraformName: refs.TFServiceAccount,
				Extractor:     refs.ExtractParamPath(refs.ClientID, true),
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
			refs.ClientID: {
				TerraformName: refs.TFServiceAccount,
				Extractor:     refs.ExtractParamPath(refs.ClientID, true),
			},
		}
		r.ExternalName.GetIDFn = refs.AccessListGetIDFn(refs.OrgID, "client_id")
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromAccessListState(refs.OrgID)
	})

	p.AddResourceConfigurator(refs.TFMCPConfig, func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "MCPConfig"
		// Atlas returns both ip_address and cidr_block for an entry, but
		// Terraform accepts exactly one of them per entry.
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"ip_access_list"},
		}
		r.ExternalName = computedKeyID(refs.MCPConfigID)
		r.TerraformPluginFrameworkIsStateEmptyFn = refs.StateEmptyWhenAttributeUnset(refs.MCPConfigID)
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_mcp_config_secret", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "MCPConfigSecret"
		r.ExternalName = computedKeyID("secret_id")
		r.TerraformPluginFrameworkIsStateEmptyFn = refs.StateEmptyWhenAttributeUnset("secret_id")
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
			refs.MCPConfigID: {
				TerraformName: refs.TFMCPConfig,
				Extractor:     refs.ExtractParamPath(refs.MCPConfigID, true),
			},
		}
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
			refs.IAMRoleID: {
				TerraformName: refs.TFProviderAccessAuthorization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_resource_policy", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "ResourcePolicy"
		r.ExternalName = providerAssignedID()
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_online_archive", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "OnlineArchive"
		r.ExternalName = importJoinedIDAssigned([]string{refs.ProjectID, refs.ClusterName, "archive_id"}, "archive_id")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			refs.ClusterName: {
				TerraformName: refs.TFCluster,
				Extractor:     refs.ExtractParamPath(refs.Name, false),
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_project_api_key", func(r *config.Resource) {
		r.ExternalName = importJoinedIDAssigned([]string{refs.APIKeyID}, refs.APIKeyID)
	})

	p.AddResourceConfigurator("mongodbatlas_access_list_api_key", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "AccessListAPIKey"
		r.ExternalName = accessListImportJoinedID([]string{refs.OrgID, refs.APIKeyID})
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"ip_address"},
		}
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
			refs.APIKeyID: {
				TerraformName: "mongodbatlas_api_key",
			},
		}
	})
}
