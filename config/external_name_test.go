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

func TestImportJoinedID_DisableNameInitializer(t *testing.T) {
	tests := []struct {
		name        string
		fields      []string
		extNameKey  string
		wantDisable bool
	}{
		{
			name:        "user-provided key in fields",
			fields:      []string{refs.ProjectID, refs.RoleName},
			extNameKey:  refs.RoleName,
			wantDisable: false,
		},
		{
			name:        "provider-assigned key not in fields",
			fields:      []string{refs.ProjectID},
			extNameKey:  "container_id",
			wantDisable: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := importJoinedID(tt.fields, "-", tt.extNameKey)
			assert.Equal(t, tt.wantDisable, e.DisableNameInitializer)
		})
	}
}

func TestImportJoinedID_GetIDFn(t *testing.T) {
	t.Run("user-provided key produces base64 state ID", func(t *testing.T) {
		e := importJoinedID([]string{refs.ProjectID, refs.RoleName}, "-", refs.RoleName)
		params := map[string]any{
			refs.ProjectID: testProjectID,
			refs.RoleName:  "cluster-monitor",
		}
		id, err := e.GetIDFn(context.Background(), "ignored", params, nil)
		require.NoError(t, err)
		decoded := decodeAtlasStateID(id)
		assert.Equal(t, testProjectID, decoded[refs.ProjectID])
		assert.Equal(t, "cluster-monitor", decoded[refs.RoleName])
	})

	t.Run("provider-assigned key included in base64 state ID", func(t *testing.T) {
		e := importJoinedID([]string{refs.ProjectID}, "-", "container_id")
		params := map[string]any{refs.ProjectID: testProjectID}
		id, err := e.GetIDFn(context.Background(), "ctr-abc123", params, nil)
		require.NoError(t, err)
		decoded := decodeAtlasStateID(id)
		assert.Equal(t, testProjectID, decoded[refs.ProjectID])
		assert.Equal(t, "ctr-abc123", decoded["container_id"])
	})

	t.Run("provider-assigned key empty returns empty", func(t *testing.T) {
		e := importJoinedID([]string{refs.ProjectID}, "-", "container_id")
		params := map[string]any{refs.ProjectID: testProjectID}
		id, err := e.GetIDFn(context.Background(), "", params, nil)
		require.NoError(t, err)
		assert.Empty(t, id)
	})
}

