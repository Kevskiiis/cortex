package main

import (
	"os"

	"cortex/internal/cli"
	"cortex/internal/cli/session"

	"golang.org/x/term"
)

var version = cli.DefaultVersion

func main() {
	if len(os.Args) > 1 {
		os.Exit(cli.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, version))
	}

	if !term.IsTerminal(int(os.Stdin.Fd())) {
		os.Exit(cli.Run([]string{"--help"}, os.Stdin, os.Stdout, os.Stderr, version))
	}

	os.Exit(session.Run(os.Stdin, os.Stdout, os.Stderr, version))
}
