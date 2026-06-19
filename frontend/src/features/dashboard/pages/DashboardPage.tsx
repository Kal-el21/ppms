import { useAuth } from "../../auth/context/AuthContext";
import { useDashboardSummary } from "../hooks/useDashboard";
import { Card, CardContent, CardHeader, CardTitle } from "../../../components/ui/card";

export default function DashboardPage() {
  const { user } = useAuth();
  const { data, isLoading } = useDashboardSummary();

  if (isLoading || !data) return <div>Loading dashboard...</div>;

  const metrics = [
    { label: "Total Projects", value: data.total_projects },
    { label: "Active Projects", value: data.active_projects },
    { label: "Completed Projects", value: data.completed_projects },
    { label: "Pending Requests", value: data.pending_requests },
    { label: "Overdue Tasks", value: data.overdue_tasks },
    { label: "Budget Usage", value: `${data.total_budget_usage_percentage.toFixed(1)}%` },
  ];

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold">Welcome, {user?.full_name}</h1>

      <div className="grid grid-cols-2 gap-4 md:grid-cols-3">
        {metrics.map((m) => (
          <Card key={m.label}>
            <CardContent className="p-4">
              <p className="text-sm text-slate-500">{m.label}</p>
              <p className="text-2xl font-semibold">{m.value}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      {data.recent_activities && data.recent_activities.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Recent Activities</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            {data.recent_activities.map((activity, idx) => (
              <div key={idx} className="flex justify-between border-b pb-2 text-sm">
                <span>
                  {activity.module} — {activity.action}
                </span>
                <span className="text-slate-400">{new Date(activity.created_at).toLocaleString()}</span>
              </div>
            ))}
          </CardContent>
        </Card>
      )}
    </div>
  );
}