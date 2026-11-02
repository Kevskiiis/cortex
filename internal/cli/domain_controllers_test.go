package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"cortex/internal/features/domain/contracts"
)

type fakeDomainControllersProvider struct {
	report  contracts.DomainControllersReport
	err     error
	request contracts.DomainControllersRequest
}

func (p *fakeDomainControllersProvider) DomainControllers(_ context.Context, request contracts.DomainControllersRequest) (contracts.DomainControllersReport, error) {
	p.request = request
	return p.report, p.err
}

func TestDomainControllersCommandRendersHumanResultAndPassesTarget(t *testing.T) {
	provider := &fakeDomainControllersProvider{report: contracts.DomainControllersReport{
		Status: "ready",
		Controllers: []contracts.DomainController{{
			Name:            "DC01.contoso.com",
			IPv4Address:     "10.0.0.10",
			OperatingSystem: "Windows Server 2022",
			Site:            "HQ",
			Enabled:         true,
			Writable:        true,
			GlobalCatalog:   true,
		}},
	}}
	var out, errOut bytes.Buffer
	root := NewRootCommandWithAllProviders(DefaultVersion, &fakeDomainInfoProvider{}, fakeIdentityProvider{}, provider)
	root.SetArgs([]string{"domain", "controllers", "--server", "DC01.contoso.com"})
	root.SetOut(&out)
	root.SetErr(&errOut)

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if provider.request.Server != "DC01.contoso.com" || provider.request.Domain != "" {
		t.Fatalf("request = %+v, want server target", provider.request)
	}
	for _, want := range []string{"Cortex domain controllers", "DC01.contoso.com", "10.0.0.10", "writable", "Found 1 domain controller(s)."} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output does not contain %q: %q", want, out.String())
		}
	}
}

func TestDomainControllersCommandRendersJSON(t *testing.T) {
	provider := &fakeDomainControllersProvider{report: contracts.DomainControllersReport{
		Status:      "empty",
		Controllers: []contracts.DomainController{},
	}}
	var out, errOut bytes.Buffer
	root := NewRootCommandWithAllProviders(DefaultVersion, &fakeDomainInfoProvider{}, fakeIdentityProvider{}, provider)
	root.SetArgs([]string{"--json", "domain", "controllers"})
	root.SetOut(&out)
	root.SetErr(&errOut)

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	var result map[string]any
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("output is not JSON: %v; output = %q", err, out.String())
	}
	if result["ok"] != true {
		t.Fatalf("JSON ok = %v, want true for empty query", result["ok"])
	}
	if strings.Contains(out.String(), "\x1b[") {
		t.Fatalf("JSON output contains ANSI escape sequence: %q", out.String())
	}
}

func TestDomainControllersCommandRejectsBothTargets(t *testing.T) {
	provider := &fakeDomainControllersProvider{}
	var out, errOut bytes.Buffer
	root := NewRootCommandWithAllProviders(DefaultVersion, &fakeDomainInfoProvider{}, fakeIdentityProvider{}, provider)
	root.SetArgs([]string{"domain", "controllers", "--domain", "contoso.com", "--server", "DC01"})
	root.SetOut(&out)
	root.SetErr(&errOut)

	err := root.Execute()
	if err == nil || err.Error() != "--domain and --server cannot be used together" {
		t.Fatalf("Execute() error = %v, want target validation error", err)
	}
}

func TestDomainControllersCommandReturnsProviderError(t *testing.T) {
	provider := &fakeDomainControllersProvider{err: errors.New("provider failed")}
	var out, errOut bytes.Buffer
	root := NewRootCommandWithAllProviders(DefaultVersion, &fakeDomainInfoProvider{}, fakeIdentityProvider{}, provider)
	root.SetArgs([]string{"domain", "controllers"})
	root.SetOut(&out)
	root.SetErr(&errOut)

	if err := root.Execute(); err == nil || err.Error() != "provider failed" {
		t.Fatalf("Execute() error = %v, want provider failed", err)
	}
}
