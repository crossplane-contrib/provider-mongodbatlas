package cloud

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	common "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/common"
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

		// ID format: {project_id}-{cluster_name}-{snapshot_id}
		// cluster_name may contain dashes, but snapshot_id is hex (no dashes),
		// so Split+last is safe here.
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 2, "project_id", "cluster_name")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromSegment("-")
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_backup_snapshot_export_bucket", func(r *config.Resource) {
		r.ShortGroup = group
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}

		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 1, "project_id")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("-", 1, 0)
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

		// ID format: {project_id}-{cluster_name}-{job_id}
		// cluster_name may contain dashes, but job_id is hex (no dashes),
		// so Split+last is safe here.
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 2, "project_id", "cluster_name")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromSegment("-")
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

		// ID format: {project_id}-{cluster_name}-{job_id}
		// cluster_name may contain dashes, but job_id is hex (no dashes),
		// so Split+last is safe here.
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 2, "project_id", "cluster_name")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromSegment("-")
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_provider_access_setup", func(r *config.Resource) {
		r.ShortGroup = group
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}

		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 2, "project_id", "provider_name")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("-", 2, 0)
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
		// ID format: {external_name}-{name}
		// external_name (provider-assigned hex ID) is at position 0 and does
		// not contain dashes. name is the suffix segment to skip from the right.
		r.ExternalName.GetIDFn = common.GetIDFromParamsAndExternalName("-", 0, "name")
		r.ExternalName.GetExternalNameFn = common.ExternalNameFromID("-", 0, 1)
	})
}
