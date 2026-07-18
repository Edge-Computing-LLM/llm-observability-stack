# llm-observability in action (k3s)

This folder provides a practical Kubernetes operations toolkit for your local `llm-observability-stack`, aligned to the topics covered across `Kubernetes in Action, 2nd Edition` code and book sections (pods, config, storage, services, ingress/gateway, controllers, and batch workloads).

## Directory layout

- `config.env.example` - base runtime variables for your local k3s namespace and services.
- `lib/common.sh` - shared shell helpers used by all kubectl scripts.
- `kubectl/` - executable shell scripts grouped by Kubernetes domain.
- `manifests/networking/` - optional networking manifests (NetworkPolicy and test pod).
- `runbooks/` - wrappers to run the read-only kubectl and native Go suites.

## Quick start

```bash
PROJECT_ROOT="${PROJECT_ROOT:-$(git rev-parse --show-toplevel 2>/dev/null || pwd)}"
cd "${PROJECT_ROOT}/jupyter-notebooks/llm-observability-in-action"
cp config.env.example config.env

# Optional: set K8S_CONTEXT in config.env if multiple contexts exist.

./runbooks/run_kubectl_suite_readonly.sh
./runbooks/run_go_suite.sh
```

## Mutation safety model

All shell scripts default to read-only behavior.

- Read-only mode: `APPLY_CHANGES=0` (default)
- Mutating mode: `APPLY_CHANGES=1`

Example (explicit mutation enablement):

```bash
APPLY_CHANGES=1 ./kubectl/04_configmaps_and_secrets.sh
```

## Script index (kubectl)

1. `01_cluster_and_api.sh` - cluster, context, API resources, explain flow.
2. `02_namespaces_and_workloads.sh` - workloads across Deployments/RS/STS/DS/Jobs/CronJobs.
3. `03_pods_lifecycle_debug.sh` - pod lifecycle, logs, events, debug entry points.
4. `04_configmaps_and_secrets.sh` - config and secret inventory with safe secret-key inspection.
5. `05_storage_and_volumes.sh` - StorageClass/PV/PVC/volume mounts and storage diagnostics.
6. `06_networking_core.sh` - Services/Endpoints/EndpointSlices/Ingress/NetworkPolicy.
7. `07_networking_advanced.sh` - Gateway API, DNS/CNI controls, node-network context.
8. `08_security_policy.sh` - ServiceAccounts/RBAC/auth checks/pod security labels.
9. `09_jobs_and_batch.sh` - Jobs/CronJobs execution and diagnostics.
10. `10_llm_observability_stack_checks.sh` - stack-specific checks for `open-webui`, `ollama`, `ollama-gateway`.

## Native Go operations

`runbooks/run_go_suite.sh` uses the repository CLI for release validation,
namespace inventory, and Service-to-Pod/Endpoint tracing. No Python dependency
is required for these operations.

## Networking manifests

- `manifests/networking/netpol.default-deny.yaml`
- `manifests/networking/netpol.allow-dns.yaml`
- `manifests/networking/netpol.allow-openwebui-to-ollama-gateway.yaml`
- `manifests/networking/test-client-pod.yaml`

Apply only when intentionally testing networking behavior:

```bash
APPLY_CHANGES=1 kubectl -n llm-observability apply -f manifests/networking/
```

## Kubernetes in Action alignment

See `BOOK_COVERAGE.md` for mapping from chapter domains into this toolkit.
