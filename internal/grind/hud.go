package grind

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// commandBar renders the vim-style `:` prompt with the user's typed buffer.
// Example:  ░▒▓  : q█  ▓▒░
func commandBar(buffer string, now time.Time) string {
	amber := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB000")).Bold(true)
	shade := lipgloss.NewStyle().Foreground(lipgloss.Color("#5a3a00"))

	cursor := " "
	if (now.UnixMilli()/500)%2 == 0 {
		cursor = "\u2588"
	}
	left := shade.Render("\u2591\u2592\u2593")
	right := shade.Render("\u2593\u2592\u2591")
	prompt := amber.Render(": ") + amber.Render(buffer) + amber.Render(cursor)
	return fmt.Sprintf("%s %s  %s", left, prompt, right)
}
