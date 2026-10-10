#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

echo "=== Running Kubernetes Code & Manifest Generators ==="
cd "${BACKEND_DIR}"

echo "1. Generating DeepCopy methods..."
go run sigs.k8s.io/controller-tools/cmd/controller-gen@v0.16.5 \
  object:headerFile="hack/boilerplate.go.txt" \
  paths="./pkg/apis/v1alpha1/..."

echo "2. Generating CRD manifests into config/crd/bases/..."
go run sigs.k8s.io/controller-tools/cmd/controller-gen@v0.16.5 \
  crd:crdVersions=v1 \
  paths="./pkg/apis/v1alpha1/..." \
  output:crd:artifacts:config="./config/crd/bases"

echo "=== Code & Manifest Generation Complete ==="
