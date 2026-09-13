# Gaps Board — Updated 2026-09-12 (proxy retired)

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

Normative docs:
- `CONTRACT-v1.md`
- `UPGRADE-v1-to-v2.md`

## PQC Track (v2)

| Item                              | Status                |
|-----------------------------------|-----------------------|
| ML-DSA-65 decision                | LOCKED                |
| Type-Prefixed / versioned design  | DOCUMENTED            |
| `verifier/v2` package             | SCAFFOLD ONLY         |
| filippo.io/mldsa integration      | NOT STARTED           |
| v2 conformance vectors            | NOT CREATED           |
| Actual PQC verification code      | NOT STARTED           |

## Still Open

| Item                        | Status                          |
|-----------------------------|---------------------------------|
| yourmodule L10 resolution   | UNRESOLVED                      |
| CP-01 → CP-08               | EMPTY                           |
| Chat proxy payload          | **RETIRED 2026-09-12** — operatorpagesllc archived path; not in this tree |
| Stripe / commercial path    | OFFER POSTED; checkout not live |
| Tagged release              | NOT DONE                        |

## Discipline

- v1 is read-only.
- All PQC work lives under the v2 track and must obey `UPGRADE-v1-to-v2.md`.
- No silent modification of frozen behavior.
- Replay harness remains the arbiter of correctness.
- Chat proxy retirement is not a verifier change.
