# 标签模块

## GET /api/v1/tags

获取全部标签列表。

### 认证

不需要。

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
curl -s http://localhost:8080/api/v1/tags
```

---

## GET /api/v1/tags/{id}

获取标签详情。

### 认证

不需要。

### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 标签 ID |

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| ID 格式无效 | 400 | 无效的标签ID |
| 标签不存在 | 404 | 标签不存在 |

### 示例

```bash
curl -s http://localhost:8080/api/v1/tags/1
```

---

## POST /api/v1/tags

创建标签。

### 认证

需要 JWT Bearer Token，且当前用户角色必须为 `admin`。

### 请求

```json
{
  "name": "放松"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 标签名称，唯一，最长 50 个字符 |

### 成功响应 (201)

```json
{
  "code": 201,
  "message": "created",
  "data": {
    "id": 1,
    "name": "放松",
    "created_at": "2026-07-07T09:00:00Z"
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| 当前用户不是管理员 | 403 | 权限不足 |
| 请求体非法 | 400 | 无效的请求体 |
| 标签名称为空 | 400 | 标签名称不能为空 |
| 标签名称超过 50 个字符 | 400 | 标签名称不能超过50个字符 |
| 标签已存在 | 409 | 标签已存在 |

### 示例

```bash
curl -s http://localhost:8080/api/v1/tags \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"放松"}'
```
