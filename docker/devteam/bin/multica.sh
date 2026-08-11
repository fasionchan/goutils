#!/bin/bash

# Environment variables:
# - MULTICA_SERVER_URL
# - MULTICA_APP_URL
# - MULTICA_WORKSPACE_ID
# - MULTICA_TOKEN
# - MULTICA_DAEMON_ID
# - MULTICA_DAEMON_DEVICE_NAME
# - MULTICA_AGENT_RUNTIME_NAME

set -euo pipefail

[ -n "${MULTICA_SERVER_URL:-}" ]   && multica config set server_url "$MULTICA_SERVER_URL"
[ -n "${MULTICA_APP_URL:-}" ]      && multica config set app_url "$MULTICA_APP_URL"
[ -n "${MULTICA_WORKSPACE_ID:-}" ] && multica config set workspace_id "$MULTICA_WORKSPACE_ID"
[ -n "${MULTICA_TOKEN:-}" ]        && multica login --token "$MULTICA_TOKEN"

[ ! -f ~/.ssh/id_ed25519 ] && ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -N "" && echo "SSH key generated:" && cat ~/.ssh/id_ed25519.pub && echo

workspace=$(pwd)

for repo in $GIT_REPOS; do
	path="${repo%%:*}"
	url="${repo#*:}"

	fullpath="$workspace/$path"
	echo "Cloning $url to $fullpath"
	makedir -p && cd "$fullpath" && git clone "$url"
done

if [ $# -gt 0 ]; then
	exec multica "$@"
else
	exec multica daemon start --foreground
fi
