# Dockerfile.be
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build backend binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /be ./cmd/be

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /be /app/be

# Port is runtime-configured, so we don't hard-code EXPOSE here
ENTRYPOINT ["/app/be"]
