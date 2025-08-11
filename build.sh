#!/bin/bash

# 获取当前脚本所在目录
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

# 构建前端
build_frontend() {
    echo "🔨 开始构建前端..."
    cd frontend || exit 1
    pnpm install
    pnpm run build
    echo "✅ 前端构建完成！"
}

# 构建后端
build_backend() {
    echo "🔨 开始构建后端..."
    cd "$SCRIPT_DIR" || exit 1
    # 用go命令构建到bin/目录下
    go build -ldflags -o bin/zmirror ./cmd/main.go
    upx -9 bin/zmirror
}

build_frontend && build_backend