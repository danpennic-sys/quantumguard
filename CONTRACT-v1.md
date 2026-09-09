# QuantumGuard Contract v1 (Frozen)

**Status:** Frozen — 2026-09-09  
**Scope:** Classical Ed25519 only  
**Invariant:** Single deterministic verifier, offline, no wall-clock, no external trust

## Evidence Bundle Format

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

### Encoding Rules (normative)

| Field          | Encoding                                      | Size / Notes                          |
|----------------|-----------------------------------------------|---------------------------------------|
| `version`      | ASCII string                                  | Must be exactly `"1"`                 |
| `payload`      | Standard base64 (RFC 4648)                    | Raw original bytes                    |
| `payload_hash` | Lowercase hex                                 | Exactly 32 bytes (SHA-256)            |
| `signature`    | Lowercase hex                                 | Exactly 64 bytes (Ed25519)            |
| `public_key`   | Lowercase hex                                 | Exactly 32 bytes (raw Ed25519 pubkey) |
| `alg`          | ASCII string                                  | Must be exactly `"ed25519"`           |

## Verdict Rules (normative)

| Condition                                              | Verdict          |
|--------------------------------------------------------|------------------|
| Bundle complete + encodings valid + hash matches + signature verifies | `PASS`     |
| Bundle complete + encodings valid but hash or signature wrong         | `FAIL`     |
| Missing fields, malformed JSON, bad encoding, wrong version, unsupported alg | `INDETERMINATE` |

## Architectural Invariants

1. Exactly one verifier implementation.
2. CLI and HTTP call the identical verifier function.
3. No network, DNS, registry, revocation, or wall-clock dependency.
4. Verification is a pure function of the evidence bytes.
5. Only three possible outputs: `PASS` | `FAIL` | `INDETERMINATE`.

## Explicit Non-Goals of v1

- Post-quantum signatures
- Multiple algorithms
- Key distribution or trust anchors
- Freshness / timestamp checks
- Revocation

Any future version (v2+) must preserve the single-verifier and determinism invariants.
