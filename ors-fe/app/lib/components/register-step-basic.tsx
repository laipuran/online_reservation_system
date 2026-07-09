import type { KeyboardEvent } from "react";

type Role = "customer" | "provider";

interface Props {
  name: string;
  email: string;
  password: string;
  role: Role;
  error: string;
  fieldErrors: { name?: string; email?: string; password?: string };
  loading: boolean;
  onNameChange: (v: string) => void;
  onEmailChange: (v: string) => void;
  onPasswordChange: (v: string) => void;
  onRoleChange: (v: Role) => void;
  onNext: () => void;
}

const ROLE_OPTIONS: { value: Role; label: string }[] = [
  { value: "customer", label: "服务体验者" },
  { value: "provider", label: "服务提供者" },
];

export default function RegisterStepBasic({
  name,
  email,
  password,
  role,
  error,
  fieldErrors,
  loading,
  onNameChange,
  onEmailChange,
  onPasswordChange,
  onRoleChange,
  onNext,
}: Props) {
  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === "Enter" && !loading) {
      onNext();
    }
  }

  return (
    <div className="space-y-4" onKeyDown={handleKeyDown}>
      <div>
        <label className="block text-sm font-medium mb-1">昵称</label>
        <input
          type="text"
          value={name}
          onChange={(e) => onNameChange(e.target.value)}
          className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800"
          placeholder="张三"
        />
        {fieldErrors.name && (
          <p className="text-red-500 text-sm mt-1">{fieldErrors.name}</p>
        )}
      </div>
      <div>
        <label className="block text-sm font-medium mb-1">邮箱</label>
        <input
          type="email"
          value={email}
          onChange={(e) => onEmailChange(e.target.value)}
          className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800"
          placeholder="your@email.com"
        />
        {fieldErrors.email && (
          <p className="text-red-500 text-sm mt-1">{fieldErrors.email}</p>
        )}
      </div>
      <div>
        <label className="block text-sm font-medium mb-1">密码</label>
        <input
          type="password"
          value={password}
          onChange={(e) => onPasswordChange(e.target.value)}
          className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800"
        />
        {fieldErrors.password && (
          <p className="text-red-500 text-sm mt-1">{fieldErrors.password}</p>
        )}
      </div>
      <div>
        <label className="block text-sm font-medium mb-3">身份</label>
        <div className="space-y-2">
          {ROLE_OPTIONS.map((opt) => (
            <label
              key={opt.value}
              className={`flex items-center gap-3 border rounded-lg px-4 py-3 cursor-pointer transition-colors ${
                role === opt.value
                  ? "border-blue-500 dark:border-blue-400 bg-blue-50 dark:bg-blue-900/20"
                  : "border-gray-300 dark:border-gray-600 hover:border-gray-400"
              }`}
            >
              <input
                type="radio"
                name="role"
                value={opt.value}
                checked={role === opt.value}
                onChange={() => onRoleChange(opt.value)}
                className="accent-blue-600"
              />
              <span className="font-medium">{opt.label}</span>
            </label>
          ))}
        </div>
      </div>
      <button
        type="button"
        onClick={onNext}
        disabled={loading}
        className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 disabled:opacity-50"
      >
        {loading ? "加载中..." : "下一步"}
      </button>
    </div>
  );
}
