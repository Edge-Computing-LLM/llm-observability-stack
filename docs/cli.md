# llm-observability CLI

`llm-observability` is the repo-local Go helper CLI. Use `edge-cli` as the
primary organization-level CLI for end-to-end installs and uninstalls.

## Architecture Boundary

- `k3s-nvidia-edge` owns the local Linux, k3s, k3s containerd, NVIDIA runtime, GPU Operator, NVIDIA device plugin, DCGM exporter, Node Feature Discovery, `RuntimeClass/nvidia`, and `nvidia.com/gpu` validation layer.
- `llm-observability-stack` owns Ollama, Open WebUI, Open WebUI Redis, OpenTelemetry Collector, optional Prometheus/Grafana, optional LangChain/FastAPI proxy, benchmark tooling, notebooks, and local model configuration.
- `edge-cli` owns cross-repository ordering: infra first, observability second.
- The CLI does not import `k3s-nvidia-edge/internal/...`.

## Build

From this repository:

```bash
go build -o bin/llm-observability ./cmd/llm-observability
```

During local sibling-repo development, `go.mod` uses:

```text
replace github.com/Edge-Computing-LLM/k3s-nvidia-edge => ../k3s-nvidia-edge
```

That keeps the long-term import path clean while allowing local changes in `k3s-nvidia-edge/pkg/edgebase` to be tested immediately. When `k3s-nvidia-edge` publishes a version tag containing `pkg/edgebase`, this can be changed to a normal tagged requirement.

## Commands

```bash
bin/llm-observability doctor
bin/llm-observability install --profile geforce-940m-k3s --skip-base --yes
bin/llm-observability status
bin/llm-observability validate
bin/llm-observability benchmark --model gemma3-1b-it-gguf-local --runs 3
bin/llm-observability uninstall --yes
bin/llm-observability print-commands --profile geforce-940m-k3s
```

## Profiles

The CLI maps profile names to existing values files:

| Profile | Values file |
|---|---|
| `geforce-940m-k3s` | `values.geforce-940m-k3s.yaml` |
| `enterprise-pilot-k3s` | `values.enterprise-pilot-k3s.yaml` |
| `cpu-k3s` | `values.cpu-k3s.yaml` |
| `local-k3s-example` | `values.local-k3s.example.yaml` |
| `local-k3s` | `values.local-k3s.yaml` |
| `full-stack-nvidia` | `values.full-stack-nvidia.example.yaml` |
| `validation-k3s` | `values.validation-k3s.yaml` |

You can also pass a values file directly:

```bash
bin/llm-observability install --profile values.geforce-940m-k3s.yaml --yes
```

Additional overrides are supported:

```bash
bin/llm-observability install \
  --profile geforce-940m-k3s \
  --values values.local-k3s.yaml \
  --set ollamaModel.gguf.hostPath=/path/to/models \
  --yes
```

For NVIDIA profiles, the CLI adds Helm safeguards so this chart does not redeploy the base GPU layer:

```text
--set gpu-operator.enabled=false
--set nvidia-device-plugin.enabled=false
--set dcgm-exporter.enabled=false
```

## Recommended Local Command

For the current Xubuntu 24 + k3s + NVIDIA GPU setup where `k3s-nvidia-edge` is already healthy:

```bash
bin/llm-observability install --profile geforce-940m-k3s --skip-base --yes
bin/llm-observability validate
```

If the base layer is not installed yet, use `edge-cli`:

```bash
edge install infra --yes
edge validate infra
edge install observability --profile geforce-940m-k3s --yes
```

By default, install and uninstall commands are dry-run unless `--yes` is provided.

## Validation

`validate` checks:

- base NVIDIA readiness for GPU profiles
- Helm release status
- Ollama service, deployment, model list, and loaded model state
- Open WebUI statefulset and service
- Open WebUI Redis deployment when present
- OpenTelemetry Collector deployment and service
- pod readiness in the namespace
- CUDA/offload evidence in Ollama logs for GPU profiles

The base CUDA validation pod from `edgebase` is dry-run unless `--yes` is provided, matching the safety behavior of `k3s-nvidia-edge`.

## Benchmark

The benchmark command wraps the existing Python client:

```bash
bin/llm-observability benchmark \
  --model gemma3-1b-it-gguf-local \
  --runs 3 \
  --prompt "Explain edge GPU observability in one sentence." \
  --output artifacts/benchmark-local.json
```

It starts a temporary `kubectl port-forward` to `svc/ollama` and then runs `benchmarks/ollama_benchmark.py`.

## Uninstall

Uninstall only the LLM stack:

```bash
bin/llm-observability uninstall --yes
```

Keep the namespace:

```bash
bin/llm-observability uninstall --keep-namespace --yes
```

Remove all layers through `edge-cli` when you need reverse-order cleanup:

```bash
edge uninstall all --yes
```

The repo-local `--with-base` flag is deprecated and returns an error.
