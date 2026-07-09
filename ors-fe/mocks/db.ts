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

  /* ── 4 大类 + 8 小类 ─────────────────────────────────── */

  db.category.create({
    id: 1,
    name: "美容美体",
    description: "美容美体服务",
    parent_id: 0,
    created_at: now,
  });
  db.category.create({
    id: 101,
    name: "面部护理",
    description: "面部护理",
    parent_id: 1,
    created_at: now,
  });
  db.category.create({
    id: 102,
    name: "身体护理",
    description: "身体护理",
    parent_id: 1,
    created_at: now,
  });
  db.category.create({
    id: 2,
    name: "舒缓理疗",
    description: "舒缓理疗服务",
    parent_id: 0,
    created_at: now,
  });
  db.category.create({
    id: 201,
    name: "推拿按摩",
    description: "推拿按摩",
    parent_id: 2,
    created_at: now,
  });
  db.category.create({
    id: 202,
    name: "足疗",
    description: "足疗",
    parent_id: 2,
    created_at: now,
  });
  db.category.create({
    id: 3,
    name: "深层面护",
    description: "深层面护服务",
    parent_id: 0,
    created_at: now,
  });
  db.category.create({
    id: 301,
    name: "清洁护理",
    description: "清洁护理",
    parent_id: 3,
    created_at: now,
  });
  db.category.create({
    id: 302,
    name: "抗衰紧致",
    description: "抗衰紧致",
    parent_id: 3,
    created_at: now,
  });
  db.category.create({
    id: 4,
    name: "疗愈静心",
    description: "疗愈静心服务",
    parent_id: 0,
    created_at: now,
  });
  db.category.create({
    id: 401,
    name: "瑜伽冥想",
    description: "瑜伽冥想",
    parent_id: 4,
    created_at: now,
  });
  db.category.create({
    id: 402,
    name: "音疗SPA",
    description: "音疗SPA",
    parent_id: 4,
    created_at: now,
  });

  /* ── 19 个服务（低分服务 id 18,19 为新增） ────────────── */

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
    category_id: 201,
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
    category_id: 101,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 3,
    title: "全身精油美体",
    description: "使用天然植物精油进行全身按摩，滋润肌肤，放松身心。",
    price: 398,
    duration_minutes: 90,
    avg_rating: 4.5,
    review_count: 76,
    status: "active",
    image_url: "",
    provider_id: 1,
    category_id: 102,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 4,
    title: "背部经络疏通",
    description: "专业背部经络疏通手法，缓解背部酸痛，促进血液循环。",
    price: 268,
    duration_minutes: 60,
    avg_rating: 4.7,
    review_count: 93,
    status: "active",
    image_url: "",
    provider_id: 1,
    category_id: 102,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 5,
    title: "韩式半永久纹眉",
    description: "自然仿真眉形设计，进口色料，持久不褪色。",
    price: 899,
    duration_minutes: 120,
    avg_rating: 4.6,
    review_count: 45,
    status: "active",
    image_url: "",
    provider_id: 2,
    category_id: 101,
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
    category_id: 401,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 7,
    title: "精油推背",
    description: "使用天然植物精油进行背部推拿，促进血液循环，缓解疲劳。",
    price: 258,
    duration_minutes: 75,
    avg_rating: 4.6,
    review_count: 89,
    status: "active",
    image_url: "",
    provider_id: 1,
    category_id: 201,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 8,
    title: "足底按摩 45 分钟",
    description: "专业足底穴位按摩，缓解脚部疲劳，改善睡眠质量。",
    price: 159,
    duration_minutes: 45,
    avg_rating: 4.4,
    review_count: 75,
    status: "active",
    image_url: "",
    provider_id: 1,
    category_id: 202,
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
    category_id: 201,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 10,
    title: "中药泡脚足疗",
    description: "精选中药材泡脚，配合足底按摩，驱寒祛湿。",
    price: 198,
    duration_minutes: 60,
    avg_rating: 4.5,
    review_count: 56,
    status: "active",
    image_url: "",
    provider_id: 1,
    category_id: 202,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 11,
    title: "小气泡深层清洁",
    description: "韩国小气泡仪器深层清洁毛孔，去除黑头粉刺。",
    price: 198,
    duration_minutes: 60,
    avg_rating: 4.5,
    review_count: 134,
    status: "active",
    image_url: "",
    provider_id: 2,
    category_id: 301,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 12,
    title: "毛孔净化管理",
    description: "深层清洁毛孔，控油消炎，改善痘痘肌。",
    price: 298,
    duration_minutes: 75,
    avg_rating: 4.6,
    review_count: 88,
    status: "active",
    image_url: "",
    provider_id: 2,
    category_id: 301,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 13,
    title: "射频紧肤",
    description: "采用射频技术紧致肌肤，淡化细纹，提升面部轮廓。",
    price: 599,
    duration_minutes: 90,
    avg_rating: 4.7,
    review_count: 192,
    status: "active",
    image_url: "",
    provider_id: 2,
    category_id: 302,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 14,
    title: "胶原蛋白导入",
    description: "进口胶原蛋白精华导入，深层补水，恢复肌肤弹性。",
    price: 399,
    duration_minutes: 60,
    avg_rating: 4.5,
    review_count: 67,
    status: "active",
    image_url: "",
    provider_id: 2,
    category_id: 302,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 15,
    title: "冥想指导课",
    description: "专业冥想导师指导呼吸与放松技巧，缓解压力焦虑。",
    price: 199,
    duration_minutes: 60,
    avg_rating: 4.5,
    review_count: 42,
    status: "active",
    image_url: "",
    provider_id: 1,
    category_id: 401,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 16,
    title: "颂钵音疗",
    description: "喜马拉雅颂钵音频共振，深层放松身心，平衡能量。",
    price: 358,
    duration_minutes: 75,
    avg_rating: 4.8,
    review_count: 103,
    status: "active",
    image_url: "",
    provider_id: 2,
    category_id: 402,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 17,
    title: "芳香SPA",
    description: "精选天然芳香精油，配合专业SPA手法，全身心放松体验。",
    price: 498,
    duration_minutes: 90,
    avg_rating: 4.7,
    review_count: 156,
    status: "active",
    image_url: "",
    provider_id: 2,
    category_id: 402,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 18,
    title: "基础足部按摩 30 分钟",
    description: "初级技师提供的基础足部按摩服务，适合想简单放松的顾客。",
    price: 99,
    duration_minutes: 30,
    avg_rating: 3.2,
    review_count: 28,
    status: "active",
    image_url: "",
    provider_id: 1,
    category_id: 202,
    created_at: now,
    updated_at: now,
  });
  db.service.create({
    id: 19,
    title: "快速面部清洁",
    description: "基础面部清洁服务，适合预算有限的顾客，效果因肤质而异。",
    price: 88,
    duration_minutes: 30,
    avg_rating: 3.5,
    review_count: 42,
    status: "active",
    image_url: "",
    provider_id: 2,
    category_id: 301,
    created_at: now,
    updated_at: now,
  });

  /* ── Tags ─────────────────────────────────────────────── */

  db.tag.create({
    id: 1,
    name: "放松",
    created_at: now,
  });
  db.tag.create({
    id: 2,
    name: "舒缓",
    created_at: now,
  });
  db.tag.create({
    id: 3,
    name: "深层清洁",
    created_at: now,
  });
  db.tag.create({
    id: 4,
    name: "减脂",
    created_at: now,
  });
  db.tag.create({
    id: 5,
    name: "增肌",
    created_at: now,
  });
  db.tag.create({
    id: 6,
    name: "美白",
    created_at: now,
  });
  db.tag.create({
    id: 7,
    name: "抗衰老",
    created_at: now,
  });
  db.tag.create({
    id: 8,
    name: "口腔健康",
    created_at: now,
  });

  db.service_tag.create({
    id: 1,
    service_id: 1,
    tag_id: 1,
  });
  db.service_tag.create({
    id: 2,
    service_id: 1,
    tag_id: 2,
  });
  db.service_tag.create({
    id: 3,
    service_id: 2,
    tag_id: 3,
  });
  db.service_tag.create({
    id: 4,
    service_id: 2,
    tag_id: 6,
  });
  db.service_tag.create({
    id: 5,
    service_id: 6,
    tag_id: 4,
  });
  db.service_tag.create({
    id: 6,
    service_id: 9,
    tag_id: 1,
  });
  db.service_tag.create({
    id: 7,
    service_id: 3,
    tag_id: 2,
  });
  db.service_tag.create({
    id: 8,
    service_id: 13,
    tag_id: 7,
  });
  db.service_tag.create({
    id: 9,
    service_id: 16,
    tag_id: 2,
  });
  db.service_tag.create({
    id: 10,
    service_id: 17,
    tag_id: 1,
  });

  /* ── Reservations（保留，引用的 service_id 1,2,6,9 未变） ── */

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

  const reviewComments: Record<number, { user_id: number; rating: number; comment: string }[]> = {
    1: [
      { user_id: 3, rating: 5, comment: "非常专业，按完之后舒服多了。" },
      { user_id: 1, rating: 4, comment: "环境很好，手法专业，下次还来。" },
      { user_id: 3, rating: 5, comment: "肩颈舒服多了，老师傅手艺好！" },
      { user_id: 1, rating: 4, comment: "价格合理，服务态度好。" },
      { user_id: 3, rating: 5, comment: "强烈推荐，长期伏案工作者的福音" },
      { user_id: 1, rating: 3, comment: "一般般，没有想象中好" },
      { user_id: 3, rating: 5, comment: "按完轻松很多，会再来的" },
    ],
    2: [
      { user_id: 1, rating: 5, comment: "做完皮肤明显变亮，效果很好" },
      { user_id: 3, rating: 4, comment: "姐姐很细心，清洁很到位" },
      { user_id: 1, rating: 5, comment: "黑头少了很多，好评" },
      { user_id: 3, rating: 4, comment: "仪器很专业，值得尝试" },
    ],
    3: [
      { user_id: 1, rating: 5, comment: "精油很好闻，全身都放松了" },
      { user_id: 3, rating: 4, comment: "手法不错，环境舒适" },
      { user_id: 1, rating: 3, comment: "价格有点贵，性价比一般" },
      { user_id: 3, rating: 5, comment: "做完皮肤滑滑的，推荐！" },
    ],
    4: [
      { user_id: 3, rating: 5, comment: "背部经络疏通效果明显" },
      { user_id: 1, rating: 4, comment: "按完很舒服，下次还来" },
      { user_id: 3, rating: 4, comment: "技师很专业" },
    ],
    5: [
      { user_id: 1, rating: 4, comment: "效果很好，颜色自然" },
      { user_id: 3, rating: 5, comment: "纹完很满意，朋友都说好看" },
      { user_id: 1, rating: 3, comment: "恢复期有点长" },
      { user_id: 3, rating: 5, comment: "很自然，不疼" },
    ],
    6: [
      { user_id: 3, rating: 5, comment: "老师非常专业，一节课下来收获很多" },
      { user_id: 1, rating: 5, comment: "环境安静，适合冥想放松" },
      { user_id: 3, rating: 4, comment: "瑜伽课程安排合理" },
    ],
    7: [
      { user_id: 1, rating: 4, comment: "精油推背很舒服" },
      { user_id: 3, rating: 5, comment: "背部轻松了很多" },
      { user_id: 1, rating: 4, comment: "技师手法到位" },
      { user_id: 3, rating: 3, comment: "价格有点小贵" },
    ],
    8: [
      { user_id: 3, rating: 4, comment: "按完脚很舒服" },
      { user_id: 1, rating: 5, comment: "穴位找得很准" },
      { user_id: 3, rating: 4, comment: "环境干净整洁" },
    ],
    9: [
      { user_id: 1, rating: 5, comment: "全身推拿非常舒服，强烈推荐" },
      { user_id: 3, rating: 5, comment: "90分钟彻底放松，物超所值" },
      { user_id: 1, rating: 4, comment: "手法好，服务周到" },
    ],
    10: [
      { user_id: 3, rating: 4, comment: "中药泡脚很舒服" },
      { user_id: 1, rating: 5, comment: "泡完全身暖和" },
      { user_id: 3, rating: 4, comment: "药材味道很好闻" },
    ],
    11: [
      { user_id: 1, rating: 5, comment: "清洁效果很明显" },
      { user_id: 3, rating: 4, comment: "做完脸很干净" },
      { user_id: 1, rating: 4, comment: "仪器操作很专业" },
      { user_id: 3, rating: 2, comment: "效果一般，没有想象中好" },
    ],
    12: [
      { user_id: 3, rating: 5, comment: "毛孔真的变小了" },
      { user_id: 1, rating: 4, comment: "痘痘改善了很多" },
      { user_id: 3, rating: 4, comment: "控油效果不错" },
    ],
    13: [
      { user_id: 1, rating: 5, comment: "做完皮肤紧致了很多" },
      { user_id: 3, rating: 4, comment: "效果很明显，会继续做" },
      { user_id: 1, rating: 5, comment: "美容师很细心" },
      { user_id: 3, rating: 3, comment: "价格偏高" },
    ],
    14: [
      { user_id: 3, rating: 5, comment: "导入完皮肤水嫩嫩的" },
      { user_id: 1, rating: 4, comment: "补水效果很好" },
      { user_id: 3, rating: 4, comment: "服务态度很好" },
    ],
    15: [
      { user_id: 1, rating: 5, comment: "冥想完心情平静了很多" },
      { user_id: 3, rating: 4, comment: "老师引导得很好" },
      { user_id: 1, rating: 4, comment: "环境很安静" },
    ],
    16: [
      { user_id: 3, rating: 5, comment: "颂钵的声音太治愈了" },
      { user_id: 1, rating: 5, comment: "做完身心都放松了" },
      { user_id: 3, rating: 4, comment: "很特别的体验" },
      { user_id: 1, rating: 5, comment: "强烈推荐压力大的人试试" },
    ],
    17: [
      { user_id: 1, rating: 5, comment: "全身SPA太享受了" },
      { user_id: 3, rating: 4, comment: "精油味道很好闻" },
      { user_id: 1, rating: 5, comment: "环境和服务都是一流的" },
    ],
    18: [
      { user_id: 3, rating: 3, comment: "时间太短了，刚有点感觉就结束了" },
      { user_id: 1, rating: 2, comment: "技师手法一般，不够专业" },
      { user_id: 3, rating: 4, comment: "价格便宜，随便按按还行" },
      { user_id: 1, rating: 3, comment: "性价比一般，一分钱一分货" },
      { user_id: 3, rating: 2, comment: "不太满意，按完没什么感觉" },
    ],
    19: [
      { user_id: 1, rating: 3, comment: "清洁不太彻底" },
      { user_id: 3, rating: 4, comment: "基础清洁够用了，价格实惠" },
      { user_id: 1, rating: 2, comment: "做完有点过敏，可能肤质不适合" },
      { user_id: 3, rating: 4, comment: "对于这个价格来说还不错" },
      { user_id: 1, rating: 3, comment: "效果一般，没有广告说的那么好" },
    ],
  };

  let reviewId = 500;
  for (const [serviceId, comments] of Object.entries(reviewComments)) {
    for (const c of comments) {
      reviewId++;
      db.review.create({
        id: reviewId,
        reservation_id: reviewId,
        user_id: c.user_id,
        service_id: Number(serviceId),
        rating: c.rating,
        comment: c.comment,
        created_at: now,
      });
    }
  }

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
