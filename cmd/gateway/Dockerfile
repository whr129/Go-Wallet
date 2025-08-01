# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum app.env ./
RUN go mod download

COPY . .
RUN go build -o gateway ./server/main.go

# Stage 2: Minimal runtime image
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/gateway .
COPY --from=builder /app/app.env .
# COPY --from=builder /app/db/migration ./db/migration

EXPOSE 8080
# Define entrypoint
ENTRYPOINT ["./gateway"]
