package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"
	"regexp"
	"strings"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/xpshim"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/password"
	"github.com/crossplane-contrib/provider-mongodbatlas/config/resources"
)

const (
	resourcePrefix = "mongodbatlas"
	modulePath     = "github.com/crossplane-contrib/provider-mongodbatlas"
)

// SkipTfResourceList resources excluded from code generation.
// - teams: deprecated alias for mongodbatlas_team (SDKv2); upstream recommends mongodbatlas_team.
var SkipTfResourceList = []string{
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

// docMarkers matches the "(Optional)" or "(Required)" markers that can precede
// the separator dash of an argument doc.
var docMarkers = regexp.MustCompile(`^(\s*\([^)]*\))*\s*$`)

// docType matches the type marker tfplugindocs puts before a description, such
// as "(String)" or "(Attributes List, Sensitive)".
var docType = regexp.MustCompile(`^\((String|Number|Boolean|Dynamic|Object|(List|Set|Map) of [^,)]+|(Attributes|Block)( (List|Set|Map))?(, Max: \d+)?)(, (Sensitive|Deprecated))*\)`)

// cleanArgumentDocs fixes the argument docs upjet puts before the schema
// description in each CRD field description.
//
// tfplugindocs writes a doc with a type marker from the schema description,
// so the marker is dropped, and so is the text when every field the doc
// applies to already prints that description. Only a deprecation, which the
// schema description does not mention, is kept.
//
// upjet drops everything up to the first "-" of a doc, expecting the
// "`name` - description" layout. The Atlas docs have no separator, so docs
// would lose the text before a hyphen ("Unique 24-hexadecimal" became
// "hexadecimal"); they get a leading "- " for upjet to drop.
func cleanArgumentDocs() ujconfig.ResourceOption {
	return func(r *ujconfig.Resource) {
		if r.MetaResource == nil {
			return
		}
		described := map[string]bool{}
		if r.TerraformResource != nil {
			schemaDescribed(r.TerraformResource.Schema, "", described)
		}
		// upjet replaces the top-level id with a schema that has no description.
		described["id"] = false
		for k, v := range r.MetaResource.ArgumentDocs {
			m := docType.FindString(v)
			text := strings.TrimSpace(v[len(m):])
			if m != "" && allDescribed(described, k) {
				text = ""
			}
			if strings.Contains(m, "Deprecated") {
				text = strings.TrimSpace("Deprecated. " + text)
			}
			if i := strings.Index(text, "-"); i > 0 && !docMarkers.MatchString(text[:i]) {
				text = "- " + text
			}
			r.MetaResource.ArgumentDocs[k] = text
		}
	}
}

// schemaDescribed records, for every field of s under prefix, whether it has a
// description.
func schemaDescribed(s map[string]*schema.Schema, prefix string, out map[string]bool) {
	for n, f := range s {
		out[prefix+n] = f.Description != ""
		if e, ok := f.Elem.(*schema.Resource); ok {
			schemaDescribed(e.Schema, prefix+n+".", out)
		}
	}
}

// allDescribed reports whether every field the doc key can apply to has a
// description. upjet applies a doc to the field with that path and, failing
// that, to fields whose path ends with it.
func allDescribed(described map[string]bool, key string) bool {
	found := false
	for path, ok := range described {
		if path == key || strings.HasSuffix(path, "."+key) {
			if !ok {
				return false
			}
			found = true
		}
	}
	return found
}

func newProvider(rootGroup string, pwGen func(string, string) ujconfig.NewInitializerFn) *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithDefaultResourceOptions(
			resetRootShortGroup(),
			cleanArgumentDocs(),
			ExternalNameConfigurations(),
		),
		ujconfig.WithExampleManifestConfiguration(ujconfig.ExampleManifestConfiguration{
			ManagedResourceNamespace: "crossplane-system",
		}),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithRootGroup(rootGroup),
		ujconfig.WithShortName("mongodbatlas"),
		ujconfig.WithSkipList(SkipTfResourceList),
		ujconfig.WithTerraformProvider(xpshim.GetSDKProvider()),
		ujconfig.WithTerraformPluginFrameworkProvider(xpshim.GetFrameworkProvider()),
		// IncludeList defaults to [".+"] (all resources).
		// Clear it to avoid overlap with the typed lists below.
		ujconfig.WithIncludeList(nil),
		ujconfig.WithTerraformPluginSDKIncludeList(TerraformSDKIncludeList()),
		ujconfig.WithTerraformPluginFrameworkIncludeList(TerraformFrameworkIncludeList()),
	)

	resources.ConfigureAI(pc)
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
