# Local k3s NVIDIA Deployment Report - 2026-07-02

## Scope

This report records the local validation of `llm-observability-stack` on the Xubuntu 24 single-node k3s workstation with NVIDIA GPU support. It also captures the comparison points from `k3s-nvidia-edge` and the local reference repositories under `/media/waqasm86/External1/Waqas-Projects`.

## Host And Cluster

- Host OS reported by Kubernetes: Ubuntu 24.04.3 LTS / Xubuntu desktop environment.
- Kernel: `6.17.0-35-generic`.
- k3s: `v1.36.2+k3s1`.
- Kubernetes node: `waqasm86-thinkpad-t450s`, `Ready`, `control-plane`.
- Container runtime: `containerd://2.3.2-k3s2`.
- NVIDIA driver: `580.95.05`.
- CUDA reported by `nvidia-smi`: `13.0`.
- GPU: NVIDIA GeForce 940M, 1 GiB VRAM.
- Kubernetes GPU resource: `nvidia.com/gpu`, capacity `1`, allocatable `1`.
- RuntimeClass: `nvidia` exists and is managed by the GPU Operator.

## Existing k3s NVIDIA State

The existing `k3s-nvidia-edge` setup is healthy and matches its documented local profile:

- Helm release `gpu-operator` is deployed in namespace `gpu-operator`.
- GPU Operator chart/app: `gpu-operator-v26.3.3` / `v26.3.3`.
- Effective GPU Operator values:
  - `driver.enabled=false`
  - `toolkit.enabled=true`
  - `gfd.enabled=false`
  - k3s containerd paths point at `/var/lib/rancher/k3s/agent/etc/containerd/config.toml` and `/run/k3s/containerd/containerd.sock`
- GPU Operator pods, NVIDIA Container Toolkit, device plugin, operator validator, and DCGM exporter are running.
- `go test ./...` passed in `k3s-nvidia-edge`.

## Reference Repositories Reviewed

The prompt path `/media/waqasm86/External1/Project-Llamatelemetry/Project-Llamatelemetry-End-to-End` was not present. The matching local reference repos used by `k3s-nvidia-edge` were found under `/media/waqasm86/External1/Waqas-Projects`:

| Area | Repository | Branch | Commit |
|---|---|---:|---:|
| Kubernetes SIGs | `node-feature-discovery` | `master` | `06aa41e40` |
| Kubernetes SIGs | `dra-driver-nvidia-gpu` | `main` | `779a7dd0` |
| CoreDNS | `coredns` | `master` | `4faf983fb` |
| Rancher/k3s | `local-path-provisioner` | `master` | `5d4bfc84` |
| Rancher/k3s | `k3s` | `main` | `5faa00674a` |
| NVIDIA | `gpu-operator` | `main` | `a37981bdf` |
| NVIDIA | `cuda-samples` | `master` | `b7c5481c` |
| NVIDIA | `DCGM` | `master` | `d646460` |
| NVIDIA | `go-dcgm` | `main` | `0740c4c` |
| NVIDIA | `libnvidia-container` | `main` | `0d1d7494` |
| NVIDIA | `nvidia-container-toolkit` | `main` | `41dd4444` |
| NVIDIA | `k8s-device-plugin` | `main` | `25e493580` |
| NVIDIA | `dcgm-exporter` | `main` | `d5e5f51` |

Observed alignment:

- `llm-observability-stack` now keeps GPU Operator, NVIDIA device plugin, and DCGM exporter out of its Helm dependency list. Those remain owned by `k3s-nvidia-edge`.
- The local deployment avoids installing duplicate NVIDIA device-plugin/operator components when the cluster already has GPU Operator.
- The local profile now follows the single-node k3s reality by using GPU capability detection instead of requiring a `node-role.kubernetes.io/worker=true` selector.

## Project Changes Applied

- Updated the local GGUF model host path from the stale `External11` path to:
  `/media/waqasm86/External1/Waqas-Projects/repos-llamatelemetry/llamatelemetry-xubuntu24/models`
- Confirmed model file:
  `gemma-3-1b-it-Q4_K_M.gguf`, size `769M`.
- Removed static worker-node selection from NVIDIA profiles so the single control-plane node can schedule Ollama.
- Kept the runtime detector overlay behavior:
  - NVIDIA mode sets `runtimeClassName: nvidia` and requests `nvidia.com/gpu: 1`.
  - CPU mode clears NVIDIA runtime, node selector, DCGM, and GPU requests.
- Disabled Open WebUI first-start external model downloads in local profiles by setting:
  - `RAG_EMBEDDING_ENGINE=ollama`
  - `RAG_EMBEDDING_MODEL=gemma3-1b-it-gguf-local`
  - `RAG_EMBEDDING_MODEL_AUTO_UPDATE=False`
  - `RAG_RERANKING_MODEL_AUTO_UPDATE=False`
  - `WHISPER_MODEL_AUTO_UPDATE=False`
  - `ENABLE_VERSION_UPDATE_CHECK=False`
- Updated repository metadata/docs to point at:
  `https://github.com/Edge-Computing-LLM/llm-observability-stack`
