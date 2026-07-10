import type { Route } from "./+types/my-reservations";

export function meta({}: Route.MetaArgs) {
  return [{ title: "我的预约 - ORS" }];
}

export default function MyReservations() {
  return (
    <div className="max-w-3xl mx-auto mt-20 px-4 text-center">
      <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-4">我的预约</h1>
      <p className="text-gray-400 dark:text-gray-500">功能开发中，敬请期待</p>
    </div>
  );
}
