# MongoDBAtlas Provider

`provider-mongodbatlas` is a [Crossplane](https://crossplane.io/) provider that
is built using [Upjet](https://github.com/crossplane/upjet) code
generation tools and exposes XRM-conformant managed resources for the
MongoDBAtlas API.

## Getting Started

Install the provider by using the following command after changing the image tag
to the [latest release](https://github.com/crossplane-contrib/provider-mongodbatlas/releases):

```
kubectl apply -f - <<EOF
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-mongodbatlas
spec:
  package: xpkg.upbound.io/crossplane-contrib/provider-mongodbatlas:v1.0.0
EOF
```

You can see the API reference [here](https://doc.crds.dev/github.com/crossplane-contrib/provider-mongodbatlas).

## Importing a resource

Refer to the [dedicated document](docs/import.md) for information about how to import existing resources.

## Developing

Run code-generation pipeline:

```console
go run cmd/generator/main.go
```

Run against a Kubernetes cluster (out of cluster):

```console
make run
```

or (deploying in-cluster):

```console
make local-deploy
```

Review your code:

```console
make reviewable
```

Build, push, and install:

```console
make all
```

Build image:

```console
make image
```

Push image:

```console
make push
```

Build binary:

```console
make build
```

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/crossplane-contrib/provider-mongodbatlas/issues).

## Contact

Please use the following to reach members of the community:

* Slack: Join our [slack channel](https://slack.crossplane.io)
* Forums:
  [crossplane-dev](https://groups.google.com/forum/#!forum/crossplane-dev)
* Twitter: [@crossplane_io](https://twitter.com/crossplane_io)
* Email: [info@crossplane.io](mailto:info@crossplane.io)

## Governance and Owners

`provider-mongodbatlas` is run according to the same
[Governance](https://github.com/crossplane/crossplane/blob/master/GOVERNANCE.md)
and [Ownership](OWNERS.md)
structure as the core Crossplane project.

## Code of Conduct

`provider-mongodbatlas` adheres to the same [Code of
Conduct](https://github.com/crossplane/crossplane/blob/master/CODE_OF_CONDUCT.md)
as the core Crossplane project.

## Licensing

`provider-mongodbatlas` is under the Apache 2.0 [license](LICENSE).
