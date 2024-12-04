## build
FROM golang:1.23.3-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/main.go

# runtime
FROM alpine:3.20.3

WORKDIR /

COPY --from=builder /app/auth-service /auth-service

# use resource secret
COPY --from=builder /app/.env /.env

ENTRYPOINT ["/auth-service"]
