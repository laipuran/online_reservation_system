import { useState, type FormEvent } from "react";
import { useNavigate, Link } from "react-router";
import { useAuth } from "../../lib/hooks/use-auth";
import { updateMyProfile, updateMyPassword } from "../../lib/api/users";
import { ApiError } from "../../lib/api/client";

export default function SettingsPage() {
  const { user, loading, setUser } = useAuth();
  const navigate = useNavigate();

  const [name, setName] = useState(user?.name ?? "");
  const [phone, setPhone] = useState(user?.phone ?? "");
  const [avatarUrl, setAvatarUrl] = useState(user?.avatar_url ?? "");
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [profileError, setProfileError] = useState("");
  const [passwordError, setPasswordError] = useState("");
  const [profileSaving, setProfileSaving] = useState(false);
  const [passwordSaving, setPasswordSaving] = useState(false);
  const [profileOk, setProfileOk] = useState(false);
  const [passwordOk, setPasswordOk] = useState(false);

  if (loading || !user) {
    return (
      <div className="flex items-center justify-center mt-20">
        <p className="text-gray-500 dark:text-gray-400">加载中...</p>
      </div>
    );
  }

  async function handleProfileSubmit(e: FormEvent) {
    e.preventDefault();
    setProfileError("");
    setProfileOk(false);

    if (!name.trim()) {
      setProfileError("昵称不能为空");
      return;
    }

    setProfileSaving(true);
    try {
      const updated = await updateMyProfile({
        name: name.trim() || undefined,
        phone: phone.trim() || undefined,
        avatar_url: avatarUrl.trim() || undefined,
      });
      setProfileOk(true);
      setUser(updated);
    } catch (err) {
      if (err instanceof ApiError) {
        setProfileError(err.message);
      } else {
        setProfileError("保存失败，请稍后重试");
      }
    } finally {
      setProfileSaving(false);
    }
  }

  async function handlePasswordSubmit(e: FormEvent) {
    e.preventDefault();
    setPasswordError("");
    setPasswordOk(false);

    if (!currentPassword) {
      setPasswordError("请输入当前密码");
      return;
    }
    if (newPassword.length < 8) {
      setPasswordError("新密码长度至少 8 位");
      return;
    }

    setPasswordSaving(true);
    try {
      await updateMyPassword({
        current_password: currentPassword,
        new_password: newPassword,
      });
      setPasswordOk(true);
      setCurrentPassword("");
      setNewPassword("");
    } catch (err) {
      if (err instanceof ApiError) {
        setPasswordError(err.message);
      } else {
        setPasswordError("修改失败，请稍后重试");
      }
    } finally {
      setPasswordSaving(false);
    }
  }

  return (
    <div className="max-w-lg mx-auto mt-8 px-4 pb-12">
      <Link
        to="/dashboard"
        className="text-sm text-blue-600 dark:text-blue-400 hover:underline"
      >
        &larr; 返回个人主页
      </Link>

      <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mt-4 mb-6">个人设置</h1>

      <form onSubmit={handleProfileSubmit} className="space-y-4 mb-10">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 border-b border-gray-200 dark:border-gray-700 pb-2">基本资料</h2>

        <div>
          <label className="block text-sm font-medium mb-1 dark:text-gray-300">邮箱</label>
          <input
            type="email"
            value={user.email}
            disabled
            className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 text-sm cursor-not-allowed"
          />
          <p className="text-xs text-gray-400 dark:text-gray-500 mt-1">邮箱不可修改</p>
        </div>

        <div>
          <label className="block text-sm font-medium mb-1 dark:text-gray-300">昵称</label>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 dark:text-gray-100 text-sm"
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1 dark:text-gray-300">手机号</label>
          <input
            type="tel"
            value={phone}
            onChange={(e) => setPhone(e.target.value)}
            className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 dark:text-gray-100 text-sm"
            placeholder="13800000000"
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1 dark:text-gray-300">头像 URL</label>
          <input
            type="url"
            value={avatarUrl}
            onChange={(e) => setAvatarUrl(e.target.value)}
            className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 dark:text-gray-100 text-sm"
            placeholder="https://example.com/avatar.png"
          />
        </div>

        {profileError && <p className="text-red-500 dark:text-red-400 text-sm">{profileError}</p>}
        {profileOk && <p className="text-green-500 dark:text-green-400 text-sm">保存成功</p>}

        <button
          type="submit"
          disabled={profileSaving}
          className="bg-blue-600 text-white px-6 py-2 rounded text-sm hover:bg-blue-700 disabled:opacity-50"
        >
          {profileSaving ? "保存中..." : "保存资料"}
        </button>
      </form>

      <form onSubmit={handlePasswordSubmit} className="space-y-4">
        <h2 className="text-lg font-semibold text-gray-900 dark:text-gray-100 border-b border-gray-200 dark:border-gray-700 pb-2">修改密码</h2>

        <div>
          <label className="block text-sm font-medium mb-1 dark:text-gray-300">当前密码</label>
          <input
            type="password"
            value={currentPassword}
            onChange={(e) => setCurrentPassword(e.target.value)}
            className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 dark:text-gray-100 text-sm"
          />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1 dark:text-gray-300">新密码</label>
          <input
            type="password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 dark:text-gray-100 text-sm"
          />
        </div>

        {passwordError && <p className="text-red-500 dark:text-red-400 text-sm">{passwordError}</p>}
        {passwordOk && <p className="text-green-500 dark:text-green-400 text-sm">密码修改成功</p>}

        <button
          type="submit"
          disabled={passwordSaving}
          className="bg-blue-600 text-white px-6 py-2 rounded text-sm hover:bg-blue-700 disabled:opacity-50"
        >
          {passwordSaving ? "修改中..." : "修改密码"}
        </button>
      </form>
    </div>
  );
}
