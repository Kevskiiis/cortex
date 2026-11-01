package ui

import (
	"io"
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	TitleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	SuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	WarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	ErrorStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
	MutedStyle   = lipgloss.NewStyle().Faint(true)
)

func ColorEnabled(out io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	file, ok := out.(*os.File)
	return ok && term.IsTerminal(int(file.Fd()))
}

func TerminalWidth(out io.Writer) int {
	file, ok := out.(*os.File)
	if !ok {
		return 80
	}
	width, _, err := term.GetSize(int(file.Fd()))
	if err != nil || width < 48 {
		return 80
	}
	return width
}

func StyleText(style lipgloss.Style, enabled bool, text string) string {
	if !enabled {
		return text
	}
	return style.Render(text)
}
