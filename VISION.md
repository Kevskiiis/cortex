# Cortex: Vision

## What it is

Cortex is a friendly command-line tool for Windows Active Directory administration. It turns everyday admin tasks, such as creating an OU or a security group, into clear, safe, guided commands that work in any Active Directory environment.

## Who it is for

- Sysadmins and help desk staff who manage Active Directory from a Windows machine.
- Junior admins who want guardrails and clear feedback.
- Experienced admins who want fast, scriptable commands.

## Principles

1. **Friendly first.** Prompts when input is missing, helpful error messages that say what to try next, and examples in every `--help`.
2. **Safe by default.** Every change shows a preview first. Nothing is applied without `--apply` or an explicit confirmation.
3. **Works anywhere.** No hardcoded OU names, group names, or naming conventions. Everything comes from a config file.
4. **Least privilege.** Cortex runs as the logged-in user and never stores credentials, so Active Directory permissions decide what it can do.
5. **Transparent.** Every action is written to an audit log.
6. **Honest.** The README states what has been tested and what hasn't.

## How it works

Cortex has two parts with a clear split:

| Go (the CLI) | PowerShell (the actions) |
|---|---|
| Commands, flags, prompts, and the interactive session | Reading and changing Active Directory |
| Input validation | `-WhatIf` previews |
| Preview, confirm, and `--apply` | The actual changes |
| Config, protected OUs, caching, logging | Returning results as JSON |

Go passes user input to scripts as separate arguments, never as a command string. Each script returns one JSON result that Go turns into friendly output.

All directory work sits behind a `Directory` interface. That allows a PowerShell backend first, a fake backend for tests, and an optional native LDAP backend later.

## How you use it

```
cortex                      opens an interactive session
cortex <command>            runs one command and exits
```

Planned commands (examples, subject to change):

```
cortex init                 detect the domain and write a config file
cortex doctor               check RSAT, PowerShell, and domain access
cortex ou list
cortex ou create            pick a parent OU, enter a name, preview, confirm
cortex group create         create a security group from a naming template
```

When Cortex runs in a script or scheduled task, it never prompts. It fails with a clear message about the missing input.

## Requirements

- Windows with PowerShell 5.1.
- RSAT (Active Directory module) and access to a domain, or a domain controller.
- Permissions in Active Directory for the actions being run.

## Roadmap

1. **Foundation:** project setup, config, `init`, and the friendly error helper.
2. **Read-only:** the `Directory` interface, a fake backend, and `ou list` through PowerShell.
3. **Discovery:** an OU picker with type-to-filter and tab completion.
4. **Changes:** create OU and create group, with preview, `--apply`, and audit logging.
5. **Session:** the interactive session, built on the same commands.
6. **Extensibility:** plugin scripts, tests (Go and Pester), and releases on GitHub.

## Non-goals (for now)

- Replacing the full Active Directory Users and Computers console.
- Running on non-Windows systems.
- Managing Entra ID or other cloud directories.
- Storing or managing credentials.

## Open decisions

- Final name: check GitHub and package managers for conflicts with "Cortex" before creating the repo.
- Session style: shell-style prompt first, with a menu on top later.