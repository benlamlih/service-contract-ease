# Stage 1: Build
FROM golang:alpine AS builder

WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy and build the application
COPY . .
RUN go build -o main ./cmd/api/main.go

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/main .
# Copy the config file
COPY config.yaml .

# Expose port if your server listens on one (optional)
EXPOSE 8080

CMD ["./main"]
