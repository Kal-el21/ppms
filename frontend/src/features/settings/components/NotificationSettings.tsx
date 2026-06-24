import { useState } from "react";
import { useAuth } from "../../auth/context/AuthContext";
import axiosInstance from "../../../api/axiosInstance";
import { useNotificationPreferences, useUpdatePreference } from "../../notifications/hooks/useNotifications";
import { Card, CardContent } from "../../../components/ui/card";
import { useToast } from "../../../components/ui/toast";

const typeMeta: Record<string, { label: string; desc: string; group: string }> = {
  REQUEST_SUBMITTED: { label: "Request Disubmit", desc: "Saat ada request baru menunggu review Anda", group: "Project Request" },
  REQUEST_APPROVED: { label: "Request Disetujui", desc: "Saat request Anda disetujui", group: "Project Request" },
  REQUEST_REJECTED: { label: "Request Ditolak/Revisi", desc: "Saat request Anda perlu revisi atau ditolak", group: "Project Request" },
  TASK_ASSIGNED: { label: "Task Ditugaskan", desc: "Saat Anda ditugaskan ke task baru", group: "Task" },
  TASK_COMPLETED: { label: "Task Selesai", desc: "Saat task di project Anda selesai", group: "Task" },
  BUDGET_WARNING: { label: "Peringatan Budget 80%", desc: "Saat penggunaan budget mencapai 80%", group: "Budget" },
  BUDGET_OVER_LIMIT: { label: "Budget Melebihi Limit", desc: "Saat penggunaan budget melebihi 100%", group: "Budget" },
  HANDOVER_SENT: { label: "Handover Dikirim", desc: "Saat dokumen dikirimkan ke Anda", group: "Handover" },
  HANDOVER_RECEIVED: { label: "Handover Diterima", desc: "Saat handover Anda dikonfirmasi diterima", group: "Handover" },
};

function Toggle({ checked, onChange, disabled = false }: { checked: boolean; onChange: (v: boolean) => void; disabled?: boolean }) {
  return (
    <button
      type="button"
      onClick={() => onChange(!checked)}
      disabled={disabled}
      className={`relative h-5 w-9 rounded-full transition-colors flex-shrink-0 ${checked ? "bg-primary-600" : "bg-surface-tertiary"} ${disabled ? "opacity-50 cursor-not-allowed" : "cursor-pointer"}`}
    >
      <span className={`absolute top-0.5 h-4 w-4 rounded-full bg-white shadow-sm transition-transform ${checked ? "translate-x-[18px]" : "translate-x-0.5"}`} />
    </button>
  );
}

export default function NotificationSettings() {
  const { user, login } = useAuth();
  const toast = useToast();
  const { data: preferences } = useNotificationPreferences();
  const { mutate: updatePref } = useUpdatePreference();
  const [togglingEmail, setTogglingEmail] = useState(false);

  const handleToggleEmail = async (enabled: boolean) => {
    setTogglingEmail(true);
    try {
      await axiosInstance.post("/me/toggle-email-notification", { enabled });
      login({ ...user!, email_notification_enabled: enabled });
      toast.success(enabled ? "Notifikasi email diaktifkan" : "Notifikasi email dinonaktifkan");
    } catch (err: any) {
      toast.error("Gagal mengubah preferensi email", err?.friendlyMessage);
    } finally {
      setTogglingEmail(false);
    }
  };

  const grouped = (preferences || []).reduce<Record<string, typeof preferences>>((acc, p) => {
    const group = typeMeta[p.type]?.group || "Lainnya";
    if (!acc[group]) acc[group] = [];
    acc[group]!.push(p);
    return acc;
  }, {});

  return (
    <div className="space-y-4">
      <Card>
        <CardContent className="pt-5">
          <p className="text-[11.5px] font-semibold uppercase tracking-wide text-ink-tertiary mb-3">Channel Notifikasi</p>
          <div className="space-y-1">
            <div className="flex items-center justify-between py-2.5">
              <div>
                <p className="text-[13px] font-medium m-0">Notifikasi In-App</p>
                <p className="text-[11.5px] text-ink-tertiary m-0 mt-0.5">Tampil di lonceng notifikasi di dalam aplikasi</p>
              </div>
              <Toggle checked={true} onChange={() => {}} />
            </div>
            <div className="flex items-center justify-between py-2.5">
              <div>
                <p className="text-[13px] font-medium m-0">Notifikasi Email</p>
                <p className="text-[11.5px] text-ink-tertiary m-0 mt-0.5">
                  Kirim ke <span className="font-medium">{user?.email}</span> via Gmail
                </p>
              </div>
              <Toggle
                checked={user?.email_notification_enabled ?? false}
                onChange={handleToggleEmail}
                disabled={togglingEmail}
              />
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="space-y-3">
        {Object.entries(grouped).map(([group, items]) => (
          <Card key={group}>
            <CardContent className="pt-5">
              <p className="text-[11.5px] font-semibold uppercase tracking-wide text-ink-tertiary mb-3">{group}</p>
              <div className="space-y-1">
                {items?.map((p) => {
                  const meta = typeMeta[p.type];
                  return (
                    <div key={p.type} className="flex items-center justify-between py-2.5">
                      <div className="min-w-0 pr-6">
                        <p className="text-[13px] font-medium m-0">{meta?.label || p.type}</p>
                        {meta?.desc && <p className="text-[11.5px] text-ink-tertiary m-0 mt-0.5">{meta.desc}</p>}
                      </div>
                      <Toggle
                        checked={p.enabled}
                        onChange={(v) => updatePref({ type: p.type, enabled: v })}
                      />
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}