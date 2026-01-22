# Use official golang image as builder
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o signal-bot .

# Use alpine for the final image
FROM alpine:latest

# Install signal-cli dependencies (Java)
RUN apk add --no-cache \
    openjdk17-jre \
    wget \
    tar \
    bash

# Install signal-cli with verification
# Note: Checksum verification recommended - get actual checksum from signal-cli releases page
RUN wget https://github.com/AsamK/signal-cli/releases/download/v0.13.1/signal-cli-0.13.1-Linux.tar.gz \
    && tar xf signal-cli-0.13.1-Linux.tar.gz -C /opt \
    && ln -sf /opt/signal-cli-0.13.1/bin/signal-cli /usr/local/bin/ \
    && rm signal-cli-0.13.1-Linux.tar.gz

# Create directory for signal data
RUN mkdir -p /signal-data

# Copy the binary from builder
COPY --from=builder /app/signal-bot /usr/local/bin/signal-bot

# Set working directory
WORKDIR /app

# Create a non-root user
RUN addgroup -g 1000 signalbot && \
    adduser -D -u 1000 -G signalbot signalbot && \
    chown -R signalbot:signalbot /signal-data

# Switch to non-root user
USER signalbot

# Default command
CMD ["signal-bot"]
