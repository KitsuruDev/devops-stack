#!/bin/sh
# Usage: smoke-color.sh <port>
set -eu
PORT="${1:?port required}"
BASE="https://127.0.0.1:${PORT}"

echo "=== Smoke on ${BASE} (TLS self-signed) ==="
curl -fsSk "${BASE}/health"
echo
curl -fsSk -o /dev/null -w "SPA HTTP %{http_code}\n" "${BASE}/"
curl -fsSk "${BASE}/api/notes" > /dev/null && echo "API OK"
echo "SMOKE OK on ${BASE}"
