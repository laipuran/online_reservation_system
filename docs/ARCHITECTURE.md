# 后端架构指南

## 1. 架构总览

本项目采用简化的分层架构（借鉴 DDD 的 Clean Architecture 思想），分为三层：

```
HTTP Handlers (internal/api/http/handler/)
    │ 调用
    ▼
Service Layer (internal/service/)
    │ 调用接口（依赖倒置）
    ▼
Repository Layer (internal/repository/)
    ├── 接口定义 (interfaces.go)
    └── PostgreSQL 实现 (postgres/)
```

### 层间依赖规则

| 层 | 允许依赖 | 禁止依赖 |
|----|---------|---------|
| `api/handler/` | `service/`, `model/`, `api/http/response/` | 直接访问 `repository/`、`postgres/` |
| `service/` | `repository/`, `model/`, `pkg/` | 直接访问 `api/`、`postgres/` |
| `repository/` | `model/` | 直接访问 `api/`、`service/` |
| `repository/postgres/` | `repository/`（实现接口）、`model/` | 直接访问 `api/`、`service/` |

## 2. 目录结构与职责

```
ors-be/
├── cmd/server/main.go          # 入口：加载配置 → 初始化依赖 → 启动 HTTP
│
├── internal/
│   ├── api/http/               # HTTP 层（处理请求/响应）
│   │   ├── server.go           # 路由注册、服务器启动
│   │   ├── handler/*.go        # Handler 函数（不含业务逻辑）
│   │   ├── middleware/*.go     # 中间件（auth、cors、logger）
│   │   └── response/response.go # 统一响应 JSON 格式
│   │
│   ├── service/*.go            # 业务逻辑层（用例编排）
│   │
│   ├── repository/             # 数据访问层
│   │   ├── interfaces.go       # 所有 repository 接口定义
│   │   ├── transaction.go      # 事务包装器
│   │   └── postgres/*.go       # PostgreSQL 具体实现
│   │
│   ├── model/*.go              # 领域模型（纯 struct，无 ORM 标签）
│   │
│   ├── config/config.go        # 配置加载（环境变量）
│   │
│   ├── middleware/              # 非 HTTP 中间件（JWT 逻辑、密码哈希）
│   │   └── auth.go             # jwt 生成/验证、bcrypt 哈希
│   │
│   └── pkg/*.go                # 工具函数（验证器、响应工具）
│
├── migrations/*.sql            # 数据库迁移（up/down 成对）
│
├── go.mod / go.sum
└── Makefile                    # 构建/运行/迁移命令
```

## 3. Handler 规范

每个 handler 是一个接收 `service` 并返回 `http.HandlerFunc` 的闭包函数。

```go
// internal/api/http/handler/auth.go
func Register(userSvc service.UserService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. 反序列化请求体
        var req RegisterRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            response.JSON(w, http.StatusBadRequest, response.Fail("无效的请求体"))
            return
        }

        // 2. 校验入参
        // 3. 调用 service
        user, err := userSvc.Register(r.Context(), req.Email, req.Password, req.Name)
        // 4. 处理错误、返回响应
    }
}
```

### Handler 原则
- 不做业务逻辑判断（交给 service）
- 不做数据库操作（交给 repository）
- 只做：解析请求、校验基本格式、调用 service、序列化响应
- 统一通过 `response.JSON(w, code, data)` 输出

## 4. Service 规范

Service 封装用例逻辑，协调多个 repository 调用。

```go
// internal/service/user.go
type UserService interface {
    Register(ctx context.Context, email, password, name string) (*model.User, error)
    Login(ctx context.Context, email, password string) (*LoginResult, error)
    GetByID(ctx context.Context, id int64) (*model.User, error)
}

type userService struct {
    userRepo repository.UserRepository
    hasher   middleware.Hasher
    tokenGen middleware.TokenGenerator
}
```

### Service 原则
- 每个方法签名包含 `ctx context.Context` 作为第一个参数
- 返回 `(*T, error)`，错误已携带业务含义
- 跨多个 repository 的操作使用 `repository.WithTx()` 包裹事务
- 不做 HTTP 层面的处理（不读写 `http.Request`/`ResponseWriter`）

## 5. Repository 规范

### 接口定义

```go
// internal/repository/interfaces.go
type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    GetByEmail(ctx context.Context, email string) (*model.User, error)
    GetByID(ctx context.Context, id int64) (*model.User, error)
    Update(ctx context.Context, user *model.User) error
}
```

### 实现规范

```go
// internal/repository/postgres/user_repo.go
type userRepo struct {
    db *sqlx.DB  // 或 *pgxpool.Pool
}

func NewUserRepo(db *sqlx.DB) repository.UserRepository {
    return &userRepo{db: db}
}
```

### Repository 原则
- 方法签名全部调用 `ctx context.Context` 透传
- 接收和返回 `model.*` 领域对象
- SQL 使用**参数化查询**防止注入
- 不包含业务逻辑

## 6. Model 规范

纯数据结构，无 ORM 标签，无业务方法。

```go
// internal/model/user.go
type User struct {
    ID           int64     `json:"id"`
    Name         string    `json:"name"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`          // 永远不序列化
    Role         string    `json:"role"`
    Phone        string    `json:"phone,omitempty"`
    AvatarURL    string    `json:"avatar_url,omitempty"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

### Model 原则
- json tag 使用 snake_case
- 敏感字段（密码哈希等）使用 `json:"-"` 防止泄露
- 可选字段使用 `omitempty`

