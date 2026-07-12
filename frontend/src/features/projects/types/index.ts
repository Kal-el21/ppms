export type ProjectStatus = "PLANNED" | "ACTIVE" | "ON_HOLD" | "COMPLETED" | "CANCELLED";
export type ProjectRole = "PROJECT_MANAGER" | "MEMBER" | "OBSERVER";
export type InitiationType = "NEW_INITIATIVE" | "RENEWAL" | "ENHANCEMENT";
export type ProjectPriority = "LOW" | "MEDIUM" | "HIGH" | "URGENT";
export type HealthStatus = "GREEN" | "YELLOW" | "RED";
export type BudgetType = "CAPEX" | "OPEX";

export interface Project {
  id: number;
  project_request_id?: number | null;
  project_code?: string;
  name: string;
  description: string;
  category: string;
  initiation_type: InitiationType | null;
  priority: ProjectPriority;
  notes: string;
  start_date: string | null;
  end_date: string | null;
  status: ProjectStatus;
  progress: number;
  health: HealthStatus;
  budget_type?: BudgetType | null;
  budget_allocated?: number;
  budget_used?: number;
  version: number;
  created_at: string;
}

export interface ProjectMember {
  id: number;
  user_id: number;
  project_role: ProjectRole;
  status: string;
  joined_at: string;
}

export interface ProjectDeadline {
  id: number;
  project_code: string;
  name: string;
  end_date: string | null;
  status: ProjectStatus;
  progress: number;
  days_remaining: number;
}

export interface CreateProjectDirectRequest {
  name: string;
  description?: string;
  category?: string;
  initiation_type?: InitiationType | null;
  priority: ProjectPriority;
  notes?: string;
  start_date?: string | null;
  end_date?: string | null;
  budget_type?: BudgetType | null;
  budget_name?: string;
  allocated_budget?: number;
}
