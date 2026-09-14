package verifier

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// CanonicalBytes returns the exact message bytes the v1 verifier signs:
// the raw decoded payload_hash (32-byte SHA-256 digest as claimed in the
// bundle). PQC-1 must sign and verify this same sequence.
//
// No alternate message construction is permitted (PQC-TRACK Rule 3).
func CanonicalBytes(data []byte) ([]byte, error) {
	var b Bundle
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("malformed JSON")
	}
	return b.CanonicalBytes()
}

// CanonicalBytes returns the raw payload_hash bytes from a parsed bundle.
func (b Bundle) CanonicalBytes() ([]byte, error) {
	if strings.TrimSpace(b.PayloadHash) == "" {
		return nil, fmt.Errorf("missing payload_hash")
	}
	payloadHash, err := hex.DecodeString(b.PayloadHash)
	if err != nil {
		return nil, fmt.Errorf("payload_hash is not valid hex")
	}
	if len(payloadHash) != sha256.Size {
		return nil, fmt.Errorf("payload_hash must be %d bytes, got %d", sha256.Size, len(payloadHash))
	}
	return payloadHash, nil
}
