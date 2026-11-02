package cli

import (
	"cortex/internal/cli/ui"
	"cortex/internal/features/doctor"
	domainfeature "cortex/internal/features/domain"
	domaincontracts "cortex/internal/features/domain/contracts"
	domaincontrollers "cortex/internal/features/domain/controllers/providers/powershell"
	domaininfo "cortex/internal/features/domain/info/providers/powershell"
	featureversion "cortex/internal/features/version"
	"cortex/internal/features/whoami"
	whoamicontacts "cortex/internal/features/whoami/contracts"
	whoamipowershell "cortex/internal/features/whoami/providers/powershell"

	"github.com/spf13/cobra"
)

const DefaultVersion = "0.1.0-dev"

func NewRootCommand(version string) *cobra.Command {
	return NewRootCommandWithAllProviders(
		version,
		domaininfo.NewProvider(),
		whoamipowershell.NewProvider(),
		domaincontrollers.NewProvider(),
	)
}

func NewRootCommandWithProvider(version string, provider domaincontracts.DomainInfoProvider) *cobra.Command {
	return NewRootCommandWithAllProviders(
		version,
		provider,
		whoamipowershell.NewProvider(),
		domaincontrollers.NewProvider(),
	)
}

func NewRootCommandWithProviders(version string, provider domaincontracts.DomainInfoProvider, identityProvider whoamicontacts.IdentityProvider) *cobra.Command {
	return NewRootCommandWithAllProviders(version, provider, identityProvider, domaincontrollers.NewProvider())
}

func NewRootCommandWithAllProviders(version string, provider domaincontracts.DomainInfoProvider, identityProvider whoamicontacts.IdentityProvider, controllerProvider domaincontracts.DomainControllersProvider) *cobra.Command {
	if version == "" {
		version = DefaultVersion
	}
	var jsonMode bool
	var versionMode bool

	root := &cobra.Command{
		Use:           "cortex",
		Short:         "A safe, friendly command-line tool for Windows Active Directory",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if versionMode {
				return ui.NewRenderer(cmd.OutOrStdout(), jsonMode).VersionFlag(version)
			}
			return cmd.Help()
		},
	}
	root.PersistentFlags().BoolVarP(&versionMode, "version", "v", false, "Print the Cortex version")
	root.PersistentFlags().BoolVar(&jsonMode, "json", false, "Emit machine-readable JSON")
	root.AddCommand(
		featureversion.NewCommand(version, &jsonMode),
		doctor.NewCommand(&jsonMode),
		domainfeature.NewCommand(&jsonMode, provider, controllerProvider),
		whoami.NewCommand(&jsonMode, identityProvider),
	)
	return root
}
