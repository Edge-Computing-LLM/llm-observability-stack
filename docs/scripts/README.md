# Replaced documentation scripts

The former Python 3.11 Kubernetes scripts were replaced by native Go commands:

| Former script | Go command |
|---|---|
| `network_inventory.py` | `bin/llm-observability network` |
| `service_path_inspector.py` | `bin/llm-observability service-path --service NAME` |
| `watch_endpoints.py` | `bin/llm-observability watch-endpoints --service NAME` |

See [Native Go Kubernetes automation](../GO-KUBERNETES-AUTOMATION.md).
