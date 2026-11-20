package clients

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	rtfake "github.com/crossplane/crossplane-runtime/v2/pkg/resource/fake"
	"github.com/crossplane/crossplane-runtime/v2/pkg/test"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	"github.com/google/go-cmp/cmp"
	tpf "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/xpshim"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	namespacedv1beta1 "github.com/crossplane-contrib/provider-mongodbatlas/apis/namespaced/v1beta1"
)

const (
	testPublicKey    = "pub"
	testPrivateKey   = "priv"
	testClientID     = "id"
	testClientSecret = "secret"
)

func TestCredentialHash_Deterministic(t *testing.T) {
	config := map[string]any{
		keyPublicKey:  "pub",
		keyPrivateKey: "priv",
	}
	h1 := credentialHash(config)
	h2 := credentialHash(config)
	assert.Equal(t, h1, h2)
}

func TestCredentialHash_DifferentForDifferentCreds(t *testing.T) {
	c1 := map[string]any{keyPublicKey: "a", keyPrivateKey: "b"}
	c2 := map[string]any{keyPublicKey: "c", keyPrivateKey: "d"}
	assert.NotEqual(t, credentialHash(c1), credentialHash(c2))
}

func TestConfigureCredentials(t *testing.T) {
	type args struct {
		creds map[string]string
	}
	type want struct {
		config map[string]any
		err    error
	}
	cases := map[string]struct {
		args args
		want want
	}{
		"ProgrammaticAPIKey": {
			args: args{creds: map[string]string{keyPublicKey: testPublicKey, keyPrivateKey: testPrivateKey}},
			want: want{config: map[string]any{keyPublicKey: testPublicKey, keyPrivateKey: testPrivateKey}},
		},
		"ServiceAccount": {
			args: args{creds: map[string]string{keyClientID: testClientID, keyClientSecret: testClientSecret}},
			want: want{config: map[string]any{keyClientID: testClientID, keyClientSecret: testClientSecret}},
		},
		"ServiceAccountTakesPrecedence": {
			args: args{creds: map[string]string{
				keyPublicKey: testPublicKey, keyPrivateKey: testPrivateKey,
				keyClientID: testClientID, keyClientSecret: testClientSecret,
			}},
			want: want{config: map[string]any{keyClientID: testClientID, keyClientSecret: testClientSecret}},
		},
		"ServiceAccountMissingSecret": {
			args: args{creds: map[string]string{keyClientID: testClientID}},
			want: want{config: map[string]any{}, err: errors.New(errServiceAccountIncomplete)},
		},
		"ServiceAccountEmptyID": {
			args: args{creds: map[string]string{keyClientID: "", keyClientSecret: testClientSecret}},
			want: want{config: map[string]any{}, err: errors.New(errServiceAccountIncomplete)},
		},
		"ServiceAccountEmptySecret": {
			args: args{creds: map[string]string{keyClientID: testClientID, keyClientSecret: ""}},
			want: want{config: map[string]any{}, err: errors.New(errServiceAccountIncomplete)},
		},
		"ServiceAccountIDAsAPIKey": {
			args: args{creds: map[string]string{keyPublicKey: "mdb_sa_id_123", keyPrivateKey: testPrivateKey}},
			want: want{config: map[string]any{}, err: errors.New(errServiceAccountAsAPIKey)},
		},
		"APIKeyMissingPrivateKey": {
			args: args{creds: map[string]string{keyPublicKey: testPublicKey}},
			want: want{config: map[string]any{}, err: errors.New(errMissingCredentials)},
		},
		"Empty": {
			args: args{creds: map[string]string{}},
			want: want{config: map[string]any{}, err: errors.New(errMissingCredentials)},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := map[string]any{}
			err := configureCredentials(got, tc.args.creds)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("configureCredentials(...): -want error, +got error:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.config, got); diff != "" {
				t.Errorf("configureCredentials(...): -want config, +got config:\n%s", diff)
			}
		})
	}
}

