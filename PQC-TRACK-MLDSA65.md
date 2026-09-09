# PQC Track — ML-DSA-65 Initialization (v1 → PQC-1)

**Status:** Initialized  
**Date:** 2026-09-09  
**Cadence:** Quiet. Deterministic. No drift.

v1 remains frozen.  
PQC work happens outside the v1 contract.

---

## 1. Purpose

Establish a parallel, non-intrusive PQC track for QuantumGuard using ML-DSA-65, without modifying or weakening the frozen v1 contract.

- v1 remains the stable verifier.
- PQC-1 becomes the experimental track.

---

## 2. Rules of the PQC Track

These rules prevent drift, hype, or accidental v1 mutation.

**Rule 1 — v1 is immutable**  
No PQC code enters v1.  
No PQC claims appear in v1 documentation.  
No PQC behavior affects v1 outputs.

**Rule 2 — PQC-1 is a separate contract**  
A new file `CONTRACT-PQC-1.md` will define PQC behavior, inputs, outputs, and invariants.  
It does not exist yet.

**Rule 3 — PQC signatures operate on the same message bytes**  
ML-DSA-65 signs the same bytes that v1 uses for its signature input (the raw `payload_hash` bytes in the current v1 design).  
No alternate message construction.

**Rule 4 — PQC verifier is deterministic**  
No network.  
No clocks.  
No external trust.  
Same discipline as v1.

**Rule 5 — PQC track must pass the replay harness**  
Conformance vectors + tests remain the arbiter of correctness for PQC-1.

---

## 3. Initial PQC-1 Scope (locked)

### Included

- ML-DSA-65 key generation
- ML-DSA-65 signing over the same message bytes used by v1
- ML-DSA-65 verification
- PQC golden fixtures
- PQC extension to the test/replay harness
- Future CLI flag: `--pqc` (not implemented yet)
- Future HTTP parameter: `?pqc=true` (not implemented yet)

### Excluded

- No PQC-only evidence formats
- No PQC-only chain entries
- No PQC-only Merkle variants
- No PQC threshold crypto
- No PQC claims in v1 marketing
- No PQC claims in v1 contract
- No modification of frozen v1 code paths

---

## 4. Placement in the Repository

```
quantumguard/
├── CONTRACT-v1.md              # frozen (read-only)
├── UPGRADE-v1-to-v2.md         # normative upgrade rules
├── PQC-TRACK-MLDSA65.md        # this document
├── verifier/                   # v1 only — do not put PQC here
├── verifier/v2/                # reserved scaffold for PQC track
│   ├── README.md
│   └── doc.go
└── ...
```

All future PQC implementation code belongs under `verifier/v2/` (or a clearly named `pqc/` package) and must not alter `verifier/verifier.go`.

---

## 5. Data Types (draft — not yet normative)

These are planning placeholders only. No code yet.

```
PQCPublicKey  — ML-DSA-65 public key bytes
PQCPrivateKey — ML-DSA-65 private key bytes
PQCSignature  — ML-DSA-65 signature bytes
```

Exact encoding, size, and serialization rules will be defined in `CONTRACT-PQC-1.md` when that document is written.

---

## 6. Core Functions (draft — not yet normative)

```
pqc_sign(message []byte, sk PQCPrivateKey) → PQCSignature
pqc_verify(message []byte, sig PQCSignature, pk PQCPublicKey) → bool
```

Invariant: the only message signed is the same byte sequence v1 uses for its signature input.

---

## 7. Golden Fixture (planned)

```
vectors/v2/pqc_mldsa65_golden.json
```

Structure (indicative):

```json
{
  "message": "hex...",
  "public_key": "hex...",
  "signature": "hex...",
  "valid": true
}
```

The test suite must reject any drift from the golden result.

---

## 8. CLI / HTTP Behavior (future)

When (and only when) implemented:

**CLI**
```
qg verify --pqc <bundle.json>
```

**HTTP**
```
POST /verify?pqc=true
```

Expected combined output shape (indicative):

```json
{
  "v1": "PASS",
  "pqc": "PASS"
}
```

v1 result remains authoritative.  
PQC result is additional information only.

---

## 9. Messaging Discipline

- QuantumGuard does **not** claim PQC support in v1.
- QuantumGuard does **not** claim ML-DSA-65 is production-ready.
- QuantumGuard does **not** claim quantum resistance.

PQC-1 is experimental, parallel, and versioned.

---

## 10. Current State

| Item                         | Status        |
|------------------------------|---------------|
| This initialization document | Written       |
| CONTRACT-PQC-1.md            | Not written   |
| ML-DSA-65 implementation     | Not started   |
| Golden fixtures              | Not created   |
| CLI / HTTP PQC flags         | Not started   |
| v1 code                      | Untouched     |

---

## 11. Next Vectors (choose deliberately)

- Draft `CONTRACT-PQC-1.md`
- Create keygen / sign / verify stubs under `verifier/v2/`
- Add PQC golden fixtures
- Extend test harness for PQC
- Or stop and hold

Quiet cadence. No drift.
