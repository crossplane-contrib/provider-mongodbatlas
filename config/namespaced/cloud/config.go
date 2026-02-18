package cloud

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_cloud_backup_schedule", func(r *config.Resource) {
		r.ShortGroup = "cloud"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})
	p.AddResourceConfigurator("mongodbatlas_cloud_backup_snapshot", func(r *config.Resource) {
		r.ShortGroup = "cloud"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}

		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			project, ok := parameters["project_id"]
			if !ok {
				return "", errors.New("project_id missing from parameters")
			}
			cluster, ok := parameters["cluster_name"]
			if !ok {
				return "", errors.New("cluster_name missing from parameters")
			}
			return fmt.Sprintf("%s-%s-%s", project, cluster, externalName), nil
		}

		r.ExternalName.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
			id, ok := tfstate["id"]
			if !ok {
				return "", errors.New("id attribute missing from state file")
			}

			idStr, ok := id.(string)
			if !ok {
				return "", errors.New("value of id needs to be string")
			}

			idSlice := strings.Split(idStr, "-")
			return idSlice[2], nil
		}
	})
	p.AddResourceConfigurator("mongodbatlas_cloud_backup_snapshot_export_bucket", func(r *config.Resource) {
		r.ShortGroup = "cloud"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}

		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			project, ok := parameters["project_id"]
			if !ok {
				return "", errors.New("project_id missing from parameters")
			}
			return fmt.Sprintf("%s-%s", project, externalName), nil
		}

		r.ExternalName.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
			id, ok := tfstate["id"]
			if !ok {
				return "", errors.New("id attribute missing from state file")
			}

			idStr, ok := id.(string)
			if !ok {
				return "", errors.New("value of id needs to be string")
			}

			idSlice := strings.Split(idStr, "-")
			return idSlice[1], nil
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_backup_snapshot_restore_job", func(r *config.Resource) {
		r.ShortGroup = "cloud"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
			"snapshot_id": {
				TerraformName: "mongodbatlas_cloud_backup_snapshot",
			},
		}

		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			project, ok := parameters["project_id"]
			if !ok {
				return "", errors.New("project_id missing from parameters")
			}
			cluster, ok := parameters["cluster_name"]
			if !ok {
				return "", errors.New("cluster_name missing from parameters")
			}
			return fmt.Sprintf("%s-%s-%s", project, cluster, externalName), nil
		}

		r.ExternalName.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
			id, ok := tfstate["id"]
			if !ok {
				return "", errors.New("id attribute missing from state file")
			}

			idStr, ok := id.(string)
			if !ok {
				return "", errors.New("value of id needs to be string")
			}

			idSlice := strings.Split(idStr, "-")
			return idSlice[2], nil
		}
	})

	p.AddResourceConfigurator("mongodbatlas_cloud_backup_snapshot_export_job", func(r *config.Resource) {
		r.ShortGroup = "cloud"
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

		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			project, ok := parameters["project_id"]
			if !ok {
				return "", errors.New("project_id missing from parameters")
			}
			cluster, ok := parameters["cluster_name"]
			if !ok {
				return "", errors.New("cluster_name missing from parameters")
			}

			return fmt.Sprintf("%s-%s-%s", project, cluster, externalName), nil
		}

		r.ExternalName.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
			id, ok := tfstate["id"]
			if !ok {
				return "", errors.New("id attribute missing from state file")
			}

			idStr, ok := id.(string)
			if !ok {
				return "", errors.New("value of id needs to be string")
			}

			idSlice := strings.Split(idStr, "-")
			return idSlice[2], nil
		}
	})
	p.AddResourceConfigurator("mongodbatlas_cloud_provider_access_setup", func(r *config.Resource) {
		r.ShortGroup = "cloud"
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}

		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			project, ok := parameters["project_id"]
			if !ok {
				return "", errors.New("project_id missing from parameters")
			}
			provider, ok := parameters["provider_name"]
			if !ok {
				return "", errors.New("provider_name missing from parameters")
			}
			return fmt.Sprintf("%s-%s-%s", project, provider, externalName), nil
		}

		r.ExternalName.GetExternalNameFn = func(tfstate map[string]any) (string, error) {
			id, ok := tfstate["id"]
			if !ok {
				return "", errors.New("id attribute missing from state file")
			}

			idStr, ok := id.(string)
			if !ok {
				return "", errors.New("value of id needs to be string")
			}

			idSlice := strings.Split(idStr, "-")
			return idSlice[2], nil
		}
	})
}
