# Start the build stage with a full Go image
FROM golang:1.25.0-alpine3.22 AS builder
WORKDIR /src

RUN apk add --no-cache build-base git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# build the root/main (binary will call cmd.Execute())
RUN CGO_ENABLED=0 GOOS=linux go build -o sse-notifications .

# Final stage
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /src/sse-notifications .
# optional config path where we might mount config.yaml
VOLUME ["/etc/sse"]
EXPOSE 8080

# Default: run server subcommand. Can be overridden in docker-compose or docker run.
CMD ["./sse-notifications", "server"]
