package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunVersionFlag(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := Run([]string{"--version"}, strings.NewReader(""), &out, &errOut, "1.2.3"); code != 0 {
		t.Fatalf("Run() code = %d, want 0; stderr = %q", code, errOut.String())
	}
	if got, want := out.String(), "Cortex 1.2.3\n"; got != want {
		t.Fatalf("version output = %q, want %q", got, want)
	}
}

func TestRunVersionCommand(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := Run([]string{"version"}, strings.NewReader(""), &out, &errOut, "1.2.3"); code != 0 {
		t.Fatalf("Run() code = %d, want 0; stderr = %q", code, errOut.String())
	}
	if got, want := out.String(), "1.2.3\n"; got != want {
		t.Fatalf("version output = %q, want %q", got, want)
	}
}

func TestRunHelp(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := Run([]string{"--help"}, strings.NewReader(""), &out, &errOut, DefaultVersion); code != 0 {
		t.Fatalf("Run() code = %d, want 0; stderr = %q", code, errOut.String())
	}
	for _, want := range []string{"Usage:", "doctor", "version"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("help output does not contain %q: %q", want, out.String())
		}
	}
}
