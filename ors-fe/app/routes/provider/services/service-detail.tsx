import { useParams, Link } from "react-router";
import { useMyProvider } from "../../../lib/hooks/use-provider";
import { useProviderServices } from "../../../lib/hooks/use-provider";

const STATUS_LABEL: Record<string, string> = {
  active: "已上架",
  inactive: "已下架",
  pending: "待审核",
  rejected: "已驳回",
};

const STATUS_CLASS: Record<string, string> = {
  active: "bg-green-100 text-green-700",
  inactive: "bg-gray-100 text-gray-500",
  pending: "bg-yellow-100 text-yellow-700",
  rejected: "bg-red-100 text-red-700",
};

export default function ProviderServiceDetailPage() {
  const { id } = useParams();
  const serviceId = Number(id);

  const { data: provider } = useMyProvider();
  const { data: servicesData, isLoading } = useProviderServices(
    provider?.id,
    {}
  );

  const service = servicesData?.items.find((s) => s.id === serviceId);

  if (isLoading) {
    return (
      <div className="flex justify-center py-20">
        <p className="text-gray-400">加载中...</p>
      </div>
    );
  }

  if (!service) {
    return (
      <div>
        <div className="mb-6">
          <Link
            to="/provider/services"
            className="text-sm text-blue-600 hover:underline"
          >
            &larr; 返回服务管理
          </Link>
        </div>
        <p className="text-gray-500">服务不存在</p>
      </div>
    );
  }

  return (
    <div>
      <div className="mb-6">
        <Link
          to="/provider/services"
          className="text-sm text-blue-600 hover:underline"
        >
          &larr; 返回服务管理
        </Link>
      </div>

      <h1 className="text-2xl font-bold mb-6">{service.title}</h1>

      <div className="border border-gray-200 dark:border-gray-700 rounded-lg p-6 space-y-4 max-w-xl">
        <div className="flex items-center gap-3">
          <span
            className={`inline-block text-xs px-2 py-0.5 rounded ${
              STATUS_CLASS[service.status] ?? ""
            }`}
          >
            {STATUS_LABEL[service.status] ?? service.status}
          </span>
        </div>

        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span className="text-gray-500">价格</span>
            <p className="font-medium">¥{service.price}</p>
          </div>
          <div>
            <span className="text-gray-500">服务时长</span>
            <p className="font-medium">{service.duration_minutes} 分钟</p>
          </div>
          <div>
            <span className="text-gray-500">分类</span>
            <p className="font-medium">{service.category.name}</p>
          </div>
          <div>
            <span className="text-gray-500">提供者</span>
            <p className="font-medium">{service.provider.business_name}</p>
          </div>
          <div>
            <span className="text-gray-500">评分</span>
            <p className="font-medium">
              {service.avg_rating.toFixed(1)} ({service.review_count} 条评价)
            </p>
          </div>
        </div>

        <div>
          <span className="text-sm text-gray-500">描述</span>
          <p className="text-sm mt-1">{service.description}</p>
        </div>

        {service.image_url && (
          <div>
            <span className="text-sm text-gray-500">图片</span>
            <img
              src={service.image_url}
              alt={service.title}
              className="mt-1 max-h-48 rounded"
            />
          </div>
        )}

        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span className="text-gray-500">创建时间</span>
            <p className="font-medium">
              {new Date(service.created_at).toLocaleString("zh-CN")}
            </p>
          </div>
          <div>
            <span className="text-gray-500">更新时间</span>
            <p className="font-medium">
              {new Date(service.updated_at).toLocaleString("zh-CN")}
            </p>
          </div>
        </div>
      </div>

      <div className="mt-6">
        <Link
          to={`/provider/services/${service.id}/edit`}
          className="inline-block bg-blue-600 text-white px-4 py-2 rounded text-sm hover:bg-blue-700"
        >
          修改服务
        </Link>
      </div>
    </div>
  );
}
