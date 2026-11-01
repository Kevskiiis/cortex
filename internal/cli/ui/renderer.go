package ui

import (
	"encoding/json"
	"fmt"
	"io"

	"cortex/internal/directory"
)

type Renderer struct {
	Out   io.Writer
	JSON  bool
	Color bool
}

type Envelope struct {
	OK    bool       `json:"ok"`
	Data  any        `json:"data,omitempty"`
	Error *ErrorInfo `json:"error,omitempty"`
	Meta  any        `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
	Next    string `json:"next,omitempty"`
}

func NewRenderer(out io.Writer, jsonOutput bool) Renderer {
	return Renderer{Out: out, JSON: jsonOutput, Color: ColorEnabled(out)}
}

func (r Renderer) Version(version string) error {
	if r.JSON {
		return r.writeEnvelope(Envelope{OK: true, Data: map[string]string{"version": version}})
	}
	_, err := fmt.Fprintln(r.Out, version)
	return err
}

func (r Renderer) VersionFlag(version string) error {
	if r.JSON {
		return r.writeEnvelope(Envelope{OK: true, Data: map[string]string{"version": version}})
	}
	_, err := fmt.Fprintln(r.Out, "Cortex "+version)
	return err
}

func (r Renderer) Error(err error) error {
	if r.JSON {
		return r.writeEnvelope(Envelope{
			OK: false,
			Error: &ErrorInfo{
				Code:    "command_failed",
				Message: err.Error(),
			},
		})
	}
	_, writeErr := fmt.Fprintln(r.Out, StyleText(ErrorStyle, r.Color, "Error: "+err.Error()))
	return writeErr
}

func (r Renderer) Doctor(checks []Check, failures int) error {
	result := DoctorResult{Checks: checks, Failures: failures}
	if r.JSON {
		var info *ErrorInfo
		if failures > 0 {
			info = &ErrorInfo{
				Code:    "doctor_failed",
				Message: fmt.Sprintf("%d prerequisite check(s) failed", failures),
				Next:    "Resolve the failed checks and run doctor again.",
			}
		}
		return r.writeEnvelope(Envelope{OK: failures == 0, Data: result, Error: info})
	}

	fmt.Fprintln(r.Out, StyleText(TitleStyle, r.Color, "Cortex doctor"))
	fmt.Fprintln(r.Out, RenderChecks(checks, r.Color, TerminalWidth(r.Out)))
	if failures > 0 {
		fmt.Fprintln(r.Out, StyleText(WarningStyle, r.Color, "Next: resolve the failed checks and run doctor again."))
		return nil
	}
	fmt.Fprintln(r.Out, StyleText(SuccessStyle, r.Color, "All checked prerequisites passed."))
	return nil
}

func (r Renderer) DomainInfo(report directory.DomainInfoReport) error {
	if r.JSON {
		var info *ErrorInfo
		if report.Error != nil {
			info = &ErrorInfo{
				Code:    report.Error.Code,
				Message: report.Error.Message,
				Detail:  report.Error.Detail,
				Next:    report.Error.Next,
			}
		}
		return r.writeEnvelope(Envelope{OK: report.Successful(), Data: report, Error: info})
	}

	fmt.Fprintln(r.Out, StyleText(TitleStyle, r.Color, "Cortex domain info"))
	if len(report.Checks) > 0 {
		checks := make([]Check, 0, len(report.Checks))
		for _, check := range report.Checks {
			checks = append(checks, Check{Name: check.Name, Status: check.Status, Detail: check.Detail, OK: check.OK})
		}
		fmt.Fprintln(r.Out, RenderChecks(checks, r.Color, TerminalWidth(r.Out)))
	}
	if report.Info != nil {
		fmt.Fprintln(r.Out)
		fmt.Fprintln(r.Out, StyleText(MutedStyle, r.Color, "Domain details"))
		fmt.Fprintln(r.Out, RenderKeyValue("DNS name", report.Info.DNSName, r.Color))
		fmt.Fprintln(r.Out, RenderKeyValue("NetBIOS", report.Info.NetBIOSName, r.Color))
		fmt.Fprintln(r.Out, RenderKeyValue("Distinguished name", report.Info.DistinguishedName, r.Color))
		fmt.Fprintln(r.Out, RenderKeyValue("Forest", report.Info.ForestName, r.Color))
		fmt.Fprintln(r.Out, RenderKeyValue("Current user", report.Info.CurrentUser, r.Color))
		fmt.Fprintln(r.Out, RenderKeyValue("Target", report.Info.Target, r.Color))
	}
	if report.Error != nil {
		fmt.Fprintln(r.Out)
		fmt.Fprintln(r.Out, StyleText(ErrorStyle, r.Color, report.Error.Message))
		if report.Error.Next != "" {
			fmt.Fprintln(r.Out, StyleText(WarningStyle, r.Color, "Next: "+report.Error.Next))
		}
		return nil
	}
	if report.Status == "partial" {
		fmt.Fprintln(r.Out, StyleText(WarningStyle, r.Color, "Domain information is partial."))
	} else {
		fmt.Fprintln(r.Out, StyleText(SuccessStyle, r.Color, "Domain information retrieved successfully."))
	}
	return nil
}

func (r Renderer) writeEnvelope(envelope Envelope) error {
	return json.NewEncoder(r.Out).Encode(envelope)
}

type Check struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail"`
	OK     bool   `json:"ok"`
}

type DoctorResult struct {
	Checks   []Check `json:"checks"`
	Failures int     `json:"failures"`
}
