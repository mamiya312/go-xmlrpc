#!/usr/bin/env python3
import sys
import xmlrpc.client

server = xmlrpc.client.ServerProxy("http://127.0.0.1:3000/RPC2")

try:
    echo = server.Test.echo("hello from python")
    if echo != "hello from python":
        raise AssertionError(f"echo mismatch: {echo!r}")

    payload = server.Test.sample(
        "edge-case",
        -2147483648,
        True,
        3.141592653589793,
        ["alpha", "beta"],
        {"name": "python"},
    )
    reply = payload.get("Reply", payload)
    if reply["Arg1"] != "edge-case":
        raise AssertionError(f"Arg1 mismatch: {reply!r}")
    if reply["Arg2"] != -2147483648:
        raise AssertionError(f"Arg2 mismatch: {reply!r}")
    if reply["Arg3"] is not True:
        raise AssertionError(f"Arg3 mismatch: {reply!r}")
    if abs(reply["Arg4"] - 3.141592653589793) >= 1e-12:
        raise AssertionError(f"Arg4 mismatch: {reply!r}")
    if reply["Arg5"] != ["alpha", "beta"]:
        raise AssertionError(f"Arg5 mismatch: {reply!r}")
    if reply["Arg6"]["name"] != "sample":
        raise AssertionError(f"Arg6 mismatch: {reply!r}")

    print({
        "python_echo": echo,
        "python_sample": reply,
    })
except Exception as exc:  # pragma: no cover
    print(f"Python XML-RPC check failed: {exc}", file=sys.stderr)
    raise
