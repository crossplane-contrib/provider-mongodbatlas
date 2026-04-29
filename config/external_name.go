package config

import (
	"context"

	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/iancoleman/strcase"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// externalNameEntry pairs an ExternalName config with the CRD field name
// to use as the initial external-name annotation (for templated configs).
type externalNameEntry struct {
	config.ExternalName
	// initCRDField is the forProvider field whose value should seed the crossplane.io/external-name
	// annotation before the first observe, replacing the default metadata.name.
	// Empty for IdentifierFromProvider resources.
	initCRDField string
}

// templated creates an externalNameEntry for template-based identifiers.
func templated(template, externalNameField string) externalNameEntry {
	return externalNameEntry{
		ExternalName: templatedStringAsIdentifier(template, externalNameField),
		initCRDField: strcase.ToLowerCamel(externalNameField),
	}
}

// identifierFromProvider creates an externalNameEntry for IdentifierFromProvider resources.
func identifierFromProvider() externalNameEntry {
	return externalNameEntry{ExternalName: config.IdentifierFromProvider}
}

// templatedStringAsIdentifier wraps config.TemplatedStringAsIdentifier with an
// empty nameField and enables the name initializer. The upstream helper
// disables the initializer when nameField is empty, which prevents Observe-only
// management policies from working because upjet skips the observe when the
// external-name annotation is absent.
//
// externalNameField is the terraform state field whose value will be used as
// the crossplane.io/external-name annotation after the first successful
// observe. Without this, the annotation would keep the initial metadata.name
// value forever.
func templatedStringAsIdentifier(template string, externalNameField string) config.ExternalName {
	e := config.TemplatedStringAsIdentifier("", template)
	e.DisableNameInitializer = false
	// Clear IdentifierFields so that parameters referenced in the template
	// (e.g. project_id) are kept in the generated CRD schema as regular
	// forProvider fields with Ref/Selector support.  The GetIDFn built by
	// TemplatedStringAsIdentifier still reads them from .parameters at
	// runtime, which is unaffected by this change.
	e.IdentifierFields = nil
	if externalNameField != "" {
		e.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
			if v, ok := tfstate[externalNameField]; ok {
				return v.(string), nil
			}
			return "", nil
		}
	}
	return e
}

