package config

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

var externalNameConfigs = map[string]config.ExternalName{
	"mongodbatlas_access_list_api_key":                                         config.IdentifierFromProvider,
	"mongodbatlas_advanced_cluster":                                            templatedStringAsIdentifier("{{ .parameters.project_id }}-{{ .parameters.name }}"),
	"mongodbatlas_alert_configuration":                                         encodedStateID([]string{refs.ProjectID}, "id"),
	"mongodbatlas_api_key_project_assignment":                                  templatedStringAsIdentifier("{{ .parameters.project_id }}/{{ .parameters.api_key_id }}"),
	"mongodbatlas_api_key":                                                     config.IdentifierFromProvider,
	"mongodbatlas_auditing":                                                    config.IdentifierFromProvider,
	"mongodbatlas_backup_compliance_policy":                                    encodedStateID([]string{refs.ProjectID}, refs.ProjectID),
	"mongodbatlas_cloud_backup_schedule":                                       encodedStateID([]string{refs.ProjectID, refs.ClusterName}, refs.ClusterName),
	"mongodbatlas_cloud_backup_snapshot_export_bucket":                         encodedStateID([]string{refs.ProjectID}, "id"),
	"mongodbatlas_cloud_backup_snapshot_export_job":                            encodedStateID([]string{refs.ProjectID, refs.ClusterName}, "export_job_id"),
	"mongodbatlas_cloud_backup_snapshot_restore_job":                           encodedStateID([]string{refs.ProjectID, refs.ClusterName}, "snapshot_restore_job_id"),
	"mongodbatlas_cloud_backup_snapshot":                                       encodedStateID([]string{refs.ProjectID, refs.ClusterName}, "snapshot_id"),
	"mongodbatlas_cloud_provider_access_authorization":                         config.IdentifierFromProvider, // doesn't support import
	"mongodbatlas_cloud_provider_access_setup":                                 encodedStateID([]string{refs.ProjectID, refs.ProviderName}, "id"),
	"mongodbatlas_cloud_user_org_assignment":                                   templatedStringAsIdentifier("{{ .parameters.org_id }}/{{ .parameters.username }}"),
	"mongodbatlas_cloud_user_project_assignment":                               templatedStringAsIdentifier("{{ .parameters.project_id }}/{{ .parameters.username }}"),
	"mongodbatlas_cloud_user_team_assignment":                                  templatedStringAsIdentifier("{{ .parameters.org_id }}/{{ .parameters.team_id }}/{{ .parameters.username }}"),
	"mongodbatlas_cluster_outage_simulation":                                   config.IdentifierFromProvider, // doesn't support import
	"mongodbatlas_cluster":                                                     encodedStateIDMapped(map[string]string{refs.ProjectID: refs.ProjectID, refs.Name: refs.ClusterName}, refs.ClusterName),
	"mongodbatlas_custom_db_role":                                              encodedStateID([]string{refs.ProjectID, refs.RoleName}, refs.RoleName),
	"mongodbatlas_custom_dns_configuration_cluster_aws":                        templatedStringAsIdentifier("{{ .parameters.project_id }}"),
	"mongodbatlas_database_user":                                               templatedStringAsIdentifier("{{ .parameters.project_id }}/{{ .parameters.username }}/{{ .parameters.auth_database_name }}"),
	"mongodbatlas_encryption_at_rest":                                          templatedStringAsIdentifier("{{ .parameters.project_id }}"),
	"mongodbatlas_encryption_at_rest_private_endpoint":                         config.IdentifierFromProvider,
	"mongodbatlas_event_trigger":                                               config.IdentifierFromProvider,
	"mongodbatlas_federated_database_instance":                                 encodedStateID([]string{refs.ProjectID, refs.Name}, refs.Name),
	"mongodbatlas_federated_query_limit":                                       encodedStateID([]string{refs.ProjectID, "tenant_name", "limit_name"}, "limit_name"),
	"mongodbatlas_federated_settings_identity_provider":                        config.IdentifierFromProvider,
	"mongodbatlas_federated_settings_org_config":                               encodedStateID([]string{"federation_settings_id", refs.OrgID}, refs.OrgID),
	"mongodbatlas_federated_settings_org_role_mapping":                         config.IdentifierFromProvider,
	"mongodbatlas_flex_cluster":                                                templatedStringAsIdentifier("{{ .parameters.project_id }}-{{ .parameters.name }}"),
	"mongodbatlas_global_cluster_config":                                       encodedStateID([]string{refs.ProjectID, refs.ClusterName}, refs.ClusterName),
	"mongodbatlas_ldap_configuration":                                          templatedStringAsIdentifier("{{ .parameters.project_id }}"),
	"mongodbatlas_ldap_verify":                                                 config.IdentifierFromProvider,
	"mongodbatlas_log_integration":                                             config.IdentifierFromProvider,
	"mongodbatlas_maintenance_window":                                          templatedStringAsIdentifier("{{ .parameters.project_id }}"),
	"mongodbatlas_mongodb_employee_access_grant":                               templatedStringAsIdentifier("{{ .parameters.project_id }}/{{ .parameters.cluster_name }}"),
	"mongodbatlas_network_container":                                           encodedStateID([]string{refs.ProjectID}, "container_id"),
	"mongodbatlas_network_peering":                                             encodedStateID([]string{refs.ProjectID, refs.ProviderName}, "peer_id"),
	"mongodbatlas_online_archive":                                              config.IdentifierFromProvider,
	"mongodbatlas_org_invitation":                                              config.IdentifierFromProvider,
	"mongodbatlas_organization":                                                config.IdentifierFromProvider,
	"mongodbatlas_private_endpoint_regional_mode":                              templatedStringAsIdentifier("{{ .parameters.project_id }}"),
	"mongodbatlas_privatelink_endpoint_service_data_federation_online_archive": templatedStringAsIdentifier("{{ .parameters.project_id }}--{{ .parameters.endpoint_id }}"),
	"mongodbatlas_privatelink_endpoint_service":                                encodedStateID([]string{refs.ProjectID, "private_link_id", "endpoint_service_id", refs.ProviderName}, "endpoint_service_id"),
	"mongodbatlas_privatelink_endpoint":                                        encodedStateID([]string{refs.ProjectID, refs.ProviderName, refs.Region}, "private_link_id"),
	"mongodbatlas_project_api_key":                                             config.IdentifierFromProvider,
	"mongodbatlas_project_invitation":                                          config.IdentifierFromProvider,
	"mongodbatlas_project_ip_access_list":                                      config.IdentifierFromProvider,
	"mongodbatlas_project_service_account_access_list_entry":                   config.IdentifierFromProvider,
	"mongodbatlas_project_service_account_secret":                              config.IdentifierFromProvider,
	"mongodbatlas_project_service_account":                                     config.IdentifierFromProvider,
	"mongodbatlas_project":                                                     config.IdentifierFromProvider,
	"mongodbatlas_push_based_log_export":                                       templatedStringAsIdentifier("{{ .parameters.project_id }}"),
	"mongodbatlas_resource_policy":                                             config.IdentifierFromProvider,
	"mongodbatlas_search_deployment":                                           templatedStringAsIdentifier("{{ .parameters.project_id }}-{{ .parameters.cluster_name }}"),
	"mongodbatlas_search_index":                                                encodedStateID([]string{refs.ProjectID, refs.ClusterName}, "index_id"), // doesn't support import
	"mongodbatlas_serverless_instance":                                         encodedStateID([]string{refs.ProjectID, refs.Name}, refs.Name),
	"mongodbatlas_service_account_access_list_entry":                           config.IdentifierFromProvider,
	"mongodbatlas_service_account_project_assignment":                          templatedStringAsIdentifier("{{ .parameters.project_id }}/{{ .parameters.client_id }}"),
	"mongodbatlas_service_account_secret":                                      config.IdentifierFromProvider,
	"mongodbatlas_service_account":                                             config.IdentifierFromProvider,
	"mongodbatlas_stream_connection":                                           templatedStringAsIdentifier("{{ .parameters.workspace_name }}-{{ .parameters.project_id }}-{{ .parameters.connection_name }}"),
	"mongodbatlas_stream_instance":                                             templatedStringAsIdentifier("{{ .parameters.project_id }}-{{ .parameters.instance_name }}"),
	"mongodbatlas_stream_privatelink_endpoint":                                 config.IdentifierFromProvider, // doesn't support import
	"mongodbatlas_stream_processor":                                            templatedStringAsIdentifier("{{ .parameters.instance_name }}-{{ .parameters.project_id }}-{{ .parameters.processor_name }}"),
	"mongodbatlas_stream_workspace":                                            templatedStringAsIdentifier("{{ .parameters.project_id }}-{{ .parameters.workspace_name }}"),
	"mongodbatlas_team_project_assignment":                                     templatedStringAsIdentifier("{{ .parameters.project_id }}/{{ .parameters.team_id }}"),
	"mongodbatlas_team":                                                        config.IdentifierFromProvider,
	"mongodbatlas_third_party_integration":                                     templatedStringAsIdentifier("{{ .parameters.project_id }}-{{ .parameters.type }}"),
	"mongodbatlas_x509_authentication_database_user":                           encodedStateID([]string{refs.ProjectID}, refs.ProjectID),
}

