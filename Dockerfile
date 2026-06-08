# syntax=docker/dockerfile:1.7-labs
FROM golang:1.26.1-alpine AS base
WORKDIR /app
COPY --exclude=./web --exclude=.git go.mod go.sum ./
RUN go mod download

FROM base AS development
RUN go install github.com/air-verse/air@latest
COPY . .
ENTRYPOINT ["air"]

FROM base AS builder
COPY --exclude=./web --exclude=.git . .
RUN go build -o /app/server ./cmd/...

FROM alpine:3.21 AS deployment
WORKDIR /app
COPY --from=builder /app/server .
RUN apk add --no-cache curl && \
    curl -1sLf 'https://dl.cloudsmith.io/public/infisical/infisical-cli/setup.alpine.sh' | sh && \
    apk add infisical
ENTRYPOINT ["/bin/sh", "-c", "infisical run --projectId=\"$INFISICAL_PROJECT_ID\" -- ./server"]
