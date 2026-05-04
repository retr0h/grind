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

// Package cmd contains the grind cobra command tree.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/retr0h/grind/internal/cli"
)

var (
	timerFlag  string
	barFlag    bool
	noBellFlag bool
)

// rootCmd has no RunE — bare `grind` falls through to the themed
// help, which auto-generates the subcommand list. Running the timer
// is an explicit `grind start`.
var rootCmd = &cobra.Command{
	Use:   "grind",
	Short: "An 8-bit retro terminal timer",
	Long: `grind is a glanceable countdown with a pixel-art coffee cup that drains
as time runs out. When the timer expires, the cup re-fills with hot pink
and pulses until you acknowledge it.`,
	Run: func(c *cobra.Command, _ []string) {
		_ = c.Help()
	},
}

// Execute runs the root command; invoked by main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&timerFlag, "timer", "25m",
		"Timer duration (Go duration: 25m, 1h30m, 45s)",
	)
	rootCmd.PersistentFlags().BoolVar(
		&barFlag, "bar", false,
		"Bar-only mode: no UI, just drive the tmux status bar via ~/.grind/state.json",
	)
	rootCmd.PersistentFlags().BoolVar(
		&noBellFlag, "no-bell", false,
		"Suppress the single terminal bell (\\a) fired when the timer expires",
	)
	rootCmd.AddCommand(statusCmd, stopCmd)

	// Wrap cobra's default help to print the themed banner above it.
	// SetHelpFunc fires for `grind --help` and for the bare-command
	// fallback alike, so the banner shows in both paths without
	// duplicating itself.
	defaultHelp := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		if c == rootCmd {
			out := c.OutOrStdout()
			_, _ = fmt.Fprintln(out)
			_, _ = fmt.Fprint(out, cli.Banner(out))
			_, _ = fmt.Fprintln(out)
		}
		defaultHelp(c, args)
	})
}
