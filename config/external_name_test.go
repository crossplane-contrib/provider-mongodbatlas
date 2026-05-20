package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/crossplane-contrib/provider-mongodbatlas/config/refs"
)

const testProjectID = "proj123"

func TestEncodeDecodeAtlasStateID(t *testing.T) {
	tests := []struct {
		name   string
		values map[string]string
	}{
		{
			name:   "single key-value",
			values: map[string]string{refs.ProjectID: "abc123"},
		},
		{
			name:   "multiple key-values",
			values: map[string]string{refs.ProjectID: "abc123", refs.RoleName: "myRole"},
		},
		{
			name:   "empty value",
			values: map[string]string{refs.ProjectID: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := encodeAtlasStateID(tt.values)
			decoded := decodeAtlasStateID(encoded)
			assert.Equal(t, tt.values, decoded)
		})
	}
}

func TestEncodeAtlasStateID_SortedKeys(t *testing.T) {
	a := encodeAtlasStateID(map[string]string{"a": "1", "b": "2"})
	b := encodeAtlasStateID(map[string]string{"b": "2", "a": "1"})
	assert.Equal(t, a, b, "key order in input should not affect output")
}

func TestEncodedStateID_DisableNameInitializer(t *testing.T) {
	tests := []struct {
		name        string
		fields      []string
		extNameKey  string
		wantDisable bool
	}{
		{
			name:        "provider-assigned key not in fields",
			fields:      []string{refs.ProjectID, refs.ProviderName, refs.Region},
			extNameKey:  "private_link_id",
			wantDisable: true,
		},
		{
			name:        "user-provided key in fields",
			fields:      []string{refs.ProjectID, refs.RoleName},
			extNameKey:  refs.RoleName,
			wantDisable: false,
		},
		{
			name:        "provider-assigned id",
			fields:      []string{refs.ProjectID},
			extNameKey:  "id",
			wantDisable: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := encodedStateID(tt.fields, tt.extNameKey)
			assert.Equal(t, tt.wantDisable, e.DisableNameInitializer)
		})
	}
}

func TestEncodedStateIDMapped_DisableNameInitializer(t *testing.T) {
	tests := []struct {
		name        string
		mapping     map[string]string
		extNameKey  string
		wantDisable bool
	}{
		{
			name:        "key is a value in mapping",
			mapping:     map[string]string{refs.ProjectID: refs.ProjectID, refs.Name: "cluster_name"},
			extNameKey:  "cluster_name",
			wantDisable: false,
		},
		{
			name:        "key not in mapping values",
			mapping:     map[string]string{refs.ProjectID: refs.ProjectID, refs.ProviderName: refs.ProviderName},
			extNameKey:  "peer_id",
			wantDisable: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := encodedStateIDMapped(tt.mapping, tt.extNameKey)
			assert.Equal(t, tt.wantDisable, e.DisableNameInitializer)
		})
	}
}

