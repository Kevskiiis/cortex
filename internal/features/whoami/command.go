package whoami

import (
	"errors"

	"cortex/internal/cli/ui"
	"cortex/internal/features/whoami/contracts"

	"github.com/spf13/cobra"
)

func NewCommand(jsonMode *bool, provider contracts.IdentityProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show the Windows identity running Cortex",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			report, err := provider.WhoAmI(cmd.Context())
			if err != nil {
				return err
			}
			_ = ui.NewRenderer(cmd.OutOrStdout(), *jsonMode).Identity(report)
			if !report.Successful() {
				if report.Error != nil {
					return &ui.ReportedError{Err: errors.New(report.Error.Message)}
				}
				return &ui.ReportedError{Err: errors.New("identity information could not be retrieved")}
			}
			return nil
		},
	}
}
