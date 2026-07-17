#!/usr/bin/env python3.11
"""Watch endpoint changes for a single service.

Requires:
  python3.11 -m pip install kubernetes

Example:
  python3.11 docs/scripts/watch_endpoints.py --namespace llm-observability --service ollama --timeout 600
"""

from __future__ import annotations

import argparse
import datetime as dt
import sys
from typing import List

from kubernetes import client, config, watch
from kubernetes.client import ApiException


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Watch endpoint changes for one service")
    parser.add_argument("--namespace", default="llm-observability", help="Namespace")
    parser.add_argument("--service", required=True, help="Service name")
    parser.add_argument("--context", default=None, help="Optional kubeconfig context")
    parser.add_argument("--kubeconfig", default=None, help="Optional kubeconfig path")
    parser.add_argument("--in-cluster", action="store_true", help="Use in-cluster auth")
    parser.add_argument("--timeout", type=int, default=300, help="Watch timeout seconds")
    return parser.parse_args()


def load_cfg(args: argparse.Namespace) -> None:
    if args.in_cluster:
        config.load_incluster_config()
    else:
        config.load_kube_config(config_file=args.kubeconfig, context=args.context)


def endpoint_snapshot(ep: client.V1Endpoints) -> str:
    blocks: List[str] = []
    for subset in ep.subsets or []:
        addrs = [a.ip for a in (subset.addresses or [])]
        not_ready = [a.ip for a in (subset.not_ready_addresses or [])]
        ports = [f"{p.name or ''}:{p.port}/{p.protocol}" for p in (subset.ports or [])]
        blocks.append(f"ready={addrs or []} not_ready={not_ready or []} ports={ports or []}")
    return " | ".join(blocks) if blocks else "<no endpoints>"


def main() -> int:
    args = parse_args()
    try:
        load_cfg(args)
    except Exception as exc:
        print(f"Failed to load kube config: {exc}", file=sys.stderr)
        return 1

    core = client.CoreV1Api()

    try:
        core.read_namespaced_service(name=args.service, namespace=args.namespace)
    except ApiException as exc:
        if exc.status == 404:
            print(f"Service not found: {args.namespace}/{args.service}", file=sys.stderr)
            return 2
        print(f"API error while checking service: {exc.status} {exc.reason}", file=sys.stderr)
        return 2

    print(f"Watching endpoints for service {args.namespace}/{args.service} (timeout={args.timeout}s)")
    print("Press Ctrl+C to stop.")

    w = watch.Watch()
    try:
        for event in w.stream(
            core.list_namespaced_endpoints,
            namespace=args.namespace,
            timeout_seconds=args.timeout,
            field_selector=f"metadata.name={args.service}",
        ):
            etype = event.get("type", "UNKNOWN")
            obj: client.V1Endpoints = event["object"]
            timestamp = dt.datetime.now(dt.timezone.utc).isoformat()
            print(f"[{timestamp}] {etype} {obj.metadata.name} -> {endpoint_snapshot(obj)}")
    except KeyboardInterrupt:
        print("Stopped by user.")
    except ApiException as exc:
        print(f"Watch API error: {exc.status} {exc.reason}", file=sys.stderr)
        return 2

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
