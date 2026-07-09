# 预约模块

## 状态流转

预约状态由后端维护，当前规则如下：

| 状态 | 说明 |
|------|------|
| `pending` | 用户创建预约后的初始状态，等待服务提供者确认 |
| `confirmed` | 服务提供者确认预约，用户可按时到店 |
| `completed` | 后台定时任务自动完成：仅当预约为 `confirmed` 且 `end_time <= now` 时更新 |
| `cancelled` | 用户取消预约，仅允许从 `pending` 或 `confirmed` 取消 |
| `rejected` | 服务提供者拒绝预约，仅允许从 `pending` 拒绝 |

后端服务启动后会每分钟扫描一次到期预约，并在启动时立即扫描一次。自动完成不会处理 `pending`、`cancelled` 或 `rejected` 状态的预约。用户创建预约、用户取消预约、预约自动完成时，会为对应服务提供者创建站内通知。

## POST /api/v1/reservations

创建预约。服务端根据服务时长计算 `end_time`，初始状态为 `pending`，并向对应服务提供者发送站内通知。

### 认证

需要 JWT Bearer Token，普通用户可调用。

### 路径参数

无。

### 查询参数

无。

### 请求体

```json
{
  "service_id": 1,
  "start_time": "2026-07-10T14:00:00Z",
  "note": "请准备热水"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| service_id | integer | 是 | 预约服务 ID |
| start_time | string | 是 | 预约开始时间，ISO 8601 格式 |
| note | string | 否 | 用户备注 |

### 成功响应 (201)

```json
{
  "code": 201,
  "message": "created",
  "data": {
    "id": 1001,
    "service": {
      "id": 1,
      "title": "肩颈按摩 60 分钟",
      "provider": {
        "id": 1,
        "business_name": "舒心养生馆"
      }
    },
    "start_time": "2026-07-10T14:00:00Z",
    "end_time": "2026-07-10T15:00:00Z",
    "status": "pending",
    "note": "请准备热水",
    "created_at": "2026-07-07T09:00:00Z"
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| 当前用户不是普通用户 | 403 | 权限不足 |
| 请求体非法或字段无效 | 400 | 无效的请求体 / 预约参数无效 |
| 服务不存在或不可预约 | 404 | 服务不存在 |
| 同一服务同一开始时间已被预约 | 409 | 该时段已被预约 |
| 服务端异常 | 500 | 预约操作失败 |

### 示例

```bash
curl -s -X POST http://localhost:8080/api/v1/reservations \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"service_id":1,"start_time":"2026-07-10T14:00:00Z","note":"请准备热水"}'
```

---

## GET /api/v1/reservations

查询当前登录用户的预约列表。

### 认证

需要 JWT Bearer Token，普通用户可调用。

### 路径参数

无。

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| status | string | 否 | 状态筛选：`pending` / `confirmed` / `completed` / `cancelled` / `rejected`，默认全部 |
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
        "id": 1001,
        "user_id": 1,
        "service_id": 1,
        "start_time": "2026-07-10T14:00:00Z",
        "end_time": "2026-07-10T15:00:00Z",
        "status": "pending",
        "created_at": "2026-07-07T09:00:00Z",
        "updated_at": "2026-07-07T09:00:00Z"
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
| 状态筛选非法 | 400 | 预约状态无效 |

### 示例

```bash
curl -s "http://localhost:8080/api/v1/reservations?status=pending&page=1&page_size=20" \
  -H "Authorization: Bearer $TOKEN"
```

---

## GET /api/v1/reservations/{id}

查询当前登录用户的一条预约详情。

### 认证

需要 JWT Bearer Token，普通用户只能查询自己的预约。

### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 预约 ID |

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
    "id": 1001,
    "user_id": 1,
    "service_id": 1,
    "start_time": "2026-07-10T14:00:00Z",
    "end_time": "2026-07-10T15:00:00Z",
    "status": "confirmed",
    "note": "请准备热水",
    "created_at": "2026-07-07T09:00:00Z",
    "updated_at": "2026-07-07T10:00:00Z"
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| ID 格式无效 | 400 | 无效的预约 ID |
| 查询他人预约 | 403 | 权限不足 |
| 预约不存在 | 404 | 预约不存在 |

### 示例

```bash
curl -s http://localhost:8080/api/v1/reservations/1001 \
  -H "Authorization: Bearer $TOKEN"
```

---

## PUT /api/v1/reservations/{id}/cancel

取消当前登录用户的预约。只允许取消 `pending` 或 `confirmed` 状态的预约，取消成功后会向对应服务提供者发送站内通知。

### 认证

需要 JWT Bearer Token，普通用户只能取消自己的预约。

### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 预约 ID |

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
    "id": 1001,
    "status": "cancelled",
    "updated_at": "2026-07-07T11:00:00Z"
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| ID 格式无效 | 400 | 无效的预约 ID |
| 查询或取消他人预约 | 403 | 权限不足 |
| 预约不存在 | 404 | 预约不存在 |
| 当前状态不可取消 | 400 | 当前预约状态不可取消 |

### 示例

```bash
curl -s -X PUT http://localhost:8080/api/v1/reservations/1001/cancel \
  -H "Authorization: Bearer $TOKEN"
```

---

## GET /api/v1/provider/reservations

服务提供者查询归属于自己服务的预约列表。

### 认证

需要 JWT Bearer Token，且当前用户角色必须为 `provider`。

### 路径参数

无。

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| status | string | 否 | 状态筛选：`pending` / `confirmed` / `completed` / `cancelled` / `rejected`，默认全部 |
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
        "id": 1001,
        "user_id": 1,
        "service_id": 1,
        "start_time": "2026-07-10T14:00:00Z",
        "end_time": "2026-07-10T15:00:00Z",
        "status": "pending",
        "created_at": "2026-07-07T09:00:00Z",
        "updated_at": "2026-07-07T09:00:00Z"
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
| 当前用户不是服务提供者 | 403 | 权限不足 |
| 状态筛选非法 | 400 | 预约状态无效 |

### 示例

```bash
curl -s "http://localhost:8080/api/v1/provider/reservations?status=pending" \
  -H "Authorization: Bearer $TOKEN"
```

---

## PUT /api/v1/provider/reservations/{id}/confirm

服务提供者确认待确认预约。预约确认后，当 `end_time <= now` 时会由后台定时任务自动更新为 `completed`，之后用户可以提交评价。

### 认证

需要 JWT Bearer Token，且当前用户角色必须为 `provider`。

### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 预约 ID |

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
    "id": 1001,
    "status": "confirmed",
    "updated_at": "2026-07-07T10:00:00Z"
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| 当前用户不是服务提供者 | 403 | 权限不足 |
| ID 格式无效 | 400 | 无效的预约 ID |
| 预约不存在或不属于该服务提供者 | 404 | 预约不存在 |
| 当前状态不可确认 | 400 | 当前预约状态不可确认 |

### 示例

```bash
curl -s -X PUT http://localhost:8080/api/v1/provider/reservations/1001/confirm \
  -H "Authorization: Bearer $TOKEN"
```

---

## PUT /api/v1/provider/reservations/{id}/reject

服务提供者拒绝待确认预约。

### 认证

需要 JWT Bearer Token，且当前用户角色必须为 `provider`。

### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 预约 ID |

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
    "id": 1001,
    "status": "rejected",
    "updated_at": "2026-07-07T10:00:00Z"
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| 当前用户不是服务提供者 | 403 | 权限不足 |
| ID 格式无效 | 400 | 无效的预约 ID |
| 预约不存在或不属于该服务提供者 | 404 | 预约不存在 |
| 当前状态不可拒绝 | 400 | 当前预约状态不可拒绝 |

### 示例

```bash
curl -s -X PUT http://localhost:8080/api/v1/provider/reservations/1001/reject \
  -H "Authorization: Bearer $TOKEN"
```
