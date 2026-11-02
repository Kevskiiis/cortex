# Cortex Directory Schema

```text
cortex/
├── .github/
│   ├── copilot-instructions.md      # Repository guidance for contributors and agents
│   └── workflows/                   # CI and release automation
│
├── cmd/
│   └── cortex/                      # The downloadable Cortex CLI application
│
├── internal/
│   ├── cli/                         # Root execution, parsing, sessions, and shared output
│   │   └── ui/                      # Shared styles, renderers, status, and tables
│   ├── features/                    # Feature-owned vertical slices
│   │   ├── doctor/                  # Local prerequisite diagnostics
│   │   ├── version/                 # Version command
│   │   ├── whoami/                  # Identity command, contract, provider, and script
│   │   └── domain/                  # Domain parent and info/controllers leaves
│   ├── config/                      # Configuration loading and validation
│   ├── safety/                      # Preview, confirmation, and protected OU rules
│   ├── audit/                       # Audit event creation and persistence
│   └── powershell/                  # PowerShell process runner and protocol handling
│       └── scripts/                 # Built-in scripts embedded in the executable
│
├── test/
│   └── integration/                 # Tests spanning the CLI and directory boundary
│
├── docs/                            # User, contributor, and operational documentation
│
├── scripts/
│   └── build/                       # Local build and release packaging helpers
│
├── CONTRAINTS.md                    # Non-negotiable safety and engineering constraints
├── README.md                        # Installation, prerequisites, and quick start
├── VISION.md                        # Product purpose and roadmap
└── ARCHITECTURE.md                  # This directory and ownership guide
```

## Ownership boundaries

- `cmd/cortex` is the application entry point. It should contain only startup and dependency wiring.
- `internal/cli` owns root execution, argument parsing, sessions, and shared output handling.
- `internal/cli/ui` owns terminal styling, human-readable formatting, JSON envelopes, and terminal capability decisions.
- `internal/features/<feature>` owns a feature's parent command, leaf commands, contracts, models, providers, scripts, and tests.
- A parent feature command only registers its leaves; each leaf owns its flags, validation, provider call, rendering, and exit status.
- `internal/features/domain/info` and `internal/features/domain/controllers` are independent leaves that share only domain contracts and common infrastructure.
- `internal/features/*/providers/powershell` adapts feature contracts to embedded PowerShell actions.
- `internal/config` loads and validates environment-specific configuration without hardcoding any organization values.
- `internal/safety` enforces preview-by-default behavior, confirmations, protected OUs, and input rules.
- `internal/audit` records actions without storing credentials or other sensitive data.
- `internal/powershell` owns process execution, timeouts, cancellation, and the JSON result protocol.
- Feature providers colocate the PowerShell scripts they embed, keeping each action beside its contract adapter.
- `internal/powershell/runner.go` executes embedded scripts through temporary files with typed process arguments and a 30-second default timeout.
- `test/integration` contains tests that verify behavior across multiple internal packages.
- `docs` contains documentation that helps users install, configure, operate, and contribute to Cortex.
- `scripts/build` contains repeatable build and packaging helpers used locally and by release automation.
- `.github/workflows` contains CI and release workflows for publicly downloadable builds.

## Design direction

Cortex is a modular monolith organized by vertical feature slices. Shared infrastructure remains horizontal, but feature commands should not import concrete providers from other features. New capabilities such as OUs, groups, SMB/NFS shares, and GPOs should add new feature packages rather than expanding a central directory interface or command file.
