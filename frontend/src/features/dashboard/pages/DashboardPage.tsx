import { useAuth } from "../../auth/context/AuthContext";
import { useDashboardSummary } from "../hooks/useDashboard";
import { Card, CardHeader, CardTitle, CardContent } from "../../../components/ui/card";
import { MetricCard } from "../../../components/ui/metric-card";
import { PageHeader } from "../../../components/shared/PageHeader";
import { Button } from "../../../components/ui/button";

export default function DashboardPage() {
  const { user } = useAuth();
  const { data, isLoading } = useDashboardSummary();

  if (isLoading || !data) {
    return <div className="text-ink-secondary text-sm">Memuat dashboard...</div>;
  }

  return (
    <div>
      <PageHeader
        title={`Selamat datang, ${user?.full_name?.split(" ")[0]}`}
        subtitle="Ringkasan portofolio proyek Anda"
        actions={
          <>
            <Button variant="outline" size="sm">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" />
                <path d="M7 10l5 5 5-5M12 15V3" />
              </svg>
              Export
            </Button>
            <Button variant="primary" size="sm">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M12 5v14M5 12h14" />
              </svg>
              Request baru
            </Button>
          </>
        }
      />

      <div className="grid grid-cols-2 md:grid-cols-4 gap-3.5 mb-7">
        <MetricCard
          label="Total Projects"
          value={data.total_projects}
          iconColor="blue"
          icon={
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M3 7l4-4h6l4 4h4v13H3z" />
            </svg>
          }
        />
        <MetricCard
          label="Active Projects"
          value={data.active_projects}
          iconColor="green"
          icon={
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="12" cy="12" r="9" />
              <path d="M9 12l2 2 4-4" />
            </svg>
          }
        />
        <MetricCard
          label="Pending Requests"
          value={data.pending_requests}
          iconColor="amber"
          delta={{ neutral: "Menunggu review" }}
          icon={
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="12" cy="12" r="9" />
              <path d="M12 7v5l3 3" />
            </svg>
          }
        />
        <MetricCard
          label="Overdue Tasks"
          value={data.overdue_tasks}
          iconColor="red"
          delta={data.overdue_tasks > 0 ? { direction: "down", text: "Perlu perhatian" } : undefined}
          icon={
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="12" cy="12" r="9" />
              <path d="M12 8v4l2.5 2.5" />
            </svg>
          }
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-[1.4fr_1fr] gap-4">
        <Card>
          <CardHeader>
            <CardTitle>Budget usage</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-baseline gap-2 mb-3">
              <span className="text-[28px] font-semibold tracking-tight">
                {data.total_budget_usage_percentage.toFixed(1)}%
              </span>
              <span className="text-xs text-ink-tertiary">rata-rata seluruh project aktif</span>
            </div>
            <div className="h-2 rounded-full bg-surface-tertiary overflow-hidden">
              <div
                className="h-full rounded-full transition-all"
                style={{
                  width: `${Math.min(data.total_budget_usage_percentage, 100)}%`,
                  background:
                    data.total_budget_usage_percentage >= 100
                      ? "linear-gradient(90deg, #2563EB, #DC2626)"
                      : data.total_budget_usage_percentage >= 80
                      ? "linear-gradient(90deg, #2563EB, #D97706)"
                      : "#2563EB",
                }}
              />
            </div>
          </CardContent>
        </Card>

        {data.recent_activities && data.recent_activities.length > 0 && (
          <Card>
            <CardHeader>
              <CardTitle>Aktivitas terbaru</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex flex-col gap-3.5">
                {data.recent_activities.map((activity, idx) => (
                  <div key={idx} className="flex gap-2.5">
                    <span className="h-1.5 w-1.5 rounded-full bg-primary-600 mt-1.5 flex-shrink-0" />
                    <div className="min-w-0">
                      <p className="text-[12.5px] m-0">
                        <span className="font-semibold">{activity.module}</span> — {activity.action}
                      </p>
                      <p className="text-[11.5px] text-ink-tertiary mt-0.5 m-0">
                        {new Date(activity.created_at).toLocaleString("id-ID")}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}