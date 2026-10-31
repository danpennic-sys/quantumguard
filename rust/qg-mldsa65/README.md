# qg-mldsa65

Real **ML-DSA-65** (FIPS 204) verification for QuantumGuard.

Contract: [`docs/QUANTUMGUARD-V2-PQC-CONTRACT.md`](../../docs/QUANTUMGUARD-V2-PQC-CONTRACT.md)

| Item | Value |
|------|--------|
| Message | 32-byte `payload_hash` (same as v1 Ed25519 input) |
| Algorithm | `ml-dsa-65` only |
| Context | empty (Pure ML-DSA) |
| `stub` | always `false` — real verify or form reject |

## Build

Requires a recent Rust stable (dependency graph uses crates that need modern Cargo).

```bash
cd rust/qg-mldsa65
cargo test
cargo build --release
```

## Library

```rust
use qg_mldsa65::verify_mldsa65;

let r = verify_mldsa65(&pk_bytes, &payload_hash_32, &sig_bytes, "ml-dsa-65");
// r.verdict: PASS | FAIL | INDETERMINATE
// r.code: OK | E_PQC_SIG_FAIL | E_PQC_MALFORMED | ...
// r.stub: false
```

## CLI

```bash
qg-mldsa65-verify --pk <hex> --sig <hex> --hash <hex32>
qg-mldsa65-verify --pk-file pk.bin --sig-file sig.bin --bundle bundle.json
```

## Status

- Real FIPS 204 path via `fips204` crate.
- Not wired into the Go CLI envelope yet.
- End-to-end PQC product coverage still 0 until production ML-DSA-65 evidence vectors exist.
- Does **not** claim quantum resistance for the Go/default path.
