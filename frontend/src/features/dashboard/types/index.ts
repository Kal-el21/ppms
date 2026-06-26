export interface RecentActivity {
  module: string;
  action: string;
  entity_type: string;
  entity_id: number | null;
  user_id: number | null;
  created_at: string;
}

export interface DashboardSummary {
  total_projects: number;
  active_projects: number;
  completed_projects: number;
  total_tasks: number;
  completed_tasks: number;
  pending_requests: number;
  overdue_tasks: number;
  total_budget_usage_percentage: number;
  recent_activities: RecentActivity[];
}