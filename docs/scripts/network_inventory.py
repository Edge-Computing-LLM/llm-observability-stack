#!/usr/bin/env python3.11
"""Print a networking inventory for a Kubernetes namespace.

Requires:
  python3.11 -m pip install kubernetes

Examples:
  python3.11 docs/scripts/network_inventory.py --namespace llm-observability
  python3.11 docs/scripts/network_inventory.py --namespace llm-observability --context my-k3s
"""

from __future__ import annotations

import argparse
import datetime as dt
import sys
from typing import Iterable, List

from kubernetes import client, config
from kubernetes.client import ApiException


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Kubernetes namespace networking inventory")
    parser.add_argument("--namespace", default="llm-observability", help="Namespace to inspect")
    parser.add_argument("--context", default=None, help="Optional kubeconfig context")
    parser.add_argument("--kubeconfig", default=None, help="Optional kubeconfig file path")
    parser.add_argument("--in-cluster", action="store_true", help="Use in-cluster auth")
    return parser.parse_args()


def load_kube_config(args: argparse.Namespace) -> None:
    if args.in_cluster:
        config.load_incluster_config()
        return
    config.load_kube_config(config_file=args.kubeconfig, context=args.context)


def fmt_ports(ports: Iterable) -> str:
    parts: List[str] = []
    for item in ports or []:
        port = getattr(item, "port", None)
        target_port = getattr(item, "target_port", None)
        protocol = getattr(item, "protocol", "TCP")
        name = getattr(item, "name", "")
        if target_port is None:
            parts.append(f"{name}:{port}/{protocol}" if name else f"{port}/{protocol}")
        else:
            left = f"{name}:{port}" if name else str(port)
            parts.append(f"{left}->{target_port}/{protocol}")
    return ", ".join(parts) if parts else "-"


def fmt_selector(selector: dict | None) -> str:
    if not selector:
        return "-"
    return ",".join(f"{k}={v}" for k, v in sorted(selector.items()))


def print_table(title: str, headers: List[str], rows: List[List[str]]) -> None:
    print(f"\n=== {title} ===")
    if not rows:
        print("<none>")
        return

    widths = [len(h) for h in headers]
    for row in rows:
        for i, val in enumerate(row):
            widths[i] = max(widths[i], len(str(val)))

    def fmt_row(values: List[str]) -> str:
        return " | ".join(str(v).ljust(widths[i]) for i, v in enumerate(values))

    print(fmt_row(headers))
    print("-+-".join("-" * w for w in widths))
    for row in rows:
        print(fmt_row(row))


def collect(namespace: str) -> int:
    core = client.CoreV1Api()
    net = client.NetworkingV1Api()
    discovery = client.DiscoveryV1Api()

    try:
        pods = core.list_namespaced_pod(namespace=namespace).items
        services = core.list_namespaced_service(namespace=namespace).items
        endpoints = core.list_namespaced_endpoints(namespace=namespace).items
        endpoint_slices = discovery.list_namespaced_endpoint_slice(namespace=namespace).items
        network_policies = net.list_namespaced_network_policy(namespace=namespace).items
    except ApiException as exc:
        print(f"API error: {exc.status} {exc.reason}", file=sys.stderr)
        print(exc.body or "", file=sys.stderr)
        return 2

    pod_rows: List[List[str]] = []
    for pod in pods:
        pod_rows.append(
            [
                pod.metadata.name,
                pod.status.phase or "",
                pod.status.pod_ip or "",
                pod.spec.node_name or "",
                pod.status.host_ip or "",
            ]
        )

    svc_rows: List[List[str]] = []
    for svc in services:
        svc_rows.append(
            [
                svc.metadata.name,
                svc.spec.type or "",
                svc.spec.cluster_ip or "",
                fmt_ports(svc.spec.ports),
                fmt_selector(svc.spec.selector),
            ]
        )

    ep_rows: List[List[str]] = []
    for ep in endpoints:
        subset_addrs: List[str] = []
        subset_ports: List[str] = []
        for subset in ep.subsets or []:
            for addr in subset.addresses or []:
                subset_addrs.append(addr.ip)
            for port in subset.ports or []:
                subset_ports.append(f"{port.name or ''}:{port.port}/{port.protocol}")
        ep_rows.append(
            [
                ep.metadata.name,
                ",".join(subset_addrs) if subset_addrs else "-",
                ",".join(subset_ports) if subset_ports else "-",
            ]
        )

    eps_rows: List[List[str]] = []
    for eps in endpoint_slices:
        addresses: List[str] = []
        for endpoint in eps.endpoints or []:
            addresses.extend(endpoint.addresses or [])
        ports = [f"{p.name or ''}:{p.port}/{p.protocol}" for p in eps.ports or []]
        eps_rows.append(
            [
                eps.metadata.name,
                eps.metadata.labels.get("kubernetes.io/service-name", "-") if eps.metadata.labels else "-",
                ",".join(addresses) if addresses else "-",
                ",".join(ports) if ports else "-",
            ]
        )

    np_rows: List[List[str]] = []
    for np in network_policies:
        ptypes = ",".join(np.spec.policy_types or []) if np.spec else ""
        pod_sel = fmt_selector(np.spec.pod_selector.match_labels if np.spec and np.spec.pod_selector else {})
        np_rows.append([np.metadata.name, ptypes or "-", pod_sel])

    print(f"Timestamp: {dt.datetime.now(dt.timezone.utc).isoformat()}")
    print(f"Namespace: {namespace}")

    print_table("Pods", ["name", "phase", "pod_ip", "node", "host_ip"], pod_rows)
    print_table("Services", ["name", "type", "cluster_ip", "ports", "selector"], svc_rows)
    print_table("Endpoints", ["name", "addresses", "ports"], ep_rows)
    print_table("EndpointSlices", ["name", "service", "addresses", "ports"], eps_rows)
    print_table("NetworkPolicies", ["name", "policy_types", "pod_selector"], np_rows)

    return 0


def main() -> int:
    args = parse_args()
    try:
        load_kube_config(args)
    except Exception as exc:  # pragma: no cover
        print(f"Failed to load kube config: {exc}", file=sys.stderr)
        return 1
    return collect(args.namespace)


if __name__ == "__main__":
    raise SystemExit(main())
