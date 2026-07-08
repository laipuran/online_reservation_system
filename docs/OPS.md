# 操作手册

## 环境要求

| 工具 | 版本 | 说明 |
|------|------|------|
| Go | 1.24+ | `go version` |
| PostgreSQL（这是我们使用的数据库） | 16+ | `psql --version` |
| migrate CLI（这是我们用来维护数据库表结构的工具，意思就是在运行软件的时候不再需要手敲 SQL 建表，节省精力，降低容错） | latest | `migrate -version` |
| Node.js | 22+ | `node --version` |
| npm | 11+ | `npm --version` |

还有一些在这里使用到的工具，在此提供说明：
- turbo: 因为前端跑 npm run dev 和 后端跑 make run 都需要 cd 到对应的目录进行操作，不仅麻烦，有时候也意识不到自己错哪里了，所以在此引入 turbo，其功能很简单，就是在项目的根目录就能快速执行各个分 repo 的指令。

### 安装 migrate CLI

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

确保 `$GOPATH/bin` 在 PATH 中（参考下方环境配置）。

---

## 环境配置（这一段可以先忽略）

参考 `internal/config/config.go`，全部通过环境变量配置：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/ors?sslmode=disable` | 数据库连接串 |
| `JWT_SECRET` | `dev-secret-do-not-use-in-production` | JWT 签名密钥（生产环境务必修改） |
| `JWT_EXPIRATION_HOURS` | `24` | Token 过期小时数 |
| `HTTP_PORT` | `8080` | 监听端口 |
| `ALLOWED_ORIGINS` | `*` | CORS 允许域名 |

建议在项目根目录创建 `.env` 文件（不会被 git 跟踪）:

```bash
# .env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/ors?sslmode=disable
JWT_SECRET=my-secret-key
JWT_EXPIRATION_HOURS=24
HTTP_PORT=8080
```

然后启动时加载：

```bash
export $(grep -v '^#' .env | xargs)
```

---

## 快速启动

### 1. 启动 PostgreSQL

**方式一：本地安装**
```bash
sudo systemctl start postgresql
sudo -u postgres createdb ors
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'postgres';"
```

**方式二：Docker（推荐开发用）**
```bash
docker run -d \
  --name ors-pg \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=ors \
  -p 5432:5432 \
  postgres:16
```

### 2. 执行数据库迁移

```bash
cd ors-be
make migrate-up
```

验证表已创建：
```bash
psql "$DATABASE_URL" -c "\dt"
```

### 3. 启动服务

**方式一：仅启动后端**
```bash
make run
```

**方式二：启动前后端（推荐开发用）**
```bash
npm run dev
```

Turbo 会同时拉起后端（`make run`）和前端开发服务器（`react-router dev`）。

看到输出如下即成功：

```
ors-be:dev: 服务启动，监听 :8080
ors-fe:dev: ➜  Local:   http://localhost:5173/
```

> 如果报 `command not found: turbo`，在项目根目录执行 `npm install`。

### 完整一键启动（Docker + 迁移 + 服务）

```bash
docker start ors-pg 2>/dev/null || docker run -d --name ors-pg -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=ors -p 5432:5432 postgres:16
sleep 2
make -C ors-be migrate-up
make -C ors-be run
```

---

## Makefile 命令

| 命令 | 说明 |
|------|------|
| `make build` | 编译二进制到 `bin/ors-be` |
| `make run` | 编译并启动服务 |
| `make clean` | 删除 `bin/` 目录 |
| `make migrate-up` | 执行所有待处理的迁移 |
| `make migrate-down` | 回滚最近一次迁移 |
| `make migrate-create` | 交互式创建新的迁移文件 |
| `make deps` | 整理依赖（`go mod tidy`） |

---

## 项目结构

```
ors-be/
├── cmd/server/main.go              # 入口
├── internal/
│   ├── api/http/                   # HTTP 层
│   │   ├── server.go               # 路由注册
│   │   ├── handler/                # 请求处理器
│   │   ├── middleware/             # HTTP 中间件
│   │   └── response/               # 统一响应格式
│   ├── service/                    # 业务逻辑层
│   ├── repository/                 # 数据访问层
│   │   ├── interfaces.go           # 接口定义
│   │   ├── transaction.go          # 事务包装
│   │   └── postgres/               # PostgreSQL 实现
│   ├── model/                      # 领域模型
│   ├── config/                     # 配置加载
│   └── auth/                       # JWT + bcrypt
├── migrations/                     # SQL 迁移
├── Makefile
└── go.mod
```

---

## Turbo 命令

项目使用 Turborepo 编排前后端。所有命令在项目根目录执行：

| 命令 | 说明 |
|------|------|
| `npm run dev` | 并行启动前后端开发模式 |
| `npm run fe-only` | 仅启动前端，启用 MSW mock（无需后端），前端 API 请求由 Mock Service Worker 拦截并提供假数据 |
| `npm run build` | 构建所有 workspace |
| `npm run lint` | 对所有 workspace 执行 lint |
| `npm run test` | 对所有 workspace 执行测试 |
| `npm run typecheck` | 对前端执行类型检查 |
| `npm run clean` | 清理所有 workspace 的构建产物 |

启动后 Turbo 会监听文件变化，修改代码时自动重新编译/热更新。

### 单独启动某个 workspace

```bash
npm run dev -w ors-fe    # 仅启动前端
npm run dev -w ors-be    # 仅启动后端
```

### 安装依赖

始终从项目根目录安装：

```bash
npm install                      # 安装全部依赖
npm install <pkg> -w ors-fe      # 为前端安装依赖
npm install <pkg> -w ors-be      # 为后端安装依赖
```

---

## API 文档

按模块分章，见 [`api/`](./api/) 目录：

| 章节 | 文件 | 内容 |
|------|------|------|
| 通用约定 | [`api/README.md`](./api/README.md) | 响应格式、状态码说明 |
| 认证模块 | [`api/auth.md`](./api/auth.md) | 注册、登录、接口调试脚本 |
