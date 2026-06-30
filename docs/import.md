# Resource Identity Reference

This document describes the minimum parameters needed to import `provider-mongodbatlas` resources.

The provider supports two ways to identify an existing resource for import:

## Option A: by parameters (recommended)

Set the required identity parameters in `spec.forProvider`.
The provider computes the Terraform ID internally and writes it into the `crossplane.io/external-name`
annotation after the first successful observe.
You do not need to set the annotation manually.

## Option B: by annotation
Set `crossplane.io/external-name` directly to the full Terraform ID.
The provider uses that value as-is and skips template rendering.
Useful when you only have the ID and do not want to populate every identity parameter, or when adopting a resource
without binding it to specific `forProvider` values.

If both are set, parameters take precedence: the rendered Terraform ID overwrites the annotation after the first observe.

Parameters marked with **ref** support `Ref`/`Selector` fields for cross-resource references (e.g. `projectIdRef`, `projectIdSelector`).

### External-name representation

After the first successful observe, the provider writes a value into the `crossplane.io/external-name` annotation.
For most resources this is the full Terraform ID (e.g. `<projectId>-<name>`).
For resources that use the Atlas `EncodeStateID` pattern internally, the annotation contains only
the human-readable key (e.g. `username` for `database_user`, `endpointId` for
`privatelink_endpoint_service_data_federation_online_archive`).
Option A (parameters) is the recommended import method for these resources because
Option B requires the annotation to contain the full compound ID in base64-encoded format, which
is an internal implementation detail of the Atlas Terraform provider.

## Importable resources

These resources use one or more parameters as their identity.
Set the parameters (directly or via a selector) and Crossplane populates the external name
automatically.

Resources with a value in the **External Name** column have a provider-assigned ID. For new
resources this is set automatically after creation. For importing existing resources, set
`crossplane.io/external-name` to that provider-assigned value.

