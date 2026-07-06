# Architecture Guide

This document explains how `llm-observability-stack` is put together, which components own which responsibilities, and how traffic moves through the local k3s deployment.

## 1. Design Goals

- Keep the stack understandable on a single local node
- Prefer reproducible local images over runtime `pip install`
- Keep Open WebUI easy to reach from the browser
- Keep internal APIs private by default and expose them only when needed
- Make observability and networking drills easy to demonstrate

## 2. Major Components

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

The root chart vendors and controls the NVIDIA and observability integration points required for local GPU-aware LLM operations:

- `gpu-operator` for clusters that need NVIDIA driver/toolkit/device-plugin/DCGM lifecycle management.
- `nvidia-device-plugin` plus `dcgm-exporter` for lightweight workstation or k3s clusters where host drivers and container runtime are already installed.
- `kube-prometheus-stack` for Prometheus Operator, Prometheus, Grafana, Alertmanager, node exporter, and kube-state-metrics.
- `opentelemetry-collector` for a directly managed OTLP collector Deployment and Service.
- `opentelemetry-operator` remains available for clusters that need an operator-managed `OpenTelemetryCollector` custom resource.

## 3. Traffic Flow

Primary user path:

1. Browser -> `open-webui` Service
2. `open-webui` pod -> `langchain-demo` Service
3. `langchain-demo` pod -> `ollama` Service
4. `langchain-demo` -> OpenTelemetry API when tracing is enabled
5. `langchain-demo` -> OpenTelemetry Collector when OpenTelemetry is enabled

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
