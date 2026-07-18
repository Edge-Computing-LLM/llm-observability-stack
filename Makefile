GO ?= go

.PHONY: build test vet helm-check check

build:
	@mkdir -p bin
	$(GO) build -o bin/llm-observability ./cmd/llm-observability
	$(GO) build -o bin/ollama-gateway ./cmd/ollama-gateway
	$(GO) build -o bin/edge-toolbox ./cmd/edge-toolbox

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

helm-check:
	helm lint .
	helm template llm-observability-stack . -f values.cpu-k3s.yaml >/dev/null
	helm template llm-observability-stack . -f values.geforce-940m-k3s.yaml >/dev/null

check: test vet helm-check build
