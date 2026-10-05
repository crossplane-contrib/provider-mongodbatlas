package config

import (
	"context"
	"encoding/base64"
	"maps"
	"slices"
	"strings"
	"testing"

	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

// TestPlainProviderIDResources covers plugin-framework resources whose
// Terraform "id" attribute is the plain provider-assigned ID. Read passes
// "id" to the Atlas API as is, so GetIDFn must not add parameter prefixes,
// and GetExternalNameFn must accept an ID without separators. The
// "{parent}-{id}" format in the Terraform docs is the import ID only;
// ImportState splits it and stores only the provider-assigned part in "id".
func TestPlainProviderIDResources(t *testing.T) {
	const atlasID = "66f1c018dba9c04e7dcfaf36"
	cases := map[string]map[string]any{
		"mongodbatlas_resource_policy": {
			"org_id": "65def6ce0f722a1507105aa5",
			"name":   "policy",
		},
		"mongodbatlas_encryption_at_rest_private_endpoint": {
			refs.ProjectID:   "65def6ce0f722a1507105aa5",
			"cloud_provider": "AZURE",
			"region_name":    "US_EAST_2",
		},
	}
	p := GetProviderNamespaced()
	for name, params := range cases {
		t.Run(name, func(t *testing.T) {
			e := p.Resources[name].ExternalName

			id, err := e.GetIDFn(context.Background(), atlasID, params, map[string]any{})
			if err != nil {
				t.Fatalf("GetIDFn: unexpected error: %v", err)
			}
			if id != atlasID {
				t.Errorf("GetIDFn: want %q, got %q", atlasID, id)
			}

			state := map[string]any{"id": atlasID}
			maps.Copy(state, params)
			extName, err := e.GetExternalNameFn(state)
			if err != nil {
				t.Fatalf("GetExternalNameFn: unexpected error: %v", err)
			}
			if extName != atlasID {
				t.Errorf("GetExternalNameFn: want %q, got %q", atlasID, extName)
			}
		})
	}
}

// encodeStateID mirrors conversion.EncodeStateID of the Atlas Terraform
// provider: sorted base64 "key:value" pairs joined by "-".
func encodeStateID(values map[string]string) string {
	enc := func(s string) string { return base64.StdEncoding.EncodeToString([]byte(s)) }
	parts := make([]string, 0, len(values))
	for _, k := range slices.Sorted(maps.Keys(values)) {
		parts = append(parts, enc(k)+":"+enc(values[k]))
	}
	return strings.Join(parts, "-")
}

// TestProjectAPIKeyID covers mongodbatlas_project_api_key. Its project IDs
// are only inside project_assignment, and Terraform writes and reads a state
// ID that holds only api_key_id.
func TestProjectAPIKeyID(t *testing.T) {
	e := GetProviderNamespaced().Resources["mongodbatlas_project_api_key"].ExternalName
	params := map[string]any{
		"description":        "key",
		"project_assignment": []any{map[string]any{refs.ProjectID: "p1", "role_names": []any{"GROUP_READ_ONLY"}}},
	}
	stateID := encodeStateID(map[string]string{"api_key_id": "key123"})

	id, err := e.GetIDFn(context.Background(), "", params, map[string]any{})
	if err != nil || id != "" {
		t.Errorf("GetIDFn before create: want empty ID and no error, got %q, %v", id, err)
	}
	id, err = e.GetIDFn(context.Background(), "key123", params, map[string]any{})
	if err != nil || id != stateID {
		t.Errorf("GetIDFn: want %q, got %q, %v", stateID, id, err)
	}
	name, err := e.GetExternalNameFn(map[string]any{"id": stateID})
	if err != nil || name != "key123" {
		t.Errorf("GetExternalNameFn: want %q, got %q, %v", "key123", name, err)
	}
}

// TestCloudUserTeamAssignmentID covers mongodbatlas_cloud_user_team_assignment.
// user_id is the required argument; username is computed, so it is not in
// forProvider before create.
func TestCloudUserTeamAssignmentID(t *testing.T) {
	e := GetProviderNamespaced().Resources["mongodbatlas_cloud_user_team_assignment"].ExternalName
	params := map[string]any{"org_id": "o1", "team_id": "t1", "user_id": "u1"}

	id, err := e.GetIDFn(context.Background(), "", params, map[string]any{})
	if err != nil || id != "o1/t1/u1" {
		t.Errorf("GetIDFn: want %q, got %q, %v", "o1/t1/u1", id, err)
	}
	state := map[string]any{"username": "user@example.com"}
	maps.Copy(state, params)
	name, err := e.GetExternalNameFn(state)
	if err != nil || name != "o1/t1/u1" {
		t.Errorf("GetExternalNameFn: want %q, got %q, %v", "o1/t1/u1", name, err)
	}
}

// TestComputedKeyFrameworkResources covers plugin-framework resources with no
// "id" attribute whose Read path needs a provider-assigned computed attribute.
// upjet never calls GetIDFn for them, so:
//   - GetExternalNameFn must read the key from the state, not from "id";
//   - SetIdentifierArgumentFn must copy the external name into the key, so
//     that an import with only the annotation can Read;
//   - an unset key must report an empty state, so that Create runs.
func TestComputedKeyFrameworkResources(t *testing.T) {
	cases := map[string]string{
		"mongodbatlas_ai_model_api_key":                    "api_key_id",
		"mongodbatlas_cloud_backup_collection_restore_job": "job_id",
		"mongodbatlas_log_integration":                     "integration_id",
		"mongodbatlas_metric_integration":                  "metric_integration_id",
		"mongodbatlas_stream_connection_failover":          "failover_connection_id",
	}
	p := GetProviderNamespaced()
	for name, key := range cases {
		t.Run(name, func(t *testing.T) {
			r := p.Resources[name]
			e := r.ExternalName

			got, err := e.GetExternalNameFn(map[string]any{refs.ProjectID: "p1", key: "k1"})
			if err != nil || got != "k1" {
				t.Errorf("GetExternalNameFn: want %q, got %q, %v", "k1", got, err)
			}

			params := map[string]any{refs.ProjectID: "p1"}
			e.SetIdentifierArgumentFn(params, "k1")
			if params[key] != "k1" {
				t.Errorf("SetIdentifierArgumentFn: want %s=%q, got %v", key, "k1", params[key])
			}
			// The state rebuilt from status.atProvider already holds the real
			// key. An older external name (log_integration used "type") must
			// not replace it.
			params = map[string]any{refs.ProjectID: "p1", key: "real"}
			e.SetIdentifierArgumentFn(params, "S3")
			if params[key] != "real" {
				t.Errorf("SetIdentifierArgumentFn with %s set: want %q kept, got %v", key, "real", params[key])
			}
			params = map[string]any{refs.ProjectID: "p1"}
			e.SetIdentifierArgumentFn(params, "")
			if _, ok := params[key]; ok {
				t.Errorf("SetIdentifierArgumentFn with empty external name: want no %s, got %v", key, params[key])
			}

			if r.TerraformPluginFrameworkIsStateEmptyFn == nil {
				t.Fatal("TerraformPluginFrameworkIsStateEmptyFn: want set, got nil")
			}
			st := tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{key: tftypes.String}},
				map[string]tftypes.Value{key: tftypes.NewValue(tftypes.String, nil)})
			empty, err := r.TerraformPluginFrameworkIsStateEmptyFn(context.Background(), st, rschema.Schema{})
			if err != nil || !empty {
				t.Errorf("TerraformPluginFrameworkIsStateEmptyFn with null %s: want true, got %v, %v", key, empty, err)
			}
		})
	}
}

