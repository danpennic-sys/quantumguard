package verifier_test

import (
	"encoding/json"
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

func TestV1ResultJSONUnchanged(t *testing.T) {
	res := verifier.VerifyJSON(load(t, "01_known_good.json"))
	raw, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["pqc"]; ok {
		t.Fatalf("v1 Result JSON must not include pqc, got %s", raw)
	}
	if string(m["verdict"]) != `"PASS"` {
		t.Fatalf("v1 verdict drifted: %s", raw)
	}
}

func TestUppercaseHex_INDETERMINATE(t *testing.T) {
	// Based on 01_known_good.json with payload_hash uppercased
	data := []byte(`{
  "alg": "ed25519",
  "payload": "UXVhbnR1bUd1YXJkIGV2aWRlbmNlIHBheWxvYWQgdjEgLSBrbm93biBnb29k",
  "payload_hash": "71B8737C3C1128BB466C53E821F5CF1A1E765B82041A60C46672226B73A16DFC",
  "public_key": "79b5562e8fe654f94078b112e8a98ba7901f853ae695bed7e0e3910bad049664",
  "signature": "a3fd916f92143149a71f02eb7d7220b432d28f4d160a1c53978966e362c0c273572627bed2915cbfefe0a05cecde42b9869a85d782f773a412c7d2759e169c00",
  "version": "1"
}`)
	res := verifier.VerifyJSON(data)
	if res.Verdict != verifier.INDETERMINATE || res.Code != verifier.CodeMalformedEncoding {
		t.Fatalf("uppercase hex: want INDETERMINATE/E_MALFORMED_ENCODING, got %s/%s", res.Verdict, res.Code)
	}
}

func TestUnknownField_INDETERMINATE(t *testing.T) {
	data := []byte(`{
  "alg": "ed25519",
  "payload": "UXVhbnR1bUd1YXJkIGV2aWRlbmNlIHBheWxvYWQgdjEgLSBrbm93biBnb29k",
  "payload_hash": "71b8737c3c1128bb466c53e821f5cf1a1e765b82041a60c46672226b73a16dfc",
  "public_key": "79b5562e8fe654f94078b112e8a98ba7901f853ae695bed7e0e3910bad049664",
  "signature": "a3fd916f92143149a71f02eb7d7220b432d28f4d160a1c53978966e362c0c273572627bed2915cbfefe0a05cecde42b9869a85d782f773a412c7d2759e169c00",
  "version": "1",
  "extra": "not allowed"
}`)
	res := verifier.VerifyJSON(data)
	if res.Verdict != verifier.INDETERMINATE || res.Code != verifier.CodeMalformedJSON {
		t.Fatalf("unknown field: want INDETERMINATE/E_MALFORMED_JSON, got %s/%s", res.Verdict, res.Code)
	}
}

func TestTrailingData_INDETERMINATE(t *testing.T) {
	data := []byte(`{"alg":"ed25519","payload":"UXVhbnR1bUd1YXJkIGV2aWRlbmNlIHBheWxvYWQgdjEgLSBrbm93biBnb29k","payload_hash":"71b8737c3c1128bb466c53e821f5cf1a1e765b82041a60c46672226b73a16dfc","public_key":"79b5562e8fe654f94078b112e8a98ba7901f853ae695bed7e0e3910bad049664","signature":"a3fd916f92143149a71f02eb7d7220b432d28f4d160a1c53978966e362c0c273572627bed2915cbfefe0a05cecde42b9869a85d782f773a412c7d2759e169c00","version":"1"}{"extra":true}`)
	res := verifier.VerifyJSON(data)
	if res.Verdict != verifier.INDETERMINATE || res.Code != verifier.CodeMalformedJSON {
		t.Fatalf("trailing data: want INDETERMINATE/E_MALFORMED_JSON, got %s/%s", res.Verdict, res.Code)
	}
}

func TestOversizedJSON_INDETERMINATE(t *testing.T) {
	// > 1 MiB raw JSON
	big := make([]byte, 1<<20+1)
	for i := range big {
		big[i] = 'a'
	}
	data := append([]byte(`{"alg":"ed25519","payload":"`), append(big, []byte(`","payload_hash":"71b8737c3c1128bb466c53e821f5cf1a1e765b82041a60c46672226b73a16dfc","public_key":"79b5562e8fe654f94078b112e8a98ba7901f853ae695bed7e0e3910bad049664","signature":"a3fd916f92143149a71f02eb7d7220b432d28f4d160a1c53978966e362c0c273572627bed2915cbfefe0a05cecde42b9869a85d782f773a412c7d2759e169c00","version":"1"}`)...)...)
	res := verifier.VerifyJSON(data)
	if res.Verdict != verifier.INDETERMINATE || res.Code != verifier.CodeMalformedEncoding {
		t.Fatalf("oversized: want INDETERMINATE/E_MALFORMED_ENCODING, got %s/%s", res.Verdict, res.Code)
	}
}
