import { useState, useEffect } from "react";
import { useNavigate } from "react-router";
import ProviderInfoForm from "../../lib/components/provider-info-form";
import { useCreateProviderProfile } from "../../lib/hooks/use-mutations";
import { useAuth } from "../../lib/hooks/use-auth";
import { ApiError } from "../../lib/api/client";

interface Fields {
  businessName: string;
  description: string;
  address: string;
  email: string;
  phone: string;
  logoUrl: string;
}

const EMPTY_FIELDS: Fields = {
  businessName: "",
  description: "",
  address: "",
  email: "",
  phone: "",
  logoUrl: "",
};

const PHONE_REGEX = /^1[3-9]\d{9}$/;

function validateProviderFields(
  fields: Fields
): Partial<Record<keyof Fields, string>> {
  const errors: Partial<Record<keyof Fields, string>> = {};
  if (!fields.businessName.trim()) errors.businessName = "请输入商家名称";
  if (!fields.description.trim()) errors.description = "请输入商家简介";
  if (!fields.address.trim()) errors.address = "请输入地址";
  if (!fields.email.trim()) errors.email = "请输入联系邮箱";
  if (fields.phone.trim() && !PHONE_REGEX.test(fields.phone.trim())) {
    errors.phone = "手机号格式不正确";
  }
  return errors;
}

export default function CompleteProfile() {
  const navigate = useNavigate();
  const { user, loading } = useAuth();
  const [fields, setFields] = useState<Fields>(EMPTY_FIELDS);
  const [errors, setErrors] = useState<
    Partial<Record<keyof Fields, string>>
  >({});
  const [apiError, setApiError] = useState("");
  const mutation = useCreateProviderProfile();

  useEffect(() => {
    if (!loading && !user) {
      navigate("/login", { replace: true });
    }
    if (!loading && user && user.role !== "provider") {
      navigate("/dashboard", { replace: true });
    }
  }, [user, loading, navigate]);

  if (loading || !user) {
    return (
      <div className="flex items-center justify-center mt-20">
        <p className="text-gray-500 dark:text-gray-400">加载中...</p>
      </div>
    );
  }

  async function handleSubmit() {
    const errs = validateProviderFields(fields);
    setErrors(errs);
    if (Object.keys(errs).length > 0) return;

    setApiError("");

    try {
      await mutation.mutateAsync({
        business_name: fields.businessName.trim(),
        description: fields.description.trim(),
        address: fields.address.trim(),
        email: fields.email.trim(),
        phone: fields.phone.trim() || undefined,
        logo_url: fields.logoUrl.trim() || undefined,
      });
      navigate("/dashboard");
    } catch (err) {
      if (err instanceof ApiError) {
        setApiError(err.message);
      } else {
        setApiError("保存失败，请稍后重试");
      }
    }
  }

  return (
    <div className="max-w-sm mx-auto mt-20 px-4">
      <h1 className="text-2xl font-bold text-center mb-2 text-gray-900 dark:text-gray-100">完善商家信息</h1>
      <p className="text-sm text-gray-500 dark:text-gray-400 text-center mb-6">
        请填写商家资料以完成注册
      </p>

      <ProviderInfoForm
        businessName={fields.businessName}
        description={fields.description}
        address={fields.address}
        email={fields.email}
        phone={fields.phone}
        logoUrl={fields.logoUrl}
        errors={errors}
        onBusinessNameChange={(v) => setFields((f) => ({ ...f, businessName: v }))}
        onDescriptionChange={(v) => setFields((f) => ({ ...f, description: v }))}
        onAddressChange={(v) => setFields((f) => ({ ...f, address: v }))}
        onEmailChange={(v) => setFields((f) => ({ ...f, email: v }))}
        onPhoneChange={(v) => setFields((f) => ({ ...f, phone: v }))}
        onLogoUrlChange={(v) => setFields((f) => ({ ...f, logoUrl: v }))}
      />

      {apiError && <p className="text-red-500 text-sm mt-4">{apiError}</p>}

      <button
        type="button"
        onClick={handleSubmit}
        disabled={mutation.isPending}
        className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 disabled:opacity-50 mt-6"
      >
        {mutation.isPending ? "保存中..." : "保存"}
      </button>
    </div>
  );
}
