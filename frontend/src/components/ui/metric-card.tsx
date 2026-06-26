import { type ReactNode } from "react";

interface MetricCardProps {
  label: string;
  value: string | number;
  icon: ReactNode;
  iconColor: "blue" | "green" | "amber" | "red" | "teal" | "indigo" | "emerald";
  delta?: { direction: "up" | "down"; text: string } | { neutral: string };
}

const iconColorMap: Record<MetricCardProps["iconColor"], string> = {
  blue: "bg-primary-50 text-primary-700 dark:bg-primary-900/40 dark:text-primary-400",
  green: "bg-success-50 text-success-700 dark:bg-success-700/20 dark:text-success-500",
  amber: "bg-warning-50 text-warning-700 dark:bg-warning-700/20 dark:text-warning-500",
  red: "bg-danger-50 text-danger-700 dark:bg-danger-900/40 dark:text-danger-400",
  teal: "bg-teal-50 text-teal-700 dark:bg-teal-900/40 dark:text-teal-400",
  indigo: "bg-indigo-50 text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-400",
  emerald: "bg-emerald-50 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-400",
};

export function MetricCard({ label, value, icon, iconColor, delta }: MetricCardProps) {
  return (
    <div className="rounded-lg border border-border bg-surface p-5">
      <div className="flex items-center justify-between mb-2.5">
        <span className="text-[12.5px] font-medium text-ink-secondary">{label}</span>
        <div className={`h-7 w-7 rounded-md flex items-center justify-center ${iconColorMap[iconColor]}`}>
          {icon}
        </div>
      </div>
      <div className="text-[28px] font-semibold leading-none tracking-tight mb-2">{value}</div>
      {delta && "neutral" in delta ? (
        <div className="text-xs text-ink-tertiary">{delta.neutral}</div>
      ) : delta ? (
        <div
          className={`text-xs font-semibold flex items-center gap-1 ${
            delta.direction === "up" ? "text-success-600" : "text-danger-600"
          }`}
        >
          {delta.direction === "up" ? (
            <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3">
              <path d="M12 19V5M5 12l7-7 7 7" />
            </svg>
          ) : (
            <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="3">
              <path d="M12 5v14M5 12l7 7 7-7" />
            </svg>
          )}
          {delta.text}
        </div>
      ) : null}
    </div>
  );
}
