# Docker镜像代理系统 (ZMirror)

一个高性能的Docker Registry代理服务，支持多源镜像聚合、访问控制和WEB管理界面。

## 功能特性

- ✅ **完全兼容** Docker Registry V2 API
- 🚀 **多源代理** 支持多个上游镜像源，按优先级自动切换
- 🔐 **访问控制** 基于白名单的镜像访问控制
- 👤 **用户管理** 支持管理员和普通用户两种角色
- 🎯 **现代化WEB界面** Vue3 + Element Plus，支持响应式设计和管理员退出
- 📊 **访问日志** 完整的访问日志记录和查询
- ⚡ **高性能** 基于Gin框架，SQLite数据库，WAL模式
- 🐳 **容器化** 完整的Docker部署支持

## 快速开始

### 方式一：直接运行

```bash
# 克隆项目
git clone <repository-url>
cd zmirror

# 构建并运行
make run
```

### 方式二：Docker运行

```bash
# 使用Docker Compose
docker-compose up -d

# 或者直接使用Docker
docker run -d \
  -p 8080:8080 \
  -v ./data:/app/data \
  --name zmirror \
  zmirror:latest
```

### 首次启动

1. 服务启动后会自动创建 `./data` 目录
2. 生成默认配置文件 `./data/config.toml`
3. 创建SQLite数据库 `./data/registry.db`
4. 创建默认管理员账户：`admin/admin123`

## 配置说明

配置文件位置：`./data/config.toml`

```toml
[server]
host = "0.0.0.0"
port = "8080"

[admin]
username = "admin"
password = "admin123"

[database]
path = "./data/registry.db"
```

## 使用方式

### 1. 管理界面

访问 `http://localhost:8080` 进入WEB管理界面

**默认管理员账户：**
- 用户名：`admin`
- 密码：`admin123`

### 2. Docker客户端配置

```bash
# 配置Docker客户端使用代理
# 方式一：直接指定代理地址
docker pull localhost:8080/library/nginx:latest

# 方式二：配置daemon.json（推荐）
# 编辑 /etc/docker/daemon.json
{
  "registry-mirrors": [
    "http://localhost:8080"
  ]
}

# 重启Docker服务
sudo systemctl restart docker
```

### 3. 用户认证

```bash
# 登录到代理服务
docker login localhost:8080

# 输入普通用户的用户名和密码
# 注意：管理员账户不能用于docker login
```

## 访问控制机制

### 白名单规则

- **匹配白名单**：允许匿名拉取镜像
- **未匹配白名单**：需要用户认证后才能拉取

示例：
```
白名单前缀：helloz
匹配镜像：helloz/nginx, helloz/redis 等
效果：无需认证即可拉取
```

### 用户类型

**重要说明：管理员和普通用户的存储方式不同！**

#### 1. 管理员用户
- **存储位置**：配置文件 `data/config.toml`
- **用途**：只能通过WEB界面登录管理系统
- **权限**：可以管理镜像源、白名单、普通用户、查看日志
- **限制**：不能用于 `docker login` 认证
- **配置示例**：
  ```toml
  [admin]
  username = "admin"
  password = "your-secure-password"
  ```

#### 2. 普通用户
- **存储位置**：SQLite数据库 `data/registry.db`
- **用途**：只能用于 `docker login` 认证拉取镜像
- **权限**：只能拉取镜像，无法访问管理界面
- **创建方式**：通过WEB管理界面添加
- **密码存储**：MD5哈希加密存储

## CDN缓存配置

### 推荐缓存的API路径

以下路径建议配置CDN缓存：

#### 长期缓存 (1年)
```
/v2/*/blobs/*
```
这些是不可变的blob数据，可以长期缓存。

#### 短期缓存 (5分钟)
```
/v2/*/manifests/*
```
Manifest文件可能会更新，建议短期缓存。

### 不应缓存的路径
```
/v2/
/v2/*/tags/list
/api/*
/static/*
```

### Nginx配置示例

