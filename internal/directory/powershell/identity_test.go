package powershell

import (
	"context"
	"errors"
	"testing"

	ps "cortex/internal/powershell"
)

func TestIdentityProcessFailureClassifiesTimeout(t *testing.T) {
	report := identityProcessFailure(&ps.ProcessError{Err: context.DeadlineExceeded})
	if report.Status != "command_timeout" {
		t.Fatalf("status = %q, want command_timeout", report.Status)
	}
}

func TestIdentityProcessFailureClassifiesPowerShellUnavailable(t *testing.T) {
	report := identityProcessFailure(&ps.ProcessError{Err: errors.New("cannot find the file")})
	if report.Status != "powershell_unavailable" {
		t.Fatalf("status = %q, want powershell_unavailable", report.Status)
	}
}
