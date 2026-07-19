package stack

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Edge-Computing-LLM/k3s-nvidia-edge/pkg/edgebase"
	"github.com/Edge-Computing-LLM/llm-observability-stack/internal/benchmark"
)

func Doctor(ctx context.Context, opts Options) error {
	r := runner(opts)
	if err := baseReady(ctx, r, opts); err != nil {
		return err
	}
	for _, step := range stackDoctorSteps(opts) {
		if err := r.Run(ctx, step); err != nil {
			return err
		}
	}
	return nil
}

func Install(ctx context.Context, opts Options) error {
	if opts.WithBase && opts.SkipBase {
		return fmt.Errorf("--with-base and --skip-base cannot both be set")
	}
	r := runner(opts)
	if opts.WithBase {
		return fmt.Errorf("--with-base is deprecated for llm-observability-stack. Use edge-cli and run edge install infra before edge install observability")
	} else if GPUProfile(opts.Profile) {
		if err := baseReady(ctx, r, opts); err != nil {
			return fmt.Errorf("base NVIDIA layer is not ready; run edge install infra or validate k3s-nvidia-edge first: %w", err)
		}
	}
	for _, step := range installSteps(opts) {
		if err := r.Run(ctx, step); err != nil {
			return err
		}
	}
	return nil
}

func Status(ctx context.Context, opts Options) error {
	r := runner(opts)
	for _, step := range statusSteps(opts) {
		if err := r.Run(ctx, step); err != nil {
			return err
		}
	}
	return nil
}

func Validate(ctx context.Context, opts Options) error {
	r := runner(opts)
	if GPUProfile(opts.Profile) {
		baseOpts := baseOptions(opts)
		baseRunner := edgebase.NewRunner(baseOpts)
		if err := edgebase.Validate(ctx, baseRunner, baseOpts); err != nil {
			return err
		}
	}
	for _, step := range validateSteps(opts) {
		if err := r.Run(ctx, step); err != nil {
			return err
		}
	}
	return nil
}

func Benchmark(ctx context.Context, opts Options) error {
	if err := runner(opts).Run(ctx, edgebase.Step{Name: "Wait for Ollama", Command: fmt.Sprintf("kubectl rollout status deploy/ollama -n %s --timeout=%s", shellQuote(opts.Namespace), shellQuote(opts.Timeout))}); err != nil {
		return err
	}
	return benchmark.WithPortForward(ctx, opts.Namespace, 15*time.Second, func() error {
		summary, err := benchmark.Run(ctx, benchmark.Config{URL: "http://127.0.0.1:11434/api/generate", Model: opts.Model, Prompt: opts.Prompt, Output: opts.Output, Runs: opts.Runs, WarmupRuns: opts.WarmupRuns, Timeout: 300 * time.Second})
		if err != nil {
			return err
		}
		if err := benchmark.Write(opts.Output, summary); err != nil {
			return err
		}
		fmt.Printf("benchmark written to %s\n", opts.Output)
		return nil
	})
}

func NetworkInventory(ctx context.Context, opts Options) error {
	return runner(opts).Run(ctx, edgebase.Step{Name: "Kubernetes network inventory", Command: fmt.Sprintf("kubectl get pods,services,endpoints,endpointslices.discovery.k8s.io,networkpolicies.networking.k8s.io -n %s -o wide", shellQuote(opts.Namespace))})
}

func ServicePath(ctx context.Context, opts Options) error {
	if strings.TrimSpace(opts.Service) == "" {
		return fmt.Errorf("--service is required")
	}
	ns, service := shellQuote(opts.Namespace), shellQuote(opts.Service)
	command := fmt.Sprintf(`set -e
kubectl get service -n %s %s -o wide
selector="$(kubectl get service -n %s %s -o go-template='{{range $key,$value := .spec.selector}}{{printf "%%s=%%s," $key $value}}{{end}}' | sed 's/,$//')"
if [ -n "$selector" ]; then kubectl get pods -n %s -l "$selector" -o wide; else echo "Service has no selector"; fi
kubectl get endpoints -n %s %s -o wide
kubectl get endpointslices.discovery.k8s.io -n %s -l kubernetes.io/service-name=%s -o wide`, ns, service, ns, service, ns, ns, service, ns, service)
	return runner(opts).Run(ctx, edgebase.Step{Name: "Kubernetes service path", Command: command})
}

