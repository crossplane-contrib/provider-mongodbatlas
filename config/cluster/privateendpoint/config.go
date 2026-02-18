package privateendpoint

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_private_endpoint_regional_mode", func(r *config.Resource) {
		r.ShortGroup = "privateendpoint"
		r.Kind = "RegionalMode"

		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_privatelink_endpoint", func(r *config.Resource) {
		r.ShortGroup = "privateendpoint"
		r.Kind = "Resource"

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
			provider, ok := parameters["provider_name"]
			if !ok {
				return "", errors.New("provider_name missing from parameters")
			}
			region, ok := parameters["region"]
			if !ok {
				return "", errors.New("region missing from parameters")
			}
			return fmt.Sprintf("%s-%s-%s-%s", project, externalName, provider, region), nil
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

	p.AddResourceConfigurator("mongodbatlas_privatelink_endpoint_service", func(r *config.Resource) {
		r.ShortGroup = "privateendpoint"
		r.Kind = "Service"

		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
			"private_link_id": {
				TerraformName: "mongodbatlas_privatelink_endpoint",
			},
		}
		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			project, ok := parameters["project_id"]
			if !ok {
				return "", errors.New("project_id missing from parameters")
			}
			link, ok := parameters["private_link_id"]
			if !ok {
				return "", errors.New("private_link_id missing from parameters")
			}
			provider, ok := parameters["provider_name"]
			if !ok {
				return "", errors.New("provider_name missing from parameters")
			}
			return fmt.Sprintf("%s--%s--%s--%s", project, link, externalName, provider), nil
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

			idSlice := strings.Split(idStr, "--")
			return idSlice[2], nil
		}
	})
}
