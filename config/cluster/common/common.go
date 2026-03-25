package common

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/v2/pkg/fieldpath"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/password"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	xpref "github.com/crossplane/crossplane-runtime/v2/pkg/reference"
	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/resource"
	corev1 "k8s.io/api/core/v1"
)

const (
	// errFmtNoAttribute is an error string for not-found attributes
	errFmtNoAttribute = `"attribute not found: %s`
	// errFmtUnexpectedType is an error string for attribute map values of unexpected type
	errFmtUnexpectedType = `unexpected type for attribute %s: Expecting a string`

	commonConfigPackagePath = "github.com/crossplane-contrib/provider-mongodbatlas/config/cluster/common"
	// ExtractResourceIDFuncPath holds the MongoDBAtlas resource ID extractor func name
	ExtractResourceIDFuncPath = commonConfigPackagePath + ".ExtractResourceID()"

	// VersionV1Alpha2 is used as minimum version for all manually configured resources.
	VersionV1Alpha2 = "v1alpha2"
	// VersionV1Alpha3 is used for resources that need a separate version to avoid type collisions.
	VersionV1Alpha3 = "v1alpha3"
)

const (
	// APISPackagePath is the package path for generated APIs root package
	APISPackagePath = "github.com/crossplane-contrib/provider-mongodbatlas/apis/cluster"
)

// GetAttributeValue reads a string attribute from the specified map
func GetAttributeValue(attrMap map[string]any, attr string) (string, error) {
	v, ok := attrMap[attr]
	if !ok {
		return "", errors.Errorf(errFmtNoAttribute, attr)
	}
	vStr, ok := v.(string)
	if !ok {
		return "", errors.Errorf(errFmtUnexpectedType, attr)
	}
	return vStr, nil
}

// ExtractResourceID extracts the value of `spec.atProvider.id`
// from a Terraformed resource. If mr is not a Terraformed
// resource, returns an empty string.
func ExtractResourceID() xpref.ExtractValueFn {
	return func(mr xpresource.Managed) string {
		tr, ok := mr.(resource.Terraformed)
		if !ok {
			return ""
		}
		return tr.GetID()
	}
}

// SetIdentifierFunc sets the identifier attribute `name` from a composite
// external-name where the identifier resides at index 0 of a colon-delimited
// string.
func SetIdentifierFunc(base map[string]any, externalName string) {
	parts := strings.Split(externalName, ":")
	base["name"] = parts[0]
}

// GetIDFromParamsAndExternalName returns a GetIDFn that builds a composite ID
// by extracting parameter values and inserting the externalName at the given
// index position, joining all parts with sep.
func GetIDFromParamsAndExternalName(sep string, externalNameIndex int, params ...string) func(context.Context, string, map[string]any, map[string]any) (string, error) {
	return func(_ context.Context, externalName string, parameters map[string]any, _ map[string]any) (string, error) {
		parts := make([]string, len(params)+1)
		parts[externalNameIndex] = externalName
		paramIdx := 0
		for i := range parts {
			if i == externalNameIndex {
				continue
			}
			v, ok := parameters[params[paramIdx]]
			if !ok {
				return "", fmt.Errorf("%s missing from parameters", params[paramIdx])
			}
			parts[i] = fmt.Sprint(v)
			paramIdx++
		}
		return strings.Join(parts, sep), nil
	}
}

// ExternalNameFromSegment returns a GetExternalNameFn that splits the
// resource ID by the given separator and returns the segment at the given
// index. If no index is provided, it returns the last segment.
//
// WARNING: This function uses strings.Split which breaks when the target
// segment itself contains the separator character (e.g., cluster names with
// dashes). For those cases, use ExternalNameFromID instead.
//
// This function is safe when the extracted segment is a fixed-format value
// that cannot contain the separator (e.g., hex IDs, enum values like AWS).
func ExternalNameFromSegment(sep string, index ...int) func(tfstate map[string]any) (string, error) {
	return func(tfstate map[string]any) (string, error) {
		idStr, err := ExtractIDFromState(tfstate)
		if err != nil {
			return "", err
		}
		parts := strings.Split(idStr, sep)
		if len(index) > 0 {
			if index[0] >= len(parts) {
				return "", fmt.Errorf("index %d out of range for ID %q split by %q into %d parts", index[0], idStr, sep, len(parts))
			}
			return parts[index[0]], nil
		}
		return parts[len(parts)-1], nil
	}
}

// ExternalNameFromID extracts the external name from a composite Terraform
// ID by skipping fixed-format segments from the left and right.
//
// skipLeft: number of segments to skip from the left (must not contain sep).
// skipRight: number of segments to skip from the right (must not contain sep).
//
// Everything between the skipped segments is returned as the external name,
// which may safely contain the separator character.
//
// Examples:
//
//	ExternalNameFromID("-", 1, 0)  on "507f-my-cluster"       -> "my-cluster"
//	ExternalNameFromID("-", 1, 1)  on "507f-my-cluster-abc12" -> "my-cluster"
//	ExternalNameFromID("-", 0, 1)  on "abc12-my-instance"     -> "abc12"
//	ExternalNameFromID("/", 1, 0)  on "507f/my-value"         -> "my-value"
func ExternalNameFromID(sep string, skipLeft, skipRight int) func(tfstate map[string]any) (string, error) {
	return func(tfstate map[string]any) (string, error) {
		idStr, err := ExtractIDFromState(tfstate)
		if err != nil {
			return "", err
		}

		result := idStr

		// Skip fixed-format segments from the left.
		for i := 0; i < skipLeft; i++ {
			idx := strings.Index(result, sep)
			if idx == -1 {
				return "", fmt.Errorf("expected at least %d left segments in ID %q with separator %q", skipLeft, idStr, sep)
			}
			result = result[idx+len(sep):]
		}

		// Skip fixed-format segments from the right.
		for i := 0; i < skipRight; i++ {
			idx := strings.LastIndex(result, sep)
			if idx == -1 {
				return "", fmt.Errorf("expected at least %d right segments in ID %q with separator %q", skipRight, idStr, sep)
			}
			result = result[:idx]
		}

		return result, nil
	}
}

