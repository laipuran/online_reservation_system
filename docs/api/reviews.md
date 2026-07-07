# 评价模块

> 本文档为 TASK-002 的 API 契约设计，接口运行时业务逻辑待后续落地。本轮仅保证契约覆盖 PRD 5.5，不表示这些接口已经注册或可调用。

## POST /api/v1/reviews

提交评价。后续实现时必须校验预约属于当前用户、预约状态为 `completed`，且一个预约只能评价一次。

### 认证

需要 JWT Bearer Token，普通用户可调用。

### 路径参数

无。

### 查询参数

无。

### 请求体

```json
{
  "reservation_id": 1001,
  "rating": 5,
  "comment": "Service was professional."
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| reservation_id | integer | 是 | 已完成预约 ID |
| rating | integer | 是 | 评分，范围 1 到 5 |
| comment | string | 否 | 评价内容 |

### 成功响应 (201)

```json
{
  "code": 201,
  "message": "created",
  "data": {
    "id": 501,
    "reservation_id": 1001,
    "user_id": 1,
    "service_id": 1,
    "rating": 5,
    "comment": "Service was professional.",
    "created_at": "2026-07-11T09:00:00Z"
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| 当前用户不是预约用户 | 403 | 权限不足 |
| 请求体非法或评分超出范围 | 400 | 评价参数无效 |
| 预约不存在 | 404 | 预约不存在 |
| 预约未完成 | 400 | 预约未完成，不能评价 |
| 预约已评价 | 409 | 该预约已评价 |

### 示例

```bash
curl -s -X POST http://localhost:8080/api/v1/reviews \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"reservation_id":1001,"rating":5,"comment":"Service was professional."}'
```

---

## GET /api/v1/services/{id}/reviews

公开查询某个服务的评价列表。

### 认证

不需要。

### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 服务 ID |

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | integer | 否 | 页码，默认 1 |
| page_size | integer | 否 | 每页数量，默认 20 |

### 请求体

无。

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "items": [
      {
        "id": 501,
        "reservation_id": 1001,
        "user_id": 1,
        "service_id": 1,
        "rating": 5,
        "comment": "Service was professional.",
        "created_at": "2026-07-11T09:00:00Z"
      }
    ],
    "page": 1,
    "page_size": 20
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| ID 格式无效 | 400 | 无效的服务 ID |
| 服务不存在 | 404 | 服务不存在 |

### 示例

```bash
curl -s "http://localhost:8080/api/v1/services/1/reviews?page=1&page_size=20"
```

---

## GET /api/v1/users/me/reviews

查询当前登录用户发表过的评价。

### 认证

需要 JWT Bearer Token，普通用户可调用。

### 路径参数

无。

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | integer | 否 | 页码，默认 1 |
| page_size | integer | 否 | 每页数量，默认 20 |

### 请求体

无。

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "items": [
      {
        "id": 501,
        "reservation_id": 1001,
        "user_id": 1,
        "service_id": 1,
        "rating": 5,
        "comment": "Service was professional.",
        "created_at": "2026-07-11T09:00:00Z"
      }
    ],
    "page": 1,
    "page_size": 20
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |

### 示例

```bash
curl -s "http://localhost:8080/api/v1/users/me/reviews?page=1&page_size=20" \
  -H "Authorization: Bearer $TOKEN"
```

---

## GET /api/v1/providers/{id}/reviews

公开查询某个服务提供者收到的评价列表。

### 认证

不需要。

### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 服务提供者 ID |

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | integer | 否 | 页码，默认 1 |
| page_size | integer | 否 | 每页数量，默认 20 |

### 请求体

无。

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "items": [
      {
        "id": 501,
        "reservation_id": 1001,
        "user_id": 1,
        "service_id": 1,
        "rating": 5,
        "comment": "Service was professional.",
        "created_at": "2026-07-11T09:00:00Z"
      }
    ],
    "page": 1,
    "page_size": 20
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| ID 格式无效 | 400 | 无效的服务提供者 ID |
| 服务提供者不存在 | 404 | 服务提供者不存在 |

### 示例

```bash
curl -s "http://localhost:8080/api/v1/providers/1/reviews?page=1&page_size=20"
```
