# kubectl Command Reference for llm-observability-stack

This reference is organized by operational workflow for the namespace `llm-observability`.

Set a convenience variable:

```bash
NS=llm-observability
```

Profile note:

- The current local example profile keeps `pythonToolbox.enabled=true` and `langsmithDashboardSeeder.enabled=false`.
- See [CONFIG-PROFILES.md](CONFIG-PROFILES.md) before assuming a workload is enabled in a given environment.

## 1. Cluster and Context Safety

```bash
kubectl config current-context
kubectl config get-contexts
kubectl cluster-info
kubectl get nodes -o wide
kubectl version --client
```

## 2. Namespace and Core Inventory

```bash
kubectl get ns
kubectl get all -n $NS
kubectl get pod,svc,pvc,configmap,secret -n $NS
kubectl get events -n $NS --sort-by=.lastTimestamp
```

## 3. Helm-Related Runtime Validation

Detect workload kind first:

```bash
kubectl get deploy,statefulset -n $NS
kubectl get deploy,statefulset -n $NS -o name | grep -E 'ollama|open-webui|langchain-demo|python-toolbox'
```

```bash
kubectl get pods -n $NS -o wide
kubectl get deploy,statefulset -n $NS
kubectl rollout status deploy/ollama -n $NS
kubectl rollout status statefulset/open-webui -n $NS
kubectl rollout status deploy/langchain-demo -n $NS
# optional
kubectl rollout status deploy/python-toolbox -n $NS
```

## 4. Service and Endpoint Checks

```bash
kubectl get svc -n $NS
kubectl get endpoints -n $NS
kubectl get endpointslices -n $NS
kubectl describe svc ollama -n $NS
kubectl describe svc open-webui -n $NS
kubectl describe svc langchain-demo -n $NS
```

## 5. Pod Detail and Scheduling Checks

```bash
kubectl get pods -n $NS -o wide
kubectl describe pod -n $NS <POD_NAME>
kubectl get pods -n $NS -o custom-columns=NAME:.metadata.name,NODE:.spec.nodeName,PHASE:.status.phase,PODIP:.status.podIP
kubectl top pods -n $NS
kubectl top nodes
```

## 6. GPU-Specific Checks

```bash
kubectl get nodes -o json | jq '.items[].status.allocatable'
kubectl get pods -n nvidia-device-plugin
kubectl logs -n nvidia-device-plugin -l app.kubernetes.io/component=device-plugin --tail=200
kubectl describe pod -n $NS deploy/ollama
kubectl exec -it -n $NS deploy/ollama -- env | grep -E 'NVIDIA|CUDA|OLLAMA'
```

## 7. Logs and Events

```bash
kubectl logs -n $NS deploy/ollama --tail=200
kubectl logs -n $NS statefulset/open-webui --tail=200
kubectl logs -n $NS deploy/langchain-demo --tail=200
# optional
kubectl logs -n $NS deploy/python-toolbox --tail=200
kubectl logs -n $NS deploy/redis --tail=200
kubectl get events -n $NS --sort-by=.lastTimestamp | tail -n 50
```

Container-specific logs:

```bash
kubectl logs -n $NS <POD_NAME> -c <CONTAINER_NAME> --tail=200
kubectl logs -n $NS <POD_NAME> -c <CONTAINER_NAME> --previous --tail=200
```

## 8. Exec and In-Pod Diagnostics

```bash
# optional
kubectl exec -it -n $NS deploy/python-toolbox -- bash
kubectl exec -it -n $NS deploy/python-toolbox -- python /workspace/examples/service_dns_check.py
kubectl exec -it -n $NS deploy/python-toolbox -- python /workspace/examples/ollama_smoke.py
kubectl exec -it -n $NS deploy/python-toolbox -- python /workspace/examples/redis_ping.py
kubectl exec -it -n $NS deploy/python-toolbox -- python /workspace/examples/langsmith_healthcheck.py
kubectl exec -it -n $NS deploy/python-toolbox -- python /workspace/examples/langsmith_inference_traces.py
kubectl exec -it -n $NS deploy/python-toolbox -- env OBS_CALL_COUNT_PER_CYCLE=4 OBS_INTERVAL_SECONDS=300 python /workspace/examples/langsmith_dashboard_seed_every_5m.py
```

CronJob checks (only when `langsmithDashboardSeeder.enabled=true`):

```bash
kubectl get cronjob -n $NS langsmith-dashboard-seeder
kubectl get jobs -n $NS --sort-by=.metadata.creationTimestamp | tail -n 10
kubectl get pods -n $NS -l job-name
```

Quick in-pod DNS/TCP checks:

```bash
kubectl exec -it -n $NS deploy/python-toolbox -- nslookup ollama
kubectl exec -it -n $NS deploy/python-toolbox -- nslookup open-webui
kubectl exec -it -n $NS deploy/python-toolbox -- nc -vz ollama 11434
kubectl exec -it -n $NS deploy/python-toolbox -- nc -vz open-webui 8080
kubectl exec -it -n $NS deploy/python-toolbox -- nc -vz langchain-demo 8000
```

## 9. Port-forward and Local API Tests

```bash
kubectl port-forward -n $NS svc/ollama 11434:11434
kubectl port-forward -n $NS svc/open-webui 8080:8080
kubectl port-forward -n $NS svc/langchain-demo 8000:8000
```

Ollama API test:

```bash
curl -s http://localhost:11434/api/tags | jq
curl -s http://localhost:11434/api/chat \
  -H 'Content-Type: application/json' \
  -d '{"model":"gemma3-1b-it-gguf-local","stream":false,"messages":[{"role":"user","content":"health check"}]}' | jq
```

LangChain demo test:

