import { type ReactNode } from "react";

interface EmptyStateProps {
  icon: ReactNode;
  title: string;
  description?: string;
  action?: ReactNode;
}

export function EmptyState({ icon, title, description, action }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center text-center py-12 px-6">
      <div className="h-11 w-11 rounded-lg bg-surface-tertiary flex items-center justify-center text-ink-tertiary mb-3">
        {icon}
      </div>
      <p className="text-[13.5px] font-medium text-ink-primary mb-1">{title}</p>
      {description && <p className="text-[12.5px] text-ink-tertiary max-w-xs mb-4">{description}</p>}
      {action}
    </div>
  );
}