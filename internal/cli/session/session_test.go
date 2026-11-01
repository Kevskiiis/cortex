package session

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecuteLineQuitIsTheOnlySessionExitCommand(t *testing.T) {
	var out, errOut bytes.Buffer

	if !executeLine("quit", &out, &errOut, "1.2.3") {
		t.Fatal("executeLine(quit) = false, want true")
	}

	if executeLine("exit", &out, &errOut, "1.2.3") {
		t.Fatal("executeLine(exit) = true, want false")
	}
	if !strings.Contains(errOut.String(), "unknown command") {
		t.Fatalf("exit error = %q, want unknown command", errOut.String())
	}
}

func TestExecuteLineUsesSharedCommandRunner(t *testing.T) {
	var out, errOut bytes.Buffer

	if executeLine("version", &out, &errOut, "1.2.3") {
		t.Fatal("executeLine(version) = true, want false")
	}
	if got, want := out.String(), "1.2.3\n"; got != want {
		t.Fatalf("version output = %q, want %q", got, want)
	}
}
