package resources

import (
	"context"
	"encoding/base64"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

// --- Base64 encode/decode (Atlas TF provider's EncodeStateID format) ---

func encodeAtlasStateID(values map[string]string) string {
	encode := func(s string) string { return base64.StdEncoding.EncodeToString([]byte(s)) }
	parts := make([]string, 0, len(values))
	for _, key := range slices.Sorted(maps.Keys(values)) {
		parts = append(parts, fmt.Sprintf("%s:%s", encode(key), encode(values[key])))
	}
	return strings.Join(parts, "-")
}

func decodeAtlasStateID(stateID string) map[string]string {
	decode := func(s string) string {
		b, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return ""
		}
		return string(b)
	}
	result := make(map[string]string)
	for part := range strings.SplitSeq(stateID, "-") {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 {
			result[decode(kv[0])] = decode(kv[1])
		}
	}
	return result
}

// --- GetIDFn factory (base64 state ID for terraform refresh) ---

func encodedStateGetIDFn(fieldMapping map[string]string, paramNames []string, externalNameKey string) func(context.Context, string, map[string]any, map[string]any) (string, error) {
	return func(_ context.Context, externalName string, parameters, _ map[string]any) (string, error) {
		rawExternalName := externalName
		// v1.1.x wrote the full encoded state ID into the external name.
		if v := decodeAtlasStateID(externalName)[externalNameKey]; v != "" {
			externalName = v
		}
		if hasAllParams(parameters, paramNames) {
			m := make(map[string]string, len(paramNames)+1)
			for _, param := range paramNames {
				stateKey := param
				if fieldMapping != nil {
					stateKey = fieldMapping[param]
				}
				m[stateKey] = parameters[param].(string)
			}
			if _, ok := m[externalNameKey]; !ok && externalName != "" {
				m[externalNameKey] = externalName
			}
			if m[externalNameKey] == "" {
				return "", nil
			}
			return encodeAtlasStateID(m), nil
		}
		if externalName != rawExternalName {
			return rawExternalName, nil
		}
		return "", fmt.Errorf("cannot determine Terraform ID: forProvider is missing %v and crossplane.io/external-name is empty or not a valid encoded state ID", paramNames)
	}
}

func accessListEncodedStateGetIDFn(prefixParams []string) func(context.Context, string, map[string]any, map[string]any) (string, error) {
	return func(_ context.Context, _ string, parameters, _ map[string]any) (string, error) {
		values := make(map[string]string, len(prefixParams)+1)
		for _, param := range prefixParams {
			v, ok := parameters[param].(string)
			if !ok || v == "" {
				return "", nil
			}
			values[param] = v
		}
		entry, ok := refs.AccessListEntry(parameters)
		if !ok {
			return "", nil
		}
		values["entry"] = entry
		return encodeAtlasStateID(values), nil
	}
}

// --- GetExternalNameFn factory ---

func encodedStateGetExternalNameFn(externalNameKey string) func(map[string]any) (string, error) {
	return func(tfstate map[string]any) (string, error) {
		id, ok := tfstate["id"].(string)
		if !ok || id == "" {
			return "", fmt.Errorf("id not found in Terraform state")
		}
		decoded := decodeAtlasStateID(id)
		if v := decoded[externalNameKey]; v != "" {
			return v, nil
		}
		return "", fmt.Errorf("key %q not found in encoded state ID", externalNameKey)
	}
}

// --- External name constructors ---

// importJoinedID builds an ExternalName for resources whose TF state ID is
// EncodeStateID over fields. The externalNameKey must appear in fields. It is
// treated as a regular forProvider parameter (user-settable, included in CRD
// schema).
func importJoinedID(fields []string, externalNameKey string) config.ExternalName {
	return buildImportJoinedID(fields, nil, externalNameKey, true)
}

// importJoinedIDAssigned is like importJoinedID but the externalNameKey is
// provider-assigned (not user-settable).
func importJoinedIDAssigned(fields []string, externalNameKey string) config.ExternalName {
	return buildImportJoinedID(fields, nil, externalNameKey, false)
}

