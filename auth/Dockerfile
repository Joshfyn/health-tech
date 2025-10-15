FROM golang:1.25.1-alpine3.22 AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the auth service binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-service ./auth/main.go

# Final stage
FROM alpine:3.22.2

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/auth-service .

# Expose the port the service runs on
EXPOSE 80
