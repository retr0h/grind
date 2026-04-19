// Copyright (c) 2026 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package grind

import (
	"fmt"
	"strings"
	"time"
)

// statusPalette isolates the styling tokens the bar renderer emits, so the
// same layout can be printed as tmux markup (`#[fg=...]`) for tmux's
// `status-right` or as raw ANSI escapes for direct-to-terminal preview
// recordings (VHS, screenshots).
type statusPalette struct {
	reset   string
	bold    string
	hotPink string
	dimPink string // expired-blink dim half
	dim     string // paused
	drained string // empty cells
	hexFG   func(string) string
}

var tmuxPalette = statusPalette{
	reset:   "#[default]",
	bold:    "#[bold]",
	hotPink: "#[fg=#ff6ec7]",
	dimPink: "#[fg=#7a3a60]",
	dim:     "#[fg=colour240]",
	drained: "#[fg=colour237]",
	hexFG:   func(hex string) string { return fmt.Sprintf("#[fg=%s]", hex) },
}

var ansiPalette = statusPalette{
	reset:   "\033[0m",
	bold:    "\033[1m",
	hotPink: "\033[38;2;255;110;199m",
	dimPink: "\033[38;2;122;58;96m",
	dim:     "\033[38;5;240m",
	drained: "\033[38;5;237m",
	hexFG:   hexToAnsi,
}

// barGradient assigns a color to each cell position along the bar. Cells
// near the start are cool green (plenty of time); cells near the end are
// pink (running out). As the filled portion grows, more warm cells light
// up — you read "running out" without reading the number.
var barGradient = []string{
	"#50fa7b", "#50fa7b", "#50fa7b",
	"#b4d97b", "#b4d97b",
	"#e5c07b", "#e5c07b", "#e5c07b",
	"#f0be72", "#f0be72",
	"#ffb86c", "#ffb86c", "#ffb86c",
	"#ff9d94", "#ff9d94",
	"#ff8cb3", "#ff8cb3",
	"#ff6ec7", "#ff6ec7", "#ff6ec7",
}

var barWidth = len(barGradient)

func hexToAnsi(hex string) string {
	if len(hex) != 7 || hex[0] != '#' {
		return "\033[38;5;189m"
	}
	var r, g, b int
	if _, err := fmt.Sscanf(hex[1:], "%02x%02x%02x", &r, &g, &b); err != nil {
		return "\033[38;5;189m"
	}
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

// EmitTmuxStatus prints one line of tmux `#[...]` markup for
// `status-right`. When no timer state exists or the owner PID is dead,
// emits nothing (tmux renders an empty slot).
func EmitTmuxStatus() {
	if out := renderStatus(tmuxPalette); out != "" {
		fmt.Print(out)
	}
}

// EmitAnsiStatus prints the same bar as EmitTmuxStatus but with raw ANSI
// escape sequences so the output renders as color in a terminal (not
// inside tmux). Intended for preview recordings and screenshots —
// tmux's `#(...)` substitution cannot interpret ANSI, so the tmux mode
// stays the default.
func EmitAnsiStatus() {
	if out := renderStatus(ansiPalette); out != "" {
		fmt.Print(out)
	}
}

func renderStatus(p statusPalette) string {
	s, err := readState()
	if err != nil {
		return ""
	}
	if !isProcessAlive(s.PID) {
		clearState()
		return ""
	}

	now := time.Now().UnixMilli()
	elapsedMs := now - s.StartMs - s.PausedForMs
	if s.PausedAtMs > 0 {
		elapsedMs = s.PausedAtMs - s.StartMs - s.PausedForMs
	}
	if elapsedMs < 0 {
		elapsedMs = 0
	}

	paused := s.PausedAtMs > 0
	expired := elapsedMs >= s.DurationMs

	progress := float64(elapsedMs) / float64(s.DurationMs)
	if expired {
		progress = 1.0
	}
	if progress > 1 {
		progress = 1
	}
	fractional := progress * float64(barWidth)
	full := int(fractional)
	partial := fractional - float64(full)

	blinkBright := (now/700)%2 == 0

	var bar strings.Builder
	for i := 0; i < barWidth; i++ {
		var char, color string
		switch {
		case i < full:
			char = "\u2593"
		case i == full && partial >= 0.5:
			char = "\u2592"
		default:
			char = "\u2591"
		}

		switch {
		case expired:
			if blinkBright {
				color = p.hotPink
			} else {
				color = p.dimPink
			}
		case paused:
			color = p.dim
		case char == "\u2591":
			color = p.drained
		default:
			color = p.hexFG(barGradient[i])
		}
		bar.WriteString(color + p.bold + char)
	}
	bar.WriteString(p.reset)

	var timeColor, timeStr string
	switch {
	case expired:
		if blinkBright {
			timeColor = p.hotPink
		} else {
			timeColor = p.dimPink
		}
		timeStr = "\u2191" + formatDuration(elapsedMs-s.DurationMs)
	case paused:
		timeColor = p.dim
		timeStr = "\u23f8 " + formatDuration(s.DurationMs-elapsedMs)
	default:
		idx := full
		if idx >= barWidth {
			idx = barWidth - 1
		}
		timeColor = p.hexFG(barGradient[idx])
		timeStr = formatDuration(s.DurationMs - elapsedMs)
	}

	return fmt.Sprintf("%s  %s%s%s%s", bar.String(), timeColor, p.bold, timeStr, p.reset)
}

func formatDuration(ms int64) string {
	total := ms / 1000
	if total < 0 {
		total = 0
	}
	m := total / 60
	s := total % 60
	if m >= 60 {
		h := m / 60
		m = m % 60
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}
