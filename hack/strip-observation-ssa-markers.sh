#!/usr/bin/env bash
set -euo pipefail

# Strip SSA merge-strategy markers from *Observation structs in generated types.
#
# upjet adds +listType=map, +listMapKey, and +kubebuilder:default markers to
# every struct variant (Init, Parameters, Observation). In Observation structs
# (written to status.atProvider) these markers cause the API server to default
# every injected "index" field to "0", which violates the map-key uniqueness
# constraint when a list has more than one item.
#
# The markers are only needed on Init/Parameters (spec-side, for SSA merges).
# Observation lists are written atomically by the controller, so they can stay
# the default listType (atomic).

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

find "${ROOT}/apis" -name 'zz_*_types.go' -print0 | while IFS= read -r -d '' file; do
  sed -i '/Observation struct {/,/^}/{
    /\/\/ +listType=map$/d
    /\/\/ +listMapKey=/d
    /\/\/ +kubebuilder:default:="0"$/d
  }' "$file"
done
