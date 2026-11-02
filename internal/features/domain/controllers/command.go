package controllers

import (
	"errors"

	"cortex/internal/cli/ui"
	"cortex/internal/features/domain/contracts"

	"github.com/spf13/cobra"
)

func NewCommand(jsonMode *bool, provider contracts.DomainControllersProvider) *cobra.Command {
	var domainTarget string
	var serverTarget string
	command := &cobra.Command{
		Use:   "controllers",
		Short: "List domain controllers and their status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if domainTarget != "" && serverTarget != "" {
				return errors.New("--domain and --server cannot be used together")
			}
			report, err := provider.DomainControllers(cmd.Context(), contracts.DomainControllersRequest{Domain: domainTarget, Server: serverTarget})
			if err != nil {
				return err
			}
			_ = ui.NewRenderer(cmd.OutOrStdout(), *jsonMode).DomainControllers(report)
			if !report.Successful() {
				if report.Error != nil {
					return &ui.ReportedError{Err: errors.New(report.Error.Message)}
				}
				return &ui.ReportedError{Err: errors.New("domain controllers could not be retrieved")}
			}
			return nil
		},
	}
	command.Flags().StringVar(&domainTarget, "domain", "", "Query controllers for a DNS domain name")
	command.Flags().StringVar(&serverTarget, "server", "", "Query through a domain controller")
	return command
}
