#!/bin/sh
# backend/healthcheck.sh
# Проверяет, что HTTP-эндпоинт /health возвращает 200 и в JSON поле "db": true.

set -eu

URL="http://127.0.0.1:8000/health"
TIMEOUT=2

# wget есть в alpine по умолчанию (busybox)
# -q  — тихий режим
# -O- — вывод в stdout
# -T  — таймаут в секундах
RESPONSE=$(wget -q -O- -T "$TIMEOUT" "$URL" 2>/dev/null) || {
    echo "healthcheck: HTTP request failed"
    exit 1
}

# Проверяем, что в ответе есть "db":true
echo "$RESPONSE" | grep -q '"db":true' || {
    echo "healthcheck: DB not ready, response=$RESPONSE"
    exit 1
}

exit 0
