// Package http provides the QuantumGuard HTTP /verify endpoint.
// It calls the exact same deterministic verifier used by the CLI.
package http

import (
	"encoding/json"
	"io"
	"net/http"

	"quantumguard/verifier"
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

func (s *Server) handleVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1 MiB limit
	if err != nil {
		writeResult(w, verifier.Result{
			Verdict: verifier.INDETERMINATE,
			Reason:  "failed to read body",
		})
		return
	}
	defer r.Body.Close()

	result := verifier.VerifyJSON(body)
	writeResult(w, result)
}

func writeResult(w http.ResponseWriter, result verifier.Result) {
	w.Header().Set("Content-Type", "application/json")

	// Map verdict to HTTP status for convenience, but the body is authoritative.
	switch result.Verdict {
	case verifier.PASS:
		w.WriteHeader(http.StatusOK)
	case verifier.FAIL:
		w.WriteHeader(http.StatusUnprocessableEntity) // 422
	default:
		w.WriteHeader(http.StatusBadRequest) // 400 for INDETERMINATE
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(result)
}
