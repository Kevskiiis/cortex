package directory

import "context"

type DomainInfoRequest struct {
	Domain string
	Server string
}

type DomainInfoProvider interface {
	DomainInfo(context.Context, DomainInfoRequest) (DomainInfoReport, error)
}

type DomainInfoReport struct {
	Status   string            `json:"status"`
	Info     *DomainInfo       `json:"info,omitempty"`
	Checks   []DiagnosticCheck `json:"checks"`
	Warnings []string          `json:"warnings,omitempty"`
	Error    *DiagnosticError  `json:"error,omitempty"`
}

type DiagnosticCheck struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail"`
	OK     bool   `json:"ok"`
}

type DiagnosticError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
	Next    string `json:"next,omitempty"`
}

type DomainInfo struct {
	DNSName           string   `json:"dns_name"`
	NetBIOSName       string   `json:"netbios_name"`
	DistinguishedName string   `json:"distinguished_name"`
	ForestName        string   `json:"forest_name"`
	CurrentUser       string   `json:"current_user"`
	ComputerName      string   `json:"computer_name"`
	ComputerDomain    string   `json:"computer_domain"`
	DomainControllers []string `json:"domain_controllers"`
	Target            string   `json:"target,omitempty"`
}

func (r DomainInfoReport) Successful() bool {
	return r.Status == "ready" || r.Status == "partial"
}
