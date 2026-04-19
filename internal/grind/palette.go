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

import "math/rand"

// maxheadroom palette (80s neon) — matches ~/git/dotfiles nvim theme.
const (
	mhOrange   = "#ffb86c"
	mhCyan     = "#00d4ff"
	mhMagenta  = "#c678dd"
	mhGreen    = "#50fa7b"
	mhYellow   = "#e5c07b"
	mhPink     = "#ff6ec7" // alert only
	mhPinkDim  = "#7a3a60" // alert pulse-off
	mhLavender = "#6272a4" // user-paused
	mhFG       = "#c0caf5"
	mhDrained  = "#3b4261"
)

// Random fill colors chosen at startup. Pink is excluded so the alert pink
// remains visually distinctive when the timer expires.
var randomFillColors = []string{
	mhOrange,
	mhCyan,
	mhMagenta,
	mhGreen,
	mhYellow,
}

func pickRandomColor() string {
	return randomFillColors[rand.Intn(len(randomFillColors))]
}

// palette captures the per-frame color choices consumed by the cup and
// timer renderers. It's derived from the timer's state (active / paused /
// expired) at each frame.
type palette struct {
	fill    string
	outline string
	drained string
	timer   string
}
