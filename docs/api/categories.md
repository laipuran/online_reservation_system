# 服务分类模块

## GET /api/v1/categories

公开查询服务分类列表。分类用于服务搜索、浏览和服务发布时选择服务归属，字段对应 PRD 3.6.3 服务分类表。

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
      "name": "医疗",
      "description": "医疗健康服务",
      "created_at": "2026-07-07T09:00:00Z"
    },
    {
      "id": 2,
      "name": "口腔护理",
      "description": "口腔检查、洁牙等服务",
      "parent_id": 1,
      "created_at": "2026-07-07T09:10:00Z"
    }
  ]
}
```

### 响应字段

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 分类 ID |
| name | string | 是 | 分类名称，如医疗、美容、健身 |
| description | string | 否 | 分类描述 |
| parent_id | integer | 否 | 父分类 ID；为空表示顶级分类 |
| created_at | string | 是 | 创建时间，ISO 8601 格式 |

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 分类列表查询失败 | 500 | 分类列表查询失败 |

### 示例

```bash
curl -s http://localhost:8080/api/v1/categories
```
