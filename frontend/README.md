# ZMirror Web 前端重构

## 项目简介

这是 ZMirror Docker镜像代理管理系统的全新前端界面，采用 Vue 3 + Element Plus + Vite 技术栈重构，提供现代化的用户体验。

## 技术栈

- **Vue 3** - 渐进式JavaScript框架
- **Element Plus** - Vue 3组件库
- **Vue Router 4** - 官方路由管理器
- **Pinia** - 状态管理
- **Axios** - HTTP客户端
- **Vite** - 前端构建工具

## 功能特性

### 🎨 现代化设计
- 精美的渐变色彩搭配
- 响应式布局，支持移动端
- 卡片式设计语言
- 流畅的动画效果

### 🔐 安全认证
- Basic认证支持
- 自动登录状态保持
- 一键退出登录功能

### 📊 仪表板
- 系统概览和统计信息
- 快速导航卡片
- 实时数据展示

### 🗂️ 管理功能
- **镜像源管理**: 添加、编辑、删除上游镜像源
- **白名单管理**: 配置允许访问的镜像前缀
- **用户管理**: 用户账户和密码管理
- **访问日志**: 详细的访问记录和筛选功能

### 🚀 开发体验
- 自动导入组件
- 热重载开发
- TypeScript支持（可选）
- 代码分割

## 快速开始

### 安装依赖

```bash
cd web-new
pnpm install
```

### 开发模式

```bash
pnpm run dev
```

访问 http://localhost:8080

### 构建生产版本

```bash
pnpm run build
```

构建结果将输出到 `../web` 目录，替换原有的静态文件。

### 预览生产版本

```bash
pnpm run preview
```

## 项目结构

```
web-new/
├── public/                 # 静态资源
├── src/
│   ├── components/         # 公共组件
│   ├── views/             # 页面组件
│   │   ├── Dashboard.vue  # 仪表板
│   │   ├── Registries.vue # 镜像源管理
│   │   ├── Whitelists.vue # 白名单管理
│   │   ├── Users.vue      # 用户管理
│   │   └── Logs.vue       # 访问日志
│   ├── stores/            # 状态管理
│   │   └── auth.js        # 认证状态
│   ├── router/            # 路由配置
│   │   └── index.js
│   ├── App.vue            # 根组件
│   ├── main.js           # 入口文件
│   └── style.css         # 全局样式
├── package.json
├── vite.config.js        # Vite配置
└── README.md
```

## API 接口

前端通过代理方式调用后端API：

- `GET /api/registries` - 获取镜像源列表
- `POST /api/registries` - 添加镜像源
- `PUT /api/registries` - 更新镜像源
- `DELETE /api/registries/:id` - 删除镜像源
- `GET /api/whitelists` - 获取白名单列表
- `POST /api/whitelists` - 添加白名单
- `DELETE /api/whitelists/:id` - 删除白名单
- `GET /api/users` - 获取用户列表
- `POST /api/users` - 添加用户
- `PUT /api/users/:id` - 更新用户
- `DELETE /api/users/:id` - 删除用户
- `GET /api/logs` - 获取访问日志

## 配置说明

### Vite 配置 (vite.config.js)

```javascript
export default defineConfig({
  server: {
    port: 8080,
    proxy: {
      '/api': {
        target: 'http://localhost:8090',  // 后端服务地址
        changeOrigin: true
      }
    }
  },
  build: {
    outDir: '../web'  // 构建输出目录
  }
})
```

### 自动导入配置

- Element Plus 组件自动导入
- Vue 3 API 自动导入
- 图标组件自动导入

## 部署说明

### 开发环境
1. 启动后端服务 (端口 8090)
2. 运行 `pnpm run dev` 启动前端开发服务器
3. 访问 http://localhost:8080

### 生产环境
1. 运行 `pnpm run build` 构建前端
2. 构建产物会自动输出到 `../web` 目录
3. 后端服务会直接提供这些静态文件

## 浏览器支持

- Chrome >= 87
- Firefox >= 78
- Safari >= 14
- Edge >= 88

## 更新记录

### v1.0.0 (2025-08-09)
- ✨ 全新的 Vue 3 + Element Plus 界面
- 🎨 现代化设计和响应式布局
- 🔐 完整的认证和权限管理
- 📊 仪表板和统计功能
- 🚀 更好的开发体验和性能

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

同主项目许可证
