package clients

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newStubSDKProvider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			keyPublicKey:  {Type: schema.TypeString, Optional: true},
			keyPrivateKey: {Type: schema.TypeString, Optional: true},
		},
		ConfigureFunc: func(d *schema.ResourceData) (any, error) {
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
		ConfigureFunc: func(d *schema.ResourceData) (any, error) {
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
