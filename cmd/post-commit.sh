#!/bin/sh
# git-ai-exporter post-commit hook
# Installed by `git-ai-exporter --install-hook`

URL=$(git config hooks.ai-exporter-url)
TOKEN=$(git config hooks.ai-exporter-token)
COUNT=$(git config hooks.ai-exporter-count)

if [ -n "$URL" ] && [ -n "$TOKEN" ]; then
  COUNT=${COUNT:-1}
  ERR=$(git-ai-exporter -r "$(git rev-parse --show-toplevel)" -n "$COUNT" --push --url "$URL" --token "$TOKEN" 2>&1 >/dev/null)
  if [ $? -ne 0 ]; then
    echo ""
    echo "  git-ai-exporter: failed to push stats"
    echo "  $ERR"
    echo ""
  fi
fi

# Chain to preserved hook if exists
LOCAL_HOOK=$(dirname "$0")/post-commit.local
if [ -x "$LOCAL_HOOK" ]; then
  exec "$LOCAL_HOOK" "$@"
fi
