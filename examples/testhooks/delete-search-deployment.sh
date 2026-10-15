#!/usr/bin/env bash
# The search Deployment must be gone before its AdvancedCluster is deleted:
# once the cluster is gone, Atlas answers every search deployment call with
# 400 ATLAS_SEARCH_DEPLOYMENT_INVALID_ATTRIBUTE and the Deployment never
# finalizes. kubectl delete waits until the finalizer is removed.
set -euo pipefail

${KUBECTL} delete deployments.search.mongodbatlas.crossplane.io deployment-example --ignore-not-found
${KUBECTL} delete deployments.search.mongodbatlas.m.crossplane.io deployment-example -n crossplane-system --ignore-not-found
