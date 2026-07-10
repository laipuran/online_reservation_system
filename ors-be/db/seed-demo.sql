-- 演示数据种子
-- 密码统一为 Abcd1234 (bcrypt hash)
-- 使用 ON CONFLICT 保证幂等，重复执行安全

-- ============================================================
-- 1. 用户账号
-- ============================================================
INSERT INTO users (id, name, email, password_hash, role, phone)
VALUES
  (1, '阳光健康',   'yangguang@test.com',  '$2a$10$HWaXgbWWdN/as92PS78ESu5bQ6.aq6Ta.7cOfzFbXoiIBRenDn96i', 'provider', '13800001001'),
  (2, '美丽人生',   'meili@test.com',      '$2a$10$HWaXgbWWdN/as92PS78ESu5bQ6.aq6Ta.7cOfzFbXoiIBRenDn96i', 'provider', '13800001002'),
  (3, '舒心家政',   'shuxin@test.com',     '$2a$10$HWaXgbWWdN/as92PS78ESu5bQ6.aq6Ta.7cOfzFbXoiIBRenDn96i', 'provider', '13800001003'),
  (4, '体验用户',   'customer@test.com',   '$2a$10$HWaXgbWWdN/as92PS78ESu5bQ6.aq6Ta.7cOfzFbXoiIBRenDn96i', 'customer', '13900001004')
ON CONFLICT (email) DO NOTHING;

-- ============================================================
-- 2. 服务商资料
-- ============================================================
INSERT INTO service_providers (id, user_id, business_name, description, address, phone, email)
VALUES
  (1, 1, '阳光健康管理中心',
   '专注经络调理与推拿按摩，十年老店，持证中医师团队，为您提供专业舒适的养生体验。',
   '北京市朝阳区建国路88号阳光大厦3层', '13800001001', 'yangguang@test.com'),
  (2, 2, '美丽人生美容会所',
   '高端皮肤管理与轻医美服务，引进国内外先进设备，资深美容师一对一服务。',
   '上海市静安区南京西路1266号恒隆广场5层', '13800001002', 'meili@test.com'),
  (3, 3, '舒心家政服务有限公司',
   '专业家政服务团队，员工均经过严格培训考核，服务贴心、价格透明。',
   '广州市天河区天河路385号太古汇2层', '13800001003', 'shuxin@test.com')
ON CONFLICT (user_id) DO NOTHING;

-- ============================================================
-- 3. 标签
-- ============================================================
INSERT INTO tags (id, name)
VALUES
  (1, '专业'),
  (2, '舒适'),
  (3, '快速'),
  (4, '实惠'),
  (5, '高品质'),
  (6, '放松'),
  (7, '舒缓'),
  (8, '深层清洁'),
  (9, '养生保健'),
  (10, '一站式')
ON CONFLICT (name) DO NOTHING;

-- ============================================================
-- 4. 服务项目 (15个，每家商户5个)
-- ============================================================

-- 商户1: 阳光健康管理中心
INSERT INTO services (id, provider_id, category_id, title, description, price, duration_minutes, image_url)
VALUES
  (1, 1, 30, '全身经络推拿',
   '专业中医师循经推拿，疏通全身经络，缓解肌肉疲劳，改善亚健康状态。适合久坐办公、经常熬夜人群。',
   198.00, 60, 'https://images.unsplash.com/photo-1544161515-4ab6ce6db874?w=800'),
  (2, 1, 32, '肩颈舒缓理疗',
   '针对肩颈僵硬的专项理疗，结合中医推拿与拉伸手法，30分钟明显缓解酸痛。赠送艾草热敷。',
   128.00, 45, 'https://images.unsplash.com/photo-1519823551278-64ac92734fb1?w=800'),
  (3, 1, 29, '足疗养生套餐',
   '中药足浴 + 足底穴位按摩 + 小腿放松，使用天然草本足浴包，促进血液循环，改善睡眠质量。',
   168.00, 60, 'https://images.unsplash.com/photo-1512295767273-ac109ac3acfa?w=800'),
  (4, 1, 10, '中医艾灸调理',
   '精选三年陈艾，针对体质定制艾灸方案，温经散寒、补气养血。含中医体质辨识服务。',
   158.00, 45, ''),
  (5, 1, 31, '精油开背护理',
   '使用进口植物精油，配合专业开背手法，深层放松背部肌肉，疏通膀胱经。赠送花茶一杯。',
   188.00, 50, 'https://images.unsplash.com/photo-1600334089648-b0d9d3028eb2?w=800')
ON CONFLICT (id) DO NOTHING;

