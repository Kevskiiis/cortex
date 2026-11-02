package powershell

import (
	"context"
	"errors"
	"testing"

	ps "cortex/internal/powershell"
)

func TestProcessFailureClassifiesTimeout(t *testing.T) {
	report := processFailure(&ps.ProcessError{Err: context.DeadlineExceeded})
	if report.Status != "command_timeout" {
		t.Fatalf("status = %q, want command_timeout", report.Status)
	}
}

func TestProcessFailureClassifiesPowerShellUnavailable(t *testing.T) {
	report := processFailure(&ps.ProcessError{Err: errors.New("cannot find the file")})
	if report.Status != "powershell_unavailable" {
		t.Fatalf("status = %q, want powershell_unavailable", report.Status)
	}
}
