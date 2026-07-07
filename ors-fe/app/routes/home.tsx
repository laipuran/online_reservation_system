import { Link } from "react-router";
import type { Route } from "./+types/home";
import { useAuth } from "../lib/hooks/use-auth";

export function meta({}: Route.MetaArgs) {
  return [{ title: "ORS - 在线预约系统" }];
}

export default function Home() {
  const { user, loading } = useAuth();

  if (loading) {
    return null;
  }

  return (
    <div className="flex flex-col items-center justify-center mt-32 px-4 text-center">
      <h1 className="text-4xl font-bold mb-4">在线预约系统</h1>
      <p className="text-gray-500 dark:text-gray-400 mb-8 max-w-md">
        欢迎使用在线预约系统，请登录或注册以继续使用。
      </p>
      {user ? (
        <Link
          to="/dashboard"
          className="bg-blue-600 text-white px-6 py-2.5 rounded hover:bg-blue-700"
        >
          进入控制台
        </Link>
      ) : (
        <div className="flex gap-4">
          <Link
            to="/login"
            className="bg-blue-600 text-white px-6 py-2.5 rounded hover:bg-blue-700"
          >
            登录
          </Link>
          <Link
            to="/register"
            className="border border-gray-300 dark:border-gray-600 px-6 py-2.5 rounded hover:bg-gray-50 dark:hover:bg-gray-800"
          >
            注册
          </Link>
        </div>
      )}
    </div>
  );
}
