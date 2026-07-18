# Quick Start

This guide is the fastest path to a working local `llm-observability-stack` deployment on a single-node k3s machine.

## 1. Prerequisites

- k3s running and reachable from `kubectl`
- `k3s-nvidia-edge` already deployed and healthy
- NVIDIA runtime configured on the node
- NVIDIA GPU allocatable in Kubernetes
- Helm 3 installed
- Docker or `nerdctl`
- local GGUF model file on host storage

Quick checks:

```bash
kubectl get nodes -o wide
kubectl get pods -n gpu-operator
kubectl describe node "$(kubectl get nodes -o jsonpath='{.items[0].metadata.name}')" | grep -A8 Allocatable
helm version
```

Read [K3S-NVIDIA-EDGE-DEPENDENCY.md](K3S-NVIDIA-EDGE-DEPENDENCY.md) if the GPU Operator, RuntimeClass, or `nvidia.com/gpu` checks are not already passing.

## 2. Prepare Local Values

For the verified local single-node k3s/NVIDIA workflow layered on top of `k3s-nvidia-edge`, use
`values.geforce-940m-k3s.yaml`. This profile deploys Ollama, Open WebUI, Open WebUI Redis, and the
OpenTelemetry collector; it does not deploy GPU Operator, NVIDIA device plugin, or DCGM exporter
workloads because those are owned by `k3s-nvidia-edge`.

Create a private override only when your host paths, secrets, or service exposure differ:

```bash
cp values.local-k3s.example.yaml values.local-k3s.yaml
```

Edit `values.local-k3s.yaml` and confirm:

- GGUF host path values point to your model directory
- the OpenTelemetry OTLP endpoint points at the in-cluster collector
- Open WebUI secret inputs are wired the way you want
- the example profile defaults still match the behavior you want locally

Do not commit `values.local-k3s.yaml`.

Profile reference:

- [CONFIG-PROFILES.md](CONFIG-PROFILES.md)

## 3. Build and Import Local Images

The GeForce 940M profile does not require local app images because Open WebUI connects directly to
Ollama. Build these images only when you intentionally enable `ollamaGateway` or `edgeToolbox` in
another profile:

```bash
./hack/build-local-image.sh ollama-gateway 0.2.0 . ollama-gateway/Dockerfile
./hack/build-local-image.sh edge-toolbox 0.2.0 . edge-toolbox/Dockerfile
```

Import them into k3s:

```bash
./hack/import-local-image-to-k3s.sh ollama-gateway 0.2.0
./hack/import-local-image-to-k3s.sh edge-toolbox 0.2.0
```

Verify:

```bash
sudo k3s ctr images ls | grep -E 'ollama-gateway|edge-toolbox'
```

## 4. Deploy the Chart

Verified local k3s/NVIDIA command:

```bash
helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.geforce-940m-k3s.yaml \
  --wait
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
- `open-webui-redis` internal `ClusterIP`
- `opentelemetry-collector` internal `ClusterIP`

## 6. Local Access

Browser:

- Open WebUI: `http://localhost:8080/`

Expose UIs and internal APIs from separate terminals:

```bash
kubectl port-forward -n llm-observability svc/open-webui 8080:8080
kubectl port-forward -n llm-observability svc/ollama 11434:11434
kubectl port-forward -n llm-observability svc/opentelemetry-collector 8888:8888 4317:4317 4318:4318
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
kubectl logs -n llm-observability deploy/ollama --tail=100
kubectl logs -n llm-observability statefulset/open-webui --tail=100
kubectl logs -n llm-observability deploy/open-webui-redis --tail=100
```

If notebook API cells fail:

- verify the required port-forwards are active
- verify `edgeToolbox.enabled` and toolbox pod health
- verify `OTEL_EXPORTER_OTLP_ENDPOINT` and `OTEL_SERVICE_NAME` for tracing notebooks

## 9. Next Reading

- [ARCHITECTURE.md](ARCHITECTURE.md)
- [OPERATIONS-RUNBOOK.md](OPERATIONS-RUNBOOK.md)
- [../jupyter-notebooks/README.md](../jupyter-notebooks/README.md)
