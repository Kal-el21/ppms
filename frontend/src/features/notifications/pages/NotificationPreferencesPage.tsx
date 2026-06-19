import { useNotificationPreferences, useUpdatePreference } from "../hooks/useNotifications";
import { Card, CardContent, CardHeader, CardTitle } from "../../../components/ui/card";

const typeLabels: Record<string, string> = {
  REQUEST_SUBMITTED: "Request Submitted",
  REQUEST_APPROVED: "Request Approved",
  REQUEST_REJECTED: "Request Rejected",
  TASK_ASSIGNED: "Task Assigned",
  TASK_COMPLETED: "Task Completed",
  BUDGET_WARNING: "Budget Warning",
  BUDGET_OVER_LIMIT: "Budget Over Limit",
  HANDOVER_SENT: "Handover Sent",
  HANDOVER_RECEIVED: "Handover Received",
};

export default function NotificationPreferencesPage() {
  const { data: preferences, isLoading } = useNotificationPreferences();
  const { mutate: updatePref } = useUpdatePreference();

  if (isLoading) return <div>Loading...</div>;

  return (
    <Card className="max-w-md">
      <CardHeader>
        <CardTitle>Notification Preferences</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        {preferences?.map((p) => (
          <div key={p.type} className="flex items-center justify-between">
            <span className="text-sm">{typeLabels[p.type] || p.type}</span>
            <input
              type="checkbox"
              checked={p.enabled}
              onChange={(e) => updatePref({ type: p.type, enabled: e.target.checked })}
              className="h-4 w-4"
            />
          </div>
        ))}
      </CardContent>
    </Card>
  );
}