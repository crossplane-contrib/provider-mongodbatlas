/*
Copyright 2021 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package clients

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"maps"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	tpf "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/xpshim"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/terraform"

	clusterv1beta1 "github.com/crossplane-contrib/provider-mongodbatlas/apis/cluster/v1beta1"
	namespacedv1beta1 "github.com/crossplane-contrib/provider-mongodbatlas/apis/namespaced/v1beta1"
)

// Credential JSON keys. They match the Terraform provider's configuration
// attributes so the values can be copied through unchanged.
const (
	keyPublicKey    = "public_key"
	keyPrivateKey   = "private_key"
	keyClientID     = "client_id"
	keyClientSecret = "client_secret"

	// serviceAccountPrefix starts every Atlas service account client ID.
	serviceAccountPrefix = "mdb_sa"
)

// credentialKeys lists every key that identifies a credential set, in a fixed
// order so credentialHash stays deterministic.
var credentialKeys = []string{keyPublicKey, keyPrivateKey, keyClientID, keyClientSecret}

const (
	// error messages
	errResolveProviderConfig    = "cannot resolve provider config"
	errNoProviderConfig         = "no providerConfigRef provided"
	errGetProviderConfig        = "cannot get referenced ProviderConfig"
	errTrackUsage               = "cannot track ProviderConfig usage"
	errExtractCredentials       = "cannot extract credentials"
	errUnmarshalCredentials     = "cannot unmarshal mongodbatlas credentials as JSON"
	errConfigureCredentials     = "cannot configure mongodbatlas credentials"
	errServiceAccountIncomplete = "service account credentials require both client_id and client_secret"
	errMissingCredentials       = "credentials must contain client_id and client_secret (service account) or public_key and private_key (programmatic API key)"
	errServiceAccountAsAPIKey   = "public_key holds a service account client ID (it starts with " + serviceAccountPrefix + "): put service account credentials in client_id and client_secret"
	errConfigureSDK             = "cannot configure SDK provider"
	errNoClient                 = "the MongoDB Atlas Terraform provider accepted the credentials but built no client"
)

// TerraformSetupBuilder returns a SetupFn for the no-fork (in-process) architecture.
// Each credential set gets one Atlas client, built by the SDKv2 provider's Configure and kept for cacheTTL.
// The client is Setup.Meta for the SDKv2 resources, and Setup.FrameworkProvider hands it to the framework resources,
// so the framework Configure that upjet calls on every reconcile builds nothing.
func TerraformSetupBuilder(sdk *schema.Provider, fw tpf.Provider, cacheTTL time.Duration) terraform.SetupFn {
	clients := newClientCache(sdk, fw, cacheTTL)
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{}

		pcSpec, err := resolveProviderConfig(ctx, client, mg)
		if err != nil {
			return terraform.Setup{}, errors.Wrap(err, errResolveProviderConfig)
		}

		data, err := resource.CommonCredentialExtractor(ctx, pcSpec.Credentials.Source, client, pcSpec.Credentials.CommonCredentialSelectors)
		if err != nil {
			return ps, errors.Wrap(err, errExtractCredentials)
		}
		creds := map[string]string{}
		if err := json.Unmarshal(data, &creds); err != nil {
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}

		ps.Configuration = map[string]any{}
		if err := configureCredentials(ps.Configuration, creds); err != nil {
			return terraform.Setup{}, errors.Wrap(err, errConfigureCredentials)
		}

		p, err := clients.get(ctx, ps.Configuration)
		if err != nil {
			return terraform.Setup{}, errors.Wrap(err, errConfigureSDK)
		}
		ps.Meta = p.client
		ps.FrameworkProvider = p

		return ps, nil
	}
}

// configureCredentials copies the Atlas credentials from the credentials JSON
// into the Terraform provider configuration. The authentication method is
// chosen from the keys present: service account credentials (client_id and
// client_secret) take precedence over programmatic API keys (public_key and
// private_key), mirroring the Terraform provider's own selection.
func configureCredentials(config map[string]any, creds map[string]string) error {
	clientID, clientSecret := creds[keyClientID], creds[keyClientSecret]
	hasID, hasSecret := clientID != "", clientSecret != ""
	switch {
	case hasID && hasSecret:
		config[keyClientID] = clientID
		config[keyClientSecret] = clientSecret
		return nil
	case hasID || hasSecret:
		return errors.New(errServiceAccountIncomplete)
	}

	publicKey, privateKey := creds[keyPublicKey], creds[keyPrivateKey]
	if publicKey == "" || privateKey == "" {
		return errors.New(errMissingCredentials)
	}
	if strings.HasPrefix(publicKey, serviceAccountPrefix) {
		return errors.New(errServiceAccountAsAPIKey)
	}
	config[keyPublicKey] = publicKey
	config[keyPrivateKey] = privateKey
	return nil
}

// clientCache holds one configured provider per credential set, keyed by credentialHash.
// It drops entries older than ttl, so a rotated or deleted credential set does not keep its client forever.
type clientCache struct {
	sdk *schema.Provider
	fw  tpf.Provider
	ttl time.Duration
	now func() time.Time

	// mu also serializes Configure: the SDKv2 provider keeps its meta in the shared instance,
	// and a stale service account must not be forgotten while another credential set configures the same client ID.
	mu      sync.Mutex
	entries map[string]cachedProvider
}

type cachedProvider struct {
	provider  *configuredProvider
	createdAt time.Time
}

func newClientCache(sdk *schema.Provider, fw tpf.Provider, ttl time.Duration) *clientCache {
	return &clientCache{sdk: sdk, fw: fw, ttl: ttl, now: time.Now, entries: map[string]cachedProvider{}}
}

// get returns the provider configured with config, configuring the SDKv2 provider on a cache miss.
func (c *clientCache) get(ctx context.Context, config map[string]any) (*configuredProvider, error) {
	h := credentialHash(config)

	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	maps.DeleteFunc(c.entries, func(_ string, e cachedProvider) bool { return now.Sub(e.createdAt) >= c.ttl })
	if e, ok := c.entries[h]; ok {
		return e.provider, nil
	}

	if id, ok := config[keyClientID].(string); ok {
		secret, _ := config[keyClientSecret].(string)
		xpshim.ForgetStaleServiceAccount(id, secret)
	}
	if diags := c.sdk.Configure(ctx, tf.NewResourceConfigRaw(config)); diags.HasError() {
		return nil, fmt.Errorf("failed to configure SDK provider: %v", diags)
	}
	// Credentials the Atlas provider's own validation rejects give no client and no error here:
	// only the framework half of the provider reports them.
	meta := c.sdk.Meta()
	if meta == nil {
		return nil, errors.New(errNoClient)
	}

	p := &configuredProvider{Provider: c.fw, client: meta}
	c.entries[h] = cachedProvider{provider: p, createdAt: now}
	return p, nil
}

// configuredProvider is the upstream framework provider, with a Configure that hands out a client built beforehand.
// The upstream Configure reads the configuration block, the MONGODB_ATLAS_* environment and, with an AWS role,
// AWS Secrets Manager; it builds a new Atlas client every time upjet calls it, which is on every reconcile.
// Everything else, the schema, the resources and the provider meta schema in particular, is upstream's own.
type configuredProvider struct {
	tpf.Provider
	client any
}

// Configure hands the prepared client to the resources and the data sources.
// Ephemeral resources, which upjet never opens, get no provider data.
func (p *configuredProvider) Configure(_ context.Context, _ tpf.ConfigureRequest, resp *tpf.ConfigureResponse) {
	resp.DataSourceData = p.client
	resp.ResourceData = p.client
}

// MetaSchema serves the upstream provider's meta schema, which its resources read on every operation.
func (p *configuredProvider) MetaSchema(ctx context.Context, req tpf.MetaSchemaRequest, resp *tpf.MetaSchemaResponse) {
	if mp, ok := p.Provider.(tpf.ProviderWithMetaSchema); ok {
		mp.MetaSchema(ctx, req, resp)
	}
}

// EphemeralResources serves the upstream ephemeral resources, so the wrapper keeps the upstream protocol surface.
func (p *configuredProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	if ep, ok := p.Provider.(tpf.ProviderWithEphemeralResources); ok {
		return ep.EphemeralResources(ctx)
	}
	return nil
}

func credentialHash(config map[string]any) string {
	h := sha256.New()
	for _, key := range credentialKeys {
		if v, ok := config[key]; ok {
			_, _ = fmt.Fprintf(h, "%s=%v;", key, v)
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func toSharedPCSpec(pc *clusterv1beta1.ProviderConfig) (*namespacedv1beta1.ProviderConfigSpec, error) {
	if pc == nil {
		return nil, nil
	}
	data, err := json.Marshal(pc.Spec)
	if err != nil {
		return nil, err
	}

	var mSpec namespacedv1beta1.ProviderConfigSpec
	err = json.Unmarshal(data, &mSpec)
	return &mSpec, err
}

func resolveProviderConfig(ctx context.Context, crClient client.Client, mg resource.Managed) (*namespacedv1beta1.ProviderConfigSpec, error) {
	switch managed := mg.(type) {
	case resource.LegacyManaged: //nolint: staticcheck
		return resolveLegacy(ctx, crClient, managed)
	case resource.ModernManaged:
		return resolveModern(ctx, crClient, managed)
	default:
		return nil, errors.New("resource is not a managed resource")
	}
}

func resolveLegacy(ctx context.Context, client client.Client, mg resource.LegacyManaged) (*namespacedv1beta1.ProviderConfigSpec, error) { //nolint: staticcheck
	configRef := mg.GetProviderConfigReference()
	if configRef == nil {
		return nil, errors.New(errNoProviderConfig)
	}
	pc := &clusterv1beta1.ProviderConfig{}
	if err := client.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetProviderConfig)
	}

	t := resource.NewLegacyProviderConfigUsageTracker(client, &clusterv1beta1.ProviderConfigUsage{})
	if err := t.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackUsage)
	}

	return toSharedPCSpec(pc)
}

func resolveModern(ctx context.Context, crClient client.Client, mg resource.ModernManaged) (*namespacedv1beta1.ProviderConfigSpec, error) {
	configRef := mg.GetProviderConfigReference()
	if configRef == nil {
		return nil, errors.New(errNoProviderConfig)
	}

	pcRuntimeObj, err := crClient.Scheme().New(namespacedv1beta1.SchemeGroupVersion.WithKind(configRef.Kind))
	if err != nil {
		return nil, errors.Wrap(err, "unknown GVK for ProviderConfig")
	}
	pcObj, ok := pcRuntimeObj.(client.Object)
	if !ok {
		// This indicates a programming error, types are not properly generated
		return nil, errors.New(" is not an Object")
	}

	// A ClusterProviderConfig is cluster-scoped. The API server ignores a
	// namespace on its lookup, but the diff server reads it from an in-memory
	// store keyed on the exact namespace, so it must be cleared here rather
	// than left to the client.
	ns := mg.GetNamespace()
	if configRef.Kind == namespacedv1beta1.ClusterProviderConfigKind {
		ns = ""
	}
	if err := crClient.Get(ctx, types.NamespacedName{Name: configRef.Name, Namespace: ns}, pcObj); err != nil {
		return nil, errors.Wrap(err, errGetProviderConfig)
	}

	var pcSpec namespacedv1beta1.ProviderConfigSpec
	pcu := &namespacedv1beta1.ProviderConfigUsage{}
	switch pc := pcObj.(type) {
	case *namespacedv1beta1.ProviderConfig:
		pcSpec = pc.Spec
		if pcSpec.Credentials.SecretRef != nil {
			pcSpec.Credentials.SecretRef.Namespace = mg.GetNamespace()
		}
	case *namespacedv1beta1.ClusterProviderConfig:
		pcSpec = pc.Spec
	default:
		return nil, errors.New("unknown provider config type")
	}
	t := resource.NewProviderConfigUsageTracker(crClient, pcu)
	if err := t.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackUsage)
	}
	return &pcSpec, nil
}
