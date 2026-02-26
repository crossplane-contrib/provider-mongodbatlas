package config

import "github.com/crossplane/upjet/v2/pkg/config"

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	"mongodbatlas_access_list_api_key":                                         config.IdentifierFromProvider,
	"mongodbatlas_advanced_cluster":                                            config.TemplatedStringAsIdentifier("name", "{{ .parameters.project_id }}-{{ .parameters.name }}"),
	"mongodbatlas_alert_configuration":                                         config.IdentifierFromProvider,
	"mongodbatlas_api_key_project_assignment":                                  config.TemplatedStringAsIdentifier("", "{{ .parameters.project_id }}/{{ .parameters.api_key_id }}"),
	"mongodbatlas_api_key":                                                     config.IdentifierFromProvider,
	"mongodbatlas_auditing":                                                    config.IdentifierFromProvider,
	"mongodbatlas_backup_compliance_policy":                                    config.ParameterAsIdentifier("project_id"),
	"mongodbatlas_cloud_backup_schedule":                                       config.TemplatedStringAsIdentifier("", "{{ .parameters.project_id }}/{{ .parameters.cluster_name }}"),
	"mongodbatlas_cloud_backup_snapshot_export_bucket":                         config.IdentifierFromProvider,
	"mongodbatlas_cloud_backup_snapshot_export_job":                            config.IdentifierFromProvider,
	"mongodbatlas_cloud_backup_snapshot_restore_job":                           config.IdentifierFromProvider,
	"mongodbatlas_cloud_backup_snapshot":                                       config.IdentifierFromProvider,
	"mongodbatlas_cloud_provider_access_authorization":                         config.IdentifierFromProvider, // doesn't support import
	"mongodbatlas_cloud_provider_access_setup":                                 config.IdentifierFromProvider,
	"mongodbatlas_cloud_user_org_assignment":                                   config.TemplatedStringAsIdentifier("", "{{ .parameters.org_id }}/{{ .parameters.username }}"),
	"mongodbatlas_cloud_user_project_assignment":                               config.TemplatedStringAsIdentifier("", "{{ .parameters.project_id }}/{{ .parameters.username }}"),
	"mongodbatlas_cloud_user_team_assignment":                                  config.TemplatedStringAsIdentifier("", "{{ .parameters.org_id }}/{{ .parameters.team_id }}/{{ .parameters.username }}"),
	"mongodbatlas_cluster_outage_simulation":                                   config.IdentifierFromProvider, // doesn't support import
	"mongodbatlas_cluster":                                                     config.TemplatedStringAsIdentifier("name", "{{ .parameters.project_id }}-{{ .parameters.name }}"),
	"mongodbatlas_custom_db_role":                                              config.TemplatedStringAsIdentifier("", "{{ .parameters.project_id }}/{{ .parameters.role_name }}"),
	"mongodbatlas_custom_dns_configuration_cluster_aws":                        config.ParameterAsIdentifier("project_id"),
	"mongodbatlas_database_user":                                               config.NameAsIdentifier,
	"mongodbatlas_encryption_at_rest":                                          config.ParameterAsIdentifier("project_id"),
	"mongodbatlas_encryption_at_rest_private_endpoint":                         config.IdentifierFromProvider,
	"mongodbatlas_event_trigger":                                               config.IdentifierFromProvider,
	"mongodbatlas_federated_database_instance":                                 config.TemplatedStringAsIdentifier("", "{{ .parameters.project_id }}--{{ .parameters.name }}"),
	"mongodbatlas_federated_query_limit":                                       config.TemplatedStringAsIdentifier("", "{{ .parameters.project_id }}--{{ .parameters.tenant_name }}--{{ .parameters.limit_name }}"),
	"mongodbatlas_federated_settings_identity_provider":                        config.IdentifierFromProvider,
	"mongodbatlas_federated_settings_org_config":                               config.TemplatedStringAsIdentifier("", "{{ .parameters.federation_settings_id }}/{{ .parameters.org_id }}"),
	"mongodbatlas_federated_settings_org_role_mapping":                         config.IdentifierFromProvider,
	"mongodbatlas_flex_cluster":                                                config.TemplatedStringAsIdentifier("name", "{{ .parameters.project_id }}-{{ .parameters.name }}"),
	"mongodbatlas_global_cluster_config":                                       config.TemplatedStringAsIdentifier("", "{{ .parameters.project_id }}/{{ .parameters.cluster_name }}"),
	"mongodbatlas_ldap_configuration":                                          config.ParameterAsIdentifier("project_id"),
	"mongodbatlas_ldap_verify":                                                 config.IdentifierFromProvider,
	"mongodbatlas_log_integration":                                             config.IdentifierFromProvider,
	"mongodbatlas_maintenance_window":                                          config.ParameterAsIdentifier("project_id"),
	"mongodbatlas_mongodb_employee_access_grant":                               config.TemplatedStringAsIdentifier("", "{{ .parameters.project_id }}/{{ .parameters.cluster_name }}"),
	"mongodbatlas_network_container":                                           config.IdentifierFromProvider,
	"mongodbatlas_network_peering":                                             config.IdentifierFromProvider,
	"mongodbatlas_online_archive":                                              config.IdentifierFromProvider,
	"mongodbatlas_org_invitation":                                              config.TemplatedStringAsIdentifier("", "{{ .parameters.org_id }}-{{ .parameters.username }}"),
	"mongodbatlas_organization":                                                config.IdentifierFromProvider,
	"mongodbatlas_private_endpoint_regional_mode":                              config.ParameterAsIdentifier("project_id"),
	"mongodbatlas_privatelink_endpoint_service_data_federation_online_archive": config.TemplatedStringAsIdentifier("", "{{ .parameters.project_id }}--{{ .parameters.endpoint_id }}"),
	"mongodbatlas_privatelink_endpoint_service":                                config.IdentifierFromProvider,
	"mongodbatlas_privatelink_endpoint":                                        config.IdentifierFromProvider,
	"mongodbatlas_project_api_key":                                             config.IdentifierFromProvider,
	"mongodbatlas_project_invitation":                                          config.TemplatedStringAsIdentifier("name", "{{ .parameters.project_id }}-{{ .parameters.username }}"),
	"mongodbatlas_project_ip_access_list":                                      config.IdentifierFromProvider,
	"mongodbatlas_project_service_account_access_list_entry":                   config.IdentifierFromProvider,
	"mongodbatlas_project_service_account_secret":                              config.IdentifierFromProvider,
	"mongodbatlas_project_service_account":                                     config.IdentifierFromProvider,
	"mongodbatlas_project":                                                     config.ParameterAsIdentifier("project_id"),
	"mongodbatlas_push_based_log_export":                                       config.ParameterAsIdentifier("project_id"),
	"mongodbatlas_resource_policy":                                             config.IdentifierFromProvider,
	"mongodbatlas_search_deployment":                                           config.TemplatedStringAsIdentifier("", "{{ .parameters.project_id }}-{{ .parameters.cluster_name }}"),
	"mongodbatlas_search_index":                                                config.IdentifierFromProvider, // doesn't support import
	"mongodbatlas_serverless_instance":                                         config.IdentifierFromProvider,
	"mongodbatlas_service_account_access_list_entry":                           config.IdentifierFromProvider,
	"mongodbatlas_service_account_project_assignment":                          config.TemplatedStringAsIdentifier("", "{{ .parameters.project_id }}/{{ .parameters.client_id }}"),
	"mongodbatlas_service_account_secret":                                      config.IdentifierFromProvider,
	"mongodbatlas_service_account":                                             config.IdentifierFromProvider,
	"mongodbatlas_stream_connection":                                           config.TemplatedStringAsIdentifier("", "{{ .parameters.workspace_name }}-{{ .parameters.project_id }}-{{ .parameters.connection_name }}"),
	"mongodbatlas_stream_instance":                                             config.TemplatedStringAsIdentifier("", "{{ .parameters.project_id }}-{{ .parameters.instance_name }}"),
	"mongodbatlas_stream_privatelink_endpoint":                                 config.IdentifierFromProvider, // doesn't support import
	"mongodbatlas_stream_processor":                                            config.TemplatedStringAsIdentifier("", "{{ .parameters.instance_name }}-{{ .parameters.project_id }}-{{ .parameters.processor_name }}"),
	"mongodbatlas_stream_workspace":                                            config.TemplatedStringAsIdentifier("", "{{ .parameters.project_id }}-{{ .parameters.workspace_name }}"),
	"mongodbatlas_team_project_assignment":                                     config.TemplatedStringAsIdentifier("", "{{ .parameters.project_id }}/{{ .parameters.team_id }}"),
	"mongodbatlas_team":                                                        config.IdentifierFromProvider,
	"mongodbatlas_third_party_integration":                                     config.TemplatedStringAsIdentifier("name", "{{ .parameters.project_id }}-{{ .parameters.type }}"),
	"mongodbatlas_x509_authentication_database_user":                           config.ParameterAsIdentifier("project_id"),
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}
