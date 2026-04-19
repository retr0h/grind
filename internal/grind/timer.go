package grind

import "time"

// timer tracks a single countdown with optional pause state. Elapsed time
// excludes any intervals during which the timer was paused.
type timer struct {
	duration  time.Duration
	startedAt time.Time
	pausedAt  time.Time     // zero unless currently paused
	pausedFor time.Duration // accumulated pause time
	fillColor string        // random maxheadroom color picked per-launch
}

func newTimer(d time.Duration) *timer {
	return &timer{
		duration:  d,
		startedAt: time.Now(),
		fillColor: pickRandomColor(),
	}
}

func (t *timer) elapsed(now time.Time) time.Duration {
	base := now
	if !t.pausedAt.IsZero() {
		base = t.pausedAt
	}
	e := base.Sub(t.startedAt) - t.pausedFor
	if e < 0 {
		return 0
	}
	return e
}

func (t *timer) remaining(now time.Time) time.Duration {
	rem := t.duration - t.elapsed(now)
	if rem < 0 {
		return 0
	}
	return rem
}

func (t *timer) expired(now time.Time) bool {
	return t.elapsed(now) >= t.duration
}

func (t *timer) expiredFor(now time.Time) time.Duration {
	if !t.expired(now) {
		return 0
	}
	return t.elapsed(now) - t.duration
}

func (t *timer) isPaused() bool {
	return !t.pausedAt.IsZero()
}

func (t *timer) pause(now time.Time) {
	if t.isPaused() {
		return
	}
	t.pausedAt = now
}

func (t *timer) resume(now time.Time) {
	if !t.isPaused() {
		return
	}
	t.pausedFor += now.Sub(t.pausedAt)
	t.pausedAt = time.Time{}
}

// fillPct returns how full the cup should be: 1.0 at start, 0 at expiry.
func (t *timer) fillPct(now time.Time) float64 {
	if t.duration <= 0 {
		return 0
	}
	pct := 1.0 - float64(t.elapsed(now))/float64(t.duration)
	if pct < 0 {
		return 0
	}
	if pct > 1 {
		return 1
	}
	return pct
}
