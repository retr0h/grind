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
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// Run starts the foreground (full-screen) timer UI. Blocks until the user
// exits via one of the vim-style exit keys (q, Q, ESC, ZZ, :q<CR>,
// Ctrl+C).
//
// When bell is true, a single BEL (\a) is emitted to the current TTY at
// the moment the timer first expires — most terminals surface that as a
// tab-activity highlight or system beep.
func Run(duration time.Duration, bell bool) error {
	tmr := newTimer(duration)

	_ = writeState(tmr)
	defer clearState()

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer func() { _ = term.Restore(fd, oldState) }()

	hideCursor()
	defer showCursor()
	defer clearScreen()

	winchCh := make(chan os.Signal, 1)
	signal.Notify(winchCh, syscall.SIGWINCH)
	defer signal.Stop(winchCh)

	doneCh := make(chan struct{})
	defer close(doneCh)

	keyCh := make(chan byte, 4)
	go func() {
		buf := make([]byte, 1)
		backoff := time.NewTimer(0)
		if !backoff.Stop() {
			<-backoff.C
		}
		defer backoff.Stop()
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				backoff.Reset(10 * time.Millisecond)
				select {
				case <-doneCh:
					return
				case <-backoff.C:
					continue
				}
			}
			if n > 0 {
				select {
				case keyCh <- buf[0]:
				case <-doneCh:
					return
				}
			}
		}
	}()

	frameTick := time.NewTicker(100 * time.Millisecond)
	defer frameTick.Stop()
	lastStateWrite := time.Now()

	var commandMode bool
	var commandBuf []byte
	var lastKey byte
	var belled bool

	draw(tmr, time.Now(), commandMode, string(commandBuf))
	for {
		select {
		case <-winchCh:
			draw(tmr, time.Now(), commandMode, string(commandBuf))
		case key := <-keyCh:
			if commandMode {
				switch key {
				case 0x1b:
					commandMode = false
					commandBuf = nil
				case '\r', '\n':
					cmd := strings.TrimSpace(string(commandBuf))
					commandMode = false
					commandBuf = nil
					if cmd == "q" || cmd == "quit" {
						return nil
					}
				case 0x7f, 8:
					if len(commandBuf) > 0 {
						commandBuf = commandBuf[:len(commandBuf)-1]
					}
				default:
					if key >= 32 && key < 127 {
						commandBuf = append(commandBuf, key)
					}
				}
				draw(tmr, time.Now(), commandMode, string(commandBuf))
				lastKey = key
				continue
			}
			switch key {
			case 0x1b, 'q', 'Q', 3:
				return nil
			case 'Z':
				if lastKey == 'Z' {
					return nil
				}
			case ' ':
				if !tmr.expired(time.Now()) {
					if tmr.isPaused() {
						tmr.resume(time.Now())
					} else {
						tmr.pause(time.Now())
					}
					_ = writeState(tmr)
					lastStateWrite = time.Now()
				}
			case ':':
				commandMode = true
				commandBuf = nil
			}
			lastKey = key
			draw(tmr, time.Now(), commandMode, string(commandBuf))
		case now := <-frameTick.C:
			if bell && !belled && tmr.expired(now) {
				fmt.Print("\a")
				belled = true
			}
			if now.Sub(lastStateWrite) >= time.Second {
				_ = writeState(tmr)
				lastStateWrite = now
			}
			draw(tmr, now, commandMode, string(commandBuf))
		}
	}
}

func draw(t *timer, now time.Time, commandMode bool, commandBuf string) {
	clearScreen()
	w, h := getTermSize()

	frame := composeFrame(t, now, commandMode, commandBuf)

	lines := strings.Split(frame, "\n")
	startRow := (h - len(lines)) / 2
	if startRow < 1 {
		startRow = 1
	}

	fmt.Printf("\033[%d;0H", startRow)
	for i, line := range lines {
		fmt.Printf("\033[%d;0H%s", startRow+i, centerText(line, w))
	}

	// Ambient CRT-style glitches outside the frame — brief white/dim
	// block flickers that disappear on the next frame.
	frameW := cupRenderedWidth()
	frameLeft := (w-frameW)/2 + 1
	if frameLeft < 1 {
		frameLeft = 1
	}
	expired := t.expired(now)
	paused := t.isPaused()
	glitch := ambientGlitch(expired, paused)
	scatterGlitches(w, h, startRow, frameLeft, frameW, len(lines), glitch)

	// Scanline-roll overlays the whole frame, so it always wins — simulates
	// classic CRT h-sync distortion.
	updateScanline(now, h, expired, paused)
	drawScanline(now, w)
}

func composeFrame(t *timer, now time.Time, commandMode bool, commandBuf string) string {
	cupW := cupRenderedWidth()
	expired := t.expired(now)
	paused := t.isPaused()

	pal := palette{
		fill:    t.fillColor,
		outline: mhFG,
		drained: mhDrained,
		timer:   mhFG,
	}

	pulseOn := (now.UnixMilli()/700)%2 == 0
	switch {
	case expired:
		pal.fill = mhPink
		pal.timer = mhPink
		if pulseOn {
			pal.outline = mhPink
		} else {
			pal.outline = mhPinkDim
		}
	case paused:
		pal.fill = mhLavender
		pal.outline = mhLavender
		pal.timer = mhLavender
	}

	cupFill := t.fillPct(now)
	if expired {
		cupFill = 1.0
	}

	glitch := ambientGlitch(expired, paused)

	var parts []string
	parts = append(parts, renderCup(cupFill, glitch, pal))
	parts = append(parts, "")

	var dur time.Duration
	if expired {
		dur = t.expiredFor(now)
	} else {
		dur = t.remaining(now)
	}
	if dur < 0 {
		dur = 0
	}
	secs := int(dur.Seconds())
	mins := secs / 60
	ss := secs % 60
	timerStr := fmt.Sprintf("%02d:%02d", mins, ss)

	timerBlock := renderTextScaled(timerStr, 2, 1)
	timerBlock = applyGlitch(timerBlock, glitch, lipgloss.Color(pal.timer))
	parts = append(parts, centerMultiline(timerBlock, cupW))

	if commandMode {
		parts = append(parts, "")
		parts = append(parts, "")
		parts = append(parts, centerLine(commandBar(commandBuf, now), cupW))
	}

	return strings.Join(parts, "\n")
}

func centerLine(text string, width int) string {
	visible := lipgloss.Width(text)
	if visible >= width {
		return text
	}
	totalPad := width - visible
	leftPad := totalPad / 2
	rightPad := totalPad - leftPad
	return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
}

func centerMultiline(block string, width int) string {
	lines := strings.Split(block, "\n")
	for i, line := range lines {
		lines[i] = centerLine(line, width)
	}
	return strings.Join(lines, "\n")
}
