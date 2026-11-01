package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"

	"cortex/internal/cli/ui"
	"cortex/internal/directory"
	directorypowershell "cortex/internal/directory/powershell"

	"github.com/spf13/cobra"
)

const DefaultVersion = "0.1.0-dev"

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
		var reported *ReportedError
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

type ReportedError struct {
	Err error
}

func (e *ReportedError) Error() string {
	return e.Err.Error()
}

func (e *ReportedError) Unwrap() error {
	return e.Err
}

func NewRootCommand(version string) *cobra.Command {
	return NewRootCommandWithProvider(version, directorypowershell.NewProvider())
}

func NewRootCommandWithProvider(version string, provider directory.DomainInfoProvider) *cobra.Command {
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
	root.AddCommand(newVersionCommand(version, &jsonMode), newDoctorCommand(&jsonMode), newDomainCommand(&jsonMode, provider))
	return root
}

func newVersionCommand(version string, jsonMode *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the Cortex version",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			_ = ui.NewRenderer(cmd.OutOrStdout(), *jsonMode).Version(version)
		},
	}
}

func newDoctorCommand(jsonMode *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check the local environment for Cortex prerequisites",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			checks, failures := runDoctor()
			_ = ui.NewRenderer(cmd.OutOrStdout(), *jsonMode).Doctor(checks, failures)
			if failures > 0 {
				return &ReportedError{Err: fmt.Errorf("%d prerequisite check(s) failed", failures)}
			}
			return nil
		},
	}
}

func newDomainCommand(jsonMode *bool, provider directory.DomainInfoProvider) *cobra.Command {
	domainCommand := &cobra.Command{
		Use:   "domain",
		Short: "Inspect the current Active Directory domain",
	}
	var domainTarget string
	var serverTarget string
	domainInfoCommand := &cobra.Command{
		Use:   "info",
		Short: "Show domain and current identity information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if domainTarget != "" && serverTarget != "" {
				return errors.New("--domain and --server cannot be used together")
			}
			report, err := provider.DomainInfo(cmd.Context(), directory.DomainInfoRequest{
				Domain: domainTarget,
				Server: serverTarget,
			})
			if err != nil {
				return err
			}
			_ = ui.NewRenderer(cmd.OutOrStdout(), *jsonMode).DomainInfo(report)
			if !report.Successful() {
				if report.Error != nil {
					return &ReportedError{Err: errors.New(report.Error.Message)}
				}
				return &ReportedError{Err: errors.New("domain information could not be retrieved")}
			}
			return nil
		},
	}
	domainInfoCommand.Flags().StringVar(&domainTarget, "domain", "", "Query a DNS domain name")
	domainInfoCommand.Flags().StringVar(&serverTarget, "server", "", "Query through a domain controller")
	domainCommand.AddCommand(domainInfoCommand)
	return domainCommand
}

func runDoctor() ([]ui.Check, int) {
	failures := 0
	checks := make([]ui.Check, 0, 4)

	if runtime.GOOS == "windows" {
		checks = append(checks, ui.Check{Name: "Platform", Status: "Windows", Detail: "Supported operating system", OK: true})
	} else {
		checks = append(checks, ui.Check{Name: "Platform", Status: runtime.GOOS, Detail: "Cortex targets Windows", OK: false})
		failures++
	}

	if _, err := exec.LookPath("powershell.exe"); err == nil {
		checks = append(checks, ui.Check{Name: "PowerShell", Status: "Found", Detail: "PowerShell executable is available", OK: true})
	} else {
		checks = append(checks, ui.Check{Name: "PowerShell", Status: "Missing", Detail: "Install Windows PowerShell 5.1", OK: false})
		failures++
	}

	checks = append(checks,
		ui.Check{Name: "RSAT", Status: "Pending", Detail: "Will be checked with directory features", OK: true},
		ui.Check{Name: "AD access", Status: "Pending", Detail: "Will be checked with directory features", OK: true},
	)
	return checks, failures
}

func jsonOutput(root *cobra.Command) bool {
	flag := root.PersistentFlags().Lookup("json")
	return flag != nil && flag.Value.String() == "true"
}
