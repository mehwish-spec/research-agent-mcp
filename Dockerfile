# Stage 1: build the Go binary
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY . .
RUN go build -o mcp-server ./cmd/server

# Stage 2: minimal runtime image
FROM alpine:3.20

RUN apk add --no-cache ca-certificates docker-cli

WORKDIR /app
COPY --from=builder /app/mcp-server .
RUN mkdir -p /data

ENV ELASTICSEARCH_URL=http://elasticsearch:9200
ENV DATA_DIR=/data

ENTRYPOINT ["./mcp-server"]
