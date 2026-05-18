# Password Management

This document describes how database user passwords are handled in the provider,
including Bring Your Own Password (BYOP), auto-generation, and rotation.

The mechanism applies to **`mongodbatlas_database_user`** (both cluster-scoped and
namespaced-scoped variants).

## How it works

Password handling is implemented as a **managed resource initializer**.
It runs once, before the first reconciliation of a newly created resource.
After initialization, the standard upjet reconciliation loop reads the password from
the referenced Secret on every cycle and passes it to the Terraform `password`
field via the sensitive-field mapping (`"password" → "passwordSecretRef"`).

The initializer (`PasswordGenerator` in `config/*/common/common.go`) supports
two modes, selected automatically based on whether `passwordSecretRef` is
already set on the resource.

## Mode 1: Bring Your Own Password (BYOP)

The user provides a pre-existing Secret containing the desired password and
references it in `spec.forProvider.passwordSecretRef`.

```yaml
apiVersion: database.mongodbatlas.m.crossplane.io/v1alpha3
kind: User
metadata:
  name: my-db-user
spec:
  forProvider:
    projectId: 00001111aaaabbbb55556666
    username: my-user
    authDatabaseName: admin
    passwordSecretRef:
      name: my-password-secret    # must already exist
      key: password
    roles:
      - roleName: readWrite
        databaseName: mydb
  providerConfigRef:
    name: default
```

**What happens:**

1. `checkBYOP()` sees `passwordSecretRef.name` is non-empty and returns `true`.
2. Initializer exits immediately. No password generation.
3. On each reconciliation, upjet reads the Secret at `passwordSecretRef` and
   sends the value as the Terraform `password` parameter.

## Mode 2: Auto-generated password

The user omits `passwordSecretRef` and sets `writeConnectionSecretToRef`
instead. The provider generates a random password and stores it in that Secret.

```yaml
apiVersion: database.mongodbatlas.m.crossplane.io/v1alpha3
kind: User
metadata:
  name: my-db-user
spec:
  forProvider:
    projectId: 00001111aaaabbbb55556666
    username: my-user
    authDatabaseName: admin
    roles:
      - roleName: readWrite
        databaseName: mydb
  writeConnectionSecretToRef:
    name: my-db-user-conn         # password will be written here under key "password"
    namespace: default
  providerConfigRef:
    name: default
```

**What happens:**

1. `checkBYOP()` sees `passwordSecretRef` is unset and returns `false`.
2. `resolveConnRef()` reads `writeConnectionSecretToRef` and gets the Secret
   name and namespace.
3. `reconcilePassword()` checks if that Secret already contains a `password` key:
   - **Secret exists with `password` key**: skips generation, just wires up
     `passwordSecretRef`.
   - **Secret missing or no `password` key**: `generateAndApply()` generates a
     random password via `crossplane-runtime/pkg/password.Generate()`, creates
     or patches the Secret with `type: connection.crossplane.io/v1alpha1`, and
     sets an owner reference to the managed resource.
4. `setPasswordSecretRef()` sets `passwordSecretRef` on the resource to point
   to the same Secret (name + `key: "password"`), then persists via
   `client.Update()`.
5. From this point on, reconciliation reads the password from `passwordSecretRef`
   like BYOP mode.

### Idempotency

If the Secret already has a `password` key (e.g. from a previous initialization
that was interrupted after writing the Secret, but before updating the resource),
the initializer reuses the existing password instead of generating a new one.

## Secret lifecycle on User deletion

What happens to the Kubernetes Secret when a `User` resource is deleted depends
on how the password was provisioned.

### Auto-generated password

The Secret created by `generateAndApply()` carries a **controller owner
reference** pointing back to the `User` resource (`config/password/password.go:94`).
When the `User` is deleted, Kubernetes' built-in garbage collector detects the owner reference and
**automatically deletes the Secret** as part of the cascade.
No custom finalizer or cleanup code is needed.

In short: **User deleted implies Secret deleted automatically.**

### BYOP (Bring Your Own Password)

The provider never sets an owner reference on a user-supplied Secret and the initializer exits before
reaching `generateAndApply()`.
Deleting the `User` resource has **no effect on the Secret**; it remains in the cluster,
managed entirely by whoever created it.

In short: **User deleted does not delete the Secret.**

### Edge case: auto-generated Secret that already existed

If `writeConnectionSecretToRef` points to a Secret that already exists *and*
already contains a `password` key, the initializer reuses it without modifying
owner references.
In this case, the Secret has **no owner reference** to the `User`, so deleting the `User`
**does not delete the Secret**; same behavior as BYOP.

## Password rotation

The provider has **no built-in rotation scheduler**. Rotation is a manual or
external-automation concern.

### Rotating a BYOP password

1. Update the Kubernetes Secret referenced by `passwordSecretRef` with the new
   password.
2. On the next reconciliation cycle, upjet detects the changed `password` value
   and calls Terraform to update the Atlas database user.

### Rotating an auto-generated password

1. Update or delete the `password` key in the Secret referenced by `writeConnectionSecretToRef`
  (which is the same Secret that `passwordSecretRef` points to after initialization).
2. Same as above — the next reconciliation pushes the new password to Atlas.

### Triggering immediate reconciliation

Crossplane reconciles managed resources on a periodic poll interval. To avoid
waiting, any metadata change (e.g. adding a label) triggers the controller's
watch and starts a reconciliation immediately (real-time compositions):

```bash
kubectl label user.database.mongodbatlas.m.crossplane.io my-db-user \
  password-rotated="$(date +%s)" --overwrite
```
