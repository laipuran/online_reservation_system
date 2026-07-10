export interface ProviderFields {
  businessName: string;
  description: string;
  address: string;
  email: string;
  phone: string;
  logoUrl: string;
}

export const EMPTY_PROVIDER_FIELDS: ProviderFields = {
  businessName: "",
  description: "",
  address: "",
  email: "",
  phone: "",
  logoUrl: "",
};

export function validateProviderFields(
  fields: ProviderFields
): Partial<Record<keyof ProviderFields, string>> {
  const errors: Partial<Record<keyof ProviderFields, string>> = {};
  if (!fields.businessName.trim()) errors.businessName = "请输入商家名称";
  if (!fields.description.trim()) errors.description = "请输入商家简介";
  if (!fields.address.trim()) errors.address = "请输入地址";
  if (!fields.email.trim()) {
    errors.email = "请输入联系邮箱";
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(fields.email.trim())) {
    errors.email = "邮箱格式不正确";
  }
  if (fields.phone.trim() && !/^1[3-9]\d{9}$/.test(fields.phone.trim())) {
    errors.phone = "手机号格式不正确";
  }
  if (fields.logoUrl.trim()) {
    try {
      new URL(fields.logoUrl.trim());
    } catch {
      errors.logoUrl = "Logo URL 格式不正确";
    }
  }
  return errors;
}
