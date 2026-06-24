import { useTheme } from "../../lib/theme";
import GlobalSearchBar from "./GlobalSearchBar";
import NotificationBell from "./NotificationBell";

export default function Navbar() {
  const { theme, toggleTheme } = useTheme();

  return (
    <header className="h-14 flex-shrink-0 bg-surface border-b border-border flex items-center justify-between px-6 sticky top-0 z-10">
      <GlobalSearchBar />

      <div className="flex items-center gap-2">
        <button
          onClick={toggleTheme}
          className="flex items-center gap-0.5 bg-surface-secondary border border-border rounded-full p-[3px] cursor-pointer"
          aria-label="Toggle theme"
        >
          <span
            className={`h-6 w-[26px] rounded-full flex items-center justify-center transition-colors ${
              theme === "light" ? "bg-surface text-primary-600 shadow-sm" : "text-ink-tertiary"
            }`}
          >
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2">
              <circle cx="12" cy="12" r="4" />
              <path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4" />
            </svg>
          </span>
          <span
            className={`h-6 w-[26px] rounded-full flex items-center justify-center transition-colors ${
              theme === "dark" ? "bg-surface text-primary-600 shadow-sm" : "text-ink-tertiary"
            }`}
          >
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2">
              <path d="M21 12.8A9 9 0 1111.2 3 7 7 0 0021 12.8z" />
            </svg>
          </span>
        </button>

        <NotificationBell />
      </div>
    </header>
  );
}