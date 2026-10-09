#!/bin/sh
set -eu

cat > /tmp/compose-smoke.yml <<'YAML'
services:
  hello:
    image: hello-world
YAML

docker compose -f /tmp/compose-smoke.yml run --rm hello
