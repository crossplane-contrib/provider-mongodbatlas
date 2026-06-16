package config

import (
	"context"
	"encoding/base64"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

// --- GetIDFn / GetExternalNameFn factories ---

// encodedStateGetIDFn returns a GetIDFn that produces base64-encoded state IDs
// matching the Atlas TF provider's conversion.EncodeStateID format.
// This is required because upjet sets d.Id() = GetIDFn result before TF Refresh,
// and TF Read calls DecodeStateID(d.Id()).
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

// --- Import wrapper ---

// wrapImportForEncodedState wraps the TF import function to decode a
// base64-encoded state ID into the plain format the original import function
// expects. This bridges the gap between GetIDFn (which must return base64 for
// TF Refresh) and the TF import function (which expects plain values).
func wrapImportForEncodedState(r *config.Resource, importOrder []string, separator string) {
	if r.TerraformResource == nil || r.TerraformResource.Importer == nil {
		return
	}
	orig := r.TerraformResource.Importer.StateContext
	if orig == nil {
		return
	}
	r.TerraformResource.Importer.StateContext = func(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
		if plainID := decodedToPlainImportID(d.Id(), importOrder, separator); plainID != "" {
			d.SetId(plainID)
		}
		return orig(ctx, d, meta)
	}
}

// decodedToPlainImportID converts a base64-encoded state ID to the plain format
// expected by TF import functions (e.g. "{project_id}-{role_name}").
func decodedToPlainImportID(encodedID string, importOrder []string, separator string) string {
	decoded := decodeAtlasStateID(encodedID)
	if len(decoded) == 0 {
		return ""
	}
	values := make([]string, 0, len(importOrder))
	for _, key := range importOrder {
		v := decoded[key]
		if v == "" {
			return ""
		}
		values = append(values, v)
	}
	return strings.Join(values, separator)
}
