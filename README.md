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
├── docs/              # domain bridges
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

## Experimental: PQC-1 (opt-in stub)

PQC-1 does **not** change v1. It is off by default (`pqc` is JSON `null`).

```bash
go run ./cmd/qg verify --pqc vectors/01_known_good.json
curl -s -X POST --data-binary @vectors/01_known_good.json 'http://localhost:8080/verify?pqc=true'
```

This path is a deterministic stub. It does not implement ML-DSA-65 and
does not claim quantum resistance. See `CONTRACT-PQC-1.md`.

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

## Domain profiles

QuantumGuard is integrity-only. Domain policy lives in profiles:

| Profile | Repo | Role |
|---------|------|------|
| Robotics Trust Profile (RTP) v1.0 | [robotics-trust-profile](https://github.com/danpennic-sys/robotics-trust-profile) | Kinematic / spatial / cross-sensor admissibility for robot GIRs |

Deep bridge (canonical GIR hashing, dual-verdict audit, actuation interlock):  
https://github.com/danpennic-sys/robotics-trust-profile/blob/main/docs/QUANTUMGUARD-BRIDGE.md

Also see `docs/RTP.md` in this repo for a short pointer.

## Operator Pack (commercial)

Source and CONTRACT-v1 remain free (MIT OR Apache-2.0).

**Operator Pack** adds signed binaries, air-gap runbooks, conformance kits, license files, and support under a commercial license.

- Private license repo (access after purchase): [danpennic-sys/quantumguard-operator](https://github.com/danpennic-sys/quantumguard-operator)
- Purchase / access: **danpennic@gmail.com**
- Claims stay honest: no PQC, no network trust, no wall-clock inside the verifier

See `MARKETING.md` for positioning.

## Non-goals (today)

- Post-quantum signatures
- Key registries / DNS / revocation
- Wall-clock or freshness checks
- Second verifier implementation
- Built-in payment integration in this public repo
- Domain geometry or kinematics (use RTP or other profiles)

## License

Licensed under either of

- Apache License, Version 2.0 ([LICENSE-APACHE](LICENSE-APACHE))
- MIT license ([LICENSE-MIT](LICENSE-MIT))

at your option.
