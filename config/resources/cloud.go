package resources

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

const groupCloud = "cloud"

func ConfigureCloud(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_backup_compliance_policy", func(r *config.Resource) {
		r.ShortGroup = groupCloud
		r.Kind = "BackupCompliancePolicy"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_backup_schedule", func(r *config.Resource) {
		r.ShortGroup = groupCloud
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			refs.ClusterName: {
				TerraformName: refs.TFCluster,
				Extractor:     refs.ExtractParamPath("name", false),
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_backup_snapshot", func(r *config.Resource) {
		r.ShortGroup = groupCloud
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			refs.ClusterName: {
				TerraformName: refs.TFCluster,
				Extractor:     refs.ExtractParamPath("name", false),
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_backup_snapshot_export_bucket", func(r *config.Resource) {
		r.ShortGroup = groupCloud
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_backup_snapshot_restore_job", func(r *config.Resource) {
		r.ShortGroup = groupCloud
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			refs.ClusterName: {
				TerraformName: refs.TFCluster,
				Extractor:     refs.ExtractParamPath("name", false),
			},
			"snapshot_id": {
				TerraformName: "mongodbatlas_cloud_backup_snapshot",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_backup_snapshot_export_job", func(r *config.Resource) {
		r.ShortGroup = groupCloud
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			refs.ClusterName: {
				TerraformName: refs.TFCluster,
				Extractor:     refs.ExtractParamPath("name", false),
			},
			"snapshot_id": {
				TerraformName: "mongodbatlas_cloud_backup_snapshot",
			},
			"export_bucket_id": {
				TerraformName: "mongodbatlas_cloud_backup_snapshot_export_bucket",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_provider_access_authorization", func(r *config.Resource) {
		r.ShortGroup = groupCloud
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			"role_id": {
				TerraformName: "mongodbatlas_cloud_provider_access_setup",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_provider_access_setup", func(r *config.Resource) {
		r.ShortGroup = groupCloud
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_user_org_assignment", func(r *config.Resource) {
		r.ShortGroup = groupCloud
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_user_project_assignment", func(r *config.Resource) {
		r.ShortGroup = groupCloud
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_user_team_assignment", func(r *config.Resource) {
		r.ShortGroup = groupCloud
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
			"team_id": {
				TerraformName: "mongodbatlas_team",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_serverless_instance", func(r *config.Resource) {
		r.ShortGroup = "serverless"
		r.TerraformResource.DeprecationMessage = "This resource is deprecated. Please use FlexCluster (mongodbatlas_flex_cluster) or AdvancedCluster (mongodbatlas_advanced_cluster) instead."
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})
}
