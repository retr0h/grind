#!/usr/bin/env bash
# grind.tmux — entry point for the grind tmux plugin.
#
# Install via TPM:
#   set -g @plugin 'retr0h/grind'
# then <prefix> I
#
# Configurable via the user's tmux.conf:
#   set -g @grind-launch-key  'g'   # <prefix> + this to start a timer
#   set -g @grind-stop-key    'G'   # <prefix> + this to stop the timer
#   set -g @grind-status-right '1'  # '1' to auto-prepend grind status to
#                                   # status-right; '0' to leave status-right
#                                   # alone and let the user place
#                                   # `#(grind status)` wherever they want.

set -u

if ! command -v grind >/dev/null 2>&1; then
    tmux display-message "grind: binary not found on \$PATH — see https://github.com/retr0h/grind"
    exit 0
fi

tmux_opt() {
    local name="$1" default="$2"
    local val
    val="$(tmux show-option -gqv "$name")"
    echo "${val:-$default}"
}

launch_key="$(tmux_opt '@grind-launch-key' 'g')"
stop_key="$(tmux_opt '@grind-stop-key' 'G')"
auto_status="$(tmux_opt '@grind-status-right' '1')"

# 1-second status refresh so the progress bar animates smoothly.
tmux set-option -g status-interval 1

# Truecolor so the gradient palette renders correctly. Idempotent — if the
# user already set these, appending `,*:RGB` does no harm.
tmux set-option -g default-terminal "tmux-256color" 2>/dev/null || true
tmux set-option -ga terminal-overrides ",*:RGB"

# Prepend grind's status output to status-right unless the user opted out.
# We check for '#(grind status)' and skip if already present.
if [ "$auto_status" = "1" ]; then
    current_sr="$(tmux show-option -gqv status-right)"
    if [[ "$current_sr" != *'grind status'* ]]; then
        tmux set-option -g status-right "#(grind status) ${current_sr}"
    fi
fi

# Launch: prompt for a Go duration, then run grind in bar mode in the
# background so the UI never takes over a pane.
tmux bind-key "$launch_key" command-prompt -p "grind:" \
    "run-shell -b 'grind --bar --timer %1'"

# Stop / dismiss: SIGTERMs the running grind process and clears state.
tmux bind-key "$stop_key" run-shell 'grind stop'
