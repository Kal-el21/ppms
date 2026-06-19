import axiosInstance from "../../../api/axiosInstance";
import type { PaginatedResponse } from "../../../types";

export interface AuditLog {
  id: number;
  user_id: number | null;
  module: string;
  action: string;
  entity_type: string;
  entity_id: number | null;
  created_at: string;
}

export const auditApi = {
  getList: async (page = 1, limit = 20, module?: string) => {
    const res = await axiosInstance.get<PaginatedResponse<AuditLog>>("/audit-logs", {
      params: { page, limit, module },
    });
    return res.data;
  },
};