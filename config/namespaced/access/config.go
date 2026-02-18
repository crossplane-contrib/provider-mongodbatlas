package access

import (
	"context"
	"errors"
	"fmt"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_access_list_api_key", func(r *config.Resource) {
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
