import { useState, useRef, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useSearch } from "../../features/search/hooks/useSearch";
import { Input } from "../ui/input";

const entityRouteMap: Record<string, string> = {
  PROJECT: "/projects",
  TASK: "/tasks",
  PROJECT_REQUEST: "/project-requests",
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

  const handleResultClick = (entityType: string, entityId: number) => {
    const base = entityRouteMap[entityType] || "/";
    navigate(`${base}/${entityId}`);
    setQuery("");
    setShowDropdown(false);
  };

  const hasResults =
    data && (data.projects.length > 0 || data.tasks.length > 0 || data.requests.length > 0);

  return (
    <div ref={containerRef} className="relative w-80">
      <Input
        placeholder="Search projects, tasks, requests..."
        value={query}
        onChange={(e) => {
          setQuery(e.target.value);
          setShowDropdown(true);
        }}
        onFocus={() => setShowDropdown(true)}
      />

      {showDropdown && query.length >= 2 && (
        <div className="absolute z-50 mt-1 w-full rounded border bg-white shadow-lg">
          {isLoading ? (
            <p className="p-3 text-sm text-slate-400">Searching...</p>
          ) : !hasResults ? (
            <p className="p-3 text-sm text-slate-400">No results found</p>
          ) : (
            <div className="max-h-96 overflow-y-auto">
              {data.projects.length > 0 && (
                <div>
                  <p className="bg-slate-50 px-3 py-1 text-xs font-medium text-slate-500">Projects</p>
                  {data.projects.map((r) => (
                    <button
                      key={`project-${r.entity_id}`}
                      onClick={() => handleResultClick("PROJECT", r.entity_id)}
                      className="block w-full px-3 py-2 text-left text-sm hover:bg-slate-50"
                    >
                      {r.title}
                    </button>
                  ))}
                </div>
              )}

              {data.tasks.length > 0 && (
                <div>
                  <p className="bg-slate-50 px-3 py-1 text-xs font-medium text-slate-500">Tasks</p>
                  {data.tasks.map((r) => (
                    <button
                      key={`task-${r.entity_id}`}
                      onClick={() => handleResultClick("TASK", r.entity_id)}
                      className="block w-full px-3 py-2 text-left text-sm hover:bg-slate-50"
                    >
                      {r.title}
                    </button>
                  ))}
                </div>
              )}

              {data.requests.length > 0 && (
                <div>
                  <p className="bg-slate-50 px-3 py-1 text-xs font-medium text-slate-500">Project Requests</p>
                  {data.requests.map((r) => (
                    <button
                      key={`request-${r.entity_id}`}
                      onClick={() => handleResultClick("PROJECT_REQUEST", r.entity_id)}
                      className="block w-full px-3 py-2 text-left text-sm hover:bg-slate-50"
                    >
                      {r.title}
                    </button>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
}