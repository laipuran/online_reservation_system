interface Props {
  businessName: string;
  description: string;
  address: string;
  email: string;
  phone: string;
  logoUrl: string;
  errors: Partial<Record<keyof Fields, string>>;
  onBusinessNameChange: (v: string) => void;
  onDescriptionChange: (v: string) => void;
  onAddressChange: (v: string) => void;
  onEmailChange: (v: string) => void;
  onPhoneChange: (v: string) => void;
  onLogoUrlChange: (v: string) => void;
}

interface Fields {
  businessName: string;
  description: string;
  address: string;
  email: string;
  phone: string;
  logoUrl: string;
}

export default function ProviderInfoForm({
  businessName,
  description,
  address,
  email,
  phone,
  logoUrl,
  errors,
  onBusinessNameChange,
  onDescriptionChange,
  onAddressChange,
  onEmailChange,
  onPhoneChange,
  onLogoUrlChange,
}: Props) {
  return (
    <div className="space-y-4">
      <div>
        <label className="block text-sm font-medium mb-1">
          商家名称 <span className="text-red-500">*</span>
        </label>
        <input
          type="text"
          value={businessName}
          onChange={(e) => onBusinessNameChange(e.target.value)}
          className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800"
          placeholder="舒心养生馆"
        />
        {errors.businessName && (
          <p className="text-red-500 text-sm mt-1">{errors.businessName}</p>
        )}
      </div>
      <div>
        <label className="block text-sm font-medium mb-1">
          商家简介 <span className="text-red-500">*</span>
        </label>
        <textarea
          value={description}
          onChange={(e) => onDescriptionChange(e.target.value)}
          rows={3}
          className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 resize-none"
          placeholder="专业按摩服务，十余年经验..."
        />
        {errors.description && (
          <p className="text-red-500 text-sm mt-1">{errors.description}</p>
        )}
      </div>
      <div>
        <label className="block text-sm font-medium mb-1">
          地址 <span className="text-red-500">*</span>
        </label>
        <input
          type="text"
          value={address}
          onChange={(e) => onAddressChange(e.target.value)}
          className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800"
          placeholder="上海市徐汇区..."
        />
        {errors.address && (
          <p className="text-red-500 text-sm mt-1">{errors.address}</p>
        )}
      </div>
      <div>
        <label className="block text-sm font-medium mb-1">
          联系邮箱 <span className="text-red-500">*</span>
        </label>
        <input
          type="email"
          value={email}
          onChange={(e) => onEmailChange(e.target.value)}
          className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800"
          placeholder="shop@example.com"
        />
        {errors.email && (
          <p className="text-red-500 text-sm mt-1">{errors.email}</p>
        )}
      </div>
      <div>
        <label className="block text-sm font-medium mb-1">联系电话</label>
        <input
          type="tel"
          value={phone}
          onChange={(e) => onPhoneChange(e.target.value)}
          className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800"
          placeholder="13800000000"
        />
        {errors.phone && (
          <p className="text-red-500 text-sm mt-1">{errors.phone}</p>
        )}
      </div>
      <div>
        <label className="block text-sm font-medium mb-1">Logo URL</label>
        <input
          type="url"
          value={logoUrl}
          onChange={(e) => onLogoUrlChange(e.target.value)}
          className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800"
          placeholder="https://example.com/logo.png"
        />
        {errors.logoUrl && (
          <p className="text-red-500 text-sm mt-1">{errors.logoUrl}</p>
        )}
      </div>
    </div>
  );
}
