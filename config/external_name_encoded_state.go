package config

import (
	"context"
	"encoding/base64"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

// --- Base64 encode/decode (Atlas TF provider's EncodeStateID format) ---

// encodeAtlasStateID replicates the Atlas TF provider's conversion.EncodeStateID.
// Base64-encodes keys and values, joins each pair with ":", joins all pairs with "-"
// (keys sorted alphabetically).
func encodeAtlasStateID(values map[string]string) string {
	encode := func(s string) string { return base64.StdEncoding.EncodeToString([]byte(s)) }
	parts := make([]string, 0, len(values))
	for _, key := range slices.Sorted(maps.Keys(values)) {
		parts = append(parts, fmt.Sprintf("%s:%s", encode(key), encode(values[key])))
	}
	return strings.Join(parts, "-")
}

// decodeAtlasStateID reverses encodeAtlasStateID, returning the key-value map.
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

// encodedStateGetIDFn returns a GetIDFn that produces base64-encoded state IDs
// matching the Atlas TF provider's conversion.EncodeStateID format.
// This is required because terraform refresh uses d.Id() from the tfstate,
// and the TF Read function calls DecodeStateID(d.Id()).
func encodedStateGetIDFn(fieldMapping map[string]string, paramNames []string, externalNameKey string) func(context.Context, string, map[string]any, map[string]any) (string, error) {
	return func(_ context.Context, externalName string, parameters, _ map[string]any) (string, error) {
		if hasAllParams(parameters, paramNames) {
			m := make(map[string]string, len(fieldMapping)+1)
			for param, stateKey := range fieldMapping {
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

// accessListEncodedStateGetIDFn produces base64-encoded state IDs for
// access-list resources where the "entry" key comes from either
// ip_address or cidr_block.
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

// encodedStateGetExternalNameFn returns a GetExternalNameFn that decodes the
// base64-encoded d.Id() from TF state and extracts the externalNameKey value.
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
