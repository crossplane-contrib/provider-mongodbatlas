package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	alertCluster "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/alert"
	cloudCluster "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/cloud"
	databaseCluster "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/database"
	federatedCluster "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/federated"
	ldapCluster "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/ldap"
	mongodbatlasCluster "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/mongodbatlas"
	networkCluster "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/network"
	privateCluster "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/privateendpoint"
	projectCluster "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/project"
	searchCluster "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/search"
	streamCluster "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/stream"

	alertNamespaced "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/alert"
	cloudnamespaced "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/cloud"
	databaseNamespaced "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/database"
	federatedNamespaced "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/federated"
	ldapNamespaced "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/ldap"
	mongodbatlasNamespaced "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/mongodbatlas"
	networkNamespaced "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/network"
	privateNamespaced "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/privateendpoint"
	projectNamespaced "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/project"
	searchNamespaced "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/search"
	streamNamespaced "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/stream"
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
		ujconfig.WithDefaultResourceOptions(
			clusterGvkOverride(),
			identifierAssignedByMongoDBAtlas(),
			clusterCommonReferencesOverride(),
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

	for _, configure := range []func(provider *ujconfig.Provider){
		alertCluster.Configure,
		cloudCluster.Configure,
		databaseCluster.Configure,
		federatedCluster.Configure,
		ldapCluster.Configure,
		mongodbatlasCluster.Configure,
		networkCluster.Configure,
		privateCluster.Configure,
		projectCluster.Configure,
		searchCluster.Configure,
		streamCluster.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}

// GetProviderNamespaced returns provider configuration
func GetProviderNamespaced() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithDefaultResourceOptions(
			namespacedGvkOverride(),
			identifierAssignedByMongoDBAtlas(),
			namespacedCommonReferencesOverride(),
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

	for _, configure := range []func(provider *ujconfig.Provider){
		alertNamespaced.Configure,
		cloudnamespaced.Configure,
		databaseNamespaced.Configure,
		federatedNamespaced.Configure,
		ldapNamespaced.Configure,
		mongodbatlasNamespaced.Configure,
		networkNamespaced.Configure,
		privateNamespaced.Configure,
		projectNamespaced.Configure,
		searchNamespaced.Configure,
		streamNamespaced.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}
