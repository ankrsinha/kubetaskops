#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="kubetaskops-dev"

echo "=== KubeTaskOps Dev Cluster Bootstrap ==="

if ! command -v kind &> /dev/null; then
  echo "Error: 'kind' CLI is not installed."
  exit 1
fi

if ! command -v kubectl &> /dev/null; then
  echo "Error: 'kubectl' CLI is not installed."
  exit 1
fi

if kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
  echo "Kind cluster '${CLUSTER_NAME}' is already running."
else
  echo "Creating Kind cluster '${CLUSTER_NAME}'..."
  kind create cluster --name "${CLUSTER_NAME}"
fi

kubectl config use-context "kind-${CLUSTER_NAME}"

echo "Cluster status:"
kubectl cluster-info --context "kind-${CLUSTER_NAME}"

echo "=== Kind Cluster '${CLUSTER_NAME}' Ready ==="
