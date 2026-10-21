/*
Copyright 2026 The Crossplane Authors.

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
	"encoding/base64"
	"fmt"
	"maps"
	"net"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	rtfake "github.com/crossplane/crossplane-runtime/v2/pkg/resource/fake"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	"github.com/crossplane/upjet/v2/pkg/diffserver"
	"github.com/crossplane/upjet/v2/pkg/terraform"
	"github.com/google/go-cmp/cmp"
	tpf "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/xpshim"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	clusterv1beta1 "github.com/crossplane-contrib/provider-mongodbatlas/apis/cluster/v1beta1"
	namespacedv1beta1 "github.com/crossplane-contrib/provider-mongodbatlas/apis/namespaced/v1beta1"
)

const (
	// atlasHost is where the Atlas SDKs send their requests by default, and
	// so where the offline client would send them if it could.
	atlasHost = "cloud.mongodb.com"
	// unguardedTestHost is the Atlas base URL host of the positive control.
	// No test can reach it: the tripwire armed by TestMain must see a request
	// for it when nothing stands in the way.
	unguardedTestHost = "cloud.mongodb.example"

	offlineTestNamespace = "team-a"
	offlineTestPCName    = "default"
	offlineTestProjectID = "5e0f6ad4f4e9f9b1c2d3e4f5"
	offlineTestOrgID     = "5e0f6ad4f4e9f9b1c2d3e4a1"
	offlineTestTeamID    = "5e0f6ad4f4e9f9b1c2d3e4b2"

	tfThirdPartyIntegration = "mongodbatlas_third_party_integration"
	tfTeam                  = "mongodbatlas_team"
	tfCluster               = "mongodbatlas_cluster"
)

// egress is the tripwire armed by TestMain, shared by every test in the
// package. Tests that read it must not run in parallel with one another.
var egress *proxyRecorder

// proxyRecorder is an HTTP proxy that records every request routed through it
// and forwards none of them. The Atlas provider's transport is built with
// http.ProxyFromEnvironment, so pointing HTTPS_PROXY at the recorder puts a
// tripwire on the process's real network boundary: nothing is stubbed or
// injected, and a request cannot leave for Atlas without being recorded here
// first.
type proxyRecorder struct {
	mu   sync.Mutex
	seen []string
}

func (p *proxyRecorder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// CONNECT carries its target in Host; an absolute-form proxy request
	// carries it in URL.
	target := r.Host
	if target == "" && r.URL != nil {
		target = r.URL.Host
	}
	p.mu.Lock()
	p.seen = append(p.seen, r.Method+" "+target)
	p.mu.Unlock()
	http.Error(w, "blocked by the egress tripwire", http.StatusBadGateway)
}

func (p *proxyRecorder) reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.seen = nil
}

// reached reports whether any recorded request targeted host.
func (p *proxyRecorder) reached(host string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, r := range p.seen {
		if strings.HasSuffix(r, " "+host) {
			return true
		}
	}
	return false
}

func TestMain(m *testing.M) {
	l, err := (&net.ListenConfig{}).Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot listen for the egress tripwire: %v\n", err)
		os.Exit(1)
	}
	egress = &proxyRecorder{}
	srv := &http.Server{Handler: egress, ReadHeaderTimeout: time.Second}
	go srv.Serve(l) //nolint:errcheck // the listener is closed when the process exits
	// net/http reads the proxy environment once, on the first request of the
	// process, so it must be set before any test sends one.
	for _, kv := range [][2]string{{"HTTPS_PROXY", "http://" + l.Addr().String()}, {"HTTP_PROXY", "http://" + l.Addr().String()}, {"NO_PROXY", ""}} {
		if err := os.Setenv(kv[0], kv[1]); err != nil {
			fmt.Fprintf(os.Stderr, "cannot set %s: %v\n", kv[0], err)
			os.Exit(1)
		}
	}
	os.Exit(m.Run())
}

// missingSecretSelectors reference a Secret that no test creates: a setup
// that resolves a ProviderConfig carrying them proves it never reads the
// credentials, as the online setup fails on the missing Secret.
func missingSecretSelectors() xpv2.CommonCredentialSelectors {
	return xpv2.CommonCredentialSelectors{
		SecretRef: &xpv2.SecretKeySelector{
			SecretReference: xpv2.SecretReference{Name: "missing-creds", Namespace: "crossplane-system"},
			Key:             "credentials",
		},
	}
}

func newTestScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	scheme := runtime.NewScheme()
	for _, add := range []func(*runtime.Scheme) error{corev1.AddToScheme, clusterv1beta1.SchemeBuilder.AddToScheme, namespacedv1beta1.SchemeBuilder.AddToScheme} {
		if err := add(scheme); err != nil {
			t.Fatalf("cannot build the test scheme: %v", err)
		}
	}
	return scheme
}

func modernManaged(kind, pcName string) xpresource.Managed {
	mg := &rtfake.ModernManaged{ObjectMeta: metav1.ObjectMeta{Name: "example", Namespace: offlineTestNamespace, UID: "test-uid"}}
	mg.SetProviderConfigReference(&xpv2.ProviderConfigReference{Kind: kind, Name: pcName})
	return mg
}

func legacyManaged(pcName string) xpresource.Managed {
	mg := &rtfake.LegacyManaged{ObjectMeta: metav1.ObjectMeta{Name: "example", UID: "test-uid"}}
	mg.SetProviderConfigReference(&xpv2.Reference{Name: pcName})
	return mg
}

// TestOfflineTerraformSetupBuilder asserts that the offline setup resolves the
// ProviderConfig a managed resource references, of every kind, never reads
// its credentials, and carries one guarded client for both halves of the
// provider.
func TestOfflineTerraformSetupBuilder(t *testing.T) {
	ctx := context.Background()
	kc := fake.NewClientBuilder().WithScheme(newTestScheme(t)).WithObjects(
		&namespacedv1beta1.ClusterProviderConfig{
			ObjectMeta: metav1.ObjectMeta{Name: offlineTestPCName},
			Spec: namespacedv1beta1.ProviderConfigSpec{
				Credentials: namespacedv1beta1.ProviderCredentials{Source: xpv2.CredentialsSourceSecret, CommonCredentialSelectors: missingSecretSelectors()},
			},
		},
		&namespacedv1beta1.ProviderConfig{
			ObjectMeta: metav1.ObjectMeta{Name: offlineTestPCName, Namespace: offlineTestNamespace},
			Spec: namespacedv1beta1.ProviderConfigSpec{
				Credentials: namespacedv1beta1.ProviderCredentials{Source: xpv2.CredentialsSourceSecret, CommonCredentialSelectors: missingSecretSelectors()},
			},
		},
		&clusterv1beta1.ProviderConfig{
			ObjectMeta: metav1.ObjectMeta{Name: offlineTestPCName},
			Spec: clusterv1beta1.ProviderConfigSpec{
				Credentials: clusterv1beta1.ProviderCredentials{Source: xpv2.CredentialsSourceSecret, CommonCredentialSelectors: missingSecretSelectors()},
			},
		},
	).Build()
	upstream := xpshim.GetFrameworkProvider()

	// The positive control for "never reads the credentials": the online
	// setup resolves the same ProviderConfig and fails on its missing Secret.
	if _, err := TerraformSetupBuilder(xpshim.GetSDKProvider(), upstream)(ctx, kc, modernManaged(namespacedv1beta1.ClusterProviderConfigKind, offlineTestPCName)); err == nil {
		t.Fatal("expected the online setup to fail on the missing credentials Secret")
	}

	cases := map[string]struct {
		mg      xpresource.Managed
		wantErr bool
	}{
		"ClusterProviderConfigFromNamespacedResource": {mg: modernManaged(namespacedv1beta1.ClusterProviderConfigKind, offlineTestPCName)},
		"ProviderConfigFromOwnNamespace":              {mg: modernManaged(namespacedv1beta1.ProviderConfigKind, offlineTestPCName)},
		"LegacyProviderConfig":                        {mg: legacyManaged(offlineTestPCName)},
		"MissingProviderConfig":                       {mg: modernManaged(namespacedv1beta1.ProviderConfigKind, "missing"), wantErr: true},
	}
	setup := OfflineTerraformSetupBuilder(upstream)
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ps, err := setup(ctx, kc, tc.mg)
			if (err != nil) != tc.wantErr {
				t.Fatalf("OfflineTerraformSetupBuilder(...)(...) error = %v, wantErr %t", err, tc.wantErr)
			}
			if err != nil {
				return
			}
			want := terraform.ProviderConfiguration{keyPublicKey: offlinePublicKey, keyPrivateKey: offlinePrivateKey}
			if diff := cmp.Diff(want, ps.Configuration); diff != "" {
				t.Errorf("unexpected Terraform provider configuration, -want +got:\n%s", diff)
			}
			p, ok := ps.FrameworkProvider.(*offlineProvider)
			if !ok || p.Provider != upstream {
				t.Fatalf("expected the upstream provider wrapped by offlineProvider, got %T", ps.FrameworkProvider)
			}
			if ps.Meta == nil || ps.Meta != p.client {
				t.Errorf("expected the SDKv2 meta and the framework provider data to be the same client, got %T and %T", ps.Meta, p.client)
			}
		})
	}
}

// offlineMeta returns the client the offline setup hands the Atlas resources.
func offlineMeta(t *testing.T) any {
	t.Helper()
	kc := fake.NewClientBuilder().WithScheme(newTestScheme(t)).WithObjects(
		&namespacedv1beta1.ClusterProviderConfig{ObjectMeta: metav1.ObjectMeta{Name: offlineTestPCName}},
	).Build()
	ps, err := OfflineTerraformSetupBuilder(xpshim.GetFrameworkProvider())(context.Background(), kc, modernManaged(namespacedv1beta1.ClusterProviderConfigKind, offlineTestPCName))
	if err != nil {
		t.Fatalf("offline setup failed: %v", err)
	}
	return ps.Meta
}

// encodeStateID encodes a composite Terraform ID the way the Atlas provider
// does: base64 key and value pairs, in key order, joined with dashes.
func encodeStateID(values map[string]string) string {
	parts := make([]string, 0, len(values))
	for _, k := range slices.Sorted(maps.Keys(values)) {
		parts = append(parts, base64.StdEncoding.EncodeToString([]byte(k))+":"+base64.StdEncoding.EncodeToString([]byte(values[k])))
	}
	return strings.Join(parts, "-")
}

// readFunc returns the read function of an SDKv2 resource. The Atlas provider
// registers its resources without timeouts.
func readFunc(t *testing.T, r *schema.Resource) schema.ReadContextFunc {
	t.Helper()
	switch {
	case r.ReadWithoutTimeout != nil:
		return r.ReadWithoutTimeout
	case r.ReadContext != nil:
		return r.ReadContext
	}
	t.Fatal("the resource has no context-aware read function")
	return nil
}

func diagnosticsText(diags diag.Diagnostics) string {
	parts := make([]string, 0, len(diags))
	for _, d := range diags {
		parts = append(parts, d.Summary+": "+d.Detail)
	}
	return strings.Join(parts, "\n")
}

// TestOfflineClientBlocksEgress asserts that a call the Atlas provider's own
// resources make on the client the offline setup hands them fails as a diff
// the server cannot compute, and that nothing reaches the process's network
// boundary. The resources are chosen for the SDK clients their reads go
// through first: the current, a pinned older one, and the legacy Atlas client.
func TestOfflineClientBlocksEgress(t *testing.T) {
	meta := offlineMeta(t)
	sdk := xpshim.GetSDKProvider()

	cases := map[string]struct {
		resource string
		config   map[string]any
		id       string
	}{
		"ThirdPartyIntegrationReadsThroughAtlasV2": {
			resource: tfThirdPartyIntegration,
			config:   map[string]any{"project_id": offlineTestProjectID, "type": "DATADOG"},
			id:       encodeStateID(map[string]string{"project_id": offlineTestProjectID, "type": "DATADOG"}),
		},
		"TeamReadsThroughAtlasV220241113": {
			resource: tfTeam,
			config:   map[string]any{"org_id": offlineTestOrgID, "name": "example"},
			id:       encodeStateID(map[string]string{"org_id": offlineTestOrgID, "id": offlineTestTeamID}),
		},
		"ClusterReadsThroughTheLegacyAtlasClient": {
			resource: tfCluster,
			config:   map[string]any{"project_id": offlineTestProjectID, "name": "example", "provider_name": "AWS"},
			id:       encodeStateID(map[string]string{"project_id": offlineTestProjectID, "cluster_name": "example", "provider_name": "AWS"}),
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			r, ok := sdk.ResourcesMap[tc.resource]
			if !ok {
				t.Fatalf("the SDKv2 provider has no resource %q", tc.resource)
			}
			d := schema.TestResourceDataRaw(t, r.Schema, tc.config)
			d.SetId(tc.id)
			egress.reset()

			diags := readFunc(t, r)(context.Background(), d, meta)
			if !diags.HasError() {
				t.Fatal("expected the read to fail as unsupported offline")
			}
			msg := diagnosticsText(diags)
			if !strings.Contains(msg, diffserver.ErrDiffComputationNotSupported.Error()) || !strings.Contains(msg, "blocked an outbound MongoDB Atlas API request") {
				t.Errorf("expected the read to fail with the offline diff server's error, got: %s", msg)
			}
			if egress.reached(atlasHost + ":443") {
				t.Fatalf("a request for %s reached the network boundary", atlasHost)
			}
		})
	}
}

// TestEgressTripwireSeesAnUnguardedClient is the positive control of the
// tripwire: the same read through the client the configured SDKv2 provider
// builds, with a key to send, is recorded at the network boundary.
func TestEgressTripwireSeesAnUnguardedClient(t *testing.T) {
	ctx := context.Background()
	sdk := xpshim.GetSDKProvider()
	cfg := sdkterraform.NewResourceConfigRaw(map[string]any{
		keyPublicKey:  "unguarded",
		keyPrivateKey: "unguarded-secret",
		"base_url":    "https://" + unguardedTestHost,
	})
	if diags := sdk.Configure(ctx, cfg); diags.HasError() {
		t.Fatalf("cannot configure the SDKv2 provider: %s", diagnosticsText(diags))
	}
	r := sdk.ResourcesMap[tfThirdPartyIntegration]
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]any{"project_id": offlineTestProjectID, "type": "DATADOG"})
	d.SetId(encodeStateID(map[string]string{"project_id": offlineTestProjectID, "type": "DATADOG"}))
	egress.reset()

	if diags := readFunc(t, r)(ctx, d, sdk.Meta()); !diags.HasError() {
		t.Fatal("expected the tripwire to fail the read")
	}
	if !egress.reached(unguardedTestHost + ":443") {
		t.Fatalf("expected a request for %s at the network boundary", unguardedTestHost)
	}
}

// TestOfflineProviderWrapsUpstream asserts that the offline provider hands
// out its client from Configure and otherwise serves what the upstream
// provider serves.
func TestOfflineProviderWrapsUpstream(t *testing.T) {
	ctx := context.Background()
	upstream := xpshim.GetFrameworkProvider()
	client := &struct{ name string }{name: "offline"}
	p := &offlineProvider{Provider: upstream, client: client}

	resp := tpf.ConfigureResponse{}
	p.Configure(ctx, tpf.ConfigureRequest{}, &resp)
	if resp.ResourceData != any(client) || resp.DataSourceData != any(client) {
		t.Errorf("expected Configure to hand the client to resources and data sources, got %T and %T", resp.ResourceData, resp.DataSourceData)
	}

	// The upstream resources read the provider meta on every operation, so
	// the wrapper must keep serving the upstream meta schema.
	mp, ok := upstream.(tpf.ProviderWithMetaSchema)
	if !ok {
		t.Fatal("the upstream provider no longer serves a provider meta schema; the forwarder in offlineProvider can go")
	}
	want, got := tpf.MetaSchemaResponse{}, tpf.MetaSchemaResponse{}
	mp.MetaSchema(ctx, tpf.MetaSchemaRequest{}, &want)
	p.MetaSchema(ctx, tpf.MetaSchemaRequest{}, &got)
	if len(want.Schema.Attributes) == 0 {
		t.Fatal("expected the upstream provider meta schema to declare attributes")
	}
	if diff := cmp.Diff(slices.Sorted(maps.Keys(want.Schema.Attributes)), slices.Sorted(maps.Keys(got.Schema.Attributes))); diff != "" {
		t.Errorf("unexpected provider meta schema attributes, -want +got:\n%s", diff)
	}

	ep, ok := upstream.(tpf.ProviderWithEphemeralResources)
	if !ok {
		t.Fatal("the upstream provider no longer serves ephemeral resources; the forwarder in offlineProvider can go")
	}
	if want, got := len(ep.EphemeralResources(ctx)), len(p.EphemeralResources(ctx)); want != got {
		t.Errorf("expected %d ephemeral resources, got %d", want, got)
	}
}
