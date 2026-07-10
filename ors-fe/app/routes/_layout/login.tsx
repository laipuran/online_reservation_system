import { useState, type FormEvent } from "react";
import { useNavigate, Link } from "react-router";
import { useLogin } from "../../lib/hooks/use-mutations";
import { ApiError } from "../../lib/api/client";

export default function Login() {
  const loginMutation = useLogin();
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");

    if (!email.trim()) {
      setError("请输入邮箱");
      return;
    }
    if (!password) {
      setError("请输入密码");
      return;
    }

    try {
      const data = await loginMutation.mutateAsync({ email: email.trim(), password });
      navigate("/dashboard", { replace: true });
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("登录失败，请稍后重试");
      }
    }
  }

  return (
    <div className="max-w-sm mx-auto mt-20 px-4">
      <h1 className="text-2xl font-bold text-center mb-6 text-gray-900 dark:text-gray-100">登录</h1>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1 dark:text-gray-300">邮箱</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 dark:text-gray-100"
            placeholder="your@email.com"
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1 dark:text-gray-300">密码</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 dark:text-gray-100"
          />
        </div>
        {error && <p className="text-red-500 text-sm">{error}</p>}
        <button
          type="submit"
          disabled={loginMutation.isPending}
          className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 disabled:opacity-50"
        >
          {loginMutation.isPending ? "登录中..." : "登录"}
        </button>
      </form>
      <p className="text-sm text-center mt-4 text-gray-500 dark:text-gray-400">
        还没有账号？{" "}
        <Link to="/register" className="text-blue-600 dark:text-blue-400 hover:underline">
          注册
        </Link>
      </p>
    </div>
  );
}
