# Python Kubernetes Scripts

These scripts use the official Python Kubernetes client (`pip install kubernetes`) and focus on networking diagnostics for `llm-observability-stack`.

They are designed to run from host or CI and do not require `python-toolbox` to be enabled in-cluster.

## Install

```bash
/usr/local/bin/python3.11 -m pip install -r docs/scripts/requirements.txt
```

## Scripts

### 1) `network_inventory.py`

Purpose:

- Print namespace networking inventory for Pods, Services, Endpoints, EndpointSlices, and NetworkPolicies.

Usage:

```bash
/usr/local/bin/python3.11 docs/scripts/network_inventory.py --namespace llm-observability
```

### 2) `service_path_inspector.py`

Purpose:

- Trace one Service to its selector, selected Pods, Endpoints, and EndpointSlices.

Usage:

```bash
/usr/local/bin/python3.11 docs/scripts/service_path_inspector.py --namespace llm-observability --service ollama
/usr/local/bin/python3.11 docs/scripts/service_path_inspector.py --namespace llm-observability --service open-webui
/usr/local/bin/python3.11 docs/scripts/service_path_inspector.py --namespace llm-observability --service langchain-demo
```

### 3) `watch_endpoints.py`

Purpose:

- Watch endpoint changes for a specific Service in near real-time.

Usage:

```bash
/usr/local/bin/python3.11 docs/scripts/watch_endpoints.py --namespace llm-observability --service ollama --timeout 600
```

## Config options

All scripts support:

- `--namespace`
- `--context`
- `--kubeconfig`
- `--in-cluster`

Use `--in-cluster` only when running script inside Kubernetes with service account auth.
