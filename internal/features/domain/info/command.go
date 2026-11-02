package info

import (
	"errors"

	"cortex/internal/cli/ui"
	"cortex/internal/features/domain/contracts"

	"github.com/spf13/cobra"
)

func NewCommand(jsonMode *bool, provider contracts.DomainInfoProvider) *cobra.Command {
	var domainTarget string
	var serverTarget string
	command := &cobra.Command{
		Use:   "info",
		Short: "Show domain and current identity information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if domainTarget != "" && serverTarget != "" {
				return errors.New("--domain and --server cannot be used together")
			}
			report, err := provider.DomainInfo(cmd.Context(), contracts.DomainInfoRequest{Domain: domainTarget, Server: serverTarget})
			if err != nil {
				return err
			}
			_ = ui.NewRenderer(cmd.OutOrStdout(), *jsonMode).DomainInfo(report)
			if !report.Successful() {
				if report.Error != nil {
					return &ui.ReportedError{Err: errors.New(report.Error.Message)}
				}
				return &ui.ReportedError{Err: errors.New("domain information could not be retrieved")}
			}
			return nil
		},
	}
	command.Flags().StringVar(&domainTarget, "domain", "", "Query a DNS domain name")
	command.Flags().StringVar(&serverTarget, "server", "", "Query through a domain controller")
	return command
}
