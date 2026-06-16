package password

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/crossplane/crossplane-runtime/v2/pkg/fieldpath"
	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
)

func TestGetPasswordRef(t *testing.T) {
	tests := []struct {
		name      string
		obj       map[string]interface{}
		path      string
		defaultNS string
		want      *passwordRef
		wantErr   bool
	}{
		{
			name: "ref set with all fields",
			obj: map[string]interface{}{
				"spec": map[string]interface{}{
					"forProvider": map[string]interface{}{
						"passwordSecretRef": map[string]interface{}{
							"name":      "my-secret",
							"namespace": "my-ns",
							"key":       "pw",
						},
					},
				},
			},
			path:      "spec.forProvider.passwordSecretRef",
			defaultNS: "default",
			want:      &passwordRef{name: "my-secret", namespace: "my-ns", key: "pw"},
		},
		{
			name: "ref set with empty namespace falls back to default",
			obj: map[string]interface{}{
				"spec": map[string]interface{}{
					"forProvider": map[string]interface{}{
						"passwordSecretRef": map[string]interface{}{
							"name": "my-secret",
							"key":  "pw",
						},
					},
				},
			},
			path:      "spec.forProvider.passwordSecretRef",
			defaultNS: "fallback-ns",
			want:      &passwordRef{name: "my-secret", namespace: "fallback-ns", key: "pw"},
		},
		{
			name: "ref set with empty key defaults to password",
			obj: map[string]interface{}{
				"spec": map[string]interface{}{
					"forProvider": map[string]interface{}{
						"passwordSecretRef": map[string]interface{}{
							"name":      "my-secret",
							"namespace": "my-ns",
						},
					},
				},
			},
			path:      "spec.forProvider.passwordSecretRef",
			defaultNS: "default",
			want:      &passwordRef{name: "my-secret", namespace: "my-ns", key: "password"},
		},
		{
			name: "ref not set (field missing)",
			obj: map[string]interface{}{
				"spec": map[string]interface{}{
					"forProvider": map[string]interface{}{},
				},
			},
			path:      "spec.forProvider.passwordSecretRef",
			defaultNS: "default",
			want:      nil,
		},
		{
			name: "ref set with empty name",
			obj: map[string]interface{}{
				"spec": map[string]interface{}{
					"forProvider": map[string]interface{}{
						"passwordSecretRef": map[string]interface{}{
							"name": "",
						},
					},
				},
			},
			path:      "spec.forProvider.passwordSecretRef",
			defaultNS: "default",
			want:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paved := fieldpath.Pave(tt.obj)
			got, err := getPasswordRef(paved, tt.path, tt.defaultNS)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReconcileExternalPassword(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))

	t.Run("Mode 2: secret has password key — BYOP skip", func(t *testing.T) {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-secret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"password": []byte("existing-pw"),
				"foo":      []byte("bar"),
			},
		}
		cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

		err := reconcileExternalPassword(context.Background(), cl, &passwordRef{
			name: "my-secret", namespace: "default", key: "password",
		})
		require.NoError(t, err)

		// Secret unchanged
		got := &corev1.Secret{}
		require.NoError(t, cl.Get(context.Background(), types.NamespacedName{Name: "my-secret", Namespace: "default"}, got))
		assert.Equal(t, []byte("existing-pw"), got.Data["password"])
		assert.Equal(t, []byte("bar"), got.Data["foo"])
	})

	t.Run("Mode 3: secret exists without password — generates without ownership", func(t *testing.T) {
		ownerUID := types.UID("xr-owner-uid")
		isController := true
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "xr-secret",
				Namespace: "default",
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: "example.org/v1",
						Kind:       "XR",
						Name:       "my-xr",
						UID:        ownerUID,
						Controller: &isController,
					},
				},
			},
			Data: map[string][]byte{
				"foo": []byte("bar"),
			},
			Type: xpresource.SecretTypeConnection,
		}
		cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

		err := reconcileExternalPassword(context.Background(), cl, &passwordRef{
			name: "xr-secret", namespace: "default", key: "password",
		})
		require.NoError(t, err)

		got := &corev1.Secret{}
		require.NoError(t, cl.Get(context.Background(), types.NamespacedName{Name: "xr-secret", Namespace: "default"}, got))
		// Password generated
		assert.NotEmpty(t, got.Data["password"])
		// Existing data preserved
		assert.Equal(t, []byte("bar"), got.Data["foo"])
		// Owner reference NOT changed — still XR's
		require.Len(t, got.OwnerReferences, 1)
		assert.Equal(t, ownerUID, got.OwnerReferences[0].UID)
	})

	t.Run("Mode 3: custom key", func(t *testing.T) {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-secret",
				Namespace: "default",
			},
			Data: map[string][]byte{
				"other": []byte("data"),
			},
		}
		cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

		err := reconcileExternalPassword(context.Background(), cl, &passwordRef{
			name: "my-secret", namespace: "default", key: "bind_password",
		})
		require.NoError(t, err)

		got := &corev1.Secret{}
		require.NoError(t, cl.Get(context.Background(), types.NamespacedName{Name: "my-secret", Namespace: "default"}, got))
		assert.NotEmpty(t, got.Data["bind_password"])
		assert.Empty(t, got.Data["password"]) // default key NOT written
	})

	t.Run("secret not found — returns error", func(t *testing.T) {
		cl := fake.NewClientBuilder().WithScheme(scheme).Build()

		err := reconcileExternalPassword(context.Background(), cl, &passwordRef{
			name: "nonexistent", namespace: "default", key: "password",
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot get password secret")
	})

	t.Run("secret exists with nil data — generates", func(t *testing.T) {
		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "empty-secret",
				Namespace: "default",
			},
		}
		cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

		err := reconcileExternalPassword(context.Background(), cl, &passwordRef{
			name: "empty-secret", namespace: "default", key: "password",
		})
		require.NoError(t, err)

		got := &corev1.Secret{}
		require.NoError(t, cl.Get(context.Background(), types.NamespacedName{Name: "empty-secret", Namespace: "default"}, got))
		assert.NotEmpty(t, got.Data["password"])
	})
}
