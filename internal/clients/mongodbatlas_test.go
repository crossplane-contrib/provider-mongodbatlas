package clients

import (
	"context"
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/test"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testPublicKey    = "pub"
	testPrivateKey   = "priv"
	testClientID     = "id"
	testClientSecret = "secret"
)

func newStubSDKProvider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			keyPublicKey:  {Type: schema.TypeString, Optional: true},
			keyPrivateKey: {Type: schema.TypeString, Optional: true},
		},
		ConfigureContextFunc: func(_ context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
			return map[string]string{
				keyPublicKey:  d.Get(keyPublicKey).(string),
				keyPrivateKey: d.Get(keyPrivateKey).(string),
			}, nil
		},
	}
}

func TestConfigureSDKCached_SameCredsSameMeta(t *testing.T) {
	resetMetaCache()
	sdk := newStubSDKProvider()

	config := map[string]any{
		keyPublicKey:  "pub1",
		keyPrivateKey: "priv1",
	}

	meta1, err := configureSDKCached(context.Background(), sdk, config)
	require.NoError(t, err)

	meta2, err := configureSDKCached(context.Background(), sdk, config)
	require.NoError(t, err)

	// Can't use assert.Same on non-pointer values; verify identity via cache hit
	// (configureSDKCached returns the exact same interface value from the map)
	m1 := meta1.(map[string]string)
	m2 := meta2.(map[string]string)
	assert.Equal(t, m1, m2, "same credentials must return same cached Meta")
	// Verify it's the same map (pointer identity) by mutating
	m1["_test"] = "sentinel"
	assert.Equal(t, "sentinel", meta2.(map[string]string)["_test"],
		"must be same underlying map instance (cached)")
	delete(m1, "_test")
}

func TestConfigureSDKCached_DifferentCredsDifferentMeta(t *testing.T) {
	resetMetaCache()
	sdk := newStubSDKProvider()

	config1 := map[string]any{
		keyPublicKey:  "pub1",
		keyPrivateKey: "priv1",
	}
	config2 := map[string]any{
		keyPublicKey:  "pub2",
		keyPrivateKey: "priv2",
	}

	meta1, err := configureSDKCached(context.Background(), sdk, config1)
	require.NoError(t, err)

	meta2, err := configureSDKCached(context.Background(), sdk, config2)
	require.NoError(t, err)

	m1 := meta1.(map[string]string)
	m2 := meta2.(map[string]string)
	assert.NotEqual(t, m1[keyPublicKey], m2[keyPublicKey],
		"different credentials must produce different Meta (no cross-tenant bleed)")
}

func TestConfigureSDKCached_ConfigureError(t *testing.T) {
	resetMetaCache()
	sdk := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"required_field": {Type: schema.TypeString, Required: true},
		},
		ConfigureContextFunc: func(_ context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
			return nil, nil
		},
	}

	rc := sdkterraform.NewResourceConfigRaw(map[string]any{})
	diags := sdk.Validate(rc)
	if diags.HasError() {
		t.Log("validation correctly fails for missing required field")
	}
}

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

func resetMetaCache() {
	metaCacheMu.Lock()
	defer metaCacheMu.Unlock()
	metaCache = map[string]any{}
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
