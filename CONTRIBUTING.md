# Contributing

Thanks for contributing to `llm-observability-stack`.

## Scope

This repository is focused on local Kubernetes observability workflows for:

- Ollama
- Open WebUI
- LangChain demo + OpenTelemetry tracing

Please keep changes aligned with local k3s reproducibility and interview-support use cases.

## Development Setup

1. Clone the repository.
2. Install prerequisites: `helm`, `kubectl`, and Docker or `nerdctl`.
3. Build chart dependencies:

```bash
helm dependency build .
```

## Values File Rules

- Keep `values.yaml` generic and non-sensitive.
- Do not commit secrets to git-tracked files.
- Use `values.local-k3s.yaml` for machine-specific overrides and secrets (this file is gitignored).
- Keep `values.local-k3s.example.yaml` as a safe template only.
- Update [docs/CONFIG-PROFILES.md](docs/CONFIG-PROFILES.md) when profile defaults or local-example expectations change.

## Validate Before PR

Run these commands before opening a PR:

```bash
helm lint .
helm template llm-observability-stack . >/tmp/rendered-default.yaml
helm template llm-observability-stack . -f values.local-k3s.example.yaml >/tmp/rendered-local.yaml
pytest -q tests/test_helm_smoke.py tests/test_langchain_demo_smoke.py
```

If your change touches local runtime behavior, include the command output and manual test notes in the PR description.

If your change touches notebooks or notebook support code:

- clear saved outputs before committing
- do not commit `.ipynb_checkpoints/`
- update [jupyter-notebooks/CATALOG.md](jupyter-notebooks/CATALOG.md) when adding or retiring notebooks
- update [docs/NOTEBOOKS-GUIDE.md](docs/NOTEBOOKS-GUIDE.md) if prerequisites or workflow assumptions change

If your change touches operational docs:

- keep [docs/QUICKSTART.md](docs/QUICKSTART.md) focused on first deploy
- keep [docs/OPERATIONS-RUNBOOK.md](docs/OPERATIONS-RUNBOOK.md) focused on day-1 operations
- keep [docs/KUBECTL-COMMAND-REFERENCE.md](docs/KUBECTL-COMMAND-REFERENCE.md) and [jupyter-notebooks/llm-observability-in-action/KUBECTL_COMMAND_CATALOG.md](jupyter-notebooks/llm-observability-in-action/KUBECTL_COMMAND_CATALOG.md) workload-kind accurate

## Local Test Matrix

Use the smallest validation set that actually exercises your change:

- Docs only:
  - verify links, commands, and profile notes manually
- Helm/chart changes:
  - `helm lint .`
  - `helm template llm-observability-stack .`
  - `helm template llm-observability-stack . -f values.local-k3s.example.yaml`
- Python or template logic changes:
  - `pytest -q tests/test_helm_smoke.py tests/test_langchain_demo_smoke.py`
- Notebook source changes:
  - `find jupyter-notebooks -type f -name '*.ipynb' ! -path '*/.ipynb_checkpoints/*' -print0 | xargs -0 -n1 jq empty`

## Documentation and Notebook Style

- Prefer relative paths and environment variables over machine-specific absolute paths.
- Prefer `PYTHON_BIN` examples over hard-coded interpreter paths in user-facing guides.
- Keep secrets out of notebook source and output.
- When describing runtime defaults, point readers to [docs/CONFIG-PROFILES.md](docs/CONFIG-PROFILES.md) instead of duplicating tables everywhere.

## Commit and PR Guidelines

- Use small, focused commits.
- Write imperative commit messages (for example: `Add Redis auth wiring for local websocket manager`).
- In the PR description include:
  - What changed
  - Why it changed
  - How it was validated
  - Any security implications

## Security

If you discover a security issue, do not open a public issue with exploit details. Follow [SECURITY.md](SECURITY.md).
