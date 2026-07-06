# k3s-nvidia-edge Dependency

`llm-observability-stack` is the LLM application and observability layer. For NVIDIA GPU deployments on local k3s, it depends on `k3s-nvidia-edge` to prepare the cluster substrate first.

Deploy in this order:

1. `k3s-nvidia-edge`
2. `llm-observability-stack`

The two projects intentionally own different layers.

## Ownership Boundary

`k3s-nvidia-edge` owns the local k3s and NVIDIA GPU layer:

- k3s server and k3s-managed containerd
- CoreDNS and local-path-provisioner through k3s
- NVIDIA Container Toolkit configuration for k3s containerd
- NVIDIA GPU Operator in namespace `gpu-operator`
- NVIDIA device plugin through GPU Operator
- NVIDIA DCGM exporter through GPU Operator
- Node Feature Discovery through GPU Operator
- `RuntimeClass/nvidia`
- Kubernetes node labels such as `nvidia.com/gpu.present=true`
- allocatable resource `nvidia.com/gpu`
- CUDA validation pod workflow

`llm-observability-stack` owns the LLM workload layer:

- Ollama and GGUF model loading
- Open WebUI browser UI
- Open WebUI Redis helper when enabled by the Open WebUI subchart
- OpenTelemetry Collector and GenAI telemetry paths
- optional Prometheus/Grafana dashboards and alerting resources
- optional FastAPI proxy, Python toolbox, benchmarks, notebooks, and diagnostics

Do not enable GPU Operator, standalone NVIDIA device plugin, or standalone DCGM exporter in `llm-observability-stack` when `k3s-nvidia-edge` is already deployed. That would duplicate the GPU substrate that the base layer already owns.

## Required Ready State

Before installing the NVIDIA profiles in this repository, the cluster should show:

```bash
helm list -n gpu-operator
kubectl get pods -n gpu-operator
kubectl get runtimeclass nvidia
kubectl get nodes --show-labels | grep nvidia.com/gpu.present=true
kubectl describe node "$(kubectl get nodes -o jsonpath='{.items[0].metadata.name}')" | grep -A8 Allocatable
```

Expected signals:

- Helm release `k3s-nvidia-edge` or `gpu-operator` is deployed in namespace `gpu-operator`.
- GPU Operator pods are `Running` or validator pods are `Completed`.
- `RuntimeClass/nvidia` exists.
- the node has `nvidia.com/gpu.present=true`.
- allocatable resources include `nvidia.com/gpu: 1` or greater.

The local GeForce profile also expects a readable GGUF model directory on the node at the path configured by `ollamaModel.gguf.hostPath` and the matching Ollama volume host path.

## Install Order

From the `k3s-nvidia-edge` repository:

```bash
cd /media/waqasm86/External1/Waqas-Projects/Project-Edge-Computing-LLM/k3s-nvidia-edge
make check
bin/k3s-nvidia-edge install --yes --sudo=false --use-local-chart --skip-base-package-install --skip-toolkit-install --skip-k3s-install
bin/k3s-nvidia-edge validate --yes
```

For a fresh host, use the full installer instead:

```bash
bin/k3s-nvidia-edge install --yes
```

Then install this repository:

```bash
cd /media/waqasm86/External1/Waqas-Projects/Project-Edge-Computing-LLM/llm-observability-stack
helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.geforce-940m-k3s.yaml \
  --wait
```

## Current GeForce Profile Behavior

`values.geforce-940m-k3s.yaml` is designed to run on top of `k3s-nvidia-edge`.

It enables:

- Ollama with `runtimeClassName: nvidia`
- `nvidia.com/gpu: 1` GPU request and limit through the Ollama subchart
- node selection on `nvidia.com/gpu.present=true`
- Open WebUI pointed directly at `http://ollama:11434`
- Open WebUI Redis from the Open WebUI subchart
- OpenTelemetry Collector

It keeps disabled:

- root-level `gpu-operator`
- root-level `nvidia-device-plugin`
- root-level `dcgm-exporter`
- root-level Redis
- LangChain demo
- Python toolbox
- etcd simulation

This keeps the GPU substrate in `k3s-nvidia-edge` and the LLM application layer in `llm-observability-stack`.

## Failure Modes

If `llm-observability-stack` is installed before `k3s-nvidia-edge` is ready, GPU profiles can fail with:

- pending Ollama pods because `RuntimeClass/nvidia` does not exist
- pending Ollama pods because `nvidia.com/gpu` is not allocatable
- scheduling failures because `nvidia.com/gpu.present=true` is missing
- model startup failures if the GGUF host path is missing

Fix the base layer first, then rerun the same Helm command for `llm-observability-stack`.
