import { request } from "./client";

export interface CategoryItem {
  id: number;
  name: string;
  description?: string;
  parent_id?: number;
  created_at: string;
}

export function fetchCategories(): Promise<CategoryItem[]> {
  return request<CategoryItem[]>("/categories");
}
