# QuantumGuard

Deterministic evidence verification. No network. No wall-clock. No external trust.

**v1 Status: FROZEN** (2026-09-09)  
See `CONTRACT-v1.md` for the normative specification.

## Contract v1 (frozen)

```json
{
  "version": "1",
  "payload": "<base64-encoded raw payload bytes>",
  "payload_hash": "<hex-encoded SHA-256 of the raw payload bytes>",
  "signature": "<hex-encoded 64-byte Ed25519 signature over the raw 32-byte payload_hash>",
  "public_key": "<hex-encoded 32-byte raw Ed25519 public key>",
  "alg": "ed25519"
}
```

### Verdict rules

| Condition                                         | Verdict        |
|---------------------------------------------------|----------------|
| Well-formed + hash matches + signature verifies   | `PASS`         |
| Well-formed but hash or signature wrong           | `FAIL`         |
| Malformed / incomplete / unsupported version/alg  | `INDETERMINATE`|

## Layout

```
quantumguard/
├── verifier/          # single deterministic verifier (shared)
├── cmd/qg/            # CLI: qg verify <bundle.json>
├── cmd/qg-http/       # HTTP server exposing POST /verify
├── http/              # HTTP handler package
├── vectors/           # conformance test vectors
├── evidence/          # (reserved)
├── README.md
└── PATCHES.md
```

## Quick start

```bash
# Run unit + equivalence tests
go test ./...

# CLI
go run ./cmd/qg verify vectors/01_known_good.json

# HTTP
go run ./cmd/qg-http &
curl -s -X POST --data-binary @vectors/01_known_good.json http://localhost:8080/verify
```

## Conformance vectors

| Vector                        | Expected        |
|-------------------------------|-----------------|
| `01_known_good.json`          | PASS            |
| `02_tampered_payload.json`    | FAIL            |
| `03_signature_failure.json`   | FAIL            |
| `04_hash_mismatch.json`       | FAIL            |
| `05_incomplete.json`          | INDETERMINATE   |
| `06_unsupported_alg.json`     | INDETERMINATE   |

CLI and HTTP must produce identical verdicts on every vector.

## Non-goals (today)

- Post-quantum signatures
- Key registries / DNS / revocation
- Wall-clock or freshness checks
- Second verifier implementation
- Payment integration

## License

Licensed under either of

- Apache License, Version 2.0 ([LICENSE-APACHE](LICENSE-APACHE))
- MIT license ([LICENSE-MIT](LICENSE-MIT))

at your option.