// TestAIModelRateLimitExternalName covers mongodbatlas_ai_model_rate_limit,
// a plugin-framework resource with no "id" attribute. The external name is
// the user-set model_group_name.
func TestAIModelRateLimitExternalName(t *testing.T) {
	e := GetProviderNamespaced().Resources["mongodbatlas_ai_model_rate_limit"].ExternalName
	name, err := e.GetExternalNameFn(map[string]any{
		refs.ProjectID: "p1", "cloud": "AWS", "geography": "US", "model_group_name": "g1",
	})
	if err != nil || name != "g1" {
		t.Errorf("GetExternalNameFn: want %q, got %q, %v", "g1", name, err)
	}
}

// TestStateIDKeysMatchTerraformRead covers SDK resources whose Terraform Read
// decodes more keys from the state ID than the external name holds.
func TestStateIDKeysMatchTerraformRead(t *testing.T) {
	p := GetProviderNamespaced()

	t.Run("mongodbatlas_x509_authentication_database_user with username", func(t *testing.T) {
		e := p.Resources["mongodbatlas_x509_authentication_database_user"].ExternalName
		id, err := e.GetIDFn(context.Background(), "p1", map[string]any{refs.ProjectID: "p1", "username": "u1"}, map[string]any{})
		want := encodeStateID(map[string]string{refs.ProjectID: "p1", "username": "u1"})
		if err != nil || id != want {
			t.Errorf("GetIDFn: want %q, got %q, %v", want, id, err)
		}
	})

	t.Run("mongodbatlas_x509_authentication_database_user without username", func(t *testing.T) {
		e := p.Resources["mongodbatlas_x509_authentication_database_user"].ExternalName
		id, err := e.GetIDFn(context.Background(), "p1", map[string]any{refs.ProjectID: "p1"}, map[string]any{})
		want := encodeStateID(map[string]string{refs.ProjectID: "p1"})
		if err != nil || id != want {
			t.Errorf("GetIDFn: want %q, got %q, %v", want, id, err)
		}
	})

	t.Run("mongodbatlas_cluster", func(t *testing.T) {
		e := p.Resources["mongodbatlas_cluster"].ExternalName
		params := map[string]any{refs.ProjectID: "p1", "name": "c1", "provider_name": "TENANT"}
		id, err := e.GetIDFn(context.Background(), "c1", params, map[string]any{})
		want := encodeStateID(map[string]string{refs.ProjectID: "p1", "cluster_name": "c1", "provider_name": "TENANT"})
		if err != nil || id != want {
			t.Errorf("GetIDFn: want %q, got %q, %v", want, id, err)
		}
		name, err := e.GetExternalNameFn(map[string]any{"id": want})
		if err != nil || name != "c1" {
			t.Errorf("GetExternalNameFn: want %q, got %q, %v", "c1", name, err)
		}
	})
}

// TestLegacyEncodedExternalName covers objects created on v1.1.x, where the
// external name was the full encoded state ID. GetIDFn must not encode it a
// second time.
func TestLegacyEncodedExternalName(t *testing.T) {
	p := GetProviderNamespaced()
	cases := map[string]struct {
		params map[string]any
		state  map[string]string
	}{
		"mongodbatlas_team": {
			params: map[string]any{"org_id": "o1", "name": "team"},
			state:  map[string]string{"org_id": "o1", "id": "t1"},
		},
		"mongodbatlas_project_api_key": {
			params: map[string]any{},
			state:  map[string]string{"api_key_id": "key123"},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			legacy := encodeStateID(tc.state)
			id, err := p.Resources[name].ExternalName.GetIDFn(context.Background(), legacy, tc.params, map[string]any{})
			if err != nil || id != legacy {
				t.Errorf("GetIDFn: want %q, got %q, %v", legacy, id, err)
			}
		})
	}
}