func TestCredentialHash_ServiceAccountKeys(t *testing.T) {
	type args struct {
		a map[string]any
		b map[string]any
	}
	cases := map[string]struct {
		args      args
		wantEqual bool
	}{
		"SameServiceAccount": {
			args: args{
				a: map[string]any{keyClientID: testClientID, keyClientSecret: testClientSecret},
				b: map[string]any{keyClientID: testClientID, keyClientSecret: testClientSecret},
			},
			wantEqual: true,
		},
		"DifferentServiceAccounts": {
			args: args{
				a: map[string]any{keyClientID: "id1", keyClientSecret: "secret1"},
				b: map[string]any{keyClientID: "id2", keyClientSecret: "secret2"},
			},
			wantEqual: false,
		},
		"ServiceAccountVersusAPIKey": {
			args: args{
				a: map[string]any{keyClientID: testClientID, keyClientSecret: testClientSecret},
				b: map[string]any{keyPublicKey: testClientID, keyPrivateKey: testClientSecret},
			},
			wantEqual: false,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			gotEqual := credentialHash(tc.args.a) == credentialHash(tc.args.b)
			if diff := cmp.Diff(tc.wantEqual, gotEqual); diff != "" {
				t.Errorf("credentialHash equality: -want, +got:\n%s", diff)
			}
		})
	}
}

// TestResolveModern asserts that a namespaced managed resource resolves a
// ClusterProviderConfig without its own namespace, and a ProviderConfig from
// its own namespace. The fake client keys objects on their exact namespace,
// like the diff server's in-memory store, so the first case fails without
// the namespace reset in resolveModern.
func TestResolveModern(t *testing.T) {
	const (
		namespace    = "team-a"
		pcName       = "default"
		clusterCreds = "cluster-creds"
		localCreds   = "namespaced-creds"
	)
	scheme := runtime.NewScheme()
	if err := namespacedv1beta1.SchemeBuilder.AddToScheme(scheme); err != nil {
		t.Fatalf("cannot build the test scheme: %v", err)
	}
	credentials := func(secret string) namespacedv1beta1.ProviderCredentials {
		return namespacedv1beta1.ProviderCredentials{
			Source: xpv2.CredentialsSourceSecret,
			CommonCredentialSelectors: xpv2.CommonCredentialSelectors{
				SecretRef: &xpv2.SecretKeySelector{
					SecretReference: xpv2.SecretReference{Name: secret, Namespace: "crossplane-system"},
					Key:             "credentials",
				},
			},
		}
	}

	type want struct {
		secret xpv2.SecretReference
		err    bool
	}
	cases := map[string]struct {
		kind string
		want want
	}{
		"ClusterProviderConfigFromNamespacedResource": {
			kind: namespacedv1beta1.ClusterProviderConfigKind,
			want: want{secret: xpv2.SecretReference{Name: clusterCreds, Namespace: "crossplane-system"}},
		},
		"ProviderConfigFromOwnNamespace": {
			kind: namespacedv1beta1.ProviderConfigKind,
			want: want{secret: xpv2.SecretReference{Name: localCreds, Namespace: namespace}},
		},
		"UnknownKind": {
			kind: "NoSuchProviderConfig",
			want: want{err: true},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			kc := fake.NewClientBuilder().WithScheme(scheme).WithObjects(
				&namespacedv1beta1.ClusterProviderConfig{
					ObjectMeta: metav1.ObjectMeta{Name: pcName},
					Spec:       namespacedv1beta1.ProviderConfigSpec{Credentials: credentials(clusterCreds)},
				},
				&namespacedv1beta1.ProviderConfig{
					ObjectMeta: metav1.ObjectMeta{Name: pcName, Namespace: namespace},
					Spec:       namespacedv1beta1.ProviderConfigSpec{Credentials: credentials(localCreds)},
				},
			).Build()
			mg := &rtfake.ModernManaged{ObjectMeta: metav1.ObjectMeta{Name: "example", Namespace: namespace, UID: "test-uid"}}
			mg.SetProviderConfigReference(&xpv2.ProviderConfigReference{Kind: tc.kind, Name: pcName})

			got, err := resolveModern(context.Background(), kc, mg)
			if (err != nil) != tc.want.err {
				t.Fatalf("resolveModern() error = %v, wantErr %v", err, tc.want.err)
			}
			if err != nil {
				return
			}
			if diff := cmp.Diff(tc.want.secret, got.Credentials.SecretRef.SecretReference); diff != "" {
				t.Errorf("resolveModern() resolved the wrong ProviderConfig (-want +got):\n%s", diff)
			}
		})
	}
}

