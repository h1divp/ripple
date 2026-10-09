cd /srv/api

IMAGE_TAG_STAGING=$IMAGE_TAG
IMAGE_TAG_PROD=$(cat .production-tag 2>/dev/null || printf 'none')
export IMAGE_TAG_STAGING IMAGE_TAG_PROD

docker compose -f docker-compose.api.yml config --quiet
docker compose -f docker-compose.api.yml pull api-staging
docker compose -f docker-compose.api.yml up -d caddy redis-staging
docker compose -f docker-compose.api.yml up -d --no-deps --force-recreate api-staging

printf '%s\n' "$IMAGE_TAG" > .staging-tag