- Updated `.gitignore` for local run outputs and Kubernetes report artifacts.

The Open WebUI environment names are documented in the official Open WebUI environment configuration reference:
`https://docs.openwebui.com/getting-started/env-configuration/`.

## Deployment

Runtime detection:

```text
Detected nvidia runtime profile
GPU nodes: waqasm86-thinkpad-t450s:1
```

Install and upgrade commands used:

```bash
./hack/bootstrap-enterprise-pilot-k3s.sh

helm upgrade llm-observability-stack . \
  -n llm-observability \
  -f values.enterprise-pilot-k3s.yaml \
  -f .generated/values.runtime-detected.yaml \
  --set namespace.create=false \
  --set kube-prometheus-stack.crds.enabled=false
```

Final Helm releases:

| Release | Namespace | Status | Chart | App |
|---|---|---|---|---|
| `gpu-operator` | `gpu-operator` | deployed | `gpu-operator-v26.3.3` | `v26.3.3` |
| `llm-observability-stack` | `llm-observability` | deployed | `llm-observability-stack-0.2.0` | `1.1.0` |

## Final Workload State

All expected pods were Running:

- `alertmanager-kube-prometheus-stack-alertmanager-0` `2/2`
- `dcgm-exporter` `1/1`
- `kube-prometheus-stack-operator` `1/1`
- `ollama-gateway` `1/1`
- `llm-observability-stack-grafana` `3/3`
- `kube-state-metrics` `1/1`
- `prometheus-node-exporter` `1/1`
- `ollama` `1/1`
- `open-webui-0` `1/1`
- `open-webui-redis` `1/1`
- `opentelemetry-collector` `1/1`
- `prometheus-kube-prometheus-stack-prometheus-0` `2/2`
- `edge-toolbox` `1/1`

## Test Results

Static and chart tests:

| Test | Result |
|---|---:|
| `go test ./...` in `k3s-nvidia-edge` | pass |
| `helm lint` enterprise NVIDIA profile | pass |
| `helm lint` CPU profile | pass |
| `helm template` enterprise NVIDIA profile | pass, 8922 rendered lines |
| `helm template` CPU profile | pass, 8500 rendered lines |
| `pytest -q` | pass, 15 tests |

Live cluster tests:

| Test | Result |
|---|---:|
| CPU pod `busybox:1.38` | pass, `cpu-smoke-ok`, `x86_64` |
| API-server CUDA pod dry run | pass |
| Standalone CUDA pod while Ollama was deployed | expected Pending due to `Insufficient nvidia.com/gpu`; the single GPU was already allocated to Ollama |
| `nvidia-smi` inside Ollama pod | pass |
| `ollama list` inside Ollama pod | pass, `gemma3-1b-it-gguf-local:latest` |
| Ollama `/api/tags` from in-cluster toolbox | HTTP 200 |
| LangChain demo `/healthz` | HTTP 200 |
| Open WebUI `/` | HTTP 200 after disabling startup downloads |
| DCGM exporter `/metrics` | HTTP 200 |
| OpenTelemetry collector `/metrics` | HTTP 200 |
| Grafana `/api/health` | HTTP 200 |
| Prometheus `/-/ready` | HTTP 200 |
| Ollama chat smoke | pass, 1.27s final run |

GPU evidence from inside the Ollama pod:

```text
NVIDIA-SMI 580.95.05
Driver Version: 580.95.05
CUDA Version: 13.0
GPU: NVIDIA GeForce 940M
Memory: 552MiB / 1024MiB during loaded model state
```

Ollama model evidence:

```text
gemma3-1b-it-gguf-local:latest    806 MB
```

## Observations

- The live k3s + NVIDIA setup is healthy.
- The project now runs on this single-node k3s control-plane laptop without requiring a separate worker label.
- The GPU profile and CPU profile both render cleanly and carry the same application path.
- The GPU cannot run a second `nvidia.com/gpu: 1` validation pod while Ollama owns the only GPU. This is correct Kubernetes resource accounting, not a failure.
- First-time image pulls were slow for Open WebUI, Grafana, Prometheus, and Alertmanager. After image caching, redeploys were much faster.
- Open WebUI 0.8.10 tried to download Hugging Face assets on first start. The local profile now disables those automatic downloads and uses the local Ollama model path for embedding configuration.
- The chart still uses `ClusterIP` services, which is appropriate for this local edge profile. Use `kubectl port-forward` for browser access.

## Access Commands

```bash
kubectl port-forward -n llm-observability svc/open-webui 8080:8080
kubectl port-forward -n llm-observability svc/llm-observability-stack-grafana 3000:80
kubectl port-forward -n llm-observability svc/kube-prometheus-stack-prometheus 9090:9090
kubectl port-forward -n llm-observability svc/ollama 11434:11434
```

## Current Status

As of this report, `llm-observability-stack` is deployed successfully in k3s with NVIDIA GPU acceleration, and the CPU fallback profile is validated by render/lint plus a live CPU-only Kubernetes smoke pod.