```bash
curl -s http://localhost:8000/healthz | jq
curl -s http://localhost:8000/invoke \
  -H 'Content-Type: application/json' \
  -d '{"prompt":"Say stack is healthy in one line"}' | jq
```

## 10. Rollout Management

```bash
kubectl rollout restart deploy/ollama -n $NS
kubectl rollout restart statefulset/open-webui -n $NS
kubectl rollout restart deploy/langchain-demo -n $NS
# optional
kubectl rollout restart deploy/python-toolbox -n $NS
kubectl rollout status deploy/ollama -n $NS
kubectl rollout status statefulset/open-webui -n $NS
```

Undo rollout (if deployment strategy/history allows):

```bash
kubectl rollout undo deploy/ollama -n $NS
kubectl rollout undo deploy/langchain-demo -n $NS
```

## 11. Config and Secret Inspection

Metadata only:

```bash
kubectl get configmap -n $NS
kubectl get secret -n $NS
kubectl describe configmap ollama-local-modelfile -n $NS
kubectl describe secret langsmith-secrets -n $NS
kubectl describe secret open-webui-secrets -n $NS
```

Decode specific key (use carefully):

```bash
kubectl get secret -n $NS open-webui-secrets -o jsonpath='{.data.WEBUI_SECRET_KEY}' | base64 -d; echo
```

## 12. Storage and PVC Troubleshooting

```bash
kubectl get pvc -n $NS
kubectl describe pvc -n $NS
kubectl get pv
kubectl describe pod -n $NS deploy/ollama | sed -n '/Volumes:/,/QoS Class/p'
kubectl describe pod -n $NS statefulset/open-webui | sed -n '/Volumes:/,/QoS Class/p'
```

## 13. Networking Deep-Dive Commands

DNS and CoreDNS:

```bash
kubectl get pods -n kube-system -l k8s-app=kube-dns
kubectl logs -n kube-system -l k8s-app=kube-dns --tail=200
kubectl exec -it -n $NS deploy/python-toolbox -- cat /etc/resolv.conf
```

Endpoints and EndpointSlices:

```bash
kubectl get endpoints -n $NS -o wide
kubectl get endpointslices -n $NS -o wide
kubectl describe endpoints -n $NS ollama
kubectl describe endpointslices -n $NS
```

Service selectors vs pod labels:

```bash
kubectl get svc -n $NS ollama -o jsonpath='{.spec.selector}' ; echo
kubectl get pods -n $NS --show-labels
kubectl get pods -n $NS -l app.kubernetes.io/name=ollama
```

## 14. Open WebUI and Redis Path Commands

```bash
kubectl get deploy,svc -n $NS | grep -E 'open-webui|redis'
kubectl logs -n $NS statefulset/open-webui --tail=200 | grep -Ei 'redis|websocket|error|warn'
kubectl logs -n $NS deploy/redis --tail=200
kubectl logs -n $NS deploy/open-webui-redis --tail=200
```

Validate only one Redis mode is active for websocket path:

```bash
kubectl get deploy -n $NS | grep -E '^redis|open-webui-redis'
kubectl get svc -n $NS | grep -E '^redis|open-webui-redis'
```

## 15. Resource YAML Capture and Diff

```bash
kubectl get deploy -n $NS ollama -o yaml > /tmp/ollama-deploy.yaml
kubectl get statefulset -n $NS open-webui -o yaml > /tmp/open-webui-sts.yaml
kubectl get svc -n $NS ollama -o yaml > /tmp/ollama-svc.yaml
kubectl get configmap -n $NS ollama-local-modelfile -o yaml > /tmp/modelfile-cm.yaml
```

## 16. Event-Driven Failure Triage Sequence

```bash
kubectl get pods -n $NS
kubectl describe pod -n $NS <FAILING_POD>
kubectl logs -n $NS <FAILING_POD> --all-containers --tail=200
kubectl get events -n $NS --sort-by=.lastTimestamp | tail -n 100
```

## 17. Cleanup and Uninstall

```bash
helm uninstall llm-observability-stack -n $NS
kubectl get all -n $NS
kubectl delete namespace $NS
```

If PVCs remain (depending on reclaim policy):

```bash
kubectl get pvc -n $NS
kubectl delete pvc -n $NS --all
```

## 18. One-liner Status Snapshot

```bash
kubectl get pods,svc,pvc -n $NS -o wide && \
kubectl get endpoints -n $NS && \
kubectl get events -n $NS --sort-by=.lastTimestamp | tail -n 20
```

## 19. Troubleshooting Command Packs

### 19.1 Ollama unavailable

```bash
kubectl get pod -n $NS -l app.kubernetes.io/name=ollama -o wide
kubectl logs -n $NS deploy/ollama --tail=200
kubectl exec -it -n $NS deploy/ollama -- ollama list
kubectl exec -it -n $NS deploy/ollama -- ls -lah /models/gguf
kubectl get endpoints -n $NS ollama -o yaml
```

### 19.2 Open WebUI cannot chat

```bash
kubectl logs -n $NS statefulset/open-webui --tail=200
kubectl get svc,endpoints -n $NS open-webui ollama
kubectl exec -it -n $NS deploy/python-toolbox -- nc -vz ollama 11434
kubectl exec -it -n $NS deploy/python-toolbox -- python /workspace/examples/ollama_smoke.py
```

### 19.3 LangChain demo failures

```bash
kubectl logs -n $NS deploy/langchain-demo --tail=200
kubectl describe deploy -n $NS langchain-demo
kubectl port-forward -n $NS svc/langchain-demo 8000:8000
curl -s http://localhost:8000/healthz | jq
curl -s http://localhost:8000/config | jq
```
