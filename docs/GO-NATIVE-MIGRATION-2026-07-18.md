# Go-native migration — 2026-07-18

Chart `0.3.0` removes Python from the deployed and validated platform path.

| Previous component | Go replacement |
|---|---|
| FastAPI/LangChain proxy | `cmd/ollama-gateway` and `internal/gateway` |
| Python diagnostic image | `cmd/edge-toolbox` and `internal/toolbox` |
| Python streaming benchmark | `internal/benchmark` through `llm-observability benchmark` |
| Python Kubernetes scripts | `network`, `service-path`, and `watch-endpoints` CLI commands |
| pytest Helm/application suite | Go tests in `tests` plus package-level tests |
| Python trace seeder | `edge-toolbox seed` CronJob command |

Values keys changed from `langchainDemo` to `ollamaGateway` and from
`pythonToolbox` to `edgeToolbox`. The optional Kubernetes workload names are
now `ollama-gateway` and `edge-toolbox`. Existing releases using those optional
components must build/import the `0.2.0` Go-native local images before enabling
them.

Jupyter notebooks remain optional Python learning assets and are excluded from
the Helm package. Vendored upstream chart maintenance scripts are not part of
the project runtime and remain unchanged.
