import { type ReactNode } from "react";

type BadgeColor = "blue" | "red" | "green" | "amber" | "gray";

interface StatusBadgeProps {
  color: BadgeColor;
  children: ReactNode;
}

const colorMap: Record<BadgeColor, string> = {
  blue: "bg-primary-50 text-primary-700 dark:bg-primary-900/40 dark:text-primary-400",
  red: "bg-danger-50 text-danger-700 dark:bg-danger-900/40 dark:text-danger-400",
  green: "bg-success-50 text-success-700 dark:bg-success-700/20 dark:text-success-500",
  amber: "bg-warning-50 text-warning-700 dark:bg-warning-700/20 dark:text-warning-500",
  gray: "bg-surface-tertiary text-ink-secondary",
};

const dotMap: Record<BadgeColor, string> = {
  blue: "bg-primary-600",
  red: "bg-danger-600",
  green: "bg-success-600",
  amber: "bg-warning-600",
  gray: "bg-ink-tertiary",
};

export function StatusBadge({ color, children }: StatusBadgeProps) {
  return (
    <span
      className={`inline-flex items-center gap-1.5 rounded-md px-2 py-0.5 text-[11.5px] font-semibold leading-relaxed ${colorMap[color]}`}
    >
      <span className={`h-1.5 w-1.5 rounded-full ${dotMap[color]}`} />
      {children}
    </span>
  );
}

// Helper: mapping status string ke warna, dipakai konsisten di semua halaman
export function getStatusColor(status: string): BadgeColor {
  const map: Record<string, BadgeColor> = {
    // Project Request
    DRAFT: "gray",
    SUBMITTED: "blue",
    UNDER_REVIEW: "amber",
    APPROVED: "green",
    REJECTED: "red",
    REVISION_REQUESTED: "amber",
    REVISED: "gray",
    // Project
    PLANNED: "gray",
    ACTIVE: "blue",
    ON_HOLD: "amber",
    COMPLETED: "green",
    CANCELLED: "red",
    // Task
    TODO: "gray",
    IN_PROGRESS: "blue",
    DONE: "green",
    // Handover
    PENDING: "amber",
    RECEIVED: "green",
    RETURNED: "red",
    // Priority
    LOW: "gray",
    MEDIUM: "blue",
    HIGH: "amber",
    URGENT: "red",
  };
  return map[status] || "gray";
}
