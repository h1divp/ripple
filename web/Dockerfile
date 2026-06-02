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
CMD ["infisical", "run", "--", "node", "build/index.js"]
