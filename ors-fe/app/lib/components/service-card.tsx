import { Link } from "react-router";
import { STATUS_CONFIG } from "../status";
import type { ServiceStatus } from "../status";

interface ServiceCardProps {
  id: number;
  title: string;
  description: string;
  price: number;
  durationMinutes: number;
  avgRating: number;
  imageUrl: string;
  status: ServiceStatus;
}

export function ServiceCard({
  id,
  title,
  description,
  price,
  durationMinutes,
  avgRating,
  imageUrl,
  status,
}: ServiceCardProps) {
  const statusCfg = STATUS_CONFIG[status] ?? STATUS_CONFIG.active;

  return (
    <Link
      to={`/services/${id}`}
      className="flex flex-col rounded-xl border border-gray-200 bg-white
                 min-h-[260px] max-h-[380px] hover:shadow-md transition-shadow
                 hover:border-gray-300 overflow-hidden
                 dark:border-gray-700 dark:bg-gray-800 dark:hover:border-gray-600"
    >
      {imageUrl ? (
        <img
          src={imageUrl}
          alt={title}
          className="w-full h-32 object-cover"
        />
      ) : (
        <div className="w-full h-32 bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
          <span className="text-gray-400 dark:text-gray-500 text-sm">{title}</span>
        </div>
      )}

      <div className="flex flex-col flex-1 p-4 pt-3">
        <div className="flex items-start justify-between gap-2">
          <h3 className="font-semibold text-base text-gray-900 dark:text-gray-100 leading-snug">
            {title}
          </h3>
          <span
            className={`shrink-0 text-xs px-1.5 py-0.5 rounded ${statusCfg.className}`}
          >
            {statusCfg.label}
          </span>
        </div>

        <p className="text-sm text-gray-500 dark:text-gray-400 mt-1.5 line-clamp-2 leading-relaxed">
          {description}
        </p>

        <div className="flex items-center gap-1 mt-2">
          <span className="text-amber-400 text-xs">★</span>
          <span className="text-amber-500 text-sm font-medium">
            综合评分 {avgRating.toFixed(1)}
          </span>
        </div>

        <div className="mt-auto flex items-center justify-between pt-3">
          <span className="text-red-500 font-bold text-lg leading-none">
            <span className="text-sm">￥</span>
            {price}
          </span>
          <span className="text-sm text-gray-400">
            {durationMinutes} 分钟
          </span>
        </div>
      </div>
    </Link>
  );
}
