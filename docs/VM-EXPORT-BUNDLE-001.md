# VM-EXPORT-BUNDLE-001 v1.0

**Status:** Design submitted. Dual hash construction (`state_hash` + `replay_merkle_root`) selected. Implementation of both paths pending. Cross-language parity still open.

**Date:** 2026-10-23  
**Owner:** Operator Pages / QuantumGuard  
**Related:** VaultMax design notes, quantumguard verifier, PACKGATE evidence packs

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
| Unimplemented checks | Must return `Pass: false` with explicit "NOT IMPLEMENTED" detail. Never silent PASS. |

## 4. Gaps

| Gap ID | Description | Status |
|--------|-------------|--------|
| GAP-005 | Build reproducibility | Design-verified — execution pending |
| GAP-011 | Trust-anchor pinning | Design-verified — execution pending |
| GAP-014 | Post-quantum readiness | Open — v1.1 required |
| GAP-019 | Cross-language root parity (Go ↔ Rust) | Open — both constructions need confirmation |
| GAP-020 | Dual-hash dispatch + omission flags | Design submitted — this document |

## 5. Non-claims

- This document is **design only**. No implementation is claimed complete.
- Cross-language parity is **not** verified.
- The Go fragments produced during design review compile in isolation; they are not yet integrated into the quantumguard verifier package.
- VaultMax Python stack (`reset_15.py` et al.) has missing modules; full-stack continuity claims remain unverifiable until those modules exist.

## 6. Next implementation steps (ordered)

1. Land dual-hash helpers + dispatch in `verifier/` (or a new `exportbundle/` package).
2. Add conformance vectors for both hash kinds and the dual-claim rejection case.
3. Obtain the corresponding Rust / substrate construction and close GAP-019.
4. Wire ADR-anchor and ledger-proof as explicit NOT-IMPLEMENTED until their specs are frozen.

---

**Authoritative status line**

VM-EXPORT-BUNDLE-001 v1.0 — Design submitted. Dual hash construction (state_hash + replay_merkle_root) selected. Implementation of both paths pending. Cross-language parity still open.
