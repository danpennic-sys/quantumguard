// Command qg is the QuantumGuard CLI.
// It calls the exact same deterministic verifier used by the HTTP endpoint.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"quantumguard/verifier"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: qg verify <bundle.json>\n")
		os.Exit(2)
	}

	cmd := os.Args[1]
	switch cmd {
	case "verify":
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "usage: qg verify <bundle.json>\n")
			os.Exit(2)
		}
		runVerify(os.Args[2])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		os.Exit(2)
	}
}

func runVerify(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read error: %v\n", err)
		os.Exit(2)
	}

	result := verifier.VerifyJSON(data)

	// Machine-readable output
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
		os.Exit(2)
	}

	switch result.Verdict {
	case verifier.PASS:
		os.Exit(0)
	case verifier.FAIL:
		os.Exit(1)
	default: // INDETERMINATE
		os.Exit(3)
	}
}
