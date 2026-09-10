#!/usr/bin/env bash
#
# Migration script: provider-mongodbatlas v0.x --> v1.x
#
# Prepares a cluster for the v1alpha2 --> v1alpha3 version bump by:
#   1. Backfilling required fields on DatabaseUser CRs
#   2. Triggering re-storage of AdvancedCluster CRs
#   3. Removing v1alpha2 from storedVersions on affected CRDs
#
# A `rollback` phase is also provided to reverse `pre`/`post` using the same
# backups, for when the v1.x provider needs to be abandoned mid-migration.
#
# BACKGROUND
# ----------
# `spec.forProvider.username` does not exist in the v0.x (v1alpha2) schema for
# `users.database.mongodbatlas.crossplane.io` - it is only added in the v1.x
# (v1alpha3) schema, which REPLACES v1alpha2 rather than adding alongside it.
# This creates a real chicken-and-egg problem:
#   - You cannot backfill `username` before v1.x's CRD is installed: the field
#     doesn't exist yet, and structural-schema pruning silently drops the patch
#     ("Warning: unknown field").
#   - You cannot install v1.x's CRD (which only declares v1alpha3) while
#     `status.storedVersions` on the live CRD still lists `v1alpha2`:
#     Kubernetes refuses ("... v1alpha2 was previously a storage version, and
#     must remain in spec.versions until a storage migration ensures no data
#     remains persisted in v1alpha2 ...").
#
# The only way through is a TRANSITIONAL CRD, applied directly with `kubectl
# apply` (bypassing Crossplane's package manager, which only ships the final
# single-version v1.x CRD), that serves BOTH v1alpha2 and v1alpha3 side by
# side. This lets existing objects be re-stored at v1alpha3 (and username
# backfilled) before v1alpha2 is dropped, satisfying Kubernetes' migration
# rule. This script builds that transitional CRD by merging the currently
# installed CRD with the v1.x CRD YAML you provide via --v1-crds-dir.
#
# ALSO IMPORTANT: pausing the leaf DatabaseUser/AdvancedCluster CR
# (`crossplane.io/paused`) only stops the *provider's* reconciler. It does
# NOT stop a Composite Resource (XR) that owns/composes that CR from
# continuing to render and re-apply its (v1alpha2-shaped, username-less)
# `base` template over it - which will silently overwrite any username you
# just backfilled. This script also finds and pauses the controlling owner
# (ownerReferences[] where controller=true) of each affected CR, if any.
#
# Usage:
#   ./scripts/migrate-to-v1.sh pre --v1-crds-dir <dir>
#   # install/activate the v1.x provider package here - see step-by-step
#   # instructions printed at the end of `pre`
#   ./scripts/migrate-to-v1.sh post
#
# To abandon the migration instead of continuing to v1.x (uses the same
# backups `pre` created - do NOT delete /tmp/crossplane-migration-backup
# until you're sure you won't need to roll back):
#   ./scripts/migrate-to-v1.sh rollback
#
# <dir> should contain the v1.x package's CRD YAMLs, e.g. by extracting them
# from the target provider image or checking out the provider-mongodbatlas
# repo at the target tag: package/crds/mongodbatlas.crossplane.io_advancedclusters.yaml
# and package/crds/database.mongodbatlas.crossplane.io_users.yaml.
#
# Dry run (no changes applied):
#   DRY_RUN=true ./scripts/migrate-to-v1.sh pre --v1-crds-dir <dir>
#   DRY_RUN=true ./scripts/migrate-to-v1.sh post
#
# DISCLAIMER: This script is provided "as is", without warranty of any kind,
# express or implied. Use at your own risk. Always test in a non-production
# environment first and ensure you have backups before running this script.
# The authors assume no liability for data loss or cluster disruption.

set -euo pipefail

BACKUP_DIR="/tmp/crossplane-migration-backup"
DRY_RUN="${DRY_RUN:-false}"

USERS_CRD="users.database.mongodbatlas.crossplane.io"
ADVCLUSTERS_CRD="advancedclusters.mongodbatlas.crossplane.io"

AFFECTED_CRDS=(
  "$ADVCLUSTERS_CRD"
  "advancedclusters.mongodbatlas.m.crossplane.io"
  "$USERS_CRD"
  "users.database.mongodbatlas.m.crossplane.io"
)

GREEN='\033[0;32m' ORANGE='\033[0;33m' RED='\033[0;31m' NC='\033[0m'
log()  { echo -e "${GREEN}[migrate]${NC} $*"; }
warn() { echo -e "${GREEN}[migrate]${NC} ${ORANGE}WARNING:${NC} $*" >&2; }
die()  { echo -e "${GREEN}[migrate]${NC} ${RED}ERROR:${NC} $*" >&2; exit 1; }

command -v kubectl >/dev/null || die "kubectl not found"
command -v jq >/dev/null      || die "jq not found"

PHASE="${1:-}"
shift || true
V1_CRDS_DIR=""
while [ $# -gt 0 ]; do
  case "$1" in
    --v1-crds-dir) V1_CRDS_DIR="$2"; shift 2 ;;
    *) die "unknown argument: $1" ;;
  esac
done

case "$PHASE" in
  pre|post|rollback) ;;
  *) die "usage: $0 <pre|post|rollback> [--v1-crds-dir <dir>]  (see script header for details)" ;;
esac

mkdir -p "$BACKUP_DIR"

