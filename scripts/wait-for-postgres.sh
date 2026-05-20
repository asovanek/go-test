#!/usr/bin/env sh
set -e

COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.dev.yml}"
MAX_ATTEMPTS="${MAX_ATTEMPTS:-30}"

attempt=0
until docker compose -f "$COMPOSE_FILE" exec -T postgres pg_isready -U app -d app >/dev/null 2>&1; do
  attempt=$((attempt + 1))
  if [ "$attempt" -ge "$MAX_ATTEMPTS" ]; then
    echo "postgres did not become ready in time" >&2
    exit 1
  fi
  sleep 1
done

echo "postgres is ready"
