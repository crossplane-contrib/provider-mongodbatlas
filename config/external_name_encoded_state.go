package config

import (
	"encoding/base64"
	"fmt"
	"maps"
	"slices"
	"strings"
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
