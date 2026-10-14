#!/bin/sh
# Canary smoke: 100 запросов к canary-backend на порту 8444.
# Падает, если успешных меньше 99.
set -u

HOST="${1:-192.168.0.200}"
OK=0

i=1
while [ "$i" -le 100 ]; do
    if curl -fsSk "https://$HOST:8444/health" >/dev/null 2>&1; then
        OK=$((OK + 1))
    fi
    i=$((i + 1))
done

echo "canary OK: $OK/100"
test "$OK" -ge 99
