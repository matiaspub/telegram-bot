FROM golang:1.21.6-alpine3.19 AS builder

COPY . /telegram-bot
WORKDIR /telegram-bot

RUN go mod download
RUN go build -o ./.bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=0 /telegram-bot/.bin/bot .
COPY --from=0 /telegram-bot/config ./config/
COPY .env ./.env

EXPOSE 80

CMD ["./bot"]