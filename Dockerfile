FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

ARG GITHUB_TOKEN

RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

COPY go.mod go.sum ./
COPY . .

RUN go mod download

RUN go build -o http-auth-example .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/http-auth-example .
COPY config.yaml ./config.yaml

EXPOSE 4222

CMD ["./http-auth-example", "serve", "--config", "config.yaml"]
