# Kubernetes in Action coverage mapping

This toolkit maps the source domains from a local `kubernetes-in-action-2nd-edition` checkout into operational scripts for your `llm-observability` k3s cluster.

## Source reviewed

- Local source checkout: set `KIA_SOURCE_ROOT` in your shell or notes if you want to record the exact path on your machine.
- Book PDF: `Kubernetes in Action_Kevin Conner.pdf`
- Chapter source folders reviewed: `Chapter02` through `Chapter18`
- High-frequency resource kinds in examples: `Pod`, `ConfigMap`, `Secret`, `PersistentVolume`, `PersistentVolumeClaim`, `Service`, `Ingress`, `Gateway`, `HTTPRoute`, `Deployment`, `StatefulSet`, `DaemonSet`, `Job`, `CronJob`

## Domain mapping

- Pod lifecycle and multi-container patterns (Ch. 2, 5, 6):
  - `kubectl/02_namespaces_and_workloads.sh`
  - `kubectl/03_pods_lifecycle_debug.sh`
- Metadata/labels/scheduling and config injection (Ch. 7, 8):
  - `kubectl/02_namespaces_and_workloads.sh`
  - `kubectl/04_configmaps_and_secrets.sh`
- Volumes and persistent storage (Ch. 9, 10):
  - `kubectl/05_storage_and_volumes.sh`
- Services, DNS, endpoints, ingress, gateway API, traffic rules (Ch. 11, 12, 13):
  - `kubectl/06_networking_core.sh`
  - `kubectl/07_networking_advanced.sh`
  - `runbooks/run_go_suite.sh` (`network` and `service-path` commands)
- Controllers and rollout management (Ch. 14, 15, 16, 17):
  - `kubectl/02_namespaces_and_workloads.sh`
  - `runbooks/run_go_suite.sh` (`status` and `validate` commands)
- Batch workloads (Ch. 18):
  - `kubectl/09_jobs_and_batch.sh`
- LLM stack-specific health integration:
  - `kubectl/10_llm_observability_stack_checks.sh`

## Notes

- This suite is optimized for your local k3s chart profile (`llm-observability` namespace).
- Scripts are read-only by default and require explicit mutation opt-in via `APPLY_CHANGES=1`.
