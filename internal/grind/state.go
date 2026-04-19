package grind

import (
	"encoding/json"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// storedState is the on-disk representation of a running timer. The main
// `grind` process writes this once per second so a separate `grind status`
// invocation (from tmux) can render the current progress.
type storedState struct {
	StartMs     int64  `json:"start_ms"`
	DurationMs  int64  `json:"duration_ms"`
	PausedAtMs  int64  `json:"paused_at_ms,omitempty"`
	PausedForMs int64  `json:"paused_for_ms,omitempty"`
	Color       string `json:"color"`
	PID         int    `json:"pid"`
}

func stateFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.TempDir()
	}
	return filepath.Join(home, ".grind", "state.json")
}

func (t *timer) toStored() storedState {
	s := storedState{
		StartMs:     t.startedAt.UnixMilli(),
		DurationMs:  int64(t.duration / time.Millisecond),
		PausedForMs: int64(t.pausedFor / time.Millisecond),
		Color:       t.fillColor,
		PID:         os.Getpid(),
	}
	if !t.pausedAt.IsZero() {
		s.PausedAtMs = t.pausedAt.UnixMilli()
	}
	return s
}

func writeState(t *timer) error {
	path := stateFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(t.toStored())
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func readState() (*storedState, error) {
	data, err := os.ReadFile(stateFilePath())
	if err != nil {
		return nil, err
	}
	var s storedState
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func clearState() {
	_ = os.Remove(stateFilePath())
}

// isProcessAlive returns true if a process with the given PID is still
// running. Used to detect and clear stale state files left behind by a
// hard-killed grind process.
func isProcessAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}

// Stop reads the state file, sends SIGTERM to the owner process, and
// removes the file. Invoked by `grind stop`.
func Stop() {
	s, err := readState()
	if err != nil {
		return
	}
	if isProcessAlive(s.PID) {
		if proc, err := os.FindProcess(s.PID); err == nil {
			_ = proc.Signal(syscall.SIGTERM)
		}
	}
	clearState()
}
