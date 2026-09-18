# QuantumGuard — Marketing Hub

**Status:** Ready for public push + commercial Operator Pack  
**Product truth:** v1 is classical (Ed25519), deterministic, offline. PQC is reserved, not shipped.

---

## 1. Positioning (one sentence)

QuantumGuard is a single deterministic evidence verifier that returns only PASS, FAIL, or INDETERMINATE — no network, no clocks, no external trust.

## 2. Core message

Most verification systems quietly depend on time, registries, or multiple code paths.  
QuantumGuard does not.

- One verifier
- Two interfaces (CLI + HTTP)
- Three possible answers
- Same bytes → same verdict, always

v1 is frozen and proven.  
Post-quantum (ML-DSA-65) is the next track, not a claim we make today.

## 3. Audience

Primary:
- Security engineers building offline / air-gapped attestation flows
- Teams that need replayable, auditable verification
- Builders who want a clean classical baseline before PQC migration

Secondary:
- Cryptography-curious developers
- Operators who hate non-deterministic tooling

## 4. What we claim (and what we don’t)

**We claim**
- Deterministic verification
- Single shared verifier for CLI and HTTP
- Explicit three-verdict model
- Fully offline operation
- Frozen, inspectable v1 contract
- Dual license (MIT OR Apache-2.0) for public source

**We do not claim**
- Post-quantum security (not shipped)
- Key distribution or trust anchors
- Freshness / timestamp validation
- Production hardening beyond the proven contract

## 5. GitHub storefront model

| Layer | Repo | Visibility | Role |
|-------|------|------------|------|
| Open core | [quantumguard](https://github.com/danpennic-sys/quantumguard) | Public | Source, CONTRACT-v1, vectors, CLI/HTTP |
| Operator Pack | [quantumguard-operator](https://github.com/danpennic-sys/quantumguard-operator) | **Private** | Signed binaries, runbooks, LICENSE.json, support |

**Fulfillment**

1. Buyer pays (Stripe Checkout or invoice).
2. Buyer gets GitHub collaborator access to `quantumguard-operator` **or** a signed tarball + `LICENSE.json`.
3. Public README links to purchase contact; private repo holds licensed assets.

**Pricing (list)**  
- Operator Pack single site: $490 one-time or $1,490/year  
- Air-gap / site license: custom  
Contact: danpennic@gmail.com

See private repo `docs/STOREFRONT.md` for operational checklist.

## 6. Short announcement (copy-paste)

**GitHub / X / LinkedIn**

QuantumGuard v1 is public.

A single deterministic evidence verifier.  
CLI and HTTP call the exact same code.  
Only three answers: PASS · FAIL · INDETERMINATE.  
No network. No wall-clock. No external trust.

Contract is frozen. Tests are green.  
PQC track (ML-DSA-65) is reserved, not invented.

Operator Pack (signed binaries + air-gap runbooks): contact danpennic@gmail.com

→ https://github.com/danpennic-sys/quantumguard

## 7. Longer announcement (blog / README hero)

Most “verifiers” are actually distributed systems in disguise.  
They call out to registries, check clocks, or maintain multiple code paths that slowly diverge.

QuantumGuard is the opposite.

v1 ships one pure function:

    evidence bytes → PASS | FAIL | INDETERMINATE

That function is used by both the CLI and the HTTP `/verify` endpoint.  
Same input always produces the same output. No network. No time. No second implementation.

The contract is frozen in `CONTRACT-v1.md`.  
Conformance vectors prove CLI ↔ HTTP equivalence.  
The gaps board in `PATCHES.md` stays honest about everything that is still missing — including post-quantum signatures.

This is the baseline.  
Everything else layers on top without breaking the invariants.

## 8. Channel plan (lightweight)

| Channel        | Action                                      | Priority |
|----------------|---------------------------------------------|----------|
| GitHub public  | Open core + storefront section in README    | Done     |
| GitHub private | Operator Pack license repo                  | Done     |
| Stripe         | Checkout for Operator Pack                  | Next     |
| X / Twitter    | Post short announcement                     | High     |
| LinkedIn       | Post longer version                         | Medium   |
| Docs / blog    | Optional deeper write-up later              | Low      |

No paid ads until first design partners. Product is the proof.

## 9. Success signal for this push

- Public repo clonable; `go test ./...` passes
- Private Operator Pack repo structured and invite-ready
- Announcement states only what is true
- No one can honestly accuse the project of claiming PQC before it exists

---

**Discipline:** Marketing must lag the product.  
Never sell the future track as present capability.
