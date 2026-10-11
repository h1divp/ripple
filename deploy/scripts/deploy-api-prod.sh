set -eu
umask 077
ENV_DIR=/srv/api/env/api

cd /srv/api
STAGING_TAG=$(cat .staging-tag)
TARGET_TAG=${IMAGE_TAG:-$STAGING_TAG}
: "${TARGET_TAG:?No staging or override image tag}"

aws s3 cp "s3://$DEPLOY_BUCKET/api/$TARGET_TAG/bundle.tgz" /tmp/api-bundle.tgz
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
  "$ENV_DIR/prod.env.sops" > "$ENV_DIR/prod.env.tmp"

unset SOPS_AGE_KEY
chmod 600 "$ENV_DIR/prod.env.tmp"
mv "$ENV_DIR/prod.env.tmp" "$ENV_DIR/prod.env"

IMAGE_TAG_PROD=$TARGET_TAG
export IMAGE_TAG_STAGING IMAGE_TAG_PROD

docker compose -f docker-compose.api.yml config --quiet
docker compose -f docker-compose.api.yml pull api-prod
docker compose -f docker-compose.api.yml up -d caddy redis-prod
docker compose -f docker-compose.api.yml up -d --no-deps --force-recreate api-prod

printf '%s\n' "$TARGET_TAG" > .prod-tag
