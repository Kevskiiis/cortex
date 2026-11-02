package doctor

import (
	"os/exec"
	"runtime"

	"cortex/internal/cli/ui"
)

func runChecks() ([]ui.Check, int) {
	failures := 0
	checks := make([]ui.Check, 0, 4)
	if runtime.GOOS == "windows" {
		checks = append(checks, ui.Check{Name: "Platform", Status: "Windows", Detail: "Supported operating system", OK: true})
	} else {
		checks = append(checks, ui.Check{Name: "Platform", Status: runtime.GOOS, Detail: "Cortex targets Windows", OK: false})
		failures++
	}
	if _, err := exec.LookPath("powershell.exe"); err == nil {
		checks = append(checks, ui.Check{Name: "PowerShell", Status: "Found", Detail: "PowerShell executable is available", OK: true})
	} else {
		checks = append(checks, ui.Check{Name: "PowerShell", Status: "Missing", Detail: "Install Windows PowerShell 5.1", OK: false})
		failures++
	}
	checks = append(checks,
		ui.Check{Name: "RSAT", Status: "Pending", Detail: "Will be checked with directory features", OK: true},
		ui.Check{Name: "AD access", Status: "Pending", Detail: "Will be checked with directory features", OK: true},
	)
	return checks, failures
}
