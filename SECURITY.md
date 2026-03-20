# Security Policy

## Supported Scope

This repository targets local/self-managed Kubernetes deployments. Security expectations are:

- No credentials in git-tracked files
- Least exposure for local services
- Safe defaults for optional components

## Reporting a Vulnerability

If you find a security issue:

1. Prefer GitHub private vulnerability reporting in this repository.
2. If private reporting is not available, open an issue without exploit details and request a private contact path.

Please include:

- Affected file(s) and configuration
- Reproduction steps
- Impact assessment
- Suggested mitigation (if available)

## Secrets Handling

- Keep secrets only in `values.local-k3s.yaml` or external Kubernetes secrets.
- Never commit real API keys, passwords, tokens, certificates, or private keys.
- Keep `values.local-k3s.example.yaml` sanitized with placeholders.
