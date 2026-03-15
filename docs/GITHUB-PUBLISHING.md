# GitHub Publishing Guide

Repository target:

- https://github.com/waqasm86/llm-observability-stack

## Local publish workflow

Run from `llm-observability-stack/`:

```bash
git init
git branch -M main
git remote add origin https://github.com/waqasm86/llm-observability-stack.git
git add .
git commit -m "Initial commit: llm-observability-stack"
git push -u origin main
```

If `origin` already exists:

```bash
git remote set-url origin https://github.com/waqasm86/llm-observability-stack.git
```

## Size guard (<100 MB total)

Check overall folder size:

```bash
du -sh .
```

Check largest files:

```bash
find . -type f -printf '%s %p\n' | sort -nr | head -30
```

Check staged payload size before commit:

```bash
git ls-files -z | xargs -0 du -ch | tail -1
```

## Secret guard

Do not commit:

- `values.local-k3s.yaml`
- `.webui_secret_key`
- `rendered.yaml`
- `.env*`

These are already excluded in `.gitignore`.
