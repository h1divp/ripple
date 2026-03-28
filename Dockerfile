FROM golang:1.26.1-alpine AS base

WORKDIR /app

COPY --exclude=./web --exclude=.git go.mod ./

RUN go mod tidy

FROM base AS development

RUN go install github.com/air-verse/air@latest

COPY . .

ENTRYPOINT ["air"]
