# 服务提供者模块

## GET /api/v1/providers/{id}

公开查询服务提供者详情。

### 认证

不需要。

### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 服务提供者 ID |

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "id": 1,
    "user_id": 2,
    "business_name": "舒心养生馆",
    "description": "专业按摩服务",
    "address": "上海市",
    "phone": "13800000000",
    "email": "shop@example.com",
    "logo_url": "https://example.com/logo.png",
    "created_at": "2026-07-07T09:00:00Z",
    "updated_at": "2026-07-07T09:00:00Z"
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| ID 格式无效 | 400 | 无效的服务提供者ID |
| 服务提供者不存在 | 404 | 服务提供者不存在 |

### 示例

```bash
curl -s http://localhost:8080/api/v1/providers/1
```

---

## POST /api/v1/providers/me

为当前登录用户创建服务提供者资料。

### 认证

需要 JWT Bearer Token，且当前用户角色必须为 `provider`。

### 请求

```json
{
  "business_name": "舒心养生馆",
  "description": "专业按摩服务",
  "address": "上海市",
  "phone": "13800000000",
  "email": "shop@example.com",
  "logo_url": "https://example.com/logo.png"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| business_name | string | 是 | 商家名称 |
| description | string | 是 | 商家简介 |
| address | string | 是 | 地址 |
| phone | string | 否 | 联系电话；如填写，只能包含数字、空格和 `-`，不能包含连续 `-`，且至少包含 1 个数字 |
| email | string | 否 | 联系邮箱，保存前自动转小写；如填写必须符合邮箱格式 |
| logo_url | string | 否 | Logo URL |

### 成功响应 (201)

```json
{
  "code": 201,
  "message": "created",
  "data": {
    "id": 1,
    "user_id": 2,
    "business_name": "舒心养生馆",
    "description": "专业按摩服务",
    "address": "上海市",
    "phone": "13800000000",
    "email": "shop@example.com",
    "logo_url": "https://example.com/logo.png",
    "created_at": "2026-07-07T09:00:00Z",
    "updated_at": "2026-07-07T09:00:00Z"
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| 当前用户不是服务提供者 | 403 | 权限不足 |
| 请求体非法 | 400 | 无效的请求体 |
| 商家名称为空 | 400 | 商家名称不能为空 |
| 邮箱格式非法 | 400 | 邮箱格式不正确 |
| 电话号码格式非法 | 400 | 电话号码格式不正确 |
| 当前用户已有服务提供者资料 | 409 | 服务提供者资料已存在 |

### 示例

```bash
curl -s http://localhost:8080/api/v1/providers/me \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"business_name":"舒心养生馆","description":"专业按摩服务","address":"上海市","phone":"13800000000","email":"shop@example.com","logo_url":"https://example.com/logo.png"}'
```

---

## GET /api/v1/providers/me

查询当前登录用户的服务提供者资料。

### 认证

需要 JWT Bearer Token，且当前用户角色必须为 `provider`。

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "id": 1,
    "user_id": 2,
    "business_name": "舒心养生馆",
    "description": "专业按摩服务",
    "address": "上海市",
    "phone": "13800000000",
    "email": "shop@example.com",
    "logo_url": "https://example.com/logo.png",
    "created_at": "2026-07-07T09:00:00Z",
    "updated_at": "2026-07-07T09:00:00Z"
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| 当前用户不是服务提供者 | 403 | 权限不足 |
| 当前用户没有服务提供者资料 | 404 | 服务提供者不存在 |

### 示例

```bash
curl -s http://localhost:8080/api/v1/providers/me \
  -H "Authorization: Bearer $TOKEN"
```

---

## PUT /api/v1/providers/me

修改当前登录用户的服务提供者资料。

### 认证

需要 JWT Bearer Token，且当前用户角色必须为 `provider`。

### 请求

```json
{
  "business_name": "舒心养生馆",
  "description": "更新后的商家简介",
  "address": "上海市徐汇区",
  "phone": "13900000000",
  "email": "new-shop@example.com",
  "logo_url": "https://example.com/new-logo.png"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| business_name | string | 是 | 商家名称 |
| description | string | 否 | 商家简介 |
| address | string | 否 | 地址 |
| phone | string | 否 | 联系电话；如填写，只能包含数字、空格和 `-`，不能包含连续 `-`，且至少包含 1 个数字 |
| email | string | 否 | 联系邮箱，保存前自动转小写；如填写必须符合邮箱格式 |
| logo_url | string | 否 | Logo URL |

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "id": 1,
    "user_id": 2,
    "business_name": "舒心养生馆",
    "description": "更新后的商家简介",
    "address": "上海市徐汇区",
    "phone": "13900000000",
    "email": "new-shop@example.com",
    "logo_url": "https://example.com/new-logo.png",
    "created_at": "2026-07-07T09:00:00Z",
    "updated_at": "2026-07-07T10:00:00Z"
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| 当前用户不是服务提供者 | 403 | 权限不足 |
| 请求体非法 | 400 | 无效的请求体 |
| 商家名称为空 | 400 | 商家名称不能为空 |
| 邮箱格式非法 | 400 | 邮箱格式不正确 |
| 电话号码格式非法 | 400 | 电话号码格式不正确 |
| 当前用户没有服务提供者资料 | 404 | 服务提供者不存在 |

### 示例

```bash
curl -s -X PUT http://localhost:8080/api/v1/providers/me \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"business_name":"舒心养生馆","description":"更新后的商家简介","address":"上海市徐汇区","phone":"13900000000","email":"new-shop@example.com","logo_url":"https://example.com/new-logo.png"}'
```