# -----------------------------------------------------------------------------
# Pause the controlling owner (XR/composite) of a CR, if it has one. Pausing
# only the leaf CR is not enough - the owning composite will keep re-applying
# its (pre-migration) base template over it on every reconcile.
# -----------------------------------------------------------------------------
# Resolve the actual API resource name (plural) for a Kind by looking it up
# in `kubectl api-resources`, rather than guessing lowercase(kind)+"s" - an
# XRD can customize spec.names.singular/plural away from that convention.
# Falls back to the lowercase guess if discovery doesn't find a match.
resolve_resource_name() {
  local kind="$1"
  local resolved
  resolved=$(kubectl api-resources --no-headers -o wide 2>/dev/null | awk -v k="$kind" '$0 ~ ("\\<" k "\\>") { for (i=1;i<=NF;i++) if ($i == k) { print $1; exit } }')
  if [ -n "$resolved" ]; then
    echo "$resolved"
  else
    echo "$kind" | tr '[:upper:]' '[:lower:]'
  fi
}

# Setting the crossplane.io/paused annotation is not synchronous: a
# reconcile that was already in flight (or queued) when the annotation was
# set can still run to completion and re-apply the pre-migration base
# template, clobbering a field we're about to backfill. Wait for the
# composite's own Synced condition to actually report ReconcilePaused
# before considering it safe - this is the same signal Crossplane itself
# uses to confirm a resource has stopped reconciling.
wait_for_pause_confirmed() {
  local resource="$1" name="$2"
  local i reason
  for i in $(seq 1 30); do
    reason=$(kubectl get "$resource" "$name" -o jsonpath='{.status.conditions[?(@.type=="Synced")].reason}' 2>/dev/null || true)
    if [ "$reason" = "ReconcilePaused" ]; then
      return 0
    fi
    sleep 2
  done
  return 1
}

pause_owner() {
  local cr_json="$1"
  local owner
  owner=$(echo "$cr_json" | jq -c '.metadata.ownerReferences[]? | select(.controller == true)' 2>/dev/null || true)
  if [ -z "$owner" ]; then
    return 0
  fi
  local owner_kind owner_name owner_api_version owner_resource
  owner_kind=$(echo "$owner" | jq -r '.kind')
  owner_name=$(echo "$owner" | jq -r '.name')
  owner_api_version=$(echo "$owner" | jq -r '.apiVersion')
  owner_resource=$(resolve_resource_name "$owner_kind")

  log "    Owning composite: $owner_kind/$owner_name ($owner_api_version) -> resource: $owner_resource"
  if [ "$DRY_RUN" = "true" ]; then
    log "    [dry-run] would pause $owner_resource/$owner_name"
  else
    if kubectl annotate "$owner_resource" "$owner_name" crossplane.io/paused=true --overwrite >/dev/null 2>&1; then
      log "    Pause annotation set on $owner_resource/$owner_name, waiting for it to take effect..."
      if wait_for_pause_confirmed "$owner_resource" "$owner_name"; then
        log "    Confirmed $owner_resource/$owner_name has stopped reconciling"
      else
        warn "    $owner_resource/$owner_name did not report ReconcilePaused within 60s - an in-flight reconcile may still land and overwrite backfilled fields. Verify manually before running 'post'."
        echo "$owner_resource/$owner_name ($owner_kind) - pause not confirmed within timeout" >> "$BACKUP_DIR/unpaused-owners.txt"
      fi
    else
      warn "    Could not pause owning composite $owner_kind/$owner_name (tried resource '$owner_resource') - pause it manually before continuing, or it may overwrite backfilled fields."
      echo "$owner_resource/$owner_name ($owner_kind)" >> "$BACKUP_DIR/unpaused-owners.txt"
    fi
  fi
}

pause_affected_and_owners() {
  local kind_crd="$1" # e.g. users.database.mongodbatlas.crossplane.io
  local items
  items=$(kubectl get "$kind_crd" -A -o json 2>/dev/null || echo '{"items":[]}')
  echo "$items" | jq -c '.items[]' 2>/dev/null | while read -r cr; do
    name=$(echo "$cr" | jq -r '.metadata.name')
    ns=$(echo "$cr" | jq -r '.metadata.namespace // empty')
    log "  Pausing $kind_crd/$name (and its owning composite, if any)"
    if [ "$DRY_RUN" = "true" ]; then
      log "    [dry-run] would pause $kind_crd/$name"
    else
      if [ -n "$ns" ]; then
        kubectl annotate "$kind_crd" "$name" -n "$ns" crossplane.io/paused=true --overwrite >/dev/null
      else
        kubectl annotate "$kind_crd" "$name" crossplane.io/paused=true --overwrite >/dev/null
      fi
      if wait_for_pause_confirmed "$kind_crd" "$name"; then
        log "    Confirmed $kind_crd/$name has stopped reconciling"
      else
        warn "    $kind_crd/$name did not report ReconcilePaused within 60s - an in-flight reconcile may still land. Verify manually before running 'post'."
        echo "$kind_crd/$name - pause not confirmed within timeout" >> "$BACKUP_DIR/unpaused-owners.txt"
      fi
    fi
    pause_owner "$cr"
  done
}

# Checks whether the pre-phase actually recorded this CR - or its owning
# composite - as NOT confirmed paused. post consults this to skip touching
# genuinely-live resources rather than racing them and relying on the
# post-patch verification alone: for a claim that's actively reconciling
# (GitOps-managed, a real owner still busy), the right move is to leave it
# alone and surface it for manual handling, not to attempt the patch and
# hope the race doesn't land.
is_unpaused() {
  local cr_json="$1" kind_crd="$2"
  [ -s "$BACKUP_DIR/unpaused-owners.txt" ] || return 1

  local name
  name=$(echo "$cr_json" | jq -r '.metadata.name')
  if grep -qF "$kind_crd/$name " "$BACKUP_DIR/unpaused-owners.txt" 2>/dev/null; then
    return 0
  fi

  local owner owner_kind owner_name owner_resource
  owner=$(echo "$cr_json" | jq -c '.metadata.ownerReferences[]? | select(.controller == true)' 2>/dev/null || true)
  if [ -n "$owner" ]; then
    owner_kind=$(echo "$owner" | jq -r '.kind')
    owner_name=$(echo "$owner" | jq -r '.name')
    owner_resource=$(resolve_resource_name "$owner_kind")
    if grep -qF "$owner_resource/$owner_name " "$BACKUP_DIR/unpaused-owners.txt" 2>/dev/null; then
      return 0
    fi
  fi

  return 1
}

