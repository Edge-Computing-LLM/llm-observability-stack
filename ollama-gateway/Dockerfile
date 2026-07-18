FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
COPY cmd/ollama-gateway ./cmd/ollama-gateway
COPY internal/gateway ./internal/gateway
RUN go mod edit -droprequire github.com/Edge-Computing-LLM/k3s-nvidia-edge \
    -dropreplace github.com/Edge-Computing-LLM/k3s-nvidia-edge \
    && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/ollama-gateway ./cmd/ollama-gateway

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/ollama-gateway /usr/local/bin/ollama-gateway
EXPOSE 8000
ENTRYPOINT ["/usr/local/bin/ollama-gateway"]
