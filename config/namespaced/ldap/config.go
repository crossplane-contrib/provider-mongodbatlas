package ldap

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_ldap_configuration", func(r *config.Resource) {
		r.ShortGroup = "ldap"
		r.Kind = "Configuration"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_ldap_verify", func(r *config.Resource) {
		r.ShortGroup = "ldap"
		r.Kind = "Verify"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			project, ok := parameters["project_id"]
			if !ok {
				return "", errors.New("org_id missing from parameters")
			}
			return fmt.Sprintf("%s-%s", project, externalName), nil
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
