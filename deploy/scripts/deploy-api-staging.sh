set -eu
umask 077
ENV_DIR=/srv/api/env/api

aws s3 cp "s3://$DEPLOY_BUCKET/api/$IMAGE_TAG/bundle.tgz" /tmp/api-bundle.tgz
mkdir -p /srv/api
tar -xzf /tmp/api-bundle.tgz -C /srv/api

mkdir -p "$ENV_DIR"

export SOPS_AGE_KEY="$(
  aws ssm get-parameter \
    --name /ripple/deploy/sops-age-key \
    --with-decryption \
    --query 'Parameter.Value' \
    --output text
)"

sops --decrypt --input-type dotenv --output-type dotenv \
  "$ENV_DIR/staging.env.sops" > "$ENV_DIR/staging.env.tmp"

unset SOPS_AGE_KEY
chmod 600 "$ENV_DIR/staging.env.tmp"
mv "$ENV_DIR/staging.env.tmp" "$ENV_DIR/staging.env"

cd /srv/api
IMAGE_TAG_STAGING=$IMAGE_TAG
IMAGE_TAG_PROD=$(cat .production-tag 2>/dev/null || printf 'none')
export IMAGE_TAG_STAGING IMAGE_TAG_PROD

docker compose -f docker-compose.api.yml config --quiet
docker compose -f docker-compose.api.yml pull api-staging
docker compose -f docker-compose.api.yml up -d caddy redis-staging
docker compose -f docker-compose.api.yml up -d --no-deps --force-recreate api-staging

printf '%s\n' "$IMAGE_TAG" > .staging-tag
