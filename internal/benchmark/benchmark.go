package benchmark

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

type Config struct {
	URL, Model, Prompt, Output string
	Runs, WarmupRuns           int
	Timeout                    time.Duration
}
type Result struct {
	Success         bool    `json:"success"`
	DurationSeconds float64 `json:"duration_seconds"`
	TTFTSeconds     float64 `json:"ttft_seconds"`
	PromptTokens    int     `json:"prompt_tokens"`
	GeneratedTokens int     `json:"generated_tokens"`
	TokensPerSecond float64 `json:"tokens_per_second"`
}
type Summary struct {
	SchemaVersion       int      `json:"schema_version"`
	CapturedAt          string   `json:"captured_at"`
	Model               string   `json:"model"`
	URL                 string   `json:"url"`
	Runs                int      `json:"runs"`
	WarmupRuns          int      `json:"warmup_runs"`
	TTFTP50Seconds      float64  `json:"ttft_p50_seconds"`
	TTFTP95Seconds      float64  `json:"ttft_p95_seconds"`
	DurationP50Seconds  float64  `json:"duration_p50_seconds"`
	DurationP95Seconds  float64  `json:"duration_p95_seconds"`
	TokensPerSecondMean float64  `json:"tokens_per_second_mean"`
	Results             []Result `json:"results"`
}

func Run(ctx context.Context, config Config) (*Summary, error) {
	if config.Runs < 1 {
		return nil, fmt.Errorf("runs must be positive")
	}
	if config.Timeout <= 0 {
		config.Timeout = 300 * time.Second
	}
	client := &http.Client{Timeout: config.Timeout}
	for i := 0; i < config.WarmupRuns; i++ {
		if _, err := runOnce(ctx, client, config); err != nil {
			return nil, fmt.Errorf("warmup %d: %w", i+1, err)
		}
	}
	results := make([]Result, 0, config.Runs)
	for i := 0; i < config.Runs; i++ {
		result, err := runOnce(ctx, client, config)
		if err != nil {
			return nil, fmt.Errorf("run %d: %w", i+1, err)
		}
		results = append(results, result)
	}
	ttft, duration := make([]float64, 0, len(results)), make([]float64, 0, len(results))
	throughput := 0.0
	for _, result := range results {
		ttft = append(ttft, result.TTFTSeconds)
		duration = append(duration, result.DurationSeconds)
		throughput += result.TokensPerSecond
	}
	return &Summary{1, time.Now().UTC().Format(time.RFC3339Nano), config.Model, config.URL, config.Runs, config.WarmupRuns, Percentile(ttft, .50), Percentile(ttft, .95), Percentile(duration, .50), Percentile(duration, .95), throughput / float64(len(results)), results}, nil
}

func runOnce(ctx context.Context, client *http.Client, config Config) (Result, error) {
	payload, _ := json.Marshal(map[string]any{"model": config.Model, "prompt": config.Prompt, "stream": true, "options": map[string]any{"temperature": 0}})
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, config.URL, bytes.NewReader(payload))
	if err != nil {
		return Result{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	started := time.Now()
	response, err := client.Do(request)
	if err != nil {
		return Result{}, err
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		detail, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		return Result{}, fmt.Errorf("Ollama status %d: %s", response.StatusCode, detail)
	}
	firstTokenAt := time.Time{}
	final := struct {
		Done            bool   `json:"done"`
		Response        string `json:"response"`
		EvalCount       int    `json:"eval_count"`
		PromptEvalCount int    `json:"prompt_eval_count"`
		EvalDuration    int64  `json:"eval_duration"`
	}{}
	scanner := bufio.NewScanner(response.Body)
	scanner.Buffer(make([]byte, 64*1024), 2<<20)
	for scanner.Scan() {
		var chunk struct {
			Done            bool   `json:"done"`
			Response        string `json:"response"`
			EvalCount       int    `json:"eval_count"`
			PromptEvalCount int    `json:"prompt_eval_count"`
			EvalDuration    int64  `json:"eval_duration"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
			return Result{}, err
		}
		if firstTokenAt.IsZero() && chunk.Response != "" {
			firstTokenAt = time.Now()
		}
		if chunk.Done {
			final = chunk
		}
	}
	if err := scanner.Err(); err != nil {
		return Result{}, err
	}
	finished := time.Now()
	if firstTokenAt.IsZero() {
		firstTokenAt = finished
	}
	evalSeconds := float64(final.EvalDuration) / 1e9
	throughput := 0.0
	if evalSeconds > 0 {
		throughput = float64(final.EvalCount) / evalSeconds
	}
	return Result{true, finished.Sub(started).Seconds(), firstTokenAt.Sub(started).Seconds(), final.PromptEvalCount, final.EvalCount, throughput}, nil
}

func Percentile(values []float64, quantile float64) float64 {
	if len(values) == 0 {
		return 0
	}
	ordered := append([]float64(nil), values...)
	sort.Float64s(ordered)
	index := int(float64(len(ordered)-1)*quantile + .5)
	if index < 0 {
		index = 0
	}
	if index >= len(ordered) {
		index = len(ordered) - 1
	}
	return ordered[index]
}
func Write(path string, summary *Summary) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func WithPortForward(ctx context.Context, namespace string, readyTimeout time.Duration, run func() error) error {
	command := exec.CommandContext(ctx, "kubectl", "port-forward", "-n", namespace, "svc/ollama", "11434:11434")
	var logs bytes.Buffer
	command.Stdout, command.Stderr = &logs, &logs
	if err := command.Start(); err != nil {
		return err
	}
	defer func() { _ = command.Process.Kill(); _ = command.Wait() }()
	deadline := time.Now().Add(readyTimeout)
	for time.Now().Before(deadline) {
		request, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1:11434/api/tags", nil)
		response, err := http.DefaultClient.Do(request)
		if err == nil {
			_ = response.Body.Close()
			if response.StatusCode < 500 {
				return run()
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
	return fmt.Errorf("port-forward did not become ready: %s", logs.String())
}
