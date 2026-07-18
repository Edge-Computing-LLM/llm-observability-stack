package toolbox

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otlptracegrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func Run(ctx context.Context, args []string, output io.Writer) error {
	if len(args) == 0 {
		return errors.New("command required: serve, dns, http, tcp, redis-ping, ollama-smoke, or seed")
	}
	switch args[0] {
	case "serve":
		return serve(ctx, args[1:])
	case "dns":
		return dns(args[1:], output)
	case "http":
		return httpCheck(ctx, args[1:], output)
	case "tcp":
		return tcpCheck(ctx, args[1:], output)
	case "redis-ping":
		return redisPing(ctx, args[1:], output)
	case "ollama-smoke":
		return ollamaSmoke(ctx, args[1:], output)
	case "seed":
		return seed(ctx, args[1:], output)
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func serve(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("serve", flag.ContinueOnError)
	address := flags.String("address", ":8080", "health server address")
	if err := flags.Parse(args); err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, "{\"status\":\"ok\",\"runtime\":\"go\"}\n")
	})
	server := &http.Server{Addr: *address, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	errCh := make(chan error, 1)
	go func() { errCh <- server.ListenAndServe() }()
	select {
	case <-ctx.Done():
		stop, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(stop)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func dns(args []string, output io.Writer) error {
	if len(args) == 0 {
		return errors.New("dns requires at least one host")
	}
	for _, host := range args {
		addresses, err := net.LookupHost(host)
		if err != nil {
			return fmt.Errorf("resolve %s: %w", host, err)
		}
		fmt.Fprintf(output, "%s\t%s\n", host, strings.Join(addresses, ","))
	}
	return nil
}

func httpCheck(ctx context.Context, args []string, output io.Writer) error {
	flags := flag.NewFlagSet("http", flag.ContinueOnError)
	timeout := flags.Duration("timeout", 10*time.Second, "request timeout")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("http requires one URL")
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, flags.Arg(0), nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: *timeout}
	started := time.Now()
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(response.Body, 500))
	fmt.Fprintf(output, "status=%d duration=%s body=%q\n", response.StatusCode, time.Since(started).Round(time.Millisecond), strings.TrimSpace(string(body)))
	if response.StatusCode >= 400 {
		return fmt.Errorf("HTTP status %d", response.StatusCode)
	}
	return nil
}

func tcpCheck(ctx context.Context, args []string, output io.Writer) error {
	flags := flag.NewFlagSet("tcp", flag.ContinueOnError)
	timeout := flags.Duration("timeout", 5*time.Second, "connect timeout")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("tcp requires host:port")
	}
	dialer := net.Dialer{Timeout: *timeout}
	started := time.Now()
	connection, err := dialer.DialContext(ctx, "tcp", flags.Arg(0))
	if err != nil {
		return err
	}
	_ = connection.Close()
	fmt.Fprintf(output, "connected=%s duration=%s\n", flags.Arg(0), time.Since(started).Round(time.Millisecond))
	return nil
}

func redisPing(ctx context.Context, args []string, output io.Writer) error {
	flags := flag.NewFlagSet("redis-ping", flag.ContinueOnError)
	address := flags.String("address", env("REDIS_HOST", "open-webui-redis")+":"+env("REDIS_PORT", "6379"), "Redis host:port")
	if err := flags.Parse(args); err != nil {
		return err
	}
	dialer := net.Dialer{Timeout: 5 * time.Second}
	connection, err := dialer.DialContext(ctx, "tcp", *address)
	if err != nil {
		return err
	}
	defer connection.Close()
	_ = connection.SetDeadline(time.Now().Add(5 * time.Second))
	if _, err := io.WriteString(connection, "*1\r\n$4\r\nPING\r\n"); err != nil {
		return err
	}
	response, err := bufio.NewReader(connection).ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Fprint(output, response)
	if !strings.HasPrefix(response, "+PONG") {
		return fmt.Errorf("unexpected Redis response %q", response)
	}
	return nil
}

