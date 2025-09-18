# light-stack

轻量级 Go 后台脚手架项目，支持多租户 SAAS 模式。

## 技术栈

- **后端**: Go 1.18+, Gin, GORM, MySQL, Redis
- **前端**: Vue3, TypeScript, Element Plus, Vite
- **认证**: JWT
- **缓存**: Redis
- **数据库**: MySQL 8.0+

## 快速开始

### 1. 环境要求

- Go 1.18+
- Node.js 16+
- MySQL 8.0+
- Redis 6.0+

### 2. 克隆项目

```bash
git clone <repository-url>
cd light-stack
```

### 3. 配置数据库

修改 `configs/config.yaml` 中的数据库配置：

```yaml
database:
  host: "localhost"
  port: "3306"
  username: "root"
  password: "your-password"
  database: "lightstack"
```

### 4. 初始化数据库

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE lightstack CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 运行数据库迁移
make migrate
```

### 5. 启动后端服务

```bash
# 安装依赖
go mod download

# 开发模式运行
make dev

# 或者构建后运行
make build
make run
```

后端服务将在 http://localhost:8080 启动

### 6. 启动前端服务

```bash
# 进入前端目录并安装依赖
cd web
npm install

# 启动开发服务器
npm run dev
```

前端服务将在 http://localhost:3000 启动

## 默认账号

- 用户名：admin
- 密码：admin123

## API 文档

启动服务后，可以通过以下接口测试：

- 健康检查: `GET /api/health`
- Ping 测试: `GET /api/v1/ping`

## 开发任务

查看 `doc/开发任务计划.md` 了解详细的开发计划和进度。

## 许可证

MIT License