# Strips fields tied to the specific object instance currently on the
# cluster (resourceVersion, uid, managedFields, ...) so this backup can
# later be fed straight to `kubectl apply -f` (as `rollback` does) without
# it being rejected as a conflicting update against whatever resourceVersion
# the object has moved on to by then.
strip_instance_metadata() {
  jq 'del(.metadata.resourceVersion, .metadata.uid, .metadata.generation, .metadata.creationTimestamp, .metadata.managedFields, .metadata.selfLink, .status)'
}

backup_users() {
  db_users=$(kubectl get "$USERS_CRD" -A -o json 2>/dev/null || echo '{"items":[]}')
  echo "$db_users" | jq -c '.items[]' 2>/dev/null | while read -r cr; do
    name=$(echo "$cr" | jq -r '.metadata.name')
    ns=$(echo "$cr" | jq -r '.metadata.namespace // "cluster-scoped"')
    if [ "$ns" = "cluster-scoped" ]; then
      kubectl get "$USERS_CRD" "$name" -o json | strip_instance_metadata > "$BACKUP_DIR/dbuser-$name.yaml"
    else
      kubectl get "$USERS_CRD" "$name" -n "$ns" -o json | strip_instance_metadata > "$BACKUP_DIR/dbuser-$ns-$name.yaml"
    fi
  done
}

backup_advancedclusters() {
  adv_clusters=$(kubectl get "$ADVCLUSTERS_CRD" -A -o json 2>/dev/null || echo '{"items":[]}')
  echo "$adv_clusters" | jq -c '.items[]' 2>/dev/null | while read -r cr; do
    name=$(echo "$cr" | jq -r '.metadata.name')
    ns=$(echo "$cr" | jq -r '.metadata.namespace // "cluster-scoped"')
    if [ "$ns" = "cluster-scoped" ]; then
      kubectl get "$ADVCLUSTERS_CRD" "$name" -o json | strip_instance_metadata > "$BACKUP_DIR/advcluster-$name.yaml"
    else
      kubectl get "$ADVCLUSTERS_CRD" "$name" -n "$ns" -o json | strip_instance_metadata > "$BACKUP_DIR/advcluster-$ns-$name.yaml"
    fi
  done
}

# -----------------------------------------------------------------------------
# Build a transitional CRD (both v1alpha2 and v1alpha3 served, v1alpha3 as
# storage) by merging the live CRD with the target v1.x CRD's v1alpha3
# version block. Requires `python3` with PyYAML.
# -----------------------------------------------------------------------------
build_transitional_crd() {
  local live_crd_name="$1" v1_crd_file="$2" out_file="$3"
  command -v python3 >/dev/null || die "python3 (with PyYAML) is required to build the transitional CRD"
  kubectl get crd "$live_crd_name" -o yaml > "$BACKUP_DIR/live-crd-$live_crd_name.yaml"
  python3 - "$BACKUP_DIR/live-crd-$live_crd_name.yaml" "$v1_crd_file" "$out_file" <<'PYEOF'
import sys, yaml

live_path, v1_path, out_path = sys.argv[1:4]

with open(live_path) as f:
    live = yaml.safe_load(f)
with open(v1_path) as f:
    v1 = yaml.safe_load(f)

if len(v1['spec']['versions']) != 1:
    sys.exit(f"expected exactly one version in {v1_path}, found {[v['name'] for v in v1['spec']['versions']]}")

new_version = v1['spec']['versions'][0]
live_version_names = [v['name'] for v in live['spec']['versions']]

# Idempotency: if a previous `pre` run already applied the transitional CRD
# (e.g. it succeeded here but died on the other CRD, or the whole phase is
# simply being re-run), the live CRD already has both versions. Detect that
# and re-emit the already-correct shape instead of failing, so `pre` is safe
# to re-run after a partial failure.
if len(live['spec']['versions']) == 2 and new_version['name'] in live_version_names:
    for v in live['spec']['versions']:
        if v['name'] == new_version['name']:
            if not (v.get('served') and v.get('storage')):
                sys.exit(
                    f"{live_path} already has {new_version['name']} but it is not "
                    f"served+storage (served={v.get('served')}, storage={v.get('storage')}) - "
                    "refusing to guess; inspect and fix manually."
                )
        elif v.get('storage'):
            sys.exit(
                f"{live_path} already has 2 versions but the old version "
                f"({v['name']}) is still marked storage=true - refusing to guess; "
                "inspect and fix manually."
            )
    print(f"{live_path} is already transitional ({live_version_names}) - reusing it as-is -> {out_path}")
    for k in ('resourceVersion', 'uid', 'creationTimestamp', 'generation', 'managedFields', 'selfLink'):
        live['metadata'].pop(k, None)
    live.pop('status', None)
    with open(out_path, 'w') as f:
        yaml.safe_dump(live, f, default_flow_style=False, sort_keys=False)
    sys.exit(0)

if len(live['spec']['versions']) != 1:
    sys.exit(f"expected exactly one existing version on {live_path}, found {live_version_names}")

old_version = live['spec']['versions'][0]

old_version['served'] = True
old_version['storage'] = False
new_version['served'] = True
new_version['storage'] = True

live['spec']['versions'] = [old_version, new_version]

for k in ('resourceVersion', 'uid', 'creationTimestamp', 'generation', 'managedFields', 'selfLink'):
    live['metadata'].pop(k, None)
live.pop('status', None)

with open(out_path, 'w') as f:
    yaml.safe_dump(live, f, default_flow_style=False, sort_keys=False)

print(f"built transitional CRD: {[v['name'] for v in live['spec']['versions']]} -> {out_path}")
PYEOF
}

