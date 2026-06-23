import { useState, useRef, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useNotifications, useMarkAsRead, useMarkAllAsRead } from "../../features/notifications/hooks/useNotifications";

export default function NotificationBell() {
  const [open, setOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate();
  const { data } = useNotifications(1, 10);
  const { mutate: markAsRead } = useMarkAsRead();
  const { mutate: markAllAsRead } = useMarkAllAsRead();

  const unreadCount = data?.meta.unread_count ?? 0;

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleClick = (notif: { id: number; is_read: boolean; action_url: string }) => {
    if (!notif.is_read) markAsRead(notif.id);
    setOpen(false);
    if (notif.action_url) navigate(notif.action_url);
  };

  return (
    <div ref={containerRef} className="relative">
      <button
        onClick={() => setOpen((v) => !v)}
        className="relative h-8 w-8 rounded-md flex items-center justify-center text-ink-secondary hover:bg-surface-secondary border border-transparent hover:border-border transition-colors"
      >
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
          <path d="M18 8a6 6 0 00-12 0c0 7-3 9-3 9h18s-3-2-3-9" />
          <path d="M13.7 21a2 2 0 01-3.4 0" />
        </svg>
        {unreadCount > 0 && (
          <span className="absolute top-1 right-1 h-[7px] w-[7px] rounded-full bg-danger-600 ring-2 ring-surface" />
        )}
      </button>

      {open && (
        <div className="absolute right-0 z-50 mt-1.5 w-80 rounded-lg border border-border bg-surface shadow-lg overflow-hidden">
          <div className="flex items-center justify-between px-3.5 py-3 border-b border-border">
            <span className="text-[13px] font-semibold">Notifications</span>
            <button
              onClick={() => markAllAsRead()}
              className="text-[11.5px] font-medium text-primary-600 hover:text-primary-700"
            >
              Tandai semua dibaca
            </button>
          </div>

          <div className="max-h-96 overflow-y-auto">
            {data?.data.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-10 px-4">
                <div className="h-9 w-9 rounded-full bg-surface-tertiary flex items-center justify-center text-ink-tertiary mb-2">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M18 8a6 6 0 00-12 0c0 7-3 9-3 9h18s-3-2-3-9" />
                  </svg>
                </div>
                <p className="text-[12.5px] text-ink-tertiary m-0">Tidak ada notifikasi</p>
              </div>
            ) : (
              data?.data.map((notif) => (
                <button
                  key={notif.id}
                  onClick={() => handleClick(notif)}
                  className={`relative block w-full text-left px-3.5 py-3 border-b border-border last:border-b-0 hover:bg-surface-secondary transition-colors ${
                    !notif.is_read ? "bg-primary-50/50 dark:bg-primary-900/10" : ""
                  }`}
                >
                  {!notif.is_read && (
                    <span className="absolute left-1.5 top-4 h-1.5 w-1.5 rounded-full bg-primary-600" />
                  )}
                  <p className="text-[12.5px] font-medium m-0 pl-3">{notif.title}</p>
                  <p className="text-[11.5px] text-ink-secondary m-0 mt-0.5 pl-3 line-clamp-2">{notif.message}</p>
                  <p className="text-[11px] text-ink-tertiary m-0 mt-1 pl-3">
                    {new Date(notif.created_at).toLocaleString("id-ID")}
                  </p>
                </button>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  );
}