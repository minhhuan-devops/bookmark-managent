# syntax=docker/dockerfile:1
# ↑ Dòng này BẮT BUỘC phải ở dòng đầu tiên.
# Nó khai báo "frontend" mà BuildKit sử dụng để parse Dockerfile.
# Cần version 1 để hỗ trợ các tính năng nâng cao như --mount=type=cache.

# ============================================
# Stage 1: Build
# ============================================
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install swag for Swagger docs generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy dependency files first (leverage Docker layer caching)
COPY go.mod go.sum ./

# Download Go modules với cache mount
# --mount=type=cache,target=/go/pkg/mod
#   → Mount một thư mục cache vào /go/pkg/mod (nơi Go lưu modules đã download)
#   → Lần build đầu: Go download modules và lưu vào cache
#   → Lần build sau: Go thấy modules đã có trong cache → skip download
#   → Khác với CACHED layer: layer cache bị mất khi COPY go.mod thay đổi,
#     nhưng cache mount LUÔN tồn tại giữa các lần build
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

# Copy source code
COPY . .

# Generate Swagger docs
RUN swag init -g cmd/api/main.go

# Build the binary with cache mounts
# --mount=type=cache,target=/go/pkg/mod
#   → Giống trên: cho go build truy cập modules đã download
#
# --mount=type=cache,target=/root/.cache/go-build
#   → Mount cache vào Go build cache (nơi Go lưu kết quả compile)
#   → Lần build đầu: Go compile TẤT CẢ packages → lưu vào cache
#   → Lần build sau: Go chỉ compile lại packages CÓ THAY ĐỔI
#   → Ví dụ: bạn chỉ sửa handler.go → Go chỉ compile lại package handler,
#     các package khác (gin, redis, service...) dùng lại từ cache
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-s -w" \
  -o /app/server \
  ./cmd/api

# Compress binary with UPX
RUN apk add --no-cache upx && upx --best --lzma /app/server

# ============================================
# Stage 2: Runtime (scratch - smallest possible)
# ============================================
FROM alpine:3.22

# Copy the compiled binary
COPY --from=builder /app/server /server

# Expose application port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/server"]
