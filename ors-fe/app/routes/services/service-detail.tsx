import { useState } from "react";
import { useParams, Link, useNavigate } from "react-router";
import { useQuery } from "@tanstack/react-query";
import { fetchServiceById, fetchServiceTags } from "../../lib/api/services";
import { fetchProvider } from "../../lib/api/providers";
import {
  fetchServiceReviews,
} from "../../lib/api/reviews";
import { useAuth } from "../../lib/hooks/use-auth";

function StarRating({ rating, size = "sm" }: { rating: number; size?: "sm" | "md" | "lg" }) {
  const sizes = { sm: "text-sm", md: "text-base", lg: "text-xl" };
  const stars = [];
  for (let i = 1; i <= 5; i++) {
    if (rating >= i) stars.push("★");
    else if (rating >= i - 0.5) stars.push("⯪");
    else stars.push("☆");
  }
  return (
    <span className={`text-amber-400 ${sizes[size]} tracking-tight`}>
      {stars.join("")}
    </span>
  );
}

function ProviderCard({ providerId }: { providerId: number }) {
  const { data: provider } = useQuery({
    queryKey: ["provider", providerId],
    queryFn: () => fetchProvider(providerId),
    enabled: !!providerId,
  });
  if (!provider) return <div className="h-20 bg-gray-50 rounded-xl animate-pulse" />;
  const initial = provider.business_name.charAt(0);
  const colorIndex = providerId % 6;
  const colors = ["bg-blue-500", "bg-emerald-500", "bg-violet-500", "bg-rose-500", "bg-amber-500", "bg-cyan-500"];
  return (
    <div className="group relative cursor-pointer">
      <div className="flex items-center gap-4 p-4 bg-white rounded-xl border border-gray-200 group-hover:rounded-b-none group-hover:border-gray-300 transition-all">
        {provider.logo_url ? (
          <img src={provider.logo_url} alt={provider.business_name} className="w-12 h-12 rounded-full object-cover shrink-0" />
        ) : (
          <div className={`w-12 h-12 ${colors[colorIndex]} rounded-full flex items-center justify-center text-white font-bold text-lg shrink-0`}>
            {initial}
          </div>
        )}
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <h2 className="font-semibold text-gray-900 truncate">{provider.business_name}</h2>
          </div>
          <div className="flex items-center gap-2 mt-0.5">
            <StarRating rating={4.5} size="sm" />
            <span className="text-xs text-gray-400">▼ 查看商家详情</span>
          </div>
        </div>
      </div>
      <div className="absolute left-0 right-0 top-full z-20 bg-white dark:bg-gray-800 border border-t-0 border-gray-300 dark:border-gray-600 rounded-b-xl p-5 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all shadow-lg">
        <div className="flex gap-4">
          {provider.logo_url ? (
            <img src={provider.logo_url} alt={provider.business_name} className="w-16 h-16 rounded-xl object-cover shrink-0" />
          ) : (
            <div className={`w-16 h-16 ${colors[colorIndex]} rounded-xl flex items-center justify-center text-white font-bold text-2xl shrink-0`}>
              {initial}
            </div>
          )}
          <div className="flex-1 space-y-2 text-sm">
            <p className="text-gray-600">{provider.description}</p>
            <div className="grid grid-cols-2 gap-x-6 gap-y-1.5 text-gray-500">
              <span>📍 {provider.address}</span>
              <span>📞 {provider.phone}</span>
              <span>✉️ {provider.email}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function ServiceDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const serviceId = Number(id);
  const { user, token, loading: authLoading } = useAuth();

  const [bookDate, setBookDate] = useState("");
  const [bookTime, setBookTime] = useState("");
  const [note, setNote] = useState("");
  const [reviewPage, setReviewPage] = useState(1);
  const [keyword, setKeyword] = useState("");

  const { data: service, isLoading: serviceLoading } = useQuery({
    queryKey: ["service", serviceId],
    queryFn: () => fetchServiceById(serviceId),
    enabled: !!serviceId,
  });

  const { data: tags = [] } = useQuery({
    queryKey: ["service-tags", serviceId],
    queryFn: () => fetchServiceTags(serviceId),
    enabled: !!serviceId,
  });

  const { data: reviewsData, isLoading: reviewsLoading } = useQuery({
    queryKey: ["service-reviews", serviceId, { page: reviewPage, page_size: 5 }],
    queryFn: () => fetchServiceReviews(serviceId, { page: reviewPage, page_size: 5 }),
    enabled: !!serviceId,
  });

  if (serviceLoading) {
    return (
      <div className="max-w-5xl mx-auto px-4 py-8">
        <div className="animate-pulse space-y-6">
          <div className="h-4 bg-gray-200 rounded w-48" />
          <div className="h-20 bg-gray-200 rounded-xl" />
          <div className="h-80 bg-gray-200 rounded-xl" />
        </div>
      </div>
    );
  }

  if (!service) {
    return (
      <div className="max-w-5xl mx-auto px-4 py-20 text-center">
        <p className="text-gray-500 text-lg">服务不存在</p>
        <Link to="/services" className="text-blue-600 hover:underline mt-4 inline-block">&larr; 返回服务列表</Link>
      </div>
    );
  }

  const reviews = reviewsData?.items ?? [];
  const reviewsTotal = reviewsData?.total ?? reviews.length;
  const reviewsPageSize = 5;
  const totalReviewPages = Math.max(1, Math.ceil(reviewsTotal / reviewsPageSize));

  const handleBooking = () => {
    if (!token) {
      navigate(`/login?redirect=/services/${serviceId}`);
      return;
    }
    if (!bookDate || !bookTime) {
      alert("请选择预约日期和时间");
      return;
    }
    const selected = new Date(`${bookDate}T${bookTime}:00`);
    if (selected <= new Date()) {
      alert("预约时间必须在当前时间之后");
      return;
    }
    const params = new URLSearchParams({ date: bookDate, time: bookTime });
    if (note) params.set("note", note);
    navigate(`/services/${serviceId}/confirm?${params.toString()}`);
  };

  const handleSearch = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && keyword.trim()) {
      navigate(`/services?keyword=${encodeURIComponent(keyword.trim())}`);
    }
  };

  const now = new Date();
  const minDate = now.toISOString().split("T")[0];
  const endTime = bookDate && bookTime
    ? new Date(`${bookDate}T${bookTime}:00`).getTime() + service.duration_minutes * 60000
    : null;

  return (
    <div className="max-w-5xl mx-auto px-4 py-6">
      <nav className="text-sm text-gray-400 mb-4">
        <Link to="/" className="hover:text-blue-600">首页</Link>
        <span className="mx-2">&gt;</span>
        <Link to="/services" className="hover:text-blue-600">预约项目</Link>
        <span className="mx-2">&gt;</span>
        <span className="text-gray-600">{service.title}</span>
      </nav>

      <div className="flex gap-4 mb-8">
        <div className="flex-[7]">
          <ProviderCard providerId={service.provider.id} />
        </div>
        <div className="flex-[3]">
          <input
            type="text"
            placeholder="搜索其他服务..."
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            onKeyDown={handleSearch}
            className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white dark:bg-gray-800 dark:text-gray-100"
          />
        </div>
      </div>

      <div className="flex gap-8">
        <div className="flex-[7] space-y-8">
          {service.image_url ? (
            <img src={service.image_url} alt={service.title} className="w-full h-72 object-cover rounded-xl" />
          ) : (
            <div className="w-full h-72 rounded-xl bg-gradient-to-br from-blue-50 via-sky-100 to-indigo-100 flex items-center justify-center">
              <div className="text-center">
                <div className="text-5xl mb-2">💆</div>
                <p className="text-gray-400 text-sm">{service.title}</p>
              </div>
            </div>
          )}

          <div>
            <div className="flex items-start justify-between gap-4 mb-4">
              <div>
                <h1 className="text-2xl font-bold text-gray-900">{service.title}</h1>
                <div className="flex items-center gap-3 mt-2">
                  <StarRating rating={service.avg_rating} size="md" />
                  <span className="text-sm text-gray-500">
                    {service.avg_rating.toFixed(1)} ({service.review_count} 条评价)
                  </span>
                  <span className="text-xs px-2 py-0.5 rounded bg-green-100 text-green-700">
                    {service.status === "active" ? "可预约" : service.status}
                  </span>
                </div>
              </div>
              <div className="text-right">
                <p className="text-3xl font-bold text-red-500">
                  <span className="text-lg">¥</span>{service.price}
                </p>
                <p className="text-sm text-gray-400">{service.duration_minutes} 分钟</p>
              </div>
            </div>
            <div className="flex gap-2 flex-wrap">
              {tags.map((tag) => (
                <span key={tag.id} className="text-xs px-3 py-1 rounded-full bg-gray-100 text-gray-600">
                  {tag.name}
                </span>
              ))}
            </div>
          </div>

          <div>
            <h3 className="font-semibold text-gray-900 mb-2">服务详情</h3>
            <p className="text-sm text-gray-600 leading-relaxed">{service.description}</p>
          </div>

            <div>
              <h3 className="font-semibold text-gray-900 mb-4">用户评价</h3>
              {reviewsTotal > 0 && (
                <div className="flex gap-8 mb-6 p-4 bg-gray-50 rounded-xl">
                  <div className="text-center shrink-0">
                    <p className="text-4xl font-bold text-gray-900">{service.avg_rating.toFixed(1)}</p>
                    <StarRating rating={service.avg_rating} size="sm" />
                    <p className="text-xs text-gray-400 mt-1">{service.review_count} 条评价</p>
                  </div>
                </div>
              )}

            {reviewsLoading ? (
              <div className="space-y-3">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="h-20 bg-gray-50 rounded-lg animate-pulse" />
                ))}
              </div>
            ) : reviews.length === 0 ? (
              <p className="text-sm text-gray-400 py-6 text-center">暂无评价</p>
            ) : (
              <div className="space-y-4">
                {reviews.map((review) => (
                  <div key={review.id} className="border-b border-gray-100 pb-4">
                    <div className="flex items-center gap-3 mb-1.5">
                      <div className="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center text-sm text-gray-500 font-medium shrink-0">
                        ユ
                      </div>
                      <div>
                        <p className="text-sm font-medium text-gray-900">用户 #{review.user_id}</p>
                        <div className="flex items-center gap-2">
                          <StarRating rating={review.rating} size="sm" />
                          <span className="text-xs text-gray-400">
                            {new Date(review.created_at).toLocaleDateString("zh-CN")}
                          </span>
                        </div>
                      </div>
                    </div>
                    <p className="text-sm text-gray-600 ml-11">{review.comment || ""}</p>
                  </div>
                ))}
                {totalReviewPages > 1 && (
                  <div className="flex items-center justify-center gap-2 pt-2">
                    <button
                      onClick={() => setReviewPage((p) => Math.max(1, p - 1))}
                      disabled={reviewPage <= 1}
                      className="text-sm px-3 py-1 rounded border border-gray-200 disabled:opacity-30 hover:bg-gray-50"
                    >
                      上一页
                    </button>
                    <span className="text-sm text-gray-400">
                      {reviewPage} / {totalReviewPages}
                    </span>
                    <button
                      onClick={() => setReviewPage((p) => Math.min(totalReviewPages, p + 1))}
                      disabled={reviewPage >= totalReviewPages}
                      className="text-sm px-3 py-1 rounded border border-gray-200 disabled:opacity-30 hover:bg-gray-50"
                    >
                      下一页
                    </button>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>

        <div className="flex-[3]">
          <div className="sticky top-24 space-y-4">
            <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-5 shadow-sm">
              <p className="text-3xl font-bold text-red-500 mb-1">
                <span className="text-lg">¥</span>{service.price}
              </p>
              <p className="text-sm text-gray-400 mb-4">{service.duration_minutes} 分钟</p>

              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">选择日期</label>
              <input
                type="date"
                value={bookDate}
                onChange={(e) => setBookDate(e.target.value)}
                min={minDate}
                className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm mb-3 focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-800 dark:text-gray-100"
              />

              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">选择时间</label>
              <input
                type="time"
                value={bookTime}
                onChange={(e) => setBookTime(e.target.value)}
                className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm mb-3 focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-800 dark:text-gray-100"
              />

              {bookDate && bookTime && endTime && (
                <p className="text-xs text-gray-500 mb-3">
                  预约时段：{bookTime} - {new Date(endTime).toTimeString().slice(0, 5)}（{service.duration_minutes} 分钟）
                </p>
              )}

              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">备注</label>
              <textarea
                value={note}
                onChange={(e) => setNote(e.target.value)}
                placeholder="如有特殊需求请备注..."
                rows={3}
                className="w-full border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 text-sm mb-4 focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none bg-white dark:bg-gray-800 dark:text-gray-100"
              />

              <button
                onClick={handleBooking}
                className="w-full bg-gradient-to-r from-blue-500 to-blue-600 text-white font-semibold py-2.5 rounded-lg hover:from-blue-600 hover:to-blue-700 transition-all"
              >
                立即预约
              </button>
            </div>

            <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-4 text-sm text-gray-500 dark:text-gray-400">
              <p className="font-medium text-gray-700 dark:text-gray-200 mb-1">{service.provider.business_name}</p>
              <div className="flex items-center gap-1">
                <StarRating rating={service.avg_rating} size="sm" />
                <span className="text-xs">{service.avg_rating.toFixed(1)}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
