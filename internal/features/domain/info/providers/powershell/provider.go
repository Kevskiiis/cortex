package powershell

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"cortex/internal/features/domain/contracts"
	psexec "cortex/internal/powershell"
)

//go:embed domain_info.ps1
var domainInfoScript []byte

type Provider struct {
	Runner psexec.Runner
}

func NewProvider() Provider {
	return Provider{Runner: psexec.NewRunner()}
}

type scriptResult struct {
	Status string                      `json:"status"`
	Info   *contracts.DomainInfo       `json:"info,omitempty"`
	Checks []contracts.DiagnosticCheck `json:"checks"`
	Error  *contracts.DiagnosticError  `json:"error,omitempty"`
}

func (p Provider) DomainInfo(ctx context.Context, request contracts.DomainInfoRequest) (contracts.DomainInfoReport, error) {
	if request.Domain != "" && request.Server != "" {
		return contracts.DomainInfoReport{}, errors.New("domain and server targets are mutually exclusive")
	}
	args := make([]string, 0, 4)
	if request.Domain != "" {
		args = append(args, "-Domain", request.Domain)
	}
	if request.Server != "" {
		args = append(args, "-Server", request.Server)
	}
	output, err := p.Runner.Run(ctx, domainInfoScript, args)
	if err != nil {
		return processFailure(err), nil
	}
	var result scriptResult
	if err := json.Unmarshal(output, &result); err != nil {
		return contracts.DomainInfoReport{Status: "invalid_result", Error: &contracts.DiagnosticError{
			Code: "invalid_result", Message: "PowerShell returned an invalid diagnostic result.", Detail: err.Error(), Next: "Run the command again and use --json for diagnostic details.",
		}}, nil
	}
	return contracts.DomainInfoReport{Status: result.Status, Info: result.Info, Checks: result.Checks, Error: result.Error}, nil
}

func processFailure(err error) contracts.DomainInfoReport {
	status, code, next, detail := "unexpected_error", "unexpected_error", "Run the command again.", ""
	var processErr *psexec.ProcessError
	if errors.As(err, &processErr) {
		detail = processErr.Error()
		lower := strings.ToLower(detail)
		switch {
		case errors.Is(processErr, context.DeadlineExceeded):
			status, code, next = "command_timeout", "command_timeout", "Check the network and try again."
		case errors.Is(processErr, context.Canceled):
			status, code, next = "cancelled", "cancelled", "Run the command again when ready."
		case strings.Contains(lower, "not recognized"), strings.Contains(lower, "cannot find the file"):
			status, code, next = "powershell_unavailable", "powershell_unavailable", "Install Windows PowerShell 5.1 and try again."
		case strings.Contains(lower, "access is denied"), strings.Contains(lower, "unauthorized"):
			status, code, next = "access_denied", "access_denied", "Use a domain-authenticated Windows account with directory read access."
		case strings.Contains(lower, "trust relationship"):
			status, code, next = "secure_channel_failed", "secure_channel_failed", "Repair the computer trust relationship before trying again."
		case strings.Contains(lower, "server is not operational"), strings.Contains(lower, "cannot contact"), strings.Contains(lower, "network"):
			status, code, next = "domain_unreachable", "domain_unreachable", "Connect to the domain network or VPN and try again."
		}
	}
	return contracts.DomainInfoReport{Status: status, Error: &contracts.DiagnosticError{Code: code, Message: fmt.Sprintf("Cortex could not complete the domain diagnostic: %s.", status), Detail: detail, Next: next}}
}
