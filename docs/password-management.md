# Password Management

This document describes how database user passwords are handled in the provider,
including Bring Your Own Password (BYOP), auto-generation, and rotation.

The mechanism applies to **`mongodbatlas_database_user`** (both cluster-scoped and
namespaced-scoped variants).

## How it works

Password handling is implemented as a **managed resource initializer**.
It runs once, before the first reconciliation of a newly created resource.
After initialization, the standard Upjet reconciliation loop reads the password
from an internal field (`passwordSecretRef`) on every cycle and passes it to the
Terraform `password` field. Users do not need to set `passwordSecretRef` directly:
the initializer wires it automatically.

The initializer (`PasswordGenerator` in `config/password/password.go`) supports two modes.
Both use **the same user-facing field**: `writeConnectionSecretToRef`.
The mode is selected automatically based on whether that Secret already contains a `password` key
when the initializer runs.

## Mode 1: Bring Your Own Password (BYOP)

Create a Secret containing a `password` key **before** creating the `User` resource.
Set `writeConnectionSecretToRef` to reference that Secret.

> **Important:** The Secret **must** have `type: connection.crossplane.io/v1alpha1`.
> Crossplane refuses to write connection details to Secrets of other types (e.g. `Opaque`),
> and you will get the error: `cannot create or update connection secret: refusing to modify uncontrolled secret of type "opaque"`.
> Kubernetes does not allow changing a Secret's type after creation — delete and recreate if needed.

### 1. Create the Secret first

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-db-user-conn
  namespace: default
type: connection.crossplane.io/v1alpha1
stringData:
  password: mySecurePassword
```

### 2. Create the User pointing to the same Secret

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
    name: my-db-user-conn
    namespace: default
  providerConfigRef:
    name: default
```

**What happens:**

1. `checkBYOP()` sees no `passwordSecretRef` on the resource and returns `false`.
2. `resolveConnRef()` reads `writeConnectionSecretToRef` and gets the Secret name and namespace.
3. `reconcilePassword()` fetches the Secret, finds an existing `password` key, and skips generation.
4. The initializer wires `passwordSecretRef` internally to point at the same Secret (`key: "password"`).
5. On each reconciliation, Upjet reads the password from `passwordSecretRef` and sends it to Atlas.

Because the Secret was not created by `generateAndApply()`, no owner reference is set.
Deleting the `User` does **not** delete the Secret.

## Mode 2: Auto-generated password

Set `writeConnectionSecretToRef` without pre-creating the Secret (or create one without a `password`
key). The provider generates a random password and stores it in that Secret.

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

1. `checkBYOP()` sees no `passwordSecretRef` on the resource and returns `false`.
2. `resolveConnRef()` reads `writeConnectionSecretToRef` and gets the Secret name and namespace.
3. `reconcilePassword()` finds the Secret missing or without a `password` key.
4. `generateAndApply()` generates a random password, creates or patches the Secret
   with `type: connection.crossplane.io/v1alpha1`, and sets a controller owner
   reference to the `User` resource (so the Secret is garbage-collected on deletion).
5. The initializer wires `passwordSecretRef` internally to point at the same Secret (`key: "password"`).
6. From this point on, reconciliation reads the password from `passwordSecretRef` like BYOP mode.

### Idempotency

If the Secret already has a `password` key (e.g. from a previous initialization
that was interrupted after writing the Secret but before updating the resource),
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

When the Secret pointed to by `writeConnectionSecretToRef` already contains a `password` key,
the initializer reuses it without calling `generateAndApply()`.
No owner reference is set on the Secret.
Deleting the `User` resource has **no effect on the Secret**; it remains in the cluster,
managed entirely by whoever created it.

In short: **User deleted does not delete the Secret.**

## Internal: `passwordSecretRef`

`passwordSecretRef` is the input channel that Upjet reads on **every reconciliation cycle** to feed
the password to Terraform. Without it, Upjet cannot pass the password to the Atlas API.

The initializer bridges `writeConnectionSecretToRef` and `passwordSecretRef`:
it resolves the password (BYOP or generated) and sets `passwordSecretRef` to point at the right Secret.
After initialization, Upjet reads from `passwordSecretRef` on every reconcile like any other sensitive field.

For simplicity, Users shouldn't need to set `passwordSecretRef` directly.

### `passwordSecretRef` explicitly set

If a user sets `passwordSecretRef` directly, the initializer is bypassed entirely.
Upjet reads the password from `passwordSecretRef` on every reconcile; `writeConnectionSecretToRef` is only
used as an output channel for connection details.

This works without errors, but has two consequences:

- If the two fields point to **different** Secrets, the password sent to Atlas comes from `passwordSecretRef`.
  Any `password` key in the `writeConnectionSecretToRef` Secret is ignored by reconciliation.
- No owner reference is set on either Secret, so **neither is garbage-collected** when the `User` is deleted.

### Edge case: neither field set

If a `User` resource has neither `passwordSecretRef` nor `writeConnectionSecretToRef`,
the initializer exits silently (no error, no generation).
Reconciliation then sends an empty password to Terraform, which may cause the Atlas API to reject the
request or create a user with no password (valid only for non-password auth types like x509 or LDAP).

## Password rotation

The provider has **no built-in rotation scheduler**. Rotation is a manual or
external-automation concern.

### Rotating a password (BYOP or auto-generated)

After initialization, both modes converge: `passwordSecretRef` points at the
`writeConnectionSecretToRef` Secret. To rotate:

1. Update the `password` key in the Secret referenced by
   `writeConnectionSecretToRef`.
2. On the next reconciliation cycle, Upjet detects the changed value and calls
   Terraform to update the Atlas database user.

### Triggering immediate reconciliation

Crossplane reconciles managed resources on a periodic poll interval. To avoid
waiting, any metadata change (e.g. adding a label) triggers the controller's
watch and starts a reconciliation immediately (real-time compositions):

```bash
kubectl label user.database.mongodbatlas.m.crossplane.io my-db-user \
  password-rotated="$(date +%s)" --overwrite
```
