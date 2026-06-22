# Pilot Report: Edge LLM Inference and Observability on Kubernetes with NVIDIA GPU

## Executive Summary

`llm-observability-stack` is a pilot-ready Edge LLM inference and observability platform for private/local LLM workloads on Kubernetes and k3s. The project packages Ollama/GGUF model serving, Open WebUI, a LangChain/LangSmith-compatible FastAPI proxy, Prometheus metrics, Grafana dashboards, OpenTelemetry Collector, Kubernetes health checks, NVIDIA GPU scheduling, and DCGM-compatible GPU monitoring into an opinionated Helm-based deployment.

The current validation is a self-deployed technical pilot, not a completed enterprise customer pilot. The repository includes a verified low-cost NVIDIA edge proof on a Lenovo ThinkPad T450s running Xubuntu 24.04 with an NVIDIA GeForce 940M, k3s, NVIDIA device plugin, Ollama `0.17.7`, and a Gemma 3 1B IT Q4_K_M GGUF model. The sanitized benchmark artifact reports TTFT p50 `0.377s`, TTFT p95 `0.381s`, end-to-end p95 `6.972s`, mean generated throughput `11.694 tokens/s`, peak observed GPU utilization `52%`, and loaded-model VRAM of `554 MiB`.

This package is intended to satisfy the NVIDIA Inception customer validation category of "Pilot report or project summary" and support "Proof of deployment" while external customer or partner validation is being pursued.

## Problem Statement

Teams deploying private LLMs on laptops, workstations, edge systems, and small GPU clusters often lack a repeatable way to answer operational questions:

- Is the model actually using the NVIDIA GPU?
- What are TTFT, latency, generated tokens per second, and error rate?
- How much GPU memory, power, temperature, and utilization does the workload consume?
- Are Open WebUI, the inference service, and observability endpoints healthy?
- Can the same local proof be repeated before moving to modern RTX workstations, NVIDIA GPU Operator/DCGM clusters, NIM, or cloud GPU platforms?

Without consistent deployment and observability, private LLM pilots become difficult to benchmark, troubleshoot, support, and justify.

## Target Users/Customers

The intended users are:

- Enterprise AI/platform teams evaluating private LLMs.
- IT teams deploying NVIDIA-powered Linux laptops and workstations.
- Field engineering teams that need local or offline AI support.
- Universities and AI labs with low-cost or mixed GPU fleets.
- OEM/SI partners validating AI workstation bundles.
- Accelerator, cloud, or GPU platform teams evaluating edge-to-cloud LLM operations.

## Customer Use Case

A representative pilot customer wants to deploy a private assistant on an NVIDIA-powered Linux workstation or small edge node. The workload uses a local GGUF model for privacy and offline operation. The customer needs to measure:

- LLM request latency and time to first token.
- Prompt and generated token volume.
- Generated tokens per second.
- Error rate and active requests.
- GPU utilization, framebuffer memory, power, and temperature.
- Kubernetes deployment health and service reachability.

The stack provides a repeatable pilot environment for deploying the model, routing user traffic through an observable proxy, visualizing performance in Grafana, and capturing benchmark evidence.

## Why This Matters for NVIDIA Inception

NVIDIA-powered edge AI deployments need operational evidence, not only model demos. This project aligns with NVIDIA Inception goals by showing a practical path from local GPU inference to measurable GPU-aware LLM operations:

- Kubernetes-native deployment on NVIDIA-enabled Linux hosts.
- `RuntimeClass: nvidia` and `nvidia.com/gpu` scheduling.
- NVIDIA device plugin and GPU Operator/DCGM-ready deployment options.
- DCGM-compatible dashboards and Prometheus alert rules.
- A planned NVIDIA NIM monitoring path through `/v1/metrics`.
- A roadmap from low-cost GeForce proof to RTX laptop, workstation, GPU Operator/DCGM, NIM, and cloud GPU validation.

The project is seeking external design-partner validation before making production or enterprise traction claims.

## Solution Overview

`llm-observability-stack` is an umbrella Helm chart that composes:

