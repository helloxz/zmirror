#!/bin/bash

# 前端开发启动脚本

echo "🚀 启动Docker镜像代理管理系统前端开发环境"

# 检查pnpm是否安装
if ! command -v pnpm &> /dev/null; then
    echo "❌ pnpm未安装，请先安装pnpm:"
    echo "npm install -g pnpm"
    exit 1
fi

# 进入frontend目录
cd "$(dirname "$0")/frontend"

# 检查依赖是否安装
if [ ! -d "node_modules" ]; then
    echo "📦 安装依赖..."
    pnpm install
fi

echo "🔧 启动开发服务器..."
echo "📱 前端地址: http://localhost:8080"
echo "🔌 后端API: http://localhost:8090 (请确保后端服务已启动)"
echo ""
echo "💡 提示:"
echo "  - 修改代码会自动热重载"
echo "  - 按 Ctrl+C 停止服务"
echo ""

# 启动开发服务器
pnpm run dev
