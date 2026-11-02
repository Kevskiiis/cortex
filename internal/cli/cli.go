package cli

import (
	"context"
	"errors"
	"io"

	"cortex/internal/cli/ui"
)

func Run(args []string, in io.Reader, out, errOut io.Writer, version string) int {
	return RunContext(context.Background(), args, in, out, errOut, version)
}

func RunContext(ctx context.Context, args []string, in io.Reader, out, errOut io.Writer, version string) int {
	root := NewRootCommand(version)
	root.SetArgs(args)
	root.SetIn(in)
	root.SetOut(out)
	root.SetErr(errOut)
	root.SetContext(ctx)

	if err := root.Execute(); err != nil {
		var reported *ui.ReportedError
		if errors.As(err, &reported) {
			return 1
		}
		if jsonOutput(root) {
			_ = ui.NewRenderer(out, true).Error(err)
		} else {
			_ = ui.NewRenderer(errOut, false).Error(err)
		}
		return 1
	}
	return 0
}
