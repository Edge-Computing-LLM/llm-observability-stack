# Go-native live validation — 2026-07-18

## Host and root-cause repair

The Xubuntu 24.04.3 host changed from network address `10.165.80.186` to
`10.53.163.158` while k3s was running. k3s retained the former node address,
which caused pod-to-API timeouts, failed health probes, metrics-server failure,
and CrashLoopBackOff in NFD, GPU Operator, kube-state-metrics, and node exporter.
Restarting k3s allowed automatic interface detection and restored the node,
service network, remotedialer, and Metrics API.

## Clean deployment

The Layer 2 release was removed first, followed by Layer 1. Both were then
installed from the local repositories without deleting retained user/model
data:

| Release | Namespace | Revision | Chart | Result |
|---|---|---:|---|---|
| `k3s-nvidia-edge` | `gpu-operator` | 1 | 0.1.0 | deployed |
| `llm-observability-stack` | `llm-observability` | 1 | 0.3.0 / app 1.2.0 | deployed |

All current project pods became Ready. The application/observability pods had
zero restarts. GPU Operator and NFD master each restarted once during the
expected k3s/containerd restart triggered by NVIDIA toolkit reconciliation.

## Verified runtime

- k3s `v1.36.2+k3s1`, node Ready at `10.53.163.158`.
- NVIDIA GeForce 940M, driver `580.95.05`, `nvidia.com/gpu: 1` allocatable.
- `RuntimeClass/nvidia`, device plugin, validator, and DCGM exporter Ready.
- Qwen local alias registered and resident `Forever`.
- Ollama reports 27% CPU / 73% GPU and 23/25 CUDA layers.
- The Go Qwen observer passed 13/13 contract checks.
- Explicit Qwen smoke inference passed in 1.909 seconds.
- Native Go benchmark: TTFT 0.359 seconds, total duration 0.566 seconds, and
  14.54 generated tokens/second for the one-run post-deployment smoke sample.
- `kubectl top node` succeeded after Metrics API recovery.
- The Helm-provisioned `Edge LLM Observability - Ubuntu + k3s + NVIDIA GPU`
  Grafana dashboard was present in the release ConfigMap.

## Validation matrix

- `go test ./...` and `go vet ./...` for all four first-party repositories.
- Go builds for `edge`, `k3s-nvidia-edge`, `llm-observability`,
  `ollama-gateway`, `edge-toolbox`, and the now-renamed `gguf-observe`.
- `helm lint` for both project charts.
- Helm rendering for default, CPU, enterprise, full-NVIDIA, local, and GeForce
  profiles.
- Static container-target builds for the Go gateway and toolbox.
- Native Go Helm/package tests, gateway proxy tests, toolbox tests, benchmark
  tests, and Qwen contract tests.
