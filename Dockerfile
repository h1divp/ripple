# syntax=docker/dockerfile:1.7-labs
FROM golang:1.26.1-alpine AS base
WORKDIR /app
COPY --exclude=./web --exclude=.git go.mod ./
RUN go mod tidy

FROM base AS development
RUN go install github.com/air-verse/air@latest
COPY . .
ENTRYPOINT ["air"]

FROM base AS builder
COPY --exclude=./web --exclude=.git . .
RUN go build -o /app/server .

FROM alpine:3.21 AS deployment
WORKDIR /app
COPY --from=builder /app/server .
ENTRYPOINT ["./server"]
