# Documentation Index

This directory contains the long-form documentation for `llm-observability-stack`. The top-level [README.md](../README.md) is the fast entry point; these documents cover architecture, deployment profiles, local operations, troubleshooting, validation, and publishing.

The documentation is organized around private LLM deployment and observability on k3s/Kubernetes with CPU-only and NVIDIA GPU paths. The verified local reference environment is Xubuntu 24 with k3s and an NVIDIA GPU.

## Start Here

1. [DEPENDENCY-AUDIT-2026-07-19.md](DEPENDENCY-AUDIT-2026-07-19.md)
2. [QWEN-1.8B-LIVE-VALIDATION-2026-07-18.md](QWEN-1.8B-LIVE-VALIDATION-2026-07-18.md)
3. [LIVE-VALIDATION-GO-NATIVE-2026-07-18.md](LIVE-VALIDATION-GO-NATIVE-2026-07-18.md)
4. [LIVE-VALIDATION-2026-07-17.md](LIVE-VALIDATION-2026-07-17.md)
5. [QUICKSTART.md](QUICKSTART.md)
6. [cli.md](cli.md)
7. [K3S-NVIDIA-EDGE-DEPENDENCY.md](K3S-NVIDIA-EDGE-DEPENDENCY.md)
8. [CONFIG-PROFILES.md](CONFIG-PROFILES.md)
9. [XUBUNTU-K3S-NVIDIA-RUNBOOK.md](XUBUNTU-K3S-NVIDIA-RUNBOOK.md)
10. [ARCHITECTURE.md](ARCHITECTURE.md)
11. [OPERATIONS-RUNBOOK.md](OPERATIONS-RUNBOOK.md)
12. [PROJECT-DOCUMENTATION.md](PROJECT-DOCUMENTATION.md)
13. [LANGUAGE-BOUNDARIES.md](LANGUAGE-BOUNDARIES.md)
14. [GO-NATIVE-MIGRATION-2026-07-18.md](GO-NATIVE-MIGRATION-2026-07-18.md)

External companion:

- [`qwen-gguf-observability`](https://github.com/Edge-Computing-LLM/qwen-gguf-observability)
  provides read-only Qwen runtime contract checks and sanitized evidence. It
  does not replace this repository's Helm, Modelfile, or benchmark assets.

Dashboard presentation is now owned by the Helm-provisioned Grafana JSON in
[`../dashboards/`](../dashboards/); there is no separate browser application or
browser-side Kubernetes access path.

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
  - Go-first runtime, optional notebook Python, Bash, and Helm/YAML ownership rules.
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
- [GO-KUBERNETES-AUTOMATION.md](GO-KUBERNETES-AUTOMATION.md)
  - Native Go CLI networking inventory, service-path, and endpoint-watch commands.

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
  - [GO-KUBERNETES-AUTOMATION.md](GO-KUBERNETES-AUTOMATION.md)

## Supporting Script Docs

- [scripts/README.md](scripts/README.md)
  - Inventory of standalone helper scripts in `docs/scripts/`.

## Related Component Docs

- [../ollama-gateway/README.md](../ollama-gateway/README.md)
- [../edge-toolbox/README.md](../edge-toolbox/README.md)
- [../hack/README.md](../hack/README.md)
- [../jupyter-notebooks/README.md](../jupyter-notebooks/README.md)
- [../jupyter-notebooks/llm-observability-in-action/README.md](../jupyter-notebooks/llm-observability-in-action/README.md)
