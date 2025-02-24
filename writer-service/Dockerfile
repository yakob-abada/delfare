# 🏗️ Stage 1: Build the Golang Application
FROM golang:1.23 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum
COPY ./app/go.mod ./app/go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY ./app .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/writer main.go

# 🏁 Stage 2: Create a minimal production image
FROM alpine:latest

# Install necessary dependencies
RUN apk --no-cache add ca-certificates

# Set working directory inside the final container
WORKDIR /root

# Copy the built binary from the builder stage
COPY --from=builder /app/writer /root/writer

# Ensure the binary has execute permissions
RUN chmod +x /root/writer

# Expose API port
EXPOSE 8080

# Run the API
CMD ["/root/writer"]