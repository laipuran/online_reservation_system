import { useState, useEffect, type FormEvent } from "react";
import { useParams, useNavigate, Link } from "react-router";
import { useMyProvider, useProviderServices, useUpdateService } from "../../../lib/hooks/use-provider";
import { useCategories } from "../../../lib/hooks/use-categories";
import { ApiError } from "../../../lib/api/client";

export default function EditServicePage() {
  const { id } = useParams();
  const serviceId = Number(id);
  const navigate = useNavigate();
  const updateMutation = useUpdateService();
  const { data: categories } = useCategories();
  const { data: provider } = useMyProvider();
  const { data: servicesData, isLoading: loadingService } = useProviderServices(
    provider?.id,
    {}
  );

  const service = servicesData?.items.find((s) => s.id === serviceId);

  const [categoryId, setCategoryId] = useState("");
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [price, setPrice] = useState("");
  const [durationMinutes, setDurationMinutes] = useState("");
  const [imageUrl, setImageUrl] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    if (service) {
      setCategoryId(String(service.category.id));
      setTitle(service.title);
      setDescription(service.description ?? "");
      setPrice(String(service.price));
      setDurationMinutes(String(service.duration_minutes));
      setImageUrl(service.image_url ?? "");
    }
  }, [service]);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");

    if (!categoryId) {
      setError("请选择服务分类");
      return;
    }
    if (!title.trim()) {
      setError("请输入服务标题");
      return;
    }
    if (!price || Number(price) < 0) {
      setError("价格不能小于 0");
      return;
    }
    if (!durationMinutes || Number(durationMinutes) <= 0) {
      setError("服务时长必须大于 0");
      return;
    }

    try {
      await updateMutation.mutateAsync({
        id: serviceId,
        data: {
          category_id: Number(categoryId),
          title: title.trim(),
          description: description.trim() || undefined,
          price: Number(price),
          duration_minutes: Number(durationMinutes),
          image_url: imageUrl.trim() || undefined,
        },
      });
      navigate("/provider/services");
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("修改失败，请稍后重试");
      }
    }
  }

  if (loadingService) {
    return (
      <div className="flex justify-center py-20">
        <p className="text-gray-400 dark:text-gray-500">加载中...</p>
      </div>
    );
  }

  if (!service) {
    return (
      <div>
        <div className="mb-6">
          <Link
            to="/provider/services"
            className="text-sm text-blue-600 dark:text-blue-400 hover:underline"
          >
            &larr; 返回服务管理
          </Link>
        </div>
        <p className="text-gray-500 dark:text-gray-400">服务不存在</p>
      </div>
    );
  }

  return (
    <div className="max-w-lg">
      <div className="mb-6">
          <Link
            to="/provider/services"
            className="text-sm text-blue-600 dark:text-blue-400 hover:underline"
          >
            &larr; 返回服务管理
          </Link>
        </div>

      <h1 className="text-2xl font-bold mb-6 text-gray-900 dark:text-gray-100">修改服务</h1>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1">
            服务分类 <span className="text-red-500">*</span>
          </label>
          <select
            value={categoryId}
            onChange={(e) => setCategoryId(e.target.value)}
            className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 text-sm"
          >
            <option value="">请选择分类</option>
            {(categories ?? []).map((c) => (
              <option key={c.id} value={c.id}>
                {c.name}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">
            服务标题 <span className="text-red-500">*</span>
          </label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 text-sm"
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">服务描述</label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 text-sm"
            rows={3}
          />
        </div>

        <div className="flex gap-4">
          <div className="flex-1">
            <label className="block text-sm font-medium mb-1">
              价格 (¥) <span className="text-red-500">*</span>
            </label>
            <input
              type="number"
              min="0"
              value={price}
              onChange={(e) => setPrice(e.target.value)}
              className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 text-sm"
            />
          </div>
          <div className="flex-1">
            <label className="block text-sm font-medium mb-1">
              服务时长 (分钟) <span className="text-red-500">*</span>
            </label>
            <input
              type="number"
              min="1"
              value={durationMinutes}
              onChange={(e) => setDurationMinutes(e.target.value)}
              className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 text-sm"
            />
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">图片 URL</label>
          <input
            type="url"
            value={imageUrl}
            onChange={(e) => setImageUrl(e.target.value)}
            className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 text-sm"
          />
        </div>

        {error && <p className="text-red-500 text-sm">{error}</p>}

        <div className="flex gap-3 pt-2">
          <button
            type="submit"
            disabled={updateMutation.isPending}
            className="bg-blue-600 text-white px-6 py-2 rounded text-sm hover:bg-blue-700 disabled:opacity-50"
          >
            {updateMutation.isPending ? "保存中..." : "保存修改"}
          </button>
          <Link
            to="/provider/services"
            className="px-6 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded hover:bg-gray-50 dark:hover:bg-gray-800"
          >
            取消
          </Link>
        </div>
      </form>
    </div>
  );
}
