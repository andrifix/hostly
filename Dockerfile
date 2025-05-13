FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

ARG GOPROXY

RUN go mod download

COPY . .

RUN mkdir ./bin

RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./hostly ./

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=builder /app/hostly /usr/local/bin/hostly
COPY --from=builder /app/Caddyfile /etc/caddy/Caddyfile

EXPOSE 80 443 443/udp

ENV XDG_DATA_HOME=/data
ENV XDG_CONFIG_HOME=/config

HEALTHCHECK CMD wget --no-verbose --tries=1 --spider http://127.0.0.1/healthz || exit 1

CMD ["hostly", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]