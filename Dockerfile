# --- Stage 1: Build ---
FROM golang:1.26-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o portfolio-tui .

# --- Stage 2: Run ---
FROM alpine:latest

# Create a non-root user for security (The K8s way!)
RUN adduser -D vaishak
USER vaishak

WORKDIR /home/vaishak

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/portfolio-tui .

# Expose the SSH port we set in ssh.go
EXPOSE 2222

# Add this to your Dockerfile
ENV TERM=xterm-256color
ENV COLORTERM=truecolor

# Run the application
CMD ["./portfolio-tui"]