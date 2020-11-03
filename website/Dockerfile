FROM golang:1.14.8-alpine3.11 AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /api
WORKDIR /api

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ./app main.go

FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /api
WORKDIR /api
COPY --from=builder /api/app .
COPY wait-for-it.sh .
RUN apk add --no-cache bash
RUN chmod +x wait-for-it.sh
EXPOSE 8000

ENTRYPOINT ["./wait-for-it.sh", "db:3306", "--", "./wait-for-it.sh", "redis:6379", "--", "./app"]