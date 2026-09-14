package v2

import (
	"encoding/json"
	"io"

	"quantumguard/verifier"
)

// Result is the PQC-1 verification outcome. It never replaces the v1 verdict.
type Result struct {
	Verdict verifier.Verdict `json:"verdict"`
	Reason  string           `json:"reason,omitempty"`
	Alg     string           `json:"alg"`
	Stub    bool             `json:"stub"`
}

// Combined is the CLI/HTTP envelope.
//
// v1 fields (verdict, reason) are produced solely by verifier.VerifyJSON.
// PQC is null when the opt-in flag is disabled, and never mutates v1.
type Combined struct {
	Verdict verifier.Verdict `json:"verdict"`
	Reason  string           `json:"reason,omitempty"`
	PQC     *Result          `json:"pqc"`
}

const (
	algMLDSA65      = "ml-dsa-65"
	stubReason      = "pqc-1 stub: ML-DSA-65 not implemented"
	canonicalPrefix = "canonical_bytes: "
)

// VerifyPQC is the deterministic PQC-1 stub.
//
// It always signs/verifies via verifier.CanonicalBytes (the same message
// bytes v1 uses as its Ed25519 input). Real ML-DSA-65 is not implemented:
// a well-formed message yields INDETERMINATE with a fixed stub reason.
func VerifyPQC(data []byte) Result {
	msg, err := verifier.CanonicalBytes(data)
	if err != nil {
		return Result{
			Verdict: verifier.INDETERMINATE,
			Reason:  canonicalPrefix + err.Error(),
			Alg:     algMLDSA65,
			Stub:    true,
		}
	}
	// CanonicalBytes already enforces 32-byte payload_hash. Re-check so
	// this path cannot skip the canonical message.
	if len(msg) != 32 {
		return Result{
			Verdict: verifier.INDETERMINATE,
			Reason:  canonicalPrefix + "payload_hash must be 32 bytes",
			Alg:     algMLDSA65,
			Stub:    true,
		}
	}
	return Result{
		Verdict: verifier.INDETERMINATE,
		Reason:  stubReason,
		Alg:     algMLDSA65,
		Stub:    true,
	}
}

// Verify runs the frozen v1 verifier, then optionally the PQC-1 stub.
// enablePQC does not change the v1 verdict, reason, or evidence semantics.
func Verify(data []byte, enablePQC bool) Combined {
	v1 := verifier.VerifyJSON(data)
	out := Combined{
		Verdict: v1.Verdict,
		Reason:  v1.Reason,
		PQC:     nil,
	}
	if enablePQC {
		p := VerifyPQC(data)
		out.PQC = &p
	}
	return out
}

// Encode writes Combined with the same indent used by CLI and HTTP.
func Encode(w io.Writer, c Combined) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(c)
}
