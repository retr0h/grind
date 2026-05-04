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

package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/retr0h/grind/internal/grind"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the timer (foreground TUI, or --bar for tmux mode)",
	Long: `Run the pixel-art coffee-cup timer in the foreground. Use --bar to
drive the tmux status line via ~/.grind/state.json instead.`,
	Args: cobra.NoArgs,
	RunE: func(_ *cobra.Command, _ []string) error {
		duration, err := time.ParseDuration(timerFlag)
		if err != nil {
			return fmt.Errorf("invalid --timer: %w", err)
		}
		if duration <= 0 {
			return fmt.Errorf("--timer must be positive")
		}
		bell := !noBellFlag
		if barFlag {
			return grind.RunBar(duration, bell)
		}
		return grind.Run(duration, bell)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