- Ollama for local GGUF-backed inference.
- Open WebUI for browser-based chat.
- `langchain-demo`, a FastAPI proxy and metrics bridge.
- Prometheus and Grafana through kube-prometheus-stack.
- OpenTelemetry Collector for OTLP traces, metrics, and logs.
- NVIDIA device plugin, GPU Operator, and DCGM exporter chart paths.
- Prometheus `ServiceMonitor`, `Probe`, and `PrometheusRule` resources.
- Grafana dashboards for LLM, benchmark, and NVIDIA GPU observability.
- Python toolbox and notebooks for in-cluster diagnostics and demonstrations.
- Benchmark scripts and sanitized artifacts for repeatable evidence capture.

## Architecture

Primary runtime path:

```text
Browser / benchmark client
        |
        v
Open WebUI or direct API client
        |
        v
langchain-demo FastAPI proxy
        |\
        | +--> Prometheus /metrics
        | +--> LangSmith traces when configured
        | +--> OpenTelemetry Collector when enabled
        v
Ollama + private GGUF model
        |
        v
NVIDIA GPU runtime / device plugin

Prometheus + Grafana + Alertmanager
        ^
        +-- ServiceMonitors, blackbox probes, benchmark metrics,
            Kubernetes metrics, and DCGM-compatible GPU metrics
```

The chart also contains optional monitoring hooks for NVIDIA NIM services exposing OpenMetrics on `/v1/metrics`.

## Deployment Environment

Current verified local environment from repository evidence:

- Host: Lenovo ThinkPad T450s on Xubuntu 24.04.
- Kubernetes: single-node k3s, combined control-plane and worker.
- GPU: NVIDIA GeForce 940M, 1 GiB VRAM, CUDA compute capability 5.0.
- Driver: `580.95.05`.
- NVIDIA device plugin: `0.18.1`.
- GPU resource: `nvidia.com/gpu: 1`.
- RuntimeClass: `nvidia`.
- Ollama: `0.17.7`.
- Model: Gemma 3 1B IT Q4_K_M GGUF, approximately 806 MB.
- Profile: `values.geforce-940m-k3s.yaml` for narrow inference proof and `values.enterprise-pilot-k3s.yaml` for the fuller local pilot stack.

## Current Project Status

Implemented:

- Helm umbrella chart and vendored dependency charts.
- Local k3s/NVIDIA deployment profiles.
- Ollama/GGUF hostPath model workflow.
- Open WebUI routing through the observable proxy.
- FastAPI proxy with LLM and HTTP Prometheus metrics.
- OpenTelemetry Collector subchart configuration.
- kube-prometheus-stack integration.
- Grafana dashboards for LLM overview, benchmark results, and NVIDIA GPU metrics.
- Prometheus alert rules for LLM error rate, p95 latency, endpoint probe failure, GPU temperature, and framebuffer memory.
- NVIDIA device plugin, GPU Operator, and DCGM exporter chart options.
- DCGM and NIM ServiceMonitor hooks.
- Benchmark script and sanitized GeForce 940M benchmark JSON.
- Jupyter notebooks and Python toolbox for diagnostics and evidence workflows.

Not yet complete:

- Named enterprise customer validation.
- Redacted customer email or approved reference quote.
- Formal SI/NCP/OEM/accelerator validation.
- Multi-device RTX/workstation/cloud benchmark matrix.
- Completed deployed NVIDIA NIM proof.
- Full production security review, SSO/RBAC, backup, retention, and support runbook evidence.

## LLM Runtime Path: Ollama / llama.cpp / GGUF

The current runtime path uses Ollama to serve a local GGUF model. The chart renders a ConfigMap-backed Modelfile from `Modelfile.gemma-3-1b-it-gguf`, mounts the local GGUF directory read-only into the Ollama pod, and creates the model at startup when configured.

The verified profiles use:

- Model name: `gemma3-1b-it-gguf-local`.
- GGUF file: `gemma-3-1b-it-Q4_K_M.gguf` in the verified local profile.
- Ollama image: `ollama/ollama:0.17.7`.
- GPU scheduling through `runtimeClassName: nvidia` and `nvidia.com/gpu`.

The project treats Ollama/GGUF as the portable local inference proof. NVIDIA NIM is a planned scale path, not the current completed inference runtime.

