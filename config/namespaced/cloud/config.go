package cloud

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

const group = "cloud"

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_backup_compliance_policy", func(r *config.Resource) {
		r.ShortGroup = group
		r.Kind = "BackupCompliancePolicy"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_backup_schedule", func(r *config.Resource) {
		r.ShortGroup = group
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})
	p.AddResourceConfigurator("mongodbatlas_cloud_backup_snapshot", func(r *config.Resource) {
		r.ShortGroup = group
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_backup_snapshot_export_bucket", func(r *config.Resource) {
		r.ShortGroup = group
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_backup_snapshot_restore_job", func(r *config.Resource) {
		r.ShortGroup = group
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
			"snapshot_id": {
				TerraformName: "mongodbatlas_cloud_backup_snapshot",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_backup_snapshot_export_job", func(r *config.Resource) {
		r.ShortGroup = group
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
			"snapshot_id": {
				TerraformName: "mongodbatlas_cloud_backup_snapshot",
			},
			"export_bucket_id": {
				TerraformName: "mongodbatlas_cloud_backup_snapshot_export_bucket",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_provider_access_setup", func(r *config.Resource) {
		r.ShortGroup = group
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_user_org_assignment", func(r *config.Resource) {
		r.ShortGroup = group
		r.References = config.References{
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_user_project_assignment", func(r *config.Resource) {
		r.ShortGroup = group
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_user_team_assignment", func(r *config.Resource) {
		r.ShortGroup = group
		r.References = config.References{
			"org_id": {
				TerraformName: "mongodbatlas_organization",
			},
			"team_id": {
				TerraformName: "mongodbatlas_team",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_serverless_instance", func(r *config.Resource) {
		r.ShortGroup = "serverless"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})
}
