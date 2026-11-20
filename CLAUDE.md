# Crossplane Provider MongoDB Atlas

Upjet-based Crossplane provider for MongoDB Atlas. Generates CRDs from the `mongodb/terraform-provider-mongodbatlas`.

## Project Structure

- `config/` — Provider configuration (external names, resource overrides, per-resource configs)
  - `config/config.go` — Resource inclusion list (`includedResources`), default `identifierFromProvider()` external name, `ExternalNameConfigurations()` (sets default for all resources), `ExternalNameConfigured()` (returns inclusion patterns)
  - `config/provider.go` — Provider setup, default resource options, resource configurator registration. Both cluster-scoped and namespace-scoped providers share a single `newProvider()` factory differing only in root group and password generator. `GetProvider` and `GetProviderNamespaced` take the Terraform provider instances from the caller, and `cmd/` builds one pair for both scopes. Configurators edit SDKv2 schemas in place, once per scope, so each edit must be idempotent (`TestScopesShareTerraformProviders` checks it).
  - `config/resources/` — Per-resource-group configurators (one file per group: `mongodbatlas.go`, `cloud.go`, `database.go`, `federated.go`, `ldap.go`, `network.go`, `privateendpoint.go`, `project.go`, `search.go`, `stream.go`, `alert.go`)
  - `config/resources/external_name.go` — All external name constructors (`templated`, `importJoinedID`, `importJoinedIDOrdered`, `importJoinedIDMapped`, `importJoinedIDHidden`, `accessListImportJoinedID`) and Atlas-specific base64 state ID encoding/decoding helpers
  - `config/refs/` — Shared constants (`refs.go`: field names, TF resource names, version constants) and reusable `GetIDFn`/`GetExternalNameFn` helpers (`helpers.go`)
  - `config/password/` — Password generation/BYOP logic with cluster-scoped (`ClusterGenerator`) and namespace-scoped (`NamespacedGenerator`) variants
- `apis/cluster/` and `apis/namespaced/` — API types. Only `zz_*` files are auto-generated (do not edit). Non-`zz_` files (e.g. `register.go`, `types.go`) are hand-maintained.
- `internal/controller/` — Generated controllers. Only `zz_*` files are auto-generated (do not edit).
- `package/crds/` — Generated CRDs (do not edit manually)
- `examples-generated/` — Generated example manifests (do not edit manually). `make generate` copies only new files into `examples/`, which is hand-maintained
- `cluster/images/provider-mongodbatlas/Dockerfile` — Provider container image

## Key Commands

- `make provider-source` — Extract the Atlas Terraform provider source into `third_party/` (gitignored) and inject the `xpshim` package plus one file into its `internal/config`; the extraction happens once, the injected files are refreshed from `hack/xpshim.go.tmpl` and `hack/xpshim_config.go.tmpl` on every run. **Required on fresh checkouts before anything compiles** (go.mod has a `replace` directive pointing at `third_party/`). Runs automatically as part of `make build` and `make generate`. Clones the `v$(TERRAFORM_PROVIDER_VERSION)` tag and records it in `third_party/terraform-provider-mongodbatlas/.provider-version`; it clones again only when that stamp differs from `TERRAFORM_PROVIDER_VERSION`.
- `make generate` — Regenerate APIs, controllers, CRDs, and examples from config
- `make build` — Build provider binaries
- `make local-deploy` — Build and deploy provider locally (requires a running cluster)
- `make e2e` — Run end-to-end tests
- `go build ./config/...` — Quick compile check for config changes

## Bumping the Atlas Terraform provider

1. Set `TERRAFORM_PROVIDER_VERSION` in the `Makefile`. Read the upstream `CHANGELOG.md` for new resources, newly sensitive attributes and resources moving between SDKv2 and the plugin framework.
2. Run `make provider-source`, `go mod tidy` and `go build ./...`. If plugin-sdk, plugin-framework or plugin-go conflict, lift them to the versions the Atlas provider requires.
3. Run `go test -run TestClassificationMatchesUpstream ./config/...`. It fails when a resource moved between SDKv2 and the plugin framework; fix `terraformSDKIncludedResources`/`terraformFrameworkIncludedResources` in `config/config.go`. It does not report new upstream resources, so add those by hand (include list, configurator, examples).
4. Run `make generate` and review `git diff package/crds/`. An attribute that upstream newly marks `sensitive` leaves `status.atProvider` and moves to the connection secret (2.22.0 did this to `Organization` `serviceAccount[].secrets`). That breaks anyone reading it from status, so put it in the release notes.
5. Run `go build ./...`, `go test ./config/... ./internal/clients/... ./cmd/...` and `make e2e`.

