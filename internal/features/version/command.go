package version

import (
	"cortex/internal/cli/ui"

	"github.com/spf13/cobra"
)

func NewCommand(version string, jsonMode *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the Cortex version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			_ = ui.NewRenderer(cmd.OutOrStdout(), *jsonMode).Version(version)
		},
	}
}
