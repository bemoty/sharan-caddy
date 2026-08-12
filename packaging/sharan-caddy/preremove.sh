#!/bin/sh
set -e

case "$1" in
  upgrade|1|failed-upgrade)
    exit 0
    ;;
esac

if command -v systemctl >/dev/null 2>&1; then
  systemctl disable --now sharan-caddy >/dev/null 2>&1 || true
fi