func WatchEndpoints(ctx context.Context, opts Options) error {
	if strings.TrimSpace(opts.Service) == "" {
		return fmt.Errorf("--service is required")
	}
	return runner(opts).Run(ctx, edgebase.Step{Name: "Watch Kubernetes endpoints", Command: fmt.Sprintf("kubectl get endpoints -n %s %s --watch --request-timeout=%s", shellQuote(opts.Namespace), shellQuote(opts.Service), shellQuote(opts.Timeout))})
}

func Uninstall(ctx context.Context, opts Options) error {
	r := runner(opts)
	for _, step := range uninstallSteps(opts) {
		if err := r.Run(ctx, step); err != nil {
			return err
		}
	}
	if opts.WithBase {
		return fmt.Errorf("--with-base is deprecated for llm-observability-stack. Use edge-cli uninstall all for reverse-order layer removal")
	}
	return nil
}

func PrintCommands(opts Options) error {
	fmt.Println("# Base readiness")
	for _, step := range baseReadySteps(opts) {
		fmt.Printf("\n## %s\n%s\n", step.Name, step.Command)
	}
	fmt.Println("\n# Install")
	for _, step := range installSteps(opts) {
		fmt.Printf("\n## %s\n%s\n", step.Name, step.Command)
	}
	fmt.Println("\n# Validate")
	for _, step := range validateSteps(opts) {
		fmt.Printf("\n## %s\n%s\n", step.Name, step.Command)
	}
	fmt.Println("\n# Benchmark")
	fmt.Println(benchmarkCommand(opts))
	fmt.Println("\n# Uninstall")
	for _, step := range uninstallSteps(opts) {
		fmt.Printf("\n## %s\n%s\n", step.Name, step.Command)
	}
	if opts.WithBase {
		fmt.Println("\n# Base uninstall")
		baseOpts := baseOptions(opts)
		edgebase.PrintCommands(baseOpts)
	}
	return nil
}

func runner(opts Options) *edgebase.Runner {
	base := edgebase.DefaultOptions()
	base.Yes = opts.Yes && !opts.DryRun
	base.RequireHostCUDA = false
	return edgebase.NewRunner(base)
}

func baseOptions(opts Options) edgebase.Options {
	base := edgebase.DefaultOptions()
	base.Yes = opts.Yes && !opts.DryRun
	base.RequireHostCUDA = false
	base.UseLocalChart = true
	base.LocalChartPath = "../k3s-nvidia-edge/charts/k3s-nvidia-edge"
	return base
}

func baseReady(ctx context.Context, r *edgebase.Runner, opts Options) error {
	for _, step := range baseReadySteps(opts) {
		if err := r.Run(ctx, step); err != nil {
			return err
		}
	}
	return nil
}

func baseReadySteps(opts Options) []edgebase.Step {
	steps := []edgebase.Step{
		{Name: "Kubernetes connectivity", Command: "kubectl cluster-info"},
		{Name: "k3s node address", Command: edgebase.NodeAddressHealthCheck()},
		{Name: "k3s nodes", Command: "kubectl get nodes -o wide"},
	}
	if GPUProfile(opts.Profile) {
		steps = append(steps,
			edgebase.Step{Name: "NVIDIA RuntimeClass", Command: "kubectl get runtimeclass nvidia"},
			edgebase.Step{Name: "NVIDIA GPU allocatable", Command: edgebase.GPUCapacityCheck()},
			edgebase.Step{Name: "GPU Operator health", Command: edgebase.GPUOperatorHealthCheck()},
		)
	}
	return steps
}

