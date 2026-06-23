import { useState, useRef, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useSearch } from "../../features/search/hooks/useSearch";

const entityRouteMap: Record<string, string> = {
  PROJECT: "/projects",
  TASK: "/tasks",
  PROJECT_REQUEST: "/project-requests",
};

const entityLabel: Record<string, string> = {
  PROJECT: "Projects",
  TASK: "Tasks",
  PROJECT_REQUEST: "Project Requests",
};

export default function GlobalSearchBar() {
  const [query, setQuery] = useState("");
  const [showDropdown, setShowDropdown] = useState(false);
  const navigate = useNavigate();
  const containerRef = useRef<HTMLDivElement>(null);

  const { data, isLoading } = useSearch(query);

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setShowDropdown(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  useEffect(() => {
    function handleKeydown(e: KeyboardEvent) {
      if ((e.metaKey || e.ctrlKey) && e.key === "k") {
        e.preventDefault();
        containerRef.current?.querySelector("input")?.focus();
      }
      if (e.key === "Escape") {
        setShowDropdown(false);
      }
    }
    document.addEventListener("keydown", handleKeydown);
    return () => document.removeEventListener("keydown", handleKeydown);
  }, []);

  const handleResultClick = (entityType: string, entityId: number) => {
    const base = entityRouteMap[entityType] || "/";
    navigate(`${base}/${entityId}`);
    setQuery("");
    setShowDropdown(false);
  };

  const groups = data
    ? [
        { key: "projects", type: "PROJECT", items: data.projects },
        { key: "tasks", type: "TASK", items: data.tasks },
        { key: "requests", type: "PROJECT_REQUEST", items: data.requests },
      ].filter((g) => g.items.length > 0)
    : [];

  return (
    <div ref={containerRef} className="relative w-[280px]">
      <div className="flex items-center gap-2 h-9 px-3 rounded-md border border-border bg-surface-secondary text-[13px] text-ink-tertiary focus-within:border-primary-500 focus-within:bg-surface focus-within:ring-2 focus-within:ring-primary-100 dark:focus-within:ring-primary-900/30 transition-all">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="flex-shrink-0">
          <circle cx="11" cy="11" r="7" />
          <path d="M21 21l-4.3-4.3" />
        </svg>
        <input
          value={query}
          onChange={(e) => {
            setQuery(e.target.value);
            setShowDropdown(true);
          }}
          onFocus={() => setShowDropdown(true)}
          placeholder="Cari project, task, request..."
          className="flex-1 bg-transparent outline-none text-ink-primary placeholder:text-ink-tertiary"
        />
        <span className="text-[10.5px] font-semibold border border-border-strong rounded px-1.5 py-0.5 text-ink-tertiary flex-shrink-0">
          ⌘K
        </span>
      </div>

      {showDropdown && query.length >= 2 && (
        <div className="absolute z-50 mt-1.5 w-[340px] rounded-lg border border-border bg-surface shadow-lg overflow-hidden">
          {isLoading ? (
            <p className="p-4 text-[12.5px] text-ink-tertiary">Mencari...</p>
          ) : groups.length === 0 ? (
            <p className="p-4 text-[12.5px] text-ink-tertiary">Tidak ada hasil untuk "{query}"</p>
          ) : (
            <div className="max-h-96 overflow-y-auto py-1.5">
              {groups.map((group) => (
                <div key={group.key} className="mb-1 last:mb-0">
                  <p className="px-3 py-1 text-[10.5px] font-semibold uppercase tracking-wide text-ink-tertiary">
                    {entityLabel[group.type]}
                  </p>
                  {group.items.map((r) => (
                    <button
                      key={`${group.type}-${r.entity_id}`}
                      onClick={() => handleResultClick(group.type, r.entity_id)}
                      className="block w-full px-3 py-2 text-left text-[13px] text-ink-primary hover:bg-surface-secondary transition-colors truncate"
                    >
                      {r.title}
                    </button>
                  ))}
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
}