import os
import time
from datetime import datetime, timezone

from langsmith import Client

from langsmith_inference_traces import (
    DEFAULT_CALL_COUNT,
    DEFAULT_ENDPOINT,
    DEFAULT_FEEDBACK_KEY,
    DEFAULT_MODEL,
    DEFAULT_PROJECT,
    DEFAULT_TIMEOUT_SECONDS,
    get_required_api_key,
    load_prompts,
    seed_inference_runs,
)

DEFAULT_INTERVAL_SECONDS = 300


def now_utc() -> str:
    return datetime.now(timezone.utc).isoformat()


def main() -> None:
    api_key = get_required_api_key()
    endpoint = os.getenv("LANGSMITH_ENDPOINT", DEFAULT_ENDPOINT)
    project = os.getenv("LANGSMITH_PROJECT", DEFAULT_PROJECT)
    model = os.getenv("OLLAMA_MODEL", DEFAULT_MODEL)
    api_url = os.getenv("OBS_INFERENCE_API_URL", "http://ollama:11434/api/chat").rstrip("/")
    timeout_seconds = int(os.getenv("OBS_TIMEOUT_SECONDS", str(DEFAULT_TIMEOUT_SECONDS)))
    feedback_key = os.getenv("LANGSMITH_FEEDBACK_KEY", DEFAULT_FEEDBACK_KEY)
    prompts = load_prompts()

    # Keeps existing behavior if only OBS_CALL_COUNT is set.
    calls_per_cycle = int(os.getenv("OBS_CALL_COUNT_PER_CYCLE", os.getenv("OBS_CALL_COUNT", str(DEFAULT_CALL_COUNT))))
    interval_seconds = int(os.getenv("OBS_INTERVAL_SECONDS", str(DEFAULT_INTERVAL_SECONDS)))
    max_cycles = int(os.getenv("OBS_MAX_CYCLES", "0"))

    client = Client(api_url=endpoint, api_key=api_key)

    print("Starting LangSmith dashboard seed scheduler.")
    print(f"Started at UTC      : {now_utc()}")
    print(f"LangSmith endpoint  : {endpoint}")
    print(f"LangSmith project   : {project}")
    print(f"Inference API URL   : {api_url}")
    print(f"Model               : {model}")
    print(f"Calls per cycle     : {calls_per_cycle}")
    print(f"Interval (seconds)  : {interval_seconds}")
    print(f"Max cycles (0=inf)  : {max_cycles}")
    print(f"Prompts loaded      : {len(prompts)}")

    cycle = 0
    while True:
        cycle += 1
        cycle_started = time.time()
        print(f"\n[{now_utc()}] Cycle {cycle} started")

        seed_inference_runs(
            client=client,
            project=project,
            model=model,
            api_url=api_url,
            timeout_seconds=timeout_seconds,
            call_count=calls_per_cycle,
            feedback_key=feedback_key,
            prompts=prompts,
        )

        if max_cycles > 0 and cycle >= max_cycles:
            print(f"[{now_utc()}] Reached OBS_MAX_CYCLES={max_cycles}. Exiting.")
            break

        elapsed = time.time() - cycle_started
        sleep_seconds = max(0, interval_seconds - int(elapsed))
        print(f"[{now_utc()}] Sleeping {sleep_seconds}s before next cycle.")
        time.sleep(sleep_seconds)


if __name__ == "__main__":
    main()
