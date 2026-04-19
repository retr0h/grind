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
	"github.com/spf13/cobra"

	"github.com/retr0h/grind/internal/grind"
)

var ansiStatusFlag bool

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Emit a single progress-bar line for tmux status-right",
	Long: `Read the current timer state from ~/.grind/state.json and emit a single
line suitable for tmux 'status-right' or 'status-left'. Prints nothing if
no timer is running.

With --ansi, print raw ANSI escape sequences instead of tmux markup. Use
this for direct-to-terminal previews (e.g. VHS recordings); tmux itself
cannot interpret ANSI inside #(...) output and should stick to the default.`,
	Run: func(_ *cobra.Command, _ []string) {
		if ansiStatusFlag {
			grind.EmitAnsiStatus()
			return
		}
		grind.EmitTmuxStatus()
	},
}

func init() {
	statusCmd.Flags().BoolVar(
		&ansiStatusFlag, "ansi", false,
		"Emit ANSI escape sequences instead of tmux #[...] markup (for preview recordings)",
	)
}
