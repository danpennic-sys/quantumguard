# QuantumGuard V2 PQC Contract

**Status:** Normative amendment to CONTRACT-PQC-1 (stub era). Not production crypto.  
**Date:** 2026-10-23  
**Prerequisite:** CONTRACT-v1.md frozen; CONTRACT-PQC-1.md remains the behavioral base.  
**License:** dual MIT OR Apache-2.0  

This document does **not** replace CONTRACT-PQC-1. It adds:

1. Stable `code` values for the PQC envelope  
2. An honest-stub clause equivalent to VM-EXPORT-BUNDLE-001 §3.1  
3. A reserved failure taxonomy for when real ML-DSA-65 lands  

It does **not** invent new algorithm families, new evidence formats, or new message constructions.

---

## 1. Intent (locked from CONTRACT-PQC-1)

| Item | Decision |
|------|----------|
| Algorithm family (stub era and first real path) | **ML-DSA-65 only** (`alg`: `"ml-dsa-65"`) |
| Other families (SLH-DSA, etc.) | Out of scope for PQC-1 / this contract |
| Bundle format | Unchanged v1 evidence JSON |
| Message bytes | `CanonicalBytes(bundle)` = raw 32-byte `payload_hash` (same as v1 Ed25519 input) |
| Opt-in | CLI `--pqc`; HTTP `?pqc=true` only |
| Relation to v1 | Additive. Never mutates v1 `verdict` / `code` / exit status |

---

## 2. Envelope (normative)

Same shape as CONTRACT-PQC-1, with `code` on both the v1 top level and the optional `pqc` object.

**Disabled (default):**

```json
{
  "verdict": "PASS",
  "code": "OK",
  "reason": "valid",
  "pqc": null
}
```

**Enabled, stub era (current implementation):**

```json
{
  "verdict": "PASS",
  "code": "OK",
  "reason": "valid",
  "pqc": {
    "verdict": "INDETERMINATE",
    "code": "E_PQC_STUB",
    "reason": "pqc-1 stub: ML-DSA-65 not implemented",
    "alg": "ml-dsa-65",
    "stub": true
  }
}
```

Rules:

- Top-level `verdict` / `code` / `reason` are **exactly** the v1 result.
- `pqc` is JSON `null` when disabled.
- When enabled, `pqc` is an object. It **never** overwrites top-level fields.
- HTTP status and CLI exit code follow **v1 only**.

---

## 3. Honest-stub clause (normative)

When PQC is enabled and the implementation has not performed a real ML-DSA-65 verification, the PQC object **MUST** report:

- `verdict`: `INDETERMINATE` (or `FAIL` if a required PQC check was demanded and cannot be satisfied — see §5)
- `code`: `E_PQC_STUB`
- `stub`: `true`

It **MUST NOT** set `pqc.verdict` to `PASS` for a cryptographic check it did not perform.

This is the same rule as VM-EXPORT-BUNDLE-001 §3.1, applied to the PQC layer.

---

## 4. Stub-era codes (normative now)

| Condition | `pqc.verdict` | `pqc.code` | `stub` |
|-----------|---------------|------------|--------|
| PQC disabled | (field absent / null) | — | — |
| PQC enabled, CanonicalBytes OK, real crypto not implemented | `INDETERMINATE` | `E_PQC_STUB` | `true` |
| PQC enabled, CanonicalBytes fails | `INDETERMINATE` | `E_PQC_CANONICAL` | `true` |
| PQC enabled, alg present and not `ml-dsa-65` | `INDETERMINATE` | `E_PQC_UNSUPPORTED_ALG` | `true` |

`alg` is always `"ml-dsa-65"` on the stub success path. Prose `reason` remains non-normative.

---

## 5. Reserved codes (real ML-DSA-65 — not implemented)

These codes are reserved. Emitting them before real verification is non-conforming.

| Condition | `pqc.verdict` | `pqc.code` | `stub` |
|-----------|---------------|------------|--------|
| Signature verifies | `PASS` | `OK` | `false` |
| Signature fails | `FAIL` | `E_PQC_SIG_FAIL` | `false` |
| Malformed PQC fields | `INDETERMINATE` | `E_PQC_MALFORMED` | `false` |

No implementation may claim `PASS` under §5 until a real ML-DSA-65 verify path exists and is covered by conformance vectors.

---

## 6. What this contract does not do

- Does not implement ML-DSA-65.
- Does not claim quantum resistance.
- Does not change CONTRACT-v1.md or v1 evidence bundles.
- Does not let PQC alter v1 PASS / FAIL / INDETERMINATE.
- Does not add PQC-only evidence formats, chains, or Merkle variants.
- Does not authorize SLH-DSA or other families under `alg: "ml-dsa-65"`.

---

## 7. Conformance order

1. This contract + CONTRACT-PQC-1 (design locked)  
2. Golden envelopes under `vectors/v2/` (must match §2 / §4)  
3. Real-crypto vectors — **only after** a non-stub implementation exists  

Do not expand `vectors/v2/` with “PASS under stub” fixtures. That would encode a gap.

---

## 8. Status line

**PQC-1 / V2 contract: stub era. Codes and honest-stub clause normative. Real ML-DSA-65 coverage: 0. Fixtures: envelope goldens only. Not production crypto.**
