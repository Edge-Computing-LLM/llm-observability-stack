package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Edge-Computing-LLM/llm-observability-stack/internal/stack"
)

const usage = `llm-observability manages the LLM application and observability layer on k3s.

Usage:
  llm-observability <command> [flags]

Commands:
  doctor          Check k3s/NVIDIA base readiness and LLM stack prerequisites
  install         Install or upgrade the llm-observability-stack Helm release
  status          Print base and LLM stack status
  validate        Run deeper release, pod, service, and optional Ollama checks
  benchmark       Run the existing Python Ollama benchmark client
  uninstall       Uninstall the LLM stack, optionally including the base layer
  print-commands  Print Helm/kubectl commands used by install/validate/uninstall

Common flags:
  --profile NAME       profile name or values file (default "geforce-940m-k3s")
  --namespace NAME     Kubernetes namespace (default "llm-observability")
  --release NAME       Helm release name (default "llm-observability-stack")
  --values FILE        additional values file; may be repeated
  --set KEY=VALUE      additional Helm --set override; may be repeated
  --with-base          deprecated; use edge-cli to install or uninstall the base layer
  --skip-base          do not install base, but still check it for GPU profiles
  --dry-run            print mutating commands without executing them
  --yes                execute mutating commands
  --timeout DURATION   rollout/Helm timeout (default 5m)

Examples:
  llm-observability doctor
  llm-observability install --profile geforce-940m-k3s --skip-base --yes
  llm-observability status
  llm-observability validate
  llm-observability benchmark --model qwen-1-8b-chat-q4-k-m-local --runs 3
  llm-observability uninstall --yes
  edge uninstall all --yes
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	cmd := strings.ToLower(os.Args[1])
	opts := stack.DefaultOptions()
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	var values repeatedFlag
	var sets repeatedFlag
	fs.StringVar(&opts.Profile, "profile", opts.Profile, "profile name or values file")
	fs.StringVar(&opts.Namespace, "namespace", opts.Namespace, "Kubernetes namespace")
	fs.StringVar(&opts.Release, "release", opts.Release, "Helm release name")
	fs.Var(&values, "values", "additional Helm values file; may be repeated")
	fs.Var(&sets, "set", "additional Helm --set override; may be repeated")
	fs.BoolVar(&opts.WithBase, "with-base", false, "deprecated; use edge-cli for k3s-nvidia-edge base operations")
	fs.BoolVar(&opts.SkipBase, "skip-base", false, "skip base installation")
	fs.BoolVar(&opts.DryRun, "dry-run", false, "print mutating commands without executing them")
	fs.BoolVar(&opts.Yes, "yes", false, "execute mutating commands")
	fs.StringVar(&opts.Timeout, "timeout", opts.Timeout, "rollout/Helm timeout")
	fs.BoolVar(&opts.KeepNamespace, "keep-namespace", false, "with uninstall: keep namespace")
	fs.StringVar(&opts.Model, "model", opts.Model, "with benchmark/validate: Ollama model")
	fs.IntVar(&opts.Runs, "runs", opts.Runs, "with benchmark: benchmark runs")
	fs.StringVar(&opts.Prompt, "prompt", opts.Prompt, "with benchmark: benchmark prompt")
	fs.StringVar(&opts.Output, "output", opts.Output, "with benchmark: output JSON path")
	fs.BoolVar(&opts.OllamaSmoke, "ollama-smoke", opts.OllamaSmoke, "with validate: run Ollama smoke test")

	if err := fs.Parse(os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	opts.ValuesFiles = values
	opts.SetValues = sets

	if cmd == "help" || cmd == "-h" || cmd == "--help" {
		fmt.Print(usage)
		return
	}

	timeout, err := time.ParseDuration(opts.Timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid --timeout: %v\n", err)
		os.Exit(2)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout+30*time.Second)
	defer cancel()

	var runErr error
	switch cmd {
	case "doctor":
		runErr = stack.Doctor(ctx, opts)
	case "install":
		runErr = stack.Install(ctx, opts)
	case "status":
		runErr = stack.Status(ctx, opts)
	case "validate":
		runErr = stack.Validate(ctx, opts)
	case "benchmark":
		runErr = stack.Benchmark(ctx, opts)
	case "uninstall":
		runErr = stack.Uninstall(ctx, opts)
	case "print-commands":
		runErr = stack.PrintCommands(opts)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s", cmd, usage)
		os.Exit(2)
	}
	if runErr != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", runErr)
		os.Exit(1)
	}
}

type repeatedFlag []string

func (f *repeatedFlag) String() string {
	return strings.Join(*f, ",")
}

func (f *repeatedFlag) Set(value string) error {
	*f = append(*f, value)
	return nil
}
