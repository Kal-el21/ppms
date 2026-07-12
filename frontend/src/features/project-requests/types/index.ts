export type RequestStatus =
  | "DRAFT"
  | "SUBMITTED"
  | "UNDER_REVIEW"
  | "APPROVED"
  | "REJECTED"
  | "REVISION_REQUESTED"
  | "REVISED";

export type InitiationType = "NEW_INITIATIVE" | "RENEWAL" | "ENHANCEMENT";
export type RequestPriority = "LOW" | "MEDIUM" | "HIGH" | "URGENT";
export type BudgetType = "CAPEX" | "OPEX";

export interface ProjectRequest {
  id: number;
  requester_id: number;
  request_number?: string;
  title: string;
  description: string;
  business_goal: string;
  expected_outcome: string;
  estimated_budget: number;
  category: string;
  initiation_type?: InitiationType | null;
  priority: RequestPriority;
  proposed_start_date?: string | null;
  proposed_end_date?: string | null;
  budget_type?: BudgetType | null;
  budget_name: string;
  notes: string;
  status: RequestStatus;
  current_revision?: number;
  submitted_at: string | null;
  approved_at?: string | null;
  rejected_at?: string | null;
  version: number;
  created_at: string;
  updated_at: string;
}

export interface RequestRevision {
  id: number;
  revision_number: number;
  title: string;
  description: string;
  business_goal: string;
  expected_outcome: string;
  estimated_budget: number;
  category: string;
  initiation_type?: InitiationType | null;
  priority: RequestPriority;
  proposed_start_date?: string | null;
  proposed_end_date?: string | null;
  budget_type?: BudgetType | null;
  budget_name: string;
  notes: string;
  revision_reason: string;
  revised_by: number;
  created_at: string;
}

export interface RequestApproval {
  id: number;
  reviewed_by: number;
  action: "APPROVED" | "REJECTED" | "REQUEST_REVISION";
  comment: string;
  project_manager_id?: number | null;
  created_at: string;
}

export interface ProjectRequestMetadataPayload {
  category?: string;
  initiation_type?: string;
  priority?: string;
  proposed_start_date?: string | null;
  proposed_end_date?: string | null;
  budget_type?: string;
  budget_name?: string;
  notes?: string;
}

export interface CreateDraftPayload extends ProjectRequestMetadataPayload {
  title: string;
  description?: string;
  business_goal?: string;
  expected_outcome?: string;
  estimated_budget: number;
}

export interface UpdateDraftPayload extends CreateDraftPayload {
  version: number;
}

export interface ReviewPayload {
  action: "APPROVED" | "REJECTED" | "REQUEST_REVISION";
  comment: string;
  project_manager_id?: number;
}

export interface RevisePayload extends CreateDraftPayload {
  revision_reason: string;
}
