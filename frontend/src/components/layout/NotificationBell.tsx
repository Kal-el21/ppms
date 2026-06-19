import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useNotifications, useMarkAsRead, useMarkAllAsRead } from "../../features/notifications/hooks/useNotifications";
import { Button } from "../ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import { Badge } from "../ui/badge";

export default function NotificationBell() {
  const [open, setOpen] = useState(false);
  const navigate = useNavigate();
  const { data } = useNotifications(1, 10);
  const { mutate: markAsRead } = useMarkAsRead();
  const { mutate: markAllAsRead } = useMarkAllAsRead();

  const unreadCount = data?.meta.unread_count ?? 0;

  const handleClick = (notif: { id: number; is_read: boolean; action_url: string }) => {
    if (!notif.is_read) {
      markAsRead(notif.id);
    }
    setOpen(false);
    if (notif.action_url) {
      navigate(notif.action_url);
    }
  };

  return (
    <DropdownMenu open={open} onOpenChange={setOpen}>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" className="relative">
          🔔
          {unreadCount > 0 && (
            <Badge variant="destructive" className="absolute -right-2 -top-2 h-5 w-5 rounded-full p-0 text-xs">
              {unreadCount > 9 ? "9+" : unreadCount}
            </Badge>
          )}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-80">
        <div className="flex items-center justify-between border-b p-3">
          <span className="font-medium">Notifications</span>
          <Button variant="ghost" size="sm" onClick={() => markAllAsRead()}>
            Mark all read
          </Button>
        </div>
        <div className="max-h-96 overflow-y-auto">
          {data?.data.length === 0 ? (
            <p className="p-4 text-center text-sm text-slate-400">No notifications</p>
          ) : (
            data?.data.map((notif) => (
              <button
                key={notif.id}
                onClick={() => handleClick(notif)}
                className={`block w-full border-b p-3 text-left text-sm hover:bg-slate-50 ${
                  !notif.is_read ? "bg-blue-50" : ""
                }`}
              >
                <p className="font-medium">{notif.title}</p>
                <p className="text-xs text-slate-500">{notif.message}</p>
                <p className="mt-1 text-xs text-slate-400">{new Date(notif.created_at).toLocaleString()}</p>
              </button>
            ))
          )}
        </div>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}