package verifier_test

import (
	"os"
	"path/filepath"
	"testing"

	"quantumguard/verifier"
)

func load(t *testing.T, name string) []byte {
	t.Helper()
	// tests run with cwd = package dir or module root; try both
	candidates := []string{
		filepath.Join("..", "vectors", name),
		filepath.Join("vectors", name),
	}
	for _, p := range candidates {
		data, err := os.ReadFile(p)
		if err == nil {
			return data
		}
	}
	t.Fatalf("could not load vector %s", name)
	return nil
}

func TestKnownGood_PASS(t *testing.T) {
	res := verifier.VerifyJSON(load(t, "01_known_good.json"))
	if res.Verdict != verifier.PASS {
		t.Fatalf("expected PASS, got %s (%s)", res.Verdict, res.Reason)
	}
}

func TestTamperedPayload_FAIL(t *testing.T) {
	res := verifier.VerifyJSON(load(t, "02_tampered_payload.json"))
	if res.Verdict != verifier.FAIL {
		t.Fatalf("expected FAIL, got %s (%s)", res.Verdict, res.Reason)
	}
}

func TestSignatureFailure_FAIL(t *testing.T) {
	res := verifier.VerifyJSON(load(t, "03_signature_failure.json"))
	if res.Verdict != verifier.FAIL {
		t.Fatalf("expected FAIL, got %s (%s)", res.Verdict, res.Reason)
	}
}

func TestHashMismatch_FAIL(t *testing.T) {
	res := verifier.VerifyJSON(load(t, "04_hash_mismatch.json"))
	if res.Verdict != verifier.FAIL {
		t.Fatalf("expected FAIL, got %s (%s)", res.Verdict, res.Reason)
	}
}

func TestIncomplete_INDETERMINATE(t *testing.T) {
	res := verifier.VerifyJSON(load(t, "05_incomplete.json"))
	if res.Verdict != verifier.INDETERMINATE {
		t.Fatalf("expected INDETERMINATE, got %s (%s)", res.Verdict, res.Reason)
	}
}

func TestUnsupportedAlg_INDETERMINATE(t *testing.T) {
	res := verifier.VerifyJSON(load(t, "06_unsupported_alg.json"))
	if res.Verdict != verifier.INDETERMINATE {
		t.Fatalf("expected INDETERMINATE, got %s (%s)", res.Verdict, res.Reason)
	}
}

func TestMalformedJSON_INDETERMINATE(t *testing.T) {
	res := verifier.VerifyJSON([]byte(`{not json`))
	if res.Verdict != verifier.INDETERMINATE {
		t.Fatalf("expected INDETERMINATE, got %s", res.Verdict)
	}
}
