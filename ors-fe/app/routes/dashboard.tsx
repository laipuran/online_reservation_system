import { useEffect } from "react";
import { useNavigate, Link } from "react-router";
import { useQuery } from "@tanstack/react-query";
import { useAuth } from "../lib/hooks/use-auth";
import { fetchMyProvider } from "../lib/api/providers";

export default function Dashboard() {
  const { user, loading } = useAuth();
  const navigate = useNavigate();

  const providerQuery = useQuery({
    queryKey: ["my-provider-profile"],
    queryFn: async () => {
      try {
        return await fetchMyProvider();
      } catch {
        return null;
      }
    },
    enabled: !loading && user?.role === "provider",
    retry: false,
  });

  useEffect(() => {
    if (!loading && !user) {
      navigate("/login", { replace: true });
    }
  }, [user, loading, navigate]);

  useEffect(() => {
    if (
      !loading &&
      user?.role === "provider" &&
      !providerQuery.isLoading &&
      providerQuery.data === null
    ) {
      navigate("/complete-profile", { replace: true });
    }
  }, [user, loading, providerQuery.isLoading, providerQuery.data, navigate]);

  if (loading || providerQuery.isLoading) {
    return (
      <div className="flex items-center justify-center mt-20">
        <p className="text-gray-500">加载中...</p>
      </div>
    );
  }

  if (!user) {
    return null;
  }

  return (
    <div className="max-w-lg mx-auto mt-20 px-4">
      <h1 className="text-2xl font-bold mb-6">控制台</h1>
      <div className="border border-gray-200 dark:border-gray-700 rounded-lg p-6 space-y-3">
        <div>
          <span className="text-sm text-gray-500">昵称</span>
          <p className="font-medium">{user.name}</p>
        </div>
        <div>
          <span className="text-sm text-gray-500">邮箱</span>
          <p className="font-medium">{user.email}</p>
        </div>
        <div>
          <span className="text-sm text-gray-500">角色</span>
          <p className="font-medium">
            {user.role === "provider" ? "服务提供者" : "服务体验者"}
          </p>
        </div>
        {user.role === "provider" && providerQuery.data && (
          <div>
            <span className="text-sm text-gray-500">商家名称</span>
            <p className="font-medium">{providerQuery.data.business_name}</p>
          </div>
        )}
      </div>

      {user.role === "provider" && (
        <div className="mt-6">
          <Link
            to="/provider/services"
            className="inline-block bg-blue-600 text-white px-6 py-2 rounded text-sm hover:bg-blue-700"
          >
            进入服务商控制台
          </Link>
        </div>
      )}
    </div>
  );
}
