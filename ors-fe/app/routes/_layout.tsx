import { Link, Outlet, useLocation, useNavigate } from "react-router";
import { useAuth } from "../lib/hooks/use-auth";
import { NotificationBell } from "../lib/components/notification-bell";
import { ThemeToggle } from "../lib/components/theme-toggle";

const GUEST_TABS = [
  { label: "首页", to: "/" },
  { label: "预约项目", to: "/services" },
];

const CUSTOMER_TABS = [
  { label: "首页", to: "/" },
  { label: "预约项目", to: "/services" },
  { label: "我的预约", to: "/dashboard" },
];

const PROVIDER_TABS = [
  { label: "服务管理", to: "/provider/services" },
  { label: "预约处理", to: "/provider/reservations" },
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

  const tabs = !user ? GUEST_TABS : user.role === "provider" ? PROVIDER_TABS : CUSTOMER_TABS;

  return (
    <div className="min-h-screen flex flex-col">
      <header className="sticky top-0 z-50 bg-white dark:bg-gray-950 border-b border-gray-200 dark:border-gray-700">
        <nav className="max-w-5xl mx-auto flex items-center justify-between px-4 h-14">
          <div className="flex items-center gap-6">
            <Link to="/" className="font-bold text-lg shrink-0">
              ORS
            </Link>
            <div className="flex items-center gap-1">
              {tabs.map((tab) => (
                <Link
                  key={tab.to}
                  to={tab.to}
                  className={`px-3 py-1.5 text-sm rounded-md transition-colors ${
                    isActive(tab.to)
                      ? "bg-blue-50 text-blue-700 font-medium dark:bg-blue-900/30 dark:text-blue-300"
                      : "text-gray-600 hover:text-gray-900 hover:bg-gray-50 dark:text-gray-400 dark:hover:text-gray-200 dark:hover:bg-gray-800"
                  }`}
                >
                  {tab.label}
                </Link>
              ))}
            </div>
          </div>
          {!loading && (
            <div className="flex items-center gap-3 shrink-0">
              <ThemeToggle />
              {user ? (
                <>
                  <Link
                    to="/dashboard"
                    className="text-sm text-gray-600 hover:underline dark:text-gray-400"
                  >
                    {user.name}
                  </Link>
                  <NotificationBell />
                  {user.avatar_url ? (
                    <img src={user.avatar_url} alt={user.name} className="w-7 h-7 rounded-full object-cover" />
                  ) : (
                    <div className="w-7 h-7 rounded-full bg-blue-500 flex items-center justify-center text-white text-sm font-medium">
                      {user.name.charAt(0)}
                    </div>
                  )}
                  <button
                    onClick={handleLogout}
                    className="text-sm text-red-500 dark:text-red-400 hover:underline"
                  >
                    退出
                  </button>
                </>
              ) : (
                <>
                  <Link
                    to="/login"
                    className="text-sm text-gray-600 hover:underline dark:text-gray-400"
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
