# 服务模块

## GET /api/v1/services

搜索/浏览已上架服务列表。

### 认证

不需要。

### Query Parameters

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 否 | 关键词搜索，匹配标题和描述 |
| category_id | integer | 否 | 按分类筛选 |
| provider_id | integer | 否 | 按服务提供者筛选 |
| min_price | number | 否 | 最低价格 |
| max_price | number | 否 | 最高价格 |
| sort_by | string | 否 | 排序字段：`price` / `rating` / `created_at`，默认 `created_at` |
| sort_order | string | 否 | `asc` / `desc`，默认 `desc` |
| page | integer | 否 | 页码，默认 1 |
| page_size | integer | 否 | 每页数量，默认 20，最大 50 |

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "items": [
      {
        "id": 1,
        "title": "肩颈按摩 60 分钟",
        "description": "缓解肩颈疲劳",
        "provider": {
          "id": 1,
          "business_name": "舒心养生馆"
        },
        "category": {
          "id": 1,
          "name": "美容"
        },
        "price": 199,
        "duration_minutes": 60,
        "image_url": "https://example.com/service.png",
        "status": "active",
        "avg_rating": 0,
        "review_count": 0,
        "created_at": "2026-07-07T09:00:00Z",
        "updated_at": "2026-07-07T09:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 20
  }
}
```

### 示例

```bash
curl -s "http://localhost:8080/api/v1/services?keyword=按摩&page=1&page_size=20"
```

---

## GET /api/v1/services/{id}

获取服务详情。

### 认证

不需要。

### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 服务 ID |

### 成功响应 (200)

```json
{
  "code": 200,
  "message": "ok",
  "data": {
    "id": 1,
    "title": "肩颈按摩 60 分钟",
    "description": "缓解肩颈疲劳",
    "provider": {
      "id": 1,
      "business_name": "舒心养生馆"
    },
    "category": {
      "id": 1,
      "name": "美容"
    },
    "price": 199,
    "duration_minutes": 60,
    "image_url": "https://example.com/service.png",
    "status": "active",
    "avg_rating": 0,
    "review_count": 0,
    "created_at": "2026-07-07T09:00:00Z",
    "updated_at": "2026-07-07T09:00:00Z"
  }
}
```

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| ID 格式无效 | 400 | 无效的服务ID |
| 服务不存在 | 404 | 服务不存在 |

---

## POST /api/v1/services

发布服务。

### 认证

需要 JWT Bearer Token，且当前用户角色必须为 `provider`。当前用户必须已经创建服务提供者资料。

### 请求

```json
{
  "category_id": 1,
  "title": "肩颈按摩 60 分钟",
  "description": "缓解肩颈疲劳",
  "price": 199,
  "duration_minutes": 60,
  "image_url": "https://example.com/service.png"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| category_id | integer | 是 | 分类 ID |
| title | string | 是 | 服务标题 |
| description | string | 否 | 服务详细描述 |
| price | number | 是 | 价格，不能小于 0 |
| duration_minutes | integer | 是 | 服务时长，必须大于 0 |
| image_url | string | 否 | 服务图片 URL |

### 成功响应 (201)

返回格式同服务详情。MVP 中服务发布后默认状态为 `active`。

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| 当前用户不是服务提供者 | 403 | 权限不足 |
| 未创建服务提供者资料 | 404 | 服务提供者不存在 |
| 请求体非法 | 400 | 无效的请求体 |
| 标题为空 | 400 | 服务标题不能为空 |
| 分类为空 | 400 | 服务分类不能为空 |
| 价格小于 0 | 400 | 服务价格不能小于0 |
| 服务时长小于等于 0 | 400 | 服务时长必须大于0 |

### 示例

```bash
curl -s http://localhost:8080/api/v1/services \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"category_id":1,"title":"肩颈按摩 60 分钟","description":"缓解肩颈疲劳","price":199,"duration_minutes":60,"image_url":"https://example.com/service.png"}'
```

---

## PUT /api/v1/services/{id}

编辑服务。

### 认证

需要 JWT Bearer Token，且当前用户角色必须为 `provider`。只能编辑自己发布的服务。

### 请求

```json
{
  "category_id": 1,
  "title": "肩颈按摩 90 分钟",
  "description": "更新后的服务说明",
  "price": 299,
  "duration_minutes": 90,
  "image_url": "https://example.com/service-new.png"
}
```

### 成功响应 (200)

返回格式同服务详情。

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| 当前用户不是服务提供者 | 403 | 权限不足 |
| 无权操作该服务 | 403 | 无权操作该服务 |
| 服务不存在 | 404 | 服务不存在 |
| 请求体非法 | 400 | 无效的请求体 |
| 参数非法 | 400 | 对应参数错误信息 |

---

## PATCH /api/v1/services/{id}/status

修改服务状态，当前支持上架/下架。

### 认证

需要 JWT Bearer Token，且当前用户角色必须为 `provider`。只能修改自己发布的服务。

### 请求

```json
{
  "status": "inactive"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| status | string | 是 | `active` / `inactive` |

### 成功响应 (200)

返回格式同服务详情。

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| 未登录或 Token 无效 | 401 | 缺少认证信息 / Token 无效或已过期 |
| 当前用户不是服务提供者 | 403 | 权限不足 |
| 无权操作该服务 | 403 | 无权操作该服务 |
| 服务不存在 | 404 | 服务不存在 |
| 状态非法 | 400 | 服务状态不正确 |

---

## GET /api/v1/services/{id}/tags

获取服务已关联的标签列表。

### 认证

不需要。

### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 服务 ID |

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

### 错误响应

| 场景 | HTTP 状态码 | message |
|------|-------------|---------|
| ID 格式无效 | 400 | 无效的服务ID |
| 服务不存在 | 404 | 服务不存在 |

---

## PUT /api/v1/services/{id}/tags

替换服务关联的标签集合。

### 认证

需要 JWT Bearer Token，且当前用户角色必须为 `provider`。只能修改自己发布的服务。

### 请求

```json
{
  "tag_ids": [1, 2, 3]
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| tag_ids | integer[] | 是 | 标签 ID 列表，重复 ID 会自动去重；空数组表示清空服务标签 |

### 成功响应 (200)

返回替换后的标签列表。

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
| 当前用户不是服务提供者 | 403 | 权限不足 |
| 无权操作该服务 | 403 | 无权操作该服务 |
| 服务不存在 | 404 | 服务不存在 |
| 标签不存在 | 404 | 标签不存在 |
| 标签 ID 非法 | 400 | 标签ID不正确 |

---

## GET /api/v1/providers/{id}/services

获取某个服务提供者的已上架服务列表。

### 认证

不需要。

### 路径参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | integer | 是 | 服务提供者 ID |

### Query Parameters

同 `GET /api/v1/services`，但 `provider_id` 会被路径参数覆盖。

### 成功响应 (200)

返回格式同服务列表。

### 示例

```bash
curl -s "http://localhost:8080/api/v1/providers/1/services?page=1&page_size=20"
```
