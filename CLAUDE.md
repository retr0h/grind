# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

grind is an 8-bit retro terminal timer. A single duration is passed via `--timer`, and the app renders a pixel-art coffee cup that drains as the clock runs out. On expiry the cup re-fills with hot pink and pulses until the user acknowledges it via a vim-style exit (`ESC`, `ZZ`, or `:q<CR>`).

## Architecture

```
grind/
├── main.go                       # 7-line entry point → cmd.Execute()
├── cmd/                          # cobra CLI tree
│   ├── root.go                   # `grind` + --timer / --bar / --no-bell flags
│   ├── status.go                 # `grind status` (+ --ansi for previews)
│   ├── stop.go                   # `grind stop`
│   └── ringbell.go               # hidden `grind ring-bell` diagnostic
├── internal/grind/               # all implementation
│   ├── timer.go                  # timer type (elapsed/remaining/paused/expired)
│   ├── palette.go                # maxheadroom colors + palette struct + random picker
│   ├── state.go                  # ~/.grind/state.json persistence + Stop()
│   ├── foreground.go             # Run() — full-screen UI, raw term, key dispatch
│   ├── background.go             # RunBar() — headless, drives tmux status bar
│   ├── bell.go                   # RingBarBell() — writes BEL to tmux client/pane ttys
│   ├── tmux.go                   # Emit{Tmux,Ansi}Status() — bar + expiry strobe
│   ├── cup.go                    # Coffee cup pixel art (22×12) + renderCup
│   ├── digits.go                 # 5×7 block glyphs (0–9, A–Z, `:`, `/`, space)
│   ├── hud.go                    # Vim-style `:` command bar with blinking cursor
│   ├── glitch.go                 # applyGlitch() + ambientGlitch()
│   ├── grid.go                   # cellW=3, cellH=2 grid primitives
│   └── terminal.go               # clearScreen / hideCursor / getTermSize / centerText
└── grind.tmux                    # tmux plugin entry point (TPM)
```

Public API of `internal/grind`:

- `Run(duration, bell)` — foreground UI (blocks until user quits)
- `RunBar(duration, bell)` — headless tmux-bar driver (blocks until SIGTERM)
- `EmitTmuxStatus()` — print one tmux `#[...]` status-right line and return
- `EmitAnsiStatus()` — same bar as ANSI escapes (for direct-terminal previews)
- `RingBarBell()` — fire a single BEL into every tmux client/pane tty
- `Stop()` — SIGTERM the running instance, clear state

Dependencies:

- `spf13/cobra` — CLI framework
- `charmbracelet/lipgloss` — colors, bold styling
- `golang.org/x/term` — raw terminal mode, terminal size detection

## Key Technical Details

- Uses raw terminal mode — all output needs `\r\n` not just `\n`
- 10 FPS frame ticker drives glitch animation and the foreground cup's 700ms pulse
- Countdown is calculated from `startedAt + pausedFor`, not a decrementing remaining
- On expiry, `fillPct` is clamped to 1.0 (cup re-fills pink)
- Foreground cup pulses the outline hot pink `#ff6ec7` ↔ dim `#7a3a60` every 700ms
- Bar mode strobes on a 1s wall-clock beat (not 700ms): tmux samples `#(grind status)` once per `status-interval`, so any sub-second period aliases. Odd seconds render bright ▓ hot pink; even seconds collapse to dim ░ pink. Runs indefinitely
- On expiry transition, `Run` writes `"\a"` to its pane's stdout; `RunBar` (no TTY) shells out to tmux and writes BEL to every client tty and every pane tty — that fires both the outer terminal's tab flash and tmux's `monitor-bell` `!` indicator. `--no-bell` suppresses both paths
- Keys read via single persistent stdin goroutine
- SIGWINCH triggers a redraw for terminal resize

## Building

```bash
go build -o grind .                  # Build binary
go run . --timer 10s                 # Quick iteration
```

## Usage

```bash
grind                                # 25 minute default
grind --timer 5m                     # 5 minute timer
grind --timer 1h30m                  # Custom duration
```

## Color Palette (Max Headroom)

```
Random fill (picked at startup):
  Orange    #ffb86c
  Cyan      #00d4ff
  Magenta   #c678dd
  Green     #50fa7b
  Yellow    #e5c07b

Alert:
  Hot pink  #ff6ec7    (expiry fill + timer)
  Dim pink  #7a3a60    (expiry pulse-off)

User paused:
  Lavender  #6272a4

Foreground / outline:
  #c0caf5
```

## Code Standards

- Follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages
- Multi-line function signatures
- golangci-lint with: errcheck, errname, govet, prealloc, predeclared, revive, staticcheck

## Roadmap

- [x] `grind status` subcommand — single-line tmux `status-right` bar sharing state via `~/.grind/state.json`
- [x] `--bar` (headless) mode + `grind stop` for tmux key-binding launches
