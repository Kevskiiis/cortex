package contracts

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

type DomainControllersRequest struct {
	Domain string
	Server string
}

type DomainControllersProvider interface {
	DomainControllers(context.Context, DomainControllersRequest) (DomainControllersReport, error)
}

type DomainControllersReport struct {
	Status      string             `json:"status"`
	Controllers []DomainController `json:"controllers"`
	Checks      []DiagnosticCheck  `json:"checks"`
	Error       *DiagnosticError   `json:"error,omitempty"`
}

type DomainController struct {
	Name            string `json:"name"`
	IPv4Address     string `json:"ipv4_address"`
	OperatingSystem string `json:"operating_system"`
	Site            string `json:"site"`
	Enabled         bool   `json:"enabled"`
	Writable        bool   `json:"writable"`
	GlobalCatalog   bool   `json:"global_catalog"`
}

func (r DomainControllersReport) Successful() bool {
	return r.Status == "ready" || r.Status == "empty"
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
