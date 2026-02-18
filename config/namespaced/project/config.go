/*
Copyright 2022 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package project

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_project", func(r *config.Resource) {
		r.References = config.References{
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_third_party_integration", func(r *config.Resource) {
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_project_ip_access_list", func(r *config.Resource) {
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"ip_address"},
		}
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
			ip, ok := parameters["ip_address"]
			if !ok {
				ip, ok = parameters["cidr_block"]
				if !ok {
					return "", errors.New("either ip_address or cidr_block parameters must be set")
				}
			}
			return fmt.Sprintf("%s-%s", project, ip), nil
		}
	})

	p.AddResourceConfigurator("mongodbatlas_third_party_integration", func(r *config.Resource) {
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_project_service_account_access_list_entry", func(r *config.Resource) {
		r.ShortGroup = "project"
		r.Kind = "ServiceAccountAccessListEntry"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}

		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			client, ok := parameters["client_id"]
			if !ok {
				return "", errors.New("client_id missing from parameters")
			}
			project, ok := parameters["project_id"]
			if !ok {
				return "", errors.New("project_id missing from parameters")
			}
			ip, ok := parameters["ip_address"]
			if !ok {
				ip, ok = parameters["cidr_block"]
				if !ok {
					return "", errors.New("either ip_address or cidr_block parameters must be set")
				}
			}
			return fmt.Sprintf("%s-%s-%s", project, client, ip), nil
		}
	})

	p.AddResourceConfigurator("mongodbatlas_project_service_account_secret", func(r *config.Resource) {
		r.ShortGroup = "project"
		r.Kind = "ServiceAccountSecret"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			client, ok := parameters["client_id"]
			if !ok {
				return "", errors.New("client_id missing from parameters")
			}
			project, ok := parameters["project_id"]
			if !ok {
				return "", errors.New("project_id missing from parameters")
			}
			return fmt.Sprintf("%s/%s/%s", project, client, externalName), nil
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

	p.AddResourceConfigurator("mongodbatlas_project_service_account", func(r *config.Resource) {
		r.ShortGroup = "project"
		r.Kind = "ServiceAccount"
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
}
