package config

import (
	"slices"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
)

// includedResources lists all TF resources the provider manages.
// Each resource's external name is configured in its AddResourceConfigurator
// callback (config/resources/*.go). This list only controls inclusion.
var includedResources = []string{
	"mongodbatlas_access_list_api_key",
	"mongodbatlas_advanced_cluster",
	"mongodbatlas_alert_configuration",
	"mongodbatlas_api_key",
	"mongodbatlas_api_key_project_assignment",
	"mongodbatlas_auditing",
	"mongodbatlas_backup_compliance_policy",
	"mongodbatlas_cloud_backup_schedule",
	"mongodbatlas_cloud_backup_snapshot",
	"mongodbatlas_cloud_backup_snapshot_export_bucket",
	"mongodbatlas_cloud_backup_snapshot_export_job",
	"mongodbatlas_cloud_backup_snapshot_restore_job",
	"mongodbatlas_cloud_provider_access_authorization",
	"mongodbatlas_cloud_provider_access_setup",
	"mongodbatlas_cloud_user_org_assignment",
	"mongodbatlas_cloud_user_project_assignment",
	"mongodbatlas_cloud_user_team_assignment",
	"mongodbatlas_cluster",
	"mongodbatlas_cluster_outage_simulation",
	"mongodbatlas_custom_db_role",
	"mongodbatlas_custom_dns_configuration_cluster_aws",
	"mongodbatlas_database_user",
	"mongodbatlas_encryption_at_rest",
	"mongodbatlas_encryption_at_rest_private_endpoint",
	"mongodbatlas_event_trigger",
	"mongodbatlas_federated_database_instance",
	"mongodbatlas_federated_query_limit",
	"mongodbatlas_federated_settings_identity_provider",
	"mongodbatlas_federated_settings_org_config",
	"mongodbatlas_federated_settings_org_role_mapping",
	"mongodbatlas_flex_cluster",
	"mongodbatlas_global_cluster_config",
	"mongodbatlas_ldap_configuration",
	"mongodbatlas_ldap_verify",
	"mongodbatlas_log_integration",
	"mongodbatlas_maintenance_window",
	"mongodbatlas_mongodb_employee_access_grant",
	"mongodbatlas_network_container",
	"mongodbatlas_network_peering",
	"mongodbatlas_online_archive",
	"mongodbatlas_org_invitation",
	"mongodbatlas_organization",
	"mongodbatlas_private_endpoint_regional_mode",
	"mongodbatlas_privatelink_endpoint",
	"mongodbatlas_privatelink_endpoint_service",
	"mongodbatlas_privatelink_endpoint_service_data_federation_online_archive",
	"mongodbatlas_project",
	"mongodbatlas_project_api_key",
	"mongodbatlas_project_invitation",
	"mongodbatlas_project_ip_access_list",
	"mongodbatlas_project_service_account",
	"mongodbatlas_project_service_account_access_list_entry",
	"mongodbatlas_project_service_account_secret",
	"mongodbatlas_push_based_log_export",
	"mongodbatlas_resource_policy",
	"mongodbatlas_search_deployment",
	"mongodbatlas_search_index",
	"mongodbatlas_serverless_instance",
	"mongodbatlas_service_account",
	"mongodbatlas_service_account_access_list_entry",
	"mongodbatlas_service_account_project_assignment",
	"mongodbatlas_service_account_secret",
	"mongodbatlas_stream_connection",
	"mongodbatlas_stream_instance",
	"mongodbatlas_stream_privatelink_endpoint",
	"mongodbatlas_stream_processor",
	"mongodbatlas_stream_workspace",
	"mongodbatlas_team",
	"mongodbatlas_team_project_assignment",
	"mongodbatlas_third_party_integration",
	"mongodbatlas_x509_authentication_database_user",
}

// identifierFromProvider wraps config.IdentifierFromProvider with the name
// initializer enabled. This seeds crossplane.io/external-name from
// metadata.name, providing a non-empty sentinel ID for the first Refresh.
func identifierFromProvider() ujconfig.ExternalName {
	e := ujconfig.IdentifierFromProvider
	e.DisableNameInitializer = false
	return e
}

// ExternalNameConfigurations applies the default external name config
// (identifierFromProvider) to every included resource. Per-resource
// configurators override this with resource-specific external names.
func ExternalNameConfigurations() ujconfig.ResourceOption {
	return func(r *ujconfig.Resource) {
		if !slices.Contains(includedResources, r.Name) {
			return
		}
		r.ExternalName = identifierFromProvider()
	}
}

// ExternalNameConfigured returns the list of all included resources
// as regex patterns for upjet's WithIncludeList.
func ExternalNameConfigured() []string {
	l := make([]string, len(includedResources))
	for i, name := range includedResources {
		l[i] = name + "$"
	}
	return l
}
