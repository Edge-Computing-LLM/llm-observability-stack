# Kubernetes Networking for llm-observability-stack

This document explains networking behavior for this project and provides practical troubleshooting workflows.

Namespace assumed in examples:

```bash
NS=llm-observability
```

Profile note:

- The current local example profile keeps `edgeToolbox.enabled=true`.
- Verify the active profile in [CONFIG-PROFILES.md](CONFIG-PROFILES.md) and with `helm get values ... -a` before assuming toolbox availability.

## 1. Networking Topology

Core in-cluster services:

- `ollama` on port `11434`
- `open-webui` on port `8080`
- `ollama-gateway` on port `8000`
- `open-webui-redis` (default mode) or `redis` (custom mode) on `6379`

### 1.1 Traffic flows

1. Browser -> `open-webui` service
2. Open WebUI app container -> `ollama-gateway` proxy (`OLLAMA_BASE_URLS=http://ollama-gateway:8000/ollama`)
3. Ollama gateway -> `ollama` service (`OLLAMA_BASE_URL`)
4. Open WebUI websocket manager -> Redis URL (`REDIS_URL` / `WEBSOCKET_REDIS_URL`)

## 2. Service Discovery and DNS

Pods resolve services by:

- Short name: `ollama`
- Namespace name: `ollama.llm-observability`
- Full FQDN: `ollama.llm-observability.svc.cluster.local`

Check DNS from toolbox pod:

```bash
kubectl exec -it -n $NS deploy/edge-toolbox -- nslookup ollama
kubectl exec -it -n $NS deploy/edge-toolbox -- nslookup open-webui
kubectl exec -it -n $NS deploy/edge-toolbox -- nslookup ollama-gateway
```

Inspect resolver config in pod:

```bash
kubectl exec -it -n $NS deploy/edge-toolbox -- cat /etc/resolv.conf
```

Inspect CoreDNS health:

```bash
kubectl get pods -n kube-system -l k8s-app=kube-dns
kubectl logs -n kube-system -l k8s-app=kube-dns --tail=200
```

## 3. Services, Selectors, Endpoints, and EndpointSlices

A Service forwards traffic only to Pods matching its selector.

### 3.1 Verify selector correctness

```bash
kubectl get svc -n $NS ollama -o yaml
kubectl get pods -n $NS --show-labels | grep ollama
kubectl get pods -n $NS -l app.kubernetes.io/name=ollama
```

### 3.2 Verify endpoints

```bash
kubectl get endpoints -n $NS ollama -o yaml
kubectl get endpointslices -n $NS -l kubernetes.io/service-name=ollama -o yaml
```

Empty endpoints means one of:

- no matching Pods
- matching Pods not Ready
- selector mismatch

## 4. ClusterIP, LoadBalancer, and Port-Forward Behavior

This project defaults to `ClusterIP` and can be locally exposed as `LoadBalancer` in local values overrides.

### 4.1 Check current service types

```bash
kubectl get svc -n $NS
```

### 4.2 Local access with port-forward

```bash
kubectl port-forward -n $NS svc/open-webui 8080:8080
kubectl port-forward -n $NS svc/ollama 11434:11434
kubectl port-forward -n $NS svc/ollama-gateway 8000:8000
```

### 4.3 LoadBalancer on local k3s

When using `LoadBalancer` in local clusters, behavior depends on installed LB implementation (for example, ServiceLB/klipper-lb in k3s).

Check assigned external IP/hostname:

```bash
kubectl get svc -n $NS -o wide
kubectl describe svc -n $NS open-webui
```

## 5. Open WebUI Websocket Networking

Open WebUI websocket support depends on Redis path consistency.

Two possible modes:

1. Subchart-managed Redis: `open-webui-redis`
2. Root custom Redis: `redis`

The stack should use one coherent websocket URL strategy at a time.

### 5.1 Verify active Redis mode

```bash
kubectl get deploy,svc -n $NS | grep -E 'open-webui-redis|redis'
kubectl logs -n $NS statefulset/open-webui --tail=200 | grep -Ei 'redis|websocket|error|warn'
```

### 5.2 Check Redis reachability from toolbox

```bash
kubectl exec -it -n $NS deploy/edge-toolbox -- nc -vz redis 6379
kubectl exec -it -n $NS deploy/edge-toolbox -- nc -vz open-webui-redis 6379
```

## 6. Pod-to-Service and Pod-to-Pod Debugging

### 6.1 Toolbox DNS + TCP + HTTP checks

