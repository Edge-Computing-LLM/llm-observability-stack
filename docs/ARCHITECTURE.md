# Architecture Guide

This document explains how `llm-observability-stack` is put together, which components own which responsibilities, and how traffic moves through the local k3s deployment.

For local NVIDIA GPU deployments, this repository is not the cluster bootstrap layer. Deploy and validate `k3s-nvidia-edge` first, then deploy `llm-observability-stack` on top of the ready k3s/NVIDIA substrate. See [k3s-nvidia-edge dependency](K3S-NVIDIA-EDGE-DEPENDENCY.md).

The Go CLI in `cmd/llm-observability` follows the same boundary. It imports `github.com/Edge-Computing-LLM/k3s-nvidia-edge/pkg/edgebase` for base-layer workflows and keeps LLM stack workflows in `internal/stack`.

## 1. Design Goals

- Keep the stack understandable on a single local node
- Prefer reproducible local images over runtime `pip install`
- Keep Open WebUI easy to reach from the browser
- Keep internal APIs private by default and expose them only when needed
- Make observability and networking drills easy to demonstrate

## 2. Major Components

### 2.0 Go CLI

Source lives in `cmd/llm-observability` and `internal/stack`.

Responsibilities:

- check base k3s/NVIDIA readiness through `edgebase`
- install and uninstall this Helm chart
- report status for Ollama, Open WebUI, Redis, OpenTelemetry Collector, and services
- validate model loading, Ollama API behavior, and CUDA/offload evidence
- wrap the existing Python benchmark client
- print the Helm and kubectl commands used under the hood

### 2.1 Root umbrella chart

The root chart owns:

- deployment composition
- values layering
- custom templates that glue the subcharts together
- observability dependencies including kube-prometheus-stack and OpenTelemetry Operator
- optional resources such as Redis, python-toolbox, and etcd simulations

Files:

- `Chart.yaml`
- `values.yaml`
- `templates/`

### 2.2 Ollama

Deployed through vendored subchart content in `charts/ollama`.

Responsibilities:

- host-mounted GGUF access
- model runtime
- optional model creation at startup through ConfigMap-backed Modelfile content

### 2.3 Open WebUI

Deployed through vendored subchart content in `charts/open-webui`.

Responsibilities:

- browser-facing chat UI
- user/session state
- Ollama-compatible request flow to the traced proxy path

### 2.4 LangChain demo / proxy

Source lives in `langchain-demo/`.

Responsibilities:

- health and config endpoints
- simple `/invoke` demo endpoint
- Ollama-compatible proxy path at `/ollama/*`
- optional OpenTelemetry-traced proxy runs for Open WebUI traffic

### 2.5 Python toolbox

Source lives in `python-toolbox/`.

Responsibilities:

- in-cluster diagnostics
- DNS and service connectivity checks
- optional OpenTelemetry helper scripts
- notebook support for cluster-side network probing

### 2.6 Redis

Optional root-level resource set. It supports Open WebUI websocket/state flows when that path is enabled for the local profile.

### 2.7 etcd simulation resources

Optional and disabled by default. These exist for troubleshooting and demo scenarios, not for the default happy path.

### 2.8 NVIDIA and observability operators

The root chart vendors NVIDIA-related charts for portability and non-local profiles, but the verified local k3s/NVIDIA path treats `k3s-nvidia-edge` as the owner of the GPU substrate. Keep GPU Operator, NVIDIA device plugin, and DCGM exporter disabled here when `k3s-nvidia-edge` is installed.

Available integration points:

- `gpu-operator` for clusters that need NVIDIA driver/toolkit/device-plugin/DCGM lifecycle management.
- `nvidia-device-plugin` plus `dcgm-exporter` for lightweight workstation or k3s clusters where host drivers and container runtime are already installed.
- `kube-prometheus-stack` for Prometheus Operator, Prometheus, Grafana, Alertmanager, node exporter, and kube-state-metrics.
- `opentelemetry-collector` for a directly managed OTLP collector Deployment and Service.
- `opentelemetry-operator` remains available for clusters that need an operator-managed `OpenTelemetryCollector` custom resource.

## 3. Traffic Flow

Primary user path for the full proxy profile:

1. Browser -> `open-webui` Service
2. `open-webui` pod -> `langchain-demo` Service
3. `langchain-demo` pod -> `ollama` Service
4. `langchain-demo` -> OpenTelemetry API when tracing is enabled
5. `langchain-demo` -> OpenTelemetry Collector when OpenTelemetry is enabled

Primary user path for `values.geforce-940m-k3s.yaml`:

1. Browser -> `open-webui` Service
2. `open-webui` pod -> `ollama` Service
3. `ollama` pod -> NVIDIA GPU through `RuntimeClass/nvidia` and `nvidia.com/gpu`
4. OpenTelemetry Collector remains available for OTLP ingest

Supporting path:

1. Notebook or operator -> `kubectl exec` or Kubernetes Python client
2. `python-toolbox` pod -> internal Services, DNS, OpenTelemetry API

## 4. Exposure Strategy

Default local pattern:

- `open-webui`: externally reachable for browser use
- `ollama`: `ClusterIP`
- `langchain-demo`: `ClusterIP`
- `python-toolbox`: no public Service, pod-only diagnostics

This keeps the local demo usable while reducing unnecessary surface area.

## 5. Configuration Ownership

There are three main configuration layers:

### 5.1 Stable defaults

- `values.yaml`
- tracked in git
- no secrets

### 5.2 Sanitized local example

- `values.local-k3s.example.yaml`
- tracked in git
- onboarding template for local machines

### 5.3 Machine-local overrides

- `values.local-k3s.yaml`
- gitignored
- contains host paths, secrets, and machine-specific overrides

## 6. Component Source Map

- `cmd/llm-observability/`: Go CLI entrypoint
- `internal/stack/`: CLI workflows for this chart and app layer
- `templates/`: root chart resources, optional OpenTelemetryCollector CR, and integration glue
- `langchain-demo/app.py`: FastAPI app and traced proxy logic
- `python-toolbox/examples/`: in-cluster helper scripts
- `hack/`: local image build/import flow
- `jupyter-notebooks/`: notebook-driven operational guides

## 7. Why the Repository Is Structured This Way

- Vendored dependency charts reduce upstream drift during local support, demos, and repeatable validation
- Local image sources make runtime behavior more repeatable than mutable pods
- Notebook and script assets are kept with the chart so the repo is self-contained for local workshops
- Optional components remain in the same repo because they are part of the troubleshooting story

## 8. Operational Boundaries

This repository is optimized for:

- local demos
- workstation troubleshooting
- single-node k3s experimentation
- observability walkthroughs

It is not trying to be a generic multi-node production platform.

## 9. Upstream Base Layer

The verified local NVIDIA flow depends on `k3s-nvidia-edge` for:

- k3s and k3s containerd runtime preparation
- GPU Operator installation in namespace `gpu-operator`
- `RuntimeClass/nvidia`
- allocatable `nvidia.com/gpu`
- DCGM exporter and device plugin lifecycle
- CUDA pod validation

`llm-observability-stack` should be installed only after those checks pass.