# -----------------------------------------------------------------------------
# Build the CRD shape needed to roll back: both the old version (from the
# pre-phase's own backup of the original, single-version CRD) and whatever
# version(s) the live CRD currently has are served, but the OLD version is
# now storage=true again. This mirrors build_transitional_crd in reverse.
#
# Narrowing spec.versions straight back down to the old version alone -
# without going through this re-widen step - looks like it should work
# (the old version's schema is unchanged) but corrupts every affected
# object: confirmed live, `kubectl get`/`patch` on an object that was ever
# served as the newer version, after that version is removed from
# spec.versions, fails with "request to convert CR from an invalid
# group/version: <new-version>" - even when no field-level write under the
# new version ever happened. Kubernetes' storedVersions/spec.versions
# migration rule exists precisely to prevent this; rolling back has to
# satisfy it in reverse, not skip it.
# -----------------------------------------------------------------------------
build_rollback_crd() {
  local live_crd_name="$1" backup_crd_file="$2" out_file="$3"
  kubectl get crd "$live_crd_name" -o yaml > "$BACKUP_DIR/rollback-live-$live_crd_name.yaml"
  python3 - "$BACKUP_DIR/rollback-live-$live_crd_name.yaml" "$backup_crd_file" "$out_file" <<'PYEOF'
import sys, yaml

live_path, backup_path, out_path = sys.argv[1:4]

with open(live_path) as f:
    live = yaml.safe_load(f)
with open(backup_path) as f:
    backup = yaml.safe_load(f)

if len(backup['spec']['versions']) != 1:
    sys.exit(f"expected exactly one version in {backup_path}, found {[v['name'] for v in backup['spec']['versions']]}")

old_version = backup['spec']['versions'][0]
old_version['served'] = True
old_version['storage'] = True

# Keep serving whatever else is currently live (so objects already stored
# under it stay readable while they're being restored) but it's no longer
# the storage version.
other_versions = [v for v in live['spec']['versions'] if v['name'] != old_version['name']]
for v in other_versions:
    v['served'] = True
    v['storage'] = False

live['spec']['versions'] = [old_version] + other_versions

for k in ('resourceVersion', 'uid', 'creationTimestamp', 'generation', 'managedFields', 'selfLink'):
    live['metadata'].pop(k, None)
live.pop('status', None)

with open(out_path, 'w') as f:
    yaml.safe_dump(live, f, default_flow_style=False, sort_keys=False)

print(f"built rollback CRD: {[v['name'] for v in live['spec']['versions']]} (storage={old_version['name']}) -> {out_path}")
PYEOF
}

drop_old_version_from_crd() {
  local crd_name="$1" old_version_name="$2" out_file="$3"
  kubectl get crd "$crd_name" -o yaml > "$out_file.orig"
  python3 - "$out_file.orig" "$old_version_name" "$out_file" <<'PYEOF'
import sys, yaml
in_path, old_name, out_path = sys.argv[1:4]
with open(in_path) as f:
    crd = yaml.safe_load(f)
crd['spec']['versions'] = [v for v in crd['spec']['versions'] if v['name'] != old_name]
for k in ('resourceVersion', 'uid', 'creationTimestamp', 'generation', 'managedFields', 'selfLink'):
    crd['metadata'].pop(k, None)
crd.pop('status', None)
with open(out_path, 'w') as f:
    yaml.safe_dump(crd, f, default_flow_style=False, sort_keys=False)
PYEOF
}

if [ "$PHASE" = "pre" ]; then
  [ -n "$V1_CRDS_DIR" ] || die "pre phase requires --v1-crds-dir <dir> pointing at the v1.x provider's package/crds/ directory"
  [ -f "$V1_CRDS_DIR/mongodbatlas.crossplane.io_advancedclusters.yaml" ] || die "missing $V1_CRDS_DIR/mongodbatlas.crossplane.io_advancedclusters.yaml"
  [ -f "$V1_CRDS_DIR/database.mongodbatlas.crossplane.io_users.yaml" ] || die "missing $V1_CRDS_DIR/database.mongodbatlas.crossplane.io_users.yaml"

  log "Pre-phase 1/3: backing up DatabaseUser and AdvancedCluster CRs..."
  if [ "$DRY_RUN" = "true" ]; then
    log "  [dry-run] would back up all DatabaseUser and AdvancedCluster CRs to $BACKUP_DIR"
  else
    backup_users
    backup_advancedclusters
  fi

  log "Pre-phase 2/3: pausing DatabaseUser/AdvancedCluster CRs and their owning composites..."
  rm -f "$BACKUP_DIR/unpaused-owners.txt"
  pause_affected_and_owners "$USERS_CRD"
  pause_affected_and_owners "$ADVCLUSTERS_CRD"

  log "Pre-phase 3/3: building and applying transitional (dual-version) CRDs..."
  if [ "$DRY_RUN" = "true" ]; then
    log "  [dry-run] would build and apply transitional CRDs for $USERS_CRD and $ADVCLUSTERS_CRD"
  else
    build_transitional_crd "$USERS_CRD" "$V1_CRDS_DIR/database.mongodbatlas.crossplane.io_users.yaml" "$BACKUP_DIR/transitional-users.yaml"
    build_transitional_crd "$ADVCLUSTERS_CRD" "$V1_CRDS_DIR/mongodbatlas.crossplane.io_advancedclusters.yaml" "$BACKUP_DIR/transitional-advancedclusters.yaml"
    kubectl apply -f "$BACKUP_DIR/transitional-users.yaml"
    kubectl apply -f "$BACKUP_DIR/transitional-advancedclusters.yaml"
  fi

  log ""
  log "Pre-phase complete. Backups saved to $BACKUP_DIR"
  log "The live CRDs now serve BOTH v1alpha2 and v1alpha3 (v1alpha3 is storage)."
  if [ -s "$BACKUP_DIR/unpaused-owners.txt" ]; then
    log ""
    warn "$(wc -l < "$BACKUP_DIR/unpaused-owners.txt" | tr -d ' ') owning composite(s) could NOT be paused automatically:"
    while read -r line; do warn "  - $line"; done < "$BACKUP_DIR/unpaused-owners.txt"
    warn "Pause these manually before running '$0 post', or they may silently overwrite backfilled fields."
  fi
  log ""
  log "IMPORTANT: do NOT install/activate the v1.x provider package yet - its"
  log "single-version CRDs would be REJECTED right now, because v1alpha2 is"
  log "still in status.storedVersions. Instead:"
  log "  1. Run: $0 post"
  log "     (this backfills username/authDatabaseName and re-stores every"
  log "     affected CR at v1alpha3, then removes v1alpha2 from storedVersions"
  log "     and drops it from spec.versions, matching the v1.x CRD shape)"
  log "  2. THEN install/activate the v1.x provider package."
  log "  3. Once it is Healthy, unpause the CRs and their owning composites"
  log "     that this script paused."
  exit 0