```bash
kubectl exec -it -n $NS deploy/edge-toolbox -- edge-toolbox dns ollama ollama-gateway open-webui open-webui-redis
kubectl exec -it -n $NS deploy/edge-toolbox -- edge-toolbox ollama-smoke
```

### 6.2 Direct health probes

```bash
kubectl exec -it -n $NS deploy/edge-toolbox -- curl -sS http://ollama:11434/api/tags
kubectl exec -it -n $NS deploy/edge-toolbox -- curl -sS http://ollama-gateway:8000/healthz
kubectl exec -it -n $NS deploy/edge-toolbox -- curl -sSI http://open-webui:8080/
```

## 7. Endpoint Troubleshooting Playbook

Run this sequence when a service is failing:

1. Check service object

```bash
kubectl get svc -n $NS <SERVICE>
kubectl describe svc -n $NS <SERVICE>
```

2. Check endpoints

```bash
kubectl get endpoints -n $NS <SERVICE> -o yaml
kubectl get endpointslices -n $NS -l kubernetes.io/service-name=<SERVICE>
```

3. Check backing pods and readiness

```bash
kubectl get pods -n $NS -l <LABEL_SELECTOR> -o wide
kubectl describe pod -n $NS <POD>
kubectl logs -n $NS <POD> --all-containers --tail=200
```

4. Check in-cluster connectivity from toolbox

```bash
kubectl exec -it -n $NS deploy/edge-toolbox -- nslookup <SERVICE>
kubectl exec -it -n $NS deploy/edge-toolbox -- nc -vz <SERVICE> <PORT>
```

5. Check namespace consistency

```bash
kubectl get all -n $NS
kubectl get all -A | grep -E 'ollama|open-webui|ollama-gateway|edge-toolbox|redis'
```

## 8. Network Policies

Current stack does not require NetworkPolicy by default, so traffic is typically unrestricted inside namespace unless cluster policies enforce restrictions.

Check for policies:

```bash
kubectl get networkpolicy -n $NS
kubectl describe networkpolicy -n $NS
```

If policies are present and traffic fails, inspect ingress/egress rules against:

- source pod labels
- destination pod labels
- destination ports
- namespace selectors

## 9. Ingress and Gateway API Notes

Open WebUI subchart supports:

- Ingress (`ingress.enabled`)
- Route (`route.enabled`)
- Managed certificate (`managedCertificate.enabled`)

Verify objects:

```bash
kubectl get ingress -n $NS
kubectl get httproute,tcproute,grpcroute -n $NS
kubectl describe ingress -n $NS
```

## 10. Networking Observability Commands

### 10.1 Live watch

```bash
kubectl get pods,svc,endpoints -n $NS -w
```

### 10.2 Event watch for networking errors

```bash
kubectl get events -n $NS --sort-by=.lastTimestamp | grep -Ei 'dns|network|connection|endpoint|probe|timeout'
```

### 10.3 Endpoint change watch

```bash
kubectl get endpoints -n $NS -w
kubectl get endpointslices -n $NS -w
```

## 11. Local vs In-Cluster Connectivity

- **In-cluster path** uses Service DNS and cluster networking.
- **Local path** for browsers/tools uses LoadBalancer, NodePort, or port-forward.

If local access fails but in-cluster works, focus on:

- service type/external IP
- host firewall
- local LB implementation
- port-forward session health

## 12. Known Networking Pitfalls for This Stack

1. Service selector mismatch causes empty endpoints.
2. Namespace mismatch between components causes DNS/service failures.
3. Redis mode mismatch (subchart Redis vs custom Redis) causes websocket instability.
4. Pod readiness failing removes endpoint addresses, appearing as intermittent service failure.
5. CoreDNS issues can mimic app/network failures.

## 13. Networking Validation Checklist

Run this after each deploy/upgrade:

```bash
kubectl get svc,endpoints -n $NS
kubectl get pods -n $NS -o wide
kubectl exec -it -n $NS deploy/edge-toolbox -- edge-toolbox dns ollama ollama-gateway open-webui open-webui-redis
kubectl exec -it -n $NS deploy/edge-toolbox -- edge-toolbox ollama-smoke
kubectl port-forward -n $NS svc/ollama-gateway 8000:8000
curl -s http://localhost:8000/healthz | jq
```

## 14. Related Docs

- `docs/KUBECTL-COMMAND-REFERENCE.md`
- `docs/GO-KUBERNETES-AUTOMATION.md`
- `docs/PROJECT-DOCUMENTATION.md`
