package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"cortex/internal/directory"
)

type fakeDomainInfoProvider struct {
	report  directory.DomainInfoReport
	err     error
	request directory.DomainInfoRequest
}

func (p *fakeDomainInfoProvider) DomainInfo(_ context.Context, request directory.DomainInfoRequest) (directory.DomainInfoReport, error) {
	p.request = request
	return p.report, p.err
}

func TestDomainInfoCommandRendersHumanResult(t *testing.T) {
	provider := &fakeDomainInfoProvider{report: directory.DomainInfoReport{
		Status: "ready",
		Checks: []directory.DiagnosticCheck{{Name: "Directory", Status: "Ready", Detail: "Query succeeded", OK: true}},
		Info:   &directory.DomainInfo{DNSName: "corp.example.com", CurrentUser: "CORP\\admin"},
	}}
	var out, errOut bytes.Buffer
	root := NewRootCommandWithProvider(DefaultVersion, provider)
	root.SetArgs([]string{"domain", "info", "--domain", "corp.example.com"})
	root.SetOut(&out)
	root.SetErr(&errOut)

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if provider.request.Domain != "corp.example.com" {
		t.Fatalf("provider domain = %q, want corp.example.com", provider.request.Domain)
	}
	for _, want := range []string{"Cortex domain info", "corp.example.com", "Domain information retrieved successfully."} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output does not contain %q: %q", want, out.String())
		}
	}
}

func TestDomainInfoCommandRendersJSONResult(t *testing.T) {
	provider := &fakeDomainInfoProvider{report: directory.DomainInfoReport{Status: "partial"}}
	var out, errOut bytes.Buffer
	root := NewRootCommandWithProvider(DefaultVersion, provider)
	root.SetArgs([]string{"--json", "domain", "info", "--server", "dc01.corp.example.com"})
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
		t.Fatalf("JSON ok = %v, want true", result["ok"])
	}
	if provider.request.Server != "dc01.corp.example.com" {
		t.Fatalf("provider server = %q, want dc01.corp.example.com", provider.request.Server)
	}
}

func TestDomainInfoCommandRejectsMultipleTargets(t *testing.T) {
	provider := &fakeDomainInfoProvider{}
	var out, errOut bytes.Buffer
	root := NewRootCommandWithProvider(DefaultVersion, provider)
	root.SetArgs([]string{"domain", "info", "--domain", "corp.example.com", "--server", "dc01.corp.example.com"})
	root.SetOut(&out)
	root.SetErr(&errOut)

	err := root.Execute()
	if err == nil || !strings.Contains(err.Error(), "cannot be used together") {
		t.Fatalf("Execute() error = %v, want mutually exclusive target error", err)
	}
}

func TestDomainInfoCommandReturnsProviderError(t *testing.T) {
	provider := &fakeDomainInfoProvider{err: errors.New("provider failed")}
	var out, errOut bytes.Buffer
	root := NewRootCommandWithProvider(DefaultVersion, provider)
	root.SetArgs([]string{"domain", "info"})
	root.SetOut(&out)
	root.SetErr(&errOut)

	if err := root.Execute(); err == nil || err.Error() != "provider failed" {
		t.Fatalf("Execute() error = %v, want provider failed", err)
	}
}
