#!/usr/bin/env bash
#
# Migration script: provider-mongodbatlas v0.x --> v1.x
#
# Prepares a cluster for the v1alpha2 --> v1alpha3 version bump by:
#   1. Backfilling required fields on DatabaseUser CRs
#   2. Triggering re-storage of AdvancedCluster CRs
#   3. Removing v1alpha2 from storedVersions on affected CRDs
#
# Run BEFORE upgrading the provider. Requires kubectl + jq.
#
# Dry run (no changes applied):
#   DRY_RUN=true ./scripts/migrate-to-v1.sh
#
# DISCLAIMER: This script is provided "as is", without warranty of any kind,
# express or implied. Use at your own risk. Always test in a non-production
# environment first and ensure you have backups before running this script.
# The authors assume no liability for data loss or cluster disruption.

set -euo pipefail

BACKUP_DIR="/tmp/crossplane-migration-backup"
DRY_RUN="${DRY_RUN:-false}"

AFFECTED_CRDS=(
  "advancedclusters.mongodbatlas.crossplane.io"
  "advancedclusters.mongodbatlas.m.crossplane.io"
  "users.database.mongodbatlas.crossplane.io"
  "users.database.mongodbatlas.m.crossplane.io"
)

GREEN='\033[0;32m' ORANGE='\033[0;33m' RED='\033[0;31m' NC='\033[0m'
log()  { echo -e "${GREEN}[migrate]${NC} $*"; }
warn() { echo -e "${GREEN}[migrate]${NC} ${ORANGE}WARNING:${NC} $*" >&2; }
die()  { echo -e "${GREEN}[migrate]${NC} ${RED}ERROR:${NC} $*" >&2; exit 1; }

command -v kubectl >/dev/null || die "kubectl not found"
command -v jq >/dev/null      || die "jq not found"

mkdir -p "$BACKUP_DIR"

# ---------------------------------------------------------------------------
# Step 1: Patch DatabaseUser CRs, backfill username from external-name
# ---------------------------------------------------------------------------
log "Step 1: Patching DatabaseUser CRs..."

db_users=$(kubectl get users.database.mongodbatlas.crossplane.io -A -o json 2>/dev/null || echo '{"items":[]}')
user_count=$(echo "$db_users" | jq '.items | length')
log "  Found $user_count DatabaseUser CR(s)"

echo "$db_users" | jq -c '.items[]' | while read -r cr; do
  name=$(echo "$cr" | jq -r '.metadata.name')
  ns=$(echo "$cr" | jq -r '.metadata.namespace // "cluster-scoped"')
  ext_name=$(echo "$cr" | jq -r '.metadata.annotations["crossplane.io/external-name"] // empty')
  existing_username=$(echo "$cr" | jq -r '.spec.forProvider.username // empty')
  existing_auth_db=$(echo "$cr" | jq -r '.spec.forProvider.authDatabaseName // empty')

  log "  Processing DatabaseUser $ns/$name (external-name: $ext_name)"

  # Backup
  if [ "$ns" = "cluster-scoped" ]; then
    kubectl get users.database.mongodbatlas.crossplane.io "$name" -o yaml > "$BACKUP_DIR/dbuser-$name.yaml"
  else
    kubectl get users.database.mongodbatlas.crossplane.io "$name" -n "$ns" -o yaml > "$BACKUP_DIR/dbuser-$ns-$name.yaml"
  fi

  # In v0.x the external-name format was: the raw username (via NameAsIdentifier)
  # In v1.0.0 the external-name format is: project_id/username/auth_database_name
  # We need to populate spec.forProvider.username if missing
  if [ -n "$existing_username" ]; then
    log "    username already set ($existing_username), skipping"
    continue
  fi

  if [ -z "$ext_name" ]; then
    warn "    No external-name annotation on $name, cannot infer username, skipping"
    continue
  fi

  # The v0.x external-name was the raw username (NameAsIdentifier)
  username="$ext_name"
  log "    Setting spec.forProvider.username = $username"

  if [ "$DRY_RUN" = "true" ]; then
    log "    [dry-run] would patch username=$username"
  else
    patch='{"spec":{"forProvider":{"username":"'"$username"'"}}}'
    if [ "$ns" = "cluster-scoped" ]; then
      kubectl patch users.database.mongodbatlas.crossplane.io "$name" \
        --type=merge -p "$patch"
    else
      kubectl patch users.database.mongodbatlas.crossplane.io "$name" \
        -n "$ns" --type=merge -p "$patch"
    fi
  fi
done

# ---------------------------------------------------------------------------
# Step 2: Touch AdvancedCluster CRs to trigger re-storage
# ---------------------------------------------------------------------------
log "Step 2: Triggering re-storage of AdvancedCluster CRs..."

adv_clusters=$(kubectl get advancedclusters.mongodbatlas.crossplane.io -A -o json 2>/dev/null || echo '{"items":[]}')
ac_count=$(echo "$adv_clusters" | jq '.items | length')
log "  Found $ac_count AdvancedCluster CR(s)"

echo "$adv_clusters" | jq -c '.items[]' | while read -r cr; do
  name=$(echo "$cr" | jq -r '.metadata.name')
  ns=$(echo "$cr" | jq -r '.metadata.namespace // "cluster-scoped"')

  log "  Touching AdvancedCluster $ns/$name"

  # Backup
  if [ "$ns" = "cluster-scoped" ]; then
    kubectl get advancedclusters.mongodbatlas.crossplane.io "$name" -o yaml > "$BACKUP_DIR/advcluster-$name.yaml"
  else
    kubectl get advancedclusters.mongodbatlas.crossplane.io "$name" -n "$ns" -o yaml > "$BACKUP_DIR/advcluster-$ns-$name.yaml"
  fi

  if [ "$DRY_RUN" = "true" ]; then
    log "    [dry-run] would touch $name"
  else
    # Add a benign annotation to force a write (re-stores at current storage version)
    ts=$(date -u +%Y-%m-%dT%H:%M:%SZ)
    patch='{"metadata":{"annotations":{"migration.crossplane.io/v1alpha3-restorage":"'"$ts"'"}}}'
    if [ "$ns" = "cluster-scoped" ]; then
      kubectl patch advancedclusters.mongodbatlas.crossplane.io "$name" \
        --type=merge -p "$patch"
    else
      kubectl patch advancedclusters.mongodbatlas.crossplane.io "$name" \
        -n "$ns" --type=merge -p "$patch"
    fi
  fi
done

# ---------------------------------------------------------------------------
# Step 3: Remove v1alpha2 from storedVersions on affected CRDs
# ---------------------------------------------------------------------------
log "Step 3: Removing v1alpha2 from storedVersions..."

for crd in "${AFFECTED_CRDS[@]}"; do
  stored=$(kubectl get crd "$crd" -o jsonpath='{.status.storedVersions}' 2>/dev/null || echo "[]")

  if echo "$stored" | grep -q 'v1alpha2'; then
    log "  $crd: replacing storedVersions with [\"v1alpha3\"]"
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
done

log ""
log "Migration complete. Backups saved to $BACKUP_DIR"
log "You can now upgrade provider-mongodbatlas to v1.x."
