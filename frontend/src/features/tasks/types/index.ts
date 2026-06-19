export type TaskPriority = "LOW" | "MEDIUM" | "HIGH" | "URGENT";
export type TaskStatus = "TODO" | "IN_PROGRESS" | "DONE" | "CANCELLED";

export interface Task {
  id: number;
  project_id: number;
  milestone_id: number | null;
  title: string;
  description: string;
  priority: TaskPriority;
  status: TaskStatus;
  progress: number;
  start_date: string | null;
  due_date: string | null;
  version: number;
  assignee_ids: number[];
}