# Project Analysis

## Scope

This chart is an application-level stack for:

- LLM serving (`ollama`)
- local chat UI (`open-webui`)
- observability smoke path (`langchain-demo` + `langsmith`)

Node-level GPU enablement remains intentionally external (NVIDIA device-plugin release).

## Current structure

- `Chart.yaml`: umbrella chart metadata + dependency wiring
- `charts/ollama`: vendored dependency
- `charts/open-webui`: vendored dependency
- `templates/`: custom glue resources (secrets, modelfile configmap, demo app)
- `values.yaml`: generic defaults
- `values.local-k3s.example.yaml`: safe template for local deployments

## Runtime data flow

1. Open WebUI sends prompts to Ollama service.
2. Ollama runs a GGUF-backed local model created from startup Modelfile.
3. LangChain demo API calls Ollama and pushes traces/metadata to LangSmith.

## Reliability notes

- Dependency vendoring reduces upstream chart drift.
- Runtime class and GPU resource values are configurable for `nvidia.com/gpu` and `nvidia.com/gpu.shared`.
- Localhost access depends on service exposure (`LoadBalancer` or port-forward).
- Startup package install in `langchain-demo` is convenient but slower/coupled to network.

## Security and publish risks reviewed

- Local secret-bearing files exist in non-example values and rendered output.
- `values.local-k3s.yaml`, `.webui_secret_key`, and `rendered.yaml` must not be committed.
- `.gitignore` now blocks these files.

## Repository size review

- Total project size is currently around `1.4M`.
- Largest files are chart docs/assets; no large model binaries are in tree.
- Size is safely below the requested `100 MB` limit.

## Recommended next hardening steps

- Replace runtime `pip install` path with a prebuilt demo image.
- Add CI checks:
  - `helm lint`
  - `helm template` smoke render
  - secret scan
  - file size guard for blobs > 10MB
