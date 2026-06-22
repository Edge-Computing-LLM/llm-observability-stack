# Deployment Proof

This document describes how to collect deployment evidence for `llm-observability-stack` before submitting the NVIDIA Inception 2026 application.

Use the exact machine, namespace, values file, and model path used for the final evidence run. Mark any command output that contains secrets as private and do not commit it.

## Repository and Environment

```bash
cd /media/waqasm86/External11/Project-Llamatelemetry/Project-Llamatelemetry-End-to-End/Project-Nvidia-Singapore-Competition-2026/llm-observability-stack

git rev-parse HEAD
git status --short
helm version
kubectl version --client
kubectl cluster-info
kubectl get nodes -o wide
kubectl get storageclass
nvidia-smi
```

## Validate the Chart

```bash
helm dependency update .
helm lint .
helm template llm-observability-stack . \
  -f values.enterprise-pilot-k3s.yaml \
  --set kube-prometheus-stack.crds.enabled=false \
  >/tmp/llm-observability-enterprise.yaml

grep -n 'opentelemetry-collector' /tmp/llm-observability-enterprise.yaml | head
grep -n 'kind: PrometheusRule' /tmp/llm-observability-enterprise.yaml
grep -n 'kind: ServiceMonitor' /tmp/llm-observability-enterprise.yaml | head
grep -n 'helm.sh/resource-policy: keep' /tmp/llm-observability-enterprise.yaml
```

For the narrow GeForce 940M Ollama-only proof:

```bash
helm template llm-observability-stack . \
  -f values.geforce-940m-k3s.yaml \
  >/tmp/llm-observability-geforce.yaml
```

## Prepare NVIDIA Runtime and Device Plugin

```bash
kubectl get runtimeclass
kubectl get runtimeclass nvidia
kubectl get pods -A | grep -Ei 'nvidia|device-plugin|dcgm'
kubectl get nodes -o jsonpath='{range .items[*]}{.metadata.name}{" gpu="}{.status.allocatable.nvidia\.com/gpu}{"\n"}{end}'
```

If the NVIDIA device plugin is not already installed in the target k3s cluster, the repository provides:

```bash
./hack/install-nvidia-device-plugin.sh
```

For the verified single-node profile:

```bash
./hack/prepare-single-node-k3s.sh
```

Verify before submission: only run these helper scripts on a cluster where changing node labels and installing the NVIDIA device plugin is expected.

## Build and Import Local Images

`langchain-demo` and `python-toolbox` are local images. Build and import them into k3s containerd before enabling those workloads:

```bash
./hack/build-local-image.sh langchain-demo 0.1.1 ./langchain-demo
./hack/import-local-image-to-k3s.sh langchain-demo 0.1.1

./hack/build-local-image.sh python-toolbox 0.2.0 ./python-toolbox
./hack/import-local-image-to-k3s.sh python-toolbox 0.2.0

sudo k3s ctr images ls | grep -E 'langchain-demo|python-toolbox'
```

## Deploy the Stack

Preferred local enterprise-pilot deployment:

```bash
./hack/bootstrap-enterprise-pilot-k3s.sh
```

Equivalent Helm command after CRDs and local images are ready:

```bash
helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.enterprise-pilot-k3s.yaml \
  --set kube-prometheus-stack.crds.enabled=false
```

Lighter first install if local images are not imported yet:

```bash
./hack/bootstrap-enterprise-pilot-k3s.sh \
  --set langchainDemo.enabled=false \
  --set pythonToolbox.enabled=false
```

GeForce 940M Ollama-only proof:

```bash
./hack/prepare-single-node-k3s.sh
./hack/install-nvidia-device-plugin.sh

helm upgrade --install llm-observability-stack . \
  -n llm-observability --create-namespace \
  -f values.geforce-940m-k3s.yaml

./hack/test-geforce-940m-inference.sh
```

## Verify Pods, Services, and Storage

```bash
NS=llm-observability

helm list -n "$NS"
helm get values llm-observability-stack -n "$NS" -a
kubectl get ns "$NS"
kubectl get all -n "$NS" -o wide
kubectl get pods,svc,pvc,configmap -n "$NS" -o wide
kubectl get events -n "$NS" --sort-by=.lastTimestamp | tail -n 50
kubectl get deploy,statefulset -n "$NS"
kubectl rollout status deploy/ollama -n "$NS" --timeout=300s
kubectl rollout status deploy/opentelemetry-collector -n "$NS" --timeout=180s
kubectl rollout status deploy/langchain-demo -n "$NS" --timeout=180s
kubectl rollout status statefulset/open-webui -n "$NS" --timeout=180s
```

