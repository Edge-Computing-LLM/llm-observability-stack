# Live stack validation — 2026-07-17

The chart was rendered and installed against a live single-node k3s cluster after
`k3s-nvidia-edge` advertised one `nvidia.com/gpu` resource.

Verified before rollout:

- NVIDIA automatic profile selection
- GPU infrastructure dependency validation
- Prometheus Operator CRD server-side application
- chart dependency build
- PVC binding through `local-path`
- GGUF host model presence (`gemma-3-1b-it-Q4_K_M.gguf`, approximately 769 MiB)
- Redis and OpenTelemetry Collector readiness
- Ollama GPU scheduling and Open WebUI scheduling
- Open WebUI health endpoint
- live Ollama inference through the locally mounted GGUF model

The first Helm client reached its 15-minute timeout while the first-time
multi-gigabyte images were still downloading. Ollama pulled in approximately 25
minutes and Open WebUI in approximately 33 minutes. No chart, storage, scheduler,
secret, or service error occurred. The interrupted revision was reconciled with a
rollback to its own revision after all resources became ready; the subsequent
unchanged idempotent CLI installs completed in under 20 seconds and final Helm
revision 4 is `deployed`.

Final runtime evidence:

- Ollama, Open WebUI, Redis, and OpenTelemetry Collector were Ready with no
  container restarts.
- Open WebUI `/health` returned `{"status":true}`.
- `ollama list` reported the local `gemma3-1b-it-gguf-local` model at 806 MB.
- The validation prompt returned exactly `validation ok`.
- Ollama detected the GeForce 940M through CUDA, offloaded 23/27 model layers,
  placed 400.8 MiB of weights on CUDA, and reported a 57% CPU / 43% GPU processor
  split for the loaded model.

Static and application validation passed independently:

- Helm lint
- default, CPU, and GeForce 940M template rendering
- Go tests and vet
- 17 Python smoke tests under `/usr/local/bin/python3.11`