fi

if [ "$PHASE" = "rollback" ]; then
  for crd in "$USERS_CRD" "$ADVCLUSTERS_CRD"; do
    [ -f "$BACKUP_DIR/live-crd-$crd.yaml" ] || die "no pre-phase CRD backup found at $BACKUP_DIR/live-crd-$crd.yaml - rollback requires the backups '$0 pre' created for $crd. Cannot roll back without them."
  done

  log "Rollback 1/3: re-widening CRDs (both versions served, v1alpha2 as storage)..."
  if [ "$DRY_RUN" = "true" ]; then
    log "  [dry-run] would re-widen $USERS_CRD and $ADVCLUSTERS_CRD, storage=v1alpha2"
  else
    build_rollback_crd "$USERS_CRD" "$BACKUP_DIR/live-crd-$USERS_CRD.yaml" "$BACKUP_DIR/rollback-transitional-users.yaml"
    build_rollback_crd "$ADVCLUSTERS_CRD" "$BACKUP_DIR/live-crd-$ADVCLUSTERS_CRD.yaml" "$BACKUP_DIR/rollback-transitional-advancedclusters.yaml"
    kubectl apply -f "$BACKUP_DIR/rollback-transitional-users.yaml"
    kubectl apply -f "$BACKUP_DIR/rollback-transitional-advancedclusters.yaml"
  fi

  log "Rollback 2/3: restoring DatabaseUser/AdvancedCluster CRs from their pre-migration backups..."
  restore_failed=""
  for f in "$BACKUP_DIR"/dbuser-*.yaml "$BACKUP_DIR"/advcluster-*.yaml; do
    [ -f "$f" ] || continue
    log "  Restoring $(basename "$f")"
    if [ "$DRY_RUN" = "true" ]; then
      log "    [dry-run] would apply $f"
    else
      # --server-side avoids a client-side-apply quirk seen live against this
      # provider's User CRD: kubectl's client-side JSON-merge-patch fallback
      # (triggered when the OpenAPI V3 path doesn't support strategic merge,
      # as with this provider's CRDs) can send a literal resourceVersion "0"
      # even though the local file has no resourceVersion at all, which the
      # apiserver rejects with "metadata.resourceVersion: Invalid value: 0".
      # Server-side apply doesn't go through that code path.
      kubectl apply --server-side --force-conflicts -f "$f" || { warn "  failed to restore $f"; restore_failed="yes"; }
    fi
  done
  [ -n "$restore_failed" ] && die "one or more CRs failed to restore from backup - inspect $BACKUP_DIR and fix manually. Do NOT proceed to narrow the CRDs back down while any object might still be storing v1alpha3 data."

  log "Rollback 3/3: removing v1alpha3 from storedVersions and dropping it from spec.versions..."
  if [ "$DRY_RUN" = "true" ]; then
    log "  [dry-run] would narrow $USERS_CRD and $ADVCLUSTERS_CRD back to v1alpha2-only"
  else
    for crd in "$USERS_CRD" "$ADVCLUSTERS_CRD"; do
      stored=$(kubectl get crd "$crd" -o jsonpath='{.status.storedVersions}') || die "cannot read storedVersions for $crd"
      if echo "$stored" | grep -q 'v1alpha3'; then
        kubectl patch crd "$crd" \
          --type='json' \
          --subresource=status \
          -p='[{"op":"replace","path":"/status/storedVersions","value":["v1alpha2"]}]'
      fi
      drop_old_version_from_crd "$crd" "v1alpha3" "$BACKUP_DIR/rollback-final-$crd.yaml"
      kubectl apply -f "$BACKUP_DIR/rollback-final-$crd.yaml"
    done
  fi

  log ""
  log "Rollback complete. $USERS_CRD and $ADVCLUSTERS_CRD now serve only"
  log "v1alpha2 again, matching the v0.x provider's shape. Remaining manual steps:"
  log "  1. Re-install/activate the v0.x provider package if you had already"
  log "     switched to v1.x."
  log "  2. Unpause the CRs and their owning composites that '$0 pre' paused"
  log "     (annotate crossplane.io/paused- to remove)."
  exit 0
fi

# ---------------------------------------------------------------------------
# Post-phase: the live CRDs now serve v1alpha3 (from the pre-phase). Backfill
# required fields, re-store every CR at v1alpha3, then finish the CRD
# transition so it matches what the v1.x provider package will install.
# ---------------------------------------------------------------------------

require_transitional_crd() {
  local crd_name="$1"
  local versions
  versions=$(kubectl get crd "$crd_name" -o jsonpath='{.spec.versions[*].name}' 2>/dev/null) \
    || die "cannot read CRD $crd_name - does it exist? Run '$0 pre --v1-crds-dir <dir>' first."
  echo "$versions" | grep -qw 'v1alpha3' \
    || die "$crd_name does not serve v1alpha3 yet - run '$0 pre --v1-crds-dir <dir>' first, it builds the transitional CRD that 'post' depends on."
}

