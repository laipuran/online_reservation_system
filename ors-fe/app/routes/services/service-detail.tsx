import { useParams } from "react-router";

export default function ServiceDetail() {
  const { id } = useParams();

  return (
    <div className="max-w-3xl mx-auto mt-20 px-4">
      <h1 className="text-2xl font-bold mb-4 text-gray-900 dark:text-gray-100">服务详情</h1>
      <p className="text-gray-500 dark:text-gray-400">
        服务 #{id} 的详情页面（待实现）
      </p>
    </div>
  );
}
