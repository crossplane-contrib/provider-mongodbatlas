package config

import (
	ujconfig "github.com/crossplane/upjet/pkg/config"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/common"
)

var gvkMap = map[string]schema.GroupVersionKind{
	"mongodbatlas_advanced_cluster": {
		Group:   "",
		Kind:    "AdvancedCluster",
		Version: common.VersionV1Alpha2,
	},
	"mongodbatlas_cluster": {
		Group:   "",
		Kind:    "Cluster",
		Version: common.VersionV1Alpha2,
	},
}

// gvkOverrides overrides the group, version and kind of the resource if it matches
// any entry in the gvkMap.
func gvkOverrides() ujconfig.ResourceOption {
	return func(r *ujconfig.Resource) {
		if r.ShortGroup == resourcePrefix {
			r.ShortGroup = ""
		}
		if gvk, ok := gvkMap[r.Name]; ok {
			r.ShortGroup = gvk.Group
			r.Kind = gvk.Kind
			r.Version = gvk.Version
		}
	}
}

// identifierAssignedByMongoDBAtlas is the most common external name
// configuration.
// The resource-specific configurations should override this whenever needed.
func identifierAssignedByMongoDBAtlas() ujconfig.ResourceOption {
	return func(r *ujconfig.Resource) {
		r.ExternalName = ujconfig.IdentifierFromProvider
	}
}

// commonReferences adds referencers for fields that are known and common among
// more than a few resources.
func commonReferences() ujconfig.ResourceOption {
	return func(r *ujconfig.Resource) {
		for k, s := range r.TerraformResource.Schema {
			// We shouldn't add referencers for status fields and sensitive fields
			// since they already have secret referencer.
			if (s.Computed && !s.Optional) || s.Sensitive {
				continue
			}

			if k == "project_id" {
				ref := ujconfig.Reference{
					Type:      common.APISPackagePath + "/mongodbatlas/v1alpha1.Project",
					Extractor: common.ExtractResourceIDFuncPath,
				}
				if r.ShortGroup == "" && r.Version == "v1alpha1" {
					ref.Type = "Project"
				}
				r.References["project_id"] = ref
			}
		}
	}
}