require_transitional_crd "$USERS_CRD"
require_transitional_crd "$ADVCLUSTERS_CRD"

log "Post-phase Step 1: Patching DatabaseUser CRs..."
rm -f "$BACKUP_DIR/user-backfill-failed.txt" "$BACKUP_DIR/advcluster-migration-failed.txt"

db_users=$(kubectl get "$USERS_CRD" -A -o json 2>/dev/null || echo '{"items":[]}')
user_count=$(echo "$db_users" | jq '.items | length')
log "  Found $user_count DatabaseUser CR(s)"

echo "$db_users" | jq -c '.items[]' | while read -r cr; do
  name=$(echo "$cr" | jq -r '.metadata.name')
  ns=$(echo "$cr" | jq -r '.metadata.namespace // "cluster-scoped"')
  ext_name=$(echo "$cr" | jq -r '.metadata.annotations["crossplane.io/external-name"] // empty')
  existing_username=$(echo "$cr" | jq -r '.spec.forProvider.username // empty')
  existing_auth_db=$(echo "$cr" | jq -r '.spec.forProvider.authDatabaseName // empty')

  log "  Processing DatabaseUser $ns/$name (external-name: $ext_name)"

  if [ -n "$existing_username" ]; then
    log "    username already set ($existing_username), skipping"
    continue
  fi

  if is_unpaused "$cr" "$USERS_CRD"; then
    warn "    $ns/$name (or its owning composite) was not confirmed paused in the pre-phase - skipping rather than racing a live reconcile. Pause it manually and re-run 'post' for this CR."
    echo "$USERS_CRD/$name ($ns) - skipped, pause not confirmed in pre-phase" >> "$BACKUP_DIR/user-backfill-failed.txt"
    continue
  fi

  if [ -z "$ext_name" ]; then
    warn "    No external-name annotation on $name, cannot infer username, skipping"
    continue
  fi

  # In v0.x the external-name format was: the raw username (via NameAsIdentifier)
  username="$ext_name"
  log "    Setting spec.forProvider.username = $username"

  patch_fields="\"username\":\"$username\""
  if [ -z "$existing_auth_db" ]; then
    log "    Setting spec.forProvider.authDatabaseName = admin (default)"
    patch_fields="$patch_fields,\"authDatabaseName\":\"admin\""
  fi
  patch="{\"spec\":{\"forProvider\":{$patch_fields}}}"

  if [ "$DRY_RUN" = "true" ]; then
    log "    [dry-run] would patch: $patch"
  else
    if [ "$ns" = "cluster-scoped" ]; then
      kubectl patch "$USERS_CRD" "$name" --type=merge -p "$patch"
    else
      kubectl patch "$USERS_CRD" "$name" -n "$ns" --type=merge -p "$patch"
    fi

    # The patch above can report success and still not stick: if the leaf CR
    # has a controlling composite (see pause_owner in the pre-phase), an
    # in-flight or unpaused reconcile of that composite can re-apply its
    # pre-v1 base template over this CR immediately afterward, silently
    # wiping the field we just set. Re-read the object a few times rather
    # than trusting the patch's own exit code, and fail loudly - not just
    # warn - if it never sticks, so this isn't discovered only after v1.x
    # is already installed and the CR fails CEL validation.
    persisted=""
    for attempt in $(seq 1 5); do
      if [ "$ns" = "cluster-scoped" ]; then
        persisted=$(kubectl get "$USERS_CRD" "$name" -o jsonpath='{.spec.forProvider.username}' 2>/dev/null || true)
      else
        persisted=$(kubectl get "$USERS_CRD" "$name" -n "$ns" -o jsonpath='{.spec.forProvider.username}' 2>/dev/null || true)
      fi
      [ "$persisted" = "$username" ] && break
      sleep 2
    done
    if [ "$persisted" = "$username" ]; then
      log "    Verified: username persisted as $username"
    else
      warn "    username on $ns/$name did not persist after patch (got \"$persisted\", wanted \"$username\") - its owning composite is likely still reconciling despite the pre-phase pause. Check $BACKUP_DIR/unpaused-owners.txt and the composite's crossplane.io/paused status before proceeding."
      echo "$USERS_CRD/$name ($ns) - username backfill did not persist (got \"$persisted\", wanted \"$username\")" >> "$BACKUP_DIR/user-backfill-failed.txt"
    fi
  fi
done

# ---------------------------------------------------------------------------
# Post-phase Step 2: Migrate AdvancedCluster region-config specs to v1alpha3
# ---------------------------------------------------------------------------
# v1alpha3 is not just a schema addition: electableSpecs, autoScaling,
# analyticsAutoScaling, readOnlySpecs and analyticsSpecs all change from a
# single-item array (v1alpha2) to a plain object (v1alpha3), and the
# top-level diskSizeGb field is removed in favor of a diskSizeGb property on
# electableSpecs/readOnlySpecs/analyticsSpecs individually. A bare "touch"
# to trigger re-storage relies on CRD version conversion to reshape this -
# but this provider has no conversion webhook, so Kubernetes' default (None)
# conversion is a pure passthrough: it re-validates the SAME array-shaped
# JSON against the new object-shaped schema, and structural-schema pruning
# silently drops the whole mismatched subtree rather than restructuring it.
# So this step reads each CR's pre-migration shape from the pre-phase's own
# backup, builds the reshaped v1alpha3 structure explicitly, and patches
# that in - it does not depend on conversion doing any work at all.
log "Post-phase Step 2: Migrating AdvancedCluster region-config specs..."

adv_clusters=$(kubectl get "$ADVCLUSTERS_CRD" -A -o json 2>/dev/null || echo '{"items":[]}')
ac_count=$(echo "$adv_clusters" | jq '.items | length')
log "  Found $ac_count AdvancedCluster CR(s)"

