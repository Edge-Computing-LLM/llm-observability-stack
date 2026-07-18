#!/usr/bin/env bash
set -euo pipefail

RUNBOOK_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
REPO_ROOT="$(cd "${RUNBOOK_DIR}/../../.." && pwd -P)"
K8S_NAMESPACE="${K8S_NAMESPACE:-llm-observability}"
K8S_RELEASE="${K8S_RELEASE:-llm-observability-stack}"
OLLAMA_GATEWAY_SERVICE="${OLLAMA_GATEWAY_SERVICE:-ollama-gateway}"

if [[ -f "${RUNBOOK_DIR}/../config.env" ]]; then
  # shellcheck disable=SC1091
  source "${RUNBOOK_DIR}/../config.env"
fi

cd "${REPO_ROOT}"
go run ./cmd/llm-observability status \
  --namespace "${K8S_NAMESPACE}" --release "${K8S_RELEASE}"
go run ./cmd/llm-observability validate \
  --namespace "${K8S_NAMESPACE}" --release "${K8S_RELEASE}"
go run ./cmd/llm-observability network --namespace "${K8S_NAMESPACE}"
for service in open-webui ollama "${OLLAMA_GATEWAY_SERVICE}"; do
  if ! kubectl get service -n "${K8S_NAMESPACE}" "${service}" >/dev/null 2>&1; then
    printf 'Skipping optional Service %s: not installed\n' "${service}"
    continue
  fi
  go run ./cmd/llm-observability service-path \
    --namespace "${K8S_NAMESPACE}" --service "${service}"
done