func ollamaSmoke(ctx context.Context, args []string, output io.Writer) error {
	flags := flag.NewFlagSet("ollama-smoke", flag.ContinueOnError)
	base := flags.String("url", env("OLLAMA_BASE_URL", "http://ollama:11434"), "Ollama base URL")
	model := flags.String("model", env("OLLAMA_MODEL", "qwen-1-8b-chat-q4-k-m-local"), "model")
	prompt := flags.String("prompt", "Reply with one short sentence saying the stack is healthy.", "prompt")
	if err := flags.Parse(args); err != nil {
		return err
	}
	payload, _ := json.Marshal(map[string]any{"model": *model, "messages": []map[string]string{{"role": "user", "content": *prompt}}, "stream": false})
	request, _ := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(*base, "/")+"/api/chat", strings.NewReader(string(payload)))
	request.Header.Set("Content-Type", "application/json")
	response, err := (&http.Client{Timeout: 180 * time.Second}).Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(response.Body, 2<<20))
	fmt.Fprintln(output, string(body))
	if response.StatusCode >= 400 {
		return fmt.Errorf("Ollama status %d", response.StatusCode)
	}
	return nil
}

func seed(ctx context.Context, args []string, output io.Writer) error {
	flags := flag.NewFlagSet("seed", flag.ContinueOnError)
	endpoint := flags.String("otlp-endpoint", env("OTEL_EXPORTER_OTLP_ENDPOINT", "opentelemetry-collector:4317"), "OTLP gRPC endpoint")
	apiURL := flags.String("url", env("OBS_INFERENCE_API_URL", "http://ollama:11434/api/chat"), "Ollama chat URL")
	model := flags.String("model", env("OLLAMA_MODEL", "qwen-1-8b-chat-q4-k-m-local"), "model")
	count := flags.Int("count", envInt("OBS_CALL_COUNT", 2), "call count")
	if err := flags.Parse(args); err != nil {
		return err
	}
	shutdown, err := tracing(ctx, *endpoint)
	if err != nil {
		return err
	}
	defer func() { _ = shutdown(context.Background()) }()
	tracer := otel.Tracer("llm-observability-stack.edge-toolbox")
	prompts := []string{"Explain Kubernetes readiness probes briefly.", "Give one use of GPU acceleration for LLM inference.", "Name useful OpenTelemetry GenAI span attributes."}
	client := &http.Client{Timeout: 180 * time.Second}
	failures := 0
	for i := 0; i < *count; i++ {
		prompt := prompts[i%len(prompts)]
		callCtx, span := tracer.Start(ctx, "ollama seed chat")
		span.SetAttributes(attribute.String("gen_ai.system", "ollama"), attribute.String("gen_ai.operation.name", "chat"), attribute.String("gen_ai.request.model", *model), attribute.Int("llm.seed.index", i+1))
		payload, _ := json.Marshal(map[string]any{"model": *model, "stream": false, "messages": []map[string]string{{"role": "user", "content": prompt}}})
		request, _ := http.NewRequestWithContext(callCtx, http.MethodPost, *apiURL, strings.NewReader(string(payload)))
		request.Header.Set("Content-Type", "application/json")
		response, requestErr := client.Do(request)
		if requestErr != nil {
			failures++
			span.RecordError(requestErr)
		} else {
			_, _ = io.Copy(io.Discard, response.Body)
			_ = response.Body.Close()
			if response.StatusCode >= 400 {
				failures++
			}
		}
		span.End()
	}
	fmt.Fprintf(output, "seeded=%d failed=%d\n", *count, failures)
	if failures > 0 {
		return fmt.Errorf("%d seed requests failed", failures)
	}
	return nil
}

func tracing(ctx context.Context, endpoint string) (func(context.Context) error, error) {
	endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	res, err := resource.Merge(resource.Default(), resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(env("OTEL_SERVICE_NAME", "edge-toolbox"))))
	if err != nil {
		return nil, err
	}
	provider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(res))
	otel.SetTracerProvider(provider)
	return provider.Shutdown, nil
}
func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func envInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err == nil {
		return value
	}
	return fallback
}
