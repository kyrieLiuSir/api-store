# 开发环境启动指南

## 环境要求

| 工具 | 版本要求 | 说明 |
|------|---------|------|
| Go | 1.22+ | 后端编译运行 |
| Bun | latest | 前端包管理和构建 |
| Node.js | 18+ | Vite 运行时需要 |
| SQLite / MySQL / PostgreSQL | 任选其一 | 数据库 |
| Redis | 可选 | 缓存，不配置时使用内存缓存 |

## 快速启动（本地开发）

### 1. 克隆项目并安装依赖

```bash
# 后端依赖
go mod download

# 前端依赖
cd web && bun install && cd ..
```

### 2. 配置环境变量

```bash
cp .env.example .env
```

编辑 `.env` 文件，最小配置（使用 SQLite，无需额外数据库）：

```env
# 不配置 SQL_DSN 时默认使用 SQLite，数据库文件�� ./new-api.db
# 如需 MySQL:
# SQL_DSN=root:password@tcp(127.0.0.1:3306)/new-api?parseTime=true
# 如需 PostgreSQL:
# SQL_DSN=postgresql://root:password@localhost:5432/new-api

# 可选：启用调试模式
DEBUG=true
```

> 默认情况下不配置任何数据库连接串，项目会自动使用 SQLite，零配置即可运行。

### 3. 启动后端

```bash
# 调试模式启动
GIN_MODE=debug go run .
ERROR_LOG_ENABLED=true SQL_DSN=postgresql://$(whoami)@localhost:5432/new-api?sslmode=disable GIN_MODE=debug go run .
```

后端默认监听 **http://localhost:3000**。

### 4. 启动前端开发服务器

```bash
cd web
bun run dev
```

前端 Vite 开发服务器默认监听 **http://localhost:5173**，已配置代理：
- `/api/*` → `http://localhost:3000`
- `/mj/*` → `http://localhost:3000`
- `/pg/*` → `http://localhost:3000`

**开发时访问 `http://localhost:5173` 即可，API 请求会自动代理到后端。**

### 5. 默认账号

首次启动后，系统会自动创建管理员账号：
- 用户名：`root`
- 密码：`123456`

> 请在首次登录后立即修改密码。

## Docker Compose 启动（完整环境）

如果你希望使用 Docker 一键启动完整环境（含 PostgreSQL + Redis）：

```bash
docker-compose up -d
```

访问 **http://localhost:3000**。

如需切换为 MySQL，参照 `docker-compose.yml` 中的注释修改即可。

## 前后端联合构建

生产构建流程（与 Dockerfile 一致）：

```bash
# 1. 构建前端
cd web
bun install
bun run build
cd ..

# 2. 构建后端（会嵌入 web/dist 到二进制中）
go build -ldflags "-s -w -X 'github.com/QuantumNous/new-api/common.Version=$(cat VERSION)'" -o new-api

# 3. 运行
./new-api
```

> Go 后端通过 `//go:embed web/dist` 将前端静态文件嵌入到二进制中，生产环境只需单个可执行文件。

## 常用环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | `3000` | 后端监听端口 |
| `GIN_MODE` | `release` | 设为 `debug` 开启调试模式 |
| `DEBUG` | `false` | 启用调试日志 |
| `SQL_DSN` | 空（SQLite） | 数据库连接串 |
| `SQLITE_PATH` | `./new-api.db` | SQLite 数据库路径 |
| `REDIS_CONN_STRING` | 空 | Redis 连接串，不配置则不启用 Redis |
| `MEMORY_CACHE_ENABLED` | `false` | 启用内存缓存 |
| `SESSION_SECRET` | 随机 | 会话密钥，多节点部署时必须统一设置 |
| `SYNC_FREQUENCY` | `60` | 数据库同步频率（秒） |
| `BATCH_UPDATE_ENABLED` | `false` | 启用批量更新 |
| `STREAMING_TIMEOUT` | `120` | 流式响应超时（秒） |
| `RELAY_TIMEOUT` | `0` | 请求超时（秒），0 表示不限制 |

完整环境变量列表参见 `.env.example`。

## 项目结构概览

```
├── main.go              # 入口文件
├── router/              # 路由定义
├── controller/          # 请求处理器
├── service/             # 业务逻辑
├── model/               # 数据模型 (GORM)
├── relay/               # AI API 代理转发
│   └── channel/         # 各厂商适配器 (openai/, claude/, gemini/ ...)
├── middleware/           # 中间件
├── setting/             # 配置管理
├── common/              # 公共工具
├── dto/                 # 数据传输对象
├── constant/            # 常量定义
├── i18n/                # 后端国际化
├── oauth/               # OAuth 实现
├── web/                 # React 前端
│   └── src/i18n/        # 前端国际化
├── .env.example         # 环境变量模板
├── docker-compose.yml   # Docker 编排
└── Dockerfile           # 容器构建
```

## 前端开发补充

```bash
cd web

# 开发服务器
bun run dev

# 代码格式化
bun run lint:fix

# ESLint 检查
bun run eslint

# i18n 相关
bun run i18n:extract   # 提取翻译 key
bun run i18n:sync      # 同步翻译文件
bun run i18n:lint      # 检查翻译完整性
```
