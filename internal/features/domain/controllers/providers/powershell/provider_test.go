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

func TestProcessFailureClassifiesDomainUnavailable(t *testing.T) {
	report := processFailure(&ps.ProcessError{Err: errors.New("server is not operational")})
	if report.Status != "domain_unreachable" {
		t.Fatalf("status = %q, want domain_unreachable", report.Status)
	}
}

func TestProcessFailureClassifiesAccessDenied(t *testing.T) {
	report := processFailure(&ps.ProcessError{Err: errors.New("access is denied")})
	if report.Status != "access_denied" {
		t.Fatalf("status = %q, want access_denied", report.Status)
	}
}
