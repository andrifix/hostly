FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

ARG GOPROXY

RUN go mod download

COPY . .

RUN mkdir ./bin

RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./hostify ./

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=builder /app/hostify /usr/local/bin/hostify
COPY --from=builder /app/Caddyfile /etc/caddy/Caddyfile
RUN mkdir "/var/log/caddy" && chown -R nobody:nobody /var/log/caddy

EXPOSE 80 443

USER nobody

HEALTHCHECK CMD wget --no-verbose --tries=1 --spider http://localhost/healthz || exit 1

CMD ["hostify", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]