export interface RecentActivity {
  module: string;
  action: string;
  entity_type: string;
  entity_id: number | null;
  user_id: number | null;
  created_at: string;
}

export interface UpcomingDeadline {
  id: number;
  project_code: string;
  name: string;
  end_date: string | null;
  status: string;
  days_remaining: number;
}

export interface DashboardSummary {
  total_projects: number;
  active_projects: number;
  on_hold_projects: number;
  completed_projects: number;
  overdue_projects: number;
  total_tasks: number;
  completed_tasks: number;
  pending_requests: number;
  overdue_tasks: number;
  total_budget_usage_percentage: number;
  by_status: Record<string, number>;
  by_initiation: Record<string, number>;
  by_priority: Record<string, number>;
  health_green: number;
  health_yellow: number;
  health_red: number;
  capex_allocated: number;
  capex_used: number;
  opex_allocated: number;
  opex_used: number;
  upcoming_deadlines: UpcomingDeadline[];
  recent_activities: RecentActivity[];
  budget_master: BudgetMasterItem[];
  absorption: AbsorptionItem[];
}

export interface BudgetMasterItem {
  budget_name: string;
  budget_type: string;
  allocated: number;
  used: number;
  remaining: number;
  usage_percentage: number;
}

export interface AbsorptionItem {
  project_code: string;
  project_name: string;
  budget_type: string;
  allocated: number;
  used: number;
  usage_percentage: number;
}
