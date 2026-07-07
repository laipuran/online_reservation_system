import { Link } from "react-router";

interface ServiceCardProps {
  id: number;
  title: string;
  description: string;
  price: number;
  durationMinutes: number;
  avgRating: number;
}

export function ServiceCard({
  id,
  title,
  description,
  price,
  durationMinutes,
  avgRating,
}: ServiceCardProps) {
  return (
    <Link
      to={`/services/${id}`}
      className="flex flex-col rounded-xl border border-gray-200 bg-white p-5 
                 min-h-[260px] max-h-[340px] hover:shadow-md transition-shadow 
                 hover:border-gray-300"
    >
      <h3 className="font-semibold text-base text-gray-900 leading-snug">
        {title}
      </h3>

      <p className="text-sm text-gray-500 mt-2 line-clamp-2 leading-relaxed">
        {description}
      </p>

      <div className="flex items-center gap-1 mt-3">
        <span className="text-amber-400 text-xs">★</span>
        <span className="text-amber-500 text-sm font-medium">
          综合评分 {avgRating.toFixed(1)}
        </span>
      </div>

      <div className="mt-auto flex items-end justify-between pt-4">
        <span className="text-red-500 font-semibold text-base leading-none">
          <span className="text-sm">￥</span>
          {price}
        </span>
        <span className="text-xs text-gray-400 leading-none">
          {durationMinutes} 分钟
        </span>
      </div>
    </Link>
  );
}
