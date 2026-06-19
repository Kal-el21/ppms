export interface Notification {
  id: number;
  type: string;
  title: string;
  message: string;
  entity_type: string;
  entity_id: number | null;
  action_url: string;
  is_read: boolean;
  created_at: string;
}

export interface NotificationPreference {
  type: string;
  enabled: boolean;
}