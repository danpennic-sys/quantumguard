# Robotics Trust Profile (RTP) and QuantumGuard

QuantumGuard does not interpret LiDAR, odometry, or joint limits.  
[RTP v1.0](https://github.com/danpennic-sys/robotics-trust-profile) does.

## Division of labor

| Component | Verdict / decision | Meaning |
|-----------|-------------------|--------|
| QuantumGuard CONTRACT-v1 | `PASS` / `FAIL` / `INDETERMINATE` | Evidence bytes authentic? |
| RTP validators | `PASS` / `FAIL` / `INDETERMINATE` | Observation admissible under policy? |
| GIS / RTP decision | `ACCEPT` / `QUARANTINE` / `REJECT` / `SUPERSEDE` | What may controllers do? |

## Integrity interlock (recommended)

Do not grant actuation permission unless:

1. Every **critical** QuantumGuard bundle verifies `PASS`, and
2. RTP `decision_state` is `ACCEPT` (or site-approved `SUPERSEDE`).

A QuantumGuard `FAIL` must not be overridden by an RTP `ACCEPT`.

## Deep documentation

Full canonicalization rules, dual-verdict audit schema, failure matrix, and worked AGV pipeline:

https://github.com/danpennic-sys/robotics-trust-profile/blob/main/docs/QUANTUMGUARD-BRIDGE.md

Examples: `docs/bridge-examples/` in the RTP repository.
