# Python Kubernetes Automation Guide (`pip install kubernetes`)

This guide shows how to automate cluster/network checks for `llm-observability-stack` using the official Kubernetes Python client.

Profile note:

- These scripts can be run from host, CI, or a temporary debug pod.
- The default local profile currently keeps `python-toolbox` disabled for lower baseline memory usage.

## 1. Install and Environment

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install kubernetes
```

Or use project script requirements:

```bash
pip install -r docs/scripts/requirements.txt
```

## 2. Authentication Modes

### 2.1 Local kubeconfig (most common)

```python
from kubernetes import config
config.load_kube_config()
```

### 2.2 Specific kubeconfig/context

```python
from kubernetes import config
config.load_kube_config(config_file="/path/to/kubeconfig", context="my-k3s")
```

### 2.3 In-cluster auth

```python
from kubernetes import config
config.load_incluster_config()
```

## 3. Ready-to-use Scripts in this Repository

- `docs/scripts/network_inventory.py`
- `docs/scripts/service_path_inspector.py`
- `docs/scripts/watch_endpoints.py`

Quick usage:

```bash
python docs/scripts/network_inventory.py --namespace llm-observability
python docs/scripts/service_path_inspector.py --namespace llm-observability --service ollama
python docs/scripts/watch_endpoints.py --namespace llm-observability --service open-webui --timeout 600
```

## 4. Practical Python Examples

### 4.1 List pods and services

```python
from kubernetes import client, config

config.load_kube_config()
ns = "llm-observability"
core = client.CoreV1Api()

pods = core.list_namespaced_pod(ns).items
services = core.list_namespaced_service(ns).items

print("Pods:")
for p in pods:
    print(f"- {p.metadata.name} phase={p.status.phase} ip={p.status.pod_ip}")

print("Services:")
for s in services:
    ports = ", ".join(str(port.port) for port in (s.spec.ports or []))
    print(f"- {s.metadata.name} type={s.spec.type} cluster_ip={s.spec.cluster_ip} ports={ports}")
```

### 4.2 Validate a service endpoint path

```python
from kubernetes import client, config

config.load_kube_config()
ns = "llm-observability"
service_name = "ollama"

core = client.CoreV1Api()
svc = core.read_namespaced_service(service_name, ns)
selector = svc.spec.selector or {}
label_query = ",".join(f"{k}={v}" for k, v in selector.items())

pods = core.list_namespaced_pod(ns, label_selector=label_query).items if label_query else []
ep = core.read_namespaced_endpoints(service_name, ns)

print(f"service={service_name} selector={selector}")
print(f"selected_pods={[p.metadata.name for p in pods]}")

addresses = []
for subset in ep.subsets or []:
    for addr in subset.addresses or []:
        addresses.append(addr.ip)

print(f"endpoint_addresses={addresses}")
```

### 4.3 Watch endpoint changes

```python
from kubernetes import client, config, watch

config.load_kube_config()
ns = "llm-observability"
service_name = "ollama"

core = client.CoreV1Api()
w = watch.Watch()

for event in w.stream(
    core.list_namespaced_endpoints,
    namespace=ns,
    field_selector=f"metadata.name={service_name}",
    timeout_seconds=300,
):
    etype = event["type"]
    obj = event["object"]
    print(etype, obj.metadata.name, obj.subsets)
```

### 4.4 List EndpointSlices for detailed networking state

```python
from kubernetes import client, config

config.load_kube_config()
ns = "llm-observability"
service_name = "open-webui"

discovery = client.DiscoveryV1Api()
items = discovery.list_namespaced_endpoint_slice(
    namespace=ns,
    label_selector=f"kubernetes.io/service-name={service_name}",
).items

for eps in items:
    addrs = []
    for endpoint in eps.endpoints or []:
        addrs.extend(endpoint.addresses or [])
    print(eps.metadata.name, addrs)
```

### 4.5 Check NetworkPolicies

```python
from kubernetes import client, config

config.load_kube_config()
ns = "llm-observability"
net = client.NetworkingV1Api()

policies = net.list_namespaced_network_policy(ns).items
for np in policies:
    types = np.spec.policy_types if np.spec else []
    print(np.metadata.name, types)
```

## 5. Usage Patterns for This Project

### 5.1 Verify post-deploy networking quickly

1. Run `network_inventory.py` to capture current state.
2. Run `service_path_inspector.py` for `ollama`, `open-webui`, `langchain-demo`.
3. Run `watch_endpoints.py` during rollout/restart events.

### 5.2 Diagnose intermittent connection failures

- Collect service + endpoints snapshots over time.
- Compare selected Pod readiness vs endpoint membership.
- Correlate endpoint changes with Kubernetes events:

```bash
kubectl get events -n llm-observability --sort-by=.lastTimestamp
```

## 6. Error Handling Recommendations

When using the Python client in production scripts:

- Catch `kubernetes.client.ApiException` and print status/reason/body.
- Retry transient HTTP status (`429`, `500`, `503`) with exponential backoff.
- Set command timeout boundaries for watch loops.
- Emit machine-readable output (JSON) for CI pipelines.

## 7. Security and Safety

- Avoid printing decoded secret values in script output.
- Use least-privilege service accounts for in-cluster execution.
- Keep kubeconfig permissions tight (`chmod 600`).
- Avoid broad wildcard RBAC if script is namespace-scoped.

## 8. Next-Level Automation Ideas

- Scheduled network baseline snapshots to file.
- Drift detection: compare expected services/selectors to live cluster.
- Automatic alert when service endpoints become empty.
- Latency checks between internal services using exec-based probes.

## 9. Related docs

- `docs/KUBERNETES-NETWORKING.md`
- `docs/KUBECTL-COMMAND-REFERENCE.md`
- `docs/scripts/README.md`
