package ui

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestRendererVersionJSONIsMachineReadable(t *testing.T) {
	var out bytes.Buffer
	if err := (Renderer{Out: &out, JSON: true}).Version("1.2.3"); err != nil {
		t.Fatalf("Version() error = %v", err)
	}

	var envelope Envelope
	if err := json.Unmarshal(out.Bytes(), &envelope); err != nil {
		t.Fatalf("Version() output is not JSON: %v", err)
	}
	if !envelope.OK {
		t.Fatal("Version() envelope OK = false, want true")
	}
	if !strings.Contains(string(out.Bytes()), `"version":"1.2.3"`) {
		t.Fatalf("Version() output = %q, want version field", out.String())
	}
}

func TestRenderChecksASCIIFallback(t *testing.T) {
	output := RenderChecks([]Check{{Name: "Platform", Status: "Windows", Detail: "Supported", OK: true}}, false, 80)
	if !strings.Contains(output, "[OK]") {
		t.Fatalf("RenderChecks() = %q, want ASCII success marker", output)
	}
	if strings.Contains(output, "✓") {
		t.Fatalf("RenderChecks() = %q, did not expect Unicode marker", output)
	}
}

func TestRenderChecksKeepsNarrowOutputWithinMinimumWidth(t *testing.T) {
	output := RenderChecks([]Check{{Name: "PowerShell", Status: "Missing", Detail: "Install Windows PowerShell 5.1", OK: false}}, false, 20)
	if !strings.Contains(output, "PowerShell") {
		t.Fatalf("RenderChecks() = %q, want check name", output)
	}
	if strings.Contains(output, "\x1b[") {
		t.Fatalf("RenderChecks() = %q, did not expect ANSI escapes", output)
	}
}

func TestRendererDoctorJSONContainsChecksAndError(t *testing.T) {
	var out bytes.Buffer
	renderer := Renderer{Out: &out, JSON: true}
	if err := renderer.Doctor([]Check{{Name: "PowerShell", Status: "Missing", Detail: "Install it", OK: false}}, 1); err != nil {
		t.Fatalf("Doctor() error = %v", err)
	}

	var envelope Envelope
	if err := json.Unmarshal(out.Bytes(), &envelope); err != nil {
		t.Fatalf("Doctor() output is not JSON: %v", err)
	}
	if envelope.OK {
		t.Fatal("Doctor() envelope OK = true, want false")
	}
	if envelope.Error == nil || envelope.Error.Code != "doctor_failed" {
		t.Fatalf("Doctor() error = %#v, want doctor_failed", envelope.Error)
	}
}
