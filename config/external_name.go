package config

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

// --- Resource map ---

var externalNameConfigs = map[string]config.ExternalName{
	"mongodbatlas_access_list_api_key":                                         accessListImportJoinedID([]string{refs.OrgID, "api_key_id"}),
	"mongodbatlas_advanced_cluster":                                            templated("{{ .parameters.project_id }}-{{ .parameters.name }}"),
	"mongodbatlas_alert_configuration":                                         importJoinedID([]string{refs.ProjectID}, "-", "id"),
	"mongodbatlas_api_key_project_assignment":                                  templated("{{ .parameters.project_id }}/{{ .parameters.api_key_id }}"),
	"mongodbatlas_api_key":                                                     importJoinedID([]string{refs.OrgID}, "-", "api_key_id"),
	"mongodbatlas_auditing":                                                    identifierFromProvider(),
	"mongodbatlas_backup_compliance_policy":                                    importJoinedID([]string{refs.ProjectID}, "-", refs.ProjectID),
	"mongodbatlas_cloud_backup_schedule":                                       importJoinedID([]string{refs.ProjectID, refs.ClusterName}, "-", refs.ClusterName),
	"mongodbatlas_cloud_backup_snapshot_export_bucket":                         importJoinedID([]string{refs.ProjectID}, "-", "id"),
	"mongodbatlas_cloud_backup_snapshot_export_job":                            importJoinedID([]string{refs.ProjectID, refs.ClusterName}, "--", "export_job_id"),
	"mongodbatlas_cloud_backup_snapshot_restore_job":                           importJoinedID([]string{refs.ProjectID, refs.ClusterName}, "-", "snapshot_restore_job_id"),
	"mongodbatlas_cloud_backup_snapshot":                                       importJoinedID([]string{refs.ProjectID, refs.ClusterName}, "-", "snapshot_id"),
	"mongodbatlas_cloud_provider_access_authorization":                         importJoinedIDMapped([]string{refs.ProjectID, "role_id"}, map[string]string{refs.ProjectID: refs.ProjectID, "role_id": "id"}, "id"),
	"mongodbatlas_cloud_provider_access_setup":                                 importJoinedID([]string{refs.ProjectID, refs.ProviderName}, "-", "id"),
	"mongodbatlas_cloud_user_org_assignment":                                   templated("{{ .parameters.org_id }}/{{ .parameters.username }}"),
	"mongodbatlas_cloud_user_project_assignment":                               templated("{{ .parameters.project_id }}/{{ .parameters.username }}"),
	"mongodbatlas_cloud_user_team_assignment":                                  templated("{{ .parameters.org_id }}/{{ .parameters.team_id }}/{{ .parameters.username }}"),
	"mongodbatlas_cluster_outage_simulation":                                   config.IdentifierFromProvider, // doesn't support import
	"mongodbatlas_cluster":                                                     importJoinedIDMapped([]string{refs.ProjectID, refs.Name}, map[string]string{refs.ProjectID: refs.ProjectID, refs.Name: refs.ClusterName}, refs.ClusterName),
	"mongodbatlas_custom_db_role":                                              importJoinedID([]string{refs.ProjectID, refs.RoleName}, "-", refs.RoleName),
	"mongodbatlas_custom_dns_configuration_cluster_aws":                        templated("{{ .parameters.project_id }}"),
	"mongodbatlas_database_user":                                               importJoinedID([]string{refs.ProjectID, "username", "auth_database_name"}, "/", "username"),
	"mongodbatlas_encryption_at_rest":                                          templated("{{ .parameters.project_id }}"),
	"mongodbatlas_encryption_at_rest_private_endpoint":                         identifierFromProvider(),
	"mongodbatlas_event_trigger":                                               importJoinedID([]string{refs.ProjectID, "app_id"}, "--", "trigger_id"),
	"mongodbatlas_federated_database_instance":                                 importJoinedID([]string{refs.ProjectID, refs.Name}, "--", refs.Name),
	"mongodbatlas_federated_query_limit":                                       importJoinedID([]string{refs.ProjectID, "tenant_name", "limit_name"}, "--", "limit_name"),
	"mongodbatlas_federated_settings_identity_provider":                        importJoinedID([]string{"federation_settings_id"}, "-", "okta_idp_id"),
	"mongodbatlas_federated_settings_org_config":                               importJoinedID([]string{"federation_settings_id", refs.OrgID}, "-", refs.OrgID),
	"mongodbatlas_federated_settings_org_role_mapping":                         importJoinedID([]string{"federation_settings_id", refs.OrgID}, "-", "role_mapping_id"),
	"mongodbatlas_flex_cluster":                                                templated("{{ .parameters.project_id }}-{{ .parameters.name }}"),
	"mongodbatlas_global_cluster_config":                                       importJoinedID([]string{refs.ProjectID, refs.ClusterName}, "-", refs.ClusterName),
	"mongodbatlas_ldap_configuration":                                          templated("{{ .parameters.project_id }}"),
	"mongodbatlas_ldap_verify":                                                 importJoinedID([]string{refs.ProjectID}, "-", "request_id"),
	"mongodbatlas_log_integration":                                             identifierFromProvider(),
	"mongodbatlas_maintenance_window":                                          templated("{{ .parameters.project_id }}"),
	"mongodbatlas_mongodb_employee_access_grant":                               templated("{{ .parameters.project_id }}/{{ .parameters.cluster_name }}"),
	"mongodbatlas_network_container":                                           importJoinedID([]string{refs.ProjectID}, "-", "container_id"),
	"mongodbatlas_network_peering":                                             importJoinedIDOrdered([]string{refs.ProjectID, refs.PeerID, refs.ProviderName}, refs.PeerID),
	"mongodbatlas_online_archive":                                              importJoinedID([]string{refs.ProjectID, refs.ClusterName}, "-", "archive_id"),
	"mongodbatlas_org_invitation":                                              importJoinedIDHidden([]string{refs.OrgID, "username"}, "-", "invitation_id"),
	"mongodbatlas_organization":                                                importJoinedID([]string{}, "-", "org_id"),
	"mongodbatlas_private_endpoint_regional_mode":                              templated("{{ .parameters.project_id }}"),
	"mongodbatlas_privatelink_endpoint_service_data_federation_online_archive": importJoinedID([]string{refs.ProjectID, "endpoint_id"}, "--", "endpoint_id"),
	"mongodbatlas_privatelink_endpoint_service":                                importJoinedID([]string{refs.ProjectID, "private_link_id", "endpoint_service_id", refs.ProviderName}, "--", "endpoint_service_id"),
	"mongodbatlas_privatelink_endpoint":                                        importJoinedIDOrdered([]string{refs.ProjectID, "private_link_id", refs.ProviderName, refs.Region}, "private_link_id"),
	"mongodbatlas_project_api_key":                                             importJoinedID([]string{refs.ProjectID}, "-", "api_key_id"),
	"mongodbatlas_project_invitation":                                          importJoinedIDHidden([]string{refs.ProjectID, "username"}, "-", "invitation_id"),
	"mongodbatlas_project_ip_access_list":                                      accessListImportJoinedID([]string{refs.ProjectID}),
	"mongodbatlas_project_service_account_access_list_entry":                   identifierFromProvider(),
	"mongodbatlas_project_service_account_secret":                              identifierFromProvider(),
	"mongodbatlas_project_service_account":                                     identifierFromProvider(),
	"mongodbatlas_project":                                                     identifierFromProvider(),
	"mongodbatlas_push_based_log_export":                                       templated("{{ .parameters.project_id }}"),
	"mongodbatlas_resource_policy":                                             identifierFromProvider(),
	"mongodbatlas_search_deployment":                                           templated("{{ .parameters.project_id }}-{{ .parameters.cluster_name }}"),
	"mongodbatlas_search_index":                                                importJoinedID([]string{refs.ProjectID, refs.ClusterName}, "-", "index_id"),
	"mongodbatlas_serverless_instance":                                         importJoinedID([]string{refs.ProjectID, refs.Name}, "-", refs.Name),
	"mongodbatlas_service_account_access_list_entry":                           identifierFromProvider(),
	"mongodbatlas_service_account_project_assignment":                          templated("{{ .parameters.project_id }}/{{ .parameters.client_id }}"),
	"mongodbatlas_service_account_secret":                                      identifierFromProvider(),
	"mongodbatlas_service_account":                                             identifierFromProvider(),
	"mongodbatlas_stream_connection":                                           templated("{{ .parameters.workspace_name }}-{{ .parameters.project_id }}-{{ .parameters.connection_name }}"),
	"mongodbatlas_stream_instance":                                             templated("{{ .parameters.project_id }}-{{ .parameters.instance_name }}"),
	"mongodbatlas_stream_privatelink_endpoint":                                 config.IdentifierFromProvider, // doesn't support import
	"mongodbatlas_stream_processor":                                            templated("{{ .parameters.instance_name }}-{{ .parameters.project_id }}-{{ .parameters.processor_name }}"),
	"mongodbatlas_stream_workspace":                                            templated("{{ .parameters.project_id }}-{{ .parameters.workspace_name }}"),
	"mongodbatlas_team_project_assignment":                                     templated("{{ .parameters.project_id }}/{{ .parameters.team_id }}"),
	"mongodbatlas_team":                                                        importJoinedID([]string{refs.OrgID}, "-", "id"),
	"mongodbatlas_third_party_integration":                                     templated("{{ .parameters.project_id }}-{{ .parameters.type }}"),
	"mongodbatlas_x509_authentication_database_user":                           importJoinedID([]string{refs.ProjectID}, "-", refs.ProjectID),
}

