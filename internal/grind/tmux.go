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

// Tmux status-right markup tokens. Tmux does not interpret raw ANSI
// escape sequences inside `#(command)` output — it only recognizes its
// own `#[fg=...]` / `#[bold]` directives. Emitting ANSI here would
// render as literal escape text in the status bar.
const (
	tmuxReset   = "#[default]"
	tmuxBold    = "#[bold]"
	tmuxHotPink = "#[fg=#ff6ec7]"
	tmuxDimPink = "#[fg=#7a3a60]"
	tmuxDim     = "#[fg=colour240]" // paused
	tmuxDrained = "#[fg=colour237]" // empty cells
)

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

func hexToTmux(hex string) string {
	return fmt.Sprintf("#[fg=%s]", hex)
}

// EmitTmuxStatus prints one line for tmux `status-right`. When no timer
// state exists or the owner PID is dead, emits nothing (tmux renders an
// empty slot).
func EmitTmuxStatus() {
	s, err := readState()
	if err != nil {
		return
	}
	if !isProcessAlive(s.PID) {
		clearState()
		return
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
				color = tmuxHotPink
			} else {
				color = tmuxDimPink
			}
		case paused:
			color = tmuxDim
		case char == "\u2591":
			color = tmuxDrained
		default:
			color = hexToTmux(barGradient[i])
		}
		bar.WriteString(color + tmuxBold + char)
	}
	bar.WriteString(tmuxReset)

	var timeColor, timeStr string
	switch {
	case expired:
		if blinkBright {
			timeColor = tmuxHotPink
		} else {
			timeColor = tmuxDimPink
		}
		timeStr = "\u2191" + formatDuration(elapsedMs-s.DurationMs)
	case paused:
		timeColor = tmuxDim
		timeStr = "\u23f8 " + formatDuration(s.DurationMs-elapsedMs)
	default:
		idx := full
		if idx >= barWidth {
			idx = barWidth - 1
		}
		timeColor = hexToTmux(barGradient[idx])
		timeStr = formatDuration(s.DurationMs - elapsedMs)
	}

	fmt.Printf("%s  %s%s%s%s", bar.String(), timeColor, tmuxBold, timeStr, tmuxReset)
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
