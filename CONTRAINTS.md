# Hard constraints

- Preview by default. Nothing changes without --apply or an explicit confirmation, and destructive actions (delete, disable) always ask first. I’d leave deletion out of the early versions.
- No credentials anywhere. Cortex runs as the logged-in user, and passwords never go in config, logs, or the repo.
- Input is never code. Validate names against an allowlist, escape distinguished-name characters, and pass values to PowerShell as separate arguments. Never use Invoke-Expression or build command strings.
- Protected OUs. Cortex refuses to touch the OUs on the protected list, with sensible defaults like Domain Controllers.
- Script integrity. Built-in scripts are embedded in the .exe. Plugin scripts live in a separate folder and show a warning before they run.
- Don’t weaken PowerShell’s safeguards. No global execution policy bypass. If a script is blocked, explain why.
- Nothing hardcoded. No lab-specific OU, group, or naming values in code.
- Targets: Windows PowerShell 5.1, and Windows only.

## Architecture practices

- Go owns the experience, PowerShell owns the work. There is no AD logic in Go.
- One script per action. Each takes typed parameters and prints exactly one JSON object, with errors in a structured field.
- Everything goes through the Directory interface, so backends can be swapped and tests need no real AD.
- Keep dependencies few, pin their versions, and check current docs before adding any.

## User experience

- Error messages answer three questions: what happened, why, and what to try next.
- Consistent command shape: cortex <noun> <verb>, such as ou create and group create.
- Every command has examples in --help, a --json option, and clear exit codes.
- Respect NO_COLOR, never rely on color alone, and show a spinner for slow actions.
- Never prompt in a non-interactive run. Fail with a message naming the missing flag.

## Reliability

- Idempotent commands. Running one twice is harmless.
- Timeouts and Ctrl+C handling on every PowerShell call.
- Cached OU data expires, and it’s cleared after any write.
- An audit log entry for every action, with no sensitive data in it.

## Testing and CI

- Go: table-driven tests against the fake directory.
- PowerShell: Pester tests with the AD commands mocked.
- CI: go vet, a linter, PSScriptAnalyzer, and tests on a Windows runner. Never run write tests against real AD from CI.
- Integration tests run by hand in your lab, after a VM snapshot.