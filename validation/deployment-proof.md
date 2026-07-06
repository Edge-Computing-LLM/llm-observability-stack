# Local Deployment Proof

This document describes how to collect local deployment evidence for `llm-observability-stack` on k3s. Use it for internal review, documentation updates, and reproducible local validation.

## Working Directory

```bash
cd /media/waqasm86/External1/Waqas-Projects/Project-Edge-Computing-LLM/llm-observability-stack
```

## Cluster and GPU Evidence

```bash
kubectl get nodes -o wide
kubectl get runtimeclass nvidia
kubectl get nodes -o jsonpath='{range .items[*]}{.metadata.name}{" gpu="}{.status.allocatable.nvidia\.com/gpu}{"\n"}{end}'
kubectl describe node waqasm86-thinkpad-t450s | grep -A8 -B5 nvidia.com/gpu
helm list -A
```

## Stack Evidence

```bash
kubectl get pods,svc,pvc -n llm-observability -o wide
kubectl get servicemonitors,probes,prometheusrules -A
kubectl logs -n llm-observability deployment/ollama --tail=100
kubectl logs -n llm-observability deployment/langchain-demo --tail=100
kubectl logs -n llm-observability statefulset/open-webui --tail=100
```

## Render Evidence

```bash
helm lint .
helm template llm-observability-stack . >/tmp/rendered-default.yaml
helm template llm-observability-stack . -f values.cpu-k3s.yaml >/tmp/rendered-cpu.yaml
helm template llm-observability-stack . -f values.geforce-940m-k3s.yaml >/tmp/rendered-geforce.yaml
helm template llm-observability-stack . \
  -f values.full-stack-nvidia.example.yaml \
  --set opentelemetry.tracing.enabled= \
  --set openWebUI.existingSecret= \
  --set open-webui.webuiSecret.existingSecretName= \
  >/tmp/rendered-full-stack-nvidia.yaml
```

## Automated Capture

```bash
./hack/capture-local-evidence.sh artifacts/evidence-local-run
```

Review the generated files before publishing or committing them. Do not commit kubeconfigs, secrets, raw private prompts, private customer evidence, model binaries, or screenshots containing credentials.

## Optional Service Checks

With port-forwards running:

```bash
curl -s http://127.0.0.1:11434/api/tags | jq
curl -s http://127.0.0.1:8000/healthz | jq
curl -s http://127.0.0.1:8000/config | jq
curl -s http://127.0.0.1:8888/metrics | head
```

## Evidence Table

| Evidence | Status | Notes |
|---|---|---|
| k3s node Ready | TBD | Capture with `kubectl get nodes -o wide`. |
| NVIDIA RuntimeClass present | TBD | Capture with `kubectl get runtimeclass nvidia`. |
| `nvidia.com/gpu` allocatable | TBD | Capture with node JSONPath command. |
| Helm release deployed | TBD | Capture with `helm list -n llm-observability`. |
| Pods Running | TBD | Capture with `kubectl get pods -n llm-observability -o wide`. |
| CPU profile renders | TBD | Capture with Helm template command. |
| NVIDIA profile renders | TBD | Capture with Helm template command. |
| Service smoke tests pass | TBD | Capture curl outputs or sanitized screenshots. |
