# 在线预约系统 代码扫描分析报告

> 生成日期：2026-07-09

---

## 一、共性问题排查

### 1.1 商家预约面板只显示数字ID

**文件**: `ors-fe/app/routes/provider/reservations/page.tsx:86-88`
**严重性**: 🔴 高

表格展示 `r.id` / `r.user_id` / `r.service_id` 三个纯数字字段，运营完全不可用。

**后端限制**: `GET /provider/reservations` 返回 `model.Reservation` 只有 `user_id` / `service_id`，无用户名称/服务标题。且后端无 `GET /users/{id}` 端点，无法反查用户姓名。

**修复方案**:
- 通过 `fetchServiceById` 获取每条预约的服务标题
- 用户名称暂时无法获取，显示 "用户 #{id}" 作为过渡

### 1.2 价格输入框上下箭头

**文件**: `ors-fe/app/routes/provider/services/new.tsx:122`
**严重性**: 🟡 中

`<input type="number">` 显示浏览器默认 spinner 箭头，对价格输入不友好（容易误触）。相同问题也存在于 `edit.tsx`。

**修复方案**: 添加 CSS 隐藏数字输入框箭头

### 1.3 注册服务端后跳转问题

**文件**: `ors-fe/app/routes/_layout/register.tsx:119-120`
**严重性**: 🔴 高

注册流程 step2 已提交了商家资料，却又跳转到 `/complete-profile` 要求再次填写，造成重复提交困惑。

同时 `login.tsx:28-29` — provider 登录后跳转到 `/provider/services` 而非 `/dashboard`，与 customer 登录逻辑不一致。

**修复方案**: 
- 注册完成后统一跳转到 `/dashboard`
- 登录后统一跳转到 `/dashboard`

### 1.4 商家界面未显示用户备注

**文件**: `ors-fe/app/routes/provider/reservations/page.tsx:81-88`
**严重性**: 🔴 高

预约表格列：`预约ID / 用户ID / 服务ID / 预约时间 / 状态 / 操作` — 缺少 `备注` 列。虽然 `ReservationItem` 类型有 `note` 字段，API 也返回，但 UI 未渲染。

**修复方案**: 新增"备注"列，显示 `r.note` 或 `-`

### 1.5 服务分类子分类

**文件**: 
- `ors-fe/app/routes/provider/services/new.tsx:84` — 分类下拉是扁平列表
- `ors-fe/app/routes/services/page.tsx:31` — `parent_id === 0` 硬编码
- `ors-fe/app/routes/services/page.tsx:25-32` — 试图用前端推导父子关系

**严重性**: 🔴 高

PRD 中 `categories` 表有 `parent_id` 支持树形层级，但前端多处硬编码 `parent_id === 0`（实际数据库用 `IS NULL` 表示一级分类）。

**修复方案**:
- `services/page.tsx`: 改用 `parent_id === null || parent_id === undefined` 判断顶级分类
- `new.tsx`: 分类选择时用 **缩进分组** 方式展示父子关系

### 1.6 预约时间未校验过去时间

**文件**: `ors-fe/app/routes/services/service-detail.tsx:178-179, 337-351`
**严重性**: 🟡 中

`minDate` 设为今日，但若用户选今日日期 + 已过去的时段，仍可提交。后端没有拒绝过去时间的校验。

**修复方案**: 在 `handleBooking` 中校验 `startTime` 是否在当前时间之后

### 1.7 黑白主题不完整

**文件**: 多个页面
**严重性**: 🟡 中

已实现 `theme-toggle.tsx` 和 `app.css` 的 dark 基础，但以下位置缺少 `dark:` 类：
- `service-detail.tsx` 中部分 `bg-white` / `border-gray-200` 等
- `complete-profile.tsx` 中部分样式
- `provider/services/edit.tsx` 中部分样式

**修复方案**: 逐页检查并补充 `dark:` 变体

### 1.8 搜索框无实际功能

**文件**: `ors-fe/app/routes/services/page.tsx:59-63`
**严重性**: 🔴 高

搜索 `<input>` 存在 UI 但没有绑定 `onChange` 查询逻辑。用户输入后不触发任何搜索行为。

