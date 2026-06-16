# Configuration Profiles

This repository uses three configuration layers. Keep their responsibilities separate to avoid drift:

1. `values.yaml`
   - Generic, git-tracked chart defaults.
   - Safe for CI, docs, and cross-machine review.
2. `values.local-k3s.example.yaml`
   - Sanitized local k3s example profile.
   - Shows the intended workstation overrides without storing live secrets.
3. `values.enterprise-pilot-k3s.yaml`
   - Verified local k3s/NVIDIA workstation profile used for the current end-to-end deployment.
   - Keeps the Ollama PVC at `5Gi` to match the existing local-path claim and avoid unsupported resize attempts.
4. `values.local-k3s.yaml`
   - Your real machine-specific override file.
   - Gitignored and expected to contain local paths, secret references, or direct secret values.

## Source of Truth

- Treat [`values.yaml`](../values.yaml) as the baseline behavior of the chart.
- Treat [`values.enterprise-pilot-k3s.yaml`](../values.enterprise-pilot-k3s.yaml) as the verified local single-node k3s/NVIDIA profile for this repository.
- Treat [`values.local-k3s.example.yaml`](../values.local-k3s.example.yaml) as a sanitized example for creating private local overrides.
- Treat `values.local-k3s.yaml` as private runtime state that may legitimately differ from the example.

## At-a-Glance Differences

| Area | `values.yaml` | `values.local-k3s.example.yaml` | `values.enterprise-pilot-k3s.yaml` |
|---|---|---|---|
| GPU resource | `nvidia.com/gpu` | `nvidia.com/gpu` | `nvidia.com/gpu` |
| GGUF host path | Placeholder path | Example host path under `/media/.../models` | Local verified host path under `/media/.../models` |
| Ollama PVC size | `20Gi` | `5Gi` | `5Gi` |
| Service exposure | `ClusterIP` | `ClusterIP` | `ClusterIP` |
| LangSmith API key | Empty string | `replace-me` placeholder | disabled |
| Open WebUI secret key | Empty string | placeholder | placeholder, chart-managed secret |
| `pythonToolbox.enabled` | `true` | `true` | `true` |
| `langsmithDashboardSeeder.enabled` | `false` | `false` | `false` |
| `etcd.enabled` | `false` | `false` | `false` |

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

Render with the verified enterprise pilot profile:

```bash
helm template llm-observability-stack . \
  -n llm-observability \
  -f values.enterprise-pilot-k3s.yaml \
  --set kube-prometheus-stack.crds.enabled=false \
  >/tmp/rendered-enterprise-pilot.yaml
```

Render with your private local profile:

```bash
helm template llm-observability-stack . -f values.local-k3s.yaml >/tmp/rendered-local-private.yaml
```

## Rules for Contributors

- Do not place live secrets in `values.yaml` or `values.local-k3s.example.yaml`.
- Keep `values.yaml` portable and workstation-agnostic.
- When the local example changes, update this file and any quick-start docs that reference the profile.