func stackDoctorSteps(opts Options) []edgebase.Step {
	ns := shellQuote(opts.Namespace)
	steps := []edgebase.Step{
		{Name: "Required commands", Command: "missing=0; for c in kubectl helm; do command -v $c >/dev/null && echo \"$c: $(command -v $c)\" || { echo \"$c: missing\"; missing=1; }; done; exit $missing"},
		{Name: "Helm release", Command: fmt.Sprintf("helm status %s -n %s || true", shellQuote(opts.Release), ns)},
		{Name: "LLM namespace", Command: fmt.Sprintf("kubectl get namespace %s || true", ns)},
		{Name: "LLM workloads", Command: fmt.Sprintf("kubectl get pods,deploy,statefulset,svc,pvc -n %s -o wide || true", ns)},
		{Name: "Ollama service", Command: fmt.Sprintf("kubectl get svc -n %s ollama || true", ns)},
		{Name: "Open WebUI service", Command: fmt.Sprintf("kubectl get svc -n %s open-webui || true", ns)},
		{Name: "OpenTelemetry Collector service", Command: fmt.Sprintf("kubectl get svc -n %s opentelemetry-collector || true", ns)},
	}
	if GPUProfile(opts.Profile) {
		steps = append(steps, edgebase.Step{Name: "Ollama NVIDIA scheduling", Command: fmt.Sprintf("kubectl get deploy -n %s ollama -o jsonpath='{.spec.template.spec.runtimeClassName}{\"\\n\"}{.spec.template.spec.containers[0].resources.limits.nvidia\\.com/gpu}{\"\\n\"}' || true", ns)})
	}
	return steps
}

func installSteps(opts Options) []edgebase.Step {
	return []edgebase.Step{
		{Name: "Apply optional monitoring CRDs", Mutating: true, Command: withRoot(`if ls charts/kube-prometheus-stack/charts/crds/crds/*.yaml >/dev/null 2>&1; then
  for crd in charts/kube-prometheus-stack/charts/crds/crds/*.yaml; do
    kubectl create --save-config=false -f "$crd" 2>/dev/null || kubectl apply --server-side -f "$crd"
  done
fi`)},
		{Name: "Helm dependency build", Mutating: true, Command: withRoot("helm dependency build .")},
		{Name: "Install llm-observability-stack", Mutating: true, Command: withRoot(helmInstallCommand(opts))},
		{Name: "Wait for Ollama", Mutating: true, Command: fmt.Sprintf("kubectl rollout status deploy/ollama -n %s --timeout=%s", shellQuote(opts.Namespace), shellQuote(opts.Timeout))},
		{Name: "Wait for Open WebUI", Mutating: true, Command: fmt.Sprintf("kubectl rollout status statefulset/open-webui -n %s --timeout=%s", shellQuote(opts.Namespace), shellQuote(opts.Timeout))},
		{Name: "Wait for OpenTelemetry Collector", Mutating: true, Command: fmt.Sprintf("kubectl rollout status deploy/opentelemetry-collector -n %s --timeout=%s || true", shellQuote(opts.Namespace), shellQuote(opts.Timeout))},
	}
}

func statusSteps(opts Options) []edgebase.Step {
	ns := shellQuote(opts.Namespace)
	return []edgebase.Step{
		{Name: "Base GPU layer", Command: "kubectl get pods -n gpu-operator -o wide || true\nkubectl get runtimeclass nvidia || true\nkubectl get nodes -o custom-columns=NAME:.metadata.name,GPU:.status.allocatable.nvidia\\\\.com/gpu || true"},
		{Name: "Helm release", Command: fmt.Sprintf("helm status %s -n %s || true", shellQuote(opts.Release), ns)},
		{Name: "LLM workloads", Command: fmt.Sprintf("kubectl get pods,deploy,statefulset,svc,pvc -n %s -o wide", ns)},
		{Name: "Services and ports", Command: fmt.Sprintf("kubectl get svc -n %s -o wide", ns)},
		{Name: "Stale or unknown pods", Command: fmt.Sprintf("kubectl get pods -n %s --no-headers | awk '$3==\"Unknown\" || $3==\"Failed\" || $3==\"Error\" || $3==\"CrashLoopBackOff\" {print}' || true", ns)},
		{Name: "Ollama loaded models", Command: fmt.Sprintf("kubectl exec -n %s deploy/ollama -- ollama ps || true", ns)},
	}
}

