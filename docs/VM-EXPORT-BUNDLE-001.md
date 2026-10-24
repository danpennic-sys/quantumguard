# VM-EXPORT-BUNDLE-001 v1.0

**Status:** Design submitted. Dual hash construction (`state_hash` + `replay_merkle_root`) selected. Vector contract published. Implementation of both paths pending. Cross-language product parity still open.

**Date:** 2026-10-23  
**Owner:** Operator Pages / QuantumGuard  
**Related:** VaultMax design notes, quantumguard verifier, PACKGATE evidence packs  
**License:** dual MIT OR Apache-2.0

---

## 1. Purpose

Define a signed, offline-verifiable export bundle format and the corresponding deterministic verifier rules. The bundle is an evidence container; the verifier produces only PASS / FAIL / INDETERMINATE with no silent success.

## 2. Dual Hash Construction (selected)

Two independent hash constructions are supported. A bundle may claim **at most one**.

| Kind | Definition | Source |
|------|------------|--------|
| `state_hash` | `SHA-256(canonical_encode(S, M, T))` | VaultMax `02_canonical_encoding.md` |
| `replay_merkle_root` | Merkle root over event list (domain-separated) | Export-bundle event-log path |

### 2.1 state_hash (VaultMax path)

- Encoding: UTF-8, lexicographically sorted keys, deterministic array order, NFC, no whitespace (RFC 8785 / JCS).
- Input: the state triple `(S, M, T)`.
- No additional domain separation prefix.

### 2.2 replay_merkle_root (event-log path)

- Leaf = hex-decoded `event.Hash` bytes (no extra SHA-256).
- Internal node = `SHA-256(0x01 || left || right)`.
- Odd leaf carried up unchanged (Bitcoin-style).

### 2.3 Verifier rules

- Claiming both roots → FAIL.
- Claiming neither and no corresponding `--no-*` flag → recorded as omitted (honest failure unless suppressed).
- Same pattern already used for ADR-anchor and ledger-proof checks.

## 3. Integrity rules already decided

| Rule | Decision |
|------|----------|
| `created_at` | Inside signed region (part of `bundle.json` → `root_hash`). Self-reported; authoritative time is the signed ledger chain. |
| Extra files | REJECT. Only files listed in the manifest (plus a small allow-list of signature/public-key files) are permitted. |
| Zip-slip | Enforced. Flat layout only; subdirectories rejected. |
| Resource limits | File count, single-file size, total uncompressed size enforced before extraction. |
| Unimplemented checks | See §3.1 Honest-stub contract. |

### 3.1 Honest-stub contract (normative)

If a check is required by the bundle or by the caller and the implementation has not yet implemented that check, the verifier **MUST** return:

- `verdict`: `FAIL`
- `code`: `NOT_IMPLEMENTED`

It **MUST NOT** return `PASS` for a check it did not perform.  
Vector `09_honest_stub` encodes this contract. It is not a coverage hole.

### 3.2 Check-order precedence (normative)

On a multi-error bundle the verifier **MUST** evaluate checks in this order and **MUST** stop at the first failure (fail-closed, first-error wins):

1. **Both hashes claimed** → `BOTH_HASHES_CLAIMED`
2. **Path validity (zip-slip)** → `ZIP_SLIP`
3. **Duplicate path** → `DUPLICATE_PATH`
4. **Extra-file rejection** → `EXTRA_FILE`
5. **Required honest stubs** (e.g. ADR anchor when required) → `NOT_IMPLEMENTED`
6. **Hash verification** (`state_hash` or `replay_merkle_root`) → `STATE_HASH_MISMATCH` / `REPLAY_MERKLE_MISMATCH` / `CREATED_AT_TAMPER` / `OK`

A third implementation that uses a different order is non-conforming on multi-error inputs even if single-error vectors pass.

## 4. Gaps

| Gap ID | Description | Status |
|--------|-------------|--------|
| GAP-005 | Build reproducibility | Design-verified — execution pending |
| GAP-011 | Trust-anchor pinning | Design-verified — execution pending |
| GAP-014 | Post-quantum readiness | Open — v1.1 required |
| GAP-019 | Cross-language product parity | Open — vector-set conform (Python+PHP) exists; full surface pending |
| GAP-020 | Dual-hash dispatch + omission flags | Design submitted |

## 5. Non-claims

- This document is design + vector contract. Full production implementation is not claimed.
- Cross-language **product** parity is not verified; **vector-set** conform (Python + PHP on the 10 cases) is.
- VaultMax Python stack (`reset_15.py` et al.) has missing modules; full-stack continuity claims remain unverifiable until those modules exist.

## 6. Next implementation steps (ordered)

1. Land dual-hash helpers + dispatch in `verifier/` (or a new `exportbundle/` package).
2. Obtain the corresponding Rust / substrate construction and close GAP-019 product surface.
3. Wire ADR-anchor and ledger-proof to real checks once their specs are frozen (until then §3.1 applies).

---

**Authoritative status line**

VM-EXPORT-BUNDLE-001 v1.0 — Design + vector contract. Dual-hash selected. Check-order and honest-stub contract normative. License dual MIT OR Apache-2.0. Vector-set conform 10/10 (Python + PHP). Cross-language product parity still open. v2 PQC: 0/2 stubs only.
