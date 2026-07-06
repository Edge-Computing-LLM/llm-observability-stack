# python-toolbox

`python-toolbox` is a minimal in-cluster helper image for diagnostics, connectivity checks, and notebook-backed Kubernetes troubleshooting.

## Purpose

This component exists to give the project a safe in-cluster execution point for:

- DNS resolution checks
- internal Service connectivity checks
- Ollama smoke tests
- Redis checks
- OpenTelemetry API health and optional trace seeding

It is intentionally lighter than a full notebook image and is meant for shell/script execution inside the cluster.

## Local Profile Status

- `pythonToolbox.enabled: true`
- continuous OpenTelemetry seeder jobs remain disabled by default

## Included Example Scripts

- `service_dns_check.py`
- `ollama_smoke.py`
- `redis_ping.py`
- `otel_genai_inference_traces.py`
- `otel_genai_trace_seed_every_5m.py`

Scripts are copied into:

- `/workspace/examples`

## Build and Refresh

Build:

```bash
../hack/build-local-image.sh python-toolbox 0.2.0 ./python-toolbox
```

Import:

```bash
../hack/import-local-image-to-k3s.sh python-toolbox 0.2.0
```

Restart:

```bash
kubectl rollout restart deploy/python-toolbox -n llm-observability
kubectl rollout status deploy/python-toolbox -n llm-observability
```

## Daily Usage

Quick shell:

```bash
kubectl exec -it -n llm-observability deploy/python-toolbox -- bash
```

Run individual scripts:

```bash
kubectl exec -it -n llm-observability deploy/python-toolbox -- \
  python /workspace/examples/service_dns_check.py

kubectl exec -it -n llm-observability deploy/python-toolbox -- \
  python /workspace/examples/ollama_smoke.py
```

Seed OpenTelemetry chart data:

```bash
kubectl exec -it -n llm-observability deploy/python-toolbox -- \
  python /workspace/examples/otel_genai_inference_traces.py
```

Run the continuous 5-minute seeder only when you explicitly want chart activity:

```bash
kubectl exec -it -n llm-observability deploy/python-toolbox -- \
  env OBS_CALL_COUNT_PER_CYCLE=4 OBS_INTERVAL_SECONDS=300 \
  python /workspace/examples/otel_genai_trace_seed_every_5m.py
```

## Operational Notes

- This image is intentionally small and does not include JupyterLab.
- The notebooks use this pod as an in-cluster probe target, especially for DNS and service reachability drills.
- If a notebook says a script is missing, rebuild and re-import the image so the running deployment matches the local source tree.
