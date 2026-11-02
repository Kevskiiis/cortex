package session

import "github.com/chzyer/readline"

func completer() readline.AutoCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("doctor"),
		readline.PcItem("domain",
			readline.PcItem("info"),
			readline.PcItem("controllers"),
		),
		readline.PcItem("version"),
		readline.PcItem("whoami"),
		readline.PcItem("help"),
		readline.PcItem("--help"),
		readline.PcItem("--version"),
		readline.PcItem("quit"),
	)
}
