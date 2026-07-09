import { useState, useEffect } from "react";
import { fetchTags, type Tag } from "../api/tags";

interface Props {
  selectedIds: number[];
  onChange: (ids: number[]) => void;
}

export default function InterestTagsPicker({ selectedIds, onChange }: Props) {
  const [tags, setTags] = useState<Tag[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchTags()
      .then(setTags)
      .finally(() => setLoading(false));
  }, []);

  function toggleTag(id: number) {
    if (selectedIds.includes(id)) {
      onChange(selectedIds.filter((i) => i !== id));
    } else {
      onChange([...selectedIds, id]);
    }
  }

  if (loading) {
    return <p className="text-gray-500 dark:text-gray-400 text-sm">加载标签中...</p>;
  }

  return (
    <div>
      <p className="text-sm text-gray-500 dark:text-gray-400 mb-3">
         选择你感兴趣的服务标签，方便我们为你推荐（可跳过）
       </p>
      <div className="flex flex-wrap gap-2">
        {tags.map((tag) => {
          const selected = selectedIds.includes(tag.id);
          return (
            <button
              key={tag.id}
              type="button"
              onClick={() => toggleTag(tag.id)}
              className={`px-3 py-1.5 rounded-full text-sm border transition-colors ${
                selected
                  ? "bg-blue-600 text-white border-blue-600"
                  : "bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 border-gray-300 dark:border-gray-600 hover:border-blue-400"
              }`}
            >
              {tag.name}
            </button>
          );
        })}
      </div>
    </div>
  );
}
