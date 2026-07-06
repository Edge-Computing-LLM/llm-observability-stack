import json
import os
import time
from typing import Any

import requests
from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.trace import Span, Status, StatusCode

DEFAULT_API_URL = "http://ollama:11434/api/chat"
DEFAULT_MODEL = "gemma3-1b-it-gguf-local"
DEFAULT_TIMEOUT_SECONDS = 180
DEFAULT_CALL_COUNT = 12
DEFAULT_OTLP_ENDPOINT = "http://opentelemetry-collector:4317"
TRACE_PREVIEW_LIMIT = 1500

DEFAULT_PROMPTS = [
    "Explain what Kubernetes readiness probe means in one short paragraph.",
    "Give one example where GPU acceleration helps inference latency.",
    "What should an OpenTelemetry GenAI span contain for local LLM inference?",
    "Write three bullet points for debugging Ollama in k3s.",
]


def truncate(value: str, limit: int = TRACE_PREVIEW_LIMIT) -> str:
    return value if len(value) <= limit else value[:limit]


def configure_tracing(endpoint: str) -> None:
    provider = TracerProvider(
        resource=Resource.create(
            {
                "service.name": os.getenv("OTEL_SERVICE_NAME", "otel-trace-seeder"),
                "service.namespace": os.getenv("OTEL_SERVICE_NAMESPACE", "llm-observability"),
                "deployment.environment": os.getenv("OTEL_DEPLOYMENT_ENVIRONMENT", "k3s-nvidia-edge"),
            }
        )
    )
    provider.add_span_processor(BatchSpanProcessor(OTLPSpanExporter(endpoint=endpoint)))
    trace.set_tracer_provider(provider)


def load_prompts() -> list[str]:
    raw_json = os.getenv("OBS_PROMPTS_JSON")
    if raw_json:
        try:
            payload = json.loads(raw_json)
        except json.JSONDecodeError as exc:
            raise SystemExit(f"OBS_PROMPTS_JSON is not valid JSON: {exc}") from exc
        if not isinstance(payload, list) or not all(isinstance(item, str) and item.strip() for item in payload):
            raise SystemExit("OBS_PROMPTS_JSON must be a non-empty JSON array of strings")
        return [item.strip() for item in payload]

    raw_split = os.getenv("OBS_PROMPTS")
    if raw_split:
        prompts = [item.strip() for item in raw_split.split("||") if item.strip()]
        if prompts:
            return prompts

    return list(DEFAULT_PROMPTS)


def extract_answer(payload: dict[str, Any]) -> str:
    message = payload.get("message")
    if isinstance(message, dict):
        content = message.get("content")
        if isinstance(content, str):
            return content
    response = payload.get("response")
    if isinstance(response, str):
        return response
    return json.dumps(payload, ensure_ascii=True)[:TRACE_PREVIEW_LIMIT]


def set_token_attributes(span: Span, payload: dict[str, Any]) -> None:
    prompt_tokens = payload.get("prompt_eval_count")
    generated_tokens = payload.get("eval_count")
    if isinstance(prompt_tokens, (int, float)) and prompt_tokens >= 0:
        span.set_attribute("gen_ai.usage.input_tokens", int(prompt_tokens))
    if isinstance(generated_tokens, (int, float)) and generated_tokens >= 0:
        span.set_attribute("gen_ai.usage.output_tokens", int(generated_tokens))


def seed_inference_runs(
    model: str,
    api_url: str,
    timeout_seconds: int,
    call_count: int,
    prompts: list[str],
) -> tuple[int, int]:
    tracer = trace.get_tracer("llm-observability-stack.otel-trace-seeder")
    success = 0
    failures = 0

    for index in range(call_count):
        prompt = prompts[index % len(prompts)]
        started = time.perf_counter()
        with tracer.start_as_current_span(
            "ollama seed chat",
            attributes={
                "gen_ai.system": "ollama",
                "gen_ai.operation.name": "chat",
                "gen_ai.request.model": model,
                "server.address": api_url,
                "llm.request.prompt_preview": truncate(prompt),
                "llm.seed.index": index + 1,
            },
        ) as span:
            try:
                response = requests.post(
                    api_url,
                    json={
                        "model": model,
                        "stream": False,
                        "messages": [{"role": "user", "content": prompt}],
                    },
                    timeout=(10, timeout_seconds),
                )
                span.set_attribute("http.response.status_code", response.status_code)
                response.raise_for_status()
                payload = response.json()
                answer = extract_answer(payload)
                elapsed_ms = (time.perf_counter() - started) * 1000
                set_token_attributes(span, payload)
                span.set_attribute("gen_ai.response.model", model)
                span.set_attribute("llm.response.latency_ms", round(elapsed_ms, 2))
                span.set_attribute("llm.response.preview", truncate(answer))
                success += 1
                print(f"[{index + 1}/{call_count}] OK latency_ms={elapsed_ms:.2f}")
            except Exception as exc:
                elapsed_ms = (time.perf_counter() - started) * 1000
                span.record_exception(exc)
                span.set_status(Status(StatusCode.ERROR, str(exc)))
                span.set_attribute("llm.response.latency_ms", round(elapsed_ms, 2))
                failures += 1
                print(f"[{index + 1}/{call_count}] ERROR error={exc}")

    print("\nCompleted OpenTelemetry GenAI trace seeding.")
    print(f"Success: {success}")
    print(f"Failed : {failures}")
    return success, failures


def main() -> None:
    endpoint = os.getenv("OTEL_EXPORTER_OTLP_ENDPOINT", DEFAULT_OTLP_ENDPOINT)
    model = os.getenv("OLLAMA_MODEL", DEFAULT_MODEL)
    api_url = os.getenv("OBS_INFERENCE_API_URL", DEFAULT_API_URL).rstrip("/")
    timeout_seconds = int(os.getenv("OBS_TIMEOUT_SECONDS", str(DEFAULT_TIMEOUT_SECONDS)))
    call_count = int(os.getenv("OBS_CALL_COUNT", str(DEFAULT_CALL_COUNT)))
    prompts = load_prompts()

    configure_tracing(endpoint)
    print(f"OTLP endpoint     : {endpoint}")
    print(f"Inference API URL : {api_url}")
    print(f"Model             : {model}")
    print(f"Calls             : {call_count}")
    print(f"Prompts loaded    : {len(prompts)}")

    seed_inference_runs(
        model=model,
        api_url=api_url,
        timeout_seconds=timeout_seconds,
        call_count=call_count,
        prompts=prompts,
    )


if __name__ == "__main__":
    main()
