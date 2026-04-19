package cmd

import (
	"github.com/spf13/cobra"

	"github.com/retr0h/grind/internal/grind"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Emit a single progress-bar line for tmux status-right",
	Long: `Read the current timer state from ~/.grind/state.json and emit a single
line suitable for tmux 'status-right' or 'status-left'. Prints nothing if
no timer is running.`,
	Run: func(_ *cobra.Command, _ []string) {
		grind.EmitTmuxStatus()
	},
}
