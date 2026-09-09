// Command qg-http runs the QuantumGuard HTTP verifier.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	qghttp "quantumguard/http"
)

func main() {
	addr := ":8080"
	if v := os.Getenv("QG_HTTP_ADDR"); v != "" {
		addr = v
	}

	srv := qghttp.NewServer()
	fmt.Printf("QuantumGuard /verify listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, srv.Handler()))
}
