package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otlptracegrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const defaultModel = "qwen-1-8b-chat-q4-k-m-local"

type Config struct {
	Address      string
	UpstreamURL  string
	Model        string
	ServiceName  string
	OTLPEndpoint string
	TraceEnabled bool
	Timeout      time.Duration
}

func ConfigFromEnv() Config {
	port := env("PORT", "8000")
	timeout, err := time.ParseDuration(env("OLLAMA_PROXY_TIMEOUT", "180s"))
	if err != nil {
		timeout = 180 * time.Second
	}
	return Config{Address: ":" + port, UpstreamURL: env("OLLAMA_UPSTREAM_BASE_URL", env("OLLAMA_BASE_URL", "http://ollama:11434")), Model: env("OLLAMA_MODEL", defaultModel), ServiceName: env("OTEL_SERVICE_NAME", "ollama-gateway"), OTLPEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), TraceEnabled: strings.EqualFold(env("OTEL_TRACES_ENABLED", "true"), "true"), Timeout: timeout}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

type Metrics struct {
	Requests *prometheus.CounterVec
	Duration *prometheus.HistogramVec
	Active   *prometheus.GaugeVec
	TTFT     *prometheus.HistogramVec
}

func newMetrics(registry *prometheus.Registry) *Metrics {
	m := &Metrics{
		Requests: prometheus.NewCounterVec(prometheus.CounterOpts{Name: "llm_observability_http_requests_total", Help: "Gateway HTTP requests."}, []string{"method", "route", "status"}),
		Duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "llm_observability_llm_request_duration_seconds", Help: "LLM request duration.", Buckets: prometheus.DefBuckets}, []string{"model", "route"}),
		Active:   prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "llm_observability_llm_active_requests", Help: "Active LLM requests."}, []string{"model", "route"}),
		TTFT:     prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: "llm_observability_time_to_first_token_seconds", Help: "Time to first upstream byte.", Buckets: prometheus.DefBuckets}, []string{"model", "route"}),
	}
	registry.MustRegister(m.Requests, m.Duration, m.Active, m.TTFT)
	return m
}

type Server struct {
	config        Config
	client        *http.Client
	mux           *http.ServeMux
	metrics       *Metrics
	registry      *prometheus.Registry
	proxy         *httputil.ReverseProxy
	shutdownTrace func(context.Context) error
}

func New(config Config) (*Server, error) {
	upstream, err := url.Parse(config.UpstreamURL)
	if err != nil {
		return nil, fmt.Errorf("invalid Ollama upstream: %w", err)
	}
	registry := prometheus.NewRegistry()
	s := &Server{config: config, client: &http.Client{Timeout: config.Timeout}, mux: http.NewServeMux(), registry: registry, metrics: newMetrics(registry), shutdownTrace: func(context.Context) error { return nil }}
	s.proxy = httputil.NewSingleHostReverseProxy(upstream)
	originalDirector := s.proxy.Director
	s.proxy.Director = func(request *http.Request) {
		originalDirector(request)
		request.URL.Path = "/" + strings.TrimPrefix(request.URL.Path, "/ollama/")
		request.Host = upstream.Host
	}
	s.proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, proxyErr error) {
		http.Error(writer, "Ollama proxy request failed: "+proxyErr.Error(), http.StatusBadGateway)
	}
	if config.TraceEnabled && config.OTLPEndpoint != "" {
		shutdown, err := configureTracing(context.Background(), config)
		if err != nil {
			return nil, err
		}
		s.shutdownTrace = shutdown
	}
	s.routes()
	return s, nil
}

func configureTracing(ctx context.Context, config Config) (func(context.Context) error, error) {
	endpoint := strings.TrimPrefix(strings.TrimPrefix(config.OTLPEndpoint, "http://"), "https://")
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(endpoint), otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	res, err := resource.Merge(resource.Default(), resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(config.ServiceName)))
	if err != nil {
		return nil, err
	}
	provider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(res))
	otel.SetTracerProvider(provider)
	return provider.Shutdown, nil
}

func (s *Server) routes() {
	s.mux.HandleFunc("/", s.root)
	s.mux.HandleFunc("/healthz", s.health)
	s.mux.HandleFunc("/config", s.configHandler)
	s.mux.Handle("/metrics", promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{}))
	s.mux.HandleFunc("/invoke", s.invoke)
	s.mux.HandleFunc("/ollama/", s.proxyHandler)
}

func (s *Server) Handler() http.Handler              { return s.observe(s.mux) }
func (s *Server) Shutdown(ctx context.Context) error { return s.shutdownTrace(ctx) }

