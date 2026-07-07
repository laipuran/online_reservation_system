# API 接口文档

## 章节

| 章节 | 内容 |
|------|------|
| [认证模块](./auth.md) | 用户注册、登录 |

## 通用约定

- **Base URL**: `http://localhost:8080/api/v1`
- **请求体格式**: `Content-Type: application/json`
- **认证方式**: JWT Bearer Token（放在 `Authorization` 请求头）

### 统一响应格式

所有接口响应均为 JSON：

```json
{
  "code": 200,
  "message": "ok",
  "data": { ... }
}
```

### HTTP 状态码说明

| 状态码 | 含义 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 认证失败（未登录或 Token 无效） |
| 409 | 资源冲突（如邮箱已注册） |
| 500 | 服务器内部错误 |