func TestEncodedStateID_GetIDFn(t *testing.T) {
	t.Run("all params present with provider-assigned key", func(t *testing.T) {
		e := encodedStateID([]string{refs.ProjectID, refs.ProviderName, refs.Region}, "private_link_id")
		params := map[string]any{
			refs.ProjectID:    testProjectID,
			refs.ProviderName: "AWS",
			refs.Region:       "us-east-1",
		}

		id, err := e.GetIDFn(context.Background(), "vpce-svc-abc", params, nil)
		require.NoError(t, err)

		decoded := decodeAtlasStateID(id)
		assert.Equal(t, testProjectID, decoded[refs.ProjectID])
		assert.Equal(t, "AWS", decoded[refs.ProviderName])
		assert.Equal(t, "us-east-1", decoded[refs.Region])
		assert.Equal(t, "vpce-svc-abc", decoded["private_link_id"])
	})

	t.Run("all params present without external name returns empty when provider-assigned key missing", func(t *testing.T) {
		e := encodedStateID([]string{refs.ProjectID, refs.ProviderName, refs.Region}, "private_link_id")
		params := map[string]any{
			refs.ProjectID:    testProjectID,
			refs.ProviderName: "AWS",
			refs.Region:       "us-east-1",
		}

		id, err := e.GetIDFn(context.Background(), "", params, nil)
		require.NoError(t, err)
		assert.Empty(t, id)
	})

	t.Run("all params present with user-provided key", func(t *testing.T) {
		e := encodedStateID([]string{refs.ProjectID, refs.RoleName}, refs.RoleName)
		params := map[string]any{
			refs.ProjectID: testProjectID,
			refs.RoleName:  "admin",
		}

		id, err := e.GetIDFn(context.Background(), "ignored", params, nil)
		require.NoError(t, err)

		decoded := decodeAtlasStateID(id)
		assert.Equal(t, "admin", decoded[refs.RoleName], "value from params takes precedence over external name")
	})

	t.Run("missing params falls back to valid encoded external name", func(t *testing.T) {
		e := encodedStateID([]string{refs.ProjectID, refs.ProviderName}, "peer_id")
		params := map[string]any{refs.ProjectID: testProjectID}

		validID := encodeAtlasStateID(map[string]string{
			refs.ProjectID:    testProjectID,
			refs.ProviderName: "AWS",
			"peer_id":         "pcx-abc",
		})
		id, err := e.GetIDFn(context.Background(), validID, params, nil)
		require.NoError(t, err)
		assert.Equal(t, validID, id)
	})

	t.Run("missing params rejects raw external name", func(t *testing.T) {
		e := encodedStateID([]string{refs.ProjectID, refs.ProviderName}, "peer_id")
		params := map[string]any{refs.ProjectID: testProjectID}

		_, err := e.GetIDFn(context.Background(), "my-resource-name", params, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not a valid encoded state ID")
	})

	t.Run("missing params and empty external name returns error", func(t *testing.T) {
		e := encodedStateID([]string{refs.ProjectID, refs.ProviderName}, "peer_id")
		params := map[string]any{refs.ProjectID: testProjectID}

		_, err := e.GetIDFn(context.Background(), "", params, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot determine Terraform ID")
	})
}

func TestEncodedStateID_GetExternalNameFn(t *testing.T) {
	t.Run("extracts key from encoded state ID", func(t *testing.T) {
		e := encodedStateID([]string{refs.ProjectID, refs.ProviderName}, "peer_id")
		stateID := encodeAtlasStateID(map[string]string{
			refs.ProjectID:    testProjectID,
			refs.ProviderName: "AWS",
			"peer_id":         "pcx-abc123",
		})

		name, err := e.GetExternalNameFn(map[string]any{"id": stateID})
		require.NoError(t, err)
		assert.Equal(t, "pcx-abc123", name)
	})

	t.Run("returns raw id when key not found in decoded state", func(t *testing.T) {
		e := encodedStateID([]string{refs.ProjectID}, "missing_key")
		stateID := encodeAtlasStateID(map[string]string{refs.ProjectID: testProjectID})

		name, err := e.GetExternalNameFn(map[string]any{"id": stateID})
		require.NoError(t, err)
		assert.Equal(t, stateID, name)
	})

	t.Run("error when id missing from state", func(t *testing.T) {
		e := encodedStateID([]string{refs.ProjectID}, "id")
		_, err := e.GetExternalNameFn(map[string]any{})
		require.Error(t, err)
	})
}

func TestEncodedStateIDMapped_GetIDFn(t *testing.T) {
	t.Run("maps param names to state keys", func(t *testing.T) {
		e := encodedStateIDMapped(
			map[string]string{refs.ProjectID: refs.ProjectID, refs.Name: "cluster_name"},
			"cluster_name",
		)
		params := map[string]any{
			refs.ProjectID: testProjectID,
			refs.Name:      "my-cluster",
		}

		id, err := e.GetIDFn(context.Background(), "ignored", params, nil)
		require.NoError(t, err)

		decoded := decodeAtlasStateID(id)
		assert.Equal(t, testProjectID, decoded[refs.ProjectID])
		assert.Equal(t, "my-cluster", decoded["cluster_name"])
	})
}

func TestHasAllParams(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]any
		fields []string
		want   bool
	}{
		{
			name:   "all present",
			params: map[string]any{"a": "1", "b": "2"},
			fields: []string{"a", "b"},
			want:   true,
		},
		{
			name:   "one missing",
			params: map[string]any{"a": "1"},
			fields: []string{"a", "b"},
			want:   false,
		},
		{
			name:   "present but empty string",
			params: map[string]any{"a": ""},
			fields: []string{"a"},
			want:   false,
		},
		{
			name:   "present but wrong type",
			params: map[string]any{"a": 42},
			fields: []string{"a"},
			want:   false,
		},
		{
			name:   "empty fields",
			params: map[string]any{},
			fields: []string{},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, hasAllParams(tt.params, tt.fields))
		})
	}
}
