import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { auditApi } from "../api/auditApi";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../../components/ui/table";
import { Badge } from "../../../components/ui/badge";

export default function AuditLogPage() {
  const [module, setModule] = useState("");
  const { data, isLoading } = useQuery({
    queryKey: ["audit-logs", module],
    queryFn: () => auditApi.getList(1, 50, module || undefined),
  });

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-semibold">Audit Logs</h1>

      <input
        placeholder="Filter by module (e.g. auth, project_request)"
        value={module}
        onChange={(e) => setModule(e.target.value)}
        className="w-64 rounded border px-3 py-2 text-sm"
      />

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
          {data?.data.map((log) => (
            <TableRow key={log.id}>
              <TableCell>
                <Badge variant="outline">{log.module}</Badge>
              </TableCell>
              <TableCell>{log.action}</TableCell>
              <TableCell>
                {log.entity_type} #{log.entity_id}
              </TableCell>
              <TableCell>{log.user_id ?? "-"}</TableCell>
              <TableCell>{new Date(log.created_at).toLocaleString()}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}