import { useState } from "react";
import { Link, useNavigate } from "react-router";
import { useMyProvider, useProviderServices, useUpdateServiceStatus } from "../../../lib/hooks/use-provider";
import { STATUS_CONFIG } from "../../../lib/status";

export default function ProviderServicesPage() {
  const navigate = useNavigate();
  const [keyword, setKeyword] = useState("");
  const [page, setPage] = useState(1);
  const [deletingId, setDeletingId] = useState<number | null>(null);
  const pageSize = 10;

  const { data: provider, isLoading: loadingProvider } = useMyProvider();
  const providerId = provider?.id;

  const { data: servicesData, isLoading: loadingServices } = useProviderServices(
    providerId,
    { keyword: keyword || undefined, page, page_size: pageSize }
  );

  const deactivateMutation = useUpdateServiceStatus();

  const services = servicesData?.items ?? [];
  const total = servicesData?.total ?? 0;
  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  function handleDelete(id: number) {
    deactivateMutation.mutate(
      { id, status: "inactive" },
      {
        onSuccess: () => {
          setDeletingId(null);
        },
      }
    );
  }

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">服务管理</h1>

      <div className="flex items-center gap-4 mb-6">
        <input
          type="text"
          placeholder="搜索服务..."
          value={keyword}
          onChange={(e) => {
            setKeyword(e.target.value);
            setPage(1);
          }}
          className="flex-1 border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 text-sm"
        />
        <Link
          to="/provider/services/new"
          className="shrink-0 bg-blue-600 text-white px-4 py-2 rounded text-sm hover:bg-blue-700"
        >
          新建服务
        </Link>
      </div>

      {loadingProvider || loadingServices ? (
        <div className="flex justify-center py-20">
          <p className="text-gray-400 dark:text-gray-500">加载中...</p>
        </div>
      ) : services.length === 0 ? (
        <div className="flex justify-center py-20">
          <p className="text-gray-400 dark:text-gray-500">暂无服务</p>
        </div>
      ) : (
        <>
          <div className="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-gray-50 dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
                  <th className="text-left px-4 py-3 font-medium text-gray-600 dark:text-gray-400">标题</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-600 dark:text-gray-400">价格</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-600 dark:text-gray-400">时长</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-600 dark:text-gray-400">状态</th>
                  <th className="text-right px-4 py-3 font-medium text-gray-600 dark:text-gray-400">操作</th>
                </tr>
              </thead>
              <tbody>
                {services.map((s) => (
                  <tr
                    key={s.id}
                    className="border-b border-gray-100 dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-gray-800/50"
                  >
                    <td className="px-4 py-3 font-medium">{s.title}</td>
                    <td className="px-4 py-3">¥{s.price}</td>
                    <td className="px-4 py-3">{s.duration_minutes} 分钟</td>
                    <td className="px-4 py-3">
                      <span
                        className={`inline-block text-xs px-2 py-0.5 rounded ${
                          (STATUS_CONFIG[s.status] ?? STATUS_CONFIG.active).className
                        }`}
                      >
                        {(STATUS_CONFIG[s.status] ?? STATUS_CONFIG.active).label}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-right space-x-2">
                      <Link
                        to={`/provider/services/${s.id}`}
                        className="text-blue-600 dark:text-blue-400 hover:underline text-xs"
                      >
                        详情
                      </Link>
                      <Link
                        to={`/provider/services/${s.id}/edit`}
                        className="text-blue-600 dark:text-blue-400 hover:underline text-xs"
                      >
                        修改
                      </Link>
                      {s.status !== "inactive" && (
                        <button
                          onClick={() => setDeletingId(s.id)}
                          className="text-red-500 hover:underline text-xs"
                        >
                          删除
                        </button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {totalPages > 1 && (
            <div className="flex items-center justify-center gap-2 mt-6">
              <button
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page <= 1}
                className="px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-800"
              >
                上一页
              </button>
              <span className="text-sm text-gray-500 dark:text-gray-400">
                {page} / {totalPages}
              </span>
              <button
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={page >= totalPages}
                className="px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-800"
              >
                下一页
              </button>
            </div>
          )}
        </>
      )}

      {deletingId !== null && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white dark:bg-gray-800 rounded-lg p-6 max-w-sm mx-4 shadow-xl">
            <h3 className="text-lg font-semibold mb-2">确认下架</h3>
            <p className="text-sm text-gray-500 dark:text-gray-400 mb-6">确定要下架此服务吗？下架后用户将无法预约。</p>
            <div className="flex justify-end gap-3">
              <button
                onClick={() => setDeletingId(null)}
                className="px-4 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded hover:bg-gray-50 dark:hover:bg-gray-700"
              >
                取消
              </button>
              <button
                onClick={() => handleDelete(deletingId)}
                disabled={deactivateMutation.isPending}
                className="px-4 py-2 text-sm bg-red-600 text-white rounded hover:bg-red-700 disabled:opacity-50"
              >
                {deactivateMutation.isPending ? "处理中..." : "确认下架"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
