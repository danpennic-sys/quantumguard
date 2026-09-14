// Package http provides the QuantumGuard HTTP /verify endpoint.
// It calls the exact same deterministic verifier used by the CLI.
package http

import (
	"io"
	"net/http"

	"quantumguard/verifier"
	v2 "quantumguard/verifier/v2"
)

// Server is a minimal HTTP server exposing /verify.
type Server struct {
	mux *http.ServeMux
}

// NewServer creates a server with the /verify route registered.
func NewServer() *Server {
	s := &Server{mux: http.NewServeMux()}
	s.mux.HandleFunc("/verify", s.handleVerify)
	s.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return s
}

// Handler returns the http.Handler for use with http.ListenAndServe or tests.
func (s *Server) Handler() http.Handler {
	return s.mux
}

func pqcEnabled(r *http.Request) bool {
	return r.URL.Query().Get("pqc") == "true"
}

func (s *Server) handleVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1 MiB limit
	if err != nil {
		out := v2.Combined{Verdict: verifier.INDETERMINATE, Reason: "failed to read body"}
		if pqcEnabled(r) {
			p := v2.VerifyPQC(nil)
			out.PQC = &p
		}
		writeResult(w, out)
		return
	}
	defer r.Body.Close()

	result := v2.Verify(body, pqcEnabled(r))
	writeResult(w, result)
}

func writeResult(w http.ResponseWriter, result v2.Combined) {
	w.Header().Set("Content-Type", "application/json")

	// HTTP status follows the v1 verdict only. The body is authoritative.
	switch result.Verdict {
	case verifier.PASS:
		w.WriteHeader(http.StatusOK)
	case verifier.FAIL:
		w.WriteHeader(http.StatusUnprocessableEntity) // 422
	default:
		w.WriteHeader(http.StatusBadRequest) // 400 for INDETERMINATE
	}

	_ = v2.Encode(w, result)
}
