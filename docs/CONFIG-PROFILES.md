# Configuration Profiles

This repository uses three configuration layers. Keep their responsibilities separate to avoid drift:

1. `values.yaml`
   - Generic, git-tracked chart defaults.
   - Safe for CI, docs, and cross-machine review.
2. `values.local-k3s.example.yaml`
   - Sanitized local k3s example profile.
   - Shows the intended workstation overrides without storing live secrets.
3. `values.local-k3s.yaml`
   - Your real machine-specific override file.
   - Gitignored and expected to contain local paths, secret references, or direct secret values.

## Source of Truth

- Treat [`values.yaml`](../values.yaml) as the baseline behavior of the chart.
- Treat [`values.local-k3s.example.yaml`](../values.local-k3s.example.yaml) as the canonical example for a local NVIDIA-enabled k3s workstation.
- Treat `values.local-k3s.yaml` as private runtime state that may legitimately differ from the example.

## At-a-Glance Differences

| Area | `values.yaml` | `values.local-k3s.example.yaml` | Why it differs |
|---|---|---|---|
| GPU resource | `nvidia.com/gpu` | `nvidia.com/gpu.shared` | The example profile is tuned for the local host's MPS-sharing setup. |
| GGUF host path | Placeholder path | Example host path under `/media/.../models` | Local workstations need a real host path for the mounted model directory. |
| Ollama service type | `ClusterIP` | `LoadBalancer` | The example profile exposes Ollama directly for host-side demos. |
| LangChain demo service type | `ClusterIP` | `LoadBalancer` | The example profile exposes the proxy directly for easier local testing. |
| Open WebUI service type | `ClusterIP` | `LoadBalancer` | The example profile expects direct browser access. |
| LangSmith API key | Empty string | `replace-me` placeholder | The example signals that the operator must provide a real secret value or existing secret. |
| Open WebUI secret key | Empty string | `replace-with-a-32-plus-char-secret` placeholder | The example makes the required secret explicit without committing a real key. |
| `pythonToolbox.enabled` | `true` | `true` | The current local workflow expects the toolbox to be available by default. |
| `langsmithDashboardSeeder.enabled` | `false` | `false` | Seeder stays off unless you intentionally enable periodic synthetic traces. |
| `etcd.enabled` | `false` | `false` | etcd simulation remains opt-in for troubleshooting exercises. |

## Recommended Workflow

Create a local file from the example:

```bash
cp values.local-k3s.example.yaml values.local-k3s.yaml
```

Then edit only the machine-specific fields in `values.local-k3s.yaml`:

- model host paths
- direct secret values or existing secret references
- local service exposure preferences
- any temporary troubleshooting toggles

## Verify Effective Runtime Values

See the fully merged release values:

```bash
helm get values llm-observability-stack -n llm-observability -a
```

Render with the example profile:

```bash
helm template llm-observability-stack . -f values.local-k3s.example.yaml >/tmp/rendered-local-example.yaml
```

Render with your private local profile:

```bash
helm template llm-observability-stack . -f values.local-k3s.yaml >/tmp/rendered-local-private.yaml
```

## Rules for Contributors

- Do not place live secrets in `values.yaml` or `values.local-k3s.example.yaml`.
- Keep `values.yaml` portable and workstation-agnostic.
- When the local example changes, update this file and any quick-start docs that reference the profile.
