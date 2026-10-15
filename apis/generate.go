//go:build generate
// +build generate

// NOTE: See the below link for details on what is happening here.
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

// Remove existing CRDs
//go:generate rm -rf ../package/crds

// Remove generated files
//go:generate bash -c "find . -iname 'zz_*' ! -iname 'zz_generated.managed*.go' -delete"
//go:generate bash -c "find . -type d -empty -delete"
//go:generate bash -c "find ../internal/controller -iname 'zz_*' -delete"
//go:generate bash -c "find ../internal/controller -type d -empty -delete"
//go:generate rm -rf ../examples-generated

// Generate documentation from Terraform docs.
//go:generate go run github.com/crossplane/upjet/v2/cmd/scraper -n ${TERRAFORM_PROVIDER_SOURCE} -r ../.work/${TERRAFORM_PROVIDER_SOURCE}/${TERRAFORM_DOCS_PATH} -o ../config/provider-metadata.yaml --prelude-xpath "//text()[contains(., \"subcategory\")]" --resource-prefix mongodbatlas

// Run Upjet generator
//go:generate go run ../cmd/generator/main.go ..

// Generate deepcopy methodsets and CRD manifests
//go:generate go run -tags generate sigs.k8s.io/controller-tools/cmd/controller-gen object:headerFile=../hack/boilerplate.go.txt paths=./... crd:allowDangerousTypes=true,crdVersions=v1 output:artifacts:config=../package/crds

// Generate crossplane-runtime methodsets (resource.Claim, etc)
//go:generate go run -tags generate github.com/crossplane/crossplane-tools/cmd/angryjet generate-methodsets --header-file=../hack/boilerplate.go.txt ./...

// Transform the generated resolvers so that cross API-group reference sources
// are resolved dynamically through the runtime scheme instead of statically
// typed imports, which would otherwise form import cycles between API groups.
// The root group's package (mongodbatlas) needs an override because its group
// is the suffix itself, not <package>.<suffix>.
//go:generate go run github.com/crossplane/upjet/v2/cmd/resolver -g mongodbatlas.crossplane.io -a github.com/crossplane-contrib/provider-mongodbatlas/internal/apis -o mongodbatlas.mongodbatlas.crossplane.io=mongodbatlas.crossplane.io -s -p ./cluster/...
//go:generate go run github.com/crossplane/upjet/v2/cmd/resolver -g mongodbatlas.m.crossplane.io -a github.com/crossplane-contrib/provider-mongodbatlas/internal/apis -o mongodbatlas.mongodbatlas.m.crossplane.io=mongodbatlas.m.crossplane.io -s -p ./namespaced/...

package apis

import (
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen" //nolint:typecheck

	_ "github.com/crossplane/crossplane-tools/cmd/angryjet" //nolint:typecheck

	_ "github.com/crossplane/upjet/v2/cmd/resolver" //nolint:typecheck
)
