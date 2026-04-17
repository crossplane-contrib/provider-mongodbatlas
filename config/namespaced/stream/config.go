package stream

import (
	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/crossplane/upjet/v2/pkg/types/comments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	common "github.com/crossplane-contrib/provider-mongodbatlas/config/namespaced/common"
)

// Configure configures the root group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("mongodbatlas_stream_connection", func(r *config.Resource) {
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
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
			common.PasswordGenerator(
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
			common.PasswordGenerator(
				"spec.forProvider.schemaRegistryAuthentication.passwordSecretRef",
				"spec.forProvider.autoGenerateSchemaRegistryPassword",
			))
	})

	p.AddResourceConfigurator("mongodbatlas_stream_instance", func(r *config.Resource) {
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_stream_instance", func(r *config.Resource) {
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_stream_processor", func(r *config.Resource) {
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})

	p.AddResourceConfigurator("mongodbatlas_stream_workspace", func(r *config.Resource) {
		r.References = config.References{
			"project_id": {
				TerraformName: "mongodbatlas_project",
			},
		}
	})
}
