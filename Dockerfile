# 构建阶段
FROM golang:1.20-alpine AS builder

# 安装必要的系统依赖
RUN apk add --no-cache gcc musl-dev sqlite-dev

# 设置工作目录
WORKDIR /app

# 复制go模块文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o zmirror ./cmd/main.go

# 运行阶段
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates sqlite tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1001 zmirror && \
    adduser -u 1001 -G zmirror -s /bin/sh -D zmirror

# 设置工作目录
WORKDIR /app

# 复制构建的应用
COPY --from=builder /app/zmirror .
COPY --from=builder /app/web ./web

# 创建数据目录
RUN mkdir -p ./data && chown -R zmirror:zmirror /app

# 切换到非root用户
USER zmirror

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/v2/ || exit 1

# 启动应用
CMD ["./zmirror"]
