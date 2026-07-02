# Redacted Feedback Template

Use this template for community, operator, lab, or internal reviewer feedback. Store private originals outside the public repository. Only commit sanitized summaries that the reviewer approved for public use.

## Feedback Metadata

| Field | Value |
|---|---|
| Date received | YYYY-MM-DD |
| Source type | community reviewer / local operator / lab user / internal reviewer |
| Public attribution allowed? | yes / no / anonymized only |
| Workload reviewed | README / runbook / Helm install / dashboard / benchmark |
| Environment | OS, Kubernetes distribution, CPU/GPU, model type |

## Approved Summary

> A reviewer from `<category>` confirmed that `<problem>` is relevant because `<reason>`, and indicated that `<evidence or feature>` would be useful for evaluating local LLM operations.

## Evidence Type

- [ ] Reviewed README or documentation.
- [ ] Reviewed demo or screenshots.
- [ ] Ran Helm template or install.
- [ ] Ran CPU-only deployment.
- [ ] Ran NVIDIA GPU deployment.
- [ ] Suggested documentation changes.
- [ ] Suggested security or operations changes.

## Constraints

- No private names unless explicitly approved.
- No raw emails or private messages.
- No proprietary prompts or data.
- No credentials, kubeconfigs, internal URLs, screenshots with secrets, or private cluster details.