// importJoinedIDMapped handles resources where forProvider param names differ
// from TF state keys (e.g. name → cluster_name, role_id → id).
func importJoinedIDMapped(paramOrder []string, fieldMapping map[string]string, externalNameKey string) config.ExternalName {
	stateKeyOrder := make([]string, 0, len(paramOrder))
	for _, p := range paramOrder {
		stateKeyOrder = append(stateKeyOrder, fieldMapping[p])
	}
	externalNameFromParams := slices.Contains(stateKeyOrder, externalNameKey)
	return buildImportJoinedID(paramOrder, fieldMapping, externalNameKey, externalNameFromParams)
}

// accessListImportJoinedID builds an ExternalName for access-list resources
// where the "entry" state key comes from either ip_address or cidr_block.
func accessListImportJoinedID(prefixParams []string) config.ExternalName {
	e := baseExternalName(false)
	e.GetIDFn = accessListEncodedStateGetIDFn(prefixParams)
	e.GetExternalNameFn = encodedStateGetExternalNameFn("entry")
	return e
}

// computedKeyID builds an ExternalName for plugin-framework resources with no
// "id" attribute whose Read needs the provider-assigned computed attribute key.
// upjet does not call GetIDFn for them.
// SetIdentifierArgumentFn copies the external name into key, so that an import with only the annotation can Read.
// Pair it with refs.StateEmptyWhenAttributeUnset(key) so that Create runs while key is still unset.
func computedKeyID(key string) config.ExternalName {
	e := baseExternalName(true)
	e.GetIDFn = config.IdentifierFromProvider.GetIDFn
	e.GetExternalNameFn = refs.ExternalNameFromStateField(key)
	setKey := refs.SetIdentifierArgument(key)
	e.SetIdentifierArgumentFn = func(base map[string]any, externalName string) {
		// Keep a key already restored from status.atProvider: an older
		// external name can hold another value (log_integration used "type").
		if v, _ := base[key].(string); v != "" {
			return
		}
		setKey(base, externalName)
	}
	return e
}

func buildImportJoinedID(fields []string, fieldMapping map[string]string, externalNameKey string, externalNameFromParams bool) config.ExternalName {
	paramFields := fields
	if !externalNameFromParams {
		filtered := make([]string, 0, len(fields))
		for _, f := range fields {
			if f != externalNameKey {
				filtered = append(filtered, f)
			}
		}
		paramFields = filtered
	}
	e := baseExternalName(!externalNameFromParams)
	e.GetIDFn = encodedStateGetIDFn(fieldMapping, paramFields, externalNameKey)
	e.GetExternalNameFn = encodedStateGetExternalNameFn(externalNameKey)
	return e
}

// templated wraps config.TemplatedStringAsIdentifier with an
// empty nameField and overrides GetIDFn so that the crossplane.io/external-name
// annotation, when set, is treated as the canonical Terraform ID.
func templated(template string) config.ExternalName {
	e := config.TemplatedStringAsIdentifier("", template)
	identifierFields := slices.Clone(e.IdentifierFields)
	e.DisableNameInitializer = false
	e.IdentifierFields = nil

	origGetIDFn := e.GetIDFn
	origGetExternalNameFn := e.GetExternalNameFn

	e.GetIDFn = func(ctx context.Context, externalName string, parameters, providerConfig map[string]any) (string, error) {
		if hasAllParams(parameters, identifierFields) {
			return origGetIDFn(ctx, externalName, parameters, providerConfig)
		}
		if externalName != "" {
			return externalName, nil
		}
		return "", fmt.Errorf("cannot determine Terraform ID: forProvider is missing %v and crossplane.io/external-name annotation is empty", identifierFields)
	}

	e.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
		if _, ok := tfstate["id"]; ok {
			return origGetExternalNameFn(tfstate)
		}
		return origGetIDFn(context.Background(), "", tfstate, nil)
	}

	return e
}

// --- Utilities ---

func baseExternalName(disableNameInit bool) config.ExternalName {
	return config.ExternalName{
		DisableNameInitializer:  disableNameInit,
		OmittedFields:           []string{},
		IdentifierFields:        nil,
		SetIdentifierArgumentFn: func(_ map[string]any, _ string) {},
	}
}

func hasAllParams(params map[string]any, fields []string) bool {
	for _, f := range fields {
		s, ok := params[f].(string)
		if !ok || s == "" {
			return false
		}
	}
	return true
}
