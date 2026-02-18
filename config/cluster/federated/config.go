package federated

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_federated_database_instance", func(r *config.Resource) {
		r.ShortGroup = "federation"
		r.Kind = "DatabaseInstance"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_federated_query_limit", func(r *config.Resource) {
		r.ShortGroup = "federation"
		r.Kind = "QueryLimit"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_privatelink_endpoint_service_data_federation_online_archive", func(r *config.Resource) {
		r.ShortGroup = "federation"
		r.Kind = "PrivateLinkEndpointService"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_federated_settings_identity_provider", func(r *config.Resource) {
		r.ShortGroup = "federation"
		r.Kind = "IdentityProvider"
		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			f, ok := parameters["federation_settings_id"]
			if !ok {
				return "", errors.New("federation_settings_id missing from parameters")
			}
			return fmt.Sprintf("%s-%s", f, externalName), nil
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

	p.AddResourceConfigurator("mongodbatlas_federated_settings_org_config", func(r *config.Resource) {
		r.ShortGroup = "federation"
		r.Kind = "OrgConfigSettings"
		r.References = config.References{
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_federated_settings_org_role_mapping", func(r *config.Resource) {
		r.ShortGroup = "federation"
		r.Kind = "RoleMapping"
		r.References = config.References{
			"federation_settings_id": {
				TerraformName: "mongodbatlas_organization",
			},
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
		}
		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			f, ok := parameters["federation_settings_id"]
			if !ok {
				return "", errors.New("federation_settings_id missing from parameters")
			}
			o, ok := parameters["org_id"]
			if !ok {
				return "", errors.New("org_id missing from parameters")
			}
			return fmt.Sprintf("%s-%s-%s", f, o, externalName), nil
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
}
