# hack Scripts

This directory contains the local image workflow helpers for `llm-observability-stack`.

## Scripts

- `build-local-image.sh`
  - builds a local image for a component directory
  - prefers `nerdctl --namespace k8s.io`
  - falls back to Docker when `nerdctl` is unavailable

- `import-local-image-to-k3s.sh`
  - saves a local image to a temporary tarball
  - imports that image into k3s containerd
  - skips import when the image already exists in containerd

- `bootstrap-enterprise-pilot-k3s.sh`
  - installs Prometheus Operator CRDs from the vendored kube-prometheus-stack chart
  - installs the local full-stack k3s profile with the vendored OpenTelemetry Collector subchart
  - calls `detect-runtime-profile.sh` by default so NVIDIA GPU nodes are used when available and CPU mode is used otherwise
  - passes any extra CLI flags through to `helm upgrade --install`

- `detect-runtime-profile.sh`
  - inspects Kubernetes node allocatable resources
  - writes `.generated/values.runtime-detected.yaml`
  - enables NVIDIA scheduling when `nvidia.com/gpu` is advertised
  - disables NVIDIA runtime and GPU requests when only CPU capacity is available
  - remains a thin repo-local Helm overlay helper; durable accelerator detection
    and ordered installation are implemented in Go by `edge-cli`

## Typical Usage

Build/import `ollama-gateway`:

```bash
./hack/build-local-image.sh ollama-gateway 0.2.0 . ollama-gateway/Dockerfile
./hack/import-local-image-to-k3s.sh ollama-gateway 0.2.0
```

Build/import `edge-toolbox`:

```bash
./hack/build-local-image.sh edge-toolbox 0.2.0 . edge-toolbox/Dockerfile
./hack/import-local-image-to-k3s.sh edge-toolbox 0.2.0
```

Bootstrap the local full-stack profile:

```bash
./hack/bootstrap-enterprise-pilot-k3s.sh \
  --set ollamaGateway.enabled=false \
  --set edgeToolbox.enabled=false
```

Force a CPU-only render or install path for validation:

```bash
./hack/detect-runtime-profile.sh --mode cpu
helm template llm-observability-stack . \
  -f values.enterprise-pilot-k3s.yaml \
  -f .generated/values.runtime-detected.yaml
```

Enable `ollamaGateway` and `edgeToolbox` only after their local images have been imported into k3s containerd.

## When To Use These Scripts

- after changing `cmd/ollama-gateway` or `internal/gateway`
- after changing scripts or dependencies in `edge-toolbox/`
- when refreshing local image tags used by the chart

## After Import

Restart the matching Kubernetes workload:

```bash
kubectl rollout restart deploy/ollama-gateway -n llm-observability
kubectl rollout restart deploy/edge-toolbox -n llm-observability
```

Then verify:

```bash
kubectl rollout status deploy/ollama-gateway -n llm-observability
kubectl rollout status deploy/edge-toolbox -n llm-observability
sudo k3s ctr images ls | grep -E 'ollama-gateway|edge-toolbox'
```
