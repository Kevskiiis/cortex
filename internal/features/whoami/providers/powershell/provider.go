package powershell

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"strings"

	"cortex/internal/features/whoami/contracts"
	psexec "cortex/internal/powershell"
)

//go:embed identity.ps1
var identityScript []byte

type Provider struct {
	Runner psexec.Runner
}

func NewProvider() Provider {
	return Provider{Runner: psexec.NewRunner()}
}

type scriptResult struct {
	Status string                      `json:"status"`
	Info   *contracts.IdentityInfo     `json:"info,omitempty"`
	Checks []contracts.DiagnosticCheck `json:"checks"`
	Error  *contracts.DiagnosticError  `json:"error,omitempty"`
}

func (p Provider) WhoAmI(ctx context.Context) (contracts.IdentityReport, error) {
	output, err := p.Runner.Run(ctx, identityScript, nil)
	if err != nil {
		return processFailure(err), nil
	}

	var result scriptResult
	if err := json.Unmarshal(output, &result); err != nil {
		return contracts.IdentityReport{
			Status: "invalid_result",
			Error: &contracts.DiagnosticError{
				Code:    "invalid_result",
				Message: "PowerShell returned an invalid identity result.",
				Detail:  err.Error(),
				Next:    "Run the command again and use --json for diagnostic details.",
			},
		}, nil
	}
	return contracts.IdentityReport{Status: result.Status, Info: result.Info, Checks: result.Checks, Error: result.Error}, nil
}

func processFailure(err error) contracts.IdentityReport {
	status := "unexpected_error"
	code := "unexpected_error"
	next := "Run the command again."
	detail := ""
	var processErr *psexec.ProcessError
	if errors.As(err, &processErr) {
		detail = processErr.Error()
		lower := strings.ToLower(detail)
		switch {
		case errors.Is(processErr, context.DeadlineExceeded):
			status, code, next = "command_timeout", "command_timeout", "Run the command again when ready."
		case errors.Is(processErr, context.Canceled):
			status, code, next = "cancelled", "cancelled", "Run the command again when ready."
		case strings.Contains(lower, "not recognized"), strings.Contains(lower, "cannot find the file"):
			status, code, next = "powershell_unavailable", "powershell_unavailable", "Install Windows PowerShell 5.1 and try again."
		}
	}
	return contracts.IdentityReport{
		Status: status,
		Error: &contracts.DiagnosticError{
			Code:    code,
			Message: "Cortex could not determine the current Windows identity.",
			Detail:  detail,
			Next:    next,
		},
	}
}
