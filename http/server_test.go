package http_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	qghttp "quantumguard/http"
	"quantumguard/verifier"
	v2 "quantumguard/verifier/v2"
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

func postVerify(t *testing.T, path string, body []byte) (v2.Combined, int) {
	t.Helper()
	srv := qghttp.NewServer()
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	rr := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rr, req)

	var res v2.Combined
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return res, rr.Code
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

			direct := verifier.VerifyJSON(data)
			viaHTTP, _ := postVerify(t, "/verify", data)
			viaPQC, codePQC := postVerify(t, "/verify?pqc=true", data)
			viaOff, codeOff := postVerify(t, "/verify", data)

			if direct.Verdict != tc.want {
				t.Fatalf("direct verifier: want %s, got %s (%s)", tc.want, direct.Verdict, direct.Reason)
			}
			if viaHTTP.Verdict != tc.want {
				t.Fatalf("HTTP verifier: want %s, got %s (%s)", tc.want, viaHTTP.Verdict, viaHTTP.Reason)
			}
			if direct.Verdict != viaHTTP.Verdict {
				t.Fatalf("CLI/HTTP mismatch: direct=%s http=%s", direct.Verdict, viaHTTP.Verdict)
			}
			if viaHTTP.PQC != nil {
				t.Fatalf("default /verify must set pqc=null, got %+v", viaHTTP.PQC)
			}
			if viaPQC.Verdict != direct.Verdict || viaPQC.Reason != direct.Reason {
				t.Fatalf("?pqc=true changed v1: %+v vs %+v", viaPQC, direct)
			}
			if viaPQC.PQC == nil {
				t.Fatal("?pqc=true must include a pqc object")
			}
			if codePQC != codeOff {
				t.Fatalf("?pqc=true must not change HTTP status: pqc=%d off=%d", codePQC, codeOff)
			}
			if viaOff.PQC != nil {
				t.Fatal("disabled path pqc must be null")
			}
		})
	}
}

func TestHTTP_PQCQueryStrict(t *testing.T) {
	data := load(t, "01_known_good.json")
	for _, path := range []string{"/verify?pqc=1", "/verify?pqc=TRUE", "/verify?pqc=yes", "/verify?pqc=false"} {
		res, _ := postVerify(t, path, data)
		if res.PQC != nil {
			t.Fatalf("%s must not enable pqc, got %+v", path, res.PQC)
		}
	}
	res, _ := postVerify(t, "/verify?pqc=true", data)
	if res.PQC == nil {
		t.Fatal("?pqc=true must enable pqc")
	}
}

func TestHTTP_CLIEnvelopeEquivalence(t *testing.T) {
	data := load(t, "01_known_good.json")
	for _, enable := range []bool{false, true} {
		path := "/verify"
		if enable {
			path = "/verify?pqc=true"
		}
		viaHTTP, _ := postVerify(t, path, data)
		direct := v2.Verify(data, enable)
		if viaHTTP.Verdict != direct.Verdict || viaHTTP.Reason != direct.Reason {
			t.Fatalf("enable=%v envelope mismatch", enable)
		}
		if (viaHTTP.PQC == nil) != (direct.PQC == nil) {
			t.Fatalf("enable=%v pqc null mismatch", enable)
		}
		if viaHTTP.PQC != nil && *viaHTTP.PQC != *direct.PQC {
			t.Fatalf("enable=%v pqc body mismatch: %+v vs %+v", enable, viaHTTP.PQC, direct.PQC)
		}
	}
}

func TestHTTP_BitForBitEncode(t *testing.T) {
	data := load(t, "01_known_good.json")
	cases := []struct {
		path   string
		enable bool
	}{
		{"/verify", false},
		{"/verify?pqc=true", true},
	}
	for _, tc := range cases {
		srv := qghttp.NewServer()
		req := httptest.NewRequest(http.MethodPost, tc.path, bytes.NewReader(data))
		rr := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rr, req)

		var want bytes.Buffer
		if err := v2.Encode(&want, v2.Verify(data, tc.enable)); err != nil {
			t.Fatal(err)
		}
		got, _ := io.ReadAll(rr.Body)
		if !bytes.Equal(got, want.Bytes()) {
			t.Fatalf("%s body drifted from Encode\n got: %s\nwant: %s", tc.path, got, want.Bytes())
		}
	}
}
