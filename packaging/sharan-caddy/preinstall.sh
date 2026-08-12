#!/bin/sh
set -e

NOLOGIN=/usr/sbin/nologin
[ -x "$NOLOGIN" ] || NOLOGIN=/sbin/nologin
[ -x "$NOLOGIN" ] || NOLOGIN=/usr/bin/nologin
[ -x "$NOLOGIN" ] || NOLOGIN=/bin/false

if ! getent group sharan-caddy >/dev/null 2>&1; then
  if command -v groupadd >/dev/null 2>&1; then
    groupadd --system sharan-caddy
  else
    addgroup -S sharan-caddy
  fi
fi

if ! getent passwd sharan-caddy >/dev/null 2>&1; then
  if command -v useradd >/dev/null 2>&1; then
    useradd --system --no-create-home --shell "$NOLOGIN" --gid sharan-caddy --home-dir /var/lib/sharan-caddy sharan-caddy
  else
    adduser -S -D -H -G sharan-caddy -s "$NOLOGIN" -h /var/lib/sharan-caddy sharan-caddy
  fi
fi
