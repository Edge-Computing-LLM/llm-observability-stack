# Project Analysis (Current State)

## Scope

This repository is a local k3s-focused LLM platform chart:

- Ollama for model serving
- Open WebUI for chat interface
- Ollama gateway API for integration and tracing checks
- Go edge toolbox for in-cluster diagnostics

Current runtime preference in local overrides:

- keep only Open WebUI externally exposed
- keep Ollama/LangChain internal and access via port-forward as needed
- disable continuous inference seeding jobs

## What Is Stable Today

- Umbrella chart and vendored dependencies are pinned and render cleanly.
- Local image workflow for `ollama-gateway` and `edge-toolbox` is documented and scriptable.
- GGUF Modelfile creation path is implemented and deployable.
- CI validates chart linting and template rendering.
- OpenTelemetry tracing is wired through the `ollama-gateway` proxy path (`/ollama/api/*`) for on-demand observability from Open WebUI traffic.

## Main Risk Areas

1. Secrets are often handled via local plaintext override files.
2. Value-key split (`openWebUI.*` vs `open-webui.*`) is easy to misconfigure.
3. Existing secret wiring for Open WebUI needs careful validation in each environment.
4. Networking and Redis websocket path can drift if both Redis modes are mixed.
5. CI does not yet run runtime or end-to-end validation.

## Recommended Hardening Priorities

1. Standardize secret handling on pre-created Kubernetes secrets and `existingSecret` references.
2. Normalize and document values ownership clearly across root chart and subchart keys.
3. Add runtime post-deploy smoke checks to CI/CD and local workflows.
4. Add stricter policy for service exposure (`LoadBalancer` only when needed).
5. Add optional NetworkPolicies for tighter internal traffic boundaries.

## Documentation Set

For full details, use:

- [README.md](README.md)
- [PROJECT-DOCUMENTATION.md](PROJECT-DOCUMENTATION.md)
- [KUBECTL-COMMAND-REFERENCE.md](KUBECTL-COMMAND-REFERENCE.md)
- [KUBERNETES-NETWORKING.md](KUBERNETES-NETWORKING.md)
- [GO-KUBERNETES-AUTOMATION.md](GO-KUBERNETES-AUTOMATION.md)
- [scripts/README.md](scripts/README.md)
