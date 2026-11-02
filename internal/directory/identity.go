package directory

import "context"

type IdentityProvider interface {
	WhoAmI(context.Context) (IdentityReport, error)
}

type IdentityReport struct {
	Status string            `json:"status"`
	Info   *IdentityInfo     `json:"info,omitempty"`
	Checks []DiagnosticCheck `json:"checks"`
	Error  *DiagnosticError  `json:"error,omitempty"`
}

type IdentityInfo struct {
	Name           string `json:"name"`
	SID            string `json:"sid"`
	Authentication string `json:"authentication"`
	ComputerName   string `json:"computer_name"`
	Domain         string `json:"domain"`
	DomainJoined   bool   `json:"domain_joined"`
	Classification string `json:"classification"`
}

func (r IdentityReport) Successful() bool {
	return r.Status == "ready" || r.Status == "local_identity" || r.Status == "domain_identity"
}
