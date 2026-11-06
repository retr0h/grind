# Development Guide

## Prerequisites

- macOS (terminal with ANSI + unicode block character support)
- [Go](https://go.dev/dl/) 1.21+
- [just](https://github.com/casey/just) — command runner
- [golangci-lint](https://golangci-lint.run/) — Go linter

## Getting Started

```bash
git clone https://github.com/retr0h/grind.git
cd grind
just fetch    # Fetch shared justfiles
just deps     # Install tool dependencies
```

## Common Commands

```bash
just deps          # Install all dependencies
just test          # Run all tests (lint + format check + unit + coverage)
just ready         # Format, lint before committing
just go-unit       # Run unit tests only
just go-vet        # Run golangci-lint
just go-fmt        # Auto-format (gofumpt + golines)
just just-fmt      # Format justfiles
```

## Running

```bash
go run . --timer 10s             # 10-second iteration check
go run .                         # 25-minute default
```

## Architecture

```
grind/
├── main.go                       # tiny — forwards to cmd.Execute()
├── cmd/                          # cobra CLI
│   ├── root.go                   # `grind` (help) + --timer / --bar / --no-bell
│   ├── start.go                  # `grind start` runs the timer
│   ├── status.go                 # `grind status` (+ --ansi for previews)
│   ├── stop.go                   # `grind stop`
│   ├── ringbell.go               # hidden `grind ring-bell` diagnostic
│   └── version.go                # `grind version`
├── internal/grind/               # implementation
│   ├── timer.go                  # Timer type + elapsed / paused / expired logic
│   ├── palette.go                # Max Headroom colors + palette struct
│   ├── state.go                  # ~/.grind/state.json persistence + Stop()
│   ├── foreground.go             # Run()     — full-screen UI
│   ├── background.go             # RunBar()  — headless tmux driver
│   ├── bell.go                   # RingBarBell(): BEL to tmux client/pane ttys
│   ├── tmux.go                   # Emit{Tmux,Ansi}Status() — gradient bar
│   ├── cup.go                    # Coffee cup pixel art + renderCup
│   ├── digits.go                 # 5×7 block glyphs
│   ├── hud.go                    # Vim-style `:` command bar
│   ├── glitch.go                 # ambientGlitch + applyGlitch
│   ├── grid.go                   # cellW=3, cellH=2 constants
│   └── terminal.go               # raw-term helpers
└── grind.tmux                    # tmux plugin entry point (TPM)
```

`internal/grind` exports six functions:

- `Run(duration, bell)` runs the foreground UI and blocks until the user quits
- `RunBar(duration, bell)` drives the tmux bar headless and blocks until SIGTERM
- `EmitTmuxStatus()` prints one tmux `#[...]` status-right line and returns
- `EmitAnsiStatus()` prints the same bar as ANSI escapes, for previews in a
  plain terminal
- `RingBarBell()` fires one BEL into every tmux client and pane tty
- `Stop()` sends SIGTERM to the running instance and clears state

`Run()` sets up raw terminal mode, launches a single-byte stdin reader
goroutine, ticks a 10 FPS frame timer, and dispatches keys. Key map:

- `q`, `Q`, `ESC`, `Ctrl+C` — quit
- `ZZ` — quit (vim)
- `SPACE` — pause / resume
- `:` — enter command mode (buffered until `<CR>`; `:q<CR>` exits)

## Timer model

`timer` (in `timer.go`) tracks:

- `startedAt` — wall-clock when launched
- `pausedAt` — non-zero while paused
- `pausedFor` — accumulated pause duration

`elapsed` = `now - startedAt - pausedFor` (clamped at `pausedAt` if paused).
`remaining` = `duration - elapsed`. `expired` when `elapsed >= duration`.
`expiredFor` = how long since expiry.

## Visual alert on expiry

When `t.expired(now)`:

- Cup `fillPct` is forced to `1.0` — it re-fills hot pink
- Palette overrides: `fill = timer = mhPink`, outline pulses between `mhPink`
  and `mhPinkDim` every 700ms
- Timer block shows `expiredFor`, counts up
- Ambient glitch rate triples (`0.05` vs `0.015`) — cup crackles

The tmux bar strobes on a 1s wall-clock beat instead of 700ms. tmux samples
`#(grind status)` once per `status-interval`, so any sub-second period aliases.
Odd seconds render bright ▓ hot pink and even seconds collapse to dim ░ pink.
The strobe runs until the timer is stopped.

On the transition to expired, `Run` writes `\a` to its pane's stdout. `RunBar`
has no TTY, because tmux detached it with `run-shell -b`, so it shells out to
tmux and writes BEL to every client tty and every pane tty. The client write
flashes the outer terminal's tab and the pane write lights tmux's `monitor-bell`
`!` indicator. `--no-bell` suppresses both paths.

## Dependencies

| Package                  | Purpose                          |
| ------------------------ | -------------------------------- |
| `spf13/cobra`            | CLI command tree                 |
| `charmbracelet/lipgloss` | Terminal styling, colors         |
| `golang.org/x/term`      | Raw terminal mode, terminal size |

## Raw Terminal Mode

grind uses `term.MakeRaw()` to put the terminal in raw mode:

- No echo (typed characters aren't displayed)
- No line buffering (each keypress is immediate)
- **Important:** `\n` does NOT include carriage return in raw mode — always use
  `\r\n`

Terminal state is restored via `defer term.Restore()` on exit.

## Block Digit Rendering

Each glyph is a 5-wide × 7-tall grid of `X`/` ` characters defined in
`digits.go`. `renderTextScaled(s, scaleW, scaleH)` expands each `X` pixel to a
`scaleW × scaleH` block of `\u2588`. The timer uses scale `2,1` — terminal cells
are taller than wide, so 2-wide scaling gives roughly square pixels.

## Sister Projects

| Project                                                        | Description                              |
| -------------------------------------------------------------- | ---------------------------------------- |
| [tlock](https://github.com/retr0h/tlock)                       | Terminal lock screen with Touch ID       |
| [osapi](https://github.com/osapi-io/osapi)                     | Linux system management REST API and CLI |
| [osapi-justfiles](https://github.com/osapi-io/osapi-justfiles) | Shared justfile recipes for Go projects  |
