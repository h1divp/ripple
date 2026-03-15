FROM golang:1.26.1-alpine AS development

WORKDIR /app

COPY --exclude=./web --exclude=.git go.mod ./

RUN go mod tidy

RUN go install github.com/air-verse/air@latest

COPY . .

ENTRYPOINT ["air"]