// --- Map entry constructors ---

// identifierFromProvider wraps config.IdentifierFromProvider with the name
// initializer enabled. This seeds crossplane.io/external-name from
// metadata.name, providing a non-empty sentinel ID for the first Refresh.
// The sentinel triggers a 404 (not found) instead of an empty-ID error,
// allowing the reconciler to proceed to Create.
func identifierFromProvider() config.ExternalName {
	e := config.IdentifierFromProvider
	e.DisableNameInitializer = false
	return e
}

// importJoinedID builds an ExternalName for resources whose TF import function
// expects plain field values joined by separator.
//
// fields lists forProvider param names whose values form the import ID.
// When externalNameKey is in fields, all values come from params.
// When externalNameKey is NOT in fields (provider-assigned), its value
// comes from the external-name annotation and is appended to the import ID.
func importJoinedID(fields []string, separator string, externalNameKey string) config.ExternalName {
	externalNameFromParams := slices.Contains(fields, externalNameKey)
	importOrder := fields
	if !externalNameFromParams {
		importOrder = append(slices.Clone(fields), externalNameKey)
	}
	m := make(map[string]string, len(fields))
	for _, f := range fields {
		m[f] = f
	}
	return buildImportJoinedID(fields, importOrder, m, separator, externalNameKey, externalNameFromParams)
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
	m := make(map[string]string, len(paramFields))
	for _, f := range paramFields {
		m[f] = f
	}
	return buildImportJoinedID(paramFields, importOrder, m, "-", externalNameKey, false)
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
// format does NOT include the provider-assigned key. The import ID contains
// only the param fields; the provider-assigned key is used only for state
// encoding and the external-name annotation.
func importJoinedIDHidden(fields []string, separator, externalNameKey string) config.ExternalName {
	m := make(map[string]string, len(fields))
	for _, f := range fields {
		m[f] = f
	}
	return buildImportJoinedID(fields, fields, m, separator, externalNameKey, false)
}

// accessListImportJoinedID builds an ExternalName for access-list resources
// where the "entry" state key comes from either ip_address or cidr_block.
func accessListImportJoinedID(prefixParams []string) config.ExternalName {
	return config.ExternalName{
		DisableNameInitializer:  false,
		OmittedFields:           []string{},
		IdentifierFields:        nil,
		SetIdentifierArgumentFn: func(_ map[string]any, _ string) {},
		GetIDFn:                 accessListEncodedStateGetIDFn(prefixParams),
		GetImportIDFn:           refs.AccessListGetIDFn(prefixParams...),
		GetExternalNameFn:       encodedStateGetExternalNameFn("entry"),
	}
}

func buildImportJoinedID(paramFields, importOrder []string, fieldMapping map[string]string, separator, externalNameKey string, externalNameFromParams bool) config.ExternalName {
	return config.ExternalName{
		DisableNameInitializer:  !externalNameFromParams,
		OmittedFields:           []string{},
		IdentifierFields:        nil,
		SetIdentifierArgumentFn: func(_ map[string]any, _ string) {},
		GetIDFn:                 encodedStateGetIDFn(fieldMapping, paramFields, externalNameKey),
		GetImportIDFn:           plainImportGetIDFn(paramFields, importOrder, separator, externalNameKey),
		GetExternalNameFn:       encodedStateGetExternalNameFn(externalNameKey),
	}
}

// plainImportGetIDFn returns a GetIDFn that produces plain import IDs by
// joining field values with separator. This is the format TF import functions
// expect (e.g. "project_id-role_name").
func plainImportGetIDFn(paramFields, importOrder []string, separator, externalNameKey string) config.GetIDFn {
	return func(_ context.Context, externalName string, parameters, _ map[string]any) (string, error) {
		if hasAllParams(parameters, paramFields) {
			return collectImportValues(importOrder, parameters, externalName, externalNameKey, separator)
		}
		// Fallback: accept a base64-encoded state ID as external name
		// (handles resources with pre-existing encoded annotations).
		if externalName != "" {
			if decoded := decodeAtlasStateID(externalName); decoded[externalNameKey] != "" {
				return externalName, nil
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

// --- templated ---

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

// --- Public API ---

// ExternalNameConfigurations applies the external name config for each resource.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		e, ok := externalNameConfigs[r.Name]
		if !ok {
			return
		}
		r.ExternalName = e
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

// --- Utilities ---

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
