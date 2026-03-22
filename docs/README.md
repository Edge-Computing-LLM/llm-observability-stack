# Documentation Index

This folder contains complete project documentation for `llm-observability-stack`.

Current local profile notes (March 2026):

- `open-webui` is exposed for browser use.
- `ollama` and `langchain-demo` are internal `ClusterIP` by default.
- `pythonToolbox` and `langsmithDashboardSeeder` are disabled by default to avoid background load.

## Core docs

- [PROJECT-DOCUMENTATION.md](PROJECT-DOCUMENTATION.md)
  - End-to-end project guide: architecture, components, values flow, deployment, operations, and troubleshooting.
- [KUBECTL-COMMAND-REFERENCE.md](KUBECTL-COMMAND-REFERENCE.md)
  - Extensive `kubectl` reference for day-0 and day-1 operations.
- [KUBERNETES-NETWORKING.md](KUBERNETES-NETWORKING.md)
  - Deep networking documentation for this stack: Services, DNS, Endpoints, traffic path, and troubleshooting.
- [PYTHON-KUBERNETES-AUTOMATION.md](PYTHON-KUBERNETES-AUTOMATION.md)
  - Python automation guide using the `kubernetes` pip package, with ready-to-run scripts.

## Existing docs

- [GITHUB-PUBLISHING.md](GITHUB-PUBLISHING.md)
  - GitHub publishing workflow and safeguards.
- [PROJECT-ANALYSIS.md](PROJECT-ANALYSIS.md)
  - Current state summary and hardening priorities.

## Scripts

- [scripts/README.md](scripts/README.md)
  - Script inventory and quick usage.
- [scripts/network_inventory.py](scripts/network_inventory.py)
  - Lists Pods, Services, Endpoints, EndpointSlices, and NetworkPolicies.
- [scripts/service_path_inspector.py](scripts/service_path_inspector.py)
  - Traces one Service to selected Pods and endpoint addresses.
- [scripts/watch_endpoints.py](scripts/watch_endpoints.py)
  - Watches endpoint changes for one Service in near real-time.

## Suggested reading order

1. `PROJECT-DOCUMENTATION.md`
2. `KUBERNETES-NETWORKING.md`
3. `KUBECTL-COMMAND-REFERENCE.md`
4. `PYTHON-KUBERNETES-AUTOMATION.md`
