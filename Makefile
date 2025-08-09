.PHONY: build run clean deps test build-frontend dev-frontend

# 默认目标
all: build

# 安装依赖
deps:
	go mod tidy
	go mod download

# 构建前端
build-frontend:
	@echo "🔨 构建前端..."
	@if [ ! -d "frontend/node_modules" ]; then \
		echo "📦 安装前端依赖..."; \
		cd frontend && pnpm install; \
	fi
	@cd frontend && pnpm run build
	@echo "✅ 前端构建完成"

# 构建应用
build: deps build-frontend
	@echo "🔨 构建后端..."
	CGO_ENABLED=1 go build -o zmirror ./cmd/main.go
	@echo "✅ 构建完成"

# 构建发布版本
build-release: deps build-frontend
	@echo "🔨 构建发布版本..."
	CGO_ENABLED=1 go build -ldflags="-w -s" -o zmirror ./cmd/main.go
	@echo "✅ 发布版本构建完成"

# 运行应用
run: build
	./zmirror

# 开发模式运行后端
dev:
	go run ./cmd/main.go

# 开发模式运行前端
dev-frontend:
	@./dev-frontend.sh

# 清理构建产物
clean:
	rm -f zmirror
	rm -rf data/
	rm -rf web/assets/
	@if [ -d "frontend/node_modules" ]; then \
		echo "🧹 清理前端依赖..."; \
		rm -rf frontend/node_modules; \
	fi

# 运行测试
test:
	go test -v ./...

# Docker构建
docker-build:
	docker build -t zmirror:latest .

# 安装到系统
install: build-release
	sudo cp zmirror /usr/local/bin/
	sudo chmod +x /usr/local/bin/zmirror

# 创建systemd服务
systemd-install: install
	sudo cp scripts/zmirror.service /etc/systemd/system/
	sudo systemctl daemon-reload
	sudo systemctl enable zmirror

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

# 显示帮助
help:
	@echo "可用的Make目标:"
	@echo "  build          - 构建应用"
	@echo "  build-release  - 构建发布版本"
	@echo "  run            - 构建并运行应用"
	@echo "  dev            - 开发模式运行"
	@echo "  clean          - 清理构建产物"
	@echo "  test           - 运行测试"
	@echo "  docker-build   - Docker构建"
	@echo "  install        - 安装到系统"
	@echo "  systemd-install- 安装systemd服务"
	@echo "  fmt            - 格式化代码"
	@echo "  lint           - 代码检查"
	@echo "  deps           - 安装依赖"
	@echo "  help           - 显示此帮助"