同时 `service-detail.tsx:172-176` 的搜索框按 Enter 会跳转到 `/services?keyword=...`，但 `services/page.tsx` 不读取 URL 参数。

**修复方案**:
- 搜索框输入后调用 API 传 `keyword` 参数
- 读取 URL 中的 `keyword` 查询参数初始化搜索框

### 1.9 头像显示不正确

**文件**: 
- `ors-fe/app/routes/_layout.tsx:76` — 显示 `user.name.charAt(0)`
- `ors-fe/app/routes/dashboard.tsx:99-101` — 同上
- `ors-fe/app/routes/services/service-detail.tsx:58` — `provider.business_name.charAt(0)`

**严重性**: 🟡 中

`user.avatar_url` 和 `provider.logo_url` 两个字段从未在前端使用过。始终显示首字符占位。

**修复方案**: 头像区域优先使用 URL，无 URL 时再 fallback 到首字符

### 1.10 手机号没有校验

**文件**: `ors-fe/app/routes/_layout/register.tsx:33-40`
**严重性**: 🟡 中

`validateProviderFields` 校验了 `businessName` / `description` / `address` / `email`，**跳过了 `phone`**。同时注册后端的 `registerRequest` 也没有 `phone` 字段。

**修复方案**: 在注册表单增加手机号输入和格式校验（前端正则）

### 1.11 无客户确认服务完成按钮

**文件**: 前后端均缺失
**严重性**: 🔴 高

后端只有定时任务自动完成预约（`end_time <= now` 变 `completed`）。PRD 中 `Reservation` 类有 `complete()` 方法。缺少 `PUT /reservations/{id}/complete` 端点，需要后端配合。

**前端方案**: 在 `ReservationCard` 的 `confirmed` 状态增加"确认完成"按钮，调用 `/reservations/{id}/complete`（后期待实现）

### 1.12 用户预约后商家无站内信

**文件**: `ors-be/internal/service/reservation.go:74-109`
**严重性**: 🔴 高

`Create()` 方法创建预约后**没有发送通知给服务提供者**。只在 confirm / reject / cancel 时发通知。PRD 时序图 3.2.3 要求 `INSERT INTO notifications` 通知提供者。

**范围**: 此问题纯后端，前端无需操作

---

## 二、后端已实现但前端未调用的接口

| 接口 | 用途 | 未使用原因 |
|------|------|-----------|
| `GET /users/me` | 获取完整用户信息 | 前端只用了 auth store 的 User 对象 |
| `PUT /users/me` | 编辑个人资料 | 无"个人设置"页面 |
| `PUT /users/me/password` | 修改密码 | 无"修改密码"页面 |
| `GET /users/me/interests` | 查看已选兴趣标签 | 注册时只设置不查看 |
| `GET /users/me/reviews` | 我的评价历史 (PRD US-11) | 无"我的评价"页面 |
| `GET /reservations/{id}` | 单条预约详情 | 未实现详情页 |
| `POST /tags` | 创建标签 | 无标签管理页面 |
| `PUT /services/{id}/tags` | 替换服务标签 | 服务编辑表单未含标签管理 |
| `GET /tags/{id}` | 查看标签详情 | 低价值接口 |
| `GET /providers/{id}/reviews` | 查看商家评价列表 | 无页面入口 |

## 三、前端调用但后端不存在的接口

| 接口 | 前端位置 | 处理方式 |
|------|---------|---------|
| `GET /services/{id}/reviews/stats` | `ors-fe/app/lib/api/reviews.ts:33-37` | 从 reviews 列表计算，移除该 API 调用 |

## 四、需要后端配合才能修复的问题

| 问题 | 所需后端变更 | 说明 |
|------|------------|------|
| 1.1 商家面板显示用户名称 | 新增 `GET /users/{id}` 端点，或 provider 预约列表返回 `user_name` | 前端拿不到用户信息 |
| 1.11 客户确认服务完成 | 新增 `PUT /reservations/{id}/complete` 端点 | PRD 有 `complete()` 方法 |
| 1.12 预约创建通知 | `reservation.go Create()` 中增加 `notifyProviderNewReservation` 调用 | PRD 时序图明确要求 |
| review 列表缺 total | `newListResponse` 增加 `total` 字段 | 分页无法显示总页数 |
