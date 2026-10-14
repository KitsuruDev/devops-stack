#!/bin/sh
# Usage: switch.sh <blue|green>
set -eu

COLOR="${1:-}"
case "$COLOR" in
    blue)  PORT=9443 ;;
    green) PORT=9445 ;;
    *) echo "Usage: $0 <blue|green>"; exit 1 ;;
esac

sudo tee /etc/nginx/conf.d/gateway.conf > /dev/null <<EOF
server {
    listen 8443 ssl;
    server_name _;

    ssl_certificate     /opt/devops-stack/tls/server.crt;
    ssl_certificate_key /opt/devops-stack/tls/server.key;
    ssl_protocols       TLSv1.2 TLSv1.3;

    location / {
        proxy_pass https://127.0.0.1:${PORT};
        proxy_ssl_verify off;
        proxy_set_header Host              \$host;
        proxy_set_header X-Real-IP         \$remote_addr;
        proxy_set_header X-Forwarded-For   \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF

sudo nginx -t
sudo nginx -s reload
echo "==> Switched to $COLOR (port $PORT)"
