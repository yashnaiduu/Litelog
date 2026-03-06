# Contributing to LiteLog

Thank you for your interest in contributing to LiteLog. This document covers how to get set up, what the contribution workflow looks like, and the standards we hold the codebase to.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Making Changes](#making-changes)
- [Commit Standards](#commit-standards)
- [Pull Request Process](#pull-request-process)
- [Reporting Bugs](#reporting-bugs)
- [Requesting Features](#requesting-features)

---

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you agree to uphold these standards.

---

## Getting Started

1. **Fork** the repository on GitHub.
2. **Clone** your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/litelog.git
   cd litelog
   ```
3. **Add the upstream remote:**
   ```bash
   git remote add upstream https://github.com/yashnaiduu/Litelog.git
   ```

---

## Development Setup

**Requirements:**
- Go 1.21+
- Git

**Build from source:**
```bash
go build -o litelog ./cmd/litelog
```

**Run tests:**
```bash
go test ./...
```

**Run the linter/vet:**
```bash
go vet ./...
```

**Start the server locally:**
```bash
./litelog start --retention 7d
```

---

## Project Structure

```
litelog/
├── cmd/litelog/       # CLI entrypoint and subcommand definitions
├── internal/          # Internal packages (not exported)
├── pkg/               # Shared packages
├── server/            # HTTP ingestion server
├── storage/           # SQLite storage engine
├── models/            # Data model definitions
├── benchmarks/        # Benchmark tests
├── assets/            # Static assets (logo, etc.)
├── website/           # Docusaurus documentation site
└── docs/              # Additional documentation
```

---

## Making Changes

1. **Sync with upstream** before starting:
   ```bash
   git fetch upstream
   git checkout main
   git merge upstream/main
   ```

2. **Create a feature branch:**
   ```bash
   git checkout -b feature/my-feature
   # or
   git checkout -b fix/some-bug
   ```

3. **Write your code.** Follow the standards below.

4. **Run tests and vet** before committing:
   ```bash
   go test ./...
   go vet ./...
   ```

---

## Code Standards

- Follow standard Go idiomatic patterns and `gofmt` formatting.
- Keep functions small and focused — single responsibility.
- Avoid unnecessary abstractions. Prefer clarity over cleverness.
- Add comments only for non-obvious logic or important edge cases.
- Never hardcode secrets or environment-specific values.
- Handle errors explicitly — no silent `_` discards for errors that matter.

---

## Commit Standards

Keep commits small, atomic, and clearly named:

```
add retention policy enforcement
fix timestamp parsing for ISO 8601 inputs
refactor ingestion pipeline batch flush logic
docs: update quick-start guide
```

- Use lowercase imperative form (`add`, `fix`, `refactor`, not `Added`, `Fixed`)
- Reference issue numbers where relevant: `fix #42: handle nil pointer in query parser`
- Avoid noisy commits like `WIP`, `minor changes`, or `fix stuff`

---

## Pull Request Process

1. Ensure all tests pass: `go test ./...`
2. Ensure `go vet ./...` produces no output.
3. Write a clear PR description:
   - **What** the change does
   - **Why** it was needed
   - **How** it was tested
4. Link any related GitHub issues.
5. Keep PRs focused — one concern per PR.
6. Be responsive to review feedback.

PRs are merged once they receive a review approval and CI passes.

---

## Reporting Bugs

Open a [GitHub Issue](https://github.com/yashnaiduu/Litelog/issues/new?template=bug_report.md) and include:

- LiteLog version (`litelog version`)
- OS and Go version
- Steps to reproduce
- Expected vs. actual behavior
- Relevant log output or error messages

---

## Requesting Features

Open a [GitHub Issue](https://github.com/yashnaiduu/Litelog/issues/new?template=feature_request.md) describing:

- The problem you are trying to solve
- Your proposed solution or approach
- Any alternatives you considered

---

## License

By contributing to LiteLog, you agree that your contributions will be licensed under the [Apache License 2.0](LICENSE).
