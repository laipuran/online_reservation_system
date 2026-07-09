import { http, HttpResponse } from "msw";
import { db } from "./db";

const API = "/api/v1";

function json(data: unknown, code = 200, message = "ok") {
  return HttpResponse.json({ code, message, data }, { status: code });
}

function err(message: string, status: number) {
  return HttpResponse.json({ code: status, message, data: null }, { status });
}

/* ── helpers ─────────────────────────────────────────────── */

function pageParams(url: URL) {
  const page = Math.max(1, Number(url.searchParams.get("page")) || 1);
  const pageSize = Math.min(50, Number(url.searchParams.get("page_size")) || 20);
  return { page, pageSize, offset: (page - 1) * pageSize };
}

function toServiceRow(s: any) {
  const p = db.provider.findFirst({ where: { id: { equals: s.provider_id } } });
  const c = db.category.findFirst({ where: { id: { equals: s.category_id } } });
  return {
    id: s.id,
    title: s.title,
    description: s.description,
    price: s.price,
    duration_minutes: s.duration_minutes,
    avg_rating: s.avg_rating,
    review_count: s.review_count,
    status: s.status,
    image_url: s.image_url,
    provider: { id: p?.id ?? 0, business_name: p?.business_name ?? "" },
    category: { id: c?.id ?? 0, name: c?.name ?? "" },
    created_at: s.created_at,
    updated_at: s.updated_at,
  };
}

function getUserId(request: Request): number | null {
  const auth = request.headers.get("Authorization");
  if (!auth) return null;
  const match = auth.match(/^Bearer mock-token-(\d+)$/);
  return match ? Number(match[1]) : null;
}

/* ── Auth ─────────────────────────────────────────────────── */

