# QuantumGuard Contract PQC-1 (Stub)

**Status:** Experimental stub — not production crypto  
**Date:** 2026-09-14  
**Prerequisite:** CONTRACT-v1.md remains frozen and authoritative

PQC-1 is an additive, opt-in layer beside v1. It does not replace v1,
does not change v1 defaults, and does not change evidence semantics.

This document is honest: ML-DSA-65 verification is **not implemented**.
The shipped path is a deterministic stub so CLI/HTTP wiring can freeze
without claiming post-quantum security.

---

## 1. Opt-in only

PQC-1 is off by default.

| Interface | Enable with | Default |
|-----------|-------------|---------|
| CLI       | `qg verify --pqc <bundle.json>` | `--pqc` absent |
| HTTP      | `POST /verify?pqc=true`         | query absent, or any value other than exactly `true` |

No other flag, header, or bundle field enables PQC-1.

---

## 2. Envelope (normative)

CLI and HTTP encode the same envelope via `verifier/v2.Encode`.

**Disabled (default):**

```json
{
  "verdict": "PASS",
  "reason": "valid",
  "pqc": null
}
```

**Enabled (`--pqc` / `?pqc=true`):**

```json
{
  "verdict": "PASS",
  "reason": "valid",
  "pqc": {
    "verdict": "INDETERMINATE",
    "reason": "pqc-1 stub: ML-DSA-65 not implemented",
    "alg": "ml-dsa-65",
    "stub": true
  }
}
```

- `verdict` and `reason` are **exactly** the v1 result from `verifier.VerifyJSON`.
- `pqc` is JSON `null` when disabled.
- `pqc` is an object when enabled. It never overwrites `verdict` / `reason`.
- HTTP status and CLI exit code follow the **v1** verdict only.

The frozen type `verifier.Result` is unchanged and does not contain `pqc`.

---

## 3. Message bytes

PQC-1 verifies over `CanonicalBytes(bundle)` — the raw 32-byte
`payload_hash` that v1 uses as its Ed25519 signature input.

No alternate message construction.

If `CanonicalBytes` fails, the stub returns `INDETERMINATE` with reason
prefixed by `canonical_bytes: `.

---

## 4. Stub verdict rules

Until real ML-DSA-65 lands:

| Condition | PQC verdict | PQC reason |
|-----------|-------------|------------|
| `CanonicalBytes` succeeds | `INDETERMINATE` | `pqc-1 stub: ML-DSA-65 not implemented` |
| `CanonicalBytes` fails    | `INDETERMINATE` | `canonical_bytes: <error>` |

`alg` is always `"ml-dsa-65"`. `stub` is always `true`.

The stub is a pure function of the evidence bytes. No network, no clock,
no external trust.

---

## 5. What this contract does not do

- Does not implement ML-DSA-65.
- Does not claim quantum resistance.
- Does not change v1 evidence bundles.
- Does not change `CONTRACT-v1.md`.
- Does not let a PQC result alter a v1 PASS / FAIL / INDETERMINATE.
- Does not add PQC-only evidence formats, chains, or Merkle variants.

---

## 6. Replay

Conformance vectors under `vectors/` remain the v1 golden set.

PQC envelope goldens:

| File | Meaning |
|------|---------|
| `vectors/v2/pqc_disabled_golden.json` | known-good, PQC off (`pqc: null`) |
| `vectors/v2/pqc_stub_golden.json`     | known-good, PQC on (stub) |

`go test ./...` is the arbiter. Any drift from these goldens is a fail.

---

## 7. Placement

All PQC implementation lives under `verifier/v2/`.
`verifier/verifier.go` is not modified by this contract.