build_advancedcluster_patch() {
  local backup_file="$1" out_file="$2" cluster_name="$3"
  python3 - "$backup_file" "$out_file" "$cluster_name" <<'PYEOF'
import sys, json, yaml

backup_path, out_path, cluster_name = sys.argv[1:4]

with open(backup_path) as f:
    cr = yaml.safe_load(f)

fp = cr.get("spec", {}).get("forProvider", {})
disk_size_gb = fp.get("diskSizeGb")

# These five fields are arrays-of-one-object in v1alpha2 and plain objects
# in v1alpha3. Only electableSpecs/readOnlySpecs/analyticsSpecs gain a
# diskSizeGb property - autoScaling/analyticsAutoScaling never had one.
ARRAY_TO_OBJECT_FIELDS = [
    "electableSpecs",
    "autoScaling",
    "analyticsAutoScaling",
    "readOnlySpecs",
    "analyticsSpecs",
]
GETS_DISK_SIZE = {"electableSpecs", "readOnlySpecs", "analyticsSpecs"}

new_replication_specs = []
for rs in fp.get("replicationSpecs", []) or []:
    new_rcs = []
    for rc in rs.get("regionConfigs", []) or []:
        new_rc = {}
        for k, v in rc.items():
            if k not in ARRAY_TO_OBJECT_FIELDS:
                new_rc[k] = v
        for field in ARRAY_TO_OBJECT_FIELDS:
            items = rc.get(field) or []
            if not items:
                continue
            if len(items) != 1:
                sys.exit(
                    f"expected exactly one item in regionConfigs[].{field} "
                    f"(v1alpha2 always emits one), found {len(items)} in {backup_path} - "
                    "refusing to guess which one is authoritative; inspect and migrate this CR manually."
                )
            obj = dict(items[0])
            if field in GETS_DISK_SIZE and disk_size_gb is not None:
                obj["diskSizeGb"] = disk_size_gb
            new_rc[field] = obj
        new_rcs.append(new_rc)
    new_rs = {k: v for k, v in rs.items() if k != "regionConfigs"}
    new_rs["regionConfigs"] = new_rcs
    new_replication_specs.append(new_rs)

# JSON merge patch: null deletes a key. diskSizeGb has no home at the
# top level in v1alpha3 - its value has already been distributed above.
# spec.forProvider.name is new in v1alpha3 and is a required parameter
# there (validated by a CEL rule, not present at all in v1alpha2) - the
# caller derives it from this CR's own external-name, matching exactly how
# the v1alpha2 mongodbatlas_advanced_cluster resource config already builds
# and parses that annotation (SetIdentifierFunc/GetExternalNameFn in
# config/mongodbatlas/config.go): "<name>:<cluster_id>".
patch = {"spec": {"forProvider": {"replicationSpecs": new_replication_specs, "diskSizeGb": None}}}
if cluster_name:
    patch["spec"]["forProvider"]["name"] = cluster_name

with open(out_path, "w") as f:
    json.dump(patch, f)
PYEOF
}

echo "$adv_clusters" | jq -c '.items[]' | while read -r cr; do
  name=$(echo "$cr" | jq -r '.metadata.name')
  ns=$(echo "$cr" | jq -r '.metadata.namespace // "cluster-scoped"')
  ext_name=$(echo "$cr" | jq -r '.metadata.annotations["crossplane.io/external-name"] // empty')
  existing_name=$(echo "$cr" | jq -r '.spec.forProvider.name // empty')

  if is_unpaused "$cr" "$ADVCLUSTERS_CRD"; then
    warn "  $ns/$name (or its owning composite) was not confirmed paused in the pre-phase - skipping rather than racing a live reconcile. Pause it manually and re-run 'post' for this CR."
    echo "$ADVCLUSTERS_CRD/$name ($ns) - skipped, pause not confirmed in pre-phase" >> "$BACKUP_DIR/advcluster-migration-failed.txt"
    continue
  fi

  if [ "$ns" = "cluster-scoped" ]; then
    backup_file="$BACKUP_DIR/advcluster-$name.yaml"
  else
    backup_file="$BACKUP_DIR/advcluster-$ns-$name.yaml"
  fi
  [ -f "$backup_file" ] || die "no pre-phase backup found at $backup_file for AdvancedCluster $ns/$name - re-run '$0 pre --v1-crds-dir <dir>' first, Step 2 depends on it to know the pre-migration shape."

  log "  Migrating AdvancedCluster $ns/$name region-config specs (from $backup_file)"

  cluster_name="$existing_name"
  if [ -z "$cluster_name" ]; then
    # external-name format is "<name>:<cluster_id>" (see config/mongodbatlas/config.go).
    cluster_name="${ext_name%%:*}"
    if [ -z "$cluster_name" ]; then
      warn "    No external-name annotation on $name, cannot infer spec.forProvider.name - the CEL required-field check on it will fail once v1.x is installed. Set it manually."
    else
      log "    Setting spec.forProvider.name = $cluster_name (from external-name)"
    fi
  fi

  patch_file="$BACKUP_DIR/advcluster-patch-$ns-$name.json"
  build_advancedcluster_patch "$backup_file" "$patch_file" "$cluster_name"

  if [ "$DRY_RUN" = "true" ]; then
    log "    [dry-run] would apply reshaped patch: $(cat "$patch_file")"
    continue
  fi

  if [ "$ns" = "cluster-scoped" ]; then
    kubectl patch "$ADVCLUSTERS_CRD" "$name" --type=merge -p "$(cat "$patch_file")"
  else
    kubectl patch "$ADVCLUSTERS_CRD" "$name" -n "$ns" --type=merge -p "$(cat "$patch_file")"
  fi

  # Verify the patch actually stuck - both the reshaped electableSpecs and
  # the backfilled name - for the same reason Step 1 verifies the username
  # patch: an owning composite re-applying its pre-v1 base template can
  # silently overwrite either.
  if [ "$ns" = "cluster-scoped" ]; then
    got_instance_size=$(kubectl get "$ADVCLUSTERS_CRD" "$name" -o jsonpath='{.spec.forProvider.replicationSpecs[0].regionConfigs[0].electableSpecs.instanceSize}' 2>/dev/null || true)
    got_name=$(kubectl get "$ADVCLUSTERS_CRD" "$name" -o jsonpath='{.spec.forProvider.name}' 2>/dev/null || true)
  else
    got_instance_size=$(kubectl get "$ADVCLUSTERS_CRD" "$name" -n "$ns" -o jsonpath='{.spec.forProvider.replicationSpecs[0].regionConfigs[0].electableSpecs.instanceSize}' 2>/dev/null || true)
    got_name=$(kubectl get "$ADVCLUSTERS_CRD" "$name" -n "$ns" -o jsonpath='{.spec.forProvider.name}' 2>/dev/null || true)
  fi
  if [ -z "$got_instance_size" ]; then
    warn "    electableSpecs on $ns/$name is empty after patching - its owning composite is likely still reconciling despite the pre-phase pause. Check $BACKUP_DIR/unpaused-owners.txt and the composite's crossplane.io/paused status before proceeding."
    echo "$ADVCLUSTERS_CRD/$name ($ns) - electableSpecs empty after migration patch" >> "$BACKUP_DIR/advcluster-migration-failed.txt"
  else
    log "    Verified: electableSpecs.instanceSize = $got_instance_size"
  fi
  if [ -n "$cluster_name" ] && [ "$got_name" != "$cluster_name" ]; then
    warn "    name on $ns/$name did not persist after patch (got \"$got_name\", wanted \"$cluster_name\") - its owning composite is likely still reconciling despite the pre-phase pause."
    echo "$ADVCLUSTERS_CRD/$name ($ns) - name backfill did not persist (got \"$got_name\", wanted \"$cluster_name\")" >> "$BACKUP_DIR/advcluster-migration-failed.txt"
  elif [ -n "$cluster_name" ]; then
    log "    Verified: name persisted as $got_name"
  fi