## NVIDIA GPU Observability Path

The repository includes several NVIDIA-aware deployment paths:

- `gpu-operator` chart dependency for clusters that need NVIDIA driver/toolkit/device plugin/DCGM lifecycle management.
- `nvidia-device-plugin` chart dependency for lightweight k3s/workstation deployments where host drivers and container runtime are already installed.
- `dcgm-exporter` chart dependency for DCGM-compatible metrics.
- `monitoring.dcgmExporter.serviceMonitor` for external DCGM exporter service discovery.
- NVIDIA GPU Grafana dashboard using `DCGM_FI_*` metrics.
- Prometheus alert rules for high GPU temperature and high framebuffer memory usage.
- GPU scheduling configuration through `nvidia.runtimeClassName`, `nvidia.gpuResourceName`, and `nvidia.gpuCount`.

## Prometheus/Grafana/OpenTelemetry Observability Stack

The observability stack includes:

- kube-prometheus-stack for Prometheus Operator, Prometheus, Grafana, Alertmanager, node exporter, and kube-state-metrics.
- `ServiceMonitor` for the `langchain-demo` `/metrics` endpoint.
- Blackbox `Probe` resources for HTTP endpoint health.
- `PrometheusRule` resources for LLM and NVIDIA GPU alerts.
- Grafana dashboard ConfigMap loading dashboards from `dashboards/*.json`.
- OpenTelemetry Collector subchart with OTLP gRPC and HTTP receivers.
- Optional operator-managed `OpenTelemetryCollector` custom resource path.

## Metrics Collected or Intended

Implemented LLM/application metrics:

- HTTP request count by method, route, and status.
- HTTP request duration.
- LLM request count by model, route, and outcome.
- Active LLM requests.
- LLM request duration.
- Time to first token.
- Average inter-token latency derived from Ollama response metadata.
- Prompt token totals.
- Generated token totals.
- Generated tokens per second.

Implemented benchmark metrics:

- Benchmark success.
- TTFT p50/p95.
- Mean tokens per second.
- End-to-end duration p95.
- Pushgateway-compatible metric publishing from `benchmarks/ollama_benchmark.py`.

GPU/DCGM metrics expected when DCGM exporter is available:

- GPU utilization.
- Framebuffer memory usage/free.
- GPU temperature.
- Power usage.
- SM and Tensor Core activity where supported.
- SM clock.

## Pilot Objectives

The pilot objectives are:

1. Deploy a private LLM stack on k3s with NVIDIA GPU scheduling.
2. Confirm local GGUF inference through Ollama.
3. Route browser/API requests through an observable proxy.
4. Capture LLM latency, TTFT, tokens, throughput, active requests, and errors.
5. Capture GPU utilization, memory, power, temperature, and related DCGM metrics when available.
6. Display LLM, benchmark, and NVIDIA GPU metrics in Grafana.
7. Produce sanitized benchmark and deployment evidence suitable for accelerator review.

## Pilot Workflow

1. Prepare the k3s node and verify NVIDIA runtime/device plugin.
2. Confirm the local GGUF model exists on host storage.
3. Build and import local `langchain-demo` and `python-toolbox` images into k3s containerd.
4. Install Prometheus Operator CRDs and the Helm chart using `values.enterprise-pilot-k3s.yaml`.
5. Verify pods, services, PVCs, RuntimeClass, GPU resource, and model availability.
6. Port-forward Ollama, Open WebUI, Grafana, and the proxy as needed.
7. Run inference tests and benchmark runs.
8. Capture Grafana screenshots, terminal outputs, logs, benchmark JSON, and GPU evidence.
9. Store sanitized artifacts under `validation/screenshots/` and `validation/benchmark-results/`.

## Proof of Deployment Checklist

- Git commit SHA captured.
- `helm lint .` passed.
- Helm template render captured for the active profile.
- `kubectl get nodes -o wide` captured.
- `kubectl get runtimeclass nvidia` captured.
- `kubectl get nodes` GPU allocatable output captured.
- `kubectl get pods,svc,pvc -n llm-observability -o wide` captured.
- Ollama model list and inference output captured.
- `langchain-demo` health and metrics endpoints captured.
- Prometheus ServiceMonitors, Probes, and PrometheusRules captured.
- Grafana dashboard screenshots captured with timestamp and visible units.
- Benchmark JSON captured and stored.
- GPU evidence from `nvidia-smi`, DCGM metrics, or both captured.

