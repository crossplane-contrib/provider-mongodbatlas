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

	"github.com/crossplane/upjet/v2/pkg/config"

	common "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/common"
)

const group = "project"

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
		r.ShortGroup = ""
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_project_invitation", func(r *config.Resource) {
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"ip_address"},
		}
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
		r.ShortGroup = group
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
		r.ShortGroup = group
		r.Kind = "ServiceAccountSecret"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("/", 2, "project_id", "client_id")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("/", 2, 0)
	})

	p.AddResourceConfigurator("mongodbatlas_project_service_account", func(r *config.Resource) {
		r.ShortGroup = group
		r.Kind = "ServiceAccount"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}

		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("/", 1, "project_id")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("/", 1, 0)
	})
}
