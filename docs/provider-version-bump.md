# Bumping the Atlas Terraform Provider Version

This provider compiles `terraform-provider-mongodbatlas` in-process via
the `xpshim` package. Bumping the upstream version requires the steps
below.

## Procedure

1. **Update version and SHA** in the Makefile:
   ```makefile
   TERRAFORM_PROVIDER_VERSION ?= <new-version>
   PROVIDER_SOURCE_SHA := <new-commit-sha>
   ```
   Get the SHA from the release tag:
   ```bash
   git ls-remote https://github.com/mongodb/terraform-provider-mongodbatlas refs/tags/v<new-version>
   ```

2. **Re-extract the provider source:**
   ```bash
   rm -rf third_party/terraform-provider-mongodbatlas
   make provider-source
   ```

3. **Verify the shim compiles** (catches signature drift in `NewSdkV2Provider`/`NewFrameworkProvider`):
   ```bash
   go build ./third_party/terraform-provider-mongodbatlas/xpshim/
   ```

4. **Resolve dependency conflicts:**
   ```bash
   go mod tidy
   go build ./...
   ```
   If plugin-sdk, plugin-framework, or plugin-go versions conflict,
   lift to the version required by the Atlas provider.

5. **Run classification test** (catches resources moving between SDKv2 and framework):
   ```bash
   go test -v -run TestClassificationMatchesUpstream ./config/...
   ```
   If it fails, update `terraformSDKIncludedResources` and
   `terraformFrameworkIncludedResources` in `config/config.go`.

6. **Regenerate:**
   ```bash
   make generate
   ```

7. **Build and run all tests:**
   ```bash
   go build ./...
   go test ./config/... ./internal/clients/...
   ```

8. **Schema diff** — review CRD changes for breaking field removals:
   ```bash
   git diff package/crds/
   ```

9. **E2E validation:**
   ```bash
   make e2e
   ```
