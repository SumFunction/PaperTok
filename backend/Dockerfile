# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 设置 Go 代理为国内镜像（加速下载）
ENV GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,direct
ENV GO111MODULE=on

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Tidy and build the application
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/server .
COPY --from=builder /app/config.yaml .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./server"]
