import { request } from "./client";

export interface Tag {
  id: number;
  name: string;
  created_at: string;
}

const MOCK_TAGS: Tag[] = [
  { id: 1, name: "放松", created_at: "2026-07-01T00:00:00Z" },
  { id: 2, name: "塑形", created_at: "2026-07-01T00:00:00Z" },
  { id: 3, name: "美白", created_at: "2026-07-01T00:00:00Z" },
  { id: 4, name: "健身", created_at: "2026-07-01T00:00:00Z" },
  { id: 5, name: "养生", created_at: "2026-07-01T00:00:00Z" },
  { id: 6, name: "护肤", created_at: "2026-07-01T00:00:00Z" },
  { id: 7, name: "按摩", created_at: "2026-07-01T00:00:00Z" },
  { id: 8, name: "瑜伽", created_at: "2026-07-01T00:00:00Z" },
  { id: 9, name: "美发", created_at: "2026-07-01T00:00:00Z" },
  { id: 10, name: "SPA", created_at: "2026-07-01T00:00:00Z" },
  { id: 11, name: "中医", created_at: "2026-07-01T00:00:00Z" },
  { id: 12, name: "减脂", created_at: "2026-07-01T00:00:00Z" },
  { id: 13, name: "抗衰", created_at: "2026-07-01T00:00:00Z" },
  { id: 14, name: "私教", created_at: "2026-07-01T00:00:00Z" },
  { id: 15, name: "理疗", created_at: "2026-07-01T00:00:00Z" },
];

const MAX_TAGS = 20;
const DISPLAY_LIMIT = 15;

export async function fetchTags(): Promise<Tag[]> {
  try {
    const tags = await request<Tag[]>("/tags");
    if (tags.length > MAX_TAGS) {
      return shuffleArray(tags).slice(0, DISPLAY_LIMIT);
    }
    return tags;
  } catch {
    return MOCK_TAGS;
  }
}

export function setUserInterests(tagIds: number[]): Promise<Tag[]> {
  return request<Tag[]>("/users/me/interests", {
    method: "PUT",
    body: JSON.stringify({ tag_ids: tagIds }),
  });
}

function shuffleArray<T>(arr: T[]): T[] {
  const copy = [...arr];
  for (let i = copy.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [copy[i], copy[j]] = [copy[j], copy[i]];
  }
  return copy;
}
