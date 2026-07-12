import { useAuth } from "../../auth/context/AuthContext";
import { useDashboardSummary } from "../hooks/useDashboard";
import { useProjectDeadlines } from "../../projects/hooks/useProjects";
import { Card, CardHeader, CardTitle, CardContent } from "../../../components/ui/card";
import { MetricCard } from "../../../components/ui/metric-card";
import { PageHeader } from "../../../components/shared/PageHeader";
import { Button } from "../../../components/ui/button";
import { StatusBadge, getStatusColor } from "../../../components/ui/status-badge";
import { DoughnutChart, BarChart, GroupedBarChart } from "../../../components/ui/charts";
import { MetricsSkeleton } from "@/components/ui/skeleton";
import { useNavigate } from "react-router-dom";

const formatCurrency = (value?: number | null) =>
  `Rp ${Math.round(value || 0).toLocaleString("id-ID")}`;

const formatDate = (value?: string | null) =>
  value ? new Date(value).toLocaleDateString("id-ID", { day: "numeric", month: "short", year: "numeric" }) : "—";

export default function DashboardPage() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const { data, isLoading } = useDashboardSummary();
  const { data: deadlines } = useProjectDeadlines("90");

  if (isLoading) {
    return (
      <div>
        <div className="mb-6">
          <div className="h-7 w-48 bg-surface-tertiary rounded-md animate-pulse mb-2" />
          <div className="h-4 w-64 bg-surface-tertiary rounded-md animate-pulse" />
        </div>
        <MetricsSkeleton />
      </div>
    );
  }

  const totalProjects = data?.total_projects ?? 0;
  const activeProjects = data?.active_projects ?? 0;
  const onHoldProjects = data?.on_hold_projects ?? 0;
  const completedProjects = data?.completed_projects ?? 0;
  const overdueProjects = data?.overdue_projects ?? 0;
  const pendingRequests = data?.pending_requests ?? 0;
  const overdueTasks = data?.overdue_tasks ?? 0;

  const statusData = [
    { label: "Planned", value: data?.by_status?.PLANNED ?? 0, color: "#94A3B8" },
    { label: "Active", value: data?.by_status?.ACTIVE ?? 0, color: "#2563EB" },
    { label: "On Hold", value: data?.by_status?.ON_HOLD ?? 0, color: "#D97706" },
    { label: "Completed", value: data?.by_status?.COMPLETED ?? 0, color: "#059669" },
    { label: "Cancelled", value: data?.by_status?.CANCELLED ?? 0, color: "#DC2626" },
  ].filter((d) => d.value > 0);

  const initiationData = [
    { label: "New", value: data?.by_initiation?.NEW_INITIATIVE ?? 0, color: "#2563EB" },
    { label: "Renewal", value: data?.by_initiation?.RENEWAL ?? 0, color: "#7C3AED" },
    { label: "Enhancement", value: data?.by_initiation?.ENHANCEMENT ?? 0, color: "#0891B2" },
  ].filter((d) => d.value > 0);

  const priorityData = [
    { label: "Low", value: data?.by_priority?.LOW ?? 0, color: "#94A3B8" },
    { label: "Medium", value: data?.by_priority?.MEDIUM ?? 0, color: "#2563EB" },
    { label: "High", value: data?.by_priority?.HIGH ?? 0, color: "#D97706" },
    { label: "Urgent", value: data?.by_priority?.URGENT ?? 0, color: "#DC2626" },
  ];

  const healthData = [
    { label: "Green", value: data?.health_green ?? 0, color: "#059669" },
    { label: "Yellow", value: data?.health_yellow ?? 0, color: "#D97706" },
    { label: "Red", value: data?.health_red ?? 0, color: "#DC2626" },
  ].filter((d) => d.value > 0);

  const budgetData = [
    { label: "CAPEX", value: Math.round(data?.capex_allocated ?? 0), color: "#2563EB" },
    { label: "CAPEX Used", value: Math.round(data?.capex_used ?? 0), color: "#1E40AF" },
    { label: "OPEX", value: Math.round(data?.opex_allocated ?? 0), color: "#7C3AED" },
    { label: "OPEX Used", value: Math.round(data?.opex_used ?? 0), color: "#5B21B6" },
  ];

  const budgetMasterData = (data?.budget_master ?? []).map((item) => ({
    label: item.budget_name,
    value: Math.round(item.allocated),
    color: item.budget_type === "CAPEX" ? "#2563EB" : "#7C3AED",
    value2: Math.round(item.used),
    color2: item.budget_type === "CAPEX" ? "#1E40AF" : "#5B21B6",
  }));

  const absorptionData = (data?.absorption ?? []).map((item) => ({
    label: item.project_code,
    value: Math.round(item.usage_percentage),
    color: item.budget_type === "CAPEX" ? "#2563EB" : "#7C3AED",
  }));

  const deadlineList = deadlines ?? data?.upcoming_deadlines ?? [];

  return (
    <div>
      <PageHeader
        title={`Selamat datang, ${user?.full_name?.split(" ")[0] ?? ""}`}
        subtitle="Ringkasan portofolio proyek Anda"
        actions={
          <Button variant="primary" size="sm" onClick={() => navigate("/project-requests/new")}>
            Request baru
          </Button>
        }
      />

      <div className="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-7 gap-3.5 mb-7">
        <MetricCard label="Total Projects" value={totalProjects} iconColor="blue"
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M3 7l4-4h6l4 4h4v13H3z"/></svg>}
        />
        <MetricCard label="Active" value={activeProjects} iconColor="green"
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="12" r="9"/><path d="M9 12l2 2 4-4"/></svg>}
        />
        <MetricCard label="On Hold" value={onHoldProjects} iconColor="amber"
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="12" r="9"/><path d="M10 9v6M14 9v6"/></svg>}
        />
        <MetricCard label="Completed" value={completedProjects} iconColor="teal"
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M20 6L9 17l-5-5"/></svg>}
        />
        <MetricCard label="Overdue" value={overdueProjects} iconColor="red"
          delta={overdueProjects > 0 ? { direction: "down", text: "Perlu perhatian" } : undefined}
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="12" r="9"/><path d="M12 8v4l2.5 2.5"/></svg>}
        />
        <MetricCard label="Pending Requests" value={pendingRequests} iconColor="amber"
          delta={{ neutral: "Menunggu review" }}
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 3"/></svg>}
        />
        <MetricCard label="Overdue Tasks" value={overdueTasks} iconColor="red"
          delta={overdueTasks > 0 ? { direction: "down", text: "Perlu perhatian" } : undefined}
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M9 11l3 3L22 4M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11"/></svg>}
        />
        <MetricCard label="CAPEX Allocated" value={formatCurrency(data?.capex_allocated)} iconColor="blue"
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><rect x="2" y="5" width="20" height="14" rx="2"/><path d="M2 10h20"/></svg>}
        />
        <MetricCard label="CAPEX Used" value={formatCurrency(data?.capex_used)} iconColor="indigo"
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M12 2v20M2 12h20"/></svg>}
        />
        <MetricCard label="OPEX Allocated" value={formatCurrency(data?.opex_allocated)} iconColor="teal"
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><rect x="2" y="5" width="20" height="14" rx="2"/><path d="M2 10h20"/></svg>}
        />
        <MetricCard label="OPEX Used" value={formatCurrency(data?.opex_used)} iconColor="emerald"
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M12 2v20M2 12h20"/></svg>}
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-4">
        <Card>
          <CardHeader><CardTitle>Status Distribution</CardTitle></CardHeader>
          <CardContent className="flex justify-center">
            <DoughnutChart data={statusData} />
          </CardContent>
        </Card>
        <Card>
          <CardHeader><CardTitle>Initiation Distribution</CardTitle></CardHeader>
          <CardContent className="flex justify-center">
            <DoughnutChart data={initiationData} />
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 mb-4">
        <Card>
          <CardHeader><CardTitle>Priority Distribution</CardTitle></CardHeader>
          <CardContent>
            <BarChart data={priorityData} />
          </CardContent>
        </Card>
        <Card className="lg:col-span-2">
          <CardHeader><CardTitle>Budget Overview (CAPEX / OPEX)</CardTitle></CardHeader>
          <CardContent>
            <BarChart data={budgetData} height={150} />
            <div className="grid grid-cols-2 gap-3 mt-3 text-[12.5px]">
              <div className="rounded-md bg-surface-secondary p-2.5">
                <p className="text-ink-tertiary mb-0.5">CAPEX</p>
                <p className="font-semibold m-0">{formatCurrency(data?.capex_allocated)}</p>
                <p className="text-ink-tertiary m-0">Used: {formatCurrency(data?.capex_used)}</p>
              </div>
              <div className="rounded-md bg-surface-secondary p-2.5">
                <p className="text-ink-tertiary mb-0.5">OPEX</p>
                <p className="font-semibold m-0">{formatCurrency(data?.opex_allocated)}</p>
                <p className="text-ink-tertiary m-0">Used: {formatCurrency(data?.opex_used)}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-4">
        <Card>
          <CardHeader><CardTitle>Budget Master (per Mata Anggaran)</CardTitle></CardHeader>
          <CardContent>
            <GroupedBarChart data={budgetMasterData} height={180} />
          </CardContent>
        </Card>
        <Card>
          <CardHeader><CardTitle>Absorption % (per Project)</CardTitle></CardHeader>
          <CardContent>
            <BarChart data={absorptionData} height={180} />
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-[1.4fr_1fr] gap-4">
        <Card>
          <CardHeader className="flex items-center justify-between">
            <CardTitle>Deadline Projects</CardTitle>
            <Button variant="ghost" size="sm" onClick={() => navigate("/projects/deadline")}>
              Lihat semua
            </Button>
          </CardHeader>
          <CardContent>
            {deadlineList.length === 0 ? (
              <p className="text-[12.5px] text-ink-tertiary">Tidak ada project yang mendekati deadline.</p>
            ) : (
              <div className="flex flex-col">
                {deadlineList.slice(0, 6).map((d) => (
                  <div
                    key={d.id}
                    className="flex items-center gap-3 py-2.5 border-b border-border last:border-b-0 cursor-pointer"
                    onClick={() => navigate(`/projects/${d.id}`)}
                  >
                    <div className="flex-1 min-w-0">
                      <p className="text-[13px] font-medium m-0 truncate">{d.name}</p>
                      <p className="text-[11.5px] text-ink-tertiary m-0">
                        {d.project_code} · {formatDate(d.end_date)}
                      </p>
                    </div>
                    <StatusBadge color={getStatusColor(d.status)}>{d.status}</StatusBadge>
                    <span
                      className={`text-xs font-semibold ${
                        d.days_remaining < 0 ? "text-danger-600" : d.days_remaining <= 30 ? "text-warning-600" : "text-ink-secondary"
                      }`}
                    >
                      {d.days_remaining < 0 ? `${Math.abs(d.days_remaining)} hari lalu` : `${d.days_remaining} hari`}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle>Health Distribution</CardTitle></CardHeader>
          <CardContent className="flex justify-center">
            <DoughnutChart data={healthData} />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
