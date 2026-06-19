export type HandoverStatus = "PENDING" | "RECEIVED" | "RETURNED";

export interface Handover {
  id: number;
  project_id: number;
  sender_id: number;
  sender_division_id: number | null;
  receiver_id: number | null;
  description: string;
  delivery_date: string | null;
  delivery_time: string | null;
  received_at: string | null;
  status: HandoverStatus;
  version: number;
  created_at: string;
}