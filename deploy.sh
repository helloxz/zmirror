#!/bin/bash

# 全栈部署脚本

set -e

echo "🚀 ZMirror Docker镜像代理管理系统部署脚本"
echo "================================================"

# 检查依赖
echo "🔍 检查系统依赖..."

# 检查Go
if ! command -v go &> /dev/null; then
    echo "❌ Go未安装，请安装Go 1.19或更高版本"
    exit 1
fi

# 检查pnpm
if ! command -v pnpm &> /dev/null; then
    echo "⚠️  pnpm未安装，尝试安装..."
    if command -v npm &> /dev/null; then
        npm install -g pnpm
    else
        echo "❌ npm未安装，请先安装Node.js和npm"
        exit 1
    fi
fi

echo "✅ 系统依赖检查完成"

# 构建项目
echo ""
echo "🔨 开始构建项目..."
make build-release

echo ""
echo "✅ 项目构建完成！"
echo ""
echo "📋 部署信息："
echo "  - 可执行文件: ./zmirror"
echo "  - Web界面: ./web/"
echo "  - 数据目录: ./data/ (首次运行时自动创建)"
echo ""
echo "🚀 启动服务："
echo "  ./zmirror"
echo ""
echo "🌐 访问地址："
echo "  - 管理界面: http://localhost:8080"
echo "  - API接口: http://localhost:8080/api"
echo "  - Docker代理: http://localhost:8080/v2"
echo ""
echo "🔐 默认账户："
echo "  - 用户名: admin"
echo "  - 密码: admin123"
echo ""
echo "💡 提示："
echo "  - 首次启动会自动创建配置文件和数据库"
echo "  - 生产环境请修改默认密码"
echo "  - 支持Docker Compose部署: docker-compose up -d"
