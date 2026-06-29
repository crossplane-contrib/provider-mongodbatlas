package refs

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

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
// Examples:
//
//	ExternalNameFromID("-", 1, 0)  on "507f-my-cluster"       -> "my-cluster"
//	ExternalNameFromID("-", 1, 1)  on "507f-my-cluster-abc12" -> "my-cluster"
//	ExternalNameFromID("/", 1, 0)  on "507f/my-value"         -> "my-value"
func ExternalNameFromID(sep string, skipLeft, skipRight int) func(tfstate map[string]any) (string, error) {
	return func(tfstate map[string]any) (string, error) {
		idStr, err := ExtractIDFromState(tfstate)
		if err != nil {
			return "", err
		}

		result := idStr

		for range skipLeft {
			idx := strings.Index(result, sep)
			if idx == -1 {
				return "", fmt.Errorf("expected at least %d left segments in ID %q with separator %q", skipLeft, idStr, sep)
			}
			result = result[idx+len(sep):]
		}

		for range skipRight {
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

// ExternalNameFromIDOrState wraps ExternalNameFromID with a fallback for
// TPF-based resources where "id" is absent from the TF schema.
func ExternalNameFromIDOrState(sep string, skipLeft, skipRight int, stateField string) func(tfstate map[string]any) (string, error) {
	fromID := ExternalNameFromID(sep, skipLeft, skipRight)
	return func(tfstate map[string]any) (string, error) {
		if _, ok := tfstate["id"]; ok {
			return fromID(tfstate)
		}
		v, ok := tfstate[stateField].(string)
		if !ok || v == "" {
			return "", fmt.Errorf("%s not found in Terraform state", stateField)
		}
		return v, nil
	}
}

// ExternalNameFromStateField returns a GetExternalNameFn that reads the given
// field directly from the Terraform state when "id" is absent.
func ExternalNameFromStateField(fields ...string) func(tfstate map[string]any) (string, error) {
	return func(tfstate map[string]any) (string, error) {
		if id, ok := tfstate["id"].(string); ok && id != "" {
			return id, nil
		}
		parts := make([]string, len(fields))
		for i, f := range fields {
			v, ok := tfstate[f].(string)
			if !ok || v == "" {
				return "", fmt.Errorf("%s not found in Terraform state", f)
			}
			parts[i] = v
		}
		return strings.Join(parts, "-"), nil
	}
}

func stateString(tfstate map[string]any, key string) (string, bool) {
	v, ok := tfstate[key].(string)
	return v, ok && v != ""
}

// AccessListEntry extracts the access list entry value from parameters,
// trying ip_address first then cidr_block.
func AccessListEntry(parameters map[string]any) (string, bool) {
	if v, ok := parameters["ip_address"].(string); ok && v != "" {
		return v, true
	}
	if v, ok := parameters["cidr_block"].(string); ok && v != "" {
		return v, true
	}
	return "", false
}

// AccessListGetIDFn returns a GetIDFn for access list resources where the IP
// is specified via either ip_address or cidr_block.
func AccessListGetIDFn(prefixParams ...string) func(context.Context, string, map[string]any, map[string]any) (string, error) {
	return func(_ context.Context, _ string, parameters map[string]any, _ map[string]any) (string, error) {
		parts := make([]string, 0, len(prefixParams)+1)
		for _, p := range prefixParams {
			v, ok := parameters[p]
			if !ok {
				return "", fmt.Errorf("%s missing from parameters", p)
			}
			parts = append(parts, fmt.Sprint(v))
		}
		entry, ok := AccessListEntry(parameters)
		if !ok {
			return "", errors.New("either ip_address or cidr_block parameters must be set")
		}
		parts = append(parts, entry)
		return strings.Join(parts, "-"), nil
	}
}

// ExtractParamPath builds an upjet extractor reference path for the given field.
func ExtractParamPath(field string, sensitive bool) string {
	return fmt.Sprintf(ExtractParamPathFmt, field, sensitive)
}

// ExternalNameFromAccessListState returns a GetExternalNameFn for access list
// entry resources. When "id" is present, returns it directly. Otherwise
// constructs "{scopeID}-{client_id}-{ip_address_or_cidr}" from state fields.
func ExternalNameFromAccessListState(scopeField string) func(tfstate map[string]any) (string, error) {
	return func(tfstate map[string]any) (string, error) {
		if id, ok := stateString(tfstate, "id"); ok {
			return id, nil
		}
		scope, ok := stateString(tfstate, scopeField)
		if !ok {
			return "", fmt.Errorf("%s not found in Terraform state", scopeField)
		}
		client, ok := stateString(tfstate, "client_id")
		if !ok {
			return "", errors.New("client_id not found in Terraform state")
		}
		ip, ok := stateString(tfstate, "ip_address")
		if !ok {
			ip, ok = stateString(tfstate, "cidr_block")
			if !ok {
				return "", errors.New("neither ip_address nor cidr_block found in Terraform state")
			}
		}
		return fmt.Sprintf("%s-%s-%s", scope, client, ip), nil
	}
}
