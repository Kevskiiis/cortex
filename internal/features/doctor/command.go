package doctor

import (
	"fmt"

	"cortex/internal/cli/ui"

	"github.com/spf13/cobra"
)

func NewCommand(jsonMode *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check the local environment for Cortex prerequisites",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			checks, failures := runChecks()
			_ = ui.NewRenderer(cmd.OutOrStdout(), *jsonMode).Doctor(checks, failures)
			if failures > 0 {
				return &ui.ReportedError{Err: fmt.Errorf("%d prerequisite check(s) failed", failures)}
			}
			return nil
		},
	}
}
