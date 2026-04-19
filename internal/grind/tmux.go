package grind

import (
	"fmt"
	"strings"
	"time"
)

// Tmux bar constants.
const (
	ansiReset = "\033[0m"
	ansiBold  = "\033[1m"

	// Expiry blink colors
	ansiHotPink = "\033[38;2;255;110;199m"
	ansiDimPink = "\033[38;2;122;58;96m"

	ansiDim     = "\033[38;5;240m" // paused
	ansiDrained = "\033[38;5;237m" // empty cells
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
				color = ansiHotPink
			} else {
				color = ansiDimPink
			}
		case paused:
			color = ansiDim
		case char == "\u2591":
			color = ansiDrained
		default:
			color = hexToAnsi(barGradient[i])
		}
		bar.WriteString(color + ansiBold + char)
	}
	bar.WriteString(ansiReset)

	var timeColor, timeStr string
	switch {
	case expired:
		if blinkBright {
			timeColor = ansiHotPink
		} else {
			timeColor = ansiDimPink
		}
		timeStr = "\u2191" + formatDuration(elapsedMs-s.DurationMs)
	case paused:
		timeColor = ansiDim
		timeStr = "\u23f8 " + formatDuration(s.DurationMs-elapsedMs)
	default:
		idx := full
		if idx >= barWidth {
			idx = barWidth - 1
		}
		timeColor = hexToAnsi(barGradient[idx])
		timeStr = formatDuration(s.DurationMs - elapsedMs)
	}

	fmt.Printf("%s  %s%s%s%s", bar.String(), timeColor, ansiBold, timeStr, ansiReset)
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
