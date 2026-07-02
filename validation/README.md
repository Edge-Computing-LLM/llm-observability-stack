# Validation Artifacts

This folder contains sanitized local validation material for `llm-observability-stack`.

The current evidence is based on a self-deployed Xubuntu 24 + k3s + NVIDIA GPU environment and chart-level CPU/GPU rendering tests. It does not claim external customer production use, third-party certification, or hardware-vendor endorsement.

## Contents

| File or folder | Purpose |
|---|---|
| `deployment-proof.md` | Commands for collecting local deployment evidence. |
| `local-run-report-2026-06-22.md` | Earlier local validation run report. |
| `pilot-report.md` | Historical local validation report, retained as a technical report. |
| `stack-one-page-summary.md` | One-page technical summary of the stack. |
| `customer-validation-summary.md` | Neutral feedback and external-review tracker. |
| `feedback-log-template.md` | Template for tracking community or operator feedback. |
| `redacted-feedback-template.md` | Template for sanitized feedback summaries. |
| `screenshots/` | Place timestamped UI and terminal screenshots here. |
| `benchmark-results/` | Place fresh benchmark JSON, summaries, and sanitized benchmark exports here. |

## Current Evidence Boundary

- The project has verified local k3s deployment evidence.
- The verified local profile has run on Xubuntu 24 with NVIDIA GPU scheduling, Ollama, Open WebUI, Prometheus, Grafana, and a Gemma 3 1B IT Q4_K_M GGUF model.
- CPU-only rendering and deployment paths are included for environments without NVIDIA GPUs.
- Private prompts, raw customer data, kubeconfigs, credentials, and model binaries must stay outside this repository.

Do not add private names, secrets, kubeconfigs, proprietary prompts, or unapproved third-party claims to this public folder.
