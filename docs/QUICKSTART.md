# Quick Start

This guide is the fastest path to a working local `llm-observability-stack` deployment on a single-node k3s machine.

## 1. Prerequisites

- k3s running and reachable from `kubectl`
- NVIDIA runtime configured on the node
- NVIDIA device plugin already healthy in-cluster
- Helm 3 installed
- Docker or `nerdctl`
- local GGUF model file on host storage

Quick checks:

```bash
kubectl get nodes -o wide
kubectl get pods -n nvidia-device-plugin
helm version
```

## 2. Prepare Local Values

For the verified local single-node k3s/NVIDIA workflow, use `values.enterprise-pilot-k3s.yaml`.

Create a private override only when your host paths, secrets, or service exposure differ:

```bash
cp values.local-k3s.example.yaml values.local-k3s.yaml
```

Edit `values.local-k3s.yaml` and confirm:

- GGUF host path values point to your model directory
- LangSmith credentials or existing secret references are set correctly
- Open WebUI secret inputs are wired the way you want
- the example profile defaults still match the behavior you want locally

Do not commit `values.local-k3s.yaml`.

Profile reference:

- [CONFIG-PROFILES.md](CONFIG-PROFILES.md)

## 3. Build and Import Local Images

Build the local images:

```bash
./hack/build-local-image.sh langchain-demo 0.1.1 ./langchain-demo
./hack/build-local-image.sh python-toolbox 0.2.0 ./python-toolbox
```

Import them into k3s:

```bash
./hack/import-local-image-to-k3s.sh langchain-demo 0.1.1
./hack/import-local-image-to-k3s.sh python-toolbox 0.2.0
```

Verify:

```bash
sudo k3s ctr images ls | grep -E 'langchain-demo|python-toolbox'
```

## 4. Deploy the Chart

Verified local k3s/NVIDIA command:

```bash
helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.enterprise-pilot-k3s.yaml \
  --set kube-prometheus-stack.crds.enabled=false
```

If you created a private `values.local-k3s.yaml`, keep the Ollama PVC size aligned with any existing `ollama` PVC before using it on the same release. The k3s `local-path` provisioner does not support resizing this PVC in place.

## 5. Verify the Deployment

```bash
kubectl get all -n llm-observability
kubectl get svc -n llm-observability
kubectl get pvc -n llm-observability
```

Typical local result:

- `open-webui` available through port-forwarding
- `ollama` internal `ClusterIP`
- `langchain-demo` internal `ClusterIP`
- `python-toolbox` running for in-cluster diagnostics

## 6. Local Access

Browser:

- Open WebUI: `http://localhost:8080/`

Expose UIs and internal APIs from separate terminals:

```bash
kubectl port-forward -n llm-observability svc/open-webui 8080:8080
kubectl port-forward -n llm-observability svc/llm-observability-stack-grafana 3000:80
kubectl port-forward -n llm-observability svc/ollama 11434:11434
kubectl port-forward -n llm-observability svc/langchain-demo 8000:8000
```

## 7. Jupyter Notebooks

Launch notebooks from the project notebook directory:

```bash
PYTHON_BIN="${PYTHON_BIN:-python3.11}"
cd jupyter-notebooks
"${PYTHON_BIN}" -m jupyter lab
```

Notebook index:

- [../jupyter-notebooks/CATALOG.md](../jupyter-notebooks/CATALOG.md)
- [../jupyter-notebooks/README.md](../jupyter-notebooks/README.md)

## 8. Minimal Troubleshooting

If the stack does not come up cleanly:

```bash
kubectl get pods -n llm-observability -o wide
kubectl logs -n llm-observability deploy/langchain-demo --tail=100
kubectl logs -n llm-observability deploy/ollama --tail=100
kubectl logs -n llm-observability statefulset/open-webui --tail=100
```

If notebook API cells fail:

- verify the required port-forwards are active
- verify `pythonToolbox.enabled` and toolbox pod health
- verify LangSmith environment variables for tracing notebooks

## 9. Next Reading

- [ARCHITECTURE.md](ARCHITECTURE.md)
- [OPERATIONS-RUNBOOK.md](OPERATIONS-RUNBOOK.md)
- [../jupyter-notebooks/README.md](../jupyter-notebooks/README.md)