export const handlers = [
  http.get(`${API}/auth/check-email`, ({ request }) => {
    const url = new URL(request.url);
    const email = url.searchParams.get("email");
    const exists = email ? !!db.user.findFirst({ where: { email: { equals: email } } }) : false;
    return json({ exists });
  }),

  http.post(`${API}/auth/register`, async ({ request }) => {
    const body: any = await request.json();
    const exists = db.user.findFirst({ where: { email: { equals: body.email } } });
    if (exists) return err("邮箱已注册", 409);
    const max = db.user.count();
    const user = db.user.create({
      id: max + 1,
      name: body.name,
      email: body.email,
      password: body.password,
      role: body.role ?? "customer",
      phone: "",
      avatar_url: "",
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    });
    return json(
      {
        user: { id: user.id, name: user.name, email: user.email, role: user.role, created_at: user.created_at, updated_at: user.updated_at },
        access_token: "mock-token-" + user.id,
      },
      201,
      "created"
    );
  }),

  http.post(`${API}/auth/login`, async ({ request }) => {
    const body: any = await request.json();
    const user = db.user.findFirst({ where: { email: { equals: body.email } } });
    if (!user || user.password !== body.password) return err("邮箱或密码错误", 401);
    return json({
      user: { id: user.id, name: user.name, email: user.email, role: user.role, created_at: user.created_at, updated_at: user.updated_at },
      access_token: "mock-token-" + user.id,
    });
  }),

  /* ── Categories ─────────────────────────────────────────── */

  http.get(`${API}/categories`, () => {
    const all = db.category.getAll();
    return json(all);
  }),

  /* ── Providers ──────────────────────────────────────────── */

  http.get(`${API}/providers/me`, ({ request }) => {
    const userId = getUserId(request);
    if (!userId) return err("缺少认证信息", 401);
    const p = db.provider.findFirst({ where: { user_id: { equals: userId } } });
    if (!p) return err("服务提供者不存在", 404);
    return json(p);
  }),

  http.post(`${API}/providers/me`, async ({ request }) => {
    const userId = getUserId(request);
    if (!userId) return err("缺少认证信息", 401);
    const user = db.user.findFirst({ where: { id: { equals: userId } } });
    if (!user || user.role !== "provider") return err("权限不足", 403);
    const exists = db.provider.findFirst({ where: { user_id: { equals: userId } } });
    if (exists) return err("服务提供者资料已存在", 409);
    const body: any = await request.json();
    const max = db.provider.count();
    const p = db.provider.create({
      id: max + 1,
      user_id: userId,
      business_name: body.business_name,
      description: body.description ?? "",
      address: body.address ?? "",
      phone: body.phone ?? "",
      email: body.email ?? "",
      logo_url: body.logo_url ?? "",
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    });
    try { localStorage.setItem("ors-provider", JSON.stringify(p)); } catch {}
    return json(p, 201, "created");
  }),

  http.put(`${API}/providers/me`, async ({ request }) => {
    const userId = getUserId(request);
    if (!userId) return err("缺少认证信息", 401);
    const p = db.provider.findFirst({ where: { user_id: { equals: userId } } });
    if (!p) return err("服务提供者不存在", 404);
    const body: any = await request.json();
    const updated = db.provider.update({
      where: { user_id: { equals: userId } },
      data: {
        business_name: body.business_name,
        description: body.description ?? "",
        address: body.address ?? "",
        phone: body.phone ?? "",
        email: body.email ?? "",
        logo_url: body.logo_url ?? "",
        updated_at: new Date().toISOString(),
      },
    });
    try { localStorage.setItem("ors-provider", JSON.stringify(updated)); } catch {}
    return json(updated);
  }),

  http.get(`${API}/providers/:id`, ({ params }) => {
    const id = Number(params.id);
    const p = db.provider.findFirst({ where: { id: { equals: id } } });
    if (!p) return err("服务提供者不存在", 404);
    return json(p);
  }),

  http.get(`${API}/providers/:id/services`, ({ params, request }) => {
    const providerId = Number(params.id);
    const url = new URL(request.url);
    const { page, pageSize, offset } = pageParams(url);
    const keyword = url.searchParams.get("keyword");
    const sortBy = url.searchParams.get("sort_by") || "created_at";
    const sortOrder = url.searchParams.get("sort_order") || "desc";

    let list = db.service.findMany({
      where: { provider_id: { equals: providerId } },
    });

    if (keyword) {
      const kw = keyword.toLowerCase();
      list = list.filter((s) => s.title.toLowerCase().includes(kw) || s.description.toLowerCase().includes(kw));
    }

    list.sort((a, b) => {
      let cmp = 0;
      if (sortBy === "price") cmp = a.price - b.price;
      else if (sortBy === "rating") cmp = a.avg_rating - b.avg_rating;
      else cmp = new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
      return sortOrder === "asc" ? cmp : -cmp;
    });

    const total = list.length;
    const items = list.slice(offset, offset + pageSize).map(toServiceRow);

    return json({ items, total, page, page_size: pageSize });
  }),

  /* ── Services ──────────────────────────────────────────── */

  http.get(`${API}/services`, ({ request }) => {
    const url = new URL(request.url);
    const { page, pageSize, offset } = pageParams(url);
    const keyword = url.searchParams.get("keyword");
    const categoryId = url.searchParams.get("category_id");
    const providerId = url.searchParams.get("provider_id");
    const minPrice = url.searchParams.get("min_price");
    const maxPrice = url.searchParams.get("max_price");
    const sortBy = url.searchParams.get("sort_by") || "created_at";
    const sortOrder = url.searchParams.get("sort_order") || "desc";

    let list = db.service.getAll();

    if (keyword) {
      const kw = keyword.toLowerCase();
      list = list.filter((s) => s.title.toLowerCase().includes(kw) || s.description.toLowerCase().includes(kw));
    }
    if (categoryId) list = list.filter((s) => s.category_id === Number(categoryId));
    if (providerId) list = list.filter((s) => s.provider_id === Number(providerId));
    if (minPrice) list = list.filter((s) => s.price >= Number(minPrice));
    if (maxPrice) list = list.filter((s) => s.price <= Number(maxPrice));

    list.sort((a, b) => {
      let cmp = 0;
      if (sortBy === "price") cmp = a.price - b.price;
      else if (sortBy === "rating") cmp = a.avg_rating - b.avg_rating;
      else cmp = new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
      return sortOrder === "asc" ? cmp : -cmp;
    });

    const total = list.length;
    const items = list.slice(offset, offset + pageSize).map(toServiceRow);

    return json({ items, total, page, page_size: pageSize });
  }),

  http.get(`${API}/services/:id`, ({ params }) => {
    const id = Number(params.id);
    const s = db.service.findFirst({ where: { id: { equals: id } } });
    if (!s) return err("服务不存在", 404);
    return json(toServiceRow(s));
  }),

  http.post(`${API}/services`, async ({ request }) => {
    const body: any = await request.json();
    const max = db.service.count();
    const now = new Date().toISOString();
    const s = db.service.create({
      id: max + 1,
      title: body.title,
      description: body.description ?? "",
      price: body.price,
      duration_minutes: body.duration_minutes,
      avg_rating: 0,
      review_count: 0,
      status: "active",
      image_url: body.image_url ?? "",
      provider_id: 1,
      category_id: body.category_id,
      created_at: now,
      updated_at: now,
    });
    return json(toServiceRow(s), 201, "created");
  }),

  http.put(`${API}/services/:id`, async ({ params, request }) => {
    const id = Number(params.id);
    const body: any = await request.json();
    const s = db.service.findFirst({ where: { id: { equals: id } } });
    if (!s) return err("服务不存在", 404);
    const updated = db.service.update({
      where: { id: { equals: id } },
      data: {
        title: body.title,
        description: body.description ?? "",
        price: body.price,
        duration_minutes: body.duration_minutes,
        image_url: body.image_url ?? "",
        category_id: body.category_id,
        updated_at: new Date().toISOString(),
      },
    });
    return json(toServiceRow(updated));
  }),

  http.patch(`${API}/services/:id/status`, async ({ params, request }) => {
    const id = Number(params.id);
    const body: any = await request.json();
    const s = db.service.findFirst({ where: { id: { equals: id } } });
    if (!s) return err("服务不存在", 404);
    const updated = db.service.update({
      where: { id: { equals: id } },
      data: { status: body.status, updated_at: new Date().toISOString() },
    });
    return json(toServiceRow(updated));
  }),

  http.get(`${API}/services/:id/tags`, ({ params }) => {
    const id = Number(params.id);
    const s = db.service.findFirst({ where: { id: { equals: id } } });
    if (!s) return err("服务不存在", 404);
    const rels = db.service_tag.findMany({ where: { service_id: { equals: id } } });
    const tags = rels.map((r) => {
      const t = db.tag.findFirst({ where: { id: { equals: r.tag_id } } });
      return t ? { id: t.id, name: t.name, created_at: t.created_at } : null;
    }).filter(Boolean);
    return json(tags);
  }),

  http.put(`${API}/services/:id/tags`, async ({ params, request }) => {
    const id = Number(params.id);
    const body: any = await request.json();
    const s = db.service.findFirst({ where: { id: { equals: id } } });
    if (!s) return err("服务不存在", 404);
    db.service_tag.deleteMany({ where: { service_id: { equals: id } } });
    let max = db.service_tag.count();
    for (const tagId of body.tag_ids) {
      max++;
      db.service_tag.create({ id: max, service_id: id, tag_id: tagId });
    }
    const rels = db.service_tag.findMany({ where: { service_id: { equals: id } } });
    const tags = rels.map((r) => {
      const t = db.tag.findFirst({ where: { id: { equals: r.tag_id } } });
      return t ? { id: t.id, name: t.name, created_at: t.created_at } : null;
    }).filter(Boolean);
    return json(tags);
  }),

  /* ── Tags ───────────────────────────────────────────────── */

  http.get(`${API}/tags`, () => {
    return json(db.tag.getAll());
  }),

  http.get(`${API}/tags/:id`, ({ params }) => {
    const id = Number(params.id);
    const t = db.tag.findFirst({ where: { id: { equals: id } } });
    if (!t) return err("标签不存在", 404);
    return json(t);
  }),

  http.post(`${API}/tags`, async ({ request }) => {
    const body: any = await request.json();
    const exists = db.tag.findFirst({ where: { name: { equals: body.name } } });
    if (exists) return err("标签已存在", 409);
    const max = db.tag.count();
    const t = db.tag.create({
      id: max + 1,
      name: body.name,
      created_at: new Date().toISOString(),
    });
    return json(t, 201, "created");
  }),

  /* ── Reviews ─────────────────────────────────────────────── */

  http.post(`${API}/reviews`, async ({ request }) => {
    const userId = getUserId(request);
    if (!userId) return err("缺少认证信息", 401);
    const body: any = await request.json();
    const { reservation_id, rating, comment } = body;
    if (!reservation_id || !rating) return err("参数不完整", 400);
    if (rating < 1 || rating > 5) return err("评分必须在 1-5 之间", 400);

    const res = db.reservation.findFirst({ where: { id: { equals: reservation_id } } });
    if (!res) return err("预约不存在", 404);
    if (res.user_id !== userId) return err("无权评价此预约", 403);
    if (res.status !== "completed") return err("仅已完成预约可评价", 400);

    const existingReview = db.review.findFirst({ where: { reservation_id: { equals: reservation_id } } });
    if (existingReview) return err("已评价过此预约", 409);

    const maxReview = db.review.count();
    const now = new Date().toISOString();
    const review = db.review.create({
      id: maxReview + 1,
      reservation_id,
      user_id: userId,
      service_id: res.service_id,
      rating,
      comment,
      created_at: now,
    });

    const svc = db.service.findFirst({ where: { id: { equals: res.service_id } } });
    if (svc) {
      const allReviews = db.review.findMany({ where: { service_id: { equals: svc.id } } });
      const totalReviews = allReviews.length;
      const avg = allReviews.reduce((sum, r) => sum + r.rating, 0) / totalReviews;
      db.service.update({
        where: { id: { equals: svc.id } },
        data: { avg_rating: Math.round(avg * 10) / 10, review_count: totalReviews, updated_at: now },
      });
    }

    const p = svc ? db.provider.findFirst({ where: { id: { equals: svc.provider_id } } }) : null;
    const notifMax = db.notification.count();
    if (p) {
      db.notification.create({
        id: notifMax + 1,
        user_id: p.user_id,
        title: "收到新评价",
        content: `您的服务「${svc?.title ?? ""}」收到了一条 ${rating}⭐ 评价。`,
        type: "system",
        is_read: false,
        created_at: now,
      });
    }

    return json(review, 201, "created");
  }),

  http.get(`${API}/services/:id/reviews`, ({ params, request }) => {
    const serviceId = Number(params.id);
    const s = db.service.findFirst({ where: { id: { equals: serviceId } } });
    if (!s) return err("服务不存在", 404);
    const url = new URL(request.url);
    const { page, pageSize, offset } = pageParams(url);

    let list = db.review.findMany({ where: { service_id: { equals: serviceId } } });
    list.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
    const items = list.slice(offset, offset + pageSize);
    return json({ items, page, page_size: pageSize });
  }),

  /* ── Reservations (provider) ────────────────────────────── */

  http.get(`${API}/provider/reservations`, ({ request }) => {
    const url = new URL(request.url);
    const { page, pageSize, offset } = pageParams(url);
    const status = url.searchParams.get("status");

    const providerServices = db.service.findMany({ where: { provider_id: { equals: 1 } } });
    const serviceIds = providerServices.map((s) => s.id);

    let list = db.reservation
      .getAll()
      .filter((r) => serviceIds.includes(r.service_id));

    if (status) list = list.filter((r) => r.status === status);

    list.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
    const items = list.slice(offset, offset + pageSize);
    return json({ items, page, page_size: pageSize });
  }),

  http.put(`${API}/provider/reservations/:id/confirm`, ({ params }) => {
    const id = Number(params.id);
    const r = db.reservation.findFirst({ where: { id: { equals: id } } });
    if (!r) return err("预约不存在", 404);
    const svc = db.service.findFirst({ where: { id: { equals: r.service_id } } });
    const updated = db.reservation.update({
      where: { id: { equals: id } },
      data: { status: "confirmed", updated_at: new Date().toISOString() },
    });
    const notifMax = db.notification.count();
    db.notification.create({
      id: notifMax + 1,
      user_id: r.user_id,
      title: "预约已确认",
      content: `您预约的「${svc?.title ?? "服务"}」已由商家确认，请按时到达。`,
      type: "reservation_confirmed",
      is_read: false,
      created_at: new Date().toISOString(),
    });
    return json(updated);
  }),

  http.put(`${API}/provider/reservations/:id/reject`, ({ params }) => {
    const id = Number(params.id);
    const r = db.reservation.findFirst({ where: { id: { equals: id } } });
    if (!r) return err("预约不存在", 404);
    const svc = db.service.findFirst({ where: { id: { equals: r.service_id } } });
    const updated = db.reservation.update({
      where: { id: { equals: id } },
      data: { status: "rejected", updated_at: new Date().toISOString() },
    });
    const notifMax = db.notification.count();
    db.notification.create({
      id: notifMax + 1,
      user_id: r.user_id,
      title: "预约已拒绝",
      content: `您预约的「${svc?.title ?? "服务"}」已被商家拒绝，请查看其他服务。`,
      type: "system",
      is_read: false,
      created_at: new Date().toISOString(),
    });
    return json(updated);
  }),

  /* ── Reservations (customer) ────────────────────────────── */

  http.post(`${API}/reservations`, async ({ request }) => {
    const userId = getUserId(request);
    if (!userId) return err("缺少认证信息", 401);
    const body: any = await request.json();
    const max = db.reservation.count();
    const now = new Date().toISOString();
    const svc = db.service.findFirst({ where: { id: { equals: body.service_id } } });
    if (!svc) return err("服务不存在", 404);
    const endTime = new Date(new Date(body.start_time).getTime() + (svc?.duration_minutes ?? 60) * 60000).toISOString();
    const r = db.reservation.create({
      id: max + 1,
      user_id: userId,
      service_id: body.service_id,
      start_time: body.start_time,
      end_time: endTime,
      status: "pending",
      note: body.note ?? "",
      created_at: now,
      updated_at: now,
    });
    const p = db.provider.findFirst({ where: { id: { equals: svc.provider_id } } });
    const notifMax = db.notification.count();
    db.notification.create({
      id: notifMax + 1,
      user_id: userId,
      title: "预约已创建",
      content: `您已成功预约「${svc.title}」，请耐心等待商家确认。`,
      type: "system",
      is_read: false,
      created_at: now,
    });
    return json({
      id: r.id,
      service: {
        id: svc.id,
        title: svc.title,
        provider: { id: p?.id ?? 0, business_name: p?.business_name ?? "" },
      },
      start_time: r.start_time,
      end_time: r.end_time,
      status: r.status,
      note: r.note,
      created_at: r.created_at,
    }, 201, "created");
  }),

  http.get(`${API}/reservations`, ({ request }) => {
    const userId = getUserId(request);
    if (!userId) return err("缺少认证信息", 401);
    const url = new URL(request.url);
    const { page, pageSize, offset } = pageParams(url);
    const status = url.searchParams.get("status");
    let list = db.reservation.findMany({ where: { user_id: { equals: userId } } });
    if (status) list = list.filter((r) => r.status === status);
    list.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
    const items = list.slice(offset, offset + pageSize);
    return json({ items, page, page_size: pageSize });
  }),

  http.get(`${API}/reservations/:id`, ({ params, request }) => {
    const userId = getUserId(request);
    if (!userId) return err("缺少认证信息", 401);
    const id = Number(params.id);
    const r = db.reservation.findFirst({ where: { id: { equals: id }, user_id: { equals: userId } } });
    if (!r) return err("预约不存在", 404);
    const svc = db.service.findFirst({ where: { id: { equals: r.service_id } } });
    const p = svc ? db.provider.findFirst({ where: { id: { equals: svc.provider_id } } }) : null;
    return json({
      ...r,
      service: svc
        ? { id: svc.id, title: svc.title, provider: { id: p?.id ?? 0, business_name: p?.business_name ?? "" } }
        : null,
    });
  }),

  http.put(`${API}/reservations/:id/cancel`, ({ params, request }) => {
    const userId = getUserId(request);
    if (!userId) return err("缺少认证信息", 401);
    const id = Number(params.id);
    const r = db.reservation.findFirst({ where: { id: { equals: id }, user_id: { equals: userId } } });
    if (!r) return err("预约不存在", 404);
    if (r.status !== "pending" && r.status !== "confirmed") return err("当前预约状态不可取消", 400);
    const updated = db.reservation.update({
      where: { id: { equals: id } },
      data: { status: "cancelled", updated_at: new Date().toISOString() },
    });
    return json(updated);
  }),

  /* ── User interests ────────────────────────────────────── */

  http.get(`${API}/users/me/interests`, ({ request }) => {
    const userId = getUserId(request);
    if (!userId) return err("缺少认证信息", 401);
    const rels = db.user_interest.findMany({ where: { user_id: { equals: userId } } });
    const tags = rels.map((r) => {
      const t = db.tag.findFirst({ where: { id: { equals: r.tag_id } } });
      return t ? { id: t.id, name: t.name, created_at: t.created_at } : null;
    }).filter(Boolean);
    return json(tags);
  }),

  http.put(`${API}/users/me/interests`, async ({ request }) => {
    const userId = getUserId(request);
    if (!userId) return err("缺少认证信息", 401);
    const body: any = await request.json();
    db.user_interest.deleteMany({ where: { user_id: { equals: userId } } });
    let max = db.user_interest.count();
    for (const tagId of body.tag_ids) {
      max++;
      db.user_interest.create({ id: max, user_id: userId, tag_id: tagId });
    }
    const rels = db.user_interest.findMany({ where: { user_id: { equals: userId } } });
    const tags = rels.map((r) => {
      const t = db.tag.findFirst({ where: { id: { equals: r.tag_id } } });
      return t ? { id: t.id, name: t.name, created_at: t.created_at } : null;
    }).filter(Boolean);
    return json(tags);
  }),

  /* ── Notifications ───────────────────────────────────── */

  http.get(`${API}/notifications`, ({ request }) => {
    const userId = getUserId(request);
    if (!userId) return err("缺少认证信息", 401);
    const url = new URL(request.url);
    const isRead = url.searchParams.get("is_read");
    let list = db.notification.findMany({ where: { user_id: { equals: userId } } });
    if (isRead === "true") list = list.filter((n) => n.is_read);
    else if (isRead === "false") list = list.filter((n) => !n.is_read);
    list.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
    return json(list);
  }),

  http.get(`${API}/notifications/unread-count`, ({ request }) => {
    const userId = getUserId(request);
    if (!userId) return err("缺少认证信息", 401);
    const list = db.notification.findMany({ where: { user_id: { equals: userId } } });
    const count = list.filter((n) => !n.is_read).length;
    return json({ count });
  }),

  http.put(`${API}/notifications/:id/read`, ({ params, request }) => {
    const userId = getUserId(request);
    if (!userId) return err("缺少认证信息", 401);
    const id = Number(params.id);
    const n = db.notification.findFirst({ where: { id: { equals: id }, user_id: { equals: userId } } });
    if (!n) return err("通知不存在", 404);
    const updated = db.notification.update({
      where: { id: { equals: id } },
      data: { is_read: true },
    });
    return json(updated);
  }),

  http.put(`${API}/notifications/read-all`, ({ request }) => {
    const userId = getUserId(request);
    if (!userId) return err("缺少认证信息", 401);
    const list = db.notification.findMany({ where: { user_id: { equals: userId } } });
    let affected = 0;
    for (const n of list) {
      if (!n.is_read) {
        db.notification.update({ where: { id: { equals: n.id } }, data: { is_read: true } });
        affected++;
      }
    }
    return json({ affected });
  }),
];
