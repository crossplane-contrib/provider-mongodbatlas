package config

import (
	"github.com/crossplane-contrib/provider-jet-mongodbatlas/config/common"
	tjconfig "github.com/crossplane/terrajet/pkg/config"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var groupKindMap = map[string]schema.GroupKind{
	"mongodbatlas_advanced_cluster": {
		Group: "",
		Kind:  "AdvancedCluster",
	},
}

// groupKindOverrides overrides the group and kind of the resource if it matches
// any entry in the GroupMap.
func groupKindOverrides() tjconfig.ResourceOption {
	return func(r *tjconfig.Resource) {
		if r.ShortGroup == resourcePrefix {
			r.ShortGroup = ""
		}
		if val, ok := groupKindMap[r.Name]; ok {
			r.ShortGroup = val.Group
			r.Kind = val.Kind
		}
	}
}

// identifierAssignedByMongoDBAtlas is the most common external name
// configuration.
// The resource-specific configurations should override this whenever needed.
func identifierAssignedByMongoDBAtlas() tjconfig.ResourceOption {
	return func(r *tjconfig.Resource) {
		r.ExternalName = tjconfig.IdentifierFromProvider
	}
}

// commonReferences adds referencers for fields that are known and common among
// more than a few resources.
func commonReferences() tjconfig.ResourceOption {
	return func(r *tjconfig.Resource) {
		for k, s := range r.TerraformResource.Schema {
			// We shouldn't add referencers for status fields and sensitive fields
			// since they already have secret referencer.
			if (s.Computed && !s.Optional) || s.Sensitive {
				continue
			}

			if k == "project_id" {
				ref := tjconfig.Reference{
					Type:      common.APISPackagePath + "/mongodbatlas/v1alpha1.Project",
					Extractor: common.ExtractResourceIDFuncPath,
				}
				if r.ShortGroup == "" {
					ref.Type = "Project"
				}
				r.References["project_id"] = ref
			}
		}
	}
}
