#!/bin/sh
set -eu

docker -v
docker compose version
systemctl is-active docker

