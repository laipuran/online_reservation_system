# 用户模块

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
