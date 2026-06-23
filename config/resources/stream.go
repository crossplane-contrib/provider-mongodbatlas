package resources

import (
	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/types/comments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

func ConfigureStream(p *config.Provider, pwGen func(string, string) config.NewInitializerFn) {
	p.AddResourceConfigurator("mongodbatlas_stream_connection", func(r *config.Resource) {
		r.ExternalName = templated("{{ .parameters.workspace_name }}-{{ .parameters.project_id }}-{{ .parameters.connection_name }}")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			refs.WorkspaceName: {
				TerraformName: refs.TFStreamWorkspace,
				Extractor:     refs.ExtractParamPath(refs.WorkspaceName, false),
			},
		}
		descAuth, _ := comments.New("If true, the authentication password will be auto-generated and"+
			" stored in the Secret referenced by the authentication.passwordSecretRef field.",
			comments.WithTFTag("-"))
		r.TerraformResource.Schema["auto_generate_password"] = &schema.Schema{
			Type:        schema.TypeBool,
			Optional:    true,
			Description: descAuth.String(),
		}
		r.InitializerFns = append(r.InitializerFns,
			pwGen(
				"spec.forProvider.authentication.passwordSecretRef",
				"spec.forProvider.autoGeneratePassword",
			))
		descSR, _ := comments.New("If true, the schema registry authentication password will be auto-generated and"+
			" stored in the Secret referenced by the schemaRegistryAuthentication.passwordSecretRef field.",
			comments.WithTFTag("-"))
		r.TerraformResource.Schema["auto_generate_schema_registry_password"] = &schema.Schema{
			Type:        schema.TypeBool,
			Optional:    true,
			Description: descSR.String(),
		}
		r.InitializerFns = append(r.InitializerFns,
			pwGen(
				"spec.forProvider.schemaRegistryAuthentication.passwordSecretRef",
				"spec.forProvider.autoGenerateSchemaRegistryPassword",
			))
	})

	p.AddResourceConfigurator("mongodbatlas_stream_instance", func(r *config.Resource) {
		r.ExternalName = templated("{{ .parameters.project_id }}-{{ .parameters.instance_name }}")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_stream_processor", func(r *config.Resource) {
		r.ExternalName = templated("{{ .parameters.instance_name }}-{{ .parameters.project_id }}-{{ .parameters.processor_name }}")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
			refs.InstanceName: {
				TerraformName: refs.TFStreamInstance,
				Extractor:     refs.ExtractParamPath(refs.InstanceName, false),
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_stream_workspace", func(r *config.Resource) {
		r.ExternalName = templated("{{ .parameters.project_id }}-{{ .parameters.workspace_name }}")
		r.References = config.References{
			refs.ProjectID: {
				TerraformName: refs.TFProject,
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_stream_privatelink_endpoint", func(r *config.Resource) {
		r.ExternalName = config.IdentifierFromProvider
	})
}
