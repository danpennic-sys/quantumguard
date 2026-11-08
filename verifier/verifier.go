// SPDX-License-Identifier: MIT OR Apache-2.0
// Copyright (c) 2026 danpennic-sys

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
//   - Malformed / incomplete / unsupported → INDETERMINATE
//   - Well-formed but cryptographically wrong → FAIL
//   - Well-formed and cryptographically valid → PASS
//
// Reason codes are stable, machine-readable, and normative for cross-implementation
// conformance. The Reason string is human-readable detail and non-normative.
package verifier

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	maxJSONBytes    = 1 << 20   // 1 MiB raw JSON
	maxPayloadBytes = 512 << 10 // 512 KiB decoded payload
	hexHashLen      = sha256.Size * 2
	hexSigLen       = ed25519.SignatureSize * 2
	hexPubLen       = ed25519.PublicKeySize * 2
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
	Payload     string `json:"payload"`
	PayloadHash string `json:"payload_hash"`
	Signature   string `json:"signature"`
	PublicKey   string `json:"public_key"`
	Alg         string `json:"alg"`
}

// Result is the deterministic verification outcome.
type Result struct {
	Verdict Verdict `json:"verdict"`
	Code    string  `json:"code"`
	Reason  string  `json:"reason,omitempty"`
}

// VerifyJSON accepts a raw JSON evidence bundle and returns a deterministic Result.
// This is the single entry point used by both CLI and HTTP.
func VerifyJSON(data []byte) Result {
	if len(data) > maxJSONBytes {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedEncoding, Reason: "JSON exceeds size limit"}
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	var b Bundle
	if err := dec.Decode(&b); err != nil {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedJSON, Reason: "malformed JSON"}
	}
	if dec.More() {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedJSON, Reason: "trailing data after JSON"}
	}
	return Verify(b)
}

// Verify performs deterministic verification of a parsed Bundle.
func Verify(b Bundle) Result {
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

	if len(b.PayloadHash) != hexHashLen || !isLowerHex(b.PayloadHash) {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedEncoding, Reason: "payload_hash must be lowercase hex of 32 bytes"}
	}
	if len(b.Signature) != hexSigLen || !isLowerHex(b.Signature) {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedEncoding, Reason: "signature must be lowercase hex of 64 bytes"}
	}
	if len(b.PublicKey) != hexPubLen || !isLowerHex(b.PublicKey) {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedEncoding, Reason: "public_key must be lowercase hex of 32 bytes"}
	}

	if len(b.Payload) > maxPayloadBytes*2 {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedEncoding, Reason: "payload exceeds size limit"}
	}

	payload, err := base64.StdEncoding.DecodeString(b.Payload)
	if err != nil {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedEncoding, Reason: "payload is not valid base64"}
	}
	if len(payload) > maxPayloadBytes {
		return Result{Verdict: INDETERMINATE, Code: CodeMalformedEncoding, Reason: "payload exceeds size limit"}
	}

	payloadHash, _ := hex.DecodeString(b.PayloadHash)
	sig, _ := hex.DecodeString(b.Signature)
	pub, _ := hex.DecodeString(b.PublicKey)

	computed := sha256.Sum256(payload)
	if subtle.ConstantTimeCompare(computed[:], payloadHash) != 1 {
		return Result{Verdict: FAIL, Code: CodeHashMismatch, Reason: "payload_hash does not match SHA-256 of payload"}
	}

	if !ed25519.Verify(ed25519.PublicKey(pub), payloadHash, sig) {
		return Result{Verdict: FAIL, Code: CodeSigFail, Reason: "signature verification failed"}
	}

	return Result{Verdict: PASS, Code: CodeOK, Reason: "valid"}
}

func isLowerHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}
