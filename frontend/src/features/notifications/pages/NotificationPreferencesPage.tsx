import { useNotificationPreferences, useUpdatePreference } from "../hooks/useNotifications";
import { Card, CardContent } from "../../../components/ui/card";
import { PageHeader } from "../../../components/shared/PageHeader";
import { EmptyState } from "../../../components/shared/EmptyState";
import { CardSkeleton } from "../../../components/ui/skeleton";

const typeMeta: Record<string, { label: string; desc: string; group: string }> = {
  REQUEST_SUBMITTED: { label: "Request Disubmit", desc: "Saat ada request baru menunggu review Anda", group: "Project Request" },
  REQUEST_APPROVED: { label: "Request Disetujui", desc: "Saat request Anda disetujui", group: "Project Request" },
  REQUEST_REJECTED: { label: "Request Ditolak", desc: "Saat request Anda ditolak permanen", group: "Project Request" },
  REVISION_REQUESTED: { label: "Perlu Revisi", desc: "Saat request Anda perlu direvisi", group: "Project Request" },
  REQUEST_REVISED: { label: "Request Direvisi", desc: "Saat request siap direview ulang", group: "Project Request" },
  TASK_ASSIGNED: { label: "Task Ditugaskan", desc: "Saat Anda ditugaskan ke task baru", group: "Task" },
  TASK_COMPLETED: { label: "Task Selesai", desc: "Saat task di project Anda selesai", group: "Task" },
  TASK_OVERDUE: { label: "Task Terlambat", desc: "Saat task yang Anda kerjakan sudah melewati deadline", group: "Task" },
  TASK_DUE_SOON: { label: "Task Akan Jatuh Tempo", desc: "Saat task yang Anda kerjakan akan jatuh tempo dalam 24 jam", group: "Task" },
  MILESTONE_COMPLETED: { label: "Milestone Selesai", desc: "Saat milestone dalam project Anda selesai", group: "Milestone" },
  PROJECT_COMPLETED: { label: "Project Selesai", desc: "Saat project yang Anda ikuti selesai", group: "Project" },
  BUDGET_WARNING: { label: "Peringatan Budget", desc: "Saat penggunaan budget mencapai 80%", group: "Budget" },
  BUDGET_OVER_LIMIT: { label: "Budget Melebihi Limit", desc: "Saat penggunaan budget melebihi 100%", group: "Budget" },
  HANDOVER_CREATED: { label: "Handover Diterima", desc: "Saat ada handover baru yang dikirimkan ke Anda", group: "Handover" },
  HANDOVER_CREATED_CONFIRM: { label: "Handover Terkirim", desc: "Saat handover Anda berhasil dibuat dan dikirim", group: "Handover" },
  HANDOVER_SENT: { label: "Handover Dikirim", desc: "Saat dokumen dikirimkan ke Anda", group: "Handover" },
  HANDOVER_RECEIVED: { label: "Handover Diterima", desc: "Saat handover Anda dikonfirmasi diterima", group: "Handover" },
};

function Toggle({ checked, onChange }: { checked: boolean; onChange: (v: boolean) => void }) {
  return (
    <button
      onClick={() => onChange(!checked)}
      className={`relative h-5 w-9 rounded-full transition-colors flex-shrink-0 ${
        checked ? "bg-primary-600" : "bg-surface-tertiary"
      }`}
    >
      <span
        className={`absolute top-0.5 h-4 w-4 rounded-full bg-white shadow-sm transition-transform ${
          checked ? "translate-x-[18px]" : "translate-x-0.5"
        }`}
      />
    </button>
  );
}

const emptyIcon = (
  <svg className="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
    <path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9M13.73 21a2 2 0 01-3.46 0" />
  </svg>
);

export default function NotificationPreferencesPage() {
  const { data: preferences, isLoading } = useNotificationPreferences();
  const { mutate: updatePref } = useUpdatePreference();

  if (isLoading) {
    return (
      <div className="max-w-2xl space-y-4">
        <PageHeader title="Notification Preferences" subtitle="Atur jenis notifikasi yang ingin Anda terima" />
        {[0, 1, 2].map((i) => (
          <CardSkeleton key={i} />
        ))}
      </div>
    );
  }

  const grouped = (preferences || []).reduce<Record<string, typeof preferences>>((acc, p) => {
    const group = typeMeta[p.type]?.group || "Lainnya";
    if (!acc[group]) acc[group] = [];
    acc[group]!.push(p);
    return acc;
  }, {});

  return (
    <div className="max-w-2xl">
      <PageHeader title="Notification Preferences" subtitle="Atur jenis notifikasi yang ingin Anda terima" />

      <div className="space-y-4">
        {Object.keys(grouped).length === 0 ? (
          <Card>
            <CardContent className="pt-5">
              <EmptyState icon={emptyIcon} title="Belum ada preferensi notifikasi" />
            </CardContent>
          </Card>
        ) : (
          Object.entries(grouped).map(([group, items]) => (
          <Card key={group}>
            <CardContent className="pt-5">
              <p className="text-[11.5px] font-semibold uppercase tracking-wide text-ink-tertiary mb-3">{group}</p>
              <div className="space-y-1">
                {items?.map((p) => {
                  const meta = typeMeta[p.type];
                  return (
                    <div key={p.type} className="flex items-center justify-between py-2.5">
                      <div className="min-w-0 pr-4">
                        <p className="text-[13px] font-medium m-0">{meta?.label || p.type}</p>
                        {meta?.desc && <p className="text-[11.5px] text-ink-tertiary m-0 mt-0.5">{meta.desc}</p>}
                      </div>
                      <Toggle checked={p.enabled} onChange={(v) => updatePref({ type: p.type, enabled: v })} />
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>
          ))
        )}
      </div>
    </div>
  );
}
