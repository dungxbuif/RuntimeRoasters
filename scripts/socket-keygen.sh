#!/usr/bin/env bash
set -euo pipefail

service="${1:-service}"
secret="$(openssl rand -hex 24)"
key="rrsock_${service}_${secret}"
hash="$(printf '%s' "$key" | shasum -a 256 | awk '{print $1}')"

printf 'SOCKET_API_KEY=%s\n' "$key"
printf 'SOCKET_INTERNAL_API_KEYS=%s:%s:events:publish\n' "$service" "$hash"
