import { createContext, useContext, useState, useCallback, type ReactNode } from "react";

type ToastType = "success" | "error" | "warning" | "info";

interface Toast {
  id: string;
  type: ToastType;
  title: string;
  description?: string;
}

interface ToastContextValue {
  toasts: Toast[];
  toast: {
    success: (title: string, description?: string) => void;
    error: (title: string, description?: string) => void;
    warning: (title: string, description?: string) => void;
    info: (title: string, description?: string) => void;
  };
  dismiss: (id: string) => void;
}

const ToastContext = createContext<ToastContextValue | undefined>(undefined);

const toastConfig: Record<ToastType, { icon: string; classes: string }> = {
  success: {
    icon: "M20 6L9 17l-5-5",
    classes:
      "bg-success-50 border-success-200 dark:bg-success-700/10 dark:border-success-700/30",
  },
  error: {
    icon: "M18 6L6 18M6 6l12 12",
    classes:
      "bg-danger-50 border-danger-200 dark:bg-danger-900/20 dark:border-danger-900/40",
  },
  warning: {
    icon: "M12 9v4M12 16h.01M10.3 3.6L1.6 18a2 2 0 001.7 3h17.4a2 2 0 001.7-3L13.7 3.6a2 2 0 00-3.4 0z",
    classes:
      "bg-warning-50 border-warning-200 dark:bg-warning-700/10 dark:border-warning-700/30",
  },
  info: {
    icon: "M12 16v-4M12 8h.01M12 22C6.5 22 2 17.5 2 12S6.5 2 12 2s10 4.5 10 10-4.5 10-10 10z",
    classes:
      "bg-primary-50 border-primary-200 dark:bg-primary-900/20 dark:border-primary-900/40",
  },
};

const iconColor: Record<ToastType, string> = {
  success: "text-success-600",
  error: "text-danger-600",
  warning: "text-warning-600",
  info: "text-primary-600",
};

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const dismiss = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const addToast = useCallback((type: ToastType, title: string, description?: string) => {
    const id = Math.random().toString(36).slice(2);
    setToasts((prev) => [...prev.slice(-4), { id, type, title, description }]);
    setTimeout(() => dismiss(id), 4000);
  }, [dismiss]);

  const toast = {
    success: (title: string, desc?: string) => addToast("success", title, desc),
    error: (title: string, desc?: string) => addToast("error", title, desc),
    warning: (title: string, desc?: string) => addToast("warning", title, desc),
    info: (title: string, desc?: string) => addToast("info", title, desc),
  };

  return (
    <ToastContext.Provider value={{ toasts, toast, dismiss }}>
      {children}
      <div className="fixed bottom-5 right-5 z-[100] flex flex-col gap-2 w-[340px]">
        {toasts.map((t) => {
          const config = toastConfig[t.type];
          return (
            <div
              key={t.id}
              className={`flex items-start gap-3 rounded-lg border px-4 py-3 shadow-md animate-in slide-in-from-bottom-2 ${config.classes}`}
            >
              <div className={`mt-0.5 flex-shrink-0 ${iconColor[t.type]}`}>
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                  <path d={config.icon} />
                </svg>
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-[13px] font-semibold text-ink-primary m-0">{t.title}</p>
                {t.description && (
                  <p className="text-[12px] text-ink-secondary m-0 mt-0.5">{t.description}</p>
                )}
              </div>
              <button
                onClick={() => dismiss(t.id)}
                className="flex-shrink-0 text-ink-tertiary hover:text-ink-secondary transition-colors ml-1"
              >
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                  <path d="M18 6L6 18M6 6l12 12" />
                </svg>
              </button>
            </div>
          );
        })}
      </div>
    </ToastContext.Provider>
  );
}

export function useToast() {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error("useToast must be used within ToastProvider");
  return ctx.toast;
}