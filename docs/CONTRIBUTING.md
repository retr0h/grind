# Contributing to grind

First off, thanks for taking the time to contribute!

## How Can I Contribute?

### Reporting Bugs

- Use the [GitHub issue tracker](https://github.com/retr0h/grind/issues) to
  report bugs
- Include your macOS version, Go version, and terminal emulator
- Include steps to reproduce the issue

### Suggesting Features

- Open an issue describing the feature you'd like to see
- Explain why this feature would be useful
- Consider whether it fits the project's scope (terminal pomodoro timer)

### Code Contributions

#### Small Fixes

Small changes like typos, grammar fixes, and formatting can be submitted
directly as a pull request.

#### Larger Changes

For bug fixes, new features, or significant changes:

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Make your changes
4. Ensure the project builds: `go build -o grind .`
5. Run the linter: `golangci-lint run`
6. Commit using [Conventional Commits](https://conventionalcommits.org/) format
7. Push to your fork and open a pull request

### Commit Messages

This project uses [Conventional Commits](https://conventionalcommits.org/).
Format: `type(scope): description`

Types: `feat`, `fix`, `docs`, `chore`, `ci`, `build`, `test`, `refactor`

Examples:

```
feat: add block glyphs for phase labels
fix: handle terminal resize during countdown
docs: update README with keybinding table
```

## Development Setup

### Prerequisites

- macOS (terminal with ANSI + unicode block support)
- Go 1.21+

### Building

```bash
git clone https://github.com/retr0h/grind.git
cd grind
go build -o grind .
```

### Testing

```bash
# Quick visual check with a short countdown
./grind --from 10s

# Normal 25-minute run
./grind
```

## Code Style

- Follow existing patterns in the codebase
- Use multi-line function signatures
- Use `\r\n` for output in raw terminal mode (not just `\n`)
- Keep the teal / gray color palette consistent