var externalNameConfigs = map[string]externalNameEntry{
	"mongodbatlas_access_list_api_key":                                         identifierFromProvider(),
	"mongodbatlas_advanced_cluster":                                            templated("{{ .parameters.project_id }}-{{ .parameters.name }}", "name"),
	"mongodbatlas_alert_configuration":                                         identifierFromProvider(),
	"mongodbatlas_api_key_project_assignment":                                  templated("{{ .parameters.project_id }}/{{ .parameters.api_key_id }}", "api_key_id"),
	"mongodbatlas_api_key":                                                     identifierFromProvider(),
	"mongodbatlas_auditing":                                                    identifierFromProvider(),
	"mongodbatlas_backup_compliance_policy":                                    templated("{{ .parameters.project_id }}", "project_id"),
	"mongodbatlas_cloud_backup_schedule":                                       templated("{{ .parameters.project_id }}/{{ .parameters.cluster_name }}", "cluster_name"),
	"mongodbatlas_cloud_backup_snapshot_export_bucket":                         identifierFromProvider(),
	"mongodbatlas_cloud_backup_snapshot_export_job":                            identifierFromProvider(),
	"mongodbatlas_cloud_backup_snapshot_restore_job":                           identifierFromProvider(),
	"mongodbatlas_cloud_backup_snapshot":                                       identifierFromProvider(),
	"mongodbatlas_cloud_provider_access_authorization":                         identifierFromProvider(), // doesn't support import
	"mongodbatlas_cloud_provider_access_setup":                                 identifierFromProvider(),
	"mongodbatlas_cloud_user_org_assignment":                                   templated("{{ .parameters.org_id }}/{{ .parameters.username }}", "username"),
	"mongodbatlas_cloud_user_project_assignment":                               templated("{{ .parameters.project_id }}/{{ .parameters.username }}", "username"),
	"mongodbatlas_cloud_user_team_assignment":                                  templated("{{ .parameters.org_id }}/{{ .parameters.team_id }}/{{ .parameters.username }}", "username"),
	"mongodbatlas_cluster_outage_simulation":                                   identifierFromProvider(), // doesn't support import
	"mongodbatlas_cluster":                                                     templated("{{ .parameters.project_id }}-{{ .parameters.name }}", "name"),
	"mongodbatlas_custom_db_role":                                              templated("{{ .parameters.project_id }}-{{ .parameters.role_name }}", "role_name"),
	"mongodbatlas_custom_dns_configuration_cluster_aws":                        templated("{{ .parameters.project_id }}", "project_id"),
	"mongodbatlas_database_user":                                               templated("{{ .parameters.project_id }}/{{ .parameters.username }}/{{ .parameters.auth_database_name }}", "username"),
	"mongodbatlas_encryption_at_rest":                                          templated("{{ .parameters.project_id }}", "project_id"),
	"mongodbatlas_encryption_at_rest_private_endpoint":                         identifierFromProvider(),
	"mongodbatlas_event_trigger":                                               identifierFromProvider(),
	"mongodbatlas_federated_database_instance":                                 templated("{{ .parameters.project_id }}--{{ .parameters.name }}", "name"),
	"mongodbatlas_federated_query_limit":                                       templated("{{ .parameters.project_id }}--{{ .parameters.tenant_name }}--{{ .parameters.limit_name }}", "limit_name"),
	"mongodbatlas_federated_settings_identity_provider":                        identifierFromProvider(),
	"mongodbatlas_federated_settings_org_config":                               templated("{{ .parameters.federation_settings_id }}/{{ .parameters.org_id }}", "org_id"),
	"mongodbatlas_federated_settings_org_role_mapping":                         identifierFromProvider(),
	"mongodbatlas_flex_cluster":                                                templated("{{ .parameters.project_id }}-{{ .parameters.name }}", "name"),
	"mongodbatlas_global_cluster_config":                                       templated("{{ .parameters.project_id }}/{{ .parameters.cluster_name }}", "cluster_name"),
	"mongodbatlas_ldap_configuration":                                          templated("{{ .parameters.project_id }}", "project_id"),
	"mongodbatlas_ldap_verify":                                                 identifierFromProvider(),
	"mongodbatlas_log_integration":                                             identifierFromProvider(),
	"mongodbatlas_maintenance_window":                                          templated("{{ .parameters.project_id }}", "project_id"),
	"mongodbatlas_mongodb_employee_access_grant":                               templated("{{ .parameters.project_id }}/{{ .parameters.cluster_name }}", "cluster_name"),
	"mongodbatlas_network_container":                                           identifierFromProvider(),
	"mongodbatlas_network_peering":                                             identifierFromProvider(),
	"mongodbatlas_online_archive":                                              identifierFromProvider(),
	"mongodbatlas_org_invitation":                                              templated("{{ .parameters.org_id }}-{{ .parameters.username }}", "username"),
	"mongodbatlas_organization":                                                identifierFromProvider(),
	"mongodbatlas_private_endpoint_regional_mode":                              templated("{{ .parameters.project_id }}", "project_id"),
	"mongodbatlas_privatelink_endpoint_service_data_federation_online_archive": templated("{{ .parameters.project_id }}--{{ .parameters.endpoint_id }}", "endpoint_id"),
	"mongodbatlas_privatelink_endpoint_service":                                identifierFromProvider(),
	"mongodbatlas_privatelink_endpoint":                                        identifierFromProvider(),
	"mongodbatlas_project_api_key":                                             identifierFromProvider(),
	"mongodbatlas_project_invitation":                                          templated("{{ .parameters.project_id }}-{{ .parameters.username }}", "username"),
	"mongodbatlas_project_ip_access_list":                                      identifierFromProvider(),
	"mongodbatlas_project_service_account_access_list_entry":                   identifierFromProvider(),
	"mongodbatlas_project_service_account_secret":                              identifierFromProvider(),
	"mongodbatlas_project_service_account":                                     identifierFromProvider(),
	"mongodbatlas_project":                                                     {ExternalName: config.ParameterAsIdentifier("project_id")},
	"mongodbatlas_push_based_log_export":                                       templated("{{ .parameters.project_id }}", "project_id"),
	"mongodbatlas_resource_policy":                                             identifierFromProvider(),
	"mongodbatlas_search_deployment":                                           templated("{{ .parameters.project_id }}-{{ .parameters.cluster_name }}", "cluster_name"),
	"mongodbatlas_search_index":                                                identifierFromProvider(), // doesn't support import
	"mongodbatlas_serverless_instance":                                         identifierFromProvider(),
	"mongodbatlas_service_account_access_list_entry":                           identifierFromProvider(),
	"mongodbatlas_service_account_project_assignment":                          templated("{{ .parameters.project_id }}/{{ .parameters.client_id }}", "client_id"),
	"mongodbatlas_service_account_secret":                                      identifierFromProvider(),
	"mongodbatlas_service_account":                                             identifierFromProvider(),
	"mongodbatlas_stream_connection":                                           templated("{{ .parameters.workspace_name }}-{{ .parameters.project_id }}-{{ .parameters.connection_name }}", "connection_name"),
	"mongodbatlas_stream_instance":                                             templated("{{ .parameters.project_id }}-{{ .parameters.instance_name }}", "instance_name"),
	"mongodbatlas_stream_privatelink_endpoint":                                 identifierFromProvider(), // doesn't support import
	"mongodbatlas_stream_processor":                                            templated("{{ .parameters.instance_name }}-{{ .parameters.project_id }}-{{ .parameters.processor_name }}", "processor_name"),
	"mongodbatlas_stream_workspace":                                            templated("{{ .parameters.project_id }}-{{ .parameters.workspace_name }}", "workspace_name"),
	"mongodbatlas_team_project_assignment":                                     templated("{{ .parameters.project_id }}/{{ .parameters.team_id }}", "team_id"),
	"mongodbatlas_team":                                                        identifierFromProvider(),
	"mongodbatlas_third_party_integration":                                     templated("{{ .parameters.project_id }}-{{ .parameters.type }}", "type"),
	"mongodbatlas_x509_authentication_database_user":                           templated("{{ .parameters.project_id }}", "project_id"),
}

