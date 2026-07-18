# Xubuntu 24 + k3s + NVIDIA GPU Runbook

This runbook is the local operating procedure for `llm-observability-stack` on the current Xubuntu 24 single-node k3s system with an NVIDIA GPU. It replaces older command notes that used stale project paths and outdated helper names.

Current verified state on July 2, 2026:

| Item | Value |
|---|---|
| Node | `waqasm86-thinkpad-t450s` |
| Node status | `Ready` |
| Kubernetes | `v1.36.2+k3s1` |
| Runtime | `containerd://2.3.2-k3s2` |
| OS image | `Ubuntu 24.04.3 LTS` |
| Kernel | `6.17.0-35-generic` |
| GPU resource | `nvidia.com/gpu: 1` |
| NVIDIA RuntimeClass | `nvidia` |
| GPU Operator release | `gpu-operator` in namespace `gpu-operator` |
| Stack release | `llm-observability-stack` in namespace `llm-observability` |

## Project Path

```bash
cd /media/waqasm86/External1/Waqas-Projects/Project-Edge-Computing-LLM/llm-observability-stack
```

## Check k3s and NVIDIA GPU State

```bash
kubectl get nodes -o wide
kubectl get pods -A
kubectl get runtimeclass nvidia
kubectl describe node waqasm86-thinkpad-t450s | grep -A8 -B5 nvidia.com/gpu
kubectl get nodes -o jsonpath='{range .items[*]}{.metadata.name}{" gpu="}{.status.allocatable.nvidia\.com/gpu}{"\n"}{end}'
```

GPU Operator and device-plugin checks:

```bash
helm list -n gpu-operator
kubectl get pods -n gpu-operator -o wide
kubectl get pods -A | grep -E 'gpu|nvidia|dcgm'
```

## Check the Running llm-observability-stack Release

```bash
helm list -n llm-observability
kubectl get pods,svc,pvc -n llm-observability -o wide
kubectl get servicemonitors,probes,prometheusrules -A
```

Expected core workloads in the current full local deployment:

- `ollama`
- `open-webui`
- `open-webui-redis`
- `ollama-gateway`
- `edge-toolbox`
- `opentelemetry-collector`
- `dcgm-exporter`
- `kube-prometheus-stack-operator`
- `prometheus-kube-prometheus-stack-prometheus`
- `alertmanager-kube-prometheus-stack-alertmanager`
- `llm-observability-stack-grafana`
- `llm-observability-stack-kube-state-metrics`
- `llm-observability-stack-prometheus-node-exporter`

## Build and Import Local Images

Use this when `ollama-gateway` or `edge-toolbox` source code changes.

```bash
docker build -t ollama-gateway:0.2.0 ./ollama-gateway
docker build -t edge-toolbox:0.2.0 ./edge-toolbox

./hack/import-local-image-to-k3s.sh ollama-gateway 0.2.0
./hack/import-local-image-to-k3s.sh edge-toolbox 0.2.0

sudo k3s ctr images list | grep -E 'ollama-gateway|edge-toolbox'
```

## Minimal GPU/Ollama Profile

The GeForce 940M has 1 GiB VRAM, so the minimal profile is the safest profile for inference testing.

Use the real local GGUF model directory:

```bash
MODEL_DIR=/media/waqasm86/External1/Waqas-Projects/repos-llamatelemetry/llamatelemetry-xubuntu24/models

helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.geforce-940m-k3s.yaml \
  --set ollamaModel.gguf.hostPath="$MODEL_DIR" \
  --set ollama.volumes[0].hostPath.path="$MODEL_DIR" \
  --set-json 'ollama.nodeSelector={"nvidia.com/gpu.present":"true"}'
```

Verify the minimal profile:

```bash
kubectl rollout status deploy/ollama -n llm-observability
kubectl rollout status deploy/opentelemetry-collector -n llm-observability
kubectl exec -n llm-observability deploy/ollama -- ollama list
./hack/test-geforce-940m-inference.sh
```

## Full Local k3s/NVIDIA Profile

This profile runs Ollama, Open WebUI, Ollama gateway, Go edge toolbox,
OpenTelemetry Collector, Prometheus, Grafana, Alertmanager, node exporter, and
kube-state-metrics. It may observe the DCGM exporter from the base layer through
ServiceMonitor resources, but it does not deploy DCGM exporter.

```bash
MODEL_DIR=/media/waqasm86/External1/Waqas-Projects/repos-llamatelemetry/llamatelemetry-xubuntu24/models

./hack/bootstrap-enterprise-pilot-k3s.sh \
  --set ollamaModel.gguf.hostPath="$MODEL_DIR" \
  --set ollama.volumes[0].hostPath.path="$MODEL_DIR" \
  --set-json 'ollama.nodeSelector={"nvidia.com/gpu.present":"true"}'
```

Manual equivalent:

```bash
helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.enterprise-pilot-k3s.yaml \
  --set kube-prometheus-stack.crds.enabled=false \
  --set ollamaModel.gguf.hostPath="$MODEL_DIR" \
  --set ollama.volumes[0].hostPath.path="$MODEL_DIR" \
  --set-json 'ollama.nodeSelector={"nvidia.com/gpu.present":"true"}'
```

