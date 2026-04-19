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
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Coffee cup pixel art, uniform 22 wide x 12 tall pixels. The body is
// centered horizontally in the glyph (cols 3–18) so the digital countdown
// rendered below it lines up with what your eye reads as "the cup" —
// otherwise the body sits off-center and the timer drifts right of the cup.
// The handle extends into cols 19–21.
//
//	X = outline, A = fill (rows eligible to drain), ' ' = transparent
var cupPixels = []string{
	"                      ",
	"   XXXXXXXXXXXXXXXX   ",
	"   XAAAAAAAAAAAAAAX   ",
	"   XAAAAAAAAAAAAAAXXXX",
	"   XAAAAAAAAAAAAAAX  X",
	"   XAAAAAAAAAAAAAAX  X",
	"   XAAAAAAAAAAAAAAXXXX",
	"   XAAAAAAAAAAAAAAX   ",
	"   XAAAAAAAAAAAAAAX   ",
	"    XAAAAAAAAAAAAX    ",
	"     XXXXXXXXXXXX     ",
	"                      ",
}

func cupRenderedWidth() int {
	if len(cupPixels) == 0 {
		return 0
	}
	return len(cupPixels[0]) * cellW
}

// renderCup returns the cup as a multi-line styled string.
// fillPct 1.0 = full, 0.0 = empty (drain happens top-down). glitch is the
// per-cell flicker probability. pal provides phase-dependent colors.
func renderCup(fillPct, glitch float64, pal palette) string {
	fillRows := make([]int, 0, len(cupPixels))
	for i, row := range cupPixels {
		if strings.ContainsRune(row, 'A') {
			fillRows = append(fillRows, i)
		}
	}
	if fillPct < 0 {
		fillPct = 0
	}
	if fillPct > 1 {
		fillPct = 1
	}
	drainedCount := int(float64(len(fillRows)) * (1.0 - fillPct))
	drained := make(map[int]bool, drainedCount)
	for i := 0; i < drainedCount && i < len(fillRows); i++ {
		drained[fillRows[i]] = true
	}

	full := strings.Repeat("\u2588", cellW)
	shade := strings.Repeat("\u2593", cellW)
	light := strings.Repeat("\u2591", cellW)

	outline := lipgloss.NewStyle().Foreground(lipgloss.Color(pal.outline))
	fill := lipgloss.NewStyle().Foreground(lipgloss.Color(pal.fill))
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color(pal.drained))

	var out strings.Builder
	for rowIdx, pixelRow := range cupPixels {
		for y := 0; y < cellH; y++ {
			for _, pix := range pixelRow {
				switch pix {
				case 'X':
					if glitch > 0 && rand.Float64() < glitch {
						out.WriteString(outline.Render(shade))
					} else {
						out.WriteString(outline.Render(full))
					}
				case 'A':
					if drained[rowIdx] {
						out.WriteString(dim.Render(light))
						continue
					}
					if glitch > 0 && rand.Float64() < glitch {
						out.WriteString(fill.Render(shade))
						continue
					}
					out.WriteString(fill.Render(full))
				default:
					out.WriteString(strings.Repeat(" ", cellW))
				}
			}
			out.WriteString("\n")
		}
	}
	return strings.TrimRight(out.String(), "\n")
}
