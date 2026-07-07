# 服务模块

## GET /api/v1/services

搜索/浏览服务列表。

### 查询参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 否 | 关键词搜索（匹配标题和描述） |
| category_id | int | 否 | 按分类筛选 |
| provider_id | int | 否 | 按提供者筛选 |
| min_price | decimal | 否 | 最低价格 |
| max_price | decimal | 否 | 最高价格 |
| sort_by | string | 否 | 排序字段（price/rating/created_at），默认 created_at |
| sort_order | string | 否 | asc / desc，默认 desc |
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 20，最大 50 |

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "items": [
      {
        "id": 1,
        "title": "肩颈按摩",
        "description": "通过揉捏按压缓解肩颈肌肉酸痛。",
        "price": 199.00,
        "duration_minutes": 60,
        "avg_rating": 4.5,
        "provider": {
          "id": 1,
          "business_name": "舒心养生馆"
        },
        "category": {
          "id": 1,
          "name": "按摩"
        }
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 参数格式错误 | 400 | 请求参数错误 |

## GET /api/v1/services/:id

获取服务详情。

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "id": 1,
    "title": "肩颈按摩",
    "description": "通过揉捏按压缓解肩颈肌肉酸痛。",
    "price": 199.00,
    "duration_minutes": 60,
    "avg_rating": 4.5,
    "provider": {
      "id": 1,
      "business_name": "舒心养生馆"
    },
    "category": {
      "id": 1,
      "name": "按摩"
    }
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 服务不存在 | 404 | 服务不存在 |
