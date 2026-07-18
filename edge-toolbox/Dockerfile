FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
COPY cmd/edge-toolbox ./cmd/edge-toolbox
COPY internal/toolbox ./internal/toolbox
RUN go mod edit -droprequire github.com/Edge-Computing-LLM/k3s-nvidia-edge \
    -dropreplace github.com/Edge-Computing-LLM/k3s-nvidia-edge \
    && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/edge-toolbox ./cmd/edge-toolbox

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/edge-toolbox /usr/local/bin/edge-toolbox
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/edge-toolbox"]
CMD ["serve"]
