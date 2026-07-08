import { Link, Outlet, useNavigate } from "react-router";
import { useAuth } from "../lib/hooks/use-auth";

export default function Layout() {
  const { user, loading, clearAuth } = useAuth();
  const navigate = useNavigate();

  function handleLogout() {
    clearAuth();
    navigate("/");
  }

  return (
    <div className="min-h-screen flex flex-col">
      <header className="border-b border-gray-200 dark:border-gray-700">
        <nav className="max-w-5xl mx-auto flex items-center justify-between px-4 h-14">
          <Link to="/" className="font-bold text-lg">
            ORS
          </Link>
          {!loading && (
            <div className="flex items-center gap-4">
              {user ? (
                <>
                  <Link
                    to="/dashboard"
                    className="text-sm text-gray-600 dark:text-gray-300 hover:underline"
                  >
                    {user.name}
                  </Link>
                  <button
                    onClick={handleLogout}
                    className="text-sm text-red-500 hover:underline"
                  >
                    退出
                  </button>
                </>
              ) : (
                <>
                  <Link
                    to="/login"
                    className="text-sm text-gray-600 dark:text-gray-300 hover:underline"
                  >
                    登录
                  </Link>
                  <Link
                    to="/register"
                    className="text-sm text-white bg-blue-600 px-3 py-1.5 rounded hover:bg-blue-700"
                  >
                    注册
                  </Link>
                </>
              )}
            </div>
          )}
        </nav>
      </header>
      <main className="flex-1">
        <Outlet />
      </main>
    </div>
  );
}