// testMeta is the client a countingSDKProvider builds: a pointer, so a test can tell a cached client from a new one.
type testMeta struct{ publicKey string }

// countingSDKProvider returns an SDKv2 provider whose Configure counts its calls and returns configure(public_key).
func countingSDKProvider(calls *int, configure func(publicKey string) (any, diag.Diagnostics)) *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			keyPublicKey:  {Type: schema.TypeString, Optional: true},
			keyPrivateKey: {Type: schema.TypeString, Optional: true},
		},
		ConfigureContextFunc: func(_ context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
			*calls++
			return configure(d.Get(keyPublicKey).(string))
		},
	}
}

func newTestMeta(publicKey string) (any, diag.Diagnostics) {
	return &testMeta{publicKey: publicKey}, nil
}

func apiKeyConfig(publicKey string) map[string]any {
	return map[string]any{keyPublicKey: publicKey, keyPrivateKey: testPrivateKey}
}

func TestClientCacheReusesOneClientPerCredentials(t *testing.T) {
	ctx := context.Background()
	calls := 0
	c := newClientCache(countingSDKProvider(&calls, newTestMeta), xpshim.GetFrameworkProvider(), time.Hour)

	a1, err := c.get(ctx, apiKeyConfig("a"))
	require.NoError(t, err)
	a2, err := c.get(ctx, apiKeyConfig("a"))
	require.NoError(t, err)
	b, err := c.get(ctx, apiKeyConfig("b"))
	require.NoError(t, err)

	assert.Same(t, a1, a2, "the same credentials must get the cached provider")
	assert.Equal(t, &testMeta{publicKey: "a"}, a1.client)
	assert.Equal(t, &testMeta{publicKey: "b"}, b.client, "other credentials must get their own client")
	assert.Equal(t, 2, calls, "Configure must run once per credential set")
}

func TestClientCacheExpiresEntries(t *testing.T) {
	ctx := context.Background()
	const ttl = time.Hour
	calls := 0
	now := time.Unix(0, 0)
	c := newClientCache(countingSDKProvider(&calls, newTestMeta), xpshim.GetFrameworkProvider(), ttl)
	c.now = func() time.Time { return now }

	first, err := c.get(ctx, apiKeyConfig("a"))
	require.NoError(t, err)
	_, err = c.get(ctx, apiKeyConfig("rotated"))
	require.NoError(t, err)

	now = now.Add(ttl - time.Second)
	cached, err := c.get(ctx, apiKeyConfig("a"))
	require.NoError(t, err)
	assert.Same(t, first, cached, "an entry younger than the TTL must be reused")

	now = now.Add(time.Second)
	rebuilt, err := c.get(ctx, apiKeyConfig("a"))
	require.NoError(t, err)
	assert.NotSame(t, first, rebuilt, "an entry as old as the TTL must be rebuilt")
	assert.Equal(t, 3, calls)
	assert.Len(t, c.entries, 1, "expired entries of other credentials must be dropped")
}

// TestClientCacheRejectsMissingClient covers credentials the Atlas provider's own validation rejects.
// Its SDKv2 Configure then returns no client and no error,
// which must fail the setup instead of caching a nil client that the SDKv2 resources panic on.
func TestClientCacheRejectsMissingClient(t *testing.T) {
	ctx := context.Background()
	calls := 0
	c := newClientCache(countingSDKProvider(&calls, func(string) (any, diag.Diagnostics) { return nil, nil }), xpshim.GetFrameworkProvider(), time.Hour)

	for range 2 {
		_, err := c.get(ctx, apiKeyConfig("a"))
		require.EqualError(t, err, errNoClient)
	}
	assert.Equal(t, 2, calls, "a missing client must not be cached")
	assert.Empty(t, c.entries)
}

func TestClientCacheConfigureError(t *testing.T) {
	ctx := context.Background()
	calls := 0
	c := newClientCache(countingSDKProvider(&calls, func(string) (any, diag.Diagnostics) { return nil, diag.Errorf("boom") }), xpshim.GetFrameworkProvider(), time.Hour)

	_, err := c.get(ctx, apiKeyConfig("a"))
	require.ErrorContains(t, err, "boom")
	assert.Empty(t, c.entries, "a failed Configure must not be cached")
}

