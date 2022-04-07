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

package config

import (
	tjconfig "github.com/crossplane/terrajet/pkg/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/crossplane-contrib/provider-jet-mongodbatlas/config/database"
	"github.com/crossplane-contrib/provider-jet-mongodbatlas/config/mongodbatlas"
	"github.com/crossplane-contrib/provider-jet-mongodbatlas/config/project"
)

const (
	resourcePrefix = "mongodbatlas"
	modulePath     = "github.com/crossplane-contrib/provider-jet-mongodbatlas"
)

// GetProvider returns provider configuration
func GetProvider(resourceMap map[string]*schema.Resource) *tjconfig.Provider {
	pc := tjconfig.NewProvider(resourceMap, resourcePrefix, modulePath,
		tjconfig.WithDefaultResourceFn(DefaultResource(
			// Note(turkenh): Some other resource configuration options rely on
			// the final version Group and Kind. So, please make sure to have
			// `groupKindOverrides()` as the first option here.
			gvkOverrides(),
			identifierAssignedByMongoDBAtlas(),
			commonReferences(),
		)),
		tjconfig.WithSkipList([]string{"mongodbatlas_encryption_at_rest", "mongodbatlas_teams"}))

	for _, configure := range []func(provider *tjconfig.Provider){
		// add custom config functions
		mongodbatlas.Configure,
		project.Configure,
		database.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}

// DefaultResource returns a DefaultResoruceFn that makes sure the original
// DefaultResource call is made with given options here.
func DefaultResource(opts ...tjconfig.ResourceOption) tjconfig.DefaultResourceFn {
	return func(name string, terraformResource *schema.Resource, orgOpts ...tjconfig.ResourceOption) *tjconfig.Resource {
		return tjconfig.DefaultResource(name, terraformResource, append(orgOpts, opts...)...)
	}
}
