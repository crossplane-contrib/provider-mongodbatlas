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
	"net/http"
	"sync"

	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/diffserver"
	"github.com/crossplane/upjet/v2/pkg/terraform"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	tpf "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/mongodb/terraform-provider-mongodbatlas/xpshim"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// Placeholder programmatic API key the offline setup configures the Atlas
	// provider with. It is never sent anywhere, because offlineTransport
	// refuses every request before it leaves the process. It exists so that
	// the Atlas client is fully built, the way the Terraform provider's
	// resources expect to find it, without reading the ProviderConfig's
	// credentials, which are assumed to be unavailable.
	offlinePublicKey  = "offline-diff-server"
	offlinePrivateKey = "offline-diff-server-placeholder" //nolint:gosec // a placeholder, never sent

	errBuildOfflineClient = "cannot build the offline Atlas client"

	// fmtErrOfflineRequestBlocked reports an Atlas API request that the
	// offline diff server refused to send. The query string is left out of
	// the message.
	fmtErrOfflineRequestBlocked = "the offline diff server blocked an outbound MongoDB Atlas API request (%s %s): computing this diff requires calling the Atlas API, which is not possible offline"
)

// OfflineTerraformSetupBuilder returns a Terraform setup that configures the
// Atlas Terraform provider without credentials and without contacting Atlas.
// Use it in place of TerraformSetupBuilder in the diff server only.
//
// The ProviderConfig a resource references is still resolved, so that a
// resource whose ProviderConfig is missing fails the way it does in a
// reconcile, but nothing in it is read: an Atlas ProviderConfig carries only
// credentials, and the Secret they would come from is never read. The setup
// carries one Atlas client, built once for the process, as the Meta the SDKv2
// resources read and as the provider data offlineProvider hands the framework
// resources. That client cannot send a request: see offlineTransport.
func OfflineTerraformSetupBuilder(fw tpf.Provider) terraform.SetupFn {
	offline := sync.OnceValues(func() (*offlineProvider, error) {
		c, err := xpshim.NewOfflineClient(offlinePublicKey, offlinePrivateKey, "", offlineTransport{})
		if err != nil {
			return nil, errors.Wrap(err, errBuildOfflineClient)
		}
		return &offlineProvider{Provider: fw, client: c}, nil
	})
	return func(ctx context.Context, c client.Client, mg resource.Managed) (terraform.Setup, error) {
		if _, err := resolveProviderConfig(ctx, c, mg); err != nil {
			return terraform.Setup{}, errors.Wrap(err, errResolveProviderConfig)
		}
		p, err := offline()
		if err != nil {
			return terraform.Setup{}, err
		}
		return terraform.Setup{
			Configuration: map[string]any{
				keyPublicKey:  offlinePublicKey,
				keyPrivateKey: offlinePrivateKey,
			},
			Meta:              p.client,
			FrameworkProvider: p,
		}, nil
	}
}

// offlineProvider is the upstream Atlas Terraform provider with its Configure
// replaced. The upstream Configure resolves the credentials from the
// configuration block, the MONGODB_ATLAS_* environment and, when an AWS role
// is configured, AWS Secrets Manager; this one reads none of them and hands
// every resource the offline client. Everything else, the schema, the
// resources and the provider meta schema in particular, is the upstream
// provider's own.
type offlineProvider struct {
	tpf.Provider
	client any
}

// Configure hands the resources and the data sources the client that cannot
// reach Atlas. Ephemeral resources, which the diff server never opens, get no
// provider data.
func (p *offlineProvider) Configure(_ context.Context, _ tpf.ConfigureRequest, resp *tpf.ConfigureResponse) {
	resp.DataSourceData = p.client
	resp.ResourceData = p.client
}

// MetaSchema serves the upstream provider's meta schema, which its resources
// read on every operation.
func (p *offlineProvider) MetaSchema(ctx context.Context, req tpf.MetaSchemaRequest, resp *tpf.MetaSchemaResponse) {
	if mp, ok := p.Provider.(tpf.ProviderWithMetaSchema); ok {
		mp.MetaSchema(ctx, req, resp)
	}
}

// EphemeralResources serves the upstream provider's ephemeral resources, so
// that the wrapped provider has the same protocol surface as the upstream one.
func (p *offlineProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	if ep, ok := p.Provider.(tpf.ProviderWithEphemeralResources); ok {
		return ep.EphemeralResources(ctx)
	}
	return nil
}

// offlineTransport refuses every HTTP request it sees, failing the diff at
// once with an error the diff server reports as DIFF_COMPUTATION_NOT_SUPPORTED,
// instead of letting the request reach Atlas or fail slowly against an
// unreachable endpoint. One Atlas resource calls the API while planning: the
// resource policy validates its policies in ModifyPlan. Every other call lives
// in the CRUD functions the diff server never runs, so for those the guard is
// a safeguard against an upstream change.
type offlineTransport struct{}

func (offlineTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, diffserver.NewDiffComputationNotSupportedError(errors.Errorf(fmtErrOfflineRequestBlocked, req.Method, req.URL.Host+req.URL.Path))
}
