# GitHub Publishing Guide

Repository target:

- https://github.com/waqasm86/llm-observability-stack

## Daily publish workflow

Run from `llm-observability-stack/`:

```bash
git checkout main
git pull --rebase origin main
```

Review changes and ensure no secrets are staged:

```bash
git status --short
git diff -- values.yaml values.local-k3s.example.yaml README.md .gitignore
```

Validate chart rendering:

```bash
helm lint .
helm template llm-observability-stack . >/tmp/rendered-default.yaml
helm template llm-observability-stack . -f values.local-k3s.example.yaml >/tmp/rendered-local.yaml
```

Commit and push:

```bash
git add README.md .gitignore values.yaml values.local-k3s.example.yaml templates/ \
  .github/ docs/ CONTRIBUTING.md SECURITY.md SUPPORT.md CODE_OF_CONDUCT.md .helmignore
git commit -m "docs: refresh GitHub guidance and harden local values workflow"
git push origin main
```

## Remote setup (first time only)

```bash
git remote add origin https://github.com/waqasm86/llm-observability-stack.git
# or, if origin already exists
git remote set-url origin https://github.com/waqasm86/llm-observability-stack.git
```

## Secret guard

Do not commit:

- `values.local-k3s.yaml`
- `.webui_secret_key`
- `.env*`
- `*.pem`, `*.key`, `*.crt`
- rendered local artifacts and private debug dumps

These are blocked by `.gitignore` and `.helmignore`.

## Size guard

Check repository size before pushing large changes:

```bash
du -sh .
find . -type f -printf '%s %p\n' | sort -nr | head -30
```
