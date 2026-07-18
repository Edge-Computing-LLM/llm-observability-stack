# Notebook Catalog

This catalog classifies the notebook material in this repository so new users can tell what is core, what is advanced, and what is historical reference.

## Core Walkthrough

These notebooks form the main guided path for the current stack:

1. `01-environment-smoke-test.ipynb`
2. `02-ollama-api-basics.ipynb`
3. `03-ollama-gateway-deep-dive.ipynb`
4. `04-opentelemetry-tracing-setup.ipynb`
5. `05-open-webui-end-to-end.ipynb`
6. `06-custom-modelfile-workflow.ipynb`
7. `07-edge-toolbox-diagnostics.ipynb`
8. `08-troubleshooting-etcd-simulations.ipynb`
9. `09-k3s-networking-deep-dive.ipynb`
10. `10-k3s-architecture-diagrams.ipynb`

## Advanced Companion Material

- `llm-observability-in-action/`
  - Script-driven companion toolkit for kubectl workflows, networking manifests, and Kubernetes Python client examples.

## Supporting Example Notebooks

These notebooks are useful as extra examples, but they are not the primary learning path:

- `llm-observability-stack-example-1.ipynb`
- `llm-observability-stack-example-2.ipynb`
- `llm-observability-stack-example-3.ipynb`
- `llm-observability-stack-example-4.ipynb`
- `llm-observability-stack-example-5.ipynb`

Treat them as reference snapshots or focused demos rather than as the canonical current workflow.

## Deprecated or Local-Only Material

- `jupyter-notebooks-2/`
  - Local duplicate working tree.
  - Ignored by git and not part of the published workflow.

Notebook checkpoint files under `.ipynb_checkpoints/` are also local-only and should not be committed.

## Prerequisite Notes

- Notebook kernels should use Python 3.11.
- Several notebooks require `kubectl port-forward` to `ollama`, `ollama-gateway`, or `open-webui`.
- `07` and `09` expect `edgeToolbox.enabled=true` in the running release or active local values.

Use [README.md](README.md) for launch instructions and [../docs/NOTEBOOKS-GUIDE.md](../docs/NOTEBOOKS-GUIDE.md) for deeper execution guidance.
