# Documentation Index

This directory contains the long-form documentation for `llm-observability-stack`. The top-level [README.md](../README.md) is the fast entry point; these documents cover architecture, deployment profiles, local operations, troubleshooting, validation, and publishing.

The documentation is organized around private LLM deployment and observability on k3s/Kubernetes with CPU-only and NVIDIA GPU paths. The verified local reference environment is Xubuntu 24 with k3s and an NVIDIA GPU.

## Start Here

1. [QWEN-1.8B-LIVE-VALIDATION-2026-07-18.md](QWEN-1.8B-LIVE-VALIDATION-2026-07-18.md)
2. [LIVE-VALIDATION-2026-07-17.md](LIVE-VALIDATION-2026-07-17.md)
3. [QUICKSTART.md](QUICKSTART.md)
4. [cli.md](cli.md)
5. [K3S-NVIDIA-EDGE-DEPENDENCY.md](K3S-NVIDIA-EDGE-DEPENDENCY.md)
6. [CONFIG-PROFILES.md](CONFIG-PROFILES.md)
7. [XUBUNTU-K3S-NVIDIA-RUNBOOK.md](XUBUNTU-K3S-NVIDIA-RUNBOOK.md)
8. [ARCHITECTURE.md](ARCHITECTURE.md)
9. [OPERATIONS-RUNBOOK.md](OPERATIONS-RUNBOOK.md)
10. [PROJECT-DOCUMENTATION.md](PROJECT-DOCUMENTATION.md)
11. [LANGUAGE-BOUNDARIES.md](LANGUAGE-BOUNDARIES.md)

External companion:

- [`qwen-gguf-observability`](https://github.com/Edge-Computing-LLM/qwen-gguf-observability)
  provides read-only Qwen runtime contract checks and sanitized evidence. It
  does not replace this repository's Helm, Modelfile, or benchmark assets.
- [`Frontend-Edge-LLM-Observability`](https://github.com/Edge-Computing-LLM/Frontend-Edge-LLM-Observability)
  owns the TypeScript/Vue presentation layer and does not receive direct
  Kubernetes credentials from this chart.

## Core Guides

- [QUICKSTART.md](QUICKSTART.md)
  - Fast local setup for k3s, values files, image build/import, install, and first validation.
- [cli.md](cli.md)
  - Go CLI build instructions, commands, profile mapping, install/validate/benchmark/uninstall examples, and two-layer architecture notes.
- [K3S-NVIDIA-EDGE-DEPENDENCY.md](K3S-NVIDIA-EDGE-DEPENDENCY.md)
  - Required install order and ownership boundary between `k3s-nvidia-edge` and this LLM stack.
- [CONFIG-PROFILES.md](CONFIG-PROFILES.md)
  - Canonical comparison of git-tracked defaults, local example values, GPU profiles, CPU profiles, and private overrides.
- [XUBUNTU-K3S-NVIDIA-RUNBOOK.md](XUBUNTU-K3S-NVIDIA-RUNBOOK.md)
  - Current local runbook for Xubuntu 24 + k3s + NVIDIA GPU, based on the active local deployment.
- [LOCAL-K3S-NVIDIA-RUNBOOK.md](LOCAL-K3S-NVIDIA-RUNBOOK.md)
  - Existing local k3s/NVIDIA command reference.
- [ARCHITECTURE.md](ARCHITECTURE.md)
  - Component ownership, request paths, service exposure, and configuration boundaries.
- [LANGUAGE-BOUNDARIES.md](LANGUAGE-BOUNDARIES.md)
  - Go, Python 3.11, Bash, Helm/YAML, and frontend ownership rules.
- [OPERATIONS-RUNBOOK.md](OPERATIONS-RUNBOOK.md)
  - Day-0 and day-1 tasks: deploy, verify, port-forward, rebuild images, debug, and clean up.
- [NOTEBOOKS-GUIDE.md](NOTEBOOKS-GUIDE.md)
  - Walkthrough of notebooks `01` through `10`, their prerequisites, and common execution pitfalls.
- [PROJECT-DOCUMENTATION.md](PROJECT-DOCUMENTATION.md)
  - Full repository documentation, component walkthroughs, and deployment model.
- [PROJECT-ANALYSIS.md](PROJECT-ANALYSIS.md)
  - Current-state summary and hardening priorities.

## Kubernetes and Automation Guides

- [KUBERNETES-NETWORKING.md](KUBERNETES-NETWORKING.md)
  - Service, EndpointSlice, DNS, ServiceLB, and traffic-path documentation for this stack.
- [KUBECTL-COMMAND-REFERENCE.md](KUBECTL-COMMAND-REFERENCE.md)
  - High-signal `kubectl` command catalog for local operations.
- [PYTHON-KUBERNETES-AUTOMATION.md](PYTHON-KUBERNETES-AUTOMATION.md)
  - Kubernetes Python client usage patterns and script-driven inspection.

## Local Validation

- [VERIFIED-LOCAL-GPU-RESULTS.md](VERIFIED-LOCAL-GPU-RESULTS.md)
  - Historical Gemma-era sanitized GPU benchmark and deployment results.
- [LOCAL-K3S-NVIDIA-REPORT-2026-07-02.md](LOCAL-K3S-NVIDIA-REPORT-2026-07-02.md)
  - Historical Gemma-era k3s/NVIDIA status report for the Xubuntu 24 system.
- [SINGLE-NODE-K3S-GEFORCE-940M.md](SINGLE-NODE-K3S-GEFORCE-940M.md)
  - Low-memory GeForce 940M profile notes and constraints.
- [../validation/README.md](../validation/README.md)
  - Validation artifact index.

## Git and Publishing

- [GITHUB-PUBLISHING.md](GITHUB-PUBLISHING.md)
  - Remote setup, safe publishing workflow, and repository hygiene guidance.

## Suggested Reading Paths

- First deploy:
  - [cli.md](cli.md)
  - [K3S-NVIDIA-EDGE-DEPENDENCY.md](K3S-NVIDIA-EDGE-DEPENDENCY.md)
  - [QUICKSTART.md](QUICKSTART.md)
  - [CONFIG-PROFILES.md](CONFIG-PROFILES.md)
  - [XUBUNTU-K3S-NVIDIA-RUNBOOK.md](XUBUNTU-K3S-NVIDIA-RUNBOOK.md)
- CPU-only or MacOS/minikube:
  - [../ReadMe-MacOS.md](../../ReadMe-MacOS.md)
  - [CONFIG-PROFILES.md](CONFIG-PROFILES.md)
- Contributor/operator:
  - [../CONTRIBUTING.md](../CONTRIBUTING.md)
  - [KUBECTL-COMMAND-REFERENCE.md](KUBECTL-COMMAND-REFERENCE.md)
  - [PYTHON-KUBERNETES-AUTOMATION.md](PYTHON-KUBERNETES-AUTOMATION.md)

## Supporting Script Docs

- [scripts/README.md](scripts/README.md)
  - Inventory of standalone helper scripts in `docs/scripts/`.

## Related Component Docs

- [../langchain-demo/README.md](../langchain-demo/README.md)
- [../python-toolbox/README.md](../python-toolbox/README.md)
- [../hack/README.md](../hack/README.md)
- [../jupyter-notebooks/README.md](../jupyter-notebooks/README.md)
- [../jupyter-notebooks/llm-observability-in-action/README.md](../jupyter-notebooks/llm-observability-in-action/README.md)