// TestClientCacheRotatedServiceAccountSecret configures the real Atlas provider against a local OAuth token endpoint.
// The upstream token source cache is keyed by client ID only and rejects a second secret until the process restarts.
// The client cache must forget the stale entry, so that a rotated secret,
// or a second ProviderConfig with another secret of the same service account, can configure.
func TestClientCacheRotatedServiceAccountSecret(t *testing.T) {
	ctx := context.Background()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/oauth/token" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"token","token_type":"Bearer","expires_in":3600}`))
	}))
	t.Cleanup(srv.Close)

	// The upstream cache is process-wide, so the client ID is unique to this test.
	const clientID = "mdb_sa_id_rotation_test"
	saConfig := func(secret string) map[string]any {
		return map[string]any{keyClientID: clientID, keyClientSecret: secret, "base_url": srv.URL}
	}
	c := newClientCache(xpshim.GetSDKProvider(), xpshim.GetFrameworkProvider(), time.Hour)

	old, err := c.get(ctx, saConfig("old"))
	require.NoError(t, err)

	// Positive control: the upstream provider on its own rejects the rotated secret.
	diags := xpshim.GetSDKProvider().Configure(ctx, sdkterraform.NewResourceConfigRaw(saConfig("rotated")))
	require.True(t, diags.HasError(), "expected the upstream provider to reject a second secret for the same client ID")
	assert.Contains(t, diags[0].Summary+diags[0].Detail, "service account credentials changed")

	rotated, err := c.get(ctx, saConfig("rotated"))
	require.NoError(t, err, "a rotated secret must configure without a restart")
	assert.NotSame(t, old.client, rotated.client)

	cached, err := c.get(ctx, saConfig("old"))
	require.NoError(t, err)
	assert.Same(t, old, cached, "the client built with the previous secret stays cached until it expires")
}

// TestTerraformSetupBuilder asserts that the setup reads the credentials from the ProviderConfig's Secret
// and hands one client, built once, to both halves of the Atlas provider.
func TestTerraformSetupBuilder(t *testing.T) {
	ctx := context.Background()
	kc := fake.NewClientBuilder().WithScheme(newTestScheme(t)).WithObjects(
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "creds", Namespace: "crossplane-system"},
			Data:       map[string][]byte{"credentials": []byte(`{"public_key":"pub","private_key":"priv"}`)},
		},
		&namespacedv1beta1.ClusterProviderConfig{
			ObjectMeta: metav1.ObjectMeta{Name: offlineTestPCName},
			Spec: namespacedv1beta1.ProviderConfigSpec{Credentials: namespacedv1beta1.ProviderCredentials{
				Source: xpv2.CredentialsSourceSecret,
				CommonCredentialSelectors: xpv2.CommonCredentialSelectors{SecretRef: &xpv2.SecretKeySelector{
					SecretReference: xpv2.SecretReference{Name: "creds", Namespace: "crossplane-system"},
					Key:             "credentials",
				}},
			}},
		},
	).Build()
	upstream := xpshim.GetFrameworkProvider()
	setup := TerraformSetupBuilder(xpshim.GetSDKProvider(), upstream, time.Hour)
	mg := modernManaged(namespacedv1beta1.ClusterProviderConfigKind, offlineTestPCName)

	first, err := setup(ctx, kc, mg)
	require.NoError(t, err)
	second, err := setup(ctx, kc, mg)
	require.NoError(t, err)

	assert.Equal(t, map[string]any{keyPublicKey: "pub", keyPrivateKey: "priv"}, map[string]any(first.Configuration))
	require.NotNil(t, first.Meta)
	assert.Same(t, first.Meta, second.Meta, "the client must be built once per credential set")

	p, ok := first.FrameworkProvider.(*configuredProvider)
	require.True(t, ok, "expected the upstream provider wrapped by configuredProvider, got %T", first.FrameworkProvider)
	assert.Equal(t, upstream, p.Provider)
	resp := tpf.ConfigureResponse{}
	p.Configure(ctx, tpf.ConfigureRequest{}, &resp)
	assert.Same(t, first.Meta, resp.ResourceData, "the framework resources must get the SDKv2 client")
	assert.Same(t, first.Meta, resp.DataSourceData)
}
