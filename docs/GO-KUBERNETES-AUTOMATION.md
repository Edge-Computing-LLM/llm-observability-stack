# Native Go Kubernetes automation

The repository-local `llm-observability` binary replaces the former Python
Kubernetes-client scripts and has no kubeconfig library dependency. It delegates
authentication and API compatibility to the installed `kubectl` binary.

```bash
go build -o bin/llm-observability ./cmd/llm-observability

bin/llm-observability network --namespace llm-observability
bin/llm-observability service-path --namespace llm-observability --service ollama
bin/llm-observability watch-endpoints --namespace llm-observability \
  --service open-webui --timeout 10m
```

All three commands are read-only. `network` prints Pods, Services, Endpoints,
EndpointSlices, and NetworkPolicies. `service-path` resolves a Service selector
to Pods and both endpoint APIs. `watch-endpoints` follows endpoint changes for
the requested bounded timeout.

For deployment and lifecycle operations, prefer the organization-level
`edge-cli`; it enforces Layer 1 before Layer 2 installation and reverse order
for removal.
