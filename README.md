# Cortex

Cortex is a friendly command-line tool for safe Windows Active Directory administration.

This README is intended to give you a clear picture of what the project is today, how it is designed, and where it is going.

> **Development status:** Cortex is actively being developed. Features, commands, and design details may change as the project progresses.

## Why Cortex

Active Directory administration often involves repetitive commands, environment-specific details, and tools that are difficult to automate safely. Cortex is being built to make those workflows easier to understand and safer to run:

- **Friendly:** clear explanations, useful next steps, and human-readable terminal output.
- **Safe by default:** read-only diagnostics first, with preview and confirmation planned for future changes.
- **Scriptable:** one-shot commands and stable JSON output for automation.
- **Least privilege:** Cortex uses the current Windows security context and does not manage or store credentials.
- **Windows-native:** designed for administrators working with Windows, PowerShell, and Active Directory.

## Current status

The initial CLI foundation is in place. It currently provides help, version reporting, local identity diagnostics, domain diagnostics, and an interactive session. Active Directory write operations are not implemented yet.

This is an early-stage project, not a replacement for Active Directory Users and Computers or a production-ready directory automation platform. The current implementation is intentionally focused on a small, testable foundation before write operations are introduced.

## Who this is for

- **Systems administrators and help desk teams** who want clearer command-line diagnostics.
- **Junior administrators** who benefit from guardrails and actionable error messages.
- **Experienced administrators** who want repeatable commands that can be used in scripts.
- **Developers and contributors** interested in Go CLIs, PowerShell integration, Windows administration, and testable system boundaries.

## Quick start

### Prerequisites

- Windows.
- Go 1.23 or newer to build from source.
- PowerShell available on the host.
- For domain information checks: a Windows security context with the appropriate domain access. RSAT and Active Directory access checks are part of the directory feature path.

### Build and run

```powershell
go build -o cortex.exe ./cmd/cortex
.\cortex.exe --help
.\cortex.exe doctor
```

For a release-ready Windows archive, see [Releases](#releases).

## Safety and privacy

Cortex is designed around the current Windows user context. It does not prompt for, accept, or store Active Directory credentials, and it does not embed organization-specific domain names, OU names, or secrets.

The current commands are diagnostic and read-only. Future directory-changing commands are intended to preview their actions and require an explicit confirmation or `--apply` before making changes. Those capabilities are not available yet, so the current status above should be treated as the source of truth.

## Build

```powershell
go build -o cortex.exe ./cmd/cortex
```

The default development version is `0.1.0-dev`. Release builds will inject a SemVer version at build time.

## One-shot commands

```text
cortex --help
cortex --version
cortex version
cortex doctor
cortex domain info
cortex domain info --domain corp.example.com
cortex domain info --server dc01.corp.example.com
cortex domain controllers
cortex domain controllers --domain corp.example.com
cortex domain controllers --server dc01.corp.example.com
cortex whoami
cortex --json doctor
cortex --json domain info
cortex --json domain controllers
cortex --json whoami
```

`cortex doctor` currently checks the operating system and PowerShell availability. RSAT and Active Directory access checks will be added with the directory features.

Human-readable output uses automatic terminal formatting with green success states, amber warnings, and Unicode status markers when supported. Set `NO_COLOR` to disable color and receive plain terminal output.

The `--json` flag emits a stable machine-readable envelope:

```json
{
	"ok": true,
	"data": {
		"version": "0.1.0-dev"
	}
}
```

JSON output never includes ANSI color codes or human-readable status text.

## Identity information

`whoami` reports the Windows security context running Cortex. It is local and read-only, so it works without Active Directory connectivity or RSAT. It reports the current user name, SID, authentication type, computer name, domain or workgroup, and domain-joined status.

The command never queries group membership, prompts for credentials, stores credentials, or changes the system.

## Domain information

`domain info` is read-only and uses the current Windows security token. It never prompts for, stores, or accepts credentials. The command checks domain membership, the current identity, PowerShell and RSAT availability, the secure channel, and domain-controller access.

If the machine is not domain joined, the current user is not domain authenticated, RSAT is unavailable, the domain cannot be reached, or access is denied, Cortex reports a friendly explanation and next step. Technical details are available in the JSON `error.detail` field.

The `--domain` and `--server` options are mutually exclusive. Both require the current Windows token to have access to the requested target.

## Domain controllers

`domain controllers` lists readable domain controllers and reports their hostname, IPv4 address, site, operating system, enabled state, writable state, and global-catalog status. It requires a domain-joined Windows computer, RSAT Active Directory tools, and directory read access.

The command is read-only and uses the current Windows token. `--domain` selects a DNS domain and `--server` selects the domain controller used for the query; the options are mutually exclusive. If no controllers are returned, Cortex reports an empty successful result rather than treating it as malformed output.

## Interactive session

Run `cortex` from an interactive terminal:

```text
cortex> doctor
cortex> version
cortex> whoami
cortex> domain controllers
cortex> quit
```

The session provides command history with 100 entries and completion for the available commands. History is stored in the user's application configuration directory. `quit` is the only explicit command that closes the session.

When standard input is not interactive, bare `cortex` displays help instead of waiting for input.

## Releases

Releases use Semantic Versioning tags such as `v0.1.0`. GitHub Actions builds a Windows archive named `cortex-v0.1.0-windows-amd64.zip` and publishes a SHA-256 checksum beside it.

To create the same archive locally from PowerShell:

```powershell
.\scripts\build\release.ps1 -Version v0.1.0
```

## Development

The project is a Go CLI with a PowerShell-backed boundary for Windows directory operations. The main areas are:

- [`cmd/cortex`](./cmd/cortex) — application entry point and dependency wiring.
- [`internal/cli`](./internal/cli) — commands, flags, prompts, sessions, and output.
- [`internal/features`](./internal/features) — feature-owned commands, contracts, providers, scripts, and tests.
- [`internal/powershell`](./internal/powershell) — process execution, timeouts, and result handling.
- [`scripts/build`](./scripts/build) — repeatable local release packaging.

Run the same checks used by continuous integration:

```powershell
go test ./...
go vet ./...
```

Before contributing, please keep changes focused, add tests for new behavior, and document user-facing changes. The project’s design direction and boundaries are described in [`VISION.md`](./VISION.md), [`ARCHITECTURE.md`](./ARCHITECTURE.md), and [`CONTRAINTS.md`](./CONTRAINTS.md).

## Roadmap

The broader direction is to add read-only directory discovery first, followed by guarded write operations with previews, confirmation, and audit logging. Planned areas include:

1. Configuration and domain discovery.
2. Read-only OU and directory queries.
3. Interactive selection and completion.
4. Preview-first OU and group changes.
5. Audit logging and broader integration coverage.

Features are added incrementally; planned items are subject to change and should not be read as currently supported commands.
