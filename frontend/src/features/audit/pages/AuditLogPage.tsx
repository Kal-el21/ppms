import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { auditApi } from "../api/auditApi";
import { StatusBadge } from "../../../components/ui/status-badge";
import { Input } from "../../../components/ui/input";
import { PageHeader } from "../../../components/shared/PageHeader";
import { EmptyState } from "../../../components/shared/EmptyState";
import { Pagination } from "../../../components/ui/pagination";
import { TableSkeleton } from "../../../components/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../../components/ui/table";

const moduleColor: Record<string, "blue" | "red" | "green" | "amber" | "gray"> = {
  auth: "blue",
  user: "blue",
  division: "gray",
  project_request: "amber",
  project: "blue",
  milestone: "blue",
  task: "blue",
  budget: "red",
  attachment: "gray",
  handover: "green",
};

function actionVariant(action: string): "blue" | "red" | "green" | "amber" | "gray" {
  if (action.includes("DELETE") || action.includes("REJECT") || action.includes("FAILED")) return "red";
  if (action.includes("APPROVE") || action.includes("CREATE") || action.includes("SUCCESS")) return "green";
  if (action.includes("REVOKE") || action.includes("WARNING")) return "amber";
  return "blue";
}

const LIMIT = 20;

export default function AuditLogPage() {
  const [module, setModule] = useState("");
  const [page, setPage] = useState(1);

  const { data, isLoading } = useQuery({
    queryKey: ["audit-logs", module, page],
    queryFn: () => auditApi.getList(page, LIMIT, module || undefined),
  });

  return (
    <div>
      <PageHeader title="Audit Logs" subtitle="Riwayat aktivitas penting di seluruh sistem" />

      <div className="mb-5 max-w-xs">
        <Input
          placeholder="Filter module (cth. auth, budget)"
          value={module}
          onChange={(e) => {
            setModule(e.target.value);
            setPage(1);
          }}
        />
      </div>

      {isLoading ? (
        <TableSkeleton rows={8} cols={5} />
      ) : !data || data.data.length === 0 ? (
        <EmptyState
          icon={
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
              <circle cx="12" cy="12" r="10" />
              <path d="M12 6v6l4 2" />
            </svg>
          }
          title="Tidak ada log ditemukan"
          description={module ? `Tidak ada aktivitas untuk module "${module}".` : undefined}
        />
      ) : (
        <>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Module</TableHead>
                <TableHead>Action</TableHead>
                <TableHead>Entity</TableHead>
                <TableHead>User ID</TableHead>
                <TableHead>Timestamp</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {data.data.map((log) => (
                <TableRow key={log.id}>
                  <TableCell>
                    <StatusBadge color={moduleColor[log.module] || "gray"}>{log.module}</StatusBadge>
                  </TableCell>
                  <TableCell>
                    <StatusBadge color={actionVariant(log.action)}>{log.action}</StatusBadge>
                  </TableCell>
                  <TableCell className="text-ink-secondary">
                    {log.entity_type} {log.entity_id ? `#${log.entity_id}` : ""}
                  </TableCell>
                  <TableCell className="text-ink-secondary">{log.user_id ?? "—"}</TableCell>
                  <TableCell className="text-ink-tertiary">
                    {new Date(log.created_at).toLocaleString("id-ID")}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
          <Pagination
            page={page}
            total={data.meta.total}
            limit={LIMIT}
            onPageChange={setPage}
          />
        </>
      )}
    </div>
  );
}