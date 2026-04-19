package grind

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Scanline-roll state — a full-width horizontal band flashes across a
// random row every N seconds to simulate classic CRT h-sync distortion.
// Held at package scope because only one foreground grind runs per
// process and the effect is global to the render loop.
var scanline struct {
	nextFire    time.Time
	activeUntil time.Time
	row         int
	expired     bool // if true, render hot pink instead of white
}

// updateScanline advances the scanline state machine. It schedules the
// next firing time on first call and on every subsequent firing. When the
// timer is paused, it quietly does nothing and lets the cooldown resume
// whenever the user unpauses.
func updateScanline(now time.Time, termH int, expired, paused bool) {
	if paused {
		return
	}
	if scanline.nextFire.IsZero() {
		scanline.nextFire = now.Add(nextScanlineDelay(expired))
		return
	}
	if now.Before(scanline.nextFire) {
		return
	}
	// Fire: pick a row, hold it visible for 100–200ms (1–2 frames).
	if termH < 4 {
		return
	}
	scanline.row = rand.Intn(termH-2) + 1
	scanline.activeUntil = now.Add(time.Duration(100+rand.Intn(100)) * time.Millisecond)
	scanline.expired = expired
	scanline.nextFire = now.Add(nextScanlineDelay(expired))
}

// nextScanlineDelay returns the random wait before the next scanline fires.
// Expired-state runs cycle much faster — the whole screen should feel
// "broken" until the user dismisses.
func nextScanlineDelay(expired bool) time.Duration {
	if expired {
		return time.Duration(1000+rand.Intn(2000)) * time.Millisecond // 1–3s
	}
	return time.Duration(8000+rand.Intn(7000)) * time.Millisecond // 8–15s
}

// drawScanline paints the active horizontal band, if any.
func drawScanline(now time.Time, termW int) {
	if now.After(scanline.activeUntil) {
		return
	}
	color := "#c0caf5" // white
	if scanline.expired {
		color = "#ff6ec7" // hot pink
	}
	band := strings.Repeat("\u2588", termW)
	styled := lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(band)
	fmt.Printf("\033[%d;0H%s", scanline.row, styled)
}

// ambientGlitch returns a per-frame glitch probability used by the cup and
// timer renderers: none when paused, elevated when expired (cup crackles),
// baseline otherwise.
func ambientGlitch(expired, paused bool) float64 {
	if paused {
		return 0
	}
	if expired {
		return 0.05
	}
	return 0.015
}

// scatterGlitches sprinkles transient white block chars at random positions
// OUTSIDE the centered cup/timer frame, creating a CRT-interference feel
// across the whole screen. Because draw() clearScreens on every frame,
// each glitch lasts exactly one frame (100ms) before disappearing — you
// see them as brief flickers, not persistent noise.
//
// termW/termH = terminal size.
// frameTop/Left = 1-indexed top-left of the cup+timer frame.
// frameW/H = frame bounding box in terminal cells.
func scatterGlitches(termW, termH, frameTop, frameLeft, frameW, frameH int, intensity float64) {
	if intensity <= 0 {
		return
	}
	// 0.015 (normal) → 0-1 per frame; 0.05 (expired) → 0-5 per frame.
	count := rand.Intn(int(intensity*100) + 1)
	whiteStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#c0caf5"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#3b4261"))

	for i := 0; i < count; i++ {
		// Up to 5 tries to land outside the frame bounding box.
		for a := 0; a < 5; a++ {
			row := rand.Intn(termH) + 1
			col := rand.Intn(termW) + 1
			if row >= frameTop && row < frameTop+frameH &&
				col >= frameLeft && col < frameLeft+frameW {
				continue
			}
			ch := "\u2593"
			if rand.Intn(3) == 0 {
				ch = "\u2588"
			}
			style := whiteStyle
			if rand.Intn(4) == 0 {
				style = dimStyle
			}
			fmt.Printf("\033[%d;%dH%s", row, col, style.Render(ch))
			break
		}
	}
}

// applyGlitch randomly replaces full-block runes in `block` with shade runes
// and colors them according to the current glitch level. Returns a styled
// (ANSI-escape-wrapped) version of the same string.
func applyGlitch(block string, glitch float64, baseColor lipgloss.Color) string {
	baseStyle := lipgloss.NewStyle().Foreground(baseColor).Bold(true)
	if glitch <= 0 {
		return baseStyle.Render(block)
	}
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	bright := lipgloss.NewStyle().Foreground(baseColor)
	var out strings.Builder
	for _, r := range block {
		if r == '\u2588' && rand.Float64() < glitch {
			switch rand.Intn(3) {
			case 0:
				out.WriteString(dim.Render("\u2591"))
			case 1:
				out.WriteString(bright.Render("\u2593"))
			default:
				out.WriteString(baseStyle.Render("\u2592"))
			}
		} else {
			out.WriteString(baseStyle.Render(string(r)))
		}
	}
	return out.String()
}
