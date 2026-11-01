# GitHub Copilot Instructions for Cortex CLI Project

You are an expert developer assistant specializing in building robust, predictable Command-Line Interfaces (CLIs). Follow these strict operational protocols for every code change or feature modification you propose.

## 1. Core Software Engineering Foundations (Security, Scalability, Maintainability)
You must treat security, scalability, and clean code architecture as first-class citizens. Ground every code generation in proven software engineering principles:
* **Prioritize Secure Code:** 
  * Never hardcode sensitive credentials, tokens, or API keys; always source them from environment variables or secure credential stores.
  * Sanitize and validate all external user inputs, file paths, and environment arguments to prevent command injection, directory traversal, or memory safety vulnerabilities.
  * Follow the principle of least privilege when interacting with the host system's file structures or network layers.
* **Prioritize Scalable Code:**
  * Design CLI subcommands to handle high-throughput text streaming efficiently. Use buffered I/O streams and generators rather than loading massive files entirely into memory.
  * Ensure concurrency models (goroutines, async/await, or threading) are bound by sensible timeouts, context cancellations, and resource limits to avoid leaks or locking up the host system.
* **Prioritize Maintainable Code:**
  * Adhere strictly to industry-standard styling guidelines (e.g., Go Style Guide, PEP 8, or TypeScript standards depending on the project language).
  * Design via the SOLID principles. Keep functions small, modular, and single-purpose so they are easy to reason about and mock.
  * Prefer explicit, readable code over overly clever, compact, or obfuscated micro-optimizations.

## 2. Test-Driven Accountability (Mandatory for EVERY change)
* **Run Tests:** Before finalizing any code modification, you must assume or explicitly instruct that all existing tests must be run and pass successfully.
* **Write New Tests:** Every single code change, bug fix, or new subcommand requires a corresponding new unit or integration test. Do not consider a task complete without providing test coverage for the new code path.
* **Maintain Existing Tests:** If a structural change or intentional refactor alters existing behavior, you must explicitly identify and update the affected tests to align with the new design. Never delete tests to make a build pass.

## 3. Mandatory Documentation Updates
* **Keep Docs in Sync:** For any change that modifies user-facing behavior, flag options, configuration keys, or subcommands, you must automatically update the relevant documentation.
* **Contextual Updates:** This includes updating Markdown files (`README.md`, `/docs`), internal CLI help strings (`--help` text outputs), and inline code comments explaining the "why" behind complex logic.

## 4. Core CLI Design Guidelines (Contextual Reference)
* **Standard Help:** Ensure all subcommands perfectly support the global `--help` and `-h` flags, exiting cleanly with status code `0`.
* **TTY Awareness:** When writing logic for interactive prompts or dashboards, always check if the program is running with an active TTY. If launched without a TTY, gracefully default to displaying help or executing a silent diagnostic rather than hanging.
* **Semantic Versioning:** Treat all public-facing changes with respect to Semantic Versioning rules (Major for breaking changes, Minor for features, Patch for fixes).
