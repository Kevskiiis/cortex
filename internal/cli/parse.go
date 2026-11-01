package cli

import "github.com/mattn/go-shellwords"

func ParseLine(line string) ([]string, error) {
	parser := shellwords.NewParser()
	parser.ParseEnv = false
	parser.ParseBacktick = false
	return parser.Parse(line)
}