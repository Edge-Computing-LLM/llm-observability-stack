# Programming language and script boundaries

This repository deliberately uses several implementation languages because it
contains a Helm application, a compatibility CLI, an instrumented API, and
operator tooling.

## Go

Use Go in `cmd/llm-observability` and `internal/stack` for durable repo-level
CLI behavior, typed options, process execution, timeouts, errors, and lifecycle
operations. New cross-repository workflows belong in `edge-cli`.

Do not create a second Go implementation of Helm templates, Python benchmark
math, or browser behavior.

## Python 3.11

Use Python 3.11 for:

- the instrumented FastAPI/LangChain service;
- benchmarks and structured result processing;
- Kubernetes API reports and notebook code;
- Helm/application tests;
- in-cluster diagnostics whose container image pins its own interpreter.

Host commands use `python3.11` and install packages directly with
`python3.11 -m pip`. This local workflow does not create virtual environments.
Container commands may use `python` when that name is supplied by the pinned
container image.

## Bash

Keep Bash in `hack/` for short, transparent composition of existing tools such
as Helm, kubectl, Docker/nerdctl, and k3s containerd import. Scripts must use
`set -euo pipefail`, quote variables, expose safe defaults, and delegate complex
logic to Go or Python 3.11.

Do not add growing state machines, JSON parsers, model policy, or duplicated
install/uninstall logic to Bash.

## Helm and YAML

Helm templates and values own Kubernetes desired state. Runtime detection may
write a generated values overlay, but it must not modify tracked defaults.

## TypeScript

Browser functionality belongs in the separate
`Frontend-Edge-LLM-Observability` repository. This chart may integrate that
frontend later, but it should not vendor a second frontend source tree.
