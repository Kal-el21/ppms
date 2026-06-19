export type MilestoneStatus = "PLANNED" | "ACTIVE" | "COMPLETED" | "CANCELLED";

export interface Milestone {
  id: number;
  project_id: number;
  name: string;
  description: string;
  order_index: number;
  start_date: string | null;
  end_date: string | null;
  status: MilestoneStatus;
  progress: number;
  version: number;
}