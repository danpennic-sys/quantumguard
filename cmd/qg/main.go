// Command qg is the QuantumGuard CLI.
// It calls the exact same deterministic verifier used by the HTTP endpoint.
package main

import (
	"fmt"
	"os"
	"strings"

	"quantumguard/verifier"
	v2 "quantumguard/verifier/v2"
)

func main() {
	path, pqc, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	os.Exit(runVerify(path, pqc))
}

func parseArgs(args []string) (path string, pqc bool, err error) {
	if len(args) == 0 {
		return "", false, fmt.Errorf("usage: qg verify [--pqc] <bundle.json>")
	}
	if args[0] != "verify" {
		return "", false, fmt.Errorf("unknown command: %s", args[0])
	}

	var files []string
	for _, a := range args[1:] {
		switch a {
		case "--pqc":
			pqc = true
		default:
			if strings.HasPrefix(a, "-") {
				return "", false, fmt.Errorf("unknown flag: %s\nusage: qg verify [--pqc] <bundle.json>", a)
			}
			files = append(files, a)
		}
	}
	if len(files) != 1 {
		return "", false, fmt.Errorf("usage: qg verify [--pqc] <bundle.json>")
	}
	return files[0], pqc, nil
}

func runVerify(path string, pqc bool) int {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read error: %v\n", err)
		return 2
	}

	result := v2.Verify(data, pqc)
	if err := v2.Encode(os.Stdout, result); err != nil {
		fmt.Fprintf(os.Stderr, "encode error: %v\n", err)
		return 2
	}

	// Exit code follows the v1 verdict only. PQC never changes it.
	switch result.Verdict {
	case verifier.PASS:
		return 0
	case verifier.FAIL:
		return 1
	default:
		return 3
	}
}
