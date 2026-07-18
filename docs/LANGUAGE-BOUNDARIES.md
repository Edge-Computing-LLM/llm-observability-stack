# Language boundaries

The deployable and testable project runtime is Go-first.

## Go

Use Go for the `llm-observability` CLI, Ollama gateway, edge toolbox,
benchmarks, Kubernetes observation workflows, OpenTelemetry instrumentation,
and automated tests. Go binaries must remain non-interactive, support bounded
timeouts, avoid secret output, and keep mutating workflows behind `--yes`.

## Helm and YAML

Helm templates and values own declarative Kubernetes resources. Do not recreate
Helm rendering or release state in Go. Layer 1 NVIDIA resources remain owned by
`k3s-nvidia-edge`; this chart owns only the application/observability layer.

## Bash

Small Bash helpers may glue established tools such as Helm, kubectl, Docker,
nerdctl, and k3s containerd together. Business logic, parsing, benchmarking,
HTTP services, evidence schemas, and test assertions belong in Go.

## Python

Python is not required by the deployed chart, host validation, benchmark,
gateway, toolbox, or runtime-contract observer. It remains only in optional
Jupyter notebooks and their exported learning exercises, because those assets
are explicitly notebook-oriented and disabled in production profiles. Vendored
upstream charts may retain their upstream maintenance scripts; this project
does not execute them.