-- 商户2: 美丽人生美容会所
INSERT INTO services (id, provider_id, category_id, title, description, price, duration_minutes, image_url)
VALUES
  (6, 2, 17, '深层清洁护理',
   '韩国小气泡 + 黑头导出 + 毛孔收缩三步曲，温和清除毛孔污垢，还原肌肤清爽透亮。',
   258.00, 60, 'https://images.unsplash.com/photo-1570172619644-dfd03ed5d881?w=800'),
  (7, 2, 18, '水光导入护理',
   '采用进口水光仪器，将玻尿酸等营养成分直达肌底，深层补水锁水，改善干燥暗沉。',
   328.00, 45, 'https://images.unsplash.com/photo-1560750588-73207b1ef5b8?w=800'),
  (8, 2, 20, '韩式半永久纹眉',
   '资深纹绣师一对一设计眉形，纯植物色乳，自然雾感效果，赠送补色一次。',
   888.00, 90, 'https://images.unsplash.com/photo-1522337360788-8b13dee7a37e?w=800'),
  (9, 2, 19, '玻尿酸导入',
   '医用级玻尿酸导入，配合专业面部按摩手法，提亮肤色、淡化细纹，效果立竿见影。',
   398.00, 60, 'https://images.unsplash.com/photo-1596755389378-c31d21fd1273?w=800'),
  (10, 2, 30, '肩颈精油按摩',
   '特色肩颈精油推拿，使用品牌精油，配合点按穴位手法，缓解肩颈僵硬与头疼。',
   198.00, 60, 'https://images.unsplash.com/photo-1552693673-1bf958298935?w=800')
ON CONFLICT (id) DO NOTHING;

-- 商户3: 舒心家政服务有限公司
INSERT INTO services (id, provider_id, category_id, title, description, price, duration_minutes, image_url)
VALUES
  (11, 3, 34, '全屋深度保洁',
   '360度无死角深度清洁，含厨房除油、卫生间除垢、窗户玻璃清洁、地板打蜡。约80-120平米。',
   299.00, 180, 'https://images.unsplash.com/photo-1581578731548-c64695cc6952?w=800'),
  (12, 3, 33, '日常保洁套餐',
   '日常家庭清洁维护，含地面清扫拖洗、家具擦拭、卫生间清洁、垃圾清理。约60-80平米。',
   159.00, 120, 'https://images.unsplash.com/photo-1527515637462-cff94eecc1ac?w=800'),
  (13, 3, 35, '空调深度清洗',
   '免拆式高温蒸汽清洗，杀菌除螨除异味，含滤网清洗。挂机/柜机通用。',
   199.00, 90, ''),
  (14, 3, 36, '收纳整理服务',
   '日本收纳师上门服务，全屋空间规划 + 物品分类整理 + 收纳方案设计。赠送收纳盒一套。',
   399.00, 180, 'https://images.unsplash.com/photo-1555041469-a586c61ea9bc?w=800'),
  (15, 3, 35, '油烟机清洗',
   '专业拆卸式油烟机深度清洗，含油网、涡轮、外壳全面清洁，消除火灾隐患。',
   149.00, 90, '')
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- 5. 服务-标签关联
-- ============================================================
INSERT INTO service_tags (service_id, tag_id)
VALUES
  -- 商户1 服务标签
  (1, 1), (1, 2), (1, 9),
  (2, 1), (2, 6), (2, 7),
  (3, 2), (3, 6), (3, 9),
  (4, 1), (4, 9),
  (5, 2), (5, 5), (5, 6),
  -- 商户2 服务标签
  (6, 1), (6, 5), (6, 8),
  (7, 1), (7, 5),
  (8, 1), (8, 5),
  (9, 1), (9, 5),
  (10, 2), (10, 6), (10, 7),
  -- 商户3 服务标签
  (11, 3), (11, 4), (11, 10),
  (12, 3), (12, 4),
  (13, 1), (13, 3),
  (14, 4), (14, 10),
  (15, 3), (15, 4)
ON CONFLICT DO NOTHING;

-- ============================================================
-- 6. 用户兴趣标签 (体验用户)
-- ============================================================
INSERT INTO user_interests (user_id, tag_id)
VALUES
  (4, 1), (4, 2), (4, 6), (4, 7), (4, 9)
ON CONFLICT DO NOTHING;

-- ============================================================
-- 7. 预约记录 (体验用户 -> 全身经络推拿, pending)
-- ============================================================
INSERT INTO reservations (id, user_id, service_id, start_time, end_time, status, note)
VALUES (
  1,
  4,
  1,
  (DATE_TRUNC('day', NOW()) + INTERVAL '1 day' + INTERVAL '10:00:00'),
  (DATE_TRUNC('day', NOW()) + INTERVAL '1 day' + INTERVAL '11:00:00'),
  'pending',
  '第一次体验，希望手法轻柔一些，谢谢。'
)
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- 8. 通知 (给商户: 收到新预约)
-- ============================================================
INSERT INTO notifications (id, user_id, title, content, type)
VALUES (
  1,
  1,
  '收到新预约',
  '用户提交了新的预约，请及时处理。',
  'system'
)
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- 9. 重置序列
-- ============================================================
SELECT setval(pg_get_serial_sequence('users',              'id'), (SELECT MAX(id) FROM users));
SELECT setval(pg_get_serial_sequence('service_providers',  'id'), (SELECT MAX(id) FROM service_providers));
SELECT setval(pg_get_serial_sequence('tags',               'id'), (SELECT MAX(id) FROM tags));
SELECT setval(pg_get_serial_sequence('services',           'id'), (SELECT MAX(id) FROM services));
SELECT setval(pg_get_serial_sequence('reservations',       'id'), (SELECT MAX(id) FROM reservations));
SELECT setval(pg_get_serial_sequence('notifications',      'id'), (SELECT MAX(id) FROM notifications));