func TestImportJoinedID_GetImportIDFn(t *testing.T) {
	t.Run("user-provided key produces plain import ID", func(t *testing.T) {
		e := importJoinedID([]string{refs.ProjectID, refs.RoleName}, "-", refs.RoleName)
		params := map[string]any{
			refs.ProjectID: testProjectID,
			refs.RoleName:  "cluster-monitor",
		}
		id, err := e.GetImportIDFn(context.Background(), "ignored", params, nil)
		require.NoError(t, err)
		assert.Equal(t, testProjectID+"-cluster-monitor", id)
	})

	t.Run("provider-assigned key appended to plain import ID", func(t *testing.T) {
		e := importJoinedID([]string{refs.ProjectID}, "-", "container_id")
		params := map[string]any{refs.ProjectID: testProjectID}
		id, err := e.GetImportIDFn(context.Background(), "ctr-abc123", params, nil)
		require.NoError(t, err)
		assert.Equal(t, testProjectID+"-ctr-abc123", id)
	})

	t.Run("double-dash separator", func(t *testing.T) {
		e := importJoinedID([]string{refs.ProjectID, "tenant_name"}, "--", "limit_name")
		params := map[string]any{refs.ProjectID: testProjectID, "tenant_name": "my-tenant"}
		id, err := e.GetImportIDFn(context.Background(), "bytesPerSecond", params, nil)
		require.NoError(t, err)
		assert.Equal(t, testProjectID+"--my-tenant--bytesPerSecond", id)
	})

	t.Run("provider-assigned key empty returns empty", func(t *testing.T) {
		e := importJoinedID([]string{refs.ProjectID}, "-", "container_id")
		params := map[string]any{refs.ProjectID: testProjectID}
		id, err := e.GetImportIDFn(context.Background(), "", params, nil)
		require.NoError(t, err)
		assert.Empty(t, id)
	})

	t.Run("missing params falls back to valid encoded external name", func(t *testing.T) {
		e := importJoinedID([]string{refs.ProjectID, refs.ProviderName}, "-", refs.PeerID)
		params := map[string]any{refs.ProjectID: testProjectID}
		validID := encodeAtlasStateID(map[string]string{
			refs.ProjectID:    testProjectID,
			refs.ProviderName: "AWS",
			refs.PeerID:       "pcx-abc",
		})
		id, err := e.GetImportIDFn(context.Background(), validID, params, nil)
		require.NoError(t, err)
		assert.Equal(t, validID, id)
	})

	t.Run("missing params rejects raw external name", func(t *testing.T) {
		e := importJoinedID([]string{refs.ProjectID, refs.ProviderName}, "-", refs.PeerID)
		params := map[string]any{refs.ProjectID: testProjectID}
		_, err := e.GetImportIDFn(context.Background(), "my-resource-name", params, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not a valid encoded state ID")
	})
}

func TestImportJoinedIDOrdered_GetIDFn(t *testing.T) {
	t.Run("produces base64 state ID", func(t *testing.T) {
		e := importJoinedIDOrdered(
			[]string{refs.ProjectID, refs.PeerID, refs.ProviderName},
			refs.PeerID,
		)
		params := map[string]any{
			refs.ProjectID:    testProjectID,
			refs.ProviderName: "AWS",
		}
		id, err := e.GetIDFn(context.Background(), "pcx-123", params, nil)
		require.NoError(t, err)
		decoded := decodeAtlasStateID(id)
		assert.Equal(t, testProjectID, decoded[refs.ProjectID])
		assert.Equal(t, "pcx-123", decoded[refs.PeerID])
		assert.Equal(t, "AWS", decoded[refs.ProviderName])
	})

	t.Run("always disables name initializer", func(t *testing.T) {
		e := importJoinedIDOrdered(
			[]string{refs.ProjectID, refs.PeerID, refs.ProviderName},
			refs.PeerID,
		)
		assert.True(t, e.DisableNameInitializer)
	})
}

func TestImportJoinedIDOrdered_GetImportIDFn(t *testing.T) {
	t.Run("provider-assigned key in middle position", func(t *testing.T) {
		e := importJoinedIDOrdered(
			[]string{refs.ProjectID, refs.PeerID, refs.ProviderName},
			refs.PeerID,
		)
		params := map[string]any{
			refs.ProjectID:    testProjectID,
			refs.ProviderName: "AWS",
		}
		id, err := e.GetImportIDFn(context.Background(), "pcx-123", params, nil)
		require.NoError(t, err)
		assert.Equal(t, testProjectID+"-pcx-123-AWS", id)
	})
}

func TestImportJoinedIDMapped_GetIDFn(t *testing.T) {
	t.Run("maps param names to state keys in base64", func(t *testing.T) {
		e := importJoinedIDMapped(
			[]string{refs.ProjectID, refs.Name},
			map[string]string{refs.ProjectID: refs.ProjectID, refs.Name: refs.ClusterName},
		)
		params := map[string]any{
			refs.ProjectID: testProjectID,
			refs.Name:      "my-cluster",
		}
		id, err := e.GetIDFn(context.Background(), "ignored", params, nil)
		require.NoError(t, err)
		decoded := decodeAtlasStateID(id)
		assert.Equal(t, testProjectID, decoded[refs.ProjectID])
		assert.Equal(t, "my-cluster", decoded[refs.ClusterName])
	})

	t.Run("does not disable name initializer when extKey in mapping values", func(t *testing.T) {
		e := importJoinedIDMapped(
			[]string{refs.ProjectID, refs.Name},
			map[string]string{refs.ProjectID: refs.ProjectID, refs.Name: refs.ClusterName},
		)
		assert.False(t, e.DisableNameInitializer)
	})
}

func TestImportJoinedIDMapped_GetImportIDFn(t *testing.T) {
	t.Run("uses param values in param order for plain import ID", func(t *testing.T) {
		e := importJoinedIDMapped(
			[]string{refs.ProjectID, refs.Name},
			map[string]string{refs.ProjectID: refs.ProjectID, refs.Name: refs.ClusterName},
		)
		params := map[string]any{
			refs.ProjectID: testProjectID,
			refs.Name:      "my-cluster",
		}
		id, err := e.GetImportIDFn(context.Background(), "ignored", params, nil)
		require.NoError(t, err)
		assert.Equal(t, testProjectID+"-my-cluster", id)
	})
}

func TestImportJoinedID_GetExternalNameFn(t *testing.T) {
	t.Run("extracts key from base64 encoded state ID", func(t *testing.T) {
		e := importJoinedID([]string{refs.ProjectID, refs.ProviderName}, "-", refs.PeerID)
		stateID := encodeAtlasStateID(map[string]string{
			refs.ProjectID:    testProjectID,
			refs.ProviderName: "AWS",
			refs.PeerID:       "pcx-abc123",
		})
		name, err := e.GetExternalNameFn(map[string]any{"id": stateID})
		require.NoError(t, err)
		assert.Equal(t, "pcx-abc123", name)
	})

	t.Run("returns raw id when key not found in decoded state", func(t *testing.T) {
		e := importJoinedID([]string{refs.ProjectID}, "-", "missing_key")
		stateID := encodeAtlasStateID(map[string]string{refs.ProjectID: testProjectID})
		name, err := e.GetExternalNameFn(map[string]any{"id": stateID})
		require.NoError(t, err)
		assert.Equal(t, stateID, name)
	})

	t.Run("error when id missing from state", func(t *testing.T) {
		e := importJoinedID([]string{refs.ProjectID}, "-", "id")
		_, err := e.GetExternalNameFn(map[string]any{})
		require.Error(t, err)
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
