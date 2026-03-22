#!/usr/bin/env python3
"""Trace a Kubernetes Service to selectors, Pods, Endpoints, and EndpointSlices.

Requires:
  pip install kubernetes

Example:
  python docs/scripts/service_path_inspector.py --namespace llm-observability --service ollama
"""

from __future__ import annotations

import argparse
import sys
from typing import Dict, List

from kubernetes import client, config
from kubernetes.client import ApiException


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Inspect end-to-end service path")
    parser.add_argument("--namespace", default="llm-observability", help="Namespace")
    parser.add_argument("--service", required=True, help="Service name")
    parser.add_argument("--context", default=None, help="Optional kubeconfig context")
    parser.add_argument("--kubeconfig", default=None, help="Optional kubeconfig file path")
    parser.add_argument("--in-cluster", action="store_true", help="Use in-cluster auth")
    return parser.parse_args()


def load_cfg(args: argparse.Namespace) -> None:
    if args.in_cluster:
        config.load_incluster_config()
    else:
        config.load_kube_config(config_file=args.kubeconfig, context=args.context)


def selector_to_query(selector: Dict[str, str]) -> str:
    if not selector:
        return ""
    return ",".join(f"{k}={v}" for k, v in selector.items())


def print_header(label: str) -> None:
    print(f"\n=== {label} ===")


def main() -> int:
    args = parse_args()
    try:
        load_cfg(args)
    except Exception as exc:
        print(f"Failed to load kube config: {exc}", file=sys.stderr)
        return 1

    core = client.CoreV1Api()
    discovery = client.DiscoveryV1Api()

    try:
        svc = core.read_namespaced_service(name=args.service, namespace=args.namespace)
    except ApiException as exc:
        if exc.status == 404:
            print(f"Service not found: {args.namespace}/{args.service}", file=sys.stderr)
            return 2
        print(f"API error: {exc.status} {exc.reason}", file=sys.stderr)
        return 2

    selector = svc.spec.selector or {}
    selector_query = selector_to_query(selector)

    print_header("Service")
    print(f"name: {svc.metadata.name}")
    print(f"namespace: {svc.metadata.namespace}")
    print(f"type: {svc.spec.type}")
    print(f"cluster_ip: {svc.spec.cluster_ip}")
    print(f"ports: {[{'name': p.name, 'port': p.port, 'targetPort': p.target_port, 'protocol': p.protocol} for p in (svc.spec.ports or [])]}")
    print(f"selector: {selector if selector else '<none>'}")

    print_header("Selected Pods")
    if not selector_query:
        print("Service has no selector (possibly externalName or manually managed endpoints).")
        pods = []
    else:
        pods = core.list_namespaced_pod(namespace=args.namespace, label_selector=selector_query).items
        if not pods:
            print("No Pods matched service selector.")
        for pod in pods:
            print(
                f"- {pod.metadata.name} phase={pod.status.phase} ready={all((cs.ready for cs in (pod.status.container_statuses or [])))} "
                f"pod_ip={pod.status.pod_ip} node={pod.spec.node_name}"
            )

    print_header("Endpoints")
    try:
        ep = core.read_namespaced_endpoints(name=args.service, namespace=args.namespace)
        if not ep.subsets:
            print("No endpoint subsets present.")
        else:
            for subset in ep.subsets:
                addrs = [a.ip for a in (subset.addresses or [])]
                not_ready = [a.ip for a in (subset.not_ready_addresses or [])]
                ports = [f"{p.name}:{p.port}/{p.protocol}" for p in (subset.ports or [])]
                print(f"addresses={addrs or []} not_ready={not_ready or []} ports={ports or []}")
    except ApiException as exc:
        print(f"Failed to read Endpoints: {exc.status} {exc.reason}")

    print_header("EndpointSlices")
    try:
        slices = discovery.list_namespaced_endpoint_slice(
            namespace=args.namespace,
            label_selector=f"kubernetes.io/service-name={args.service}",
        ).items
        if not slices:
            print("No EndpointSlices found for service.")
        for eps in slices:
            addresses: List[str] = []
            for endpoint in eps.endpoints or []:
                addresses.extend(endpoint.addresses or [])
            ports = [f"{p.name}:{p.port}/{p.protocol}" for p in eps.ports or []]
            print(f"- {eps.metadata.name} addresses={addresses or []} ports={ports or []}")
    except ApiException as exc:
        print(f"Failed to list EndpointSlices: {exc.status} {exc.reason}")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