The caches under `third_party/` and `.work/` follow `TERRAFORM_PROVIDER_VERSION` on their own: `provider-source` checks its stamp, `pull-docs` moves the `.work/mongodb/mongodbatlas` docs clone to the new tag, and the schema step runs `terraform init -upgrade` so the old `.work/terraform/.terraform.lock.hcl` cannot block it. A docs clone left on the old tag produces old field descriptions and no examples for new resources.

## Code Generation (upjet)

All resources are code-generated by upjet. The generation pipeline:

1. **Resource inclusion** (`config/config.go`) defines which TF resources to include and sets a default `identifierFromProvider()` external name for all
2. **Resource configurators** (`config/resources/*.go`) set per-resource overrides: `ShortGroup`, `Kind`, `Version`, `References`, `LateInitializer`, `r.ExternalName`, custom `GetIDFn`/`GetExternalNameFn`. Each resource sets its own external name via constructors from `config/resources/external_name.go`.
3. **Default resource options** (`config/provider.go`) apply to all resources: ShortGroup reset and external name configuration
4. `make generate` runs upjet to produce `zz_*` files in `apis/` and `internal/controller/`, plus `package/crds/` and `examples-generated/`

### Important upjet behaviors

- `TemplatedStringAsIdentifier("", tmpl)` parses `{{ .parameters.X }}` from the template and adds those parameters to `IdentifierFields`. Parameters in `IdentifierFields` are **removed from the generated CRD schema** (no `forProvider` field, no Ref/Selector). Our `templated()` wrapper clears `IdentifierFields` to prevent this.
- `DisableNameInitializer` is checked at code generation time and baked into `zz_controller.go`. You must run `make generate` after changing it.
- When changing a resource's `Version` (e.g., v1alpha2 → v1alpha3), delete the old `apis/*/GROUP/OLD_VERSION/` directory before running `make generate`, or stale files will cause build errors.
- Two provider scopes exist (cluster and namespaced) that generate into separate API packages. Both share a single set of resource configurators in `config/resources/`; only the root group and password generator differ (see `newProvider()` in `config/provider.go`).
- A sensitive object attribute (e.g. the list `service_account[*].secrets[*]`) goes to the connection secret one field at a time, and `GetSensitiveAttributes` fails the reconcile when any field is not a string, null included. SDKv2 resources store empty strings and are safe; check what a plugin-framework resource's Read stores before relying on it.
- upjet's `getDescription` drops everything up to the first `-` of an argument doc, expecting `` `name` - description ``. The Atlas docs use the tfplugindocs layout without that separator, so `cleanArgumentDocs()` in `config/provider.go` prefixes `- ` to docs whose first hyphen is inside the text (`24-hexadecimal`, `cloud-agnostic`). Without it those CRD descriptions lose their start. It also drops the tfplugindocs type marker (`(String)`, `(Attributes List, Sensitive)`). tfplugindocs writes those docs from the schema description, which upjet prints after the doc, so the doc text is emptied when every field it applies to has a schema description; a `Deprecated` marker becomes `Deprecated.`. Two kinds of field print no schema description even when `schema.json` has one, so their doc must stay: nested objects converted from the plugin framework, and the top-level `id`, which `pkg/pipeline/crd.go` replaces after the resource options run.
- An upjet bump that only changes runtime code (`pkg/controller`, `pkg/resource`, `pkg/terraform`) needs no `make generate`; changes under `pkg/config`, `pkg/types` or `pkg/pipeline` do. Check with `git diff --stat OLD NEW -- pkg` in an upjet clone.

## External Name Patterns

Each resource sets `r.ExternalName` in its configurator using constructors from `config/resources/external_name.go`. The default (set by `config/config.go`) is `identifierFromProvider()`. Resources that don't override it use provider-assigned IDs.

