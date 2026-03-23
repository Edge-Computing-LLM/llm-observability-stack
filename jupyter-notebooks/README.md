# llm-observability-stack Jupyter Tutorials

This directory contains a full 8-notebook progression for your local k3s-based `llm-observability-stack`, from basic checks to advanced troubleshooting and operations.

## Notebook Sequence

1. **01-environment-smoke-test.ipynb**  
   Validates Python 3.11 kernel, host tooling (`kubectl`, `helm`, `k3s`, `nvidia-smi`), Kubernetes health, memory/GPU visibility, and active values profile.

2. **02-ollama-api-basics.ipynb**  
   Covers direct Ollama API usage (`/api/tags`, `/api/chat`, streaming mode), multi-prompt benchmarking, and latency charting.

3. **03-langchain-proxy-deep-dive.ipynb**  
   Explores `langchain-demo` proxy endpoints (`/healthz`, `/config`, `/invoke`, `/ollama/api/*`) and compares direct vs proxy latency.

4. **04-langsmith-tracing-setup.ipynb**  
   Validates LangSmith credentials, generates traced inference calls, queries runs from your project, and plots observability metrics.

5. **05-open-webui-end-to-end.ipynb**  
   Validates Browser -> Open WebUI -> LangChain proxy -> Ollama -> LangSmith flow with manual browser prompts plus post-run trace analysis.

6. **06-custom-modelfile-workflow.ipynb**  
   Walks through GGUF-backed Modelfile inspection, tuning, optional model creation via Ollama API, and benchmark comparison.

7. **07-python-toolbox-diagnostics.ipynb**  
   Covers toolbox deployment checks, in-cluster diagnostic script execution, pod memory profiling, and GPU correlation during inference.

8. **08-troubleshooting-etcd-simulations.ipynb**  
   Provides advanced health matrix diagnostics, memory pressure analysis, Helm render drills for optional components, and ops runbook patterns.

## Python 3.11 Requirement

All notebooks are authored for **`/usr/local/bin/python3.11`** and use the `python311` kernelspec.

Register kernel if needed:

```bash
/usr/local/bin/python3.11 -m pip install ipykernel
/usr/local/bin/python3.11 -m ipykernel install --user --name python311 --display-name "Python 3.11"
```

Launch Jupyter:

```bash
cd /media/waqasm86/External1/Project-Nvidia-Office/Project-Llamatelemetry/langchain-kubernetes-jupyterlab/llm-observability-stack/jupyter-notebooks
/usr/local/bin/python3.11 -m jupyter lab
```

## Port-Forward Expectations

Several notebooks assume local access to internal `ClusterIP` services. Keep these in separate terminals when needed:

```bash
kubectl port-forward -n llm-observability svc/ollama 11434:11434
kubectl port-forward -n llm-observability svc/langchain-demo 8000:8000
```
