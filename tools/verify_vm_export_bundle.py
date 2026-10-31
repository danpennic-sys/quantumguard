#!/usr/bin/env python3
# SPDX-License-Identifier: MIT OR Apache-2.0
# Copyright (c) 2026 danpennic-sys
"""
VM-EXPORT-BUNDLE-001 — fail-closed verifier skeleton (verify-only).

Status: design-aid only. Not a production implementation.
Emits machine-readable JSON. Never claims PASS for unimplemented checks.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import sys
from pathlib import Path
from typing import Any


def jcs_like_canonical(obj: Any) -> bytes:
    """Minimal deterministic JSON for the skeleton (sorted keys, no whitespace).
    Production implementations MUST use full RFC 8785 / JCS.
    """
    return json.dumps(obj, sort_keys=True, separators=(",", ":"), ensure_ascii=False).encode("utf-8")


def compute_state_hash(state: Any) -> str:
    canonical = jcs_like_canonical(state)
    return hashlib.sha256(canonical).hexdigest()


def compute_replay_merkle_root(events: list[dict]) -> str:
    """Domain-separated Merkle (0x01 prefix), Bitcoin-style odd carry."""
    if not events:
        return ""
    leaves = []
    for e in events:
        h = bytes.fromhex(e["hash"])
        leaves.append(h)
    while len(leaves) > 1:
        next_level = []
        for i in range(0, len(leaves), 2):
            if i + 1 == len(leaves):
                next_level.append(leaves[i])
            else:
                h = hashlib.sha256(b"\x01" + leaves[i] + leaves[i + 1]).digest()
                next_level.append(h)
        leaves = next_level
    return leaves[0].hex()


def verify_case(input_data: dict) -> dict:
    """Return {verdict, code, detail?, computed_*?} — fail-closed."""
    has_state = bool(input_data.get("state_hash"))
    has_merkle = bool(input_data.get("replay_merkle_root"))
    if has_state and has_merkle:
        return {"verdict": "FAIL", "code": "BOTH_HASHES_CLAIMED"}

    files = input_data.get("files") or []
    seen = set()
    for f in files:
        name = f.get("name", "")
        if "/" in name or "\\" in name or name.startswith("..") or name in (".", "..", ""):
            return {
                "verdict": "FAIL",
                "code": "ZIP_SLIP",
                "detail": f'invalid entry name "{name}" (subdirectories not allowed)',
            }
        if name in seen:
            return {
                "verdict": "FAIL",
                "code": "DUPLICATE_PATH",
                "detail": f"path appears more than once: {name}",
            }
        seen.add(name)

    extracted = input_data.get("extracted_files")
    if extracted is not None:
        allowed = {"bundle.json", "signature.ed25519", "public_key.pem"} | {
            f["name"] for f in files
        }
        for name in extracted:
            if name not in allowed:
                return {
                    "verdict": "FAIL",
                    "code": "EXTRA_FILE",
                    "detail": f"unauthorized file present in bundle: {name}",
                }

    if input_data.get("require_adr_anchor"):
        return {
            "verdict": "FAIL",
            "code": "NOT_IMPLEMENTED",
            "detail": "adr_anchor: NOT IMPLEMENTED in v1.0 — see VM-EXPORT-BUNDLE-001 §3.1",
        }

    if has_state:
        state = input_data.get("state")
        if state is None:
            return {"verdict": "FAIL", "code": "STATE_HASH_MISMATCH", "detail": "missing state"}
        computed = compute_state_hash(state)
        claimed = input_data["state_hash"]
        if computed != claimed:
            outer_ca = input_data.get("created_at")
            state_ca = state.get("created_at")
            if outer_ca and state_ca and outer_ca != state_ca:
                return {
                    "verdict": "FAIL",
                    "code": "CREATED_AT_TAMPER",
                    "detail": "created_at is inside the signed region; any modification invalidates root_hash/signature",
                }
            return {"verdict": "FAIL", "code": "STATE_HASH_MISMATCH"}
        return {
            "verdict": "PASS",
            "code": "OK",
            "computed_state_hash": computed,
        }

    if has_merkle:
        events = input_data.get("events") or []
        computed = compute_replay_merkle_root(events)
        claimed = input_data["replay_merkle_root"]
        if computed != claimed:
            return {"verdict": "FAIL", "code": "REPLAY_MERKLE_MISMATCH"}
        return {
            "verdict": "PASS",
            "code": "OK",
            "computed_replay_merkle_root": computed,
        }

    return {
        "verdict": "FAIL",
        "code": "NOT_IMPLEMENTED",
        "detail": "no hash claimed and no --no-* suppression",
    }


def main() -> int:
    parser = argparse.ArgumentParser(description="VM-EXPORT-BUNDLE-001 verify-only skeleton")
    parser.add_argument("input", type=Path, help="Path to input.json fixture")
    parser.add_argument("--expected", type=Path, help="Optional expected.json for self-check")
    args = parser.parse_args()

    data = json.loads(args.input.read_text())
    result = verify_case(data)

    print(json.dumps(result, indent=2, sort_keys=True))

    if args.expected:
        expected = json.loads(args.expected.read_text())
        if result.get("verdict") != expected.get("verdict") or result.get("code") != expected.get("code"):
            print("SELF-CHECK FAILED", file=sys.stderr)
            return 2

    return 0 if result.get("verdict") == "PASS" else 1


if __name__ == "__main__":
    sys.exit(main())
