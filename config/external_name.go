package config

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

// --- Types ---

// externalNameEntry holds a config.ExternalName plus optional import format
// info for resources whose TF provider uses conversion.EncodeStateID.
// GetIDFn returns base64-encoded state IDs (required for TF Read/Refresh),
// while importOrder+separator describe the plain format the TF import function
// expects. ExternalNameConfigurations wraps the import function to translate
// between the two.
type externalNameEntry struct {
	config.ExternalName
	importOrder []string // state keys in TF import format order; nil = no import wrapping
	separator   string   // separator used by TF import function ("-" or "--")
}

// --- Resource map ---

var externalNameConfigs = map[string]externalNameEntry{
	"mongodbatlas_access_list_api_key":                                         identifierFromProvider(),
	"mongodbatlas_advanced_cluster":                                            templated("{{ .parameters.project_id }}-{{ .parameters.name }}"),
	"mongodbatlas_alert_configuration":                                         importJoinedID([]string{refs.ProjectID}, "-", "id"),
	"mongodbatlas_api_key_project_assignment":                                  templated("{{ .parameters.project_id }}/{{ .parameters.api_key_id }}"),
	"mongodbatlas_api_key":                                                     identifierFromProvider(),
	"mongodbatlas_auditing":                                                    identifierFromProvider(),
	"mongodbatlas_backup_compliance_policy":                                    importJoinedID([]string{refs.ProjectID}, "-", refs.ProjectID),
	"mongodbatlas_cloud_backup_schedule":                                       importJoinedID([]string{refs.ProjectID, refs.ClusterName}, "-", refs.ClusterName),
	"mongodbatlas_cloud_backup_snapshot_export_bucket":                         importJoinedID([]string{refs.ProjectID}, "-", "id"),
	"mongodbatlas_cloud_backup_snapshot_export_job":                            importJoinedID([]string{refs.ProjectID, refs.ClusterName}, "--", "export_job_id"),
	"mongodbatlas_cloud_backup_snapshot_restore_job":                           importJoinedID([]string{refs.ProjectID, refs.ClusterName}, "-", "snapshot_restore_job_id"),
	"mongodbatlas_cloud_backup_snapshot":                                       importJoinedID([]string{refs.ProjectID, refs.ClusterName}, "-", "snapshot_id"),
	"mongodbatlas_cloud_provider_access_authorization":                         identifierFromProvider(), // doesn't support import
	"mongodbatlas_cloud_provider_access_setup":                                 importJoinedID([]string{refs.ProjectID, refs.ProviderName}, "-", "id"),
	"mongodbatlas_cloud_user_org_assignment":                                   templated("{{ .parameters.org_id }}/{{ .parameters.username }}"),
	"mongodbatlas_cloud_user_project_assignment":                               templated("{{ .parameters.project_id }}/{{ .parameters.username }}"),
	"mongodbatlas_cloud_user_team_assignment":                                  templated("{{ .parameters.org_id }}/{{ .parameters.team_id }}/{{ .parameters.username }}"),
	"mongodbatlas_cluster_outage_simulation":                                   identifierFromProvider(), // doesn't support import
	"mongodbatlas_cluster":                                                     importJoinedIDMapped([]string{refs.ProjectID, refs.Name}, map[string]string{refs.ProjectID: refs.ProjectID, refs.Name: refs.ClusterName}),
	"mongodbatlas_custom_db_role":                                              importJoinedID([]string{refs.ProjectID, refs.RoleName}, "-", refs.RoleName),
	"mongodbatlas_custom_dns_configuration_cluster_aws":                        templated("{{ .parameters.project_id }}"),
	"mongodbatlas_database_user":                                               templated("{{ .parameters.project_id }}/{{ .parameters.username }}/{{ .parameters.auth_database_name }}"),
	"mongodbatlas_encryption_at_rest":                                          templated("{{ .parameters.project_id }}"),
	"mongodbatlas_encryption_at_rest_private_endpoint":                         identifierFromProvider(),
	"mongodbatlas_event_trigger":                                               identifierFromProvider(),
	"mongodbatlas_federated_database_instance":                                 importJoinedID([]string{refs.ProjectID, refs.Name}, "--", refs.Name),
	"mongodbatlas_federated_query_limit":                                       importJoinedID([]string{refs.ProjectID, "tenant_name", "limit_name"}, "--", "limit_name"),
	"mongodbatlas_federated_settings_identity_provider":                        identifierFromProvider(),
	"mongodbatlas_federated_settings_org_config":                               importJoinedID([]string{"federation_settings_id", refs.OrgID}, "-", refs.OrgID),
	"mongodbatlas_federated_settings_org_role_mapping":                         identifierFromProvider(),
	"mongodbatlas_flex_cluster":                                                templated("{{ .parameters.project_id }}-{{ .parameters.name }}"),
	"mongodbatlas_global_cluster_config":                                       importJoinedID([]string{refs.ProjectID, refs.ClusterName}, "-", refs.ClusterName),
	"mongodbatlas_ldap_configuration":                                          templated("{{ .parameters.project_id }}"),
	"mongodbatlas_ldap_verify":                                                 identifierFromProvider(),
	"mongodbatlas_log_integration":                                             identifierFromProvider(),
	"mongodbatlas_maintenance_window":                                          templated("{{ .parameters.project_id }}"),
	"mongodbatlas_mongodb_employee_access_grant":                               templated("{{ .parameters.project_id }}/{{ .parameters.cluster_name }}"),
	"mongodbatlas_network_container":                                           importJoinedID([]string{refs.ProjectID}, "-", "container_id"),
	"mongodbatlas_network_peering":                                             importJoinedIDOrdered([]string{refs.ProjectID, refs.PeerID, refs.ProviderName}, refs.PeerID),
	"mongodbatlas_online_archive":                                              identifierFromProvider(),
	"mongodbatlas_org_invitation":                                              identifierFromProvider(),
	"mongodbatlas_organization":                                                identifierFromProvider(),
	"mongodbatlas_private_endpoint_regional_mode":                              templated("{{ .parameters.project_id }}"),
	"mongodbatlas_privatelink_endpoint_service_data_federation_online_archive": templated("{{ .parameters.project_id }}--{{ .parameters.endpoint_id }}"),
	"mongodbatlas_privatelink_endpoint_service":                                importJoinedID([]string{refs.ProjectID, "private_link_id", "endpoint_service_id", refs.ProviderName}, "--", "endpoint_service_id"),
	"mongodbatlas_privatelink_endpoint":                                        importJoinedIDOrdered([]string{refs.ProjectID, "private_link_id", refs.ProviderName, refs.Region}, "private_link_id"),
	"mongodbatlas_project_api_key":                                             identifierFromProvider(),
	"mongodbatlas_project_invitation":                                          identifierFromProvider(),
	"mongodbatlas_project_ip_access_list":                                      identifierFromProvider(),
	"mongodbatlas_project_service_account_access_list_entry":                   identifierFromProvider(),
	"mongodbatlas_project_service_account_secret":                              identifierFromProvider(),
	"mongodbatlas_project_service_account":                                     identifierFromProvider(),
	"mongodbatlas_project":                                                     identifierFromProvider(),
	"mongodbatlas_push_based_log_export":                                       templated("{{ .parameters.project_id }}"),
	"mongodbatlas_resource_policy":                                             identifierFromProvider(),
	"mongodbatlas_search_deployment":                                           templated("{{ .parameters.project_id }}-{{ .parameters.cluster_name }}"),
	"mongodbatlas_search_index":                                                importJoinedID([]string{refs.ProjectID, refs.ClusterName}, "-", "index_id"), // doesn't support import
	"mongodbatlas_serverless_instance":                                         importJoinedID([]string{refs.ProjectID, refs.Name}, "-", refs.Name),
	"mongodbatlas_service_account_access_list_entry":                           identifierFromProvider(),
	"mongodbatlas_service_account_project_assignment":                          templated("{{ .parameters.project_id }}/{{ .parameters.client_id }}"),
	"mongodbatlas_service_account_secret":                                      identifierFromProvider(),
	"mongodbatlas_service_account":                                             identifierFromProvider(),
	"mongodbatlas_stream_connection":                                           templated("{{ .parameters.workspace_name }}-{{ .parameters.project_id }}-{{ .parameters.connection_name }}"),
	"mongodbatlas_stream_instance":                                             templated("{{ .parameters.project_id }}-{{ .parameters.instance_name }}"),
	"mongodbatlas_stream_privatelink_endpoint":                                 identifierFromProvider(), // doesn't support import
	"mongodbatlas_stream_processor":                                            templated("{{ .parameters.instance_name }}-{{ .parameters.project_id }}-{{ .parameters.processor_name }}"),
	"mongodbatlas_stream_workspace":                                            templated("{{ .parameters.project_id }}-{{ .parameters.workspace_name }}"),
	"mongodbatlas_team_project_assignment":                                     templated("{{ .parameters.project_id }}/{{ .parameters.team_id }}"),
	"mongodbatlas_team":                                                        identifierFromProvider(),
	"mongodbatlas_third_party_integration":                                     templated("{{ .parameters.project_id }}-{{ .parameters.type }}"),
	"mongodbatlas_x509_authentication_database_user":                           importJoinedID([]string{refs.ProjectID}, "-", refs.ProjectID),
}

