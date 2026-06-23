import { useState, type ReactNode } from "react";

interface Tab {
  key: string;
  label: string;
  count?: number;
}

interface TabsProps {
  tabs: Tab[];
  defaultTab?: string;
  children: (activeTab: string) => ReactNode;
}

export function Tabs({ tabs, defaultTab, children }: TabsProps) {
  const [active, setActive] = useState(defaultTab || tabs[0]?.key);

  return (
    <div>
      <div className="flex items-center gap-1 border-b border-border mb-5">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActive(tab.key)}
            className={`relative px-3.5 py-2.5 text-[13.5px] font-medium transition-colors ${
              active === tab.key ? "text-primary-700 dark:text-primary-400" : "text-ink-secondary hover:text-ink-primary"
            }`}
          >
            <span className="flex items-center gap-1.5">
              {tab.label}
              {tab.count !== undefined && (
                <span
                  className={`text-[11px] font-semibold px-1.5 py-0.5 rounded-full ${
                    active === tab.key
                      ? "bg-primary-100 text-primary-700 dark:bg-primary-900/40 dark:text-primary-400"
                      : "bg-surface-tertiary text-ink-tertiary"
                  }`}
                >
                  {tab.count}
                </span>
              )}
            </span>
            {active === tab.key && (
              <span className="absolute left-0 right-0 -bottom-px h-[2px] bg-primary-600 rounded-full" />
            )}
          </button>
        ))}
      </div>
      {children(active)}
    </div>
  );
}