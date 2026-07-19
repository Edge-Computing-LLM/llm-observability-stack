# llm-observability-stack

Kubernetes-native observability, benchmarking, and operations tooling for private LLM inference on local edge systems.

Preferred organization CLI: [`edge-cli`](https://github.com/Edge-Computing-LLM/edge-cli).
Repo-local legacy/helper CLI documentation: [docs/cli.md](docs/cli.md).

This repository packages a Helm-based application and observability stack for k3s and Kubernetes with Ollama/GGUF model serving, Open WebUI, a native Go OpenTelemetry GenAI gateway, Prometheus, Grafana, OpenTelemetry Collector, blackbox probes, benchmark metrics, and NVIDIA/DCGM-compatible dashboards.

The repository also includes a Go CLI named `llm-observability` for repo-local helper workflows. New end-to-end installs should use `edge-cli`, which deploys `k3s-nvidia-edge` first and then this chart.

GitHub repository: <https://github.com/Edge-Computing-LLM/llm-observability-stack>

## Required Base Layer For Local NVIDIA k3s

For local NVIDIA GPU deployments, deploy [`k3s-nvidia-edge`](https://github.com/Edge-Computing-LLM/k3s-nvidia-edge) first. This repository expects the GPU substrate to already exist before GPU profiles such as `values.geforce-940m-k3s.yaml` are installed.

`k3s-nvidia-edge` owns k3s, k3s containerd NVIDIA runtime wiring, GPU Operator, NVIDIA device plugin, DCGM exporter, Node Feature Discovery, `RuntimeClass/nvidia`, and the allocatable `nvidia.com/gpu` resource. `llm-observability-stack` then deploys Ollama, Open WebUI, OpenTelemetry, dashboards, benchmarks, and application-level observability on top of that base layer.

An empty local k3s cluster with only CoreDNS, local-path-provisioner, and other default k3s system components is a valid starting point before the base layer is deployed. Run `edge install infra` or validate `k3s-nvidia-edge` before installing GPU profiles from this repository.

Read the full dependency guide before installing GPU profiles:

- [k3s-nvidia-edge dependency](docs/K3S-NVIDIA-EDGE-DEPENDENCY.md)

## What This Stack Provides

- Local private LLM serving through Ollama and legally obtained GGUF models.
- Kubernetes deployment through Helm with k3s-friendly profiles.
- NVIDIA GPU scheduling with `runtimeClassName: nvidia` and `nvidia.com/gpu` when a GPU is available.
- Optional CPU validation profiles for development clusters without NVIDIA GPUs.
- Open WebUI for browser-based interaction with local models.
- A native Go Ollama gateway with streaming, OpenTelemetry traces, and LLM request metrics for TTFT, latency, active requests, and errors.
- Prometheus, Grafana, Alertmanager, kube-state-metrics, node exporter, ServiceMonitors, probes, and alert rules.
- OpenTelemetry Collector endpoints for OTLP traces, metrics, and logs.
- Optional native Go diagnostics including DNS/HTTP/TCP checks, Redis checks, OpenTelemetry seeding, and benchmark reporting.

## Verified Local NVIDIA GPU Deployment

The current local deployment target is a single-node Xubuntu 24 system running k3s with an NVIDIA GPU. The verified low-memory edge profile has been tested on:

- Host: ThinkPad T450s on Xubuntu 24.
- GPU: NVIDIA GeForce 940M, 1 GiB VRAM, CUDA compute capability 5.0.
- k3s node: combined control-plane and worker.
- NVIDIA device plugin resource: `nvidia.com/gpu: 1`.
- RuntimeClass: `nvidia`.
- Model profiles: Qwen 1.8B Chat Q4_K_M, Gemma 3 1B IT Q4_K_M, and
  Ollama `llama3.2:1b`, deployed one at a time with partial CUDA offload.

Measured after deployment, warmup, and exact-response, arithmetic, and
translation inference checks:

| Metric | Result |
|---|---:|
| Model size | Approximately 1.2 GB |
| CUDA layers | 23/25 |
| Processor split | 27% CPU / 73% GPU |
| Context / batch | 256 / 1 |
| Observed throughput | 9.75-15.78 tokens/s |
| VRAM usage | 824 MiB used / 152 MiB free |
| Residency | `Forever` |

Evidence and reproduction:

- [Single-node k3s GeForce 940M guide](docs/SINGLE-NODE-K3S-GEFORCE-940M.md)
- [Local k3s NVIDIA deployment report - 2026-07-02](docs/LOCAL-K3S-NVIDIA-REPORT-2026-07-02.md)
- [Live layered validation - 2026-07-08](docs/LIVE-VALIDATION-2026-07-08.md)
- [Verified local GPU results](docs/VERIFIED-LOCAL-GPU-RESULTS.md)
- [Xubuntu k3s NVIDIA runbook](docs/XUBUNTU-K3S-NVIDIA-RUNBOOK.md)
- [Sanitized benchmark artifact](artifacts/geforce-940m-benchmark.json)
- [GeForce 940M Helm profile](values.geforce-940m-k3s.yaml)
- [Multi-model GGUF runtime evidence companion](https://github.com/Edge-Computing-LLM/gguf-observability)
- [Sequential three-model deployment and validation](docs/MULTI-MODEL-LOW-VRAM.md)

These numbers prove constrained local edge feasibility. They do not claim enterprise load, concurrency, fleet reliability, or production readiness.

The companion repository performs read-only runtime contract checks and
captures sanitized evidence. This chart remains the source of truth for the
Modelfile, Helm values, model lifecycle, and workload configuration.

The former standalone TypeScript/Vue dashboard has been migrated into the
chart-owned, Helm-provisioned Grafana dashboard
[`dashboards/edge-llm-observability.json`](dashboards/edge-llm-observability.json).
The Grafana version uses Prometheus, DCGM, and kube-state-metrics directly and
keeps the complete dashboard definition reproducible with the Helm release.

## Who This Is For

- Developers running private LLMs on local Linux systems.
- Platform teams evaluating local LLM observability on k3s or Kubernetes.
- IT and field engineering teams that need repeatable offline/private AI deployments.
- Labs using low-cost CPU and GPU systems for model-serving experiments.
- Operators who need a local-first path from CPU-only testing to NVIDIA GPU acceleration.

## What This Is Not

- Not a generic cloud-only LLM observability SaaS.
- Not a replacement for OpenTelemetry, Grafana, Prometheus, DCGM, or NIM.
- Not a claim that every laptop GPU is suitable for production LLM inference.
- Not a repository for committing GGUF model binaries, kubeconfigs, credentials, or secrets.

## Platform Components

- Vendored Helm charts for Ollama, Open WebUI, kube-prometheus-stack, OpenTelemetry Collector, and OpenTelemetry Operator.
- Native Go OpenTelemetry GenAI-instrumented Ollama gateway with Prometheus metrics.
- TTFT, latency, token, throughput, active-request, HTTP, and error telemetry.
- Optional kube-prometheus-stack, Grafana, Alertmanager, node exporter, and kube-state-metrics from the root umbrella chart.
- OpenTelemetry Collector endpoint for OTLP traces, metrics, and logs, with an optional operator-managed collector path.
- Blackbox endpoint probes and Prometheus alert rules.
- NVIDIA DCGM dashboard and external DCGM ServiceMonitor integration.
- A comprehensive Edge LLM dashboard for live GPU metrics, workload
  readiness, service inventory, the validated Qwen profile, and telemetry
  readiness.
- NVIDIA NIM `/v1/metrics` ServiceMonitor path for environments that use NIM.
- Pushgateway-compatible benchmark reporting.
- Optional Go edge toolbox, Redis, OpenTelemetry seeder, and etcd failure simulation.

## Runtime Architecture

```text
User or benchmark client
        |
        v
Open WebUI / Go Ollama gateway
        |                \
        |                 +--> OpenTelemetry GenAI traces
        |                 +--> Prometheus /metrics
        v
Ollama + private GGUF model       Optional NVIDIA NIM
        |                              |
        +---------- NVIDIA GPU --------+
                         |
                  DCGM / GPU metrics

Prometheus + Grafana + Alertmanager
        ^
        +-- ServiceMonitors, probes, benchmark Pushgateway, Kubernetes metrics
```

The verified laptop profile uses Ollama/GGUF. The same observability contract can be used on larger local RTX workstations with the NVIDIA substrate prepared by `k3s-nvidia-edge`.

## Repository Layout

```text
llm-observability-stack/
├── Chart.yaml
├── values.yaml
├── values.validation-k3s.yaml
├── values.geforce-940m-k3s.yaml
├── values.enterprise-pilot-k3s.yaml
├── values.full-stack-nvidia.example.yaml
├── values.cpu-k3s.yaml
├── values.local-k3s.example.yaml
├── artifacts/                     # sanitized public benchmark evidence
├── cmd/                            # Go CLI, gateway, and toolbox entrypoints
├── dashboards/                    # LLM, benchmark, and NVIDIA GPU dashboards
├── internal/                       # CLI, gateway, toolbox, benchmark packages
├── templates/                     # application monitoring and security manifests
├── charts/                        # vendored dependency charts
├── ollama-gateway/                # native Go gateway image definition
├── edge-toolbox/                  # native Go in-cluster diagnostics image
├── docs/                          # architecture, operations, and local runbooks
├── hack/                          # validation, device-plugin, and evidence scripts
└── tests/                         # Helm and application smoke tests
```

Build the CLI:

```bash
go build -o bin/llm-observability ./cmd/llm-observability
```

Preferred local CLI flow from the organization control plane:

```bash
edge install all --accelerator auto --yes
edge validate observability
```

Repo-local helper flow when `k3s-nvidia-edge` is already healthy:

```bash
bin/llm-observability doctor
bin/llm-observability install --profile geforce-940m-k3s --skip-base --yes
bin/llm-observability validate
```

## Prerequisites

- Linux host or cluster with k3s/Kubernetes reachable through `kubectl`.
- Helm 3 or 4.
- For local NVIDIA k3s GPU profiles: `k3s-nvidia-edge` deployed and validated first.
- NVIDIA driver and NVIDIA Container Toolkit for GPU profiles.
- `RuntimeClass/nvidia` and `nvidia.com/gpu` provided by `k3s-nvidia-edge` for GPU mode.
- A legally obtained GGUF model available on node storage.
- Go 1.25 or newer for CLI, gateway, toolbox, benchmark, and tests.

Quick checks:

```bash
kubectl get nodes -o wide
helm list -n gpu-operator
kubectl get pods -n gpu-operator
kubectl get runtimeclass nvidia
kubectl get nodes -o jsonpath='{range .items[*]}{.metadata.name}{" gpu="}{.status.allocatable.nvidia\.com/gpu}{"\n"}{end}'
helm version
```

The local bootstrap helper detects the Kubernetes runtime before installing. It uses NVIDIA mode when Kubernetes advertises `nvidia.com/gpu`; otherwise it writes a CPU-only overlay and runs the same edge LLM observability path without NVIDIA runtime or GPU resource requests.

The organization CLI exposes the same policy directly. The full installer uses
host detection, while `edge install observability --accelerator auto --yes` uses
Kubernetes allocatable GPU capacity. Explicit `cpu` and `nvidia` modes are also
available for deterministic automation.

## Quick Start

### A. Minimal validation profile

```bash
helm template llm-observability-stack . \
  -f values.validation-k3s.yaml
```

### B. Verified GeForce 940M edge profile

Review the machine-specific model host path before using this profile on another system. The profile schedules on nodes with `nvidia.com/gpu.present=true`, which supports a single-node k3s control-plane/worker laptop without requiring a separate worker label.

This profile uses the locally retained Qwen 1.8B Chat Q4_K_M GGUF and
`Modelfile.qwen-1.8b-chat-q4_K_M`. On the 1 GiB GeForce 940M it pins 23/25
layers to CUDA, limits batch size to 1, uses a 256-token context, and keeps the
model loaded indefinitely. The measured steady-state allocation is 824 MiB VRAM.

Preferred: deploy and validate the base layer through `edge-cli` first:

```bash
edge install infra --yes
edge validate infra
```

Then deploy the LLM stack:

```bash
cd /media/waqasm86/External1/Waqas-Projects/Project-Linux-Kubernetes-Nvidia/Project-Edge-Computing-LLM/llm-observability-stack

helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.geforce-940m-k3s.yaml

./hack/test-geforce-940m-inference.sh
```

Select Gemma or Llama by adding the matching overlay after the base GeForce
values file:

```bash
# Gemma 3 1B IT Q4_K_M
helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.geforce-940m-k3s.yaml \
  -f values.gemma-3-1b-geforce-940m-k3s.yaml

# Llama 3.2 1B
helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.geforce-940m-k3s.yaml \
  -f values.llama3.2-1b-geforce-940m-k3s.yaml
```

Use only one profile at a time on the 1 GiB GPU. The Modelfile `num_gpu`
setting controls partial layer offload; read-only runtime evidence enforces a
900 MiB total observed VRAM ceiling.

### C. Full observability NVIDIA profile

```bash
helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.full-stack-nvidia.example.yaml
```

This installs the LLM and observability layer only. It does not install GPU
Operator, NVIDIA device plugin, or DCGM exporter. Use private values files or
existing Kubernetes Secrets for OpenTelemetry and Open WebUI secrets. Never
commit secrets.

### D. Local full-stack k3s profile

This profile is tailored for the verified local single-node k3s/NVIDIA GPU workstation. It uses the vendored OpenTelemetry Collector subchart, keeps external-facing services as `ClusterIP`, and keeps the existing Ollama `local-path` PVC at `5Gi`.

```bash
helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.enterprise-pilot-k3s.yaml \
  --set kube-prometheus-stack.crds.enabled=false
```

Import the local `ollama-gateway` and `edge-toolbox` images into k3s containerd before enabling those two workloads.

For a guided local setup, use:

```bash
./hack/bootstrap-enterprise-pilot-k3s.sh
```

To inspect the generated runtime overlay without installing:

```bash
./hack/detect-runtime-profile.sh
cat .generated/values.runtime-detected.yaml
```

To force CPU mode for validation:

```bash
./hack/detect-runtime-profile.sh --mode cpu
helm template llm-observability-stack . \
  -f values.enterprise-pilot-k3s.yaml \
  -f .generated/values.runtime-detected.yaml \
  --set kube-prometheus-stack.crds.enabled=false
```

Do not switch an existing release from `values.enterprise-pilot-k3s.yaml` to a private profile that changes the `ollama` PVC size unless you intentionally recreate or migrate the PVC. k3s `local-path` storage does not resize that claim in place.

## Access and Benchmarking

```bash
kubectl get pods -n llm-observability -o wide
kubectl port-forward -n llm-observability svc/ollama 11434:11434
kubectl port-forward -n llm-observability \
  svc/llm-observability-stack-grafana 3000:80
```

For the GeForce 940M profile, open <http://127.0.0.1:3000> and select
**Edge LLM Observability - Ubuntu + k3s + NVIDIA GPU**. See
[`dashboards/README.md`](dashboards/README.md) for provisioning, credentials,
and dashboard-as-code guidance.

Run the native Go benchmark (it manages a temporary port-forward):

```bash
bin/llm-observability benchmark \
  --model qwen-1-8b-chat-q4-k-m-local \
  --runs 10 \
  --output artifacts/benchmark-local.json
```

Only sanitized evidence intended for publication should be committed.

## Validation

```bash
helm lint .
helm template llm-observability-stack . >/tmp/rendered-default.yaml
helm template llm-observability-stack . \
  -f values.geforce-940m-k3s.yaml >/tmp/rendered-geforce.yaml
helm template llm-observability-stack . \
  -f values.full-stack-nvidia.example.yaml \
  --set opentelemetry.tracing.enabled= \
  --set openWebUI.existingSecret= \
  --set open-webui.webuiSecret.existingSecretName= \
  >/tmp/rendered-full-stack-nvidia.yaml

go test ./...
go vet ./...
./hack/validate-local-stack.sh
./hack/validate-local-stack.sh --strict-gpu
```

The strict GPU check requires an active cluster with an allocatable NVIDIA GPU.

## Local Runbooks

- [Xubuntu k3s NVIDIA runbook](docs/XUBUNTU-K3S-NVIDIA-RUNBOOK.md)
- [Local k3s NVIDIA runbook](docs/LOCAL-K3S-NVIDIA-RUNBOOK.md)
- [Operations runbook](docs/OPERATIONS-RUNBOOK.md)
- [Configuration profiles](docs/CONFIG-PROFILES.md)
- [k3s-nvidia-edge dependency](docs/K3S-NVIDIA-EDGE-DEPENDENCY.md)
- [GitHub publishing guide](docs/GITHUB-PUBLISHING.md)

## Security and Evidence Boundaries

- Use `existingSecret` references or private ignored values files.
- Keep prompt and response capture disabled or redacted for confidential workloads.
- Do not commit model binaries, kubeconfigs, private evidence, credentials, or TLS keys.
- Treat host-path model mounts and `local-path` persistence as local edge-reference defaults, not universal production storage.
- Complete TLS, SSO/RBAC, backup, retention, network-policy, and threat-model review before production use.

## Troubleshooting

```bash
kubectl get pods -A -o wide
kubectl describe pod -n llm-observability -l app.kubernetes.io/name=ollama
kubectl logs -n llm-observability deployment/ollama --tail=200
kubectl get pods -n gpu-operator
kubectl get nodes -o json | jq '.items[].status.allocatable'
watch -n 0.5 nvidia-smi
```

The first Ollama image pull can be several gigabytes and may exceed a short Helm wait timeout. Once cached, rerun the same `helm upgrade --install` command to reconcile the release.

## Documentation

Start with [docs/README.md](docs/README.md), then use:

- [Architecture](docs/ARCHITECTURE.md)
- [Configuration profiles](docs/CONFIG-PROFILES.md)
- [k3s-nvidia-edge dependency](docs/K3S-NVIDIA-EDGE-DEPENDENCY.md)
- [Quickstart](docs/QUICKSTART.md)
- [Operations runbook](docs/OPERATIONS-RUNBOOK.md)
- [Xubuntu k3s NVIDIA runbook](docs/XUBUNTU-K3S-NVIDIA-RUNBOOK.md)
- [Complete project documentation](docs/PROJECT-DOCUMENTATION.md)

## Project Status

`llm-observability-stack` is an open-source local LLM observability reference implementation with verified single-node k3s/NVIDIA evidence and CPU-only deployment support. The next hardening areas are modern RTX benchmarking, multi-node testing, security review, backup/restore validation, and production-specific access control.
