#!/bin/sh
# git-ai-exporter post-commit hook
# Installed by `git-ai-exporter --install-hook`
#
# Configuration hierarchy (priority high to low):
#   1. Project-level: hooks.ai-exporter-url / hooks.ai-exporter-token
#   2. Group-level:   hooks.ai-exporter.group -> hooks.ai-exporter.groups.<group>.url / token
#
# Setup commands:
#   # Global: configure groups
#   git config --global hooks.ai-exporter.groups.work.url https://work-dashboard.com/api/collect
#   git config --global hooks.ai-exporter.groups.work.token YOUR_TOKEN
#
#   # Per project: join a group
#   git config hooks.ai-exporter.group work
#
#   # Or per project: override completely
#   git config hooks.ai-exporter.url https://custom.com/api/collect
#   git config hooks.ai-exporter.token YOUR_TOKEN

resolve_config() {
    # 1. Check project-level URL/token first
    URL=$(git config hooks.ai-exporter-url)
    TOKEN=$(git config hooks.ai-exporter-token)

    if [ -n "$URL" ] && [ -n "$TOKEN" ]; then
        return
    fi

    # 2. Check group config from global
    GROUP=$(git config hooks.ai-exporter.group)
    if [ -n "$GROUP" ]; then
        GROUP_URL=$(git config --global hooks.ai-exporter.groups."$GROUP".url)
        GROUP_TOKEN=$(git config --global hooks.ai-exporter.groups."$GROUP".token)

        # Only use group values if project-level is missing
        URL=${URL:-$GROUP_URL}
        TOKEN=${TOKEN:-$GROUP_TOKEN}
    fi

    # 3. If still empty, skip silently
    if [ -z "$URL" ] || [ -z "$TOKEN" ]; then
        exit 0
    fi
}

resolve_config

COUNT=$(git config hooks.ai-exporter-count)
COUNT=${COUNT:-1}

HOOK_DIR=$(dirname "$0")
LOG="$HOOK_DIR/git-ai-exporter.log"

git-ai-exporter -r "$(git rev-parse --show-toplevel)" -n "$COUNT" --push --detach --url "$URL" --token "$TOKEN" >/dev/null 2>>"$LOG"

LOCAL_HOOK=$HOOK_DIR/post-commit.local
if [ -x "$LOCAL_HOOK" ]; then
  exec "$LOCAL_HOOK" "$@"
fi
