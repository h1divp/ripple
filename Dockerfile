FROM golang:latest

WORKDIR /app

COPY --exclude=./web --exclude=.git go.mod ./

RUN go mod tidy

RUN go install github.com/air-verse/air@latest

COPY . .

ENTRYPOINT ["air"]
