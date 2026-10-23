# VM-EXPORT-BUNDLE-001 Conformance Vectors

**Status:** Design + vector contract submitted. Implementation pending.  
**Related:** `docs/VM-EXPORT-BUNDLE-001.md`

These vectors are the executable contract for the dual-hash export-bundle verifier.
They are language-neutral. Any conforming implementation (Go, Rust, Python, …)
must produce the listed verdict and error code for each case.

## Manifest

See `vectors/VM-EXPORT-BUNDLE-001/manifest.json`.

## Case schema

Each case under `vectors/VM-EXPORT-BUNDLE-001/cases/<case_id>/` contains:

| File | Purpose |
|------|---------|
| `input.json` | Minimal fixture describing the bundle claim under test |
| `expected.json` | Required verdict, error code, and any computed hashes |

## Verdict values

- `PASS`
- `FAIL`
- `INDETERMINATE` (reserved; not used in the current set)

## Stable error codes

| Code | Meaning |
|------|---------|
| `OK` | All required checks passed |
| `STATE_HASH_MISMATCH` | Claimed state_hash does not match recomputed value |
| `REPLAY_MERKLE_MISMATCH` | Claimed replay_merkle_root does not match recomputed value |
| `BOTH_HASHES_CLAIMED` | Bundle claims both state_hash and replay_merkle_root |
| `EXTRA_FILE` | Unauthorized file present in the bundle |
| `ZIP_SLIP` | Path traversal or subdirectory detected |
| `DUPLICATE_PATH` | Same path appears more than once |
| `CREATED_AT_TAMPER` | created_at altered after signing (root_hash/signature break) |
| `NOT_IMPLEMENTED` | Required check is an honest stub |
| `CANONICAL_ORDER` | JSON key order / whitespace variation changes the hash |

## Non-claims

- These vectors do not constitute a complete implementation.
- Cross-language parity is still open.
- The optional verifier skeleton is verify-only and fail-closed.
