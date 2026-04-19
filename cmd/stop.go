package cmd

import (
	"github.com/spf13/cobra"

	"github.com/retr0h/grind/internal/grind"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the running timer (send SIGTERM + clear state)",
	Run: func(_ *cobra.Command, _ []string) {
		grind.Stop()
	},
}
