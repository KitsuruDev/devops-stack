#!/bin/sh
# Usage: smoke-color.sh <port>
set -eu
PORT="${1:?port required}"
BASE="http://127.0.0.1:${PORT}"

echo "=== Smoke on ${BASE} ==="
curl -fsS "${BASE}/health"
echo
curl -fsS -o /dev/null -w "SPA HTTP %{http_code}\n" "${BASE}/"
curl -fsS "${BASE}/api/notes" > /dev/null && echo "API OK"
