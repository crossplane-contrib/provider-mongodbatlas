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

// SkipTfResourceList resources excluded from code generation.
// - encryption_at_rest: historically broken under CLI mode (state drift on key fields); re-evaluate under no-fork.
// - teams: deprecated alias for mongodbatlas_team (SDKv2); upstream recommends mongodbatlas_team.
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

func newProvider(rootGroup string, pwGen func(string, string) ujconfig.NewInitializerFn) *ujconfig.Provider {
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
		ujconfig.WithRootGroup(rootGroup),
		ujconfig.WithShortName("mongodbatlas"),
		ujconfig.WithSkipList(SkipTfResourceList),
	)

	resources.ConfigureAlert(pc)
	resources.ConfigureCloud(pc)
	resources.ConfigureDatabase(pc, pwGen)
	resources.ConfigureFederated(pc)
	resources.ConfigureLDAP(pc, pwGen)
	resources.ConfigureMongoDBAtlas(pc)
	resources.ConfigureNetwork(pc)
	resources.ConfigurePrivateEndpoint(pc)
	resources.ConfigureProject(pc)
	resources.ConfigureSearch(pc)
	resources.ConfigureStream(pc, pwGen)

	pc.ConfigureResources()
	return pc
}

// GetidentifierFromProvider returns the cluster-scoped provider configuration.
func GetidentifierFromProvider() *ujconfig.Provider {
	return newProvider("mongodbatlas.crossplane.io", password.ClusterGenerator)
}

// GetProviderNamespaced returns the namespace-scoped provider configuration.
func GetProviderNamespaced() *ujconfig.Provider {
	return newProvider("mongodbatlas.m.crossplane.io", password.NamespacedGenerator)
}
