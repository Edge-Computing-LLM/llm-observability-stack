# ollama-gateway

`ollama-gateway` is a native Go HTTP gateway for the local Ollama service. It
provides `/healthz`, `/config`, `/invoke`, `/ollama/*`, and `/metrics`, preserves
streaming responses, and exports OpenTelemetry traces when an OTLP endpoint is
configured.

Build the local image from the repository root:

```bash
./hack/build-local-image.sh ollama-gateway 0.2.0 . ollama-gateway/Dockerfile
./hack/import-local-image-to-k3s.sh ollama-gateway 0.2.0
```

The image contains one non-root static Go binary and no Python runtime.