## 7. 响应格式

所有 API 统一响应格式：

```json
// 成功
{ "code": 200, "message": "ok", "data": { ... } }

// 失败
{ "code": 400, "message": "邮箱已注册", "data": null }
```

```go
// internal/api/http/response/response.go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func JSON(w http.ResponseWriter, httpStatus int, resp Response) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(httpStatus)
    json.NewEncoder(w).Encode(resp)
}

func Ok(data interface{}) Response {
    return Response{Code: 200, Message: "ok", Data: data}
}

func Fail(msg string) Response {
    return Response{Code: 400, Message: msg, Data: nil}
}

func Unauthorized(msg string) Response {
    return Response{Code: 401, Message: msg, Data: nil}
}

func Error(httpStatus int, msg string) Response {
    return Response{Code: httpStatus, Message: msg, Data: nil}
}
```

## 8. 错误处理

Service 层返回的 error 通过 errors.Is/errors.As 判断类型：

```go
// 定义错误哨兵
var (
    ErrEmailAlreadyRegistered = errors.New("邮箱已注册")
    ErrInvalidCredentials     = errors.New("邮箱或密码错误")
    ErrUserNotFound           = errors.New("用户不存在")
)
```

Handler 层统一处理：

```go
if errors.Is(err, service.ErrEmailAlreadyRegistered) {
    response.JSON(w, http.StatusConflict, response.Error(409, err.Error()))
    return
}
if errors.Is(err, service.ErrInvalidCredentials) {
    response.JSON(w, http.StatusUnauthorized, response.Unauthorized(err.Error()))
    return
}
```

## 9. 配置管理

使用环境变量（12-Factor App），不引入 viper。

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/ors?sslmode=disable` | 数据库连接 |
| `JWT_SECRET` | （必填） | JWT 签名密钥 |
| `JWT_EXPIRATION_HOURS` | `24` | Token 过期时间 |
| `HTTP_PORT` | `8080` | 监听端口 |
| `ALLOWED_ORIGINS` | `*` | CORS 允许的域名 |

```go
// internal/config/config.go
type Config struct {
    DatabaseURL        string
    JWTSecret          string
    JWTExpirationHours int
    HTTPPort           string
}

func Load() *Config {
    return &Config{
        DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ors?sslmode=disable"),
        JWTSecret:          getEnv("JWT_SECRET", ""),  // 生产环境必填
        JWTExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", 24),
        HTTPPort:           getEnv("HTTP_PORT", "8080"),
    }
}
```

## 10. 依赖注入

在 `main.go` 中手工装配依赖（不使用 DI 框架）：

```go
func main() {
    cfg := config.Load()

    db, err := postgres.Connect(cfg.DatabaseURL)
    // ...

    userRepo := postgres.NewUserRepo(db)
    userSvc := service.NewUserService(userRepo, hasher, tokenGen)

    router := chi.NewRouter()
    router.Route("/api/v1", func(r chi.Router) {
        handler.RegisterAuthRoutes(r, userSvc)
        // 后续添加更多路由...
    })

    http.ListenAndServe(":"+cfg.HTTPPort, router)
}
```

## 11. 测试约定

- 单元测试文件名 `*_test.go` 与被测文件同级
- 表驱动测试（Table Driven Tests）
- Repository 测试使用 `testcontainers-go` 启动真实 PostgreSQL
- 测试函数命名：`TestXxx_ShouldYyy_WhenZzz`

## 12. 命名规范

| 项 | 规范 | 示例 |
|----|------|------|
| 包名 | 全小写，单数 | `service`、`repository` |
| 文件 | snake_case | `user_repo.go`、`auth_handler.go` |
| 接口 | 方法名 + `er` 后缀 | `UserRepository`、`UserService` |
| 接口文件 | 以 `er` 结尾的单数 | `interfaces.go` |
| 结构体 | 小写开头（包内私有） | `userService`、`userRepo` |
| 构造函数 | `New` + 类型名 | `NewUserService`、`NewUserRepo` |
| 变量 | 驼峰 | `userSvc`、`userRepo` |
| JSON 字段 | snake_case | `avatar_url`、`password_hash` |

## 13. 迁移文件规范

- 使用原始 SQL，`up`/`down` 成对
- 文件名格式：`<序号>_<描述>.up.sql`
- 示例：`001_create_users.up.sql` / `001_create_users.down.sql`
- 业务表迁移序号跟随 PRD 表编号：PRD `3.6.x` 表使用三位前缀 `00x` / `0xx`，例如 `3.6.1 users -> 001_create_users`、`3.6.2 service_providers -> 002_create_service_providers`、`3.6.3 categories -> 003_create_categories`
- 通过 Makefile 的 `make migrate-up` / `make migrate-down` 执行
- 迁移工具：`golang-migrate/migrate`

## 14. 添加新功能的步骤

1. **Model**：在 `internal/model/` 定义数据结构
2. **Repository 接口**：在 `internal/repository/interfaces.go` 添加接口
3. **Repository 实现**：在 `internal/repository/postgres/` 实现
4. **Service**：在 `internal/service/` 编写用例逻辑
5. **Handler**：在 `internal/api/http/handler/` 编写 HTTP 处理
6. **路由注册**：在 `internal/api/http/server.go` 注册路由
7. **迁移**：在 `migrations/` 添加 SQL 迁移文件
8. **测试**：添加单元/集成测试

---

> 本文档是项目后端开发的规范依据，所有代码提交需遵循上述约定。
