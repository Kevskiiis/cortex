package ui

import (
	"fmt"
	"strings"

	domaincontracts "cortex/internal/features/domain/contracts"

	"github.com/charmbracelet/lipgloss"
)

func RenderChecks(checks []Check, color bool, width int) string {
	if width < 48 {
		width = 48
	}
	nameWidth := 16
	statusWidth := 10
	detailWidth := width - nameWidth - statusWidth - 6
	if detailWidth < 16 {
		detailWidth = 16
	}

	lines := make([]string, 0, len(checks))
	for _, check := range checks {
		icon, style := "!", WarningStyle
		if check.OK {
			icon, style = "✓", SuccessStyle
		}
		if !color {
			if check.OK {
				icon = "[OK]"
			} else {
				icon = "[WARN]"
			}
		}

		name := lipgloss.NewStyle().Width(nameWidth).Render(check.Name)
		status := lipgloss.NewStyle().Width(statusWidth).Render(StyleText(style, color, icon+" "+check.Status))
		detail := lipgloss.NewStyle().Width(detailWidth).Render(check.Detail)
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Top, "  ", name, status, detail))
	}
	return strings.Join(lines, "\n")
}

func RenderKeyValue(key, value string, color bool) string {
	return fmt.Sprintf("  %s %s", StyleText(MutedStyle, color, key+":"), value)
}

func RenderDomainControllers(controllers []domaincontracts.DomainController, color bool, width int) string {
	if len(controllers) == 0 {
		return "  No domain controllers found."
	}
	if width < 48 {
		width = 48
	}

	lines := make([]string, 0, len(controllers)+1)
	if width < 80 {
		for _, controller := range controllers {
			state := controllerState(controller)
			lines = append(lines, fmt.Sprintf("  %s  %s  %s",
				truncate(controller.Name, 22),
				truncate(controller.IPv4Address, 15),
				truncate(state, 18)))
		}
		return strings.Join(lines, "\n")
	}

	nameWidth := 28
	ipWidth := 16
	siteWidth := 20
	stateWidth := 18
	osWidth := width - nameWidth - ipWidth - siteWidth - stateWidth - 10
	if osWidth < 16 {
		osWidth = 16
	}
	header := fmt.Sprintf("  %-*s %-*s %-*s %-*s %s", nameWidth, "Name", ipWidth, "IPv4", siteWidth, "Site", stateWidth, "State", "Operating system")
	lines = append(lines, StyleText(MutedStyle, color, header))
	for _, controller := range controllers {
		lines = append(lines, fmt.Sprintf("  %-*s %-*s %-*s %-*s %s",
			nameWidth, truncate(controller.Name, nameWidth),
			ipWidth, truncate(controller.IPv4Address, ipWidth),
			siteWidth, truncate(controller.Site, siteWidth),
			stateWidth, truncate(controllerState(controller), stateWidth),
			truncate(controller.OperatingSystem, osWidth)))
	}
	return strings.Join(lines, "\n")
}

func controllerState(controller domaincontracts.DomainController) string {
	state := "disabled"
	if controller.Enabled {
		state = "enabled"
	}
	if controller.Writable {
		state += ", writable"
	} else {
		state += ", read-only"
	}
	if controller.GlobalCatalog {
		state += ", GC"
	}
	return state
}

func truncate(value string, width int) string {
	runes := []rune(value)
	if len(runes) <= width {
		return value
	}
	if width <= 1 {
		return string(runes[:width])
	}
	return string(runes[:width-1]) + "…"
}
