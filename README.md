# Cortex

Cortex is a friendly command-line tool for safe Windows Active Directory administration.

## Current status

The initial CLI foundation is in place. It currently provides help, version reporting, environment diagnostics, and an interactive session. Active Directory operations are not implemented yet.

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
cortex --json doctor
cortex --json domain info
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

## Domain information

`domain info` is read-only and uses the current Windows security token. It never prompts for, stores, or accepts credentials. The command checks domain membership, the current identity, PowerShell and RSAT availability, the secure channel, and domain-controller access.

If the machine is not domain joined, the current user is not domain authenticated, RSAT is unavailable, the domain cannot be reached, or access is denied, Cortex reports a friendly explanation and next step. Technical details are available in the JSON `error.detail` field.

The `--domain` and `--server` options are mutually exclusive. Both require the current Windows token to have access to the requested target.

## Interactive session

Run `cortex` from an interactive terminal:

```text
cortex> doctor
cortex> version
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
