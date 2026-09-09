# Gaps Board — Updated 2026-09-09 (post-freeze)

## Frozen / Done

| Item                              | Status      |
|-----------------------------------|-------------|
| QuantumGuard Contract v1          | **FROZEN**  |
| Shared deterministic verifier     | DONE        |
| CLI (`qg verify`)                 | DONE        |
| HTTP `/verify`                    | DONE        |
| CLI ↔ HTTP equivalence            | PROVEN      |
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

## Still Open (unchanged)

| Item                        | Status                          |
|-----------------------------|---------------------------------|
| yourmodule L10 resolution   | UNRESOLVED                      |
| CP-01 → CP-08               | EMPTY                           |
| Chat proxy payload          | NIL / BROKEN                    |
| Stripe / commercial path    | NOT STARTED                     |
| Tagged release              | NOT DONE                        |

## Discipline

- v1 is read-only.
- All PQC work lives under the v2 track and must obey `UPGRADE-v1-to-v2.md`.
- No silent modification of frozen behavior.
- Replay harness remains the arbiter of correctness.
