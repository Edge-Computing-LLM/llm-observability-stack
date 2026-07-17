package stack

import (
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	Profile       string
	Namespace     string
	Release       string
	ValuesFiles   []string
	SetValues     []string
	WithBase      bool
	SkipBase      bool
	DryRun        bool
	Yes           bool
	Timeout       string
	KeepNamespace bool
	Model         string
	Runs          int
	Prompt        string
	Output        string
	OllamaSmoke   bool
}

func DefaultOptions() Options {
	return Options{
		Profile:     "geforce-940m-k3s",
		Namespace:   "llm-observability",
		Release:     "llm-observability-stack",
		Timeout:     "5m",
		Model:       "qwen-1-8b-chat-q4-k-m-local",
		Runs:        3,
		Prompt:      "Explain GPU observability in one concise sentence.",
		Output:      "artifacts/benchmark-local.json",
		OllamaSmoke: true,
	}
}

func RepoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "Chart.yaml")); err == nil {
			if _, err := os.Stat(filepath.Join(wd, "templates")); err == nil {
				return wd, nil
			}
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return "", os.ErrNotExist
		}
		wd = parent
	}
}

func ProfileValuesFile(profile string) (string, bool) {
	if profile == "" {
		profile = DefaultOptions().Profile
	}
	if strings.HasSuffix(profile, ".yaml") || strings.Contains(profile, "/") {
		return profile, true
	}
	profiles := map[string]string{
		"default":              "values.yaml",
		"local-k3s":            "values.local-k3s.yaml",
		"local-k3s-example":    "values.local-k3s.example.yaml",
		"enterprise-pilot-k3s": "values.enterprise-pilot-k3s.yaml",
		"geforce-940m-k3s":     "values.geforce-940m-k3s.yaml",
		"cpu-k3s":              "values.cpu-k3s.yaml",
		"validation-k3s":       "values.validation-k3s.yaml",
		"full-stack-nvidia":    "values.full-stack-nvidia.example.yaml",
	}
	file, ok := profiles[profile]
	return file, ok
}

func GPUProfile(profile string) bool {
	file, _ := ProfileValuesFile(profile)
	name := strings.ToLower(profile + " " + file)
	return !strings.Contains(name, "cpu")
}
