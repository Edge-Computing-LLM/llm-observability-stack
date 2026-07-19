package tests

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func root(t *testing.T) string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate test source")
	}
	return filepath.Dir(filepath.Dir(file))
}
func helm(t *testing.T, args ...string) (string, error) {
	t.Helper()
	if _, err := exec.LookPath("helm"); err != nil {
		t.Skip("helm is unavailable")
	}
	command := exec.Command("helm", args...)
	command.Dir = root(t)
	output, err := command.CombinedOutput()
	return string(output), err
}
func requireContains(t *testing.T, value string, expected ...string) {
	t.Helper()
	for _, item := range expected {
		if !strings.Contains(value, item) {
			t.Errorf("render is missing %q", item)
		}
	}
}
func requireAbsent(t *testing.T, value string, forbidden ...string) {
	t.Helper()
	for _, item := range forbidden {
		if strings.Contains(value, item) {
			t.Errorf("render unexpectedly contains %q", item)
		}
	}
}

func TestLocalProfileUsesNativeGoWorkloads(t *testing.T) {
	manifest, err := helm(t, "template", "llm-observability-stack", ".", "-f", "values.local-k3s.example.yaml")
	if err != nil {
		t.Fatal(manifest)
	}
	requireContains(t, manifest, "name: ollama", "name: open-webui", "name: ollama-gateway", "name: edge-toolbox", "/usr/local/bin/ollama-gateway", "/usr/local/bin/edge-toolbox")
	requireAbsent(t, manifest, "python", "uvicorn", "app.py", "langchain")
}

func TestGeForceProfileContract(t *testing.T) {
	manifest, err := helm(t, "template", "llm-observability-stack", ".", "-f", "values.geforce-940m-k3s.yaml")
	if err != nil {
		t.Fatal(manifest)
	}
	requireContains(t, manifest, "qwen-1.8b-chat-q4_K_M.gguf", "nvidia.com/gpu.present: \"true\"", "nvidia.com/gpu: 1", "OLLAMA_KEEP_ALIVE", "PARAMETER num_gpu 23", "PARAMETER num_ctx 256", "PARAMETER num_batch 1", "helm.sh/resource-policy: keep", "name: open-webui", "name: open-webui-redis", "name: llm-observability-dashboards", "edge-llm-observability.json", "DCGM_FI_DEV_MEM_COPY_UTIL", "kind: ServiceMonitor")
	requireAbsent(t, manifest, "name: ollama-gateway", "name: edge-toolbox", "/bin/ollama rm", "kind: ClusterPolicy")
}

func TestGeForceModelOverlays(t *testing.T) {
	tests := []struct {
		name       string
		valuesFile string
		expected   []string
		absent     []string
	}{
		{
			name:       "gemma local GGUF",
			valuesFile: "values.gemma-3-1b-geforce-940m-k3s.yaml",
			expected:   []string{"gemma3-1b-it-gguf-local", "FROM /models/gguf/gemma-3-1b-it-Q4_K_M.gguf", "PARAMETER num_gpu 23", "PARAMETER num_ctx 256", "PARAMETER num_batch 1"},
		},
		{
			name:       "llama registry model",
			valuesFile: "values.llama3.2-1b-geforce-940m-k3s.yaml",
			expected:   []string{"llama3-2-1b-local", "llama3.2:1b", "FROM llama3.2:1b", "PARAMETER num_gpu 8", "PARAMETER num_ctx 256", "PARAMETER num_batch 1"},
			absent:     []string{"FROM /models/gguf/qwen-1.8b-chat-q4_K_M.gguf"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			manifest, err := helm(t, "template", "llm-observability-stack", ".", "-f", "values.geforce-940m-k3s.yaml", "-f", test.valuesFile)
			if err != nil {
				t.Fatal(manifest)
			}
			requireContains(t, manifest, test.expected...)
			requireAbsent(t, manifest, test.absent...)
		})
	}
}

func TestCPUProfileHasNoNVIDIAScheduling(t *testing.T) {
	manifest, err := helm(t, "template", "llm-observability-stack", ".", "-f", "values.cpu-k3s.yaml")
	if err != nil {
		t.Fatal(manifest)
	}
	requireContains(t, manifest, "name: ollama", "name: ollama-gateway", "/usr/local/bin/ollama-gateway")
	requireAbsent(t, manifest, "runtimeClassName: nvidia", "runtimeClassName: \"nvidia\"", "nvidia.com/gpu: 1", "nvidia.com/gpu.present", "python")
}

func TestDestructiveModelCleanupIsRejected(t *testing.T) {
	output, err := helm(t, "template", "llm-observability-stack", ".", "--set", "ollama.ollama.models.clean=true")
	if err == nil || !strings.Contains(output, "models.clean must stay false") {
		t.Fatalf("expected cleanup rejection, err=%v output=%s", err, output)
	}
}
func TestBaseLayerChartsAreRejected(t *testing.T) {
	for _, key := range []string{"gpu-operator.enabled=true", "nvidia-device-plugin.enabled=true", "dcgm-exporter.enabled=true"} {
		output, err := helm(t, "template", "llm-observability-stack", ".", "--set", key)
		if err == nil || !strings.Contains(output, "k3s-nvidia-edge") {
			t.Fatalf("expected %s rejection, err=%v output=%s", key, err, output)
		}
	}
}
func TestSecretMismatchIsRejected(t *testing.T) {
	output, err := helm(t, "template", "llm-observability-stack", ".", "--set", "openWebUI.existingSecret=legacy-secret", "--set", "open-webui.webuiSecret.existingSecretName=subchart-secret")
	if err == nil || !strings.Contains(output, "Secret name mismatch") {
		t.Fatalf("expected secret rejection, err=%v output=%s", err, output)
	}
}

func TestChartPackageIsSmallAndExcludesDevelopmentSources(t *testing.T) {
	directory := t.TempDir()
	output, err := helm(t, "package", ".", "-d", directory)
	if err != nil {
		t.Fatal(output)
	}
	matches, _ := filepath.Glob(filepath.Join(directory, "llm-observability-stack-*.tgz"))
	if len(matches) != 1 {
		t.Fatalf("package matches: %#v", matches)
	}
	info, err := os.Stat(matches[0])
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() >= 3_000_000 {
		t.Fatalf("chart package is %d bytes", info.Size())
	}
	file, err := os.Open(matches[0])
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	compressed, err := gzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}
	defer compressed.Close()
	archive := tar.NewReader(compressed)
	names := map[string]bool{}
	for {
		header, err := archive.Next()
		if err != nil {
			break
		}
		names[header.Name] = true
	}
	for _, forbidden := range []string{"llm-observability-stack/.git/", "llm-observability-stack/docs/", "llm-observability-stack/tests/", "llm-observability-stack/internal/", "llm-observability-stack/ollama-gateway/", "llm-observability-stack/edge-toolbox/"} {
		for name := range names {
			if strings.HasPrefix(name, forbidden) {
				t.Errorf("package contains %s", name)
			}
		}
	}
	for _, required := range []string{"llm-observability-stack/templates/ollama-gateway-deployment.yaml", "llm-observability-stack/templates/edge-toolbox-deployment.yaml", "llm-observability-stack/dashboards/edge-llm-observability.json"} {
		if !names[required] {
			t.Errorf("package is missing %s", required)
		}
	}
}
