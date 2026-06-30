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
		if externalName != "" {
			if decoded := decodeAtlasStateID(externalName); decoded[externalNameKey] != "" {
				return externalName, nil
			}
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
		return id, nil
	}
}

// --- External name constructors ---

// importJoinedID builds an ExternalName for resources whose TF import function
// expects plain field values joined by separator.
func importJoinedID(fields []string, separator string, externalNameKey string) config.ExternalName {
	externalNameFromParams := slices.Contains(fields, externalNameKey)
	importOrder := fields
	if !externalNameFromParams {
		importOrder = append(slices.Clone(fields), externalNameKey)
	}
	return buildImportJoinedID(fields, importOrder, nil, separator, externalNameKey, externalNameFromParams)
}

// importJoinedIDOrdered handles resources where the provider-assigned key
// appears at a non-trailing position in the TF import format.
func importJoinedIDOrdered(importOrder []string, externalNameKey string) config.ExternalName {
	paramFields := make([]string, 0, len(importOrder)-1)
	for _, f := range importOrder {
		if f != externalNameKey {
			paramFields = append(paramFields, f)
		}
	}
	return buildImportJoinedID(paramFields, importOrder, nil, "-", externalNameKey, false)
}

// importJoinedIDMapped handles resources where forProvider param names differ
// from TF state keys (e.g. name → cluster_name, role_id → id).
func importJoinedIDMapped(paramOrder []string, fieldMapping map[string]string, externalNameKey string) config.ExternalName {
	stateKeyOrder := make([]string, 0, len(paramOrder))
	for _, p := range paramOrder {
		stateKeyOrder = append(stateKeyOrder, fieldMapping[p])
	}
	externalNameFromParams := slices.Contains(stateKeyOrder, externalNameKey)
	return buildImportJoinedID(paramOrder, paramOrder, fieldMapping, "-", externalNameKey, externalNameFromParams)
}

// importJoinedIDHidden builds an ExternalName for resources whose TF import
// format does NOT include the provider-assigned key.
func importJoinedIDHidden(fields []string, separator, externalNameKey string) config.ExternalName {
	return buildImportJoinedID(fields, fields, nil, separator, externalNameKey, false)
}

// accessListImportJoinedID builds an ExternalName for access-list resources
// where the "entry" state key comes from either ip_address or cidr_block.
func accessListImportJoinedID(prefixParams []string) config.ExternalName {
	e := baseExternalName(false)
	e.GetIDFn = accessListEncodedStateGetIDFn(prefixParams)
	e.GetImportIDFn = refs.AccessListGetIDFn(prefixParams...)
	e.GetExternalNameFn = encodedStateGetExternalNameFn("entry")
	return e
}

func buildImportJoinedID(paramFields, importOrder []string, fieldMapping map[string]string, separator, externalNameKey string, externalNameFromParams bool) config.ExternalName {
	e := baseExternalName(!externalNameFromParams)
	e.GetIDFn = encodedStateGetIDFn(fieldMapping, paramFields, externalNameKey)
	e.GetImportIDFn = plainImportGetIDFn(paramFields, importOrder, separator, externalNameKey)
	e.GetExternalNameFn = encodedStateGetExternalNameFn(externalNameKey)
	return e
}

func plainImportGetIDFn(paramFields, importOrder []string, separator, externalNameKey string) config.GetIDFn {
	return func(_ context.Context, externalName string, parameters, _ map[string]any) (string, error) {
		if hasAllParams(parameters, paramFields) {
			return collectImportValues(importOrder, parameters, externalName, externalNameKey, separator)
		}
		if externalName != "" {
			if decoded := decodeAtlasStateID(externalName); decoded[externalNameKey] != "" {
				values := make([]string, 0, len(importOrder))
				for _, field := range importOrder {
					values = append(values, decoded[field])
				}
				return strings.Join(values, separator), nil
			}
		}
		return "", fmt.Errorf("cannot determine Terraform ID: forProvider is missing %v and crossplane.io/external-name is empty or not a valid encoded state ID", paramFields)
	}
}

func collectImportValues(importOrder []string, parameters map[string]any, externalName, externalNameKey, separator string) (string, error) {
	values := make([]string, 0, len(importOrder))
	for _, field := range importOrder {
		if v, ok := parameters[field].(string); ok && v != "" {
			values = append(values, v)
		} else if field == externalNameKey && externalName != "" {
			values = append(values, externalName)
		} else if field == externalNameKey {
			return "", nil
		}
	}
	if len(values) == len(importOrder) {
		return strings.Join(values, separator), nil
	}
	return "", nil
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
