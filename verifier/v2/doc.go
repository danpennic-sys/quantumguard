// Package v2 is the QuantumGuard PQC-1 track (ML-DSA-65).
//
// PQC-1 is an additive, opt-in layer. It does not replace v1.
// The current implementation is a deterministic stub: it binds
// verification to verifier.CanonicalBytes (the same 32-byte
// payload_hash v1 signs) and returns INDETERMINATE until real
// ML-DSA-65 lands.
//
// See:
//   - ../../CONTRACT-v1.md          (frozen classical contract)
//   - ../../CONTRACT-PQC-1.md       (opt-in stub contract)
//   - ../../UPGRADE-v1-to-v2.md     (normative upgrade rules)
package v2
