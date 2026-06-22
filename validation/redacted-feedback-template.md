# Redacted Feedback Template

Use this template for customer, partner, accelerator, cloud provider, or community reviewer feedback. Store private originals outside the public repository. Only commit sanitized, permissioned summaries.

## Feedback Record

| Field | Value |
|---|---|
| Date received | `YYYY-MM-DD` |
| Source type | customer / partner / accelerator / cloud provider / community reviewer |
| Organization | `Redacted` or approved organization name |
| Contact role/title | `Redacted` or approved title |
| Contact name | `Redacted` unless explicitly approved |
| Channel | Email / LinkedIn / Slack / call / form / other |
| Permission to use in NVIDIA application? | Yes / No / Pending |
| Permission scope | Anonymous / role only / organization name / direct quote / reference contact |
| Related use case | Private LLM / edge AI / GPU observability / Kubernetes deployment / other |

## Redacted Quote or Summary

> Paste approved redacted quote here.

If direct quote permission is not available, use an anonymized summary:

> A reviewer from `<category>` confirmed that `<problem>` is relevant because `<reason>`, and indicated that `<evidence or feature>` would be useful for evaluating an Edge LLM pilot.

## Value Confirmed

Check all that apply:

- [ ] Problem is relevant.
- [ ] Would consider a technical pilot.
- [ ] Wants benchmark evidence.
- [ ] Wants Grafana/Prometheus dashboards.
- [ ] Wants GPU utilization and memory evidence.
- [ ] Wants NVIDIA NIM or GPU Operator validation.
- [ ] Wants Kubernetes/OpenShift packaging review.
- [ ] Wants security/data-flow review.
- [ ] Will introduce another reviewer or partner.

## Reviewer Notes

```text
Paste sanitized notes here. Remove names, emails, tokens, hostnames, customer names, and confidential details.
```

## Follow-Up Actions

- [ ] Send pilot report.
- [ ] Send one-page summary.
- [ ] Send screenshots.
- [ ] Schedule demo.
- [ ] Ask for permissioned quote.
- [ ] Ask for pilot criteria.
- [ ] Ask for hardware/cloud validation support.

## Evidence Attachment Index

| Artifact | Location | Public or private? | Notes |
|---|---|---|---|
| Original email/message | Private path or inbox reference | Private | Do not commit. |
| Redacted screenshot/PDF | `validation/screenshots/<file>` | Public if sanitized | Optional. |
| Approved quote | This file or private evidence folder | Depends on permission | Confirm scope. |
