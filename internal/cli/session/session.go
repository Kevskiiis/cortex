package session

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cortex/internal/cli"
	"cortex/internal/cli/ui"

	"github.com/chzyer/readline"
)

const HistoryLimit = 100

func Run(in io.ReadCloser, out, errOut io.Writer, version string) int {
	historyFile := historyPath()
	if historyFile != "" {
		_ = os.MkdirAll(filepath.Dir(historyFile), 0o700)
	}

	lineReader, err := readline.NewEx(&readline.Config{
		Prompt:          prompt(),
		HistoryFile:     historyFile,
		HistoryLimit:    HistoryLimit,
		AutoComplete:    completer(),
		InterruptPrompt: "^C",
		Stdin:           in,
		Stdout:          out,
		Stderr:          errOut,
	})
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}
	defer lineReader.Close()

	for {
		line, err := lineReader.Readline()
		if err == readline.ErrInterrupt {
			continue
		}
		if err == io.EOF {
			return 0
		}
		if err != nil {
			fmt.Fprintln(errOut, err)
			return 1
		}
		if executeLine(line, out, errOut, version) {
			return 0
		}
	}
}

func executeLine(line string, out, errOut io.Writer, version string) bool {
	args, err := cli.ParseLine(line)
	if err != nil {
		fmt.Fprintf(errOut, "parse error: %v\n", err)
		return false
	}
	if len(args) == 0 {
		return false
	}
	if len(args) == 1 && args[0] == "quit" {
		return true
	}

	_ = cli.Run(args, nil, out, errOut, version)
	return false
}

func prompt() string {
	return ui.StyleText(ui.TitleStyle, ui.ColorEnabled(os.Stdout), "cortex> ")
}

func completer() readline.AutoCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("doctor"),
		readline.PcItem("version"),
		readline.PcItem("help"),
		readline.PcItem("--help"),
		readline.PcItem("--version"),
		readline.PcItem("quit"),
	)
}

func historyPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(configDir, "Cortex", "history")
}
