package ui

import (
	"fmt"
	"strings"

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
