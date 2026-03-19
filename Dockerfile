# Stage 1: Build
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Chỉ copy go.mod & go.sum trước để cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o api cmd/api/*.go

# Stage 2: Minimal runtime
FROM alpine:latest
WORKDIR /app

# Copy binary từ builder
COPY --from=builder /app/api ./

# Copy file cấu hình nếu có (env, config.yaml, v.v.)
# COPY config.yaml ./

# Chạy app
CMD ["./api", "server"]
