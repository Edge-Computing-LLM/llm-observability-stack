# Sequential GGUF validation on a 1 GiB NVIDIA GPU

This runbook deploys one model at a time on the GeForce 940M reference system.
It uses partial CUDA layer offload and keeps the remaining model layers in
system RAM.

## Why the limit is observational

Kubernetes requests `nvidia.com/gpu: 1`, which assigns the whole GPU. The
GeForce 940M does not support MIG memory partitions, and Ollama does not expose
a per-container hard VRAM quota. Each Modelfile therefore pins `num_gpu`,
`num_ctx`, and `num_batch`, while `gguf-observability` rejects total
`nvidia-smi` usage above 900 MiB.

This is a guardrail, not a hardware-enforced memory sandbox. Display-server and
other host allocations count toward the measured total.

## Profiles

| Model | Values files | Ollama alias |
|---|---|---|
| Qwen 1.8B Chat Q4_K_M | `values.geforce-940m-k3s.yaml` | `qwen-1-8b-chat-q4-k-m-local` |
| Gemma 3 1B IT Q4_K_M | base plus `values.gemma-3-1b-geforce-940m-k3s.yaml` | `gemma3-1b-it-gguf-local` |
| Llama 3.2 1B | base plus `values.llama3.2-1b-geforce-940m-k3s.yaml` | `llama3-2-1b-local` |

Qwen and Gemma use read-only host-mounted GGUF files. Llama derives from the
official Ollama `llama3.2:1b` registry model and inherits its maintained chat
template.

## Safe sequence

1. Deploy and validate `k3s-nvidia-edge`.
2. Install one profile with `helm upgrade --install`.
3. Wait for `deployment/ollama` and confirm exactly one entry in `ollama ps`.
4. Run a fixed public smoke prompt.
5. Capture selected facts with `gguf-observe`; never persist the response.
6. Change to the next overlay. `OLLAMA_MAX_LOADED_MODELS=1` prevents a second
   runner from remaining resident.
7. Confirm total observed VRAM remains at or below 900 MiB after every warmup.

Example observer call:

```bash
GGUF_MODEL=gemma3-1b-it-gguf-local \
GGUF_VRAM_CEILING_MIB=900 \
GGUF_EXPECTED_NUM_GPU=23 \
GGUF_EXPECTED_NUM_CTX=256 \
GGUF_EXPECTED_NUM_BATCH=1 \
../gguf-observability/bin/gguf-observe validate
```

The Ollama volume is retained across Helm upgrades. Do not delete the PVC or
host GGUF directory during a profile change.

## Validated results

On July 19, 2026 all three aliases passed registration, residency, GPU-active,
keep-alive, effective-parameter, workload, and 900 MiB ceiling checks:

| Model | CPU/GPU split | Total observed VRAM | Fixed smoke duration |
|---|---:|---:|---:|
| Gemma 3 1B IT | 62% / 38% | 450 MiB | 1.297 s |
| Qwen 1.8B Chat | 27% / 73% | 824 MiB | 1.932 s |
| Llama 3.2 1B | 52% / 48% | 541 MiB | 1.591 s |

These are local observations, not portable performance claims. Sanitized
evidence is owned by `gguf-observability`.