func ExtractIDFromState(tfstate map[string]any) (string, error) {
	id, ok := tfstate["id"]
	if !ok {
		return "", errors.New("id attribute missing from state file")
	}
	idStr, ok := id.(string)
	if !ok {
		return "", errors.New("value of id needs to be string")
	}
	return idStr, nil
}

// PasswordSecretRefSetter is implemented by managed resources that expose a
// PasswordSecretRef field in their forProvider spec.
type PasswordSecretRefSetter interface {
	SetPasswordSecretRef(ref *v1.SecretKeySelector)
}

// PasswordGenerator returns an InitializerFn with the following behavior:
//   - BYOP mode: if byopSecretRefPath is already set on the managed resource,
//     the referenced secret is used as-is and no password is generated.
//   - Auto-generate mode: if writeConnectionSecretPath is set, a password is
//     generated and written to that secret under the "password" key. The
//     byopSecretRefPath field is then set on the managed resource to point to
//     the same secret so Terraform can inject it as the password parameter.
func PasswordGenerator(byopSecretRefPath, writeConnectionSecretPath string) config.NewInitializerFn {
	return func(cl client.Client) managed.Initializer {
		return managed.InitializerFn(func(ctx context.Context, mg xpresource.Managed) error {
			paved, err := fieldpath.PaveObject(mg)
			if err != nil {
				return fmt.Errorf("cannot pave object: %w", err)
			}

			// BYOP: if passwordSecretRef is already set, use it as-is.
			byopSel := &v1.SecretKeySelector{}
			if err := paved.GetValueInto(byopSecretRefPath, byopSel); err == nil && byopSel.Name != "" {
				return nil
			} else if xpresource.Ignore(fieldpath.IsNotFound, err) != nil {
				return fmt.Errorf("cannot read %s: %w", byopSecretRefPath, err)
			}

			// Auto-generate: check writeConnectionSecretToRef.
			connRef := &v1.SecretReference{}
			if err := paved.GetValueInto(writeConnectionSecretPath, connRef); err != nil {
				if fieldpath.IsNotFound(err) {
					return nil
				}
				return fmt.Errorf("cannot read %s: %w", writeConnectionSecretPath, err)
			}
			if connRef.Name == "" {
				return nil
			}

			ns := connRef.Namespace
			if ns == "" {
				ns = mg.GetNamespace()
			}

			const passwordKey = "password"
			s := &corev1.Secret{}
			getErr := cl.Get(ctx, types.NamespacedName{Namespace: ns, Name: connRef.Name}, s)
			if xpresource.IgnoreNotFound(getErr) != nil {
				return fmt.Errorf("cannot get connection secret: %w", getErr)
			}
			if getErr == nil && len(s.Data[passwordKey]) != 0 {
				// Password already exists; ensure passwordSecretRef points to it.
				return setPasswordSecretRef(ctx, cl, mg, connRef.Name, ns, passwordKey)
			}

			// Generate and write the password.
			pw, err := password.Generate()
			if err != nil {
				return fmt.Errorf("cannot generate password: %w", err)
			}
			s.SetName(connRef.Name)
			s.SetNamespace(ns)
			if !meta.WasCreated(s) {
				meta.AddOwnerReference(s, meta.AsOwner(meta.TypedReferenceTo(mg, mg.GetObjectKind().GroupVersionKind())))
			}
			if s.Data == nil {
				s.Data = make(map[string][]byte, 1)
			}
			s.Data[passwordKey] = []byte(pw)
			if err := xpresource.NewAPIPatchingApplicator(cl).Apply(ctx, s); err != nil {
				return fmt.Errorf("cannot apply password secret: %w", err)
			}
			return setPasswordSecretRef(ctx, cl, mg, connRef.Name, ns, passwordKey)
		})
	}
}

// setPasswordSecretRef sets the passwordSecretRef on the managed resource to
// point to the given secret/key, then persists the change via client.Update.
func setPasswordSecretRef(ctx context.Context, cl client.Client, mg xpresource.Managed, name, namespace, key string) error {
	setter, ok := mg.(PasswordSecretRefSetter)
	if !ok {
		return nil
	}
	setter.SetPasswordSecretRef(&v1.SecretKeySelector{
		SecretReference: v1.SecretReference{Name: name, Namespace: namespace},
		Key:             key,
	})
	if err := cl.Update(ctx, mg); err != nil {
		return fmt.Errorf("cannot update managed resource with password secret ref: %w", err)
	}
	return nil
}
