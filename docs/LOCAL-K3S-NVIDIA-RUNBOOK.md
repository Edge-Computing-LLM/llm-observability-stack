# Local k3s NVIDIA Runbook

This runbook deploys `llm-observability-stack` on a local Xubuntu 24 host with k3s, NVIDIA GPU support, local GGUF storage, Ollama, Open WebUI, Prometheus/Grafana, DCGM exporter, and the vendored OpenTelemetry Collector chart.

## 1. Enter the Project

```bash
cd /media/waqasm86/External1/Waqas-Projects/Project-Edge-Computing-LLM/llm-observability-stack
```

## 2. Verify Host Prerequisites

```bash
nvidia-smi
helm version
kubectl version --client
kubectl get nodes -o wide
kubectl get storageclass
```

For k3s, `local-path` storage should exist. The GPU must be visible on the host before Kubernetes can schedule Ollama with `nvidia.com/gpu`.

## 3. Verify NVIDIA Runtime and Device Plugin

```bash
kubectl get runtimeclass
kubectl get pods -A | grep -Ei 'nvidia|device-plugin|dcgm'
kubectl get nodes -o jsonpath='{range .items[*]}{.metadata.name}{" gpu="}{.status.allocatable.nvidia\.com/gpu}{"\n"}{end}'
```

If the device plugin is not installed, install it using the local helper:

```bash
./hack/install-nvidia-device-plugin.sh
```

## 4. Verify the Local GGUF Model

The default local enterprise profile expects:

```bash
ls -lh /media/waqasm86/External1/Waqas-Projects/repos-llamatelemetry/llamatelemetry-xubuntu24/models/gemma-3-1b-it-Q4_K_M.gguf
```

The chart mounts that host directory read-only into Ollama at `/models/gguf`, and the Ollama PVC is annotated with `helm.sh/resource-policy: keep`.

## 5. Build and Import Local Images

The `langchain-demo` and `python-toolbox` images are local project images. Build and import them into k3s containerd:

```bash
./hack/build-local-image.sh langchain-demo 0.1.1 ./langchain-demo
./hack/import-local-image-to-k3s.sh langchain-demo 0.1.1

./hack/build-local-image.sh python-toolbox 0.2.0 ./python-toolbox
./hack/import-local-image-to-k3s.sh python-toolbox 0.2.0
```

Verify:

```bash
sudo k3s ctr images ls | grep -E 'langchain-demo|python-toolbox'
```

The images that must be built locally before enabling `langchainDemo` and `pythonToolbox` are:

```text
langchain-demo:0.1.1
python-toolbox:0.2.0
```

The platform images rendered by the enterprise profile are:

```text
docker.io/grafana/grafana:13.0.2
docker.io/library/busybox:1.38.0
ghcr.io/jkroepke/kube-webhook-certgen:1.8.3
ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-k8s:0.153.0
ghcr.io/open-webui/open-webui:0.8.10
nvcr.io/nvidia/k8s/dcgm-exporter:4.5.3-4.8.2-distroless
ollama/ollama:0.17.7
quay.io/kiwigrid/k8s-sidecar:2.7.3
quay.io/prometheus-operator/prometheus-operator:v0.91.0
quay.io/prometheus/alertmanager:v0.33.0
quay.io/prometheus/node-exporter:v1.11.1-distroless
quay.io/prometheus/prometheus:v3.12.0-distroless
redis:7.4.2-alpine3.21
registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.19.1
```

## 6. Validate Helm Rendering

```bash
helm dependency update .
helm lint .
helm template llm-observability-stack . -f values.enterprise-pilot-k3s.yaml >/tmp/llm-observability-enterprise.yaml
grep -n 'opentelemetry-collector' /tmp/llm-observability-enterprise.yaml | head
grep -n 'helm.sh/resource-policy: keep' /tmp/llm-observability-enterprise.yaml
```

## 7. Install the Stack

Use the bootstrap helper to apply Prometheus Operator CRDs first, then install the umbrella chart:

```bash
./hack/bootstrap-enterprise-pilot-k3s.sh
```

For a lighter first install without local app images:

```bash
./hack/bootstrap-enterprise-pilot-k3s.sh \
  --set langchainDemo.enabled=false \
  --set pythonToolbox.enabled=false
```

After importing local images, enable the workloads with:

```bash
helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.enterprise-pilot-k3s.yaml \
  --set kube-prometheus-stack.crds.enabled=false
```

## 8. Watch Rollout

```bash
kubectl get pods,svc,pvc -n llm-observability -o wide
kubectl rollout status deploy/ollama -n llm-observability --timeout=300s
kubectl rollout status deploy/opentelemetry-collector -n llm-observability --timeout=180s
kubectl rollout status deploy/langchain-demo -n llm-observability --timeout=180s
```

## 9. Verify Ollama and the Local Model

```bash
kubectl exec -n llm-observability deploy/ollama -- ls -lh /models/gguf
kubectl exec -n llm-observability deploy/ollama -- ollama list
```

Port-forward Ollama:

```bash
kubectl port-forward -n llm-observability svc/ollama 11434:11434
```

From another terminal:

```bash
curl -s http://127.0.0.1:11434/api/tags | jq
curl -s http://127.0.0.1:11434/api/generate \
  -H 'Content-Type: application/json' \
  -d '{"model":"gemma3-1b-it-gguf-local","prompt":"Reply with one short sentence.","stream":false}' | jq
```

## 10. Verify OpenTelemetry Collector

```bash
kubectl get deploy,svc,configmap -n llm-observability | grep opentelemetry-collector
kubectl logs -n llm-observability deploy/opentelemetry-collector --tail=100
kubectl port-forward -n llm-observability svc/opentelemetry-collector 4317:4317 4318:4318 8888:8888
curl -s http://127.0.0.1:8888/metrics | head
```

`langchain-demo` sends OTLP to `http://opentelemetry-collector:4317` when OpenTelemetry is enabled.

## 11. Access Open WebUI and Grafana

```bash
kubectl port-forward -n llm-observability svc/open-webui 8080:8080
kubectl port-forward -n llm-observability svc/llm-observability-stack-grafana 3000:80
```

Open:

```text
http://127.0.0.1:8080
http://127.0.0.1:3000
```

Grafana local credentials are pinned by this chart:

```text
username: admin
password: admin
```

## 12. Run Tests

```bash
pytest -q tests
./hack/competition-validate.sh
./hack/competition-validate.sh --strict-gpu
```

## 13. Uninstall Without Deleting Models

```bash
helm uninstall llm-observability-stack -n llm-observability
kubectl get pvc -n llm-observability
```

Do not delete the Ollama PVC or the host GGUF directory unless you intentionally want to remove model state. The GGUF hostPath mount is read-only, and Helm keeps the Ollama PVC.
