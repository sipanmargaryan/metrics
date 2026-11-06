#!/usr/bin/env python3
"""Send 100 API requests to generate Prometheus/Grafana traffic (stdlib only)."""

import argparse
from collections import Counter
import json
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
import uuid


class NoRedirects(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        return None


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--base-url", default="http://localhost:8080")
    parser.add_argument("--delay", type=float, default=0.35, help="seconds between requests (default: 0.35)")
    args = parser.parse_args()
    url = urllib.parse.urlsplit(args.base_url)
    if url.scheme not in ("http", "https") or not url.netloc or url.query or url.fragment:
        parser.error("--base-url must be an HTTP(S) URL without a query or fragment")
    if not 0 <= args.delay < float("inf"):
        parser.error("--delay must be a finite, non-negative number")

    opener = urllib.request.build_opener(NoRedirects())
    statuses = Counter()
    endpoints = Counter()
    sent = 0
    unexpected = 0

    def send(method, path, expected, payload=None):
        nonlocal sent, unexpected
        if sent:
            time.sleep(args.delay)
        data = None if payload is None else json.dumps(payload).encode()
        request = urllib.request.Request(
            args.base_url.rstrip("/") + path,
            data=data,
            method=method,
            headers={"Content-Type": "application/json"},
        )
        sent += 1
        start = time.perf_counter()
        try:
            try:
                response = opener.open(request, timeout=10)
            except urllib.error.HTTPError as error:
                response = error  # Expected 400/404/409 responses still count.
            with response:
                status = response.code
                body = response.read().decode()
        except (urllib.error.URLError, TimeoutError, OSError) as error:
            raise RuntimeError("Cannot reach API: {}. Start it with make up.".format(error)) from error

        statuses[status] += 1
        endpoints["{} {}".format(method, path)] += 1
        elapsed = (time.perf_counter() - start) * 1000
        print("[{}/100] {} {} -> {} ({:.1f} ms)".format(sent, method, path, status, elapsed), flush=True)
        if status != expected:
            unexpected += 1
            print("  Expected {}; response: {}".format(expected, body[:300]), file=sys.stderr)
        try:
            return status, json.loads(body)
        except json.JSONDecodeError:
            return status, None

    # These two setup requests are included in the total of 100.
    name = "metrics-load-" + uuid.uuid4().hex[:12]
    status, warehouse = send("POST", "/api/v1/warehouses", 201, {"name": name})
    if status != 201 or not isinstance(warehouse, dict) or not warehouse.get("id"):
        raise RuntimeError("Warehouse setup failed. Check migrations and API logs.")
    product_payload = {"name": name, "warehouse_id": warehouse["id"], "quantity": 1000, "price": 1999}
    status, product = send("POST", "/api/v1/products", 201, product_payload)
    if status != 201 or not isinstance(product, dict) or not product.get("id"):
        raise RuntimeError("Product setup failed. Check migrations and API logs.")

    cases = [
        ("GET", "/api/v1/warehouses", 200, None),
        ("GET", "/api/v1/products", 200, None),
        ("POST", "/api/v1/orders", 201, {"product_id": product["id"], "quantity": 1}),
        ("POST", "/api/v1/products", 400, {}),
        ("POST", "/api/v1/warehouses", 409, {"name": name}),
        ("GET", "/api/v1/load-test-not-found", 404, None),
        ("POST", "/api/v1/products", 409, product_payload),
    ]
    for i in range(98):
        send(*cases[i % len(cases)])

    print("\nCompleted {} requests; {} unexpected statuses.".format(sent, unexpected))
    print("HTTP statuses: " + ", ".join("{}={}".format(code, count) for code, count in sorted(statuses.items())))
    for endpoint, count in sorted(endpoints.items()):
        print("  {}: {}".format(endpoint, count))
    print("Test records remain in the database. Allow 15 seconds for the next Prometheus scrape.")
    return 1 if unexpected else 0


if __name__ == "__main__":
    try:
        sys.exit(main())
    except (RuntimeError, KeyboardInterrupt) as error:
        print("\nStopped: {}".format(error or "interrupted"), file=sys.stderr)
        sys.exit(1)
