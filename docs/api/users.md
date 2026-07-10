# 用户模块

## GET /api/v1/users/me

获取当前登录用户信息。

### 认证

需要 JWT Bearer Token。

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "id": 1,
    "name": "张三",
    "email": "zhangsan@example.com",
    "role": "customer",
    "phone": "13800000000",
    "avatar_url": "https://example.com/avatar.png",
    "created_at": "2026-07-07T09:00:00Z",
    "updated_at": "2026-07-07T09:00:00Z"
  }
}
```

---

## PUT /api/v1/users/me

修改当前登录用户资料。该接口不修改邮箱、角色和密码。

### 认证

需要 JWT Bearer Token。

### 请求

```json
{
  "name": "张三",
  "phone": "13800000000",
  "avatar_url": "https://example.com/avatar.png"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 昵称 |
| phone | string | 否 | 电话号码；如填写，只能包含数字、空格和 `-`，不能包含连续 `-`，且至少包含 1 个数字 |
| avatar_url | string | 否 | 头像 URL |

### 成功响应 (200)

返回修改后的用户信息。

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| 请求体非法 | 400 | 无效的请求体 |
| 昵称为空 | 400 | 昵称不能为空 |
| 电话号码格式非法 | 400 | 电话号码格式不正确 |
| 用户不存在 | 404 | 用户不存在 |

---

## PUT /api/v1/users/me/password

修改当前登录用户密码。

### 认证

需要 JWT Bearer Token。

### 请求

```json
{
  "current_password": "password123",
  "new_password": "newpass123"
}
```

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "message": "密码修改成功"
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| 请求体非法 | 400 | 无效的请求体 |
| 当前密码错误 | 401 | 当前密码错误 |
| 新密码不足 8 位 | 400 | 密码长度至少8位 |
| 用户不存在 | 404 | 用户不存在 |

---

## GET /api/v1/users/me/interests

获取当前登录用户的兴趣标签。

### 认证

需要 JWT Bearer Token。

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "name": "放松",
      "created_at": "2026-07-07T09:00:00Z"
    }
  ]
}
```

### 示例

```bash
curl -s http://localhost:8080/api/v1/users/me/interests \
  -H "Authorization: Bearer $TOKEN"
```

---

## PUT /api/v1/users/me/interests

替换当前登录用户的兴趣标签集合。

### 认证

需要 JWT Bearer Token。

### 请求

```json
{
  "tag_ids": [1, 2, 3]
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tag_ids | integer[] | 是 | 标签 ID 列表，重复 ID 会自动去重；空数组表示清空兴趣标签 |

### 成功响应 (200)

返回替换后的兴趣标签列表。

```json
{
  "code": 200,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "name": "放松",
      "created_at": "2026-07-07T09:00:00Z"
    }
  ]
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| 请求体非法 | 400 | 无效的请求体 |
| 标签 ID 非法 | 400 | 标签ID不正确 |
| 标签不存在 | 404 | 标签不存在 |

### 示例

```bash
curl -s -X PUT http://localhost:8080/api/v1/users/me/interests \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"tag_ids":[1,2,3]}'
```