func validateSteps(opts Options) []edgebase.Step {
	ns := shellQuote(opts.Namespace)
	steps := []edgebase.Step{
		{Name: "Helm release deployed", Command: fmt.Sprintf("helm status %s -n %s", shellQuote(opts.Release), ns)},
		{Name: "Expected services", Command: fmt.Sprintf("kubectl get svc -n %s ollama open-webui opentelemetry-collector", ns)},
		{Name: "Ollama rollout", Command: fmt.Sprintf("kubectl rollout status deploy/ollama -n %s --timeout=%s", ns, shellQuote(opts.Timeout))},
		{Name: "Open WebUI rollout", Command: fmt.Sprintf("kubectl rollout status statefulset/open-webui -n %s --timeout=%s", ns, shellQuote(opts.Timeout))},
		{Name: "Open WebUI Redis if enabled", Command: fmt.Sprintf("if kubectl get deploy -n %s open-webui-redis >/dev/null 2>&1; then kubectl rollout status deploy/open-webui-redis -n %s --timeout=%s; fi", ns, ns, shellQuote(opts.Timeout))},
		{Name: "OpenTelemetry Collector rollout", Command: fmt.Sprintf("kubectl rollout status deploy/opentelemetry-collector -n %s --timeout=%s", ns, shellQuote(opts.Timeout))},
		{Name: "All pods ready", Command: fmt.Sprintf("kubectl wait --for=condition=Ready pod --all -n %s --timeout=%s", ns, shellQuote(opts.Timeout))},
		{Name: "GGUF model mount", Command: fmt.Sprintf("kubectl exec -n %s deploy/ollama -- sh -c 'ls -lh /models/gguf || true'", ns)},
		{Name: "Ollama models", Command: fmt.Sprintf("kubectl exec -n %s deploy/ollama -- ollama list", ns)},
		{Name: "Ollama residency", Command: fmt.Sprintf("kubectl exec -n %s deploy/ollama -- ollama ps", ns)},
	}
	if opts.OllamaSmoke {
		steps = append(steps, edgebase.Step{Name: "Ollama smoke test", Command: fmt.Sprintf("kubectl exec -n %s deploy/ollama -- ollama run %s %s", ns, shellQuote(opts.Model), shellQuote("Reply with exactly: validation ok"))})
	}
	if GPUProfile(opts.Profile) {
		steps = append(steps, edgebase.Step{Name: "Ollama CUDA evidence", Command: fmt.Sprintf("kubectl logs -n %s deploy/ollama --tail=-1 | grep -Ei 'CUDA|offload|model weights|gpu memory'", ns)})
	}
	return steps
}

func uninstallSteps(opts Options) []edgebase.Step {
	steps := []edgebase.Step{
		{Name: "Uninstall llm-observability-stack", Mutating: true, Command: fmt.Sprintf("helm uninstall %s -n %s --wait || true", shellQuote(opts.Release), shellQuote(opts.Namespace))},
	}
	if !opts.KeepNamespace {
		steps = append(steps, edgebase.Step{Name: "Delete LLM namespace", Mutating: true, Command: fmt.Sprintf("kubectl delete namespace %s --ignore-not-found", shellQuote(opts.Namespace))})
	}
	return steps
}

func helmInstallCommand(opts Options) string {
	valuesFile, ok := ProfileValuesFile(opts.Profile)
	if !ok {
		valuesFile = opts.Profile
	}
	args := []string{
		"helm upgrade --install",
		shellQuote(opts.Release),
		".",
		"-n", shellQuote(opts.Namespace),
		"--create-namespace",
		"-f", shellQuote(valuesFile),
		"--wait",
		"--timeout", shellQuote(opts.Timeout),
	}
	for _, file := range opts.ValuesFiles {
		args = append(args, "-f", shellQuote(file))
	}
	for _, set := range opts.SetValues {
		args = append(args, "--set", shellQuote(set))
	}
	if GPUProfile(opts.Profile) {
		args = append(args,
			"--set", "gpu-operator.enabled=false",
			"--set", "nvidia-device-plugin.enabled=false",
			"--set", "dcgm-exporter.enabled=false",
		)
	}
	return strings.Join(args, " ")
}

func benchmarkCommand(opts Options) string {
	return fmt.Sprintf("llm-observability benchmark --namespace %s --model %s --runs %d --prompt %s --output %s", shellQuote(opts.Namespace), shellQuote(opts.Model), opts.Runs, shellQuote(opts.Prompt), shellQuote(opts.Output))
}

func withRoot(command string) string {
	return fmt.Sprintf("cd %s\n%s", shellQuote(mustRepoRoot()), command)
}

func mustRepoRoot() string {
	root, err := RepoRoot()
	if err != nil {
		return "."
	}
	return root
}
