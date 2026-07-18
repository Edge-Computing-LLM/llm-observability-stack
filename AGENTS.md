# Repository instructions

This repository owns the Layer 2 Helm release: Ollama/GGUF, Open WebUI,
OpenTelemetry, Prometheus/Grafana, the Go Ollama gateway, Go edge toolbox,
benchmarks, and dashboards. Layer 1 NVIDIA resources belong to
`k3s-nvidia-edge` and must not be duplicated.

The deployable runtime and automated tests are Go-first. Python is allowed only
inside optional Jupyter learning assets and untouched vendored upstream charts.
Before completing a change run `gofmt`, `go test ./...`, `go vet ./...`,
`helm lint .`, and render the CPU and GeForce profiles. Never commit model
weights, API keys, Kubernetes Secrets, kubeconfigs, or private prompt content.
Keep mutations behind `--yes` and retain model-cleanup rejection safeguards.
