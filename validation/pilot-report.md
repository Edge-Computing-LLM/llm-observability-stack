# Local Validation Report

This report summarizes the current local technical validation of `llm-observability-stack`. The filename is retained for compatibility with older links, but the content now describes local stack validation only.

## Summary

`llm-observability-stack` deploys local LLM inference and observability components on k3s/Kubernetes. The verified local environment is Xubuntu 24 with k3s, NVIDIA GPU scheduling, Ollama, Open WebUI, a native Go Ollama gateway, OpenTelemetry Collector, Prometheus, Grafana, Alertmanager, DCGM exporter, and diagnostic workloads.

## Verified Local Environment

| Area | Value |
|---|---|
| OS | Xubuntu/Ubuntu 24.04 series |
| Kubernetes | k3s `v1.36.2+k3s1` |
| Runtime | containerd through k3s |
| GPU path | NVIDIA RuntimeClass `nvidia`, node advertises `nvidia.com/gpu: 1` |
| Namespace | `llm-observability` |
| Release | `llm-observability-stack` |
| Main profiles | `values.geforce-940m-k3s.yaml`, `values.enterprise-pilot-k3s.yaml`, `values.cpu-k3s.yaml`, `values.full-stack-nvidia.example.yaml` |

## What Was Validated

- Helm chart linting and template rendering.
- GPU profile rendering with NVIDIA runtime and GPU resource requests.
- CPU-only profile rendering without NVIDIA runtime or GPU resource requests.
- Local k3s deployment with full observability services.
- Ollama model-serving path with local GGUF model mount.
- Open WebUI startup with local auto-download behavior disabled.
- Prometheus/Grafana monitoring resources and DCGM-compatible GPU observability path.

## Known Limits

- The 1 GiB GeForce 940M profile is useful for small-model local validation only.
- Local validation is not a production readiness claim.
- Production use still requires TLS, SSO/RBAC, secrets management, backup/restore validation, retention policy, network policies, and workload-specific capacity testing.

## Primary References

- `README.md`
- `docs/XUBUNTU-K3S-NVIDIA-RUNBOOK.md`
- `docs/LOCAL-K3S-NVIDIA-REPORT-2026-07-02.md`
- `docs/VERIFIED-LOCAL-GPU-RESULTS.md`
- `validation/deployment-proof.md`
- `hack/validate-local-stack.sh`
- `hack/capture-local-evidence.sh`