| Resource | Required Identity Parameters | External Name |
|----------|------------------------------|---------------|
| `mongodbatlas_access_list_api_key` | `orgId` (**ref**), `apiKeyId` (**ref**), `ipAddress` or `cidrBlock` | |
| `mongodbatlas_advanced_cluster` | `projectId` (**ref**), `name` | |
| `mongodbatlas_alert_configuration` | `projectId` (**ref**) | `alertId` |
| `mongodbatlas_api_key_project_assignment` | `projectId` (**ref**), `apiKeyId` (**ref**) | |
| `mongodbatlas_api_key` | `orgId` (**ref**) | `apiKeyId` |
| `mongodbatlas_auditing` | `projectId` (**ref**) | |
| `mongodbatlas_backup_compliance_policy` | `projectId` (**ref**) | |
| `mongodbatlas_cloud_backup_schedule` | `projectId` (**ref**), `clusterName` | |
| `mongodbatlas_cloud_backup_snapshot_export_bucket` | `projectId` (**ref**) | `bucketId` |
| `mongodbatlas_cloud_backup_snapshot_export_job` | `projectId` (**ref**), `clusterName` | `jobId` |
| `mongodbatlas_cloud_backup_snapshot_restore_job` | `projectId` (**ref**), `clusterName` | `jobId` |
| `mongodbatlas_cloud_backup_snapshot` | `projectId` (**ref**), `clusterName` | `snapshotId` |
| `mongodbatlas_cloud_provider_access_setup` | `projectId` (**ref**), `providerName` | `roleId` |
| `mongodbatlas_cloud_user_org_assignment` | `orgId` (**ref**), `username` | |
| `mongodbatlas_cloud_user_project_assignment` | `projectId` (**ref**), `username` | |
| `mongodbatlas_cloud_user_team_assignment` | `orgId` (**ref**), `teamId` (**ref**), `username` | |
| `mongodbatlas_cluster` | `projectId` (**ref**), `name` | |
| `mongodbatlas_custom_db_role` | `projectId` (**ref**), `roleName` | |
| `mongodbatlas_custom_dns_configuration_cluster_aws` | `projectId` (**ref**) | |
| `mongodbatlas_database_user` | `projectId` (**ref**), `username`, `authDatabaseName` | |
| `mongodbatlas_encryption_at_rest_private_endpoint` | `projectId` (**ref**), `cloudProvider` | `endpointId` |
| `mongodbatlas_encryption_at_rest` | `projectId` (**ref**) | |
| `mongodbatlas_event_trigger` | `projectId` (**ref**), `appId` | `triggerId` |
| `mongodbatlas_federated_database_instance` | `projectId` (**ref**), `name` | |
| `mongodbatlas_federated_query_limit` | `projectId` (**ref**), `tenantName`, `limitName` | |
| `mongodbatlas_federated_settings_identity_provider` | `federationSettingsId` | `idpId` |
| `mongodbatlas_federated_settings_org_config` | `federationSettingsId`, `orgId` (**ref**) | |
| `mongodbatlas_federated_settings_org_role_mapping` | `federationSettingsId` (**ref**), `orgId` (**ref**) | `roleMappingId` |
| `mongodbatlas_flex_cluster` | `projectId` (**ref**), `name` | |
| `mongodbatlas_global_cluster_config` | `projectId` (**ref**), `clusterName` | |
| `mongodbatlas_ldap_configuration` | `projectId` (**ref**) | |
| `mongodbatlas_ldap_verify` | `projectId` (**ref**) | `requestId` |
| `mongodbatlas_log_integration` | `projectId` (**ref**) | `type` |
| `mongodbatlas_maintenance_window` | `projectId` (**ref**) | |
| `mongodbatlas_mongodb_employee_access_grant` | `projectId` (**ref**), `clusterName` | |
| `mongodbatlas_network_container` | `projectId` (**ref**) | `containerId` |
| `mongodbatlas_network_peering` | `projectId` (**ref**), `providerName` | `peeringId` |
| `mongodbatlas_online_archive` | `projectId` (**ref**), `clusterName` | `archiveId` |
| `mongodbatlas_org_invitation` | `orgId` (**ref**), `username` | |
| `mongodbatlas_organization` | n/a (requires ID) | `orgId` |
| `mongodbatlas_private_endpoint_regional_mode` | `projectId` (**ref**) | |
| `mongodbatlas_privatelink_endpoint_service_data_federation_online_archive` | `projectId` (**ref**), `endpointId` | |
| `mongodbatlas_privatelink_endpoint_service` | `projectId` (**ref**), `privateLinkId` (**ref**), `providerName` | `endpointServiceId` |
| `mongodbatlas_privatelink_endpoint` | `projectId` (**ref**), `providerName`, `region` | `endpointId` |
| `mongodbatlas_project_api_key` | `projectId` | |
| `mongodbatlas_project_invitation` | `projectId` (**ref**), `username` | |
| `mongodbatlas_project_ip_access_list` | `projectId` (**ref**), `ipAddress` or `cidrBlock` | |
| `mongodbatlas_project_service_account_access_list_entry` | `projectId` (**ref**), `clientId`, `ipAddress` or `cidrBlock` | |
| `mongodbatlas_project_service_account_secret` | `projectId` (**ref**), `clientId` | `secretId` |
| `mongodbatlas_project_service_account` | `projectId` (**ref**) | `clientId` |
| `mongodbatlas_project` | | `projectId` |
| `mongodbatlas_push_based_log_export` | `projectId` (**ref**) | |
| `mongodbatlas_resource_policy` | `orgId` (**ref**) | `policyId` |
| `mongodbatlas_search_deployment` | `projectId` (**ref**), `clusterName` | |
| `mongodbatlas_serverless_instance` | `projectId` (**ref**), `name` | provider-assigned hex ID |
| `mongodbatlas_service_account_access_list_entry` | `orgId` (**ref**), `clientId`, `ipAddress` or `cidrBlock` | |
| `mongodbatlas_service_account_project_assignment` | `projectId` (**ref**), `clientId` | |
| `mongodbatlas_service_account_secret` | `orgId` (**ref**), `clientId` | `secretId` |
| `mongodbatlas_service_account` | `orgId` (**ref**) | `clientId` |
| `mongodbatlas_stream_connection` | `workspaceName`, `projectId` (**ref**), `connectionName` | |
| `mongodbatlas_stream_instance` | `projectId` (**ref**), `instanceName` | |
| `mongodbatlas_stream_processor` | `instanceName`, `projectId` (**ref**), `processorName` | |
| `mongodbatlas_stream_workspace` | `projectId` (**ref**), `workspaceName` | |
| `mongodbatlas_team_project_assignment` | `projectId`, `teamId` | |
| `mongodbatlas_team` | `orgId` (**ref**) | `teamId` |
| `mongodbatlas_third_party_integration` | `projectId` (**ref**), `type` | |
| `mongodbatlas_x509_authentication_database_user` | `projectId` (**ref**) | |

## Non-importable

The following resources do not support import:

- `mongodbatlas_cloud_provider_access_authorization`
- `mongodbatlas_cluster_outage_simulation`
- `mongodbatlas_search_index`
- `mongodbatlas_stream_privatelink_endpoint`


## Terraform ID format (for Option B)

When importing via the `crossplane.io/external-name` annotation, the value must be the
full Terraform ID. The format mirrors the identity-parameter list, joined by the
separator below.

