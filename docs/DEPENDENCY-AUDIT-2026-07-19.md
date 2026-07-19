# Dependency audit: 2026-07-19

The local reference repositories and GitHub releases were fetched before this
comparison. Layer 1 is current; several Layer 2 charts have newer upstream
releases.

| Component | Repository version | Current upstream | Decision |
|---|---:|---:|---|
| k3s | v1.36.2+k3s1 | v1.36.2+k3s1 | Current |
| NVIDIA GPU Operator | v26.3.3 | v26.3.3 | Current |
| Ollama Helm chart | 1.50.0 | 1.67.0 | Migration required |
| Ollama image | 0.17.7 | 0.32.1 | Migration required |
| Open WebUI Helm chart | 12.10.0 | 15.2.0 | Migration required |
| Open WebUI image | 0.8.10 | 0.10.2 | Migration required |
| kube-prometheus-stack | 86.2.3 | 87.17.0 | Review CRD/rule changes |
| OpenTelemetry Collector chart | 0.158.1 | 0.165.0 | Migration required |
| OpenTelemetry Operator chart | 0.115.0 | 0.120.0 | Review CRD changes |

Sources: [k3s releases](https://github.com/k3s-io/k3s/releases),
[GPU Operator releases](https://github.com/NVIDIA/gpu-operator/releases),
[Ollama releases](https://github.com/ollama/ollama/releases),
[Ollama chart](https://github.com/otwld/ollama-helm),
[Open WebUI releases](https://github.com/open-webui/open-webui/releases),
[Open WebUI charts](https://github.com/open-webui/helm-charts),
[Prometheus Community charts](https://github.com/prometheus-community/helm-charts),
and [OpenTelemetry charts](https://github.com/open-telemetry/opentelemetry-helm-charts).

## Validation result

The existing pins passed module verification, formatting, unit tests, race
tests, vet, `govulncheck`, builds, Helm lint, and default/CPU/local/GeForce/full
NVIDIA rendering with Helm 4.2.3. GitHub's current `main` workflow is also
green.

The live cluster was not suitable for a safe multi-component upgrade because a
post-reboot stale k3s node address and Flannel interface were breaking pod
networking. Upgrading the charts during that fault would mix configuration and
version variables and would not produce trustworthy evidence. The versions
above are therefore recorded as a dedicated migration, not changed in place.

Layer 2 now depends on the Layer 1 commit that detects stale node networking and
unhealthy GPU Operator pods. Once the host network is repaired, upgrade and
validate the components in this order: observability CRDs/operators,
Prometheus, Ollama, then Open WebUI. Preserve the GeForce 940M resource limits,
model-cleanup rejection, existing PVCs, and rollback values between steps.
