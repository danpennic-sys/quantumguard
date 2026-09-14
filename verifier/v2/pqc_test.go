package v2_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"quantumguard/verifier"
	v2 "quantumguard/verifier/v2"
)

func load(t *testing.T, name string) []byte {
	t.Helper()
	candidates := []string{
		filepath.Join("..", "..", "vectors", name),
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

func loadV2(t *testing.T, name string) []byte {
	t.Helper()
	candidates := []string{
		filepath.Join("..", "..", "vectors", "v2", name),
		filepath.Join("vectors", "v2", name),
	}
	for _, p := range candidates {
		data, err := os.ReadFile(p)
		if err == nil {
			return data
		}
	}
	t.Fatalf("could not load v2 vector %s", name)
	return nil
}

func TestVerifyPQC_UsesCanonicalBytes(t *testing.T) {
	good := v2.VerifyPQC(load(t, "01_known_good.json"))
	if good.Verdict != verifier.INDETERMINATE {
		t.Fatalf("stub must be INDETERMINATE, got %s", good.Verdict)
	}
	if good.Reason != "pqc-1 stub: ML-DSA-65 not implemented" {
		t.Fatalf("unexpected stub reason: %s", good.Reason)
	}
	if good.Alg != "ml-dsa-65" || !good.Stub {
		t.Fatalf("stub metadata: alg=%s stub=%v", good.Alg, good.Stub)
	}

	incomplete := v2.VerifyPQC(load(t, "05_incomplete.json"))
	if incomplete.Verdict != verifier.INDETERMINATE {
		t.Fatalf("incomplete stub must be INDETERMINATE, got %s", incomplete.Verdict)
	}
	if incomplete.Reason == good.Reason {
		t.Fatal("PQC stub did not consult canonical_bytes: incomplete and known-good share a reason")
	}
	if wantPrefix := "canonical_bytes: "; len(incomplete.Reason) < len(wantPrefix) || incomplete.Reason[:len(wantPrefix)] != wantPrefix {
		t.Fatalf("incomplete reason must be canonical_bytes-prefixed, got %q", incomplete.Reason)
	}
}

func TestVerifyPQC_Deterministic(t *testing.T) {
	data := load(t, "01_known_good.json")
	a := v2.VerifyPQC(data)
	b := v2.VerifyPQC(data)
	if a != b {
		t.Fatalf("PQC stub is not deterministic: %+v vs %+v", a, b)
	}
}

func TestPQCNeverChangesV1(t *testing.T) {
	cases := []struct {
		file string
		want verifier.Verdict
	}{
		{"01_known_good.json", verifier.PASS},
		{"02_tampered_payload.json", verifier.FAIL},
		{"03_signature_failure.json", verifier.FAIL},
		{"04_hash_mismatch.json", verifier.FAIL},
		{"05_incomplete.json", verifier.INDETERMINATE},
		{"06_unsupported_alg.json", verifier.INDETERMINATE},
	}
	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			data := load(t, tc.file)
			direct := verifier.VerifyJSON(data)
			off := v2.Verify(data, false)
			on := v2.Verify(data, true)

			if direct.Verdict != tc.want {
				t.Fatalf("v1 drifted: want %s got %s (%s)", tc.want, direct.Verdict, direct.Reason)
			}
			if off.Verdict != direct.Verdict || off.Reason != direct.Reason {
				t.Fatalf("disabled envelope changed v1: %+v vs %+v", off, direct)
			}
			if on.Verdict != direct.Verdict || on.Reason != direct.Reason {
				t.Fatalf("enabled envelope changed v1: %+v vs %+v", on, direct)
			}
			if off.PQC != nil {
				t.Fatalf("disabled pqc must be null, got %+v", off.PQC)
			}
			if on.PQC == nil {
				t.Fatal("enabled pqc must be present")
			}
		})
	}
}

func TestGolden_DisabledNull(t *testing.T) {
	got := v2.Verify(load(t, "01_known_good.json"), false)
	var buf bytes.Buffer
	if err := v2.Encode(&buf, got); err != nil {
		t.Fatal(err)
	}
	want := loadV2(t, "pqc_disabled_golden.json")
	if !bytes.Equal(buf.Bytes(), want) {
		t.Fatalf("disabled golden drift\n got: %s\nwant: %s", buf.Bytes(), want)
	}
}

func TestGolden_StubEnabled(t *testing.T) {
	got := v2.Verify(load(t, "01_known_good.json"), true)
	var buf bytes.Buffer
	if err := v2.Encode(&buf, got); err != nil {
		t.Fatal(err)
	}
	want := loadV2(t, "pqc_stub_golden.json")
	if !bytes.Equal(buf.Bytes(), want) {
		t.Fatalf("enabled golden drift\n got: %s\nwant: %s", buf.Bytes(), want)
	}
}

func TestEncode_PQCNullJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := v2.Encode(&buf, v2.Verify(load(t, "01_known_good.json"), false)); err != nil {
		t.Fatal(err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	if string(raw["pqc"]) != "null" {
		t.Fatalf("disabled pqc JSON must be null, got %s", raw["pqc"])
	}
}
