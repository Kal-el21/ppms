export interface ApprovalWorkflow {
  id: number;
  name: string;
  is_active: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface ApprovalLevel {
  id: number;
  workflow_id: number;
  level_order: number;
  role_required: string;
  created_at?: string;
}
