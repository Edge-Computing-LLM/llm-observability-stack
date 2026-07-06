package stack

import "testing"

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

func stringsContains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return sub == ""
}
