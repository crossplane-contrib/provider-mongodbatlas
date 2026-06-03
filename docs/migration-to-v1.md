# Migration Guide to v1

> **Disclaimer:** This script and guide are provided "as is", without warranty of any kind, express
> or implied. Use at your own risk. Always test in a non-production environment first and ensure you
> have backups before running migration operations against your cluster.

## Breaking Changes

v1.0.0 upgrades these CRDs from `v1alpha2` to `v1alpha3`:

| CRD | Breaking Change |
|-----|----------------|
| `advancedclusters.mongodbatlas.crossplane.io` | Version bump only (schema compatible) |
| `users.database.mongodbatlas.crossplane.io` | `spec.forProvider.username` is now a required field. Previously inferred from `metadata.name` via `crossplane.io/external-name` |

## Prerequisites

- `kubectl` access to the cluster running `provider-mongodbatlas` v0.x.
- `jq` installed.
- Provider v1.x package ready but **not yet installed**.

## Migration Steps

Run the migration script **before** upgrading the provider:

```bash
./scripts/migrate-to-v1.sh
```

The script will:
1. Patch all `DatabaseUser` CRs to add `spec.forProvider.username` from `crossplane.io/external-name`.
2. Patch all `DatabaseUser` CRs to set `spec.forProvider.authDatabaseName` from `crossplane.io/external-name` (if missing).
3. Touch every `AdvancedCluster` CR to trigger re-storage at the current version.
4. Patch `status.storedVersions` on both CRDs to remove `v1alpha2`.

After running the script, upgrade the provider to v1.x.

## Dry Run

Preview changes without modifying anything:

```bash
DRY_RUN=true ./scripts/migrate-to-v1.sh
```

Example output:

```
[migrate] Step 1: Patching DatabaseUser CRs...
[migrate]   Found 2 DatabaseUser CR(s)
[migrate]   Processing default/my-db-user (external-name: admin)
[migrate]     [dry-run] would patch username=admin
[migrate]   Processing default/my-other-user (external-name: readonly)
[migrate]     [dry-run] would patch username=readonly
[migrate] Step 2: Triggering re-storage of AdvancedCluster CRs...
[migrate]   Found 1 AdvancedCluster CR(s)
[migrate]   Touching default/my-cluster
[migrate]     [dry-run] would touch my-cluster
[migrate] Step 3: Removing v1alpha2 from storedVersions...
[migrate]   advancedclusters.mongodbatlas.crossplane.io: removing v1alpha2 from storedVersions
[migrate]     [dry-run] would set storedVersions=["v1alpha3"]
[migrate]   users.database.mongodbatlas.crossplane.io: removing v1alpha2 from storedVersions
[migrate]     [dry-run] would set storedVersions=["v1alpha3"]
[migrate]
[migrate] Migration complete. Backups saved to /tmp/crossplane-migration-backup
[migrate] You can now upgrade provider-mongodbatlas to v1.x.
```

## Rollback

The script creates a backup of each modified CR in `/tmp/crossplane-migration-backup/` before patching. To restore:

```bash
for f in /tmp/crossplane-migration-backup/*.yaml; do
  kubectl apply -f "$f"
done
```