## Current Validation Status

Current status: self-deployed technical pilot and proof of deployment are available. External customer validation remains open.

Available repository evidence:

- `README.md` NVIDIA Inception positioning and verified GeForce 940M proof.
- `docs/competition/VERIFIED-LOCAL-RESULTS.md`.
- `artifacts/geforce-940m-benchmark.json`.
- `values.geforce-940m-k3s.yaml`.
- `values.enterprise-pilot-k3s.yaml`.
- `hack/test-geforce-940m-inference.sh`.
- `hack/bootstrap-enterprise-pilot-k3s.sh`.
- `dashboards/llm-overview.json`.
- `dashboards/nvidia-gpu.json`.
- `dashboards/benchmark-results.json`.
- `templates/prometheus-rules.yaml`.
- `templates/nvidia-servicemonitors.yaml`.

## Known Limitations

- No confirmed paying customer is claimed.
- No named enterprise design partner is claimed.
- No approved partner/OEM/SI/NCP/accelerator validation is claimed.
- Current verified GPU proof is a constrained GeForce 940M profile, not a modern RTX or data-center GPU benchmark.
- The current benchmark uses low concurrency and a small model.
- NIM is planned/optional and has not been presented as a completed deployment proof.
- Security hardening is partial and requires environment-specific review before production use.
- Production readiness, support SLAs, compliance, SSO/RBAC, backup, retention, and multi-user policy are not yet claimed.

## Planned NVIDIA NIM Integration

The repository includes an optional NIM monitoring path through `monitoring.nim.serviceMonitor`, intended for services that expose OpenMetrics on `/v1/metrics`. The planned NIM integration is:

- Deploy NIM on a supported NVIDIA GPU environment.
- Configure the ServiceMonitor selector and port for the NIM service.
- Compare Ollama/GGUF and NIM observability contracts.
- Add NIM latency, request, and GPU correlation evidence to Grafana/Prometheus dashboards.
- Produce a separate NIM deployment proof with hardware, driver, CUDA, GPU Operator/DCGM, and model details.

This is a roadmap item until a dated NIM deployment and benchmark artifact are captured.

## Roadmap Toward Enterprise Pilot

1. Capture a modern RTX laptop benchmark with 1B, 3B, and 7B GGUF models.
2. Capture a workstation benchmark with sustained load and concurrency.
3. Validate GPU Operator/DCGM on a supported NVIDIA cluster.
4. Complete an optional NIM deployment proof and metrics comparison.
5. Recruit a design partner with a private/local LLM use case.
6. Agree success criteria, data boundaries, and allowed evidence.
7. Produce a redacted customer pilot report or approved quote.
8. Conduct security and data-flow review.

## Success Criteria

Technical pilot success criteria:

- Stack deploys through Helm without manual manifest rewriting.
- NVIDIA GPU is visible to Kubernetes as `nvidia.com/gpu`.
- Ollama model loads and completes inference.
- Open WebUI can route requests through the observable proxy.
- Prometheus scrapes LLM proxy metrics.
- Grafana dashboards load and show LLM/benchmark/GPU panels.
- Benchmark script produces JSON with TTFT, latency, and throughput.
- Evidence can be captured without exposing secrets or private model data.

Enterprise pilot success criteria:

- Named or responsibly anonymized customer workload.
- Approved data and privacy boundary.
- Measurable before/after operational result.
- Redacted confirmation of value from customer or partner.
- Reproducible deployment and benchmark evidence on target NVIDIA hardware.

## Conclusion

`llm-observability-stack` is a credible, pilot-ready EdgeLLM observability platform with repository-backed proof of local NVIDIA GPU inference and monitoring. It is strongest today as a self-deployed technical pilot and proof-of-deployment package. The next milestone is external validation from a design partner, accelerator, cloud provider, OEM/SI, or enterprise reviewer, followed by modern NVIDIA GPU benchmarks and optional NIM deployment evidence.