```nginx
server {
    listen 80;
    server_name your-registry.example.com;
    
    location / {
        proxy_pass http://zmirror:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # 长期缓存 blob 数据
    location ~ ^/v2/.*/blobs/ {
        proxy_pass http://zmirror:8080;
        proxy_cache_valid 200 1y;
        proxy_cache_key $uri;
        add_header X-Cache-Status $upstream_cache_status;
    }
    
    # 短期缓存 manifest 数据
    location ~ ^/v2/.*/manifests/ {
        proxy_pass http://zmirror:8080;
        proxy_cache_valid 200 5m;
        proxy_cache_key $uri;
        add_header X-Cache-Status $upstream_cache_status;
    }
}
```

## API文档

### Docker Registry V2 API

完全兼容Docker Registry V2 API规范：

- `GET /v2/` - API版本检查
- `GET /v2/{name}/tags/list` - 列出标签
- `GET /v2/{name}/manifests/{reference}` - 获取manifest
- `GET /v2/{name}/blobs/{digest}` - 获取blob数据

### 管理API

需要管理员认证，基于HTTP Basic Auth：

#### 镜像源管理
- `GET /api/registries` - 获取所有镜像源
- `POST /api/registries` - 创建镜像源
- `PUT /api/registries` - 更新镜像源
- `DELETE /api/registries/{id}` - 删除镜像源

#### 白名单管理
- `GET /api/whitelists` - 获取所有白名单
- `POST /api/whitelists` - 创建白名单
- `DELETE /api/whitelists/{id}` - 删除白名单

#### 用户管理
- `GET /api/users` - 获取所有用户
- `POST /api/users` - 创建用户
- `DELETE /api/users/{id}` - 删除用户

#### 访问日志
- `GET /api/logs` - 获取访问日志
- `DELETE /api/logs` - 清空访问日志

## 开发说明

### 目录结构

```
zmirror/
├── cmd/
│   └── main.go              # 应用入口
├── internal/
│   ├── model/               # 数据模型
│   ├── service/             # 业务逻辑
│   ├── handler/             # HTTP处理器
│   └── middleware/          # 中间件
├── web/
│   └── index.html           # WEB管理界面
├── docs/                    # 文档
├── data/                    # 数据目录（运行时生成）
│   ├── config.toml         # 配置文件
│   └── registry.db         # SQLite数据库
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

### 构建命令

```bash
# 安装依赖
make deps

# 开发模式运行
make dev

# 构建应用
make build

# 构建发布版本
make build-release

# 运行测试
make test

# 代码格式化
make fmt
```

### 数据库模型

主要数据表：

1. **users** - 用户表
2. **registries** - 镜像源表
3. **whitelists** - 白名单表
4. **access_logs** - 访问日志表

## 部署建议

### 生产环境配置

1. **修改默认密码**
   ```toml
   [admin]
   username = "admin"
   password = "your-secure-password"
   ```

2. **配置反向代理**
   - 使用Nginx或其他反向代理
   - 配置SSL/TLS证书
   - 启用缓存策略

3. **数据备份**
   - 定期备份 `./data` 目录
   - 特别是 `registry.db` 数据库文件

4. **监控和日志**
   - 配置容器日志收集
   - 监控服务健康状态
   - 设置告警机制

### Docker Compose配置

```yaml
version: '3.8'

services:
  zmirror:
    image: zmirror:latest
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data
    environment:
      - TZ=Asia/Shanghai
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8080/v2/"]
      interval: 30s
      timeout: 5s
      retries: 3
```

## 故障排除

### 常见问题

1. **无法连接上游镜像源**
   - 检查网络连接
   - 验证镜像源URL配置
   - 查看服务日志

2. **认证失败**
   - 确认用户名密码正确
   - 检查用户类型（管理员不能用于docker login）
   - 验证密码是否已修改

3. **缓存问题**
   - 清理Docker客户端缓存
   - 重启Docker daemon
   - 检查CDN缓存配置

### 日志查看

```bash
# Docker容器日志
docker logs zmirror

# 应用访问日志
# 可通过WEB界面查看，或直接查询数据库
```

## 许可证

本项目采用MIT许可证，详见LICENSE文件。

## 贡献

欢迎提交Issue和Pull Request来改进本项目。

## 更新日志

### v1.0.0
- 初始版本发布
- 支持Docker Registry V2 API
- 实现多源代理和访问控制
- 提供WEB管理界面
