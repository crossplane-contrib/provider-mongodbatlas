package mongodbatlas

import (
	"context"
	"errors"
	"fmt"

	"github.com/crossplane/upjet/v2/pkg/config"

	common "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/common"
)

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_cluster", func(r *config.Resource) {
		r.Version = common.VersionV1Alpha2
		r.UseAsync = true
		r.Kind = "Cluster"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_advanced_cluster", func(r *config.Resource) {
		r.Version = common.VersionV1Alpha3
		r.ShortGroup = ""
		r.Kind = "AdvancedCluster"
		r.UseAsync = true
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_flex_cluster", func(r *config.Resource) {
		r.ShortGroup = ""
		r.UseAsync = true
		r.Kind = "FlexCluster"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cluster_outage_simulation", func(r *config.Resource) {
		r.ShortGroup = "cluster"
		r.Kind = "OutageSimulation"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_mongodb_employee_access_grant", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "EmployeeAccessGrant"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_auditing", func(r *config.Resource) {
		r.UseAsync = true
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_team", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "Team"
		r.References = config.References{
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
		}

		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 1, "org_id")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("-", 1, 0)
	})

	p.AddResourceConfigurator("mongodbatlas_team_project_assignment", func(r *config.Resource) {
		r.ShortGroup = "team"
		r.Kind = "ProjectAssignment"
		r.References = config.References{
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_api_key", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "APIKey"
		r.References = config.References{
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
		}

		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 1, "org_id")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("-", 1, 0)
	})

	p.AddResourceConfigurator("mongodbatlas_api_key_project_assignment", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "APIKeyProjectAssignment"
		r.References = config.References{
			"api_key_id": {
				TerraformName: "mongodbatlas_api_key",
			},
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_encryption_at_rest", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "EncryptionAtRest"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_encryption_at_rest_private_endpoint", func(r *config.Resource) {
		r.ShortGroup = "encryptionatrest"
		r.Kind = "PrivateEndpoint"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 2, "project_id", "cloud_provider")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("-", 2, 0)
	})

	p.AddResourceConfigurator("mongodbatlas_event_trigger", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "EventTrigger"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}

		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 2, "project_id", "app_id")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("-", 2, 0)
	})

	// Configure configures the root group
	p.AddResourceConfigurator("mongodbatlas_global_cluster_config", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "GlobalClusterConfig"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_log_integration", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "LogIntegration"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}

		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("/", 1, "project_id")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("/", 1, 0)
	})

	p.AddResourceConfigurator("mongodbatlas_organization", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "Organization"
	})

	p.AddResourceConfigurator("mongodbatlas_org_invitation", func(r *config.Resource) {
		r.ShortGroup = "org"
		r.Kind = "Invitation"
		r.TerraformResource.DeprecationMessage = "This resource is deprecated. Migrate to mongodbatlas_cloud_user_org_assignment for managing organization membership."

		r.References = config.References{
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_custom_dns_configuration_cluster_aws", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "CustomDNSConfigurationClusterAWS"

		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_maintenance_window", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "MaintenanceWindow"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_service_account_project_assignment", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "ServiceAccountProjectAssignment"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_service_account", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "ServiceAccount"
		r.References = config.References{
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
		}

		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("/", 1, "org_id")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("/", 1, 0)
	})

	p.AddResourceConfigurator("mongodbatlas_service_account_secret", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "ServiceAccountSecret"
		r.References = config.References{
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
		}

		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("/", 2, "org_id", "client_id")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("/", 2, 0)
	})

	p.AddResourceConfigurator("mongodbatlas_service_account_access_list_entry", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "ServiceAccountAccessListEntry"
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"ip_address"},
		}
		r.References = config.References{
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
		}

		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			org, ok := parameters["org_id"]
			if !ok {
				return "", errors.New("org_id missing from parameters")
			}
			client, ok := parameters["client_id"]
			if !ok {
				return "", errors.New("client_id missing from parameters")
			}
			ip, ok := parameters["ip_address"]
			if !ok {
				ip, ok = parameters["cidr_block"]
				if !ok {
					return "", errors.New("either ip_address or cidr_block parameters must be set")
				}
			}
			return fmt.Sprintf("%s-%s-%s", org, client, ip), nil
		}
	})

	p.AddResourceConfigurator("mongodbatlas_push_based_log_export", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "PushBasedLogExport"
		r.TerraformResource.DeprecationMessage = "This resource is deprecated and will be removed in the next major version. Please use mongodbatlas_log_integration instead."
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_resource_policy", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "ResourcePolicy"
		r.References = config.References{
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
		}
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 1, "org_id")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("-", 1, 0)
	})

	p.AddResourceConfigurator("mongodbatlas_online_archive", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "OnlineArchive"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}

		// ID format: {project_id}-{cluster_name}-{archive_id}
		// cluster_name may contain dashes, but the extracted archive_id is a hex
		// ID that never contains dashes, so Split+last is safe here.
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 2, "project_id", "cluster_name")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromSegment("-")
	})

	p.AddResourceConfigurator("mongodbatlas_access_list_api_key", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "AccessListAPIKey"
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"ip_address"},
		}
		r.References = config.References{
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
			"api_key_id": {
				TerraformName: "mongodbatlas_api_key",
			},
		}

		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			org, ok := parameters["org_id"]
			if !ok {
				return "", errors.New("org_id missing from parameters")
			}
			api_key, ok := parameters["api_key_id"]
			if !ok {
				return "", errors.New("api_key_id missing from parameters")
			}
			ip, ok := parameters["ip_address"]
			if !ok {
				ip, ok = parameters["cidr_block"]
				if !ok {
					return "", errors.New("either ip_address or cidr_block parameters must be set")
				}
			}
			return fmt.Sprintf("%s-%s-%s", org, api_key, ip), nil
		}
	})
}
