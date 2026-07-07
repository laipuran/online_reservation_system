# 操作手册

## 环境要求

| 工具 | 版本 | 说明 |
|------|------|------|
| Go | 1.24+ | `go version` |
| PostgreSQL | 16+ | `psql --version` |
| migrate CLI | latest | `migrate -version` |

### 安装 migrate CLI

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

确保 `$GOPATH/bin` 在 PATH 中（参考下方环境配置）。

---

## 环境配置

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

```bash
make run
```

看到输出 `服务启动，监听 :8080` 即成功。

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

## API 文档

按模块分章，见 [`api/`](./api/) 目录：

| 章节 | 文件 | 内容 |
|------|------|------|
| 通用约定 | [`api/README.md`](./api/README.md) | 响应格式、状态码说明 |
| 认证模块 | [`api/auth.md`](./api/auth.md) | 注册、登录、接口调试脚本 |
