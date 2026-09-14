# QuantumGuard Verifier v2 (PQC Track)

**Status:** Deterministic stub — ML-DSA-65 crypto not implemented  
**Algorithm target:** ML-DSA-65  
**Governing documents:** `../../CONTRACT-PQC-1.md`, `../../UPGRADE-v1-to-v2.md`

## Rules

- v1 code remains completely untouched (`verifier/verifier.go`).
- PQC-1 is opt-in (`--pqc` / `?pqc=true`). Default `pqc` is JSON `null`.
- The PQC result never changes the v1 verdict, HTTP status, or CLI exit code.
- Message bytes come only from `verifier.CanonicalBytes` (raw `payload_hash`).
- No hybrid signatures. No algorithm negotiation.
- Same three-verdict model. Same determinism / offline invariants.

## Current contents

- `pqc.go` — stub verifier + combined CLI/HTTP envelope
- Replay goldens under `../../vectors/v2/`

## Not yet

- Real ML-DSA-65 (e.g. filippo.io/mldsa)
- Production PQC evidence fields
