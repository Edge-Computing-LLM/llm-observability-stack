# Contributing

Thanks for contributing to `llm-observability-stack`.

## Scope

This repository is focused on local Kubernetes observability workflows for:

- Ollama
- Open WebUI
- LangChain demo + LangSmith tracing

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

## Validate Before PR

Run these commands before opening a PR:

```bash
helm lint .
helm template llm-observability-stack . >/tmp/rendered-default.yaml
helm template llm-observability-stack . -f values.local-k3s.example.yaml >/tmp/rendered-local.yaml
```

If your change touches local runtime behavior, include the command output and manual test notes in the PR description.

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