If `pythonToolbox.enabled=true`:

```bash
kubectl rollout status deploy/python-toolbox -n "$NS" --timeout=180s
```

## Verify GPU Visibility

```bash
nvidia-smi
kubectl get runtimeclass nvidia
kubectl get nodes -o custom-columns=NAME:.metadata.name,GPU_CAPACITY:.status.capacity.nvidia\\.com/gpu,GPU_ALLOCATABLE:.status.allocatable.nvidia\\.com/gpu
kubectl describe deploy/ollama -n "$NS"
kubectl logs -n "$NS" deploy/ollama --tail=200 | grep -Ei 'cuda|nvidia|gpu|vram|compute' || true
```

For the verified GeForce profile:

```bash
./hack/test-geforce-940m-inference.sh
```

## Verify Prometheus Targets and Monitoring Resources

```bash
kubectl get servicemonitors,probes,prometheusrules -n "$NS"
kubectl get servicemonitors,probes,prometheusrules -A
kubectl describe servicemonitor langchain-demo -n "$NS"
kubectl describe prometheusrule llm-observability-stack -n "$NS"
kubectl get pods -n "$NS" | grep -Ei 'prometheus|grafana|alertmanager|operator'
```

Prometheus UI access:

```bash
kubectl port-forward -n "$NS" svc/llm-observability-stack-prometheus 9090:9090
```

Verify before submission: service names can vary if Helm release names or fullname overrides change. If the command fails, run:

```bash
kubectl get svc -n "$NS" | grep -Ei 'prometheus|grafana|alertmanager'
```

Prometheus queries to capture:

```text
up
llm_observability_llm_requests_total
llm_observability_time_to_first_token_seconds_count
llm_observability_generated_tokens_per_second_bucket
DCGM_FI_DEV_GPU_UTIL
DCGM_FI_DEV_FB_USED
```

`DCGM_FI_*` metrics require DCGM exporter or GPU Operator/DCGM to be available and scraped.

## Verify Grafana Dashboards

```bash
kubectl port-forward -n "$NS" svc/llm-observability-stack-grafana 3000:80
```

Open:

```text
http://127.0.0.1:3000
username: admin
password: admin
```

Dashboards expected from repository:

- `LLM Observability Overview`
- `NVIDIA GPU Observability`
- `LLM Competition Benchmark`

If dashboard discovery fails, verify the dashboard ConfigMap:

```bash
kubectl get configmap -n "$NS" | grep -i dashboard
kubectl describe configmap llm-observability-dashboards -n "$NS"
```

## Verify OpenTelemetry Collector

```bash
kubectl get deploy,svc,configmap -n "$NS" | grep opentelemetry-collector
kubectl logs -n "$NS" deploy/opentelemetry-collector --tail=100
kubectl port-forward -n "$NS" svc/opentelemetry-collector 4317:4317 4318:4318 8888:8888
curl -s http://127.0.0.1:8888/metrics | head
```

## Run an Inference Test

Port-forward Ollama:

```bash
kubectl port-forward -n "$NS" svc/ollama 11434:11434
```

From another terminal:

```bash
curl -s http://127.0.0.1:11434/api/tags | jq

curl -s http://127.0.0.1:11434/api/generate \
  -H 'Content-Type: application/json' \
  -d '{"model":"gemma3-1b-it-gguf-local","prompt":"Reply with one short sentence for deployment proof.","stream":false}' | jq
```

Proxy path:

```bash
kubectl port-forward -n "$NS" svc/langchain-demo 8000:8000

curl -s http://127.0.0.1:8000/healthz | jq
curl -s http://127.0.0.1:8000/config | jq
curl -s http://127.0.0.1:8000/invoke \
  -H 'Content-Type: application/json' \
  -d '{"prompt":"Say the stack is healthy in one line."}' | jq
curl -s http://127.0.0.1:8000/metrics | grep 'llm_observability' | head
```

Open WebUI:

```bash
kubectl port-forward -n "$NS" svc/open-webui 8080:8080
```

Open `http://127.0.0.1:8080` and send a short prompt. Capture a screenshot only if it contains no secrets or private prompt data.

## Run a Benchmark

With Ollama port-forward active:

```bash
./benchmarks/ollama_benchmark.py \
  --model gemma3-1b-it-gguf-local \
  --warmup-runs 1 \
  --runs 10 \
  --output validation/benchmark-results/benchmark-$(date -u +%Y%m%dT%H%M%SZ).json
```

For a direct comparison with the existing sanitized artifact, see:

```text
artifacts/geforce-940m-benchmark.json
```

If a Pushgateway is configured, the benchmark script supports:

```bash
./benchmarks/ollama_benchmark.py \
  --model gemma3-1b-it-gguf-local \
  --warmup-runs 1 \
  --runs 10 \
  --pushgateway http://<pushgateway-host>:9091 \
  --output validation/benchmark-results/benchmark-pushgateway-$(date -u +%Y%m%dT%H%M%SZ).json
```

Verify before submission: only use Pushgateway commands if a Pushgateway exists in the deployment.

## Capture Logs

```bash
mkdir -p validation/benchmark-results/logs

kubectl logs -n "$NS" deploy/ollama --tail=300 > validation/benchmark-results/logs/ollama-$(date -u +%Y%m%dT%H%M%SZ).log
kubectl logs -n "$NS" deploy/langchain-demo --tail=300 > validation/benchmark-results/logs/langchain-demo-$(date -u +%Y%m%dT%H%M%SZ).log
kubectl logs -n "$NS" deploy/opentelemetry-collector --tail=300 > validation/benchmark-results/logs/opentelemetry-collector-$(date -u +%Y%m%dT%H%M%SZ).log
kubectl logs -n "$NS" statefulset/open-webui --tail=300 > validation/benchmark-results/logs/open-webui-$(date -u +%Y%m%dT%H%M%SZ).log
```

Review logs before sharing. Remove tokens, prompts, private hostnames, or customer data.

## Screenshot Checklist

Place sanitized screenshots in `validation/screenshots/`:

- `01-kubectl-pods-services-pvc.png`: terminal showing pods, services, PVCs in `llm-observability`.
- `02-gpu-runtime.png`: terminal showing `nvidia-smi`, `RuntimeClass`, and `nvidia.com/gpu` allocatable.
- `03-ollama-inference.png`: terminal showing `/api/tags` and successful inference output.
- `04-langchain-demo-metrics.png`: terminal/browser showing `/healthz` and `/metrics`.
- `05-prometheus-targets.png`: Prometheus targets or query results.
- `06-grafana-llm-overview.png`: LLM Observability Overview dashboard.
- `07-grafana-nvidia-gpu.png`: NVIDIA GPU Observability dashboard.
- `08-grafana-benchmark.png`: LLM Competition Benchmark dashboard.
- `09-open-webui-chat.png`: Open WebUI chat path with non-sensitive demo prompt.
- `10-benchmark-json.png`: terminal or file view showing benchmark summary.

## Evidence Table

Fill this table for the final NVIDIA Inception submission.

| Evidence item | Value / link / screenshot | Notes |
|---|---|---|
| Evidence date UTC | `TBD` | Fill manually. |
| Git commit SHA | `TBD` | `git rev-parse HEAD`. |
| Host machine | `TBD` | Example: Lenovo ThinkPad T450s. |
| OS | `TBD` | Example: Xubuntu 24.04. |
| Kubernetes distribution/version | `TBD` | `kubectl version` and k3s version. |
| Helm version | `TBD` | `helm version`. |
| NVIDIA GPU SKU/count | `TBD` | `nvidia-smi`. |
| Driver/CUDA visibility | `TBD` | `nvidia-smi`; pod logs. |
| RuntimeClass | `TBD` | `kubectl get runtimeclass nvidia`. |
| GPU allocatable | `TBD` | `kubectl get nodes ... nvidia.com/gpu`. |
| Values profile | `TBD` | Example: `values.enterprise-pilot-k3s.yaml`. |
| Model | `TBD` | Example: `gemma3-1b-it-gguf-local`. |
| Inference result | `TBD` | Link screenshot/log. |
| TTFT p50/p95 | `TBD` | Benchmark JSON. |
| Latency p50/p95 | `TBD` | Benchmark JSON. |
| Tokens per second | `TBD` | Benchmark JSON. |
| GPU utilization | `TBD` | DCGM/Grafana/nvidia-smi. |
| GPU memory | `TBD` | DCGM/Grafana/nvidia-smi. |
| Prometheus targets screenshot | `TBD` | Place in `validation/screenshots/`. |
| Grafana LLM dashboard screenshot | `TBD` | Place in `validation/screenshots/`. |
| Grafana GPU dashboard screenshot | `TBD` | Place in `validation/screenshots/`. |
| Open WebUI screenshot | `TBD` | Optional, non-sensitive. |
| Logs captured | `TBD` | Sanitized only. |
| External reviewer/customer feedback | `TBD` | Use `redacted-feedback-template.md`. |
