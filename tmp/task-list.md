# 前端待办清单

> 标记说明：
> - ✅ 前端可独立修复
> - 🔶 需要后端先完成，前端做预备
> - ❌ 纯后端

---

## Phase 1: 🔴 高优先级（可用性阻塞）

### 1. ✅ 商家预约面板显示服务名称和备注
- [x] 修改 `useProviderReservations` hook，拉取服务信息填充 title
- [x] 表格列改为：预约时间 / 服务名称 / 备注 / 状态 / 操作
- [x] 新增"备注"列显示 `r.note`

### 2. ✅ 搜索框连接查询
- [x] `services/page.tsx` 搜索框绑定状态，输入触发送 keyword 参数
- [x] 从 URL 读取 `keyword` 查询参数初始化搜索框

### 3. ✅ 服务分类子分类层级
- [x] 修复 `parent_id === 0` → 改用 `parent_id == null`
- [x] `new.tsx` 分类选择器用缩进展示父子层级

### 4. ✅ 注册跳转修复
- [x] provider 注册成功跳转到 `/dashboard` 而非 `/complete-profile`
- [x] `login.tsx` 登录后统一跳转到 `/dashboard`

### 5. ✅ review stats 不存在的接口处理
- [x] 删除 `fetchServiceReviewStats` 调用
- [x] 从 reviews 列表中计算评分分布
- [x] 处理 `total` 字段缺失的兼容

---

## Phase 2: 🟡 中优先级（体验优化）

### 6. ✅ 价格输入框隐藏上下箭头
- [x] 添加全局 CSS 隐藏 `input[type=number]` 的 spinner
- [x] 同时应用到 `new.tsx` 和 `edit.tsx`

### 7. ✅ 手机号校验
- [x] `register.tsx` provider 表单增加手机号格式校验
- [x] `complete-profile.tsx` 同步增加

### 8. ✅ 头像显示修复
- [x] layout / dashboard 使用 `user.avatar_url`
- [x] `ProviderCard` 使用 `provider.logo_url`
- [x] 无 URL 时 fallback 首字符

### 9. ✅ 黑白主题补全
- [x] 逐页补充缺失的 `dark:` 样式类

### 10. ✅ 预约时间校验
- [x] `service-detail.tsx handleBooking` 校验 startTime > now

---

## Phase 3: 🔶 需要后端配合

### 11. 🔶 商家面板显示用户姓名
前端已获取服务标题，但用户姓名需后端暴露 `GET /users/{id}` 或 provider 预约返回 `user_name`
- [x] 当前方案：显示"用户 #{id}"
- [ ] 后端完成后再切换到用户姓名

### 12. 🔶 客户确认服务完成按钮
- [x] `ReservationCard` 的 `confirmed` 状态增加"确认完成"按钮 UI
- [ ] 后端实现 `PUT /reservations/{id}/complete` 后可用

---

## 后端未实现的前端对接接口

### 需要前端去掉的调用
- [x] `GET /services/{id}/reviews/stats` — 改为从 reviews 列表计算

### 后端已实现但前端未用的接口（待后续功能开发）
| 接口 | 用途 |
|------|------|
| `PUT /users/me` | 个人资料编辑页 |
| `PUT /users/me/password` | 修改密码页 |
| `GET /users/me/reviews` | 我的评价历史页 |
| `GET /providers/{id}/reviews` | 商家评价列表页 |
| `PUT /services/{id}/tags` | 服务标签管理 |
