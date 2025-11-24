/*
Copyright 2022 Upbound Inc.
*/

package config

import "github.com/crossplane/upjet/v2/pkg/config"

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
	"mongodbatlas_access_list_api_key":                     idWithStub(),
	"mongodbatlas_advanced_cluster":                        idWithStub(),
	"mongodbatlas_alert_configuration":                     idWithStub(),
	"mongodbatlas_api_key":                                 idWithStub(),
	"mongodbatlas_auditing":                                idWithStub(),
	"mongodbatlas_backup_compliance_policy":                idWithStub(),
	"mongodbatlas_cloud_backup_schedule":                   idWithStub(),
	"mongodbatlas_cloud_backup_snapshot":                   idWithStub(),
	"mongodbatlas_cloud_backup_snapshot_export_bucket":     idWithStub(),
	"mongodbatlas_cloud_backup_snapshot_export_job":        idWithStub(),
	"mongodbatlas_cloud_backup_snapshot_restore_job":       idWithStub(),
	"mongodbatlas_cloud_provider_access":                   idWithStub(),
	"mongodbatlas_cloud_provider_snapshot":                 idWithStub(),
	"mongodbatlas_cloud_provider_snapshot_backup_policy":   idWithStub(),
	"mongodbatlas_cloud_provider_snapshot_restore_job":     idWithStub(),
	"mongodbatlas_cluster":                                 idWithStub(),
	"mongodbatlas_custom_db_role":                          idWithStub(),
	"mongodbatlas_custom_dns_configuration_cluster_aws":    idWithStub(),
	"mongodbatlas_data_lake":                               idWithStub(),
	"mongodbatlas_database_user":                           idWithStub(),
	// "mongodbatlas_encryption_at_rest":                      idWithStub(),
	"mongodbatlas_event_trigger":                           idWithStub(),
	"mongodbatlas_federated_settings_identity_provider":    idWithStub(),
	"mongodbatlas_federated_settings_org_config":           idWithStub(),
	"mongodbatlas_federated_settings_org_role_mapping":     idWithStub(),
	"mongodbatlas_global_cluster_config":                   idWithStub(),
	"mongodbatlas_ldap_configuration":                      idWithStub(),
	"mongodbatlas_ldap_verify":                             idWithStub(),
	"mongodbatlas_maintenance_window":                      idWithStub(),
	"mongodbatlas_network_container":                       idWithStub(),
	"mongodbatlas_network_peering":                         idWithStub(),
	"mongodbatlas_online_archive":                          idWithStub(),
	"mongodbatlas_org_invitation":                          idWithStub(),
	"mongodbatlas_private_endpoint_regional_mode":          idWithStub(),
	"mongodbatlas_private_ip_mode":                         idWithStub(),
	"mongodbatlas_privatelink_endpoint":                    idWithStub(),
	"mongodbatlas_privatelink_endpoint_serverless":         idWithStub(),
	"mongodbatlas_privatelink_endpoint_service":            idWithStub(),
	"mongodbatlas_privatelink_endpoint_service_adl":        idWithStub(),
	"mongodbatlas_privatelink_endpoint_service_serverless": idWithStub(),
	"mongodbatlas_project":                                 idWithStub(),
	"mongodbatlas_project_api_key":                         idWithStub(),
	"mongodbatlas_project_invitation":                      idWithStub(),
	"mongodbatlas_project_ip_access_list":                  idWithStub(),
	"mongodbatlas_search_index":                            idWithStub(),
	"mongodbatlas_serverless_instance":                     idWithStub(),
	// "mongodbatlas_teams":                                   idWithStub(),
	"mongodbatlas_third_party_integration":                 idWithStub(),
	"mongodbatlas_x509_authentication_database_user":       idWithStub(),
}

func idWithStub() config.ExternalName {
	e := config.IdentifierFromProvider
	e.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
		en, _ := config.IDAsExternalName(tfstate)
		return en, nil
	}
	return e
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
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
