package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	qghttp "quantumguard/http"
	"quantumguard/verifier"
)

func load(t *testing.T, name string) []byte {
	t.Helper()
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

func postVerify(t *testing.T, body []byte) verifier.Result {
	t.Helper()
	srv := qghttp.NewServer()
	req := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	var res verifier.Result
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return res
}

func TestHTTP_MatchesVerifier(t *testing.T) {
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

			// Direct verifier
			direct := verifier.VerifyJSON(data)

			// HTTP path
			viaHTTP := postVerify(t, data)

			if direct.Verdict != tc.want {
				t.Fatalf("direct verifier: want %s, got %s (%s)", tc.want, direct.Verdict, direct.Reason)
			}
			if viaHTTP.Verdict != tc.want {
				t.Fatalf("HTTP verifier: want %s, got %s (%s)", tc.want, viaHTTP.Verdict, viaHTTP.Reason)
			}
			if direct.Verdict != viaHTTP.Verdict {
				t.Fatalf("CLI/HTTP mismatch: direct=%s http=%s", direct.Verdict, viaHTTP.Verdict)
			}
		})
	}
}
