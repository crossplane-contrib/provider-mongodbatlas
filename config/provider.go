package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/password"
	"github.com/crossplane-contrib/provider-mongodbatlas/config/resources"
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

// resetRootShortGroup clears the ShortGroup for resources in the root API group.
// upjet defaults ShortGroup to the resource prefix; resources without a sub-group
// must have it cleared so they land in the root mongodbatlas.crossplane.io group.
func resetRootShortGroup() ujconfig.ResourceOption {
	return func(r *ujconfig.Resource) {
		if r.ShortGroup == resourcePrefix {
			r.ShortGroup = ""
		}
	}
}

// GetidentifierFromProvider returns provider configuration
func GetidentifierFromProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithDefaultResourceOptions(
			resetRootShortGroup(),
			ExternalNameConfigurations(),
		),
		ujconfig.WithExampleManifestConfiguration(ujconfig.ExampleManifestConfiguration{
			ManagedResourceNamespace: "crossplane-system",
		}),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithRootGroup("mongodbatlas.crossplane.io"),
		ujconfig.WithShortName("mongodbatlas"),
		ujconfig.WithSkipList(SkipTfResourceList),
	)

	resources.ConfigureAlert(pc)
	resources.ConfigureCloud(pc)
	resources.ConfigureDatabase(pc, password.ClusterGenerator)
	resources.ConfigureFederated(pc)
	resources.ConfigureLDAP(pc, password.ClusterGenerator)
	resources.ConfigureMongoDBAtlas(pc)
	resources.ConfigureNetwork(pc)
	resources.ConfigurePrivateEndpoint(pc)
	resources.ConfigureProject(pc)
	resources.ConfigureSearch(pc)
	resources.ConfigureStream(pc, password.ClusterGenerator)

	pc.ConfigureResources()
	return pc
}

// GetProviderNamespaced returns provider configuration
func GetProviderNamespaced() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithDefaultResourceOptions(
			resetRootShortGroup(),
			ExternalNameConfigurations(),
		),
		ujconfig.WithExampleManifestConfiguration(ujconfig.ExampleManifestConfiguration{
			ManagedResourceNamespace: "crossplane-system",
		}),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithRootGroup("mongodbatlas.m.crossplane.io"),
		ujconfig.WithShortName("mongodbatlas"),
		ujconfig.WithSkipList(SkipTfResourceList),
	)

	resources.ConfigureAlert(pc)
	resources.ConfigureCloud(pc)
	resources.ConfigureDatabase(pc, password.NamespacedGenerator)
	resources.ConfigureFederated(pc)
	resources.ConfigureLDAP(pc, password.NamespacedGenerator)
	resources.ConfigureMongoDBAtlas(pc)
	resources.ConfigureNetwork(pc)
	resources.ConfigurePrivateEndpoint(pc)
	resources.ConfigureProject(pc)
	resources.ConfigureSearch(pc)
	resources.ConfigureStream(pc, password.NamespacedGenerator)

	pc.ConfigureResources()
	return pc
}
