import { useServices } from "../lib/hooks/use-services";
import { ServiceCard } from "../lib/components/service-card";
import type { Route } from "./+types/home";

export function meta({}: Route.MetaArgs) {
  return [{ title: "ORS - 在线预约系统" }];
}

export default function Home() {
  const { data, isLoading } = useServices();
  const services = data?.items ?? [];

  return (
    <div className="max-w-5xl mx-auto px-4 py-10">
      <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-2">推荐服务</h1>
      <p className="text-sm text-gray-500 dark:text-gray-400 mb-8">
        发现附近的热门服务，立即预约
      </p>

      {isLoading ? (
        <div className="flex items-center justify-center py-20">
          <p className="text-gray-400 dark:text-gray-500">加载中...</p>
        </div>
      ) : services.length === 0 ? (
        <div className="flex items-center justify-center py-20">
          <p className="text-gray-400 dark:text-gray-500">暂无服务</p>
        </div>
      ) : (
        <div className="grid grid-cols-3 gap-6">
          {services.map((s) => (
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
      )}
    </div>
  );
}
