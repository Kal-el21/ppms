export type EntityType = "PROJECT_REQUEST" | "PROJECT" | "MILESTONE" | "TASK" | "BUDGET_TRANSACTION" | "HANDOVER";

export interface Attachment {
  id: number;
  entity_type: EntityType;
  entity_id: number;
  original_name: string;
  file_size: number;
  mime_type: string;
  version: number;
  created_at: string;
}