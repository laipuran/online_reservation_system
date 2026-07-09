import { useState } from "react";
import { useNavigate, Link } from "react-router";
import RegisterStepBasic from "../../lib/components/register-step-basic";
import ProviderInfoForm from "../../lib/components/provider-info-form";
import InterestTagsPicker from "../../lib/components/interest-tags-picker";
import {
  useRegister,
  useCreateProviderProfile,
  useSetUserInterests,
} from "../../lib/hooks/use-mutations";
import { ApiError } from "../../lib/api/client";
import {
  validateProviderFields,
  EMPTY_PROVIDER_FIELDS,
  type ProviderFields,
} from "../../lib/validation";

type Role = "customer" | "provider";

export default function Register() {
  const navigate = useNavigate();

  const [step, setStep] = useState(1);
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState<Role>("customer");
  const [error, setError] = useState("");

  const [providerFields, setProviderFields] = useState<ProviderFields>(EMPTY_PROVIDER_FIELDS);
  const [providerErrors, setProviderErrors] = useState<Partial<Record<keyof ProviderFields, string>>>({});
  const [interestIds, setInterestIds] = useState<number[]>([]);

  const registerMutation = useRegister();
  const createProfileMutation = useCreateProviderProfile();
  const setInterestsMutation = useSetUserInterests();

  const submitting =
    registerMutation.isPending ||
    createProfileMutation.isPending ||
    setInterestsMutation.isPending;

  function fieldErrors() {
    const errs: { name?: string; email?: string; password?: string } = {};
    if (!name.trim()) errs.name = "请输入昵称";
    if (!email.trim()) {
      errs.email = "请输入邮箱";
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.trim())) {
      errs.email = "邮箱格式不正确";
    }
    if (password.length < 8) errs.password = "密码长度至少 8 位";
    return errs;
  }

  function handleNext() {
    setError("");

    const errs = fieldErrors();
    if (Object.keys(errs).length > 0) {
      const first = errs.name || errs.email || errs.password || "";
      setError(first);
      return;
    }

    setStep(2);
  }

  function handleBack() {
    setStep(1);
  }

  function handleFinalSubmit() {
    if (role === "provider") {
      const errs = validateProviderFields(providerFields);
      setProviderErrors(errs);
      if (Object.keys(errs).length > 0) return;
    }

    handleRegister();
  }

  async function handleRegister() {
    setError("");

    try {
      await registerMutation.mutateAsync({ name: name.trim(), email: email.trim(), password, role });

      if (role === "provider") {
        await createProfileMutation.mutateAsync({
          business_name: providerFields.businessName.trim(),
          description: providerFields.description.trim(),
          address: providerFields.address.trim(),
          email: providerFields.email.trim(),
          phone: providerFields.phone.trim() || undefined,
          logo_url: providerFields.logoUrl.trim() || undefined,
        });
      }

      if (role === "customer" && interestIds.length > 0) {
        await setInterestsMutation.mutateAsync(interestIds);
      }

      if (role === "provider") {
        navigate("/dashboard", { replace: true });
      } else {
        navigate("/dashboard", { replace: true });
      }
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError("注册失败，请稍后重试");
      }
    }
  }

  const step1FieldErrors = fieldErrors();
  // Only show inline errors after user has attempted to proceed (error is set)
  const inlineErrors = error
    ? step1FieldErrors
    : {};

  return (
    <div className="max-w-sm mx-auto mt-20 px-4">
      <h1 className="text-2xl font-bold text-center mb-6">注册</h1>

      {/* step indicator */}
      <div className="flex items-center justify-center gap-2 mb-6">
        <span
          className={`w-2.5 h-2.5 rounded-full ${
            step === 1 ? "bg-blue-600" : "bg-gray-300 dark:bg-gray-600"
          }`}
        />
        <span className="text-gray-300 dark:text-gray-600">—</span>
        <span
          className={`w-2.5 h-2.5 rounded-full ${
            step === 2 ? "bg-blue-600" : "bg-gray-300 dark:bg-gray-600"
          }`}
        />
      </div>

      {step === 1 && (
        <RegisterStepBasic
          name={name}
          email={email}
          password={password}
          role={role}
          error={error}
          fieldErrors={inlineErrors}
          loading={submitting}
          onNameChange={setName}
          onEmailChange={setEmail}
          onPasswordChange={setPassword}
          onRoleChange={setRole}
          onNext={handleNext}
        />
      )}

      {step === 2 && (
        <div className="space-y-6">
          {role === "customer" && (
            <InterestTagsPicker
              selectedIds={interestIds}
              onChange={setInterestIds}
            />
          )}

          {role === "provider" && (
            <ProviderInfoForm
              businessName={providerFields.businessName}
              description={providerFields.description}
              address={providerFields.address}
              email={providerFields.email}
              phone={providerFields.phone}
              logoUrl={providerFields.logoUrl}
              errors={providerErrors}
              onBusinessNameChange={(v) =>
                setProviderFields((f) => ({ ...f, businessName: v }))
              }
              onDescriptionChange={(v) =>
                setProviderFields((f) => ({ ...f, description: v }))
              }
              onAddressChange={(v) =>
                setProviderFields((f) => ({ ...f, address: v }))
              }
              onEmailChange={(v) =>
                setProviderFields((f) => ({ ...f, email: v }))
              }
              onPhoneChange={(v) =>
                setProviderFields((f) => ({ ...f, phone: v }))
              }
              onLogoUrlChange={(v) =>
                setProviderFields((f) => ({ ...f, logoUrl: v }))
              }
            />
          )}

          {error && <p className="text-red-500 text-sm">{error}</p>}

          <div className="flex gap-3">
            <button
              type="button"
              onClick={handleBack}
              disabled={submitting}
              className="flex-1 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 py-2 rounded hover:bg-gray-50 dark:hover:bg-gray-800 disabled:opacity-50"
            >
              上一步
            </button>
            <button
              type="button"
              onClick={handleFinalSubmit}
              disabled={submitting}
              className="flex-1 bg-blue-600 text-white py-2 rounded hover:bg-blue-700 disabled:opacity-50"
            >
              {submitting ? "注册中..." : "提交注册"}
            </button>
          </div>
        </div>
      )}

      <p className="text-sm text-center mt-4 text-gray-500 dark:text-gray-400">
        已有账号？{" "}
        <Link to="/login" className="text-blue-600 dark:text-blue-400 hover:underline">
          登录
        </Link>
      </p>
    </div>
  );
}
