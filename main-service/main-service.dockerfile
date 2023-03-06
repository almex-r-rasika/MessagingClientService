# base go image
FROM golang:1.18-alpine as builder

RUN apk --update add \
    go \
    musl-dev
RUN apk --update add \
    util-linux-dev
RUN apk add --no-cache tzdata
RUN apk --update --no-cache add curl
RUN apk add --no-cache ca-certificates

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=1 go build -o mainApp ./cmd/api

RUN chmod +x /app/mainApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/mainApp /app

CMD [ "/app/mainApp" ]