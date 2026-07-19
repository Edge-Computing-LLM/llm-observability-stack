# Qwen 1.8B live validation — 2026-07-18

## Scope

This report records the migration of the live single-node Ubuntu 24.04 + k3s
cluster from the local Gemma model to Ollama
`qwen:1.8b-chat-q4_K_M`. The node uses an NVIDIA GeForce 940M with 1 GiB VRAM.

The upstream tag is a Qwen 1.5 1.84B-parameter Q4_K_M model with an approximately
1.2 GB GGUF payload. Its Tongyi Qianwen Research License permits non-commercial
research/evaluation; commercial use requires a separate license.

## Local GGUF

The official Ollama tag was pulled into the retained Ollama PVC. Its model layer
was copied to the host model directory as:

```text
/media/waqasm86/External1/Waqas-Projects/repos-llamatelemetry/llamatelemetry-xubuntu24/models/qwen-1.8b-chat-q4_K_M.gguf
```

Verified payload:

- size: approximately 1.2 GB;
- SHA-256: `ef0125bcc77278420b64229fc29e235435c517f67ab7aa8546a9f5f7be644cef`;
- hash matched Ollama's `sha256-ef0125...` model layer exactly.

The chart mounts the host directory read-only and creates the runtime alias
`qwen-1-8b-chat-q4-k-m-local` from
`Modelfile.qwen-1.8b-chat-q4_K_M`.

## VRAM tuning

The requested 850 GB limit was interpreted as 850 MiB because the physical GPU
has only 1,024 MiB. The upstream model cannot fit fully in VRAM and must use
CPU/GPU split inference.

| Configuration | Result |
|---|---|
| Automatic, batch 512, context 1024 | Failed safely: CUDA graph allocation exceeded 1 GiB |
| 6 GPU layers, batch 512 | Failed safely: graph buffer still too large |
| 13 GPU layers, batch 64 | Failed safely: graph buffer still too large |
| 1 GPU layer, batch 1, context 256 | Passed; 80 MiB VRAM |
| 24 GPU layers, batch 1, context 256 | Passed; 860 MiB VRAM, above ceiling |
| 23 GPU layers, batch 1, context 256 | Passed; 824 MiB VRAM, selected |

Final Ollama evidence:

```text
PROCESSOR: 27%/73% CPU/GPU
CONTEXT:   256
UNTIL:     Forever
GPU VRAM:  824 MiB used, 152 MiB free
layers:    23/25 repeating layers offloaded to CUDA
weights:   734.9 MiB CUDA, 277.5 MiB CPU
KV cache:  2.9 MiB CUDA
graph:     8.0 MiB CUDA
```

`OLLAMA_KEEP_ALIVE=-1` and `OLLAMA_MAX_LOADED_MODELS=1` keep this Qwen runner
resident and prevent a second model from competing for the GPU.

## Deployment results

- Gemma was stopped and removed from the Ollama registry; the original read-only
  host GGUF was not deleted.
- The clean redeployment is Helm revision 1 with chart `0.3.0` and status
  `deployed`.
- Ollama was recreated with zero restarts and automatically created/warmed the
  local Qwen alias.
- Ollama, Open WebUI, Redis, and OpenTelemetry Collector were Ready.
- Open WebUI `/health` returned `{"status":true}`.
- `edge install observability --accelerator auto --yes` completed, including its
  Qwen smoke prompt.
- `ollama ps` reported `Forever` after all tests.

## Inference tests

| Test | Response | Throughput |
|---|---|---:|
| Exact validation | `validation ok` | 15.78 tokens/s |
| Arithmetic (`2 + 2`) | `2 plus 2 equals 4.` | 9.75 tokens/s |
| English-to-Chinese translation | `边缘计算的中文翻译是“边缘计算”。` | 11.47 tokens/s |

Static validation also passed:

- Helm dependency build and lint;
- GeForce profile render with `num_gpu=23`, `num_batch=1`, and `num_ctx=256`;
- native Go Helm smoke tests;
- edge-cli Go tests.

The Qwen model was deliberately left loaded in GPU memory after validation.

Ongoing read-only contract validation and sanitized evidence capture are kept
in the complementary
[`gguf-observability`](https://github.com/Edge-Computing-LLM/gguf-observability)
repository. This stack remains the deployment and model-configuration owner.
