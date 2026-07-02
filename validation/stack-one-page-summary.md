# llm-observability-stack One-Page Summary

`llm-observability-stack` is a Helm-based local LLM observability stack for k3s and Kubernetes. It packages Ollama/GGUF model serving, Open WebUI, a LangChain-compatible FastAPI proxy, Prometheus metrics, Grafana dashboards, OpenTelemetry Collector, endpoint probes, benchmark tooling, and optional NVIDIA GPU monitoring.

## Problem

Private/local LLM deployments are hard to operate when inference latency, time to first token, throughput, GPU utilization, endpoint health, and Kubernetes deployment state are spread across disconnected tools.

## Solution

The stack provides one reproducible deployment path for local inference and observability:

- Ollama serves local GGUF models.
- Open WebUI provides browser access.
- `langchain-demo` exposes a FastAPI proxy with LLM metrics.
- Prometheus and Grafana collect and visualize application, cluster, benchmark, and GPU metrics.
- OpenTelemetry Collector accepts OTLP telemetry.
- CPU-only and NVIDIA GPU profiles support development, local validation, and edge deployment.

## Verified Local State

- Xubuntu 24 single-node k3s environment.
- NVIDIA RuntimeClass `nvidia`.
- `nvidia.com/gpu: 1` advertised by the node.
- GPU Operator deployed separately in namespace `gpu-operator`.
- `llm-observability-stack` deployed in namespace `llm-observability`.
- Full local stack pods Running, including Ollama, Open WebUI, LangChain proxy, OpenTelemetry Collector, Prometheus, Grafana, Alertmanager, DCGM exporter, and Python toolbox.

## Key References

- `README.md`
- `docs/XUBUNTU-K3S-NVIDIA-RUNBOOK.md`
- `docs/LOCAL-K3S-NVIDIA-REPORT-2026-07-02.md`
- `docs/VERIFIED-LOCAL-GPU-RESULTS.md`
- `values.geforce-940m-k3s.yaml`
- `values.cpu-k3s.yaml`
- `values.full-stack-nvidia.example.yaml`
- `hack/validate-local-stack.sh`
- `hack/capture-local-evidence.sh`

## Current Limits

- Local evidence is not a production-readiness claim.
- 1 GiB GPUs are suitable only for very small local models.
- Production deployments still need security review, TLS/SSO/RBAC, backup/restore validation, data-retention policy, and workload-specific capacity testing.
