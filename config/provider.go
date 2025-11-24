/*
Copyright 2021 Upbound Inc.
*/

package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	databaseCluster "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/database"
	mongodbatlasCluster "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/mongodbatlas"
	projectCluster "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/project"

	databaseNamespaced "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/database"
	mongodbatlasNamespaced "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/mongodbatlas"
	projectNamespaced "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/project"
)

const (
	resourcePrefix = "mongodbatlas"
	modulePath     = "github.com/crossplane-contrib/provider-mongodbatlas"
)

var SkipTfResourceList = []string{
	"mongodbatlas_encryption_at_rest",
	"mongodbatlas_teams",
}

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithSkipList(SkipTfResourceList),
		// ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithShortName("mongodbatlas"),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithRootGroup("mongodbatlas.crossplane.io"),
		ujconfig.WithDefaultResourceOptions(
			clusterGvkOverride(),
			identifierAssignedByMongoDBAtlas(),
			clusterCommonReferencesOverride(),
		))

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		databaseCluster.Configure,
		mongodbatlasCluster.Configure,
		projectCluster.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}

// GetProviderNamespaced returns provider configuration
func GetProviderNamespaced() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithSkipList(SkipTfResourceList),
		// ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithShortName("mongodbatlas"),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithRootGroup("mongodbatlas.m.crossplane.io"),
		ujconfig.WithDefaultResourceOptions(
			namespacedGvkOverride(),
			identifierAssignedByMongoDBAtlas(),
			namespacedCommonReferencesOverride(),
		))

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		databaseNamespaced.Configure,
		mongodbatlasNamespaced.Configure,
		projectNamespaced.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}