- `identifierFromProvider()` — Wraps `config.IdentifierFromProvider` with `DisableNameInitializer = false`. Provider-assigned ID; all parameters stay in the CRD schema.
- `providerAssignedID()` — `config.IdentifierFromProvider` (name initializer off) with `refs.NotFoundWhenPathParamEmpty`, for plugin-framework resources whose Terraform ID is the plain provider-assigned ID and whose Read fails client-side on an empty ID. Use it instead of the default when other resources reference the target by external name (a `metadata.name` placeholder gets cached by early references, as with Project) or when Atlas answers a malformed ID with 400 instead of 404 (ResourcePolicy).
- `templated(tmpl)` — Template-based ID (e.g. `"{{ .parameters.project_id }}/{{ .parameters.name }}"`). Clears `IdentifierFields` so template parameters remain in the CRD schema with Ref/Selector support.
- `importJoinedID(fields, externalNameKey)` — For resources whose TF state ID is `EncodeStateID` over `fields`. Uses base64-encoded state IDs at runtime (`encodedStateGetIDFn`/`encodedStateGetExternalNameFn`). No import ID is configured: the no-fork connectors never run `terraform import`.
- `importJoinedIDOrdered(importOrder, externalNameKey)` — Like `importJoinedID` but for resources where the provider-assigned key appears at a non-trailing position.
- `importJoinedIDMapped(paramOrder, fieldMapping)` — Like `importJoinedID` but with param→stateKey name mapping (e.g. `name` → `cluster_name`).
- `importJoinedIDAssigned(fields, externalNameKey)` — Like `importJoinedID` but marks `externalNameKey` as hidden (not in CRD schema). Used when the key is provider-assigned and shouldn't be user-settable.
- `accessListImportJoinedID(prefixParams)` — For access-list resources with ip_address/cidr_block/aws_security_group fallback logic.
- `computedKeyID(key)` — For plugin-framework resources with no `id` attribute whose Read needs a provider-assigned computed attribute (`key`). Reads the external name from `key` in state and copies the annotation into `key` for imports. Pair it with `r.TerraformPluginFrameworkIsStateEmptyFn = refs.StateEmptyWhenAttributeUnset(key)` so that Create runs. Only for a computed key: when the attribute the Read needs is required (user-supplied, like `endpoint_id` on `mongodbatlas_stream_workspace_private_endpoint`), use `templated` over the upstream `ImportState` `idAttributes` joined with `/`.

### How `templated()` differs from upstream

`config.TemplatedStringAsIdentifier("", tmpl)` would normally add template parameters to `IdentifierFields`, removing them from the CRD schema. Our wrapper:
1. Passes an empty `nameField` to avoid tying the external name to `metadata.name`.
2. Sets `DisableNameInitializer = false` so the name initializer runs (required for Observe-only policies).
3. Clears `IdentifierFields` so parameters like `project_id` stay in `forProvider` with Ref/Selector support.
4. Wraps `GetIDFn` to accept the external-name annotation as canonical ID when present.
5. Wraps `GetExternalNameFn` with a fallback to re-evaluate the template from state fields when `id` is absent.

## Atlas clients

`TerraformSetupBuilder` (`internal/clients/mongodbatlas.go`) builds one Atlas client per credential set with the SDKv2 provider's `Configure` and keeps it for `--client-cache-ttl` (`CLIENT_CACHE_TTL`, default 30m). The framework resources get the same client through `configuredProvider`, which replaces the upstream `Configure`, so the framework `Configure` that upjet calls on every reconcile builds nothing.

- Upstream caches service account token sources by client ID only and rejects any other secret for that ID until the process restarts. On a cache miss, `xpshim.ForgetStaleServiceAccount` (injected into the upstream `internal/config` from `hack/xpshim_config.go.tmpl`) drops the entry when its secret differs, so rotated secrets and two ProviderConfigs sharing a service account both work.
- For credentials its own validation rejects (service account credentials in `public_key`, for instance), the upstream SDKv2 `Configure` returns no client and no error. A nil meta fails the setup and is not cached; `configureCredentials` rejects an `mdb_sa` `public_key` with its own message first.

## Offline diff server

`provider internal diff-server [--network tcp|unix] [--address :9099]` serves upjet's diff gRPC service (`proto/diff/v1alpha1`, `pkg/diffserver`), which the Crossplane CLI drives to simulate changes. `start` is the default command, so the image keeps its behaviour when invoked with flags alone. The server never contacts the MongoDB Atlas API or a cluster:

- It is dispatched before `ctrl.GetConfig()` in `cmd/provider/main.go`, so no kubeconfig is needed; ProviderConfigs come from each request's object store.
- `internal/clients/offline.go` replaces the Terraform setup: the ProviderConfig is resolved (a missing one fails like in a reconcile) but its credentials are never read. One Atlas client, built by `xpshim.NewOfflineClient` in `hack/xpshim.go.tmpl` with a placeholder API key, serves as the SDKv2 `Meta` and as the provider data of the framework resources, and every request its Atlas SDK clients would send is refused by a transport that returns `DIFF_COMPUTATION_NOT_SUPPORTED` (the Realm client, which no resource uses while planning, is left without keys and fails with its own error before sending). The upstream framework `Configure`, which reads the `MONGODB_ATLAS_*` environment and can call AWS Secrets Manager, never runs: `configuredProvider` (shared with the online setup) wraps the upstream provider with a `Configure` that hands out the offline client.
- `xpshim.NewOfflineClient` lives in the injected template because the Atlas client type is in the Terraform provider's `internal/config`. It builds the client the way the configured provider does, checks that the admin SDK clients share the one HTTP client, and swaps that client's transport (the legacy `Atlas` client hides its HTTP client, so only the tripwire test guards it). Edit `hack/xpshim.go.tmpl`, then run `make provider-source` to refresh the injected copy (`make build` and `make generate` do that too; a bare `go test` does not).
- A resource whose plan needs the API (the ResourcePolicy validates its policies in `ModifyPlan`) fails with `FailedPrecondition` instead of reaching out. Every other resource plans offline.
- `cmd/provider/diffserver_test.go` drives the real server over a unix socket; `internal/clients/offline_test.go` proves no request reaches the process's network boundary with an `HTTPS_PROXY` tripwire and a positive control, through the resources' own read functions.
- `resolveModern` clears the namespace for a `ClusterProviderConfig`: the server's in-memory store keys objects on the exact namespace, unlike the API server.
- `package/crossplane.yaml` declares the `DiffServer` capability. The CLI only runs an image whose package metadata carries it (`crossplane resource simulate` reports "package does not declare the DiffServer capability" and falls back to a client-side diff otherwise) and starts the container with `--debug internal diff-server --address :9099`.
- To try it: `crossplane resource simulate -f resources.yaml --providers providers.yaml`, with a Provider manifest whose `spec.package` is the built xpkg pushed to a registry the Docker daemon can pull from. The CLI still needs a kubeconfig: it reads the live resources and the referenced ProviderConfigs from the cluster (the ProviderConfig CRDs and objects must exist there), and only the provider discovery is replaced by `--providers`.

