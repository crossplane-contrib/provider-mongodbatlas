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
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/upbound/upjet/pkg/config"

	"github.com/crossplane-contrib/provider-jet-mongodbatlas/config/database"
	"github.com/crossplane-contrib/provider-jet-mongodbatlas/config/mongodbatlas"
	"github.com/crossplane-contrib/provider-jet-mongodbatlas/config/project"
)

const (
	resourcePrefix = "mongodbatlas"
	modulePath     = "github.com/crossplane-contrib/provider-jet-mongodbatlas"
)

//go:embed schema.json
var providerSchema string

// GetProvider returns provider configuration
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, nil,
		ujconfig.WithSkipList([]string{"mongodbatlas_encryption_at_rest", "mongodbatlas_teams"}),
		ujconfig.WithDefaultResourceOptions(
			gvkOverrides(),
			identifierAssignedByMongoDBAtlas(),
			commonReferences(),
		),
		ujconfig.WithRootGroup("mongodbatlas.jet.crossplane.io"), // keep the old terrajet naming
	)

	for _, configure := range []func(provider *ujconfig.Provider){
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
