import { useNavigate, useParams, Link } from "react-router";
import { useQuery, useMutation } from "@tanstack/react-query";
import { fetchServiceById } from "../../lib/api/services";
import { createReservation } from "../../lib/api/reservations";
import { useAuth } from "../../lib/hooks/use-auth";

export default function ConfirmPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const serviceId = Number(id);
  const { user, token, loading: authLoading } = useAuth();

  const { data: service, isLoading } = useQuery({
    queryKey: ["service", serviceId],
    queryFn: () => fetchServiceById(serviceId),
    enabled: !!serviceId,
  });

  const bookingMutation = useMutation({
    mutationFn: () => {
      const startTime = `${searchParams.get("date")}T${searchParams.get("time")}:00`;
      return createReservation({
        service_id: serviceId,
        start_time: startTime,
        note: searchParams.get("note") || undefined,
      });
    },
    onSuccess: () => {
      navigate("/dashboard");
    },
  });

  const searchParams = new URLSearchParams(
    typeof window !== "undefined" ? window.location.search : ""
  );
  const date = searchParams.get("date") || "";
  const time = searchParams.get("time") || "";
  const note = searchParams.get("note") || "";

  if (authLoading) {
    return (
      <div className="max-w-lg mx-auto px-4 py-20 text-center">
        <p className="text-gray-400">加载中...</p>
      </div>
    );
  }

  if (!token) {
    navigate(`/login?redirect=/services/${serviceId}`);
    return null;
  }

  if (isLoading) {
    return (
      <div className="max-w-lg mx-auto px-4 py-20">
        <div className="animate-pulse space-y-4">
          <div className="h-6 bg-gray-200 rounded w-48" />
          <div className="h-40 bg-gray-200 rounded-xl" />
          <div className="h-10 bg-gray-200 rounded-lg" />
        </div>
      </div>
    );
  }

  if (!service) {
    return (
      <div className="max-w-lg mx-auto px-4 py-20 text-center">
        <p className="text-gray-500">服务不存在</p>
        <Link to="/services" className="text-blue-600 hover:underline mt-4 inline-block">&larr; 返回服务列表</Link>
      </div>
    );
  }

  if (!date || !time) {
    return (
      <div className="max-w-lg mx-auto px-4 py-20 text-center">
        <p className="text-amber-600 mb-2">缺少预约时间信息</p>
        <p className="text-sm text-gray-400 mb-4">请从服务详情页选择时间后再来确认</p>
        <Link to={`/services/${serviceId}`} className="text-blue-600 hover:underline">&larr; 返回服务详情</Link>
      </div>
    );
  }

  const startDateTime = new Date(`${date}T${time}:00`);
  const endDateTime = new Date(startDateTime.getTime() + service.duration_minutes * 60000);
  const formatTime = (d: Date) =>
    d.toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" });
  const formatDate = (d: Date) =>
    d.toLocaleDateString("zh-CN", { year: "numeric", month: "long", day: "numeric", weekday: "short" });

  const handlePay = () => {
    bookingMutation.mutate();
  };

  return (
    <div className="max-w-lg mx-auto px-4 py-8">
      <Link to={`/services/${serviceId}`} className="text-sm text-blue-600 hover:underline">
        &larr; 返回服务详情
      </Link>

      <h1 className="text-xl font-bold text-gray-900 mt-4 mb-6">确认预约信息</h1>

      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <div className="p-5 space-y-4">
          <div className="flex items-start gap-4">
            {service.image_url ? (
              <img src={service.image_url} alt={service.title} className="w-20 h-20 rounded-lg object-cover shrink-0" />
            ) : (
              <div className="w-20 h-20 rounded-lg bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center text-2xl shrink-0">
                💆
              </div>
            )}
            <div>
              <h2 className="font-semibold text-gray-900">{service.title}</h2>
              <p className="text-sm text-gray-500 mt-0.5">{service.provider.business_name}</p>
              <p className="text-2xl font-bold text-red-500 mt-1">
                <span className="text-sm">¥</span>{service.price}
              </p>
            </div>
          </div>

          <div className="border-t border-gray-100 pt-4 space-y-3 text-sm">
            <div className="flex justify-between">
              <span className="text-gray-500">预约日期</span>
              <span className="text-gray-900 font-medium">{formatDate(startDateTime)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">服务时段</span>
              <span className="text-gray-900 font-medium">
                {formatTime(startDateTime)} — {formatTime(endDateTime)}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-500">服务时长</span>
              <span className="text-gray-900 font-medium">{service.duration_minutes} 分钟</span>
            </div>
            {note && (
              <div className="flex justify-between">
                <span className="text-gray-500">备注</span>
                <span className="text-gray-900 font-medium max-w-[200px] text-right">{note}</span>
              </div>
            )}
            <div className="flex justify-between border-t border-gray-100 pt-3">
              <span className="text-gray-900 font-semibold">合计</span>
              <span className="text-xl font-bold text-red-500">¥{service.price}</span>
            </div>
          </div>
        </div>
      </div>

      <button
        onClick={handlePay}
        disabled={bookingMutation.isPending}
        className="w-full mt-6 bg-gradient-to-r from-red-500 to-red-600 text-white font-semibold py-3 rounded-xl text-lg hover:from-red-600 hover:to-red-700 transition-all disabled:opacity-50 shadow-sm"
      >
        {bookingMutation.isPending ? "提交中..." : "确认并支付"}
      </button>

      <p className="text-xs text-gray-400 text-center mt-3">
        * 模拟支付，不会产生真实扣款
      </p>
    </div>
  );
}
