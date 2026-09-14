package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseArgs(t *testing.T) {
	cases := []struct {
		args    []string
		path    string
		pqc     bool
		wantErr bool
	}{
		{[]string{"verify", "bundle.json"}, "bundle.json", false, false},
		{[]string{"verify", "--pqc", "bundle.json"}, "bundle.json", true, false},
		{[]string{"verify", "bundle.json", "--pqc"}, "bundle.json", true, false},
		{[]string{"verify"}, "", false, true},
		{[]string{"verify", "--pqc"}, "", false, true},
		{[]string{"other", "bundle.json"}, "", false, true},
		{[]string{"verify", "--nope", "bundle.json"}, "", false, true},
		{[]string{"verify", "a.json", "b.json"}, "", false, true},
	}
	for _, tc := range cases {
		path, pqc, err := parseArgs(tc.args)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("args %v: expected error", tc.args)
			}
			continue
		}
		if err != nil {
			t.Fatalf("args %v: unexpected error %v", tc.args, err)
		}
		if path != tc.path || pqc != tc.pqc {
			t.Fatalf("args %v: path=%q pqc=%v want path=%q pqc=%v", tc.args, path, pqc, tc.path, tc.pqc)
		}
	}
}

func TestRunVerify_ExitFollowsV1(t *testing.T) {
	good := filepath.Join("..", "..", "vectors", "01_known_good.json")
	if _, err := os.Stat(good); err != nil {
		good = filepath.Join("vectors", "01_known_good.json")
	}

	// Redirect stdout so the test doesn't spam.
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }()
	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer devnull.Close()

	os.Stdout = devnull
	if code := runVerify(good, false); code != 0 {
		t.Fatalf("known-good without --pqc: exit %d", code)
	}
	if code := runVerify(good, true); code != 0 {
		t.Fatalf("known-good with --pqc must not change v1 exit, got %d", code)
	}
}
