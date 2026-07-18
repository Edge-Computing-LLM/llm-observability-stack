package stack

import (
	"context"
	"testing"
)

func TestProfileValuesFile(t *testing.T) {
	tests := map[string]string{
		"":                     "values.geforce-940m-k3s.yaml",
		"geforce-940m-k3s":     "values.geforce-940m-k3s.yaml",
		"enterprise-pilot-k3s": "values.enterprise-pilot-k3s.yaml",
		"cpu-k3s":              "values.cpu-k3s.yaml",
		"custom.yaml":          "custom.yaml",
	}
	for input, want := range tests {
		got, ok := ProfileValuesFile(input)
		if !ok {
			t.Fatalf("ProfileValuesFile(%q) returned !ok", input)
		}
		if got != want {
			t.Fatalf("ProfileValuesFile(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestGPUProfile(t *testing.T) {
	if !GPUProfile("geforce-940m-k3s") {
		t.Fatalf("geforce profile should be a GPU profile")
	}
	if GPUProfile("cpu-k3s") {
		t.Fatalf("cpu profile should not be a GPU profile")
	}
}

func TestGeForceDefaultsUseQwenAndNativeGoBenchmark(t *testing.T) {
	opts := DefaultOptions()
	if opts.Model != "qwen-1-8b-chat-q4-k-m-local" {
		t.Fatalf("default model = %q, want Qwen local alias", opts.Model)
	}
	if !stringsContains(benchmarkCommand(opts), "llm-observability benchmark") {
		t.Fatalf("benchmark command must use the native Go CLI")
	}
}

func TestHelmInstallDisablesBaseLayerChartsForGPUProfiles(t *testing.T) {
	cmd := helmInstallCommand(DefaultOptions())
	for _, want := range []string{
		"values.geforce-940m-k3s.yaml",
		"gpu-operator.enabled=false",
		"nvidia-device-plugin.enabled=false",
		"dcgm-exporter.enabled=false",
	} {
		if !stringsContains(cmd, want) {
			t.Fatalf("install command missing %q in:\n%s", want, cmd)
		}
	}
}

func TestInstallWithBaseIsDeprecated(t *testing.T) {
	opts := DefaultOptions()
	opts.WithBase = true
	opts.DryRun = true
	err := Install(context.Background(), opts)
	if err == nil || !stringsContains(err.Error(), "edge-cli") {
		t.Fatalf("Install with --with-base error = %v, want edge-cli deprecation", err)
	}
}

func stringsContains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return sub == ""
}
