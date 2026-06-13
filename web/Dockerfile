FROM node:25-alpine AS base
WORKDIR /app
RUN npm install -g pnpm
COPY package*.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY . .

FROM base AS development
EXPOSE 5173
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

FROM base AS builder
RUN pnpm run build

FROM node:25-alpine AS deployment
WORKDIR /app
COPY --from=builder /app/build ./build
COPY --from=builder /app/node_modules ./node_modules
RUN npm install -g @infisical/cli
CMD ["/bin/sh", "-c", "\
  INFISICAL_TOKEN=$(infisical login \
    --method=universal-auth \
    --client-id=\"$INFISICAL_CLIENT_ID\" \
    --client-secret=\"$INFISICAL_CLIENT_SECRET\" \
    --silent --plain) && \
  export INFISICAL_TOKEN && \
  exec infisical run \
    --projectId=\"$INFISICAL_PROJECT_ID\" \
    --env=\"$INFISICAL_ENV\" \
    --path=\"$INFISICAL_PATH\" \
    -- node build/index.js"]
