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

type fakeIdentityProvider struct {
	report directory.IdentityReport
	err    error
}

func (p fakeIdentityProvider) WhoAmI(context.Context) (directory.IdentityReport, error) {
	return p.report, p.err
}

func TestWhoAmICommandRendersHumanResult(t *testing.T) {
	provider := fakeIdentityProvider{report: directory.IdentityReport{
		Status: "domain_identity",
		Checks: []directory.DiagnosticCheck{{Name: "Identity", Status: "Ready", Detail: "CONTOSO\\jsmith", OK: true}},
		Info: &directory.IdentityInfo{
			Name:           `CONTOSO\jsmith`,
			SID:            "S-1-5-21-123",
			Authentication: "domain",
			ComputerName:   "ADMIN-PC-01",
			Domain:         "contoso.com",
			DomainJoined:   true,
		},
	}}
	var out, errOut bytes.Buffer
	root := NewRootCommandWithProviders(DefaultVersion, &fakeDomainInfoProvider{}, provider)
	root.SetArgs([]string{"whoami"})
	root.SetOut(&out)
	root.SetErr(&errOut)

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"Cortex identity", `CONTOSO\jsmith`, "contoso.com", "Identity information retrieved successfully."} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output does not contain %q: %q", want, out.String())
		}
	}
}

func TestWhoAmICommandRendersJSONResult(t *testing.T) {
	provider := fakeIdentityProvider{report: directory.IdentityReport{
		Status: "local_identity",
		Info:   &directory.IdentityInfo{Name: `WORKSTATION\jsmith`, Classification: "local_identity"},
	}}
	var out, errOut bytes.Buffer
	root := NewRootCommandWithProviders(DefaultVersion, &fakeDomainInfoProvider{}, provider)
	root.SetArgs([]string{"--json", "whoami"})
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
}

func TestWhoAmICommandRendersFailureAndReturnsError(t *testing.T) {
	provider := fakeIdentityProvider{report: directory.IdentityReport{
		Status: "identity_unavailable",
		Error: &directory.DiagnosticError{
			Code:    "identity_unavailable",
			Message: "Windows did not provide a current user identity.",
		},
	}}
	var out, errOut bytes.Buffer
	root := NewRootCommandWithProviders(DefaultVersion, &fakeDomainInfoProvider{}, provider)
	root.SetArgs([]string{"whoami"})
	root.SetOut(&out)
	root.SetErr(&errOut)

	err := root.Execute()
	if err == nil || !strings.Contains(err.Error(), "current user identity") {
		t.Fatalf("Execute() error = %v, want identity failure", err)
	}
	if !strings.Contains(out.String(), "Windows did not provide") {
		t.Fatalf("output = %q, want friendly failure", out.String())
	}
}

func TestWhoAmICommandReturnsProviderError(t *testing.T) {
	provider := fakeIdentityProvider{err: errors.New("provider failed")}
	var out, errOut bytes.Buffer
	root := NewRootCommandWithProviders(DefaultVersion, &fakeDomainInfoProvider{}, provider)
	root.SetArgs([]string{"whoami"})
	root.SetOut(&out)
	root.SetErr(&errOut)

	if err := root.Execute(); err == nil || err.Error() != "provider failed" {
		t.Fatalf("Execute() error = %v, want provider failed", err)
	}
}
