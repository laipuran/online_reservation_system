import { useEffect } from "react";
import { useNavigate } from "react-router";
import { useAuth } from "../lib/hooks/use-auth";

export default function Dashboard() {
  const { user, loading } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (!loading && !user) {
      navigate("/login", { replace: true });
    }
  }, [user, loading, navigate]);

  if (loading) {
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
          <p className="font-medium">{user.role}</p>
        </div>
      </div>
    </div>
  );
}
