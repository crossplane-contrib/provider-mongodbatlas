#!/usr/bin/env bash
# Injects an identifying label into every generated CRD manifest so the
# provider's controller-runtime cache can filter CRDs server-side when
# --filter-crd-cache is enabled (see cmd/provider/main.go).
# Runs as part of `make generate` (generate.done hook).
#
# Plain text insertion (no YAML re-serialization) to preserve
# controller-gen's formatting exactly. Relies on the stable manifest
# structure: a single top-level `metadata:` block containing a
# `  name: <crd-name>` line. Idempotent.
set -euo pipefail

cd "$(dirname "$0")/.."

label_key="app.kubernetes.io/component"
label_value="provider-mongodbatlas"

for f in package/crds/*.yaml; do
	grep -q "^    ${label_key}: ${label_value}$" "$f" && continue
	awk -v key="$label_key" -v value="$label_value" '
		/^  name: / && !done {
			print
			print "  labels:"
			print "    " key ": " value
			done = 1
			next
		}
		{ print }
	' "$f" >"$f.tmp" && mv "$f.tmp" "$f"
done
