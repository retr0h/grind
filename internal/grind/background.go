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
func RunBar(duration time.Duration) error {
	tmr := newTimer(duration)

	if err := writeState(tmr); err != nil {
		return err
	}
	defer clearState()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-sigCh:
			return nil
		case <-tick.C:
			_ = writeState(tmr)
		}
	}
}
