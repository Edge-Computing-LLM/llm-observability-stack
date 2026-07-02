# External Feedback Summary

This file tracks sanitized external feedback for `llm-observability-stack`. It is intentionally neutral: it does not claim customer production use, paid deployments, certification, or vendor endorsement.

## Current Status

- Local self-deployed validation: available.
- k3s/NVIDIA GPU deployment evidence: available.
- CPU-only chart rendering path: available.
- Sanitized benchmark artifact: available.
- External production validation: not claimed.

## Evidence Already Available

- `docs/XUBUNTU-K3S-NVIDIA-RUNBOOK.md`
- `docs/LOCAL-K3S-NVIDIA-REPORT-2026-07-02.md`
- `docs/VERIFIED-LOCAL-GPU-RESULTS.md`
- `artifacts/geforce-940m-benchmark.json`
- `hack/validate-local-stack.sh`
- `hack/capture-local-evidence.sh`

## Feedback To Collect

- Whether the deployment flow is clear enough for another local k3s operator.
- Which metrics matter most for local LLM operations: TTFT, throughput, latency, GPU memory, GPU utilization, endpoint health, or Kubernetes readiness.
- Whether CPU-only fallback instructions are sufficient for non-NVIDIA environments.
- What security controls are required before the stack can be used with private workloads.

## Sanitized Feedback Template

```text
Date:
Reviewer category:
Environment:
What was reviewed:
Useful parts:
Confusing parts:
Missing requirements:
Permission to quote or summarize:
```

Store private originals outside the repository. Commit only summaries that are sanitized and explicitly permitted for public use.
