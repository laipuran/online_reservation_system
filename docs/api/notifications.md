# 通知模块

> 本文档为 TASK-002 的 API 契约设计，接口运行时业务逻辑待后续落地。本轮仅保证契约覆盖 PRD 5.6，不表示这些接口已经注册或可调用。

## GET /api/v1/notifications

查询当前登录用户的站内通知列表。

### 认证

需要 JWT Bearer Token。

### 路径参数

无。

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| is_read | boolean | 否 | 是否按已读状态筛选，默认全部 |
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
        "id": 1,
        "user_id": 1,
        "title": "Reservation confirmed",
        "content": "Your reservation has been confirmed.",
        "type": "reservation_confirmed",
        "is_read": false,
        "created_at": "2026-07-07T09:00:00Z"
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
| is_read 参数非法 | 400 | 通知筛选参数无效 |

### 示例

```bash
curl -s "http://localhost:8080/api/v1/notifications?is_read=false&page=1&page_size=20" \
  -H "Authorization: Bearer $TOKEN"
```

---

## GET /api/v1/notifications/unread-count

查询当前登录用户未读通知数。

### 认证

需要 JWT Bearer Token。

### 路径参数

无。

### 查询参数

无。

### 请求体

无。

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "count": 3
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |

### 示例

```bash
curl -s http://localhost:8080/api/v1/notifications/unread-count \
  -H "Authorization: Bearer $TOKEN"
```

---

## PUT /api/v1/notifications/{id}/read

标记当前登录用户的一条通知为已读。

### 认证

需要 JWT Bearer Token。用户只能标记自己的通知。

### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 通知 ID |

### 查询参数

无。

### 请求体

无。

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "id": 1,
    "user_id": 1,
    "title": "Reservation confirmed",
    "content": "Your reservation has been confirmed.",
    "type": "reservation_confirmed",
    "is_read": true,
    "created_at": "2026-07-07T09:00:00Z"
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| ID 格式无效 | 400 | 无效的通知 ID |
| 通知不存在或不属于当前用户 | 404 | 通知不存在 |

### 示例

```bash
curl -s -X PUT http://localhost:8080/api/v1/notifications/1/read \
  -H "Authorization: Bearer $TOKEN"
```

---

## PUT /api/v1/notifications/read-all

标记当前登录用户的全部未读通知为已读。

### 认证

需要 JWT Bearer Token。

### 路径参数

无。

### 查询参数

无。

### 请求体

无。

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "updated_count": 3
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |

### 示例

```bash
curl -s -X PUT http://localhost:8080/api/v1/notifications/read-all \
  -H "Authorization: Bearer $TOKEN"
```
