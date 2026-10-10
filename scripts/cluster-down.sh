#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="kubetaskops-dev"

echo "=== Tearing down Kind Cluster '${CLUSTER_NAME}' ==="

if kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
  kind delete cluster --name "${CLUSTER_NAME}"
  echo "Cluster '${CLUSTER_NAME}' deleted successfully."
else
  echo "Cluster '${CLUSTER_NAME}' does not exist."
fi
