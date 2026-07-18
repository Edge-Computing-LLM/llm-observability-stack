# Grafana dashboards

The JSON files in this directory are the source of truth for the Grafana
dashboards installed by this Helm chart. `templates/grafana-dashboards.yaml`
packages every `dashboards/*.json` file into a labelled ConfigMap, and the
Grafana sidecar from `kube-prometheus-stack` provisions them automatically.

## Edge LLM dashboard migration

`edge-llm-observability.json` recreates the information formerly rendered by
the standalone `Frontend-Edge-LLM-Observability` Vue application:

- live NVIDIA DCGM utilization, memory, temperature, clock, PCIe, model,
  driver, host, PCI bus, and UUID information;
- workload readiness for the `gpu-operator` and `llm-observability`
  namespaces, backed by kube-state-metrics;
- the cluster-internal service inventory and ports;
- the validated Qwen 1.8B Chat Q4_K_M runtime profile;
- DCGM, Ollama/Qwen, Open WebUI, and OpenTelemetry readiness views;
- clear provenance notes that distinguish live Prometheus data from the
  chart-owned reference profile.

The migration was verified against the final frontend source commit
`13a4fc8e18d91d7602ed04aee883e1b4b8515a20` before the standalone repository
was retired.

The GeForce 940M profile enables a compact Prometheus, Grafana,
kube-state-metrics, node-exporter, and Prometheus Operator deployment. It
disables Alertmanager, default dashboards, default rules, and k3s-incompatible
control-plane scrapes to fit the single-node development machine.

## Reproduce locally

```bash
helm dependency build .
helm upgrade --install llm-observability-stack . \
  --namespace llm-observability \
  --create-namespace \
  -f values.geforce-940m-k3s.yaml \
  --wait --timeout 15m

kubectl -n llm-observability port-forward \
  svc/llm-observability-stack-grafana 3000:80
```

Open <http://127.0.0.1:3000>, sign in with the local profile credentials, and
select **Edge LLM Observability - Ubuntu + k3s + NVIDIA GPU**. The bundled
credentials are only appropriate while Grafana remains a local ClusterIP
accessed through `kubectl port-forward`; replace them before any ingress or
shared access is enabled.

Do not edit a provisioned dashboard only in the Grafana UI. Export reviewed
changes as JSON into this directory so the Helm installation remains
reproducible.
