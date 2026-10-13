#!/usr/bin/env bash
# Loads the Atlas sample datasets into every AdvancedCluster labeled
# testing.upbound.io/load-sample-data=true, so resources that need existing
# collections (OnlineArchive, search Index) can be created. A loaded cluster
# gets the testing.upbound.io/sample-data-loaded annotation and is skipped by
# later runs of this hook in the same batch. Needs the Atlas CLI (atlas) and
# jq; credentials come from MONGODB_ATLAS_* or UPTEST_CLOUD_CREDENTIALS.
set -euo pipefail

export_from_credentials() {
  local value
  value=$(jq -r --arg key "$2" '.[$key] // empty' <<<"${UPTEST_CLOUD_CREDENTIALS}")
  if [ -n "${value}" ]; then
    export "$1=${value}"
  fi
}

if [ -z "${MONGODB_ATLAS_CLIENT_ID:-}${MONGODB_ATLAS_PUBLIC_API_KEY:-}" ]; then
  export_from_credentials MONGODB_ATLAS_CLIENT_ID client_id
  export_from_credentials MONGODB_ATLAS_CLIENT_SECRET client_secret
  export_from_credentials MONGODB_ATLAS_PUBLIC_API_KEY public_key
  export_from_credentials MONGODB_ATLAS_PRIVATE_API_KEY private_key
fi

found=0
for resource in advancedclusters.mongodbatlas.crossplane.io advancedclusters.mongodbatlas.m.crossplane.io; do
  while read -r -u 3 namespace name; do
    [ -n "${name}" ] || continue
    found=1
    ref=("${resource}/${name}")
    if [ "${namespace}" != "<none>" ]; then
      ref+=(-n "${namespace}")
    fi
    if [ -n "$(${KUBECTL} get "${ref[@]}" -o jsonpath='{.metadata.annotations.testing\.upbound\.io/sample-data-loaded}')" ]; then
      echo "Sample data already loaded into ${ref[*]}"
      continue
    fi
    echo "Waiting for ${ref[*]} to become Ready..."
    ${KUBECTL} wait "${ref[@]}" --for=condition=Ready --timeout=30m
    cluster=$(${KUBECTL} get "${ref[@]}" -o jsonpath='{.spec.forProvider.name}')
    project=$(${KUBECTL} get "${ref[@]}" -o jsonpath='{.spec.forProvider.projectId}')
    echo "Loading sample data into Atlas cluster ${cluster}..."
    job=$(atlas clusters sampleData load "${cluster}" --projectId "${project}" -o json | jq -r '._id // empty')
    if [ -z "${job}" ]; then
      echo "Could not start the sample data load for Atlas cluster ${cluster}" >&2
      exit 1
    fi
    atlas clusters sampleData watch "${job}" --projectId "${project}"
    ${KUBECTL} annotate "${ref[@]}" testing.upbound.io/sample-data-loaded=true --overwrite
  done 3< <(${KUBECTL} get "${resource}" --all-namespaces -l testing.upbound.io/load-sample-data=true \
    -o custom-columns=NAMESPACE:.metadata.namespace,NAME:.metadata.name --no-headers 2>/dev/null || true)
done

if [ "${found}" -eq 0 ]; then
  echo "No AdvancedCluster labeled testing.upbound.io/load-sample-data=true found" >&2
  exit 1
fi
