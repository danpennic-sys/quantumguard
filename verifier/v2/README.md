# QuantumGuard Verifier v2 (PQC Track)

**Status:** Scaffold only — not implemented  
**Algorithm target:** ML-DSA-65  
**Governing document:** `../../UPGRADE-v1-to-v2.md`

## Rules

- v1 code remains completely untouched.
- This package will eventually handle ML-DSA-65 bundles.
- Until the implementation is complete and proven, any v2 bundle MUST be treated as `INDETERMINATE` by the main verifier.
- No hybrid signatures.
- No algorithm negotiation.
- Same three-verdict model.
- Same single-verifier architecture (this will be integrated, not run in parallel).

## Planned contents (future)

- ML-DSA-65 verification path
- Type-prefixed or versioned bundle parsing
- v2-specific conformance vectors under `../../vectors/v2/`
- Integration into the single top-level `VerifyJSON` entry point

## Current state

Empty. No crypto. No dependencies added yet.
