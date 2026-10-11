set -eu
umask 077
ENV_DIR=/srv/web/env/web

cd /srv/web
STAGING_TAG=$(cat .staging-tag)
TARGET_TAG=${IMAGE_TAG:-$STAGING_TAG}
: "${TARGET_TAG:?No staging or override image tag}"

aws s3 cp "s3://$DEPLOY_BUCKET/web/$TARGET_TAG/bundle.tgz" /tmp/web-bundle.tgz
mkdir -p /srv/web
tar -xzf /tmp/web-bundle.tgz -C /srv/web

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
IMAGE_TAG_STAGING=$STAGING_TAG
export IMAGE_TAG_STAGING IMAGE_TAG_PROD

docker compose -f docker-compose.web.yml config --quiet
docker compose -f docker-compose.web.yml pull web-prod
docker compose -f docker-compose.web.yml up -d --no-deps --force-recreate web-prod

# Make sure we're using the correct tag before setting .prod-tag. set -e will cause this to exit on failure.
test -n "$(docker compose -f docker-compose.web.yml ps -q web-prod)"

printf '%s\n' "$TARGET_TAG" > .prod-tag
