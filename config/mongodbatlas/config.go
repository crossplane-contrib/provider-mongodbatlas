/*
Copyright 2021 The Crossplane Authors.

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

package mongodbatlas

import (
	"github.com/crossplane/terrajet/pkg/config"

	"github.com/crossplane-contrib/provider-jet-mongodbatlas/config/common"
)

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_project", func(r *config.Resource) {
		r.ExternalName = config.IdentifierFromProvider
	})
	p.AddResourceConfigurator("mongodbatlas_cluster", func(r *config.Resource) {
		r.ExternalName = config.IdentifierFromProvider
		r.References = config.References{
			"project_id": config.Reference{
				Type:      "Project",
				Extractor: common.ExtractResourceIDFuncPath,
			},
		}
		r.UseAsync = true
	})
	p.AddResourceConfigurator("mongodbatlas_advanced_cluster", func(r *config.Resource) {
		r.ShortGroup = ""
		r.Kind = "AdvancedCluster"
		r.ExternalName = config.IdentifierFromProvider
		r.References = config.References{
			"project_id": config.Reference{
				Type:      "Project",
				Extractor: common.ExtractResourceIDFuncPath,
			},
		}
		r.UseAsync = true
	})
	p.AddResourceConfigurator("mongodbatlas_project_ip_access_list", func(r *config.Resource) {
		r.ExternalName = config.IdentifierFromProvider
		r.References = config.References{
			"project_id": config.Reference{
				Type:      "Project",
				Extractor: common.ExtractResourceIDFuncPath,
			},
		}
	})
}
