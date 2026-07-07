# 认证模块

## POST /api/v1/auth/register

用户注册。

### 请求

```json
{
  "email": "zhangsan@example.com",
  "password": "Abcd1234",
  "name": "张三",
  "role": "customer"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| email | string | 是 | 邮箱，自动转小写，唯一约束 |
| password | string | 是 | 密码，长度至少 8 位 |
| name | string | 是 | 昵称 |
| role | string | 否 | 用户角色，支持 `customer` / `provider`，默认 `customer` |

### 成功响应 (201)

```json
{
  "code": 201,
  "message": "created",
  "data": {
    "user": {
      "id": 1,
      "name": "张三",
      "email": "zhangsan@example.com",
      "role": "customer",
      "created_at": "2026-07-07T09:00:00Z",
      "updated_at": "2026-07-07T09:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 邮箱格式无效 | 400 | 邮箱格式不正确 |
| 密码不足 8 位 | 400 | 密码长度至少8位 |
| 昵称为空 | 400 | 昵称不能为空 |
| 用户角色无效 | 400 | 用户角色不正确 |
| 邮箱已注册 | 409 | 邮箱已注册 |
| 请求体非法 | 400 | 无效的请求体 |

### 示例

```bash
curl -s http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"zhangsan@example.com","password":"Abcd1234","name":"张三","role":"customer"}'
```

注册服务提供者账号时，将 `role` 设为 `provider`。

---

## POST /api/v1/auth/login

用户登录。

### 请求

```json
{
  "email": "zhangsan@example.com",
  "password": "Abcd1234"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| email | string | 是 | 登录邮箱 |
| password | string | 是 | 密码 |

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "user": {
      "id": 1,
      "name": "张三",
      "email": "zhangsan@example.com",
      "role": "customer",
      "created_at": "2026-07-07T09:00:00Z",
      "updated_at": "2026-07-07T09:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 邮箱或密码错误 | 401 | 邮箱或密码错误 |
| 请求体非法 | 400 | 无效的请求体 |

### 示例

```bash
curl -s http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"zhangsan@example.com","password":"Abcd1234"}'
```
