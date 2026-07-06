#!/bin/sh
# git-ai-exporter post-commit hook
# Installed by `git-ai-exporter --install-hook`

URL=$(git config hooks.ai-exporter-url)
TOKEN=$(git config hooks.ai-exporter-token)

if [ -z "$URL" ] || [ -z "$TOKEN" ]; then
  echo ""
  echo "  git-ai-exporter: hooks not configured."
  echo "  Run the following commands to set up:"
  echo ""
  echo "    git config hooks.ai-exporter-url https://your-dashboard.com/api/collect"
  echo "    git config hooks.ai-exporter-token your-api-token"
  echo ""
  exit 0
fi

COUNT=$(git config hooks.ai-exporter-count)
COUNT=${COUNT:-1}

HOOK_DIR=$(dirname "$0")
LOG="$HOOK_DIR/git-ai-exporter.log"

git-ai-exporter -r "$(git rev-parse --show-toplevel)" -n "$COUNT" --push --detach --url "$URL" --token "$TOKEN" >/dev/null 2>>"$LOG"

LOCAL_HOOK=$HOOK_DIR/post-commit.local
if [ -x "$LOCAL_HOOK" ]; then
  exec "$LOCAL_HOOK" "$@"
fi
