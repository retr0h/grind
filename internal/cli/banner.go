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

// Package cli holds CLI-only output helpers (banner, accent, mute) —
// never imported by the running TUI under internal/grind/, only by
// cmd/. Uses lipgloss so we share termenv's NO_COLOR / TTY / color-
// profile detection with the TUI rather than reinventing it.
package cli

import (
	"io"

	"github.com/charmbracelet/lipgloss"
)

// mhOrange is the maxheadroom-palette accent shared with the TUI
// (internal/grind/palette.go). Duplicated here because that package's
// constants are unexported; the source of truth is still palette.go.
const mhOrange = "#ffb86c"

var (
	muteStyle   = lipgloss.NewStyle().Faint(true)
	accentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(mhOrange))
)

// Mute returns s rendered as secondary text (faint).
func Mute(_ io.Writer, s string) string { return muteStyle.Render(s) }

// Accent returns s rendered in the project accent (mhOrange).
func Accent(_ io.Writer, s string) string { return accentStyle.Render(s) }

// Banner returns the GRIND block-letter logo. Top line muted (the
// silhouette), bottom line accent (the weight) — line-level coloring
// matches the install summary so the curl|bash banner and `grind
// --help` look the same.
func Banner(_ io.Writer) string {
	const top = "█▀▀ █▀█ █ █▄█ █▀▄"
	const bot = "█▄█ █▀▄ █ █░█ █▄▀"
	return muteStyle.Render(top) + "\n" + accentStyle.Render(bot) + "\n"
}