// templatedStringAsIdentifier wraps config.TemplatedStringAsIdentifier with an
// empty nameField and overrides GetIDFn so that the crossplane.io/external-name
// annotation, when set, is treated as the canonical Terraform ID.
// The template is only rendered as a fallback when the annotation is absent and all
// parameter fields referenced by the template are populated in spec.forProvider.
//
// That is:
//   - annotation set, forProvider unset: annotation used as TF ID (import).
//   - annotation unset, forProvider fields populated: template rendered from params.
//   - both set: template rendered (parameters take precedence). Either way,
//     the annotation will be overwritten with the rendered/imported TF ID
//     after the first successful observe via upstream's GetExternalNameFn.
func templatedStringAsIdentifier(template string) config.ExternalName {
	e := config.TemplatedStringAsIdentifier("", template)
	identifierFields := slices.Clone(e.IdentifierFields)
	e.DisableNameInitializer = false
	// Clear IdentifierFields so that template parameters (e.g. project_id)
	// remain in the generated CRD schema as regular forProvider fields with
	// Ref/Selector support. The wrapped GetIDFn still reads them at runtime.
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
		// TPF-based resources in the Atlas TF provider don't include "id" in
		// their schema. Reconstruct the external name by rendering the template
		// with parameter values read directly from the Terraform state.
		return origGetIDFn(context.Background(), "", tfstate, nil)
	}

	return e
}

