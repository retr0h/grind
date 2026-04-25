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
	"os/signal"
	"syscall"
	"time"
)

// RunBar starts a headless "bar-only" timer. No terminal rendering, no key
// handling — just a ticker that keeps `~/.grind/state.json` fresh so the
// tmux status bar (driven by `grind status`) can pick up progress. Exits
// cleanly on SIGTERM/SIGINT (e.g. from `grind stop` or <prefix> G).
//
// When bell is true, the first tick that observes an expired timer fires
// a single BEL into the launching tmux pane so the window/tab lights up.
func RunBar(duration time.Duration, bell bool) error {
	tmr := newTimer(duration)

	if err := writeState(tmr); err != nil {
		return err
	}
	defer clearState()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(sigCh)

	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	belled := false
	for {
		select {
		case <-sigCh:
			return nil
		case now := <-tick.C:
			if bell && !belled && tmr.expired(now) {
				RingBarBell()
				belled = true
			}
			_ = writeState(tmr)
		}
	}
}
