package powershell

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"cortex/internal/directory"
	psexec "cortex/internal/powershell"
)

//go:embed domain_controllers.ps1
var domainControllersScript []byte

type domainControllersScriptResult struct {
	Status      string                       `json:"status"`
	Controllers []directory.DomainController `json:"controllers"`
	Checks      []directory.DiagnosticCheck  `json:"checks"`
	Error       *directory.DiagnosticError   `json:"error,omitempty"`
}

func (p Provider) DomainControllers(ctx context.Context, request directory.DomainControllersRequest) (directory.DomainControllersReport, error) {
	if request.Domain != "" && request.Server != "" {
		return directory.DomainControllersReport{}, errors.New("domain and server targets are mutually exclusive")
	}

	args := make([]string, 0, 4)
	if request.Domain != "" {
		args = append(args, "-Domain", request.Domain)
	}
	if request.Server != "" {
		args = append(args, "-Server", request.Server)
	}

	output, err := p.Runner.Run(ctx, domainControllersScript, args)
	if err != nil {
		return controllersProcessFailure(err), nil
	}

	var result domainControllersScriptResult
	if err := json.Unmarshal(output, &result); err != nil {
		return directory.DomainControllersReport{
			Status: "invalid_result",
			Error: &directory.DiagnosticError{
				Code:    "invalid_result",
				Message: "PowerShell returned an invalid domain-controller result.",
				Detail:  err.Error(),
				Next:    "Run the command again and use --json for diagnostic details.",
			},
		}, nil
	}
	return directory.DomainControllersReport{
		Status:      result.Status,
		Controllers: result.Controllers,
		Checks:      result.Checks,
		Error:       result.Error,
	}, nil
}

func controllersProcessFailure(err error) directory.DomainControllersReport {
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
	return directory.DomainControllersReport{
		Status: status,
		Error: &directory.DiagnosticError{
			Code:    code,
			Message: fmt.Sprintf("Cortex could not list domain controllers: %s.", status),
			Detail:  detail,
			Next:    next,
		},
	}
}
