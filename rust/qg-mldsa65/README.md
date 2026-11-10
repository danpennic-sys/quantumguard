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

Rust **1.85.0**, matching `rust-version` and `.github/workflows/qg-mldsa65-once.yml`.

```bash
cd rust/qg-mldsa65
cargo fmt --check
cargo test
cargo build --release
```

`cargo test` runs the unit tests and `tests/nist_acvp.rs`.

## NIST ACVP

Vendored files are the gzipped `internalProjection.json` projections from
`usnistgov/ACVP-Server` @ `975de31eb83d87039ec88934fdc47d8c312b892d`, as carried
by `fips204` 0.5.0. The suite fails if any case is skipped.

| File | Cases | Through this crate |
|------|------:|--------------------|
| keyGen | 75 | no (key generation is not this API) |
| sigGen | 360 | 60 ML-DSA-65 pure + HashML-DSA verifies |
| sigGen-tr1 | 720 | 120 ML-DSA-65 pure + HashML-DSA verifies |
| sigVer | 180 | 30 ML-DSA-65 pure (`verify_external`) + HashML-DSA (`verify_prehash`) |

The other ML-DSA-65 sigVer cases are internal and external-mu. They run through
the `fips204` 0.5.0 test hooks, not through `verify_mldsa65`. ML-DSA-44 and
ML-DSA-87 run the same way. This is not a CMVP certificate.

`verify_mldsa65` stays the QuantumGuard entry: 32-byte payload hash, empty context.

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

- Real FIPS 204 path via `fips204` 0.5.0.
- NIST ACVP evidence is in `tests/nist_acvp.rs` (recorded on `77a9205`, GitHub Actions run `42904351983`).
- Not wired into the Go CLI or HTTP verifier.
- The product `--pqc` path is still the deterministic stub. It does not claim quantum resistance.
