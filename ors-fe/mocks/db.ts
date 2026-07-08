import { factory, primaryKey } from "@mswjs/data";

export const db = factory({
  user: {
    id: primaryKey(Number),
    name: String,
    email: String,
    password: String,
    role: String,
    phone: String,
    avatar_url: String,
    created_at: String,
    updated_at: String,
  },
  provider: {
    id: primaryKey(Number),
    user_id: Number,
    business_name: String,
    description: String,
    address: String,
    phone: String,
    email: String,
    logo_url: String,
    created_at: String,
    updated_at: String,
  },
  category: {
    id: primaryKey(Number),
    name: String,
    description: String,
    parent_id: Number,
    created_at: String,
  },
  service: {
    id: primaryKey(Number),
    title: String,
    description: String,
    price: Number,
    duration_minutes: Number,
    avg_rating: Number,
    review_count: Number,
    status: String,
    image_url: String,
    provider_id: Number,
    category_id: Number,
    created_at: String,
    updated_at: String,
  },
  tag: {
    id: primaryKey(Number),
    name: String,
    created_at: String,
  },
  service_tag: {
    id: primaryKey(Number),
    service_id: Number,
    tag_id: Number,
  },
  reservation: {
    id: primaryKey(Number),
    user_id: Number,
    service_id: Number,
    start_time: String,
    end_time: String,
    status: String,
    note: String,
    created_at: String,
    updated_at: String,
  },
  review: {
    id: primaryKey(Number),
    reservation_id: Number,
    user_id: Number,
    service_id: Number,
    rating: Number,
    comment: String,
    created_at: String,
  },
  user_interest: {
    id: primaryKey(Number),
    user_id: Number,
    tag_id: Number,
  },
  notification: {
    id: primaryKey(Number),
    user_id: Number,
    title: String,
    content: String,
    type: String,
    is_read: Boolean,
    created_at: String,
  },
});

const now = "2026-07-08T00:00:00Z";

