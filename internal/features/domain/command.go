package domain

import (
	"cortex/internal/features/domain/contracts"
	"cortex/internal/features/domain/controllers"
	"cortex/internal/features/domain/info"

	"github.com/spf13/cobra"
)

func NewCommand(jsonMode *bool, infoProvider contracts.DomainInfoProvider, controllersProvider contracts.DomainControllersProvider) *cobra.Command {
	command := &cobra.Command{
		Use:   "domain",
		Short: "Inspect the current Active Directory domain",
	}
	command.AddCommand(
		info.NewCommand(jsonMode, infoProvider),
		controllers.NewCommand(jsonMode, controllersProvider),
	)
	return command
}