func (s *Server) Run(ctx context.Context) error {
	server := &http.Server{Addr: s.config.Address, Handler: s.Handler(), ReadHeaderTimeout: 10 * time.Second}
	errCh := make(chan error, 1)
	go func() { errCh <- server.ListenAndServe() }()
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = s.Shutdown(shutdownCtx)
		return server.Shutdown(shutdownCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

type responseStatus struct {
	http.ResponseWriter
	code int
}

func (w *responseStatus) WriteHeader(code int) { w.code = code; w.ResponseWriter.WriteHeader(code) }
func (w *responseStatus) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (s *Server) observe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/metrics" {
			next.ServeHTTP(writer, request)
			return
		}
		started := time.Now()
		route := request.URL.Path
		if strings.HasPrefix(route, "/ollama/") {
			route = "ollama_proxy"
		}
		ctx, span := otel.Tracer(s.config.ServiceName).Start(request.Context(), request.Method+" "+route)
		defer span.End()
		status := &responseStatus{ResponseWriter: writer, code: http.StatusOK}
		next.ServeHTTP(status, request.WithContext(ctx))
		span.SetAttributes(attribute.Int("http.response.status_code", status.code), attribute.String("http.request.method", request.Method), attribute.String("url.path", request.URL.Path))
		s.metrics.Requests.WithLabelValues(request.Method, route, strconv.Itoa(status.code)).Inc()
		s.metrics.Duration.WithLabelValues(s.config.Model, route).Observe(time.Since(started).Seconds())
	})
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}
func (s *Server) root(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}
	writeJSON(writer, 200, map[string]string{"name": "Edge Ollama Go Gateway", "health": "/healthz", "invoke": "/invoke", "config": "/config", "ollama_proxy": "/ollama/api/*"})
}
func (s *Server) health(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, 200, map[string]any{"status": "ok", "model": s.config.Model, "ollama_base_url": s.config.UpstreamURL, "ollama_upstream_base_url": s.config.UpstreamURL, "otel_traces_enabled": s.config.TraceEnabled, "otel_exporter_otlp_endpoint": s.config.OTLPEndpoint})
}
func (s *Server) configHandler(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, 200, map[string]any{"model": s.config.Model, "ollama_base_url": s.config.UpstreamURL, "ollama_upstream_base_url": s.config.UpstreamURL, "otel_service_name": s.config.ServiceName, "otel_exporter_otlp_endpoint": s.config.OTLPEndpoint})
}

type promptInput struct {
	Prompt string `json:"prompt"`
	System string `json:"system"`
}

func (s *Server) invoke(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var input promptInput
	if err := json.NewDecoder(io.LimitReader(request.Body, 1<<20)).Decode(&input); err != nil || strings.TrimSpace(input.Prompt) == "" {
		http.Error(writer, "prompt is required", http.StatusBadRequest)
		return
	}
	prompt := input.Prompt
	if input.System != "" {
		prompt = "System: " + input.System + "\n\nUser: " + input.Prompt
	}
	payload, _ := json.Marshal(map[string]any{"model": s.config.Model, "prompt": prompt, "stream": false})
	upstreamRequest, err := http.NewRequestWithContext(request.Context(), http.MethodPost, strings.TrimRight(s.config.UpstreamURL, "/")+"/api/generate", strings.NewReader(string(payload)))
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}
	upstreamRequest.Header.Set("Content-Type", "application/json")
	s.metrics.Active.WithLabelValues(s.config.Model, "invoke").Inc()
	defer s.metrics.Active.WithLabelValues(s.config.Model, "invoke").Dec()
	response, err := s.client.Do(upstreamRequest)
	if err != nil {
		http.Error(writer, "Ollama request failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 16<<20))
	if err != nil {
		http.Error(writer, err.Error(), 502)
		return
	}
	if response.StatusCode >= 400 {
		writer.Header().Set("Content-Type", response.Header.Get("Content-Type"))
		writer.WriteHeader(response.StatusCode)
		_, _ = writer.Write(body)
		return
	}
	var result struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(writer, "invalid Ollama response", 502)
		return
	}
	writeJSON(writer, 200, map[string]string{"response": result.Response, "model": s.config.Model, "ollama_base_url": s.config.UpstreamURL})
}

type firstWrite struct {
	http.ResponseWriter
	once    sync.Once
	observe func()
}

func (w *firstWrite) Write(data []byte) (int, error) {
	w.once.Do(w.observe)
	return w.ResponseWriter.Write(data)
}
func (w *firstWrite) WriteHeader(code int) { w.once.Do(w.observe); w.ResponseWriter.WriteHeader(code) }
func (w *firstWrite) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
func (s *Server) proxyHandler(writer http.ResponseWriter, request *http.Request) {
	started := time.Now()
	s.metrics.Active.WithLabelValues(s.config.Model, "ollama_proxy").Inc()
	defer s.metrics.Active.WithLabelValues(s.config.Model, "ollama_proxy").Dec()
	wrapped := &firstWrite{ResponseWriter: writer, observe: func() {
		s.metrics.TTFT.WithLabelValues(s.config.Model, "ollama_proxy").Observe(time.Since(started).Seconds())
	}}
	s.proxy.ServeHTTP(wrapped, request)
}

func LogConfig(config Config) {
	slog.Info("starting Go Ollama gateway", "address", config.Address, "upstream", config.UpstreamURL, "model", config.Model, "tracing", config.TraceEnabled)
}
