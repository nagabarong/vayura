# syntax=docker/dockerfile:1

# Gunakan base image golang
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod dan go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build binary
RUN go build -o server ./cmd/server

# Stage runtime (lebih ringan)
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/server .
COPY Uploads ./Uploads

# Expose port default
EXPOSE 8080

# Jalankan server
CMD ["./server"]
