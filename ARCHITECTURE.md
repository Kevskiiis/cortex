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
│   ├── cli/                         # Commands, prompts, flags, and human/JSON output
│   │   └── ui/                      # Shared styles, renderers, status, and tables
│   ├── directory/                   # Active Directory contracts and implementations
│   │   ├── powershell/              # PowerShell-backed operations and action scripts
│   │   └── fake/                    # In-memory directory implementation for tests
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
- `internal/cli` owns command definitions, prompts, flags, exit codes, and output formatting.
- `internal/cli/ui` owns terminal styling, human-readable formatting, JSON envelopes, and terminal capability decisions.
- `internal/directory` owns the directory contract and keeps Active Directory details behind a small boundary.
- `internal/directory` also owns the local identity contract used by `whoami`; identity discovery does not require Active Directory access.
- `internal/directory/powershell` implements directory operations through PowerShell and colocates action scripts that must be embedded with their provider.
- `internal/directory/powershell/domain_info.go` adapts the read-only `domain info` action and classifies expected PowerShell failures.
- `internal/directory/powershell/domain_controllers.go` adapts the read-only `domain controllers` action; its script returns controller records as one structured JSON object.
- `internal/directory/fake` provides a predictable backend for unit tests without a real domain.
- `internal/config` loads and validates environment-specific configuration without hardcoding any organization values.
- `internal/safety` enforces preview-by-default behavior, confirmations, protected OUs, and input rules.
- `internal/audit` records actions without storing credentials or other sensitive data.
- `internal/powershell` owns process execution, timeouts, cancellation, and the JSON result protocol.
- `internal/powershell/scripts` contains built-in PowerShell scripts that are embedded into the executable for normal downloads.
- `internal/powershell/runner.go` executes embedded scripts through temporary files with typed process arguments and a 30-second default timeout.
- `test/integration` contains tests that verify behavior across multiple internal packages.
- `docs` contains documentation that helps users install, configure, operate, and contribute to Cortex.
- `scripts/build` contains repeatable build and packaging helpers used locally and by release automation.
- `.github/workflows` contains CI and release workflows for publicly downloadable builds.

## Design direction

Cortex starts as a small modular monolith. The CLI, safety rules, configuration, audit logging, and directory access are separated by responsibility, but remain in one executable. New abstractions such as LDAP, plugins, caching, or interactive sessions should be added only when the corresponding feature is implemented.
