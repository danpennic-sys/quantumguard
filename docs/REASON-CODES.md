# Stable reason codes

**Status:** Normative for cross-implementation conformance.  
Prose `reason` / `detail` fields are non-normative.

The JSON field name is `code` (already shipped in CLI/HTTP and VM-EXPORT expected files).

## Core QuantumGuard v1

| case | verdict | code | notes |
|------|---------|------|-------|
| 01_known_good | PASS | OK | |
| 02_tampered_payload | FAIL | E_HASH_MISMATCH | |
| 03_signature_failure | FAIL | E_SIG_FAIL | |
| 04_hash_mismatch | FAIL | E_HASH_MISMATCH | **duplicate path of 02** — same check fires first; not a distinct code |
| 05_incomplete | INDETERMINATE | E_INCOMPLETE | |
| 06_unsupported_alg | INDETERMINATE | E_UNSUPPORTED_ALG | |

Also defined: `E_UNSUPPORTED_VERSION`, `E_MALFORMED_JSON`, `E_MALFORMED_ENCODING`.

## VM-EXPORT-BUNDLE-001

| case | verdict | code |
|------|---------|------|
| 01_valid_minimal | PASS | OK |
| 02_state_hash_mismatch | FAIL | STATE_HASH_MISMATCH |
| 03_replay_merkle_mismatch | FAIL | REPLAY_MERKLE_MISMATCH |
| 04_both_hashes_claimed | FAIL | BOTH_HASHES_CLAIMED |
| 05_extra_file | FAIL | EXTRA_FILE |
| 06_zip_slip | FAIL | ZIP_SLIP |
| 07_duplicate_path | FAIL | DUPLICATE_PATH |
| 08_created_at_tampered | FAIL | CREATED_AT_TAMPER |
| 09_honest_stub | FAIL | NOT_IMPLEMENTED |
| 10_canonical_order | PASS | OK |

## Cross-language note

A second implementation (e.g. PHP) that matches `verdict` + `code` on all vector cases has **vector-set conform** with the reference. That is not full product parity until both implement the same complete surface.
