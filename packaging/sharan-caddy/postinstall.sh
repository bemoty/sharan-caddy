#!/bin/sh
set -e

chown root:sharan-caddy /etc/sharan-caddy
chmod 0750 /etc/sharan-caddy
chown sharan-caddy:sharan-caddy /var/lib/sharan-caddy
chmod 0750 /var/lib/sharan-caddy

if command -v systemctl >/dev/null 2>&1; then
  systemctl daemon-reload >/dev/null 2>&1 || true
  NEXT="systemctl enable --now sharan-caddy"
else
  NEXT="start sharan-caddy using your init system"
fi

echo "sharan-caddy installed. Edit /etc/sharan-caddy/Caddyfile, then: $NEXT"
