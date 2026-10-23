// Package verifier implements the single deterministic QuantumGuard evidence verifier.
//
// Contract v1 (locked):
//
//	{
//	  "version": "1",
//	  "payload": "<base64-encoded raw payload bytes>",
//	  "payload_hash": "<hex-encoded SHA-256 of the raw payload bytes>",
//	  "signature": "<hex-encoded 64-byte Ed25519 signature over the raw 32-byte payload_hash>",
//	  "public_key": "<hex-encoded 32-byte raw Ed25519 public key>",
//	  "alg": "ed25519"
//	}
//
// Verdict rules (no network, no wall-clock, no external trust):
//
//	- Malformed / incomplete / unsupported → INDETERMINATE
//	- Well-formed but cryptographically wrong → FAIL
//	- Well-formed and cryptographically valid → PASS
//
// Reason codes are stable, machine-readable, and normative for cross-implementation
// conformance. The Reason string is human-readable detail and non-normative.
package verifier

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// Verdict is the only allowed output of verification.
type Verdict string

const (
	PASS          Verdict = "PASS"
	FAIL          Verdict = "FAIL"
	INDETERMINATE Verdict = "INDETERMINATE"
)

// Stable reason codes (normative). Prose Reason is non-normative detail.
const (
	CodeOK                = "OK"
	CodeHashMismatch      = "E_HASH_MISMATCH"
	CodeSigFail           = "E_SIG_FAIL"
	CodeIncomplete        = "E_INCOMPLETE"
	CodeUnsupportedAlg    = "E_UNSUPPORTED_ALG"
	CodeUnsupportedVer    = "E_UNSUPPORTED_VERSION"
	CodeMalformedJSON     = "E_MALFORMED_JSON"
	CodeMalformedEncoding = "E_MALFORMED_ENCODING"
)

// Bundle is the v1 evidence bundle.
type Bundle struct {
	Version     string `json:"version"`
	Payload     string `json:"payload"`      // base64
	PayloadHash string `json:"payload_hash"` // hex
	Signature   string `json:"signature"`    // hex
	PublicKey   string `json:"public_key"`   // hex, 32-byte raw Ed25519
	Alg         string `json:"alg"`
}

// Result is the deterministic verification outcome.
type Result struct {
	Verdict Verdict `json:"verdict"`
	Code    string  `json:"code"`             // stable, normative
	Reason  string  `json:"reason,omitempty"` // prose, non-normative
}

// VerifyJSON accepts a raw JSON evidence bundle and returns a deterministic Result.
// This is the single entry point used by both CLI and HTTP.
func VerifyJSON(data []byte) Result {
	var b Bundle
	if err := json.Unmarshal(data, &b); err != nil {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedJSON, Reason: "malformed JSON"}
	}
	return Verify(b)
}

// Verify performs deterministic verification of a parsed Bundle.
func Verify(b Bundle) Result {
	// --- Structural / completeness checks → INDETERMINATE ---
	if strings.TrimSpace(b.Version) == "" ||
		strings.TrimSpace(b.Payload) == "" ||
		strings.TrimSpace(b.PayloadHash) == "" ||
		strings.TrimSpace(b.Signature) == "" ||
		strings.TrimSpace(b.PublicKey) == "" ||
		strings.TrimSpace(b.Alg) == "" {
		return Result{Verdict: INDETERMINATE, Code: CodeIncomplete, Reason: "incomplete bundle: missing required field"}
	}

	if b.Version != "1" {
		return Result{Verdict: INDETERMINATE, Code: CodeUnsupportedVer, Reason: fmt.Sprintf("unsupported version: %q", b.Version)}
	}

	if b.Alg != "ed25519" {
		return Result{Verdict: INDETERMINATE, Code: CodeUnsupportedAlg, Reason: fmt.Sprintf("unsupported alg: %q", b.Alg)}
	}

	// Decode fields. Any decode failure → INDETERMINATE (malformed encoding).
	payload, err := base64.StdEncoding.DecodeString(b.Payload)
	if err != nil {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedEncoding, Reason: "payload is not valid base64"}
	}

	payloadHash, err := hex.DecodeString(b.PayloadHash)
	if err != nil {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedEncoding, Reason: "payload_hash is not valid hex"}
	}
	if len(payloadHash) != sha256.Size {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedEncoding, Reason: fmt.Sprintf("payload_hash must be %d bytes, got %d", sha256.Size, len(payloadHash))}
	}

	sig, err := hex.DecodeString(b.Signature)
	if err != nil {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedEncoding, Reason: "signature is not valid hex"}
	}
	if len(sig) != ed25519.SignatureSize {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedEncoding, Reason: fmt.Sprintf("signature must be %d bytes, got %d", ed25519.SignatureSize, len(sig))}
	}

	pub, err := hex.DecodeString(b.PublicKey)
	if err != nil {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedEncoding, Reason: "public_key is not valid hex"}
	}
	if len(pub) != ed25519.PublicKeySize {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedEncoding, Reason: fmt.Sprintf("public_key must be %d bytes, got %d", ed25519.PublicKeySize, len(pub))}
	}

	// --- Cryptographic checks → FAIL or PASS ---

	// 1. Hash must match
	computed := sha256.Sum256(payload)
	if !hmacEqual(computed[:], payloadHash) {
		return Result{Verdict: FAIL, Code: CodeHashMismatch, Reason: "payload_hash does not match SHA-256 of payload"}
	}

	// 2. Signature must verify over the raw payload_hash bytes
	if !ed25519.Verify(ed25519.PublicKey(pub), payloadHash, sig) {
		return Result{Verdict: FAIL, Code: CodeSigFail, Reason: "signature verification failed"}
	}

	return Result{Verdict: PASS, Code: CodeOK, Reason: "valid"}
}

// hmacEqual is a constant-time comparison.
func hmacEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var v byte
	for i := 0; i < len(a); i++ {
		v |= a[i] ^ b[i]
	}
	return v == 0
}
