# Live Validation - 2026-07-08

Validated as Layer 2 on local Xubuntu 24, single-node k3s, NVIDIA GeForce 940M.

Prerequisite:

```bash
edge validate infra
```

Tested commands:

```bash
edge install observability --profile geforce-940m-k3s --yes
edge validate observability
edge status
edge logs --tail 30
edge uninstall observability --yes
edge install observability --profile geforce-940m-k3s --yes
```

Results:

- Deployed Ollama, Open WebUI, Open WebUI Redis, and OpenTelemetry Collector in `llm-observability`.
- Ollama used `runtimeClassName: nvidia` and `nvidia.com/gpu: 1`.
- Ollama loaded `gemma3-1b-it-gguf-local` and smoke validation returned `validation ok`.
- No GPU Operator, NVIDIA device plugin, or DCGM exporter resources were rendered by this chart.
- Helm package did not contain old `charts/gpu-operator`, `charts/nvidia-device-plugin`, or `charts/dcgm-exporter` dependencies.
- DCGM exporter remained owned by Layer 1 in namespace `gpu-operator`.

The GeForce 940M profile is valid for the tested low-VRAM workstation flow when Layer 1 is ready first.
