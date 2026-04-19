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
	"os"
	"os/exec"
	"strings"
)

// RingBarBell fires a single terminal BEL (0x07) on behalf of a headless
// `grind --bar` process. The bar daemon doesn't own a TTY — it was
// detached by tmux's `run-shell -b` — so `fmt.Print("\a")` goes nowhere.
//
// Two complementary writes:
//  1. Every attached tmux client's tty (the outer terminal emulator);
//     this flashes the window tab in iTerm2 / Terminal.app.
//  2. Every pane's tty; feeding BEL into a pane's output makes tmux's
//     own `monitor-bell` light up the window tab in the status line.
//
// Either write can fail silently; we'd rather stay quiet than spam
// errors from a one-shot alert.
//
// Exported so a small subcommand (`grind ring-bell`) can fire it from
// the CLI for quick verification.
func RingBarBell() {
	if os.Getenv("TMUX") == "" {
		return
	}
	for _, tty := range tmuxPaths("list-clients", "#{client_tty}") {
		writeBEL(tty)
	}
	for _, tty := range tmuxPaths("list-panes", "-a", "-F", "#{pane_tty}") {
		writeBEL(tty)
	}
}

func tmuxPaths(args ...string) []string {
	out, err := exec.Command("tmux", args...).Output()
	if err != nil {
		return nil
	}
	var paths []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			paths = append(paths, line)
		}
	}
	return paths
}

func writeBEL(path string) {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()
	_, _ = f.Write([]byte{0x07})
}
