# Operations Runbook

This runbook collects the day-0 and day-1 commands you are most likely to need while operating `llm-observability-stack` on a local k3s machine.

## 1. Before You Change Anything

Check current state:

```bash
kubectl get all -n llm-observability
kubectl get svc -n llm-observability
helm list -n llm-observability
git status --short
```

## 2. Build and Refresh Local Images

### Rebuild `ollama-gateway`

```bash
./hack/build-local-image.sh ollama-gateway 0.2.0 . ollama-gateway/Dockerfile
./hack/import-local-image-to-k3s.sh ollama-gateway 0.2.0
kubectl rollout restart deploy/ollama-gateway -n llm-observability
kubectl rollout status deploy/ollama-gateway -n llm-observability
```

### Rebuild `edge-toolbox`

```bash
./hack/build-local-image.sh edge-toolbox 0.2.0 . edge-toolbox/Dockerfile
./hack/import-local-image-to-k3s.sh edge-toolbox 0.2.0
kubectl rollout restart deploy/edge-toolbox -n llm-observability
kubectl rollout status deploy/edge-toolbox -n llm-observability
```

## 3. Install, Upgrade, and Roll Back

### Install or upgrade

```bash
helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.enterprise-pilot-k3s.yaml \
  --set kube-prometheus-stack.crds.enabled=false
```

Use a private `values.local-k3s.yaml` only when your host paths, secrets, or service exposure differ. Keep `ollama.persistentVolume.size` aligned with any existing `ollama` PVC; k3s `local-path` storage does not resize that claim in place.

### Inspect release values

```bash
helm get values llm-observability-stack -n llm-observability -a
```

### Review history and rollback

```bash
helm history llm-observability-stack -n llm-observability
helm rollback llm-observability-stack <REVISION> -n llm-observability
```

### Uninstall

```bash
helm uninstall llm-observability-stack -n llm-observability
```

## 4. Health Checks

### Basic workload checks

```bash
kubectl get pods -n llm-observability -o wide
kubectl get svc -n llm-observability
kubectl get pvc -n llm-observability
kubectl top pods -n llm-observability
```

### Service-specific logs

```bash
kubectl logs -n llm-observability deploy/ollama-gateway --tail=100
kubectl logs -n llm-observability deploy/ollama --tail=100
kubectl logs -n llm-observability statefulset/open-webui --tail=100
kubectl logs -n llm-observability deploy/edge-toolbox --tail=100
```

## 5. Access Internal APIs

```bash
kubectl port-forward -n llm-observability svc/ollama 11434:11434
kubectl port-forward -n llm-observability svc/ollama-gateway 8000:8000
```

Use these for:

- direct Ollama API tests
- Ollama gateway notebook cells
- OpenTelemetry traced proxy requests from the host

## 6. Jupyter Notebook Operations

Launch:

```bash
PYTHON_BIN="${PYTHON_BIN:-python3.11}"
cd jupyter-notebooks
"${PYTHON_BIN}" -m jupyter lab
```

Useful pairings:

- `01` before any major change
- `07` when validating `edge-toolbox`
- `09` when validating cluster networking

If notebook cells fail:

- confirm the required port-forwards
- confirm the release is installed
- confirm the toolbox pod is running
- confirm OpenTelemetry environment variables for tracing notebooks

## 7. In-Cluster Toolbox Checks

Run a diagnostic command. The hardened distroless image intentionally has no
interactive shell:

```bash
kubectl exec -n llm-observability deploy/edge-toolbox -- edge-toolbox dns ollama ollama-gateway
```

Run helper scripts:

```bash
kubectl exec -it -n llm-observability deploy/edge-toolbox -- edge-toolbox dns ollama ollama-gateway open-webui open-webui-redis
kubectl exec -it -n llm-observability deploy/edge-toolbox -- edge-toolbox ollama-smoke
kubectl exec -it -n llm-observability deploy/edge-toolbox -- edge-toolbox redis-ping
```

## 8. GPU Checks

```bash
watch -n 0.5 nvidia-smi
kubectl get pods -n nvidia-device-plugin
kubectl get nodes -o json | jq '.items[0].status.allocatable'
```

If scheduling is failing, inspect the node plugin runtime behavior and confirm the resource name requested in `values.enterprise-pilot-k3s.yaml` or your private local override.

## 9. Troubleshooting Patterns

### `open-webui` works but notebooks fail

- likely missing `kubectl port-forward` for `ollama` and/or `ollama-gateway`

### `ollama-gateway` is unhealthy

- inspect logs
- confirm the local image tag imported into k3s matches the chart values
- confirm OpenTelemetry env inputs are valid when tracing is enabled

### `ollama` is healthy but model calls fail

- verify GGUF host path mounts
- verify Modelfile render values
- verify the model exists through `/api/tags`

### `edge-toolbox` is present but scripts are missing

- rebuild/import the local `edge-toolbox` image
- restart the toolbox deployment

## 10. Cleanup and Hygiene

Do not commit:

- `values.local-k3s.yaml`
- notebook checkpoint directories
- generated notebook assets
- rendered manifests
- secret material
- large model binaries

Before publishing:

```bash
git status --short
helm lint .
go test ./...
```
