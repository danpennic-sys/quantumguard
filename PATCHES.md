# Gaps Board — Updated 2026-10-28 (checkout live)

Checked against `main` and tag `v1.0.0` (`28045276`).
`go test -race ./...` re-run locally on go1.23.12 linux/amd64 on 2026-10-28: PASS.
The six v1 vectors were re-run through `qg verify`: 01 PASS, 02–04 FAIL, 05–06 INDETERMINATE.

## Frozen / Done

| Item                              | Status      |
|-----------------------------------|-------------|
| QuantumGuard Contract v1          | **FROZEN**  |
| Shared deterministic verifier     | DONE        |
| CLI (`qg verify`)                 | DONE        |
| HTTP `/verify`                    | DONE        |
| CLI ↔ HTTP equivalence            | BY CONSTRUCTION |
| Conformance vectors (v1)          | DONE        |
| Three-verdict model               | LOCKED      |
| v1 → v2 Upgrade Contract          | **WRITTEN** |
| Tagged release `v1.0.0`           | **DONE** (published 2026-09-28). Unsigned binaries + `SHA256SUMS.txt`. Not HSM-signed. Do not move this tag. |
| Go CI on `main`                   | **DONE** (`.github/workflows/ci.yml`). Node publish workflow removed 2026-10-26. |
| Operator Pack checkout (one-time) | **LIVE** — https://buy.stripe.com/5kQaEYc9AewBaJleoN0ZW01 ($490) |
| Operator Pack checkout (yearly)   | **LIVE** — https://buy.stripe.com/cNibJ24H8gEJcRteoN0ZW02 ($1,490/year) |

Normative docs:
- `CONTRACT-v1.md`
- `UPGRADE-v1-to-v2.md`

`main` has moved past `v1.0.0` (reason codes, VM-EXPORT-BUNDLE-001, the Rust ML-DSA crate, CI). That delta is not a new release. Cutting another tag is a product decision, not a docs fix.

## PQC Track (v2)

| Item                              | Status                |
|-----------------------------------|-----------------------|
| ML-DSA-65 decision                | LOCKED                |
| Type-Prefixed / versioned design  | DOCUMENTED            |
| `CONTRACT-PQC-1.md`               | WRITTEN (stub)        |
| `verifier/v2` package             | STUB SHIPPED          |
| CLI `--pqc` / HTTP `?pqc=true`    | SHIPPED (opt-in stub) |
| Envelope goldens (`vectors/v2`)   | SHIPPED               |
| filippo.io/mldsa integration      | NOT STARTED           |
| Actual ML-DSA-65 verification     | NOT STARTED           |

## Still Open

| Item                        | Status                                                                 |
|-----------------------------|------------------------------------------------------------------------|
| Chat proxy payload          | **RETIRED 2026-09-12** — operatorpagesllc archived path; not in this tree |
| Next tag after `v1.0.0`     | NOT CUT. Do not retag or move `v1.0.0`.                                |

## Discipline

- v1 is read-only.
- All PQC work lives under the v2 track and must obey `UPGRADE-v1-to-v2.md`.
- No silent modification of frozen behavior.
- Replay harness remains the arbiter of correctness.
- Chat proxy retirement is not a verifier change.
- One verdict path. `rust/qg-mldsa65` and the VM-EXPORT tools do not replace `qg`.
- Fulfillment remains manual (collaborator invite or signed drop). No webhook auto-mint.
