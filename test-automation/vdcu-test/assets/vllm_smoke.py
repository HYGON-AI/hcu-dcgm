#!/usr/bin/env python3
"""vLLM OpenAI-compatible service smoke test for a single vDCU.

The script intentionally uses only the Python standard library so it can run
inside the vLLM image without adding a requests dependency.
"""
import argparse
import json
import subprocess
import sys
import time
import urllib.error
import urllib.request


def request(url, method="GET", payload=None, timeout=10):
    body = None
    headers = {}
    if payload is not None:
        body = json.dumps(payload).encode("utf-8")
        headers["Content-Type"] = "application/json"
    req = urllib.request.Request(url, data=body, headers=headers, method=method)
    with urllib.request.urlopen(req, timeout=timeout) as response:
        return json.loads(response.read().decode("utf-8"))


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--model", required=True)
    parser.add_argument("--host", default="127.0.0.1")
    parser.add_argument("--port", type=int, default=8000)
    parser.add_argument("--max-model-len", type=int, default=2048)
    parser.add_argument("--startup-timeout", type=int, default=300)
    args = parser.parse_args()

    command = [
        "vllm", "serve", args.model, "--host", args.host,
        "--port", str(args.port), "--max-model-len", str(args.max_model_len),
        "--gpu-memory-utilization", "0.85",
    ]
    try:
        # vLLM can emit a large startup log; discard it so the service cannot
        # block on a full pipe while this probe polls the HTTP endpoint.
        server = subprocess.Popen(command, stdout=subprocess.DEVNULL,
                                  stderr=subprocess.STDOUT, text=True)
    except OSError as exc:
        print("VLLM_FAIL " + json.dumps({"error": str(exc)}), flush=True)
        return 2

    base = "http://%s:%d" % (args.host, args.port)
    started = time.time()
    model_id = None
    try:
        while time.time() - started < args.startup_timeout:
            if server.poll() is not None:
                break
            try:
                models = request(base + "/v1/models")
                if models.get("data"):
                    model_id = models["data"][0].get("id") or args.model
                    break
            except (urllib.error.URLError, TimeoutError, ValueError):
                time.sleep(2)
        if not model_id:
            server.terminate()
            try:
                server.wait(timeout=20)
            except subprocess.TimeoutExpired:
                server.kill()
            print("VLLM_FAIL " + json.dumps({"error": "startup timeout",
                                               "log": ""}), flush=True)
            return 1

        result = request(base + "/v1/completions", "POST", {
            "model": model_id,
            "prompt": "用一句话介绍 GPU 虚拟化。",
            "max_tokens": 32,
            "temperature": 0,
        }, timeout=60)
        choices = result.get("choices") or []
        text = choices[0].get("text", "") if choices else ""
        print("VLLM_RESULT " + json.dumps({
            "model": model_id,
            "text": text[:200],
            "completion_tokens": result.get("usage", {}).get("completion_tokens", 0),
        }, ensure_ascii=False), flush=True)
        return 0 if text else 1
    finally:
        server.terminate()
        try:
            server.wait(timeout=20)
        except subprocess.TimeoutExpired:
            server.kill()


if __name__ == "__main__":
    sys.exit(main())
