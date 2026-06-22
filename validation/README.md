# NVIDIA Inception 2026 Validation Package

This folder contains customer validation and deployment evidence material for the `llm-observability-stack` NVIDIA Inception 2026 application.

The current evidence package is based on a self-deployed technical pilot and repository-verified proof of deployment. It does not claim a paying customer, formal enterprise pilot, approved partner validation, or NVIDIA/Lenovo certification. The materials are written to support the NVIDIA Inception customer validation requirement through the strongest currently available categories:

- Pilot report or project summary.
- Proof of deployment.
- Technical validation evidence from the repository.
- Templates for future redacted customer or partner feedback.

## Contents

| File or folder | Purpose |
|---|---|
| `pilot-report.md` | Polished pilot report for Edge LLM inference and observability on Kubernetes with NVIDIA GPU. |
| `deployment-proof.md` | Practical evidence capture guide with Helm, kubectl, GPU, Prometheus, Grafana, inference, logs, and screenshot commands. |
| `customer-validation-summary.md` | NVIDIA-facing summary of which validation requirement is satisfied today, what is missing, and next external validation actions. |
| `nvidia-inception-one-page-summary.md` | One-page competition summary for application, partner, or accelerator review. |
| `partner-outreach-log-template.md` | Tracking table for customer, partner, accelerator, cloud provider, and community outreach. |
| `redacted-feedback-template.md` | Template for safely capturing redacted email, LinkedIn, or reviewer feedback. |
| `screenshots/` | Place timestamped UI and terminal screenshots here before submission. |
| `benchmark-results/` | Place fresh benchmark JSON, summaries, and sanitized benchmark exports here before submission. |

## Evidence Boundary

Repository evidence currently supports the following claims:

- The project is a pilot-ready, self-deployed EdgeLLM observability platform.
- The stack deploys through Helm on k3s/Kubernetes.
- The verified local profile has run on a Lenovo ThinkPad T450s with NVIDIA GeForce 940M, NVIDIA device plugin, Ollama, and a Gemma 3 1B IT Q4_K_M GGUF model.
- The repository contains Prometheus metrics, Grafana dashboards, Prometheus alert rules, OpenTelemetry Collector configuration, NVIDIA/DCGM monitoring hooks, benchmark tooling, and sanitized benchmark evidence.
- NVIDIA NIM is represented as a planned/optional integration path through `/v1/metrics` ServiceMonitor configuration, not as a completed deployed NIM proof.

Do not add private customer names, raw secrets, kubeconfigs, proprietary prompts, or unapproved partner claims to this public folder.
