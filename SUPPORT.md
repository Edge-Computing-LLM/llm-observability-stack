# Support

## Getting Help

Use GitHub Issues for setup, debugging, and enhancement questions:

- Repository: https://github.com/waqasm86/llm-observability-stack
- New issue: https://github.com/waqasm86/llm-observability-stack/issues/new

## What to Include

Please include:

- Kubernetes distro and version (for example `k3s v1.30.x`)
- GPU/device-plugin details (`nvidia.com/gpu` vs `nvidia.com/gpu.shared`)
- Helm command used
- Values files used (`values.yaml`, `values.local-k3s.yaml` override summary)
- Relevant `kubectl get pods -n llm-observability` and service output
- Relevant logs from failing components

Do not include secrets, tokens, or passwords.

## Troubleshooting First

Before opening an issue, run:

```bash
helm lint .
helm template llm-observability-stack . >/tmp/rendered-default.yaml
helm template llm-observability-stack . -f values.local-k3s.example.yaml >/tmp/rendered-local.yaml
```
