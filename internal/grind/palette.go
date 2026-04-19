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