## Conventions

- Every resource must have an `AddResourceConfigurator` in its `config/resources/*.go` with explicit `r.ExternalName` (unless default `identifierFromProvider()` is correct) and `r.References` for `project_id` (or other cross-resource refs)
- `resetRootShortGroup()` only clears ShortGroup for 2-word TF resources (where the default equals `"mongodbatlas"`). Resources with 3+ words in their TF name (e.g. `mongodbatlas_advanced_cluster`) **must** set `r.ShortGroup = ""` explicitly in their configurator to land in the root API group.
- Resources in sub-groups set `r.ShortGroup = "groupname"` explicitly in their configurator
- Version constants (`VersionV1Alpha2`, `VersionV1Alpha3`) live in `config/refs/refs.go`
- Use `refs.ExtractParamPath(field, isObservation)` for cross-resource extractor references (`isObservation=true` reads `status.atProvider`, `false` reads `spec.forProvider`) instead of `fmt.Sprintf(refs.ExtractParamPathFmt, ...)`
- Use `refs.AccessListGetIDFn(prefixParams...)` for access-list resources with ip_address/cidr_block fallback
- Examples in `examples/` must stay uptestable: set `projectId`/`orgId` (and nested ones like `targetProjectId`, `groupId`) to `${data.mongodbatlas.project_id}`/`${data.mongodbatlas.org_id}` from the uptest datasource, never a hardcoded ID or a selector without a matching object of the reference's target kind in the same file. Other datasource keys: `mongodbatlas.user_id` (active org member), `mongodbatlas.aws_container_region` (AWS region with no network container in the project), the Datadog keys `datadog.api_key`/`datadog.region`/`datadog.otlp_metrics_endpoint` (the examples create their own Secret from them) and the optional `mongodbatlas.federation_settings_id` (federated examples keep the placeholder behind manual-intervention, so a federated org can drop the annotation and run them locally). No Terraform expressions (`${var.*}`, `${aws_*}`): a field that held one gets a selector, adding the missing `r.References` entry in `config/` when needed. New files that `make generate` copies into `examples/` come from the Terraform docs and break these rules (`${var.*}` values, a `projectIdSelector` with nothing to match, unprefixed names), so rewrite them before committing
- Example clusters use the cheapest tier the resource supports (Flex, else M10) in AWS `US_EAST_1`, and `${Rand.RFC1123Subdomain}` names (uptest gives each occurrence its own value, so only use it where the name is read through a selector). Manifests that need infrastructure outside Atlas (AWS, Azure, Confluent, Kafka, Datadog, LDAP, App Services) or that cannot be created or cleaned up carry `upjet.upbound.io/manual-intervention` with the reason, on every manifest of the file
- Examples must be safe to run in one uptest batch: prefix every `metadata.name` and `testing.upbound.io/example-name` label (and the selectors that match it) with the file name, e.g. `onlinearchive-example`. Files with clusters set `uptest.upbound.io/timeout` (3600, or 5400 with backup snapshots). Clusters that need the Atlas sample datasets carry the label `testing.upbound.io/load-sample-data: "true"` and the dependent resource sets `uptest.upbound.io/pre-assert-hook: ../../../testhooks/load-sample-data.sh` (needs `atlas` and `jq`; hooks must be executable)
- Keep `examples/cluster` and `examples/namespaced` in sync: the namespaced file is the cluster one with `.m.crossplane.io` API groups, `metadata.namespace: crossplane-system`, and no `namespace` in secret references
