package directory

import "context"

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
