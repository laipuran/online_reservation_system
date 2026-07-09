import { useState, useRef, useEffect, useMemo } from "react";
import { fetchTags, createTag, type Tag } from "../api/tags";
import { X } from "lucide-react";

interface TagInputProps {
  selectedTags: Tag[];
  onChange: (tags: Tag[]) => void;
}

export function TagInput({ selectedTags, onChange }: TagInputProps) {
  const [input, setInput] = useState("");
  const [allTags, setAllTags] = useState<Tag[]>([]);
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    fetchTags().then(setAllTags).catch(() => {});
  }, []);

  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  const suggestions = useMemo(() => {
    if (!input.trim()) return [];
    const lower = input.toLowerCase();
    const selectedIds = new Set(selectedTags.map((t) => t.id));
    return allTags.filter(
      (t) => t.name.includes(lower) && !selectedIds.has(t.id)
    );
  }, [input, allTags, selectedTags]);

  function addTag(tag: Tag) {
    if (selectedTags.some((t) => t.id === tag.id)) return;
    onChange([...selectedTags, tag]);
    setInput("");
    setOpen(false);
  }

  function removeTag(tagId: number) {
    onChange(selectedTags.filter((t) => t.id !== tagId));
  }

  function handleKeyDown(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key === "Enter" && input.trim()) {
      e.preventDefault();
      const existing = allTags.find(
        (t) => t.name.toLowerCase() === input.trim().toLowerCase()
      );
      if (existing) {
        addTag(existing);
      } else {
        onChange([
          ...selectedTags,
          { id: 0, name: input.trim(), created_at: "" },
        ]);
      }
      setInput("");
    }
    if (e.key === "Backspace" && !input && selectedTags.length > 0) {
      removeTag(selectedTags[selectedTags.length - 1].id);
    }
  }

  return (
    <div ref={ref} className="relative">
      <div className="flex flex-wrap gap-1.5 mb-1.5">
        {selectedTags.map((tag) => (
          <span
            key={tag.id || tag.name}
            className="inline-flex items-center gap-1 text-xs px-2 py-1 rounded-full bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300"
          >
            {tag.name}
            <button
              type="button"
              onClick={() => removeTag(tag.id)}
              className="hover:bg-blue-200 dark:hover:bg-blue-800 rounded-full p-0.5"
            >
              <X className="w-3 h-3" />
            </button>
          </span>
        ))}
      </div>
      <input
        type="text"
        value={input}
        onChange={(e) => {
          setInput(e.target.value);
          setOpen(true);
        }}
        onFocus={() => setOpen(true)}
        onKeyDown={handleKeyDown}
        placeholder={selectedTags.length === 0 ? "输入标签名，回车添加..." : "继续添加标签..."}
        className="w-full border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 dark:text-gray-100 text-sm"
      />
      {open && suggestions.length > 0 && (
        <div className="absolute left-0 right-0 top-full mt-1 z-10 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-lg max-h-40 overflow-y-auto">
          {suggestions.map((tag) => (
            <button
              key={tag.id}
              type="button"
              onClick={() => addTag(tag)}
              className="w-full text-left px-3 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-700 dark:text-gray-200"
            >
              {tag.name}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