| Resource | Annotation value (`crossplane.io/external-name`) |
|----------|--------------------------------------------------|
| `mongodbatlas_advanced_cluster` | `<projectId>-<name>` |
| `mongodbatlas_api_key_project_assignment` | `<projectId>/<apiKeyId>` |
| `mongodbatlas_backup_compliance_policy` | `<projectId>` |
| `mongodbatlas_cloud_backup_schedule` | `<projectId>/<clusterName>` |
| `mongodbatlas_cloud_user_org_assignment` | `<orgId>/<username>` |
| `mongodbatlas_cloud_user_project_assignment` | `<projectId>/<username>` |
| `mongodbatlas_cloud_user_team_assignment` | `<orgId>/<teamId>/<username>` |
| `mongodbatlas_cluster` | `<projectId>-<name>` |
| `mongodbatlas_custom_db_role` | `<projectId>-<roleName>` |
| `mongodbatlas_custom_dns_configuration_cluster_aws` | `<projectId>` |
| `mongodbatlas_database_user` | `<projectId>/<username>/<authDatabaseName>` |
| `mongodbatlas_encryption_at_rest` | `<projectId>` |
| `mongodbatlas_federated_database_instance` | `<projectId>--<name>` |
| `mongodbatlas_federated_query_limit` | `<projectId>--<tenantName>--<limitName>` |
| `mongodbatlas_federated_settings_org_config` | `<federationSettingsId>/<orgId>` |
| `mongodbatlas_flex_cluster` | `<projectId>-<name>` |
| `mongodbatlas_global_cluster_config` | `<projectId>/<clusterName>` |
| `mongodbatlas_ldap_configuration` | `<projectId>` |
| `mongodbatlas_maintenance_window` | `<projectId>` |
| `mongodbatlas_mongodb_employee_access_grant` | `<projectId>/<clusterName>` |
| `mongodbatlas_org_invitation` | `<orgId>-<username>` |
| `mongodbatlas_private_endpoint_regional_mode` | `<projectId>` |
| `mongodbatlas_privatelink_endpoint_service_data_federation_online_archive` | `<projectId>--<endpointId>` |
| `mongodbatlas_project_invitation` | `<projectId>-<username>` |
| `mongodbatlas_push_based_log_export` | `<projectId>` |
| `mongodbatlas_search_deployment` | `<projectId>-<clusterName>` |
| `mongodbatlas_service_account_project_assignment` | `<projectId>/<clientId>` |
| `mongodbatlas_stream_connection` | `<workspaceName>-<projectId>-<connectionName>` |
| `mongodbatlas_stream_instance` | `<projectId>-<instanceName>` |
| `mongodbatlas_stream_processor` | `<instanceName>-<projectId>-<processorName>` |
| `mongodbatlas_stream_workspace` | `<projectId>-<workspaceName>` |
| `mongodbatlas_team_project_assignment` | `<projectId>/<teamId>` |
| `mongodbatlas_third_party_integration` | `<projectId>-<type>` |
| `mongodbatlas_x509_authentication_database_user` | `<projectId>` |

Resources not listed here either use a provider-assigned ID (see the **External Name**
column above — set the annotation to that value) or use a single parameter as ID
(e.g. `mongodbatlas_project` uses `projectId`).

## Examples

### Option A — import by parameters

Create a manifest with:

1. The **required identity parameters** listed in the table above (under `spec.forProvider`).
2. `managementPolicies: ["Observe"]` so Crossplane reads the remote state without modifying it.
3. A `providerConfigRef` pointing to valid Atlas credentials.

```yaml
apiVersion: database.mongodbatlas.m.crossplane.io/v1alpha2
kind: User
metadata:
  name: my-db-user            # any name you choose for the Crossplane resource
spec:
  forProvider:
    authDatabaseName: admin    # must match the existing user's auth database
    projectId: 00001111aaaabbbb55556666
    username: my-user          # must match the existing username in Atlas
  managementPolicies:
    - Observe                  # read-only: Crossplane will not create or modify the user
  providerConfigRef:
    name: default
    kind: ClusterProviderConfig
```

After applying, Crossplane will:
- Render the Terraform ID from the parameters (`00001111aaaabbbb55556666/my-user/admin`)
  and write it to `crossplane.io/external-name` after the first successful observe.
- Populate `status.atProvider` with the full remote state (roles, scopes, etc.).
- Report the resource as `Ready` and `Synced` once the observe succeeds.

### Option B — import by external-name annotation

Set the annotation to the full Terraform ID and omit the identity parameters from
`spec.forProvider`. CRD schema validation still applies, so any field marked
`required` in the generated CRD must be present — check the CRD before relying on
this path for a given resource.

```yaml
apiVersion: database.mongodbatlas.m.crossplane.io/v1alpha2
kind: User
metadata:
  name: my-db-user
  annotations:
    crossplane.io/external-name: 00001111aaaabbbb55556666/my-user/admin
spec:
  forProvider: {}              # template parameters omitted; ID comes from annotation
  managementPolicies:
    - Observe
  providerConfigRef:
    name: default
    kind: ClusterProviderConfig
```

The provider uses the annotation value directly as the Terraform ID. After observe,
`status.atProvider` is populated with the full remote state (which includes
`projectId`, `username`, `authDatabaseName`, etc.).
