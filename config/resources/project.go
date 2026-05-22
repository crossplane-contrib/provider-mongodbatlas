package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

func ConfigureProject(p *config.Provider) {
	p.AddResourceConfigurator(refs.TFProject, func(r *config.Resource) {
		r.References = config.References{
			refs.OrgID: {
				TerraformName: refs.TFOrganization,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_third_party_integration", func(r *config.Resource) {
		r.ShortGroup = ""
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_project_invitation", func(r *config.Resource) {
		r.TerraformResource.DeprecationMessage = "This resource is deprecated. Migrate to mongodbatlas_cloud_user_project_assignment for managing project membership."
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"ip_address"},
		}
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_project_ip_access_list", func(r *config.Resource) {
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"ip_address"},
		}
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			project, ok := parameters[refs.ProjectID]
			if !ok {
				return "", errors.New("project_id missing from parameters")
			}
			ip, ok := parameters["ip_address"]
			if !ok {
				ip, ok = parameters["cidr_block"]
				if !ok {
					return "", errors.New("either ip_address or cidr_block parameters must be set")
				}
			}
			return fmt.Sprintf("%s-%s", project, ip), nil
		}
	})

	p.AddResourceConfigurator("mongodbatlas_third_party_integration", func(r *config.Resource) {
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_project_service_account_access_list_entry", func(r *config.Resource) {
		r.ShortGroup = "project"
		r.Kind = "ServiceAccountAccessListEntry"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
		r.ExternalName.GetIDFn = func(_ context.Context, externalName string, parameters map[string]any, setup map[string]any) (string, error) {
			client, ok := parameters["client_id"]
			if !ok {
				return "", errors.New("client_id missing from parameters")
			}
			project, ok := parameters[refs.ProjectID]
			if !ok {
				return "", errors.New("project_id missing from parameters")
			}
			ip, ok := parameters["ip_address"]
			if !ok {
				ip, ok = parameters["cidr_block"]
				if !ok {
					return "", errors.New("either ip_address or cidr_block parameters must be set")
				}
			}
			return fmt.Sprintf("%s-%s-%s", project, client, ip), nil
		}
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromAccessListState(refs.ProjectID)
	})

	p.AddResourceConfigurator("mongodbatlas_project_service_account_secret", func(r *config.Resource) {
		r.ShortGroup = "project"
		r.Kind = "ServiceAccountSecret"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
		r.ExternalName.GetIDFn = refs.GetIDFromParamsAndExternalName("/", 2, refs.ProjectID, "client_id")
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromIDOrState("/", 2, 0, "secret_id")
	})

	p.AddResourceConfigurator("mongodbatlas_project_service_account", func(r *config.Resource) {
		r.ShortGroup = "project"
		r.Kind = "ServiceAccount"
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
		r.ExternalName.GetIDFn = refs.GetIDFromParamsAndExternalName("/", 1, refs.ProjectID)
		r.ExternalName.GetExternalNameFn = refs.ExternalNameFromIDOrState("/", 1, 0, "client_id")
	})
}
