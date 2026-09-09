# QuantumGuard Upgrade Contract: v1 → v2

**Status:** Normative  
**Date:** 2026-09-09  
**Prerequisite:** CONTRACT-v1.md is frozen and read-only

## Purpose

Define exactly what may change and what must never change when moving from classical Ed25519 (v1) to post-quantum ML-DSA-65 (v2).

## Permanently Frozen (must never change)

These properties are architectural invariants. Violating any of them is a breaking change that requires a new major version and a new project decision.

1. **Single verifier implementation**  
   There is exactly one verification function that both CLI and HTTP call. No parallel or alternate verifiers.

2. **Determinism**  
   Verification is a pure function of the evidence bytes. Same input bytes → same verdict, always.

3. **Offline / no external trust**  
   No network, DNS, registry, revocation service, or wall-clock dependency may influence the verdict.

4. **Three-verdict model only**  
   The only legal outputs are: `PASS` | `FAIL` | `INDETERMINATE`.

5. **CLI ↔ HTTP equivalence**  
   Both interfaces must produce bit-for-bit identical results on the same evidence.

6. **Replay harness as arbiter**  
   Conformance vectors + `go test ./...` remain the sole authority on correctness.

## Allowed to Change in v2

| Area                    | Rule                                                                 |
|-------------------------|----------------------------------------------------------------------|
| Signature algorithm     | May add ML-DSA-65 (and only ML-DSA-65 for v2)                        |
| Public key encoding     | New encoding rules for ML-DSA-65 public keys                         |
| Signature encoding      | New encoding rules for ML-DSA-65 signatures                          |
| Bundle format           | Must be distinguishable from v1 (Type-Prefixed or version field)     |
| Key / signature sizes   | Expected to be larger; size limits may be raised                     |
| Supported `alg` values  | v1: `"ed25519"` only. v2: `"ed25519"` + `"ml-dsa-65"`                |

## Version Discrimination Rules

- v1 bundles continue to use `"version": "1"` and `"alg": "ed25519"`.
- v2 bundles MUST be unambiguously distinguishable from v1.
- Recommended approach: Type-Prefixed envelope or a new top-level `"version": "2"`.
- The single verifier MUST reject any ambiguous or mixed-version bundle as `INDETERMINATE`.

## Compatibility Requirements

- A v2 verifier MUST still correctly verify pure v1 (Ed25519) bundles.
- A v1-only verifier MUST treat any v2 bundle as `INDETERMINATE` (unsupported), never as `FAIL`.
- No silent fallback or algorithm negotiation is permitted.

## Explicit Non-Goals of the v1→v2 Transition

- Hybrid signatures (classical + PQC in one signature)
- Multiple PQC algorithms in v2
- Changing the three-verdict model
- Introducing any network or time dependency
- Breaking CLI ↔ HTTP equivalence

## Acceptance Criteria for v2

v2 is only considered ready when:

1. All existing v1 vectors still produce the exact same verdicts.
2. New v2 (ML-DSA-65) vectors exist and pass.
3. CLI and HTTP produce identical results on every vector (v1 and v2).
4. `go test ./...` is green.
5. The upgrade rules in this document are satisfied.
