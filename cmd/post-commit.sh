#!/bin/sh
# git-ai-exporter post-commit hook
# Installed by `git-ai-exporter --install-hook`
#
# Configuration hierarchy (priority high to low):
#   1. Project-level: hooks.ai-exporter-url / hooks.ai-exporter-token
#   2. Group-level:   hooks.ai-exporter.group -> hooks.ai-exporter.groups.<group>.url / token
#   3. Global default: hooks.ai-exporter.groups.default.url / token
#
# Setup commands:
#   # Global: configure groups
#   git config --global hooks.ai-exporter.groups.default.url https://your-dashboard.com/api/collect
#   git config --global hooks.ai-exporter.groups.default.token YOUR_TOKEN
#
#   # Per project: join a group
#   git config hooks.ai-exporter.group work
#
#   # Or per project: override completely
#   git config hooks.ai-exporter.url https://custom.com/api/collect
#   git config hooks.ai-exporter.token YOUR_TOKEN

COUNT=$(git config hooks.ai-exporter-count)
COUNT=${COUNT:-1}

HOOK_DIR=$(dirname "$0")
LOG="$HOOK_DIR/git-ai-exporter.log"

git-ai-exporter -r "$(git rev-parse --show-toplevel)" -n "$COUNT" --push --detach >/dev/null 2>>"$LOG"

LOCAL_HOOK=$HOOK_DIR/post-commit.local
if [ -x "$LOCAL_HOOK" ]; then
  exec "$LOCAL_HOOK" "$@"
fi
