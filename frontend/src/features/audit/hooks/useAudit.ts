import { useQuery } from "@tanstack/react-query";
import { auditApi } from "../api/auditApi";

export function useAuditList(page = 1, limit = 20, module?: string) {
  return useQuery({
    queryKey: ["audit-logs", module, page],
    queryFn: () => auditApi.getList(page, limit, module || undefined),
  });
}

export default useAuditList;
