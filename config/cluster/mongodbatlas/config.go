package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_cluster", func(r *config.Resource) {
		r.UseAsync = true
		r.Kind = "Cluster"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_advanced_cluster", func(r *config.Resource) {
		r.ShortGroup = "mongodbatlas"
		r.Kind = "AdvancedCluster"
		r.UseAsync = true
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_flex_cluster", func(r *config.Resource) {
		r.UseAsync = true
		r.ShortGroup = "mongodbatlas"
		r.Kind = "FlexCluster"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cluster_outage_simulation", func(r *config.Resource) {
		r.ShortGroup = "mongodbatlas"
		r.Kind = "ClusterOutageSimulation"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_mongodb_employee_access_grant", func(r *config.Resource) {
		r.ShortGroup = "mongodbatlas"
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
		r.ShortGroup = "mongodbatlas"
		r.Kind = "Team"
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
			return fmt.Sprintf("%s-%s", org, externalName), nil
		}

		r.ExternalName.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
			id, ok := tfstate["id"]
			if !ok {
				return "", errors.New("id attribute missing from state file")
			}

			idStr, ok := id.(string)
			if !ok {
				return "", errors.New("value of id needs to be string")
			}

			idSlice := strings.Split(idStr, "-")
			return idSlice[1], nil
		}
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
		r.ShortGroup = "mongodbatlas"
		r.Kind = "APIKey"
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
			return fmt.Sprintf("%s-%s", org, externalName), nil
		}

		r.ExternalName.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
			id, ok := tfstate["id"]
			if !ok {
				return "", errors.New("id attribute missing from state file")
			}

			idStr, ok := id.(string)
			if !ok {
				return "", errors.New("value of id needs to be string")
			}

			idSlice := strings.Split(idStr, "-")
			return idSlice[1], nil
		}
	})

	p.AddResourceConfigurator("mongodbatlas_api_key_project_assignment", func(r *config.Resource) {
		r.ShortGroup = "mongodbatlas"
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

	p.AddResourceConfigurator("mongodbatlas_event_trigger", func(r *config.Resource) {
		r.ShortGroup = "mongodbatlas"
		r.Kind = "EventTrigger"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}

		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			project, ok := parameters["project_id"]
			if !ok {
				return "", errors.New("project_id missing from parameters")
			}
			app, ok := parameters["app_id"]
			if !ok {
				return "", errors.New("app_id missing from parameters")
			}
			return fmt.Sprintf("%s-%s-%s", project, app, externalName), nil
		}

		r.ExternalName.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
			id, ok := tfstate["id"]
			if !ok {
				return "", errors.New("id attribute missing from state file")
			}

			idStr, ok := id.(string)
			if !ok {
				return "", errors.New("value of id needs to be string")
			}

			idSlice := strings.Split(idStr, "-")
			return idSlice[2], nil
		}
	})

	p.AddResourceConfigurator("mongodbatlas_log_integration", func(r *config.Resource) {
		r.ShortGroup = "mongodbatlas"
		r.Kind = "LogIntegration"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}

		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			project, ok := parameters["project_id"]
			if !ok {
				return "", errors.New("project_id missing from parameters")
			}
			return fmt.Sprintf("%s/%s", project, externalName), nil
		}

		r.ExternalName.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
			id, ok := tfstate["id"]
			if !ok {
				return "", errors.New("id attribute missing from state file")
			}

			idStr, ok := id.(string)
			if !ok {
				return "", errors.New("value of id needs to be string")
			}

			idSlice := strings.Split(idStr, "/")
			return idSlice[1], nil
		}
	})

	p.AddResourceConfigurator("mongodbatlas_custom_dns_configuration_cluster_aws", func(r *config.Resource) {
		r.ShortGroup = "mongodbatlas"
		r.Kind = "CustomDNSConfigurationClusterAWS"

		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_maintenance_window", func(r *config.Resource) {
		r.ShortGroup = "mongodbatlas"
		r.Kind = "MaintenanceWindow"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_service_account_project_assignment", func(r *config.Resource) {
		r.ShortGroup = "mongodbatlas"
		r.Kind = "ServiceAccountProjectAssignment"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_service_account", func(r *config.Resource) {
		r.ShortGroup = "mongodbatlas"
		r.Kind = "ServiceAccount"
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
			return fmt.Sprintf("%s/%s", org, externalName), nil
		}
		r.ExternalName.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
			id, ok := tfstate["id"]
			if !ok {
				return "", errors.New("id attribute missing from state file")
			}

			idStr, ok := id.(string)
			if !ok {
				return "", errors.New("value of id needs to be string")
			}

			idSlice := strings.Split(idStr, "/")
			return idSlice[1], nil
		}
	})

	p.AddResourceConfigurator("mongodbatlas_service_account_secret", func(r *config.Resource) {
		r.ShortGroup = "mongodbatlas"
		r.Kind = "ServiceAccountSecret"
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

			return fmt.Sprintf("%s/%s/%s", org, client, externalName), nil
		}
		r.ExternalName.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
			id, ok := tfstate["id"]
			if !ok {
				return "", errors.New("id attribute missing from state file")
			}

			idStr, ok := id.(string)
			if !ok {
				return "", errors.New("value of id needs to be string")
			}

			idSlice := strings.Split(idStr, "/")
			return idSlice[2], nil
		}
	})

	p.AddResourceConfigurator("mongodbatlas_service_account_access_list_entry", func(r *config.Resource) {
		r.ShortGroup = "mongodbatlas"
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
		r.ShortGroup = "mongodbatlas"
		r.Kind = "PushBasedLogExport"
		r.TerraformResource.DeprecationMessage = "This resource is deprecated and will be removed in the next major version. Please use mongodbatlas_log_integration instead."
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_resource_policy", func(r *config.Resource) {
		r.ShortGroup = "mongodbatlas"
		r.Kind = "ResourcePolicy"
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
			return fmt.Sprintf("%s-%s", org, externalName), nil
		}

		r.ExternalName.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
			id, ok := tfstate["id"]
			if !ok {
				return "", errors.New("id attribute missing from state file")
			}

			idStr, ok := id.(string)
			if !ok {
				return "", errors.New("value of id needs to be string")
			}

			idSlice := strings.Split(idStr, "-")
			return idSlice[1], nil
		}
	})

	p.AddResourceConfigurator("mongodbatlas_online_archive", func(r *config.Resource) {
		r.ShortGroup = "mongodbatlas"
		r.Kind = "OnlineArchive"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}

		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			project, ok := parameters["project_id"]
			if !ok {
				return "", errors.New("project_id missing from parameters")
			}
			cluster, ok := parameters["cluster_name"]
			if !ok {
				return "", errors.New("cluster_name missing from parameters")
			}

			return fmt.Sprintf("%s-%s-%s", project, cluster, externalName), nil
		}

		r.ExternalName.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
			id, ok := tfstate["id"]
			if !ok {
				return "", errors.New("id attribute missing from state file")
			}

			idStr, ok := id.(string)
			if !ok {
				return "", errors.New("value of id needs to be string")
			}

			idSlice := strings.Split(idStr, "-")
			return idSlice[2], nil
		}
	})
}
