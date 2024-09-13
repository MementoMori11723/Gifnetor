# Use the official Go base image for building the app
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go app into the working directory
COPY app.go .

# Build the Go app with optimizations for a small binary
RUN go build -ldflags="-s -w" -o app app.go

# Create a new, smaller image for running the app
FROM alpine:latest

# Install required dependencies: ffmpeg and ca-certificates
RUN apk --no-cache add ca-certificates ffmpeg

# Set working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Create necessary directories
RUN mkdir -p ./uploads ./gifs

# Expose the port the app runs on
EXPOSE 8080

# Command to run the Go app
CMD ["./app"]

