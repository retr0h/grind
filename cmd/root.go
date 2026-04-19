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
	"time"

	"github.com/spf13/cobra"

	"github.com/retr0h/grind/internal/grind"
)

var (
	timerFlag string
	barFlag   bool
)

var rootCmd = &cobra.Command{
	Use:   "grind",
	Short: "An 8-bit retro terminal timer",
	Long: `grind is a glanceable countdown with a pixel-art coffee cup that drains
as time runs out. When the timer expires, the cup re-fills with hot pink
and pulses until you acknowledge it.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		duration, err := time.ParseDuration(timerFlag)
		if err != nil {
			return fmt.Errorf("invalid --timer: %w", err)
		}
		if duration <= 0 {
			return fmt.Errorf("--timer must be positive")
		}
		if barFlag {
			return grind.RunBar(duration)
		}
		return grind.Run(duration)
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
	rootCmd.AddCommand(statusCmd, stopCmd)
}
