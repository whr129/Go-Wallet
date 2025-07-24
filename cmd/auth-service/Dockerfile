# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum auth.env ./
RUN go mod download

COPY . .
RUN go build -o auth-service ./server/main.go

# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/auth-service .
COPY --from=builder /app/auth.env .
COPY --from=builder /app/db/migration ./db/migration

EXPOSE 8080
# Define entrypoint
ENTRYPOINT ["./auth-service"]
