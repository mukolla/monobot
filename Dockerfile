FROM golang:1.20.4-alpine3.18 as build

COPY . /github.com/mukolla/monobot/
WORKDIR /github.com/mukolla/monobot/

RUN go mod download
RUN go mod tidy
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/mukolla/monobot/bin/bot .
COPY --from=0 /github.com/mukolla/monobot/configs configs/

EXPOSE 80

CMD ["./bot"]