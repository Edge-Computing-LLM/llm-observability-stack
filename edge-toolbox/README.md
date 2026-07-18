# edge-toolbox

`edge-toolbox` is the native Go in-cluster diagnostic image. Its static binary
supports `dns`, `http`, `tcp`, `redis-ping`, `ollama-smoke`, and OpenTelemetry
`seed` commands. The default `serve` command exposes `/healthz` on port 8080.

Examples:

```bash
kubectl exec -n llm-observability deploy/edge-toolbox -- edge-toolbox dns ollama open-webui
kubectl exec -n llm-observability deploy/edge-toolbox -- edge-toolbox http http://ollama:11434/api/tags
kubectl exec -n llm-observability deploy/edge-toolbox -- edge-toolbox redis-ping
kubectl exec -n llm-observability deploy/edge-toolbox -- edge-toolbox ollama-smoke
kubectl exec -n llm-observability deploy/edge-toolbox -- edge-toolbox seed --count 2
```

Build from the repository root with `edge-toolbox/Dockerfile`. The resulting
image contains no shell, package manager, or Python runtime.

```bash
./hack/build-local-image.sh edge-toolbox 0.2.0 . edge-toolbox/Dockerfile
./hack/import-local-image-to-k3s.sh edge-toolbox 0.2.0
```
