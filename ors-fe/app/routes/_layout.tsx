import { Link, Outlet, useLocation, useNavigate } from "react-router";
import { useAuth } from "../lib/hooks/use-auth";
import { NotificationBell } from "../lib/components/notification-bell";

const TABS = [
  { label: "首页", to: "/" },
  { label: "预约项目", to: "/services" },
  { label: "我的预约", to: "/my-reservations" },
];

export default function Layout() {
  const { user, loading, clearAuth } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  function isActive(path: string) {
    if (path === "/") return location.pathname === "/";
    return location.pathname.startsWith(path);
  }

  function handleLogout() {
    clearAuth();
    navigate("/");
  }

  return (
    <div className="min-h-screen flex flex-col">
      <header className="border-b border-gray-200 dark:border-gray-700">
        <nav className="max-w-5xl mx-auto flex items-center justify-between px-4 h-14">
          <div className="flex items-center gap-6">
            <Link to="/" className="font-bold text-lg shrink-0">
              ORS
            </Link>
            <div className="flex items-center gap-1">
              {TABS.map((tab) => (
                <Link
                  key={tab.to}
                  to={tab.to}
                  className={`px-3 py-1.5 text-sm rounded-md transition-colors ${
                    isActive(tab.to)
                      ? "bg-blue-50 text-blue-700 font-medium"
                      : "text-gray-600 hover:text-gray-900 hover:bg-gray-50"
                  }`}
                >
                  {tab.label}
                </Link>
              ))}
            </div>
          </div>
          {!loading && (
            <div className="flex items-center gap-4 shrink-0">
              {user ? (
                <>
                  <Link
                    to="/dashboard"
                    className="text-sm text-gray-600 hover:underline"
                  >
                    {user.name}
                  </Link>
                  <NotificationBell />
                  <div className="w-7 h-7 rounded-full bg-blue-500 flex items-center justify-center text-white text-sm font-medium">
                    {user.name.charAt(0)}
                  </div>
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
                    className="text-sm text-gray-600 hover:underline"
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