## CPU-Only Fallback

Use CPU mode when NVIDIA GPU runtime is missing or when validating chart behavior on systems without NVIDIA hardware.

```bash
./hack/detect-runtime-profile.sh --mode cpu

helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.cpu-k3s.yaml \
  --set kube-prometheus-stack.crds.enabled=false
```

Render-only CPU validation:

```bash
helm template llm-observability-stack . \
  -f values.cpu-k3s.yaml \
  >/tmp/rendered-cpu.yaml
```

## Port-Forward Services

Run each command in its own terminal when you need local access.

```bash
kubectl port-forward -n llm-observability svc/ollama 11434:11434
kubectl port-forward -n llm-observability svc/ollama-gateway 8000:8000
kubectl port-forward -n llm-observability svc/open-webui 8080:8080
kubectl port-forward -n llm-observability svc/llm-observability-stack-grafana 3000:80
kubectl port-forward -n llm-observability svc/opentelemetry-collector 4317:4317 4318:4318 8888:8888
```

Useful local URLs:

- Ollama: `http://localhost:11434`
- Ollama gateway: `http://localhost:8000`
- Open WebUI: `http://localhost:8080`
- Grafana: `http://localhost:3000`
- OpenTelemetry Collector metrics: `http://localhost:8888/metrics`

## Service Tests

Ollama tags:

```bash
curl -s http://127.0.0.1:11434/api/tags | jq
```

Ollama generation:

```bash
curl -s http://127.0.0.1:11434/api/generate \
  -H 'Content-Type: application/json' \
  -d '{"model":"qwen-1-8b-chat-q4-k-m-local:latest","prompt":"Reply with one short sentence.","stream":false}' | jq
```

Ollama gateway:

```bash
curl -s http://127.0.0.1:8000/healthz | jq
curl -s http://127.0.0.1:8000/config | jq
curl -s http://127.0.0.1:8000/ollama/api/tags | jq
```

OpenTelemetry Collector metrics:

```bash
curl -s http://127.0.0.1:8888/metrics | head
```

## Repository Validation

```bash
helm lint .
helm template llm-observability-stack . >/tmp/rendered-default.yaml
helm template llm-observability-stack . -f values.geforce-940m-k3s.yaml >/tmp/rendered-geforce.yaml
helm template llm-observability-stack . -f values.cpu-k3s.yaml >/tmp/rendered-cpu.yaml
helm template llm-observability-stack . \
  -f values.full-stack-nvidia.example.yaml \
  --set opentelemetry.tracing.enabled= \
  --set openWebUI.existingSecret= \
  --set open-webui.webuiSecret.existingSecretName= \
  >/tmp/rendered-full-stack-nvidia.yaml

go test ./...
./hack/validate-local-stack.sh --strict-gpu
```

## Benchmark

Keep benchmark outputs under `artifacts/` and commit only sanitized evidence.

```bash
bin/llm-observability benchmark \
  --model qwen-1-8b-chat-q4-k-m-local:latest \
  --warmup-runs 1 \
  --runs 3 \
  --output artifacts/local-benchmark.json
```

## Go Edge Toolbox Checks

Run these after the full local profile is running:

```bash
kubectl exec -it -n llm-observability deploy/edge-toolbox -- edge-toolbox dns ollama ollama-gateway open-webui open-webui-redis
kubectl exec -it -n llm-observability deploy/edge-toolbox -- edge-toolbox ollama-smoke
kubectl exec -it -n llm-observability deploy/edge-toolbox -- edge-toolbox redis-ping
```

## Capture Local Evidence

```bash
./hack/capture-local-evidence.sh artifacts/evidence-local-run
```

Review every generated file before sharing or committing. The helper does not collect Secrets or ConfigMaps, but local paths, node names, and cluster metadata may still be sensitive.

## Clean Up

Uninstall only the application stack:

```bash
helm uninstall llm-observability-stack -n llm-observability
kubectl delete namespace llm-observability
```

Keep the GPU Operator installed if other local GPU workloads depend on it. Remove it only when resetting the whole k3s/NVIDIA environment:

```bash
helm uninstall gpu-operator -n gpu-operator
kubectl delete namespace gpu-operator
```

## Troubleshooting

Pending pods:

```bash
kubectl get pods -n llm-observability -o wide
kubectl describe pod -n llm-observability <pod-name>
```

GPU scheduling:

```bash
kubectl get runtimeclass nvidia
kubectl describe node waqasm86-thinkpad-t450s | grep -A8 -B5 nvidia.com/gpu
kubectl logs -n gpu-operator -l app=nvidia-operator-validator --tail=200
```

Ollama model path:

```bash
kubectl describe pod -n llm-observability -l app.kubernetes.io/name=ollama
kubectl exec -n llm-observability deploy/ollama -- ls -lah /models/gguf
kubectl exec -n llm-observability deploy/ollama -- ollama list
```

Open WebUI startup:

```bash
kubectl logs -n llm-observability statefulset/open-webui --tail=200
```

Open WebUI should not try to download embedding, reranking, version-check, or Whisper assets during local startup. The values files set those auto-update paths off for local reliability.
