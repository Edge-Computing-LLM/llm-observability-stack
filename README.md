# llm-observability-stack

Umbrella Helm chart for a local single-node stack:

- `k3s`
- NVIDIA GPU runtime + device plugin (installed separately)
- `ollama`
- `open-webui`
- `langchain` demo API with `langsmith` tracing

This repository is tuned for local Xubuntu/k3s workflows and GGUF models mounted from host storage.

GitHub repository: https://github.com/waqasm86/llm-observability-stack

## What this deploys

- Vendored `ollama` chart (`charts/ollama`)
- Vendored `open-webui` chart (`charts/open-webui`)
- Local Modelfile ConfigMap for GGUF-backed Ollama model creation
- LangChain demo API deployment/service (observability test endpoint)
- Optional secret creation for LangSmith and Open WebUI

## Architecture

1. Ollama runs as a Kubernetes deployment with NVIDIA runtime class.
2. A hostPath GGUF directory is mounted into Ollama (`/models/gguf`).
3. A templated Modelfile builds a local model (default: `gemma3-1b-it-gguf-local`).
4. Open WebUI connects to Ollama via in-cluster service DNS (`http://ollama:11434`).
5. LangChain demo API invokes Ollama and emits traces to LangSmith.

## Prerequisites

- Kubernetes: `k3s` (single node is supported)
- NVIDIA stack: container toolkit + working GPU runtime + NVIDIA device plugin in cluster
- `helm` and `kubectl` installed on host
- Local GGUF model file available on host filesystem

## Quick start

Use the example values file as your starting point:

```bash
cp values.local-k3s.example.yaml values.local-k3s.yaml
```

Edit `values.local-k3s.yaml` and set:

- `ollamaModel.gguf.hostPath` to your GGUF parent directory
- `ollama.volumes[0].hostPath.path` to the same directory
- LangSmith credentials (or existing secret reference)
- Open WebUI secret key (32+ chars)

Then install/upgrade:

```bash
helm dependency build .
helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.local-k3s.yaml
```

## Localhost access

For local browser/API access, set service type to `LoadBalancer` (already set in local k3s values):

- `ollama.service.type: LoadBalancer`
- `open-webui.service.type: LoadBalancer`

Endpoints:

- Open WebUI: `http://localhost:8080/`
- Ollama API: `http://localhost:11434/`
- LangChain demo (cluster service): `http://langchain-demo.llm-observability.svc.cluster.local:8000`

## Model workflow

1. Put your GGUF file on host.
2. Mount the GGUF directory to `/models/gguf` in Ollama.
3. Use `__GGUF_PATH__` in Modelfile `FROM` line (chart replaces it at render time).
4. Ollama creates the model at startup from ConfigMap-backed Modelfile.

## Repository hygiene

This repository intentionally excludes local secrets and generated files via `.gitignore`, including:

- `values.local-k3s.yaml`
- `.webui_secret_key`
- `rendered.yaml`
- model binaries like `*.gguf`

Use only `values.local-k3s.example.yaml` in git.

## Size policy

Target repository size: `< 100 MB`.

Current chart source footprint is small (well under this threshold). To verify before pushing:

```bash
du -sh .
find . -type f -printf '%s %p\n' | sort -nr | head -30
```

## Troubleshooting

- If browser access fails on localhost, check service types:
  - `kubectl get svc -n llm-observability open-webui ollama`
- If pods are healthy but services are `ClusterIP`, switch to `LoadBalancer` or use `kubectl port-forward`.
- If GPU scheduling fails, verify:
  - `kubectl get nodes -o json | jq '.items[0].status.allocatable'`
  - `kubectl get pods -n nvidia-device-plugin`