// ExternalNameConfigurations applies the external name config for each resource.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if entry, ok := externalNameConfigs[r.Name]; ok {
			r.ExternalName = entry.ExternalName
		}
	}
}

// ExternalNameInitializers registers an initializer for resources that seed the
// external-name annotation from a forProvider field instead of metadata.name.
func ExternalNameInitializers() config.ResourceOption {
	return func(r *config.Resource) {
		if entry, ok := externalNameConfigs[r.Name]; ok && entry.initCRDField != "" {
			r.InitializerFns = append(r.InitializerFns, initFromField(entry.initCRDField))
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(externalNameConfigs))
	i := 0
	for name := range externalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}

// initFromField returns a NewInitializerFn that seeds the external-name
// annotation from forProvider.crdField before NameAsExternalName runs.
func initFromField(crdField string) config.NewInitializerFn {
	return func(c client.Client) managed.Initializer {
		return managed.InitializerFn(func(ctx context.Context, mg resource.Managed) error {
			if meta.GetExternalName(mg) != "" {
				return nil
			}
			u, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(mg)
			val, _, _ := unstructured.NestedString(u, "spec", "forProvider", crdField)
			if val == "" {
				return nil
			}
			meta.SetExternalName(mg, val)
			return c.Update(ctx, mg)
		})
	}
}
