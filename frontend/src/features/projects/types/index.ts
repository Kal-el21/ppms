export type ProjectStatus = "PLANNED" | "ACTIVE" | "ON_HOLD" | "COMPLETED" | "CANCELLED";
export type ProjectRole = "PROJECT_MANAGER" | "MEMBER" | "OBSERVER";

export interface Project {
  id: number;
  name: string;
  description: string;
  start_date: string | null;
  end_date: string | null;
  status: ProjectStatus;
  progress: number;
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