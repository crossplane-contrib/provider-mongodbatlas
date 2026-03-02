# Resource Identity Reference

This document describes the minimum parameters needed to import `provider-mongodbatlas` resources.

Typically you use the `crossplane.io/external-name` annotation to tell Crossplane which is the ID of
a given resource, so it can internally _import_ the resource.

In this provider, this is managed automatically by Crossplane: you do not need to set it manually.
Just set the required parameters in the resource `spec` and Crossplane handles the rest.

Parameters marked with **ref** support `Ref`/`Selector` fields for cross-resource references (e.g. `projectIdRef`, `projectIdSelector`).

## Importable resources

These resources use one or more parameters as their identity.
Set the parameters (directly or via a selector) and Crossplane populates the external name
automatically.

Resources with a value in the **External Name** column have a provider-assigned ID. For new
resources this is set automatically after creation. For importing existing resources, set
`crossplane.io/external-name` to that provider-assigned value.

| Resource | Required Identity Parameters | External Name |
|----------|------------------------------|---------------|
| `mongodbatlas_access_list_api_key` | `org_id` (**ref**), `api_key_id` (**ref**), `ip_address` or `cidr_block` | |
| `mongodbatlas_advanced_cluster` | `project_id` (**ref**), `name` | |
| `mongodbatlas_alert_configuration` | `project_id` (**ref**) | `alert_id` |
| `mongodbatlas_api_key_project_assignment` | `project_id` (**ref**), `api_key_id` (**ref**) | |
| `mongodbatlas_api_key` | `org_id` (**ref**) | `api_key_id` |
| `mongodbatlas_auditing` | `project_id` (**ref**) | |
| `mongodbatlas_backup_compliance_policy` | `project_id` (**ref**) | |
| `mongodbatlas_cloud_backup_schedule` | `project_id` (**ref**), `cluster_name` | |
| `mongodbatlas_cloud_backup_snapshot_export_bucket` | `project_id` (**ref**) | `bucket_id` |
| `mongodbatlas_cloud_backup_snapshot_export_job` | `project_id` (**ref**), `cluster_name`, `snapshot_id` (**ref**), `export_bucket_id` (**ref**) | `job_id` |
| `mongodbatlas_cloud_backup_snapshot_restore_job` | `project_id` (**ref**), `cluster_name`, `snapshot_id` (**ref**) | `job_id` |
| `mongodbatlas_cloud_backup_snapshot` | `project_id` (**ref**), `cluster_name` | `snapshot_id` |
| `mongodbatlas_cloud_provider_access_setup` | `project_id` (**ref**), `provider_name` | `role_id` |
| `mongodbatlas_cloud_user_org_assignment` | `org_id` (**ref**), `username` | |
| `mongodbatlas_cloud_user_project_assignment` | `project_id` (**ref**), `username` | |
| `mongodbatlas_cloud_user_team_assignment` | `org_id` (**ref**), `team_id` (**ref**), `username` | |
| `mongodbatlas_cluster` | `project_id` (**ref**), `name` | |
| `mongodbatlas_custom_db_role` | `project_id` (**ref**), `role_name` | |
| `mongodbatlas_custom_dns_configuration_cluster_aws` | `project_id` (**ref**) | |
| `mongodbatlas_database_user` | `project_id` (**ref**), `username`, `auth_database_name` | `username` |
| `mongodbatlas_encryption_at_rest_private_endpoint` | `project_id` (**ref**), `cloud_provider` | `endpoint_id` |
| `mongodbatlas_encryption_at_rest` | `project_id` (**ref**) | |
| `mongodbatlas_event_trigger` | `project_id` (**ref**), `app_id` | `trigger_id` |
| `mongodbatlas_federated_database_instance` | `project_id` (**ref**), `name` | |
| `mongodbatlas_federated_query_limit` | `project_id` (**ref**), `tenant_name`, `limit_name` | |
| `mongodbatlas_federated_settings_identity_provider` | `federation_settings_id` | `idp_id` |
| `mongodbatlas_federated_settings_org_config` | `federation_settings_id`, `org_id` (**ref**) | |
| `mongodbatlas_federated_settings_org_role_mapping` | `federation_settings_id` (**ref**), `org_id` (**ref**) | `role_mapping_id` |
| `mongodbatlas_flex_cluster` | `project_id` (**ref**), `name` | |
| `mongodbatlas_global_cluster_config` | `project_id` (**ref**), `cluster_name` | |
| `mongodbatlas_ldap_configuration` | `project_id` (**ref**) | |
| `mongodbatlas_ldap_verify` | `project_id` (**ref**) | `request_id` |
| `mongodbatlas_log_integration` | `project_id` (**ref**) | `type` |
| `mongodbatlas_maintenance_window` | `project_id` (**ref**) | |
| `mongodbatlas_mongodb_employee_access_grant` | `project_id` (**ref**), `cluster_name` | |
| `mongodbatlas_network_container` | `project_id` (**ref**) | `container_id` |
| `mongodbatlas_network_peering` | `project_id` (**ref**), `container_id` (**ref**), `provider_name` | `peering_id` |
| `mongodbatlas_online_archive` | `project_id` (**ref**), `cluster_name` | `archive_id` |
| `mongodbatlas_org_invitation` | `org_id` (**ref**), `username` | |
| `mongodbatlas_organization` | n/a (requires ID) | `org_id` |
| `mongodbatlas_private_endpoint_regional_mode` | `project_id` (**ref**) | |
| `mongodbatlas_privatelink_endpoint_service_data_federation_online_archive` | `project_id` (**ref**), `endpoint_id` | |
| `mongodbatlas_privatelink_endpoint_service` | `project_id` (**ref**), `private_link_id` (**ref**), `provider_name` | `endpoint_service_id` |
| `mongodbatlas_privatelink_endpoint` | `project_id` (**ref**), `provider_name`, `region` | `endpoint_id` |
| `mongodbatlas_project_api_key` | `project_id` | |
| `mongodbatlas_project_invitation` | `project_id` (**ref**), `username` | |
| `mongodbatlas_project_ip_access_list` | `project_id` (**ref**), `ip_address` or `cidr_block` | |
| `mongodbatlas_project_service_account_access_list_entry` | `project_id` (**ref**), `client_id`, `ip_address` or `cidr_block` | |
| `mongodbatlas_project_service_account_secret` | `project_id` (**ref**), `client_id` | `secret_id` |
| `mongodbatlas_project_service_account` | `project_id` (**ref**) | `client_id` |
| `mongodbatlas_project` | `org_id` (**ref**) | |
| `mongodbatlas_push_based_log_export` | `project_id` (**ref**) | |
| `mongodbatlas_resource_policy` | `org_id` (**ref**) | `policy_id` |
| `mongodbatlas_search_deployment` | `project_id` (**ref**), `cluster_name` | |
| `mongodbatlas_serverless_instance` | `project_id` (**ref**), `name` | provider-assigned hex ID |
| `mongodbatlas_service_account_access_list_entry` | `org_id` (**ref**), `client_id`, `ip_address` or `cidr_block` | |
| `mongodbatlas_service_account_project_assignment` | `project_id` (**ref**), `client_id` | |
| `mongodbatlas_service_account_secret` | `org_id` (**ref**), `client_id` | `secret_id` |
| `mongodbatlas_service_account` | `org_id` (**ref**) | `client_id` |
| `mongodbatlas_stream_connection` | `workspace_name`, `project_id` (**ref**), `connection_name` | |
| `mongodbatlas_stream_instance` | `project_id` (**ref**), `instance_name` | |
| `mongodbatlas_stream_processor` | `instance_name`, `project_id` (**ref**), `processor_name` | |
| `mongodbatlas_stream_workspace` | `project_id` (**ref**), `workspace_name` | |
| `mongodbatlas_team_project_assignment` | `project_id`, `team_id` | |
| `mongodbatlas_team` | `org_id` (**ref**) | `team_id` |
| `mongodbatlas_third_party_integration` | `project_id` (**ref**), `type` | |
| `mongodbatlas_x509_authentication_database_user` | `project_id` (**ref**) | |

## Non-importable

The following resources do not support import:

- `mongodbatlas_cloud_provider_access_authorization`
- `mongodbatlas_cluster_outage_simulation`
- `mongodbatlas_search_index`
- `mongodbatlas_stream_privatelink_endpoint`