// --- Map entry constructors ---

// identifierFromProvider wraps config.IdentifierFromProvider as an
// externalNameEntry so every entry in externalNameConfigs has the same type.
// No import wrapping needed: these resources use plain IDs for both state and import.
func identifierFromProvider() externalNameEntry {
	return externalNameEntry{ExternalName: config.IdentifierFromProvider}
}

// templated wraps templatedStringAsIdentifier as an externalNameEntry.
// No import wrapping needed: these resources use plain template-rendered IDs
// for both state and import.
func templated(tmpl string) externalNameEntry {
	return externalNameEntry{ExternalName: templatedStringAsIdentifier(tmpl)}
}

// importJoinedID builds an externalNameEntry for resources whose TF import
// function expects plain field values joined by separator, but whose TF
// Read/Refresh expects a base64-encoded state ID.
//
// GetIDFn returns base64 (for d.Id() / Refresh).
// importOrder+separator describe the plain format for import wrapping.
func importJoinedID(fields []string, separator string, externalNameKey string) externalNameEntry {
	m := make(map[string]string, len(fields))
	for _, f := range fields {
		m[f] = f
	}
	return importJoinedIDCore(m, fields, separator, externalNameKey)
}

// importJoinedIDOrdered handles resources where the provider-assigned key
// appears at a non-trailing position in the TF import format.
func importJoinedIDOrdered(importOrder []string, externalNameKey string) externalNameEntry {
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
	return externalNameEntry{
		ExternalName: config.ExternalName{
			DisableNameInitializer:  true,
			OmittedFields:           []string{},
			IdentifierFields:        nil,
			SetIdentifierArgumentFn: func(_ map[string]any, _ string) {},
			GetIDFn:                 encodedStateGetIDFn(m, paramFields, externalNameKey),
			GetExternalNameFn:       encodedStateGetExternalNameFn(externalNameKey),
		},
		importOrder: importOrder,
		separator:   "-",
	}
}

