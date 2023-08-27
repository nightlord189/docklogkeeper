FROM golang:1.21.0-alpine3.18

WORKDIR /app

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main

ENTRYPOINT ["/main"]