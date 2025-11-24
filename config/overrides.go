package config

import (
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	"k8s.io/apimachinery/pkg/runtime/schema"

	commonCluster "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/common"
	commonNamespaced "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/common"
)

var clusterGkvOverrideMap = map[string]schema.GroupVersionKind{
	"mongodbatlas_advanced_cluster": {
		Group:   "",
		Kind:    "AdvancedCluster",
		Version: commonCluster.VersionV1Alpha2,
	},
	"mongodbatlas_cluster": {
		Group:   "",
		Kind:    "Cluster",
		Version: commonCluster.VersionV1Alpha2,
	},
}

var namespacedGkvOverrideMap = map[string]schema.GroupVersionKind{
	"mongodbatlas_advanced_cluster": {
		Group:   "",
		Kind:    "AdvancedCluster",
		Version: commonNamespaced.VersionV1Alpha2,
	},
	"mongodbatlas_cluster": {
		Group:   "",
		Kind:    "Cluster",
		Version: commonNamespaced.VersionV1Alpha2,
	},
}

// clusterGvkOverride overrides the group, version and kind of the resource if it matches
// any entry in the gvkMap.
func clusterGvkOverride() ujconfig.ResourceOption {
	return func(r *ujconfig.Resource) {
		if r.ShortGroup == resourcePrefix {
			r.ShortGroup = ""
		}
		if gvk, ok := clusterGkvOverrideMap[r.Name]; ok {
			r.ShortGroup = gvk.Group
			r.Kind = gvk.Kind
			r.Version = gvk.Version
		}
	}
}

// clusterCommonReferencesOverride adds referencers for fields that are known and common among
// more than a few resources.
func clusterCommonReferencesOverride() ujconfig.ResourceOption {
	return func(r *ujconfig.Resource) {
		for k, s := range r.TerraformResource.Schema {
			// We shouldn't add referencers for status fields and sensitive fields
			// since they already have secret referencer.
			if (s.Computed && !s.Optional) || s.Sensitive {
				continue
			}

			// Most resources reference a "project_id". This is an ad-hoc way of
			// TODO: Set this field only in the config of the resources actually using it.
			if k == "project_id" {
				ref := ujconfig.Reference{
					Type:      commonCluster.APISPackagePath + "/mongodbatlas/v1alpha1.Project",
					Extractor: commonCluster.ExtractResourceIDFuncPath,
				}
				if r.ShortGroup == "" && r.Version == "v1alpha1" {
					ref.Type = "Project"
				}
				r.References["project_id"] = ref
			}
		}
	}
}


// namespacedGvkOverride overrides the group, version and kind of the resource if it matches
// any entry in the gvkMap.
func namespacedGvkOverride() ujconfig.ResourceOption {
	return func(r *ujconfig.Resource) {
		if r.ShortGroup == resourcePrefix {
			r.ShortGroup = ""
		}
		if gvk, ok := namespacedGkvOverrideMap[r.Name]; ok {
			r.ShortGroup = gvk.Group
			r.Kind = gvk.Kind
			r.Version = gvk.Version
		}
	}
}

// namespacedCommonReferencesOverride adds referencers for fields that are known and common among
// more than a few resources.
func namespacedCommonReferencesOverride() ujconfig.ResourceOption {
	return func(r *ujconfig.Resource) {
		for k, s := range r.TerraformResource.Schema {
			// We shouldn't add referencers for status fields and sensitive fields
			// since they already have secret referencer.
			if (s.Computed && !s.Optional) || s.Sensitive {
				continue
			}

			// Most resources reference a "project_id". This is an ad-hoc way of
			// TODO: Set this field only in the config of the resources actually using it.
			if k == "project_id" {
				ref := ujconfig.Reference{
					Type:      commonNamespaced.APISPackagePath + "/mongodbatlas/v1alpha1.Project",
					Extractor: commonNamespaced.ExtractResourceIDFuncPath,
				}
				if r.ShortGroup == "" && r.Version == "v1alpha1" {
					ref.Type = "Project"
				}
				r.References["project_id"] = ref
			}
		}
	}
}
