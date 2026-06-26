export type RequestStatus =
  | "DRAFT"
  | "SUBMITTED"
  | "UNDER_REVIEW"
  | "APPROVED"
  | "REJECTED"
  | "REVISION_REQUESTED"
  | "REVISED";

export interface ProjectRequest {
  id: number;
  requester_id: number;
  title: string;
  description: string;
  business_goal: string;
  expected_outcome: string;
  estimated_budget: number;
  status: RequestStatus;
  submitted_at: string | null;
  version: number;
  created_at: string;
  updated_at: string;
}

export interface RequestRevision {
  id: number;
  revision_number: number;
  title: string;
  description: string;
  estimated_budget: number;
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

export interface CreateDraftPayload {
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
