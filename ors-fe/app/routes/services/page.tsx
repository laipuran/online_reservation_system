import { useState, useMemo } from "react";
import { useSearchParams } from "react-router";
import { useCategories } from "../../lib/hooks/use-categories";
import { useServices } from "../../lib/hooks/use-services";
import { ServiceCard } from "../../lib/components/service-card";
import type { Route } from "./+types/page";

export function meta({}: Route.MetaArgs) {
  return [{ title: "自助预约项目 - ORS" }];
}

function shuffle<T>(arr: T[]): T[] {
  const a = [...arr];
  for (let i = a.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [a[i], a[j]] = [a[j], a[i]];
  }
  return a;
}

export default function ServicesPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const { data: categories = [] } = useCategories();
  const urlKeyword = searchParams.get("keyword") || "";
  const [keyword, setKeyword] = useState(urlKeyword);
  const { data: servicesData, isLoading } = useServices({
    keyword: keyword || undefined,
    page_size: 50,
  });
  const [selectedParentId, setSelectedParentId] = useState<number | null>(null);

  const parentCategories = useMemo(
    () => categories.filter((c) => c.parent_id == null),
    [categories]
  );

  const childIdsByParent = useMemo(() => {
    const map = new Map<number, number[]>();
    for (const p of parentCategories) {
      const ids = categories
        .filter((c) => c.parent_id === p.id)
        .map((c) => c.id);
      ids.push(p.id);
      map.set(p.id, ids);
    }
    return map;
  }, [categories, parentCategories]);

  const services = servicesData?.items ?? [];

  const filtered = useMemo(() => {
    if (selectedParentId === null) return services;
    const allowed = childIdsByParent.get(selectedParentId) ?? [selectedParentId];
    return services.filter((s) => allowed.includes(s.category.id));
  }, [services, selectedParentId, childIdsByParent]);

  const shuffled = useMemo(() => shuffle(filtered), [filtered]);

  function handleSearch(value: string) {
    setKeyword(value);
    setSearchParams(value ? { keyword: value } : {}, { replace: true });
  }

  return (
    <div className="max-w-5xl mx-auto px-4 py-8">
      <div className="flex items-start justify-between mb-2">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">自助预约项目</h1>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">发现优质服务，轻松预约</p>
        </div>
        <input
          type="text"
          placeholder="搜索服务..."
          value={keyword}
          onChange={(e) => handleSearch(e.target.value)}
          className="border border-gray-300 dark:border-gray-600 rounded-lg px-4 py-2 text-sm w-64 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
        />
      </div>

      <div className="flex gap-2 mt-6 mb-8 flex-wrap">
        <button
          onClick={() => setSelectedParentId(null)}
          className={`px-4 py-1.5 rounded-full text-sm font-medium transition-colors ${
            selectedParentId === null
              ? "bg-blue-600 text-white"
              : "bg-gray-100 text-gray-600 hover:bg-gray-200"
          }`}
        >
          全部项目
        </button>
        {parentCategories.map((p) => (
          <button
            key={p.id}
            onClick={() => setSelectedParentId(p.id)}
            className={`px-4 py-1.5 rounded-full text-sm font-medium transition-colors ${
              selectedParentId === p.id
                ? "bg-blue-600 text-white"
              : "bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-700"
            }`}
          >
            {p.name}
          </button>
        ))}
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center py-20">
          <p className="text-gray-400 dark:text-gray-500">加载中...</p>
        </div>
      ) : filtered.length === 0 ? (
        <div className="flex items-center justify-center py-20">
          <p className="text-gray-400 dark:text-gray-500">暂无服务</p>
        </div>
      ) : selectedParentId === null ? (
        <div className="grid grid-cols-3 gap-6">
          {shuffled.map((s) => (
            <ServiceCard
              key={s.id}
              id={s.id}
              title={s.title}
              description={s.description ?? ""}
              price={s.price}
              durationMinutes={s.duration_minutes}
              avgRating={s.avg_rating}
              imageUrl={s.image_url ?? ""}
              status={s.status ?? "active"}
            />
          ))}
        </div>
      ) : (
        <section>
          <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
            {categories.find((c) => c.id === selectedParentId)?.name ?? ""}
          </h2>
          <div className="grid grid-cols-3 gap-6">
            {filtered.map((s) => (
              <ServiceCard
                key={s.id}
                id={s.id}
                title={s.title}
                description={s.description ?? ""}
                price={s.price}
                durationMinutes={s.duration_minutes}
                avgRating={s.avg_rating}
                imageUrl={s.image_url ?? ""}
                status={s.status ?? "active"}
              />
            ))}
          </div>
        </section>
      )}
    </div>
  );
}