// encodeAtlasStateID replicates the encoding used by the Atlas TF provider's
// conversion.EncodeStateID (internal/common/conversion).
// It base64-encodes keys and values, joins each pair with ":", and joins all pairs with "-"
// (keys sorted alphabetically).
// That is, given {refs.ProjectID: "abc123", refs.RoleName: "myRole"}:
//
// 1. Sort keys alphabetically -> project_id, role_name
// 2. Base64-encode each key AND value:
//   - project_id -> cHJvamVjdF9pZA==
//   - abc123 -> YWJjMTIz
//   - role_name -> cm9sZV9uYW1l
//   - myRole -> bXlSb2xl
//
// 3. Join each key:value pair with :
//   - cHJvamVjdF9pZA==:YWJjMTIz
//   - cm9sZV9uYW1l:bXlSb2xl
//
// 4. Join all pairs with -
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

// encodedStateID is a convenience wrapper for encodedStateIDMapped where param names and state keys match.
func encodedStateID(fields []string, externalNameKey string) config.ExternalName {
	m := make(map[string]string, len(fields))
	for _, f := range fields {
		m[f] = f
	}
	return encodedStateIDMapped(m, externalNameKey)
}

// encodedStateIDMapped builds a config.ExternalName for Atlas resources whose
// TF provider uses conversion.EncodeStateID / DecodeStateID for d.Id().
//
// fieldMapping maps forProvider param names to EncodeStateID map keys
// (e.g. {refs.Name: refs.ClusterName} when the TF schema field differs from the
// state key). externalNameKey is the decoded state key whose value is stored
// in the crossplane.io/external-name annotation after the first successful
// observe.
//
// If externalNameKey is already in fieldMapping, all encoded keys come from
// forProvider parameters. If not (provider-assigned / computed field), the
// external name annotation supplies the missing value when reconstructing
// the encoded state ID.
func encodedStateIDMapped(fieldMapping map[string]string, externalNameKey string) config.ExternalName {
	paramNames := slices.Sorted(maps.Keys(fieldMapping))
	externalNameFromParams := slices.Contains(slices.Collect(maps.Values(fieldMapping)), externalNameKey)
	return config.ExternalName{
		DisableNameInitializer:  !externalNameFromParams,
		OmittedFields:           []string{},
		IdentifierFields:        nil,
		SetIdentifierArgumentFn: func(_ map[string]any, _ string) {},
		GetIDFn:                 encodedStateGetIDFn(fieldMapping, paramNames, externalNameKey),
		GetExternalNameFn:       encodedStateGetExternalNameFn(externalNameKey),
	}
}

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
				return "", fmt.Errorf("cannot determine Terraform ID: %s is not yet available", externalNameKey)
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

// hasAllParams returns true if every field in fields is present and non-empty
// in params.
func hasAllParams(params map[string]any, fields []string) bool {
	for _, f := range fields {
		s, ok := params[f].(string)
		if !ok || s == "" {
			return false
		}
	}
	return true
}

// ExternalNameConfigurations applies the external name config for each resource.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := externalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, 0, len(externalNameConfigs))
	for name := range externalNameConfigs {
		l = append(l, name+"$")
	}
	return l
}
