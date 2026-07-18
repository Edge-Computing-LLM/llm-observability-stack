# Local Validation Run Report: 2026-06-22

This report summarizes a local end-to-end validation run for `llm-observability-stack` on the Xubuntu 24/k3s/NVIDIA GPU workstation.

## Run Metadata

| Item | Value |
|---|---|
| Captured at | 2026-06-22 17:01 UTC |
| Repository commit tested | `7a4181ff54d31a4ec4977278c805055aad7f96b6` |
| Host | `waqasm86-thinkpad-t450s` |
| OS/kernel | Ubuntu/Xubuntu 24.04.3 LTS, Linux `6.17.0-35-generic` |
| Kubernetes | k3s `v1.35.5+k3s1` |
| Helm | `v4.1.4` |
| NVIDIA GPU | NVIDIA GeForce 940M, 1024 MiB VRAM |
| NVIDIA driver | `580.95.05` |
| RuntimeClass | `nvidia` present |
| Kubernetes GPU resource | `nvidia.com/gpu=1` allocatable |
| Helm profile | `values.enterprise-pilot-k3s.yaml` |
| Helm release | `llm-observability-stack`, namespace `llm-observability`, revision `2` |

## Static Validation

| Check | Result |
|---|---|
| `helm dependency build .` | Passed |
| `helm lint .` | Passed |
| `helm template` default profile | Passed |
| `helm template -f values.local-k3s.example.yaml` | Passed |
| `helm template -f values.geforce-940m-k3s.yaml` | Passed |
| `helm template -f values.full-stack-nvidia.example.yaml` with empty secret overrides | Passed |
| `python3 -m pytest -q tests` | Passed, 12 tests |
| `./hack/validate-local-stack.sh --strict-gpu` | Passed |

## Deployment Result

The full local k3s/NVIDIA profile was installed/upgraded successfully:

```bash
helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.enterprise-pilot-k3s.yaml \
  --set kube-prometheus-stack.crds.enabled=false \
  --wait --timeout 10m
```

All primary workloads reached Running/Ready state:

- Ollama
- Open WebUI
- Open WebUI Redis
- LangChain demo proxy
- Python toolbox
- OpenTelemetry Collector
- DCGM exporter
- kube-prometheus-stack operator
- Prometheus
- Alertmanager
- Grafana
- kube-state-metrics
- Prometheus node exporter

## LLM Runtime Validation

Ollama created and served the expected local GGUF-backed model:

| Item | Result |
|---|---|
| Model | `gemma3-1b-it-gguf-local:latest` |
| Format | GGUF |
| Family | Gemma 3 |
| Quantization | `Q4_K_M` |
| Model size reported by Ollama | 806,059,328 bytes |
| Direct in-cluster Ollama smoke test | Passed |
| Direct `ollama run` inference | Passed |
| LangChain proxy `/healthz` | HTTP 200 |
| LangChain proxy `/config` | HTTP 200 |
| LangChain proxy `/metrics` | HTTP 200 |
| LangChain proxy `/invoke` | HTTP 200 |

Observed proxy response:

```text
The proxy path is currently healthy and functioning as expected.
```

## NVIDIA GPU Validation

Kubernetes and Ollama both detected the GPU path:

- `nvidia-smi` reported NVIDIA GeForce 940M with driver `580.95.05`.
- Kubernetes reported `nvidia.com/gpu=1` allocatable.
- Ollama logs reported CUDA inference compute on `NVIDIA GeForce 940M`, compute capability `5.0`.
- Ollama loaded CUDA backend from `/usr/lib/ollama/cuda_v12/libggml-cuda.so`.
- `nvidia-smi` showed 554 MiB GPU memory in use after model load.

DCGM exporter was deployed and served metrics through `http://dcgm-exporter:9400/metrics`, including `DCGM_FI_DEV_*` metrics labelled with the `ollama` pod.

## Observability Validation

The following monitoring resources were present:

- `ServiceMonitor/dcgm-exporter`
- `ServiceMonitor/ollama-gateway`
- `ServiceMonitor/opentelemetry-collector`
- kube-prometheus-stack ServiceMonitors
- `Probe/llm-stack-http`
- `PrometheusRule/llm-observability-stack`

OpenTelemetry Collector started successfully and reported:

```text
Everything is ready. Begin running and processing data.
```

Open WebUI initially refused HTTP while downloading first-start Hugging Face assets. After the download completed, `http://open-webui:8080/` returned HTTP 200.

## Benchmark Summary

Benchmark command:

```bash
bin/llm-observability benchmark \
  --model gemma3-1b-it-gguf-local \
  --warmup-runs 1 \
  --runs 3 \
  --output validation/benchmark-results/benchmark-20260622T170052Z.json
```

The raw JSON output is intentionally ignored by Git under `validation/benchmark-results/`. Sanitized summary:

| Metric | Result |
|---|---:|
| Warmup runs | 1 |
| Measured runs | 3 |
| TTFT p50 | 0.591 s |
| TTFT p95 | 2.148 s |
| End-to-end p50 | 7.986 s |
| End-to-end p95 | 8.563 s |
| Mean generated throughput | 10.668 tokens/s |

## Notes and Caveats

- This is a local self-deployed validation run, not an external customer deployment.
- Open WebUI first startup can take several minutes because it downloads embedding assets from Hugging Face.
- Raw benchmark JSON, logs, screenshots, chart archives, caches, and local image artifacts remain ignored by Git unless explicitly sanitized and force-added.
- The deployment uses the local GeForce 940M edge profile and should not be represented as modern RTX, NIM, or production-scale evidence.