export function seed() {
  if (db.user.count() > 0) return;

  db.user.create({
    id: 1,
    name: "张三",
    email: "zhangsan@example.com",
    password: "Abcd1234",
    role: "customer",
    phone: "13800000001",
    avatar_url: "",
    created_at: now,
    updated_at: now,
  });
  db.user.create({
    id: 2,
    name: "李四",
    email: "lisi@example.com",
    password: "Abcd1234",
    role: "provider",
    phone: "13800000002",
    avatar_url: "",
    created_at: now,
    updated_at: now,
  });
  db.user.create({
    id: 3,
    name: "王五",
    email: "wangwu@example.com",
    password: "Abcd1234",
    role: "customer",
    phone: "13800000003",
    avatar_url: "",
    created_at: now,
    updated_at: now,
  });
  db.user.create({
    id: 4,
    name: "赵六",
    email: "zhaoliu@example.com",
    password: "Abcd1234",
    role: "provider",
    phone: "13800000004",
    avatar_url: "",
    created_at: now,
    updated_at: now,
  });

  db.provider.create({
    id: 1,
    user_id: 2,
    business_name: "舒心养生馆",
    description: "专业养生按摩服务，十年老店",
    address: "上海市徐汇区",
    phone: "13800000002",
    email: "shuxin@example.com",
    logo_url: "",
    created_at: now,
    updated_at: now,
  });
  db.provider.create({
    id: 2,
    user_id: 4,
    business_name: "美颜坊",
    description: "专业美容护肤服务",
    address: "上海市浦东新区",
    phone: "13800000004",
    email: "meiyan@example.com",
    logo_url: "",
    created_at: now,
    updated_at: now,
  });

  db.category.create({ id: 1,  name: "医疗健康", description: "体检、中医理疗、康复等服务",           parent_id: 0, created_at: now });
  db.category.create({ id: 2,  name: "口腔护理", description: "洁牙、矫正、种植等口腔服务",           parent_id: 0, created_at: now });
  db.category.create({ id: 3,  name: "美容护肤", description: "皮肤管理、轻医美、纹绣等服务",          parent_id: 0, created_at: now });
  db.category.create({ id: 4,  name: "美发造型", description: "剪发、染烫、头皮护理等",               parent_id: 0, created_at: now });
  db.category.create({ id: 5,  name: "健身运动", description: "私教、瑜伽、团课等健身服务",            parent_id: 0, created_at: now });
  db.category.create({ id: 6,  name: "按摩推拿", description: "足疗、全身按摩、SPA 水疗等",            parent_id: 0, created_at: now });
  db.category.create({ id: 7,  name: "家政服务", description: "保洁、家电清洗、搬家等",               parent_id: 0, created_at: now });
  db.category.create({ id: 8,  name: "宠物服务", description: "洗护、美容、寄养等服务",               parent_id: 0, created_at: now });
  db.category.create({ id: 9,  name: "体检中心", description: "常规体检、入职体检、高端体检套餐",      parent_id: 1, created_at: now });
  db.category.create({ id: 10, name: "中医理疗", description: "针灸、拔罐、艾灸、中药调理",           parent_id: 1, created_at: now });

  db.service.create({
    id: 1,
    title: "肩颈按摩 60 分钟",
    description: "通过揉捏按压缓解肩颈肌肉酸痛，适合长期伏案工作者。",
    price: 199,
    duration_minutes: 60,
    avg_rating: 4.5,
    review_count: 128,
    status: "active",
    image_url: "",
    provider_id: 1,
    category_id: 6,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 2,
    title: "深层清洁护肤",
    description: "使用专业仪器深层清洁毛孔，去除黑头粉刺，提亮肤色。",
    price: 298,
    duration_minutes: 90,
    avg_rating: 4.8,
    review_count: 256,
    status: "active",
    image_url: "",
    provider_id: 2,
    category_id: 3,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 3,
    title: "私人健身指导",
    description: "一对一专业健身教练指导，量身定制训练计划。",
    price: 399,
    duration_minutes: 60,
    avg_rating: 4.3,
    review_count: 64,
    status: "pending",
    image_url: "",
    provider_id: 1,
    category_id: 5,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 4,
    title: "精油推背",
    description: "使用天然植物精油进行背部推拿，促进血液循环，缓解疲劳。",
    price: 258,
    duration_minutes: 75,
    avg_rating: 4.6,
    review_count: 89,
    status: "rejected",
    image_url: "",
    provider_id: 1,
    category_id: 6,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 5,
    title: "水光针护理",
    description: "进口水光针仪器导入玻尿酸精华，深层补水保湿。",
    price: 599,
    duration_minutes: 120,
    avg_rating: 4.7,
    review_count: 192,
    status: "inactive",
    image_url: "",
    provider_id: 2,
    category_id: 3,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 6,
    title: "瑜伽私教课",
    description: "专业瑜伽导师一对一教学，涵盖哈他瑜伽、流瑜伽等多种流派。",
    price: 329,
    duration_minutes: 90,
    avg_rating: 4.9,
    review_count: 312,
    status: "active",
    image_url: "",
    provider_id: 1,
    category_id: 5,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 7,
    title: "足底按摩 45 分钟",
    description: "专业足底穴位按摩，缓解脚部疲劳，改善睡眠质量。",
    price: 159,
    duration_minutes: 45,
    avg_rating: 4.4,
    review_count: 75,
    status: "active",
    image_url: "",
    provider_id: 1,
    category_id: 6,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 8,
    title: "韩式半永久纹眉",
    description: "自然仿真眉形设计，进口色料，持久不褪色。",
    price: 899,
    duration_minutes: 120,
    avg_rating: 4.6,
    review_count: 45,
    status: "active",
    image_url: "",
    provider_id: 2,
    category_id: 3,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 9,
    title: "全身推拿 90 分钟",
    description: "从头到脚全身推拿，疏通经络，缓解全身疲劳。",
    price: 358,
    duration_minutes: 90,
    avg_rating: 4.7,
    review_count: 210,
    status: "active",
    image_url: "",
    provider_id: 1,
    category_id: 6,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 10,
    title: "洁牙护理",
    description: "专业超声波洁牙，去除牙结石和牙菌斑。",
    price: 199,
    duration_minutes: 45,
    avg_rating: 4.5,
    review_count: 168,
    status: "active",
    image_url: "",
    provider_id: 2,
    category_id: 2,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 11,
    title: "入职体检套餐",
    description: "包含常规血检、尿检、胸透、心电图等基础项目。",
    price: 299,
    duration_minutes: 120,
    avg_rating: 4.2,
    review_count: 56,
    status: "active",
    image_url: "",
    provider_id: 1,
    category_id: 9,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 12,
    title: "艾灸调理",
    description: "传统艾灸疗法，温经散寒，调理气血，适合亚健康人群。",
    price: 168,
    duration_minutes: 60,
    avg_rating: 4.6,
    review_count: 34,
    status: "active",
    image_url: "",
    provider_id: 1,
    category_id: 10,
    created_at: now,
    updated_at: now,
  });

  db.tag.create({ id: 1, name: "放松", created_at: now });
  db.tag.create({ id: 2, name: "舒缓", created_at: now });
  db.tag.create({ id: 3, name: "深层清洁", created_at: now });
  db.tag.create({ id: 4, name: "减脂", created_at: now });
  db.tag.create({ id: 5, name: "增肌", created_at: now });
  db.tag.create({ id: 6, name: "美白", created_at: now });
  db.tag.create({ id: 7, name: "抗衰老", created_at: now });
  db.tag.create({ id: 8, name: "口腔健康", created_at: now });

  db.service_tag.create({ id: 1, service_id: 1, tag_id: 1 });
  db.service_tag.create({ id: 2, service_id: 1, tag_id: 2 });
  db.service_tag.create({ id: 3, service_id: 2, tag_id: 3 });
  db.service_tag.create({ id: 4, service_id: 2, tag_id: 6 });
  db.service_tag.create({ id: 5, service_id: 3, tag_id: 4 });
  db.service_tag.create({ id: 6, service_id: 3, tag_id: 5 });
  db.service_tag.create({ id: 7, service_id: 6, tag_id: 4 });
  db.service_tag.create({ id: 8, service_id: 8, tag_id: 6 });
  db.service_tag.create({ id: 9, service_id: 8, tag_id: 7 });
  db.service_tag.create({ id: 10, service_id: 10, tag_id: 8 });

  db.reservation.create({
    id: 1001,
    user_id: 1,
    service_id: 1,
    start_time: "2026-07-10T14:00:00Z",
    end_time: "2026-07-10T15:00:00Z",
    status: "pending",
    note: "请准备热水",
    created_at: now,
    updated_at: now,
  });
  db.reservation.create({
    id: 1002,
    user_id: 1,
    service_id: 2,
    start_time: "2026-07-11T10:00:00Z",
    end_time: "2026-07-11T11:30:00Z",
    status: "confirmed",
    note: "",
    created_at: now,
    updated_at: now,
  });
  db.reservation.create({
    id: 1003,
    user_id: 3,
    service_id: 6,
    start_time: "2026-07-12T09:00:00Z",
    end_time: "2026-07-12T10:30:00Z",
    status: "pending",
    note: "请准备瑜伽垫",
    created_at: now,
    updated_at: now,
  });
  db.reservation.create({
    id: 1004,
    user_id: 3,
    service_id: 1,
    start_time: "2026-07-09T16:00:00Z",
    end_time: "2026-07-09T17:00:00Z",
    status: "completed",
    note: "",
    created_at: now,
    updated_at: now,
  });
  db.reservation.create({
    id: 1005,
    user_id: 1,
    service_id: 9,
    start_time: "2026-07-13T14:00:00Z",
    end_time: "2026-07-13T15:30:00Z",
    status: "pending",
    note: "最近腰疼",
    created_at: now,
    updated_at: now,
  });

  db.review.create({
    id: 501,
    reservation_id: 1004,
    user_id: 3,
    service_id: 1,
    rating: 5,
    comment: "非常专业，按完之后舒服多了。",
    created_at: now,
  });

  db.notification.create({
    id: 1,
    user_id: 1,
    title: "预约已确认",
    content: "您预约的「肩颈按摩 60 分钟」已由舒心养生馆确认，请按时到达。",
    type: "reservation_confirmed",
    is_read: false,
    created_at: now,
  });
  db.notification.create({
    id: 2,
    user_id: 1,
    title: "预约即将开始",
    content: "您预约的「深层清洁护肤」将于 2026-07-11 10:00 开始，请提前准备。",
    type: "reservation_reminder",
    is_read: false,
    created_at: now,
  });
  db.notification.create({
    id: 3,
    user_id: 1,
    title: "预约已创建",
    content: "您已成功预约「全身推拿 90 分钟」，请耐心等待商家确认。",
    type: "system",
    is_read: true,
    created_at: now,
  });
}
