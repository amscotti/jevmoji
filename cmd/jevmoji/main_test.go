package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	body := "# comment line\n\nTYPESAFE_API_KEY=abc123\nQUOTED=\"yes\"\nSINGLE='it works'\nEXPORTED=export-value\nBROKEN LINE\nEMPTY=\n"
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write env file: %v", err)
	}

	t.Setenv("PRESERVED", "from-environment")
	t.Setenv("TYPESAFE_API_KEY", "from-environment")

	if err := loadEnvFile(path); err != nil {
		t.Fatalf("loadEnvFile: %v", err)
	}

	if got := os.Getenv("TYPESAFE_API_KEY"); got != "from-environment" {
		t.Errorf("existing variable must win, got %q", got)
	}
	if got := os.Getenv("QUOTED"); got != "yes" {
		t.Errorf("QUOTED = %q, want %q", got, "yes")
	}
	if got := os.Getenv("SINGLE"); got != "it works" {
		t.Errorf("SINGLE = %q, want %q", got, "it works")
	}
	if got := os.Getenv("EXPORTED"); got != "export-value" {
		t.Errorf("EXPORTED = %q, want %q", got, "export-value")
	}
	if _, exists := os.LookupEnv("EMPTY"); !exists {
		t.Error("EMPTY should be set to the empty string")
	}
	if _, exists := os.LookupEnv("BROKEN LINE"); exists {
		t.Error("lines without '=' must be ignored")
	}
}

func TestLoadEnvFileMissing(t *testing.T) {
	if err := loadEnvFile(filepath.Join(t.TempDir(), "nope.env")); !os.IsNotExist(err) {
		t.Errorf("want ErrNotExist, got %v", err)
	}
}

func TestLoadEnvFileExportAndQuotes(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	body := "export JEVMOJI_EXPORTED=hi\nJEVMOJI_MISMATCHED=\"oops\nJEVMOJI_INNER=it's fine\nJEVMOJI_HASHED=value # kept verbatim\n"
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write env file: %v", err)
	}
	if err := loadEnvFile(path); err != nil {
		t.Fatalf("loadEnvFile: %v", err)
	}
	if got := os.Getenv("JEVMOJI_EXPORTED"); got != "hi" {
		t.Errorf("JEVMOJI_EXPORTED = %q, want %q", got, "hi")
	}
	if got := os.Getenv("JEVMOJI_MISMATCHED"); got != `"oops` {
		t.Errorf("JEVMOJI_MISMATCHED = %q, want it untouched", got)
	}
	if got := os.Getenv("JEVMOJI_INNER"); got != "it's fine" {
		t.Errorf("JEVMOJI_INNER = %q", got)
	}
	if got := os.Getenv("JEVMOJI_HASHED"); got != "value # kept verbatim" {
		t.Errorf("JEVMOJI_HASHED = %q (no inline-comment stripping by design)", got)
	}
}

func TestUnquote(t *testing.T) {
	cases := map[string]string{
		`"yes"`:      "yes",
		`'yes'`:      "yes",
		`"oops`:      `"oops`,
		`oops"`:      `oops"`,
		`""`:         "",
		`"it's"`:     "it's",
		`it's fine`:  "it's fine",
		`plain`:      "plain",
		``:           "",
		`"a'b"`:      "a'b",
		`'a"b'`:      `a"b`,
		`"'nested'"`: `'nested'`,
	}
	for in, want := range cases {
		if got := unquote(in); got != want {
			t.Errorf("unquote(%q) = %q, want %q", in, got, want)
		}
	}
}