done

# Tracked in separate files (rather than one shared list) so this message
# names the resource type that actually failed - Step 1 (DatabaseUser) and
# Step 2 (AdvancedCluster) failures are unrelated and were previously easy
# to conflate: a DatabaseUser-only failure here used to be reported as an
# "AdvancedCluster region-config migration" failure, sending whoever's
# debugging it to inspect the wrong resource.
post_failures=()
[ -s "$BACKUP_DIR/user-backfill-failed.txt" ] && post_failures+=("DatabaseUser username/authDatabaseName backfill - see $BACKUP_DIR/user-backfill-failed.txt")
[ -s "$BACKUP_DIR/advcluster-migration-failed.txt" ] && post_failures+=("AdvancedCluster region-config migration - see $BACKUP_DIR/advcluster-migration-failed.txt")
if [ "${#post_failures[@]}" -gt 0 ]; then
  warn "Post-phase verification failures:"
  for f in "${post_failures[@]}"; do warn "  - $f"; done
  die "one or more backfills did not persist. Do not install/activate the v1.x provider package until these are resolved."
fi

# ---------------------------------------------------------------------------
# Post-phase Step 3: Remove v1alpha2 from storedVersions on affected CRDs,
# then drop v1alpha2 from spec.versions entirely so the CRD shape matches
# what the v1.x provider package will install (which only declares v1alpha3).
# ---------------------------------------------------------------------------
log "Post-phase Step 3: Finishing CRD version transition..."

for crd in "${AFFECTED_CRDS[@]}"; do
  get_err=$(kubectl get crd "$crd" 2>&1 >/dev/null) && get_status=0 || get_status=$?
  if [ "$get_status" -ne 0 ]; then
    if echo "$get_err" | grep -qi 'NotFound'; then
      log "  $crd: not installed, skipping"
      continue
    fi
    die "cannot check CRD $crd: $get_err"
  fi

  stored=$(kubectl get crd "$crd" -o jsonpath='{.status.storedVersions}') || die "cannot read storedVersions for $crd"
  if echo "$stored" | grep -q 'v1alpha2'; then
    log "  $crd: removing v1alpha2 from status.storedVersions"
    if [ "$DRY_RUN" = "true" ]; then
      log "    [dry-run] would set storedVersions=[\"v1alpha3\"]"
    else
      kubectl patch crd "$crd" \
        --type='json' \
        --subresource=status \
        -p='[{"op":"replace","path":"/status/storedVersions","value":["v1alpha3"]}]'
    fi
  else
    log "  $crd: v1alpha2 not in storedVersions, skipping"
  fi

  versions=$(kubectl get crd "$crd" -o jsonpath='{.spec.versions[*].name}') || die "cannot read spec.versions for $crd"
  if echo "$versions" | grep -qw 'v1alpha2'; then
    log "  $crd: dropping v1alpha2 from spec.versions"
    if [ "$DRY_RUN" = "true" ]; then
      log "    [dry-run] would drop v1alpha2 from spec.versions"
    else
      drop_old_version_from_crd "$crd" "v1alpha2" "$BACKUP_DIR/final-$crd.yaml"
      kubectl apply -f "$BACKUP_DIR/final-$crd.yaml"
    fi
  else
    log "  $crd: v1alpha2 already absent from spec.versions, skipping"
  fi
done

log ""
log "Post-phase complete. Backups (from the pre-phase run) are in $BACKUP_DIR"
log "The CRDs now match the shape the v1.x provider package will install."
log "You can now install/activate the v1.x provider package."
log "Once it is Healthy, unpause the CRs and their owning composites that"
log "were paused in the pre-phase (annotate crossplane.io/paused- to remove)."
