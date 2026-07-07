import { useQuery } from "@tanstack/react-query";
import { fetchServices, type ServiceItem, type ServiceQueryParams } from "../api/service-items";

const MOCK_SERVICES: ServiceItem[] = [
  {
    id: 1,
    title: "肩颈按摩",
    description: "通过揉捏按压缓解肩颈肌肉酸痛，适合长期伏案工作者。",
    price: 199,
    duration_minutes: 60,
    avg_rating: 4.5,
    provider: { id: 1, business_name: "舒心养生馆" },
    category: { id: 1, name: "按摩" },
  },
  {
    id: 2,
    title: "深层清洁护肤",
    description: "使用专业仪器深层清洁毛孔，去除黑头粉刺，提亮肤色。",
    price: 298,
    duration_minutes: 90,
    avg_rating: 4.8,
    provider: { id: 2, business_name: "美颜坊" },
    category: { id: 2, name: "美容" },
  },
  {
    id: 3,
    title: "私人健身指导",
    description: "一对一专业健身教练指导，量身定制训练计划。",
    price: 399,
    duration_minutes: 60,
    avg_rating: 4.3,
    provider: { id: 3, business_name: "FitZone 健身" },
    category: { id: 3, name: "健身" },
  },
  {
    id: 4,
    title: "精油推背",
    description: "使用天然植物精油进行背部推拿，促进血液循环，缓解疲劳。",
    price: 258,
    duration_minutes: 75,
    avg_rating: 4.6,
    provider: { id: 1, business_name: "舒心养生馆" },
    category: { id: 1, name: "按摩" },
  },
  {
    id: 5,
    title: "水光针护理",
    description: "进口水光针仪器导入玻尿酸精华，深层补水保湿。",
    price: 599,
    duration_minutes: 120,
    avg_rating: 4.7,
    provider: { id: 2, business_name: "美颜坊" },
    category: { id: 2, name: "美容" },
  },
  {
    id: 6,
    title: "瑜伽私教课",
    description: "专业瑜伽导师一对一教学，涵盖哈他瑜伽、流瑜伽等多种流派。",
    price: 329,
    duration_minutes: 90,
    avg_rating: 4.9,
    provider: { id: 3, business_name: "FitZone 健身" },
    category: { id: 3, name: "健身" },
  },
];

function getMockServices(params: ServiceQueryParams) {
  const filtered = params.keyword
    ? MOCK_SERVICES.filter(
        (s) =>
          s.title.includes(params.keyword!) ||
          s.description.includes(params.keyword!)
      )
    : MOCK_SERVICES;
  return {
    items: filtered,
    total: filtered.length,
    page: 1,
    page_size: filtered.length,
  };
}

export function useServices(params: ServiceQueryParams = {}) {
  return useQuery({
    queryKey: ["services", params],
    queryFn: () =>
      fetchServices(params).catch(() => getMockServices(params)),
    placeholderData: getMockServices(params),
  });
}