// importJoinedIDMapped handles resources where forProvider param names differ
// from TF state keys (e.g. name → cluster_name).
func importJoinedIDMapped(paramOrder []string, fieldMapping map[string]string) externalNameEntry {
	stateKeyOrder := make([]string, 0, len(paramOrder))
	for _, p := range paramOrder {
		stateKeyOrder = append(stateKeyOrder, fieldMapping[p])
	}
	return importJoinedIDCore(fieldMapping, stateKeyOrder, "-", refs.ClusterName)
}

func importJoinedIDCore(fieldMapping map[string]string, importOrder []string, separator, externalNameKey string) externalNameEntry {
	paramNames := slices.Sorted(maps.Keys(fieldMapping))
	externalNameFromParams := slices.Contains(slices.Collect(maps.Values(fieldMapping)), externalNameKey)
	fullImportOrder := importOrder
	if !externalNameFromParams && !slices.Contains(importOrder, externalNameKey) {
		fullImportOrder = append(slices.Clone(importOrder), externalNameKey)
	}
	return externalNameEntry{
		ExternalName: config.ExternalName{
			DisableNameInitializer:  !externalNameFromParams,
			OmittedFields:           []string{},
			IdentifierFields:        nil,
			SetIdentifierArgumentFn: func(_ map[string]any, _ string) {},
			GetIDFn:                 encodedStateGetIDFn(fieldMapping, paramNames, externalNameKey),
			GetExternalNameFn:       encodedStateGetExternalNameFn(externalNameKey),
		},
		importOrder: fullImportOrder,
		separator:   separator,
	}
}

// --- templatedStringAsIdentifier ---

// templatedStringAsIdentifier wraps config.TemplatedStringAsIdentifier with an
// empty nameField and overrides GetIDFn so that the crossplane.io/external-name
// annotation, when set, is treated as the canonical Terraform ID.
func templatedStringAsIdentifier(template string) config.ExternalName {
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

// ExternalNameConfigurations applies the external name config for each resource
// and wraps the TF import function for encoded-state resources.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		e, ok := externalNameConfigs[r.Name]
		if !ok {
			return
		}
		r.ExternalName = e.ExternalName
		if len(e.importOrder) > 0 {
			wrapImportForEncodedState(r, e.importOrder, e.separator)
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
