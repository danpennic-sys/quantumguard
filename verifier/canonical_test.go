package verifier_test

import (
	"encoding/hex"
	"testing"

	"quantumguard/verifier"
)

func TestCanonicalBytes_MatchesV1SignedMessage(t *testing.T) {
	data := load(t, "01_known_good.json")
	got, err := verifier.CanonicalBytes(data)
	if err != nil {
		t.Fatalf("CanonicalBytes: %v", err)
	}
	wantHex := "71b8737c3c1128bb466c53e821f5cf1a1e765b82041a60c46672226b73a16dfc"
	if hex.EncodeToString(got) != wantHex {
		t.Fatalf("canonical bytes drifted from v1 payload_hash\n got %s\nwant %s", hex.EncodeToString(got), wantHex)
	}
	if len(got) != 32 {
		t.Fatalf("canonical bytes must be 32 bytes, got %d", len(got))
	}

	res := verifier.VerifyJSON(data)
	if res.Verdict != verifier.PASS {
		t.Fatalf("v1 must still PASS known-good, got %s", res.Verdict)
	}

	got2, err := verifier.CanonicalBytes(data)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(got) != hex.EncodeToString(got2) {
		t.Fatal("CanonicalBytes is not deterministic")
	}
}

func TestCanonicalBytes_Incomplete(t *testing.T) {
	_, err := verifier.CanonicalBytes(load(t, "05_incomplete.json"))
	if err == nil {
		t.Fatal("expected error for missing payload_hash")
	}
}

func TestCanonicalBytes_MalformedJSON(t *testing.T) {
	_, err := verifier.CanonicalBytes([]byte(`{not json`))
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}
