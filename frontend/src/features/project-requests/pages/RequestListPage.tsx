import { useNavigate } from "react-router-dom";
import { useRequestList } from "../hooks/useRequests";
import { Button } from "../../../components/ui/button";
import { Badge } from "../../../components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../../components/ui/table";
import type { RequestStatus } from "../types";

const statusVariant: Record<RequestStatus, "default" | "destructive" | "outline" | "secondary"> = {
  DRAFT: "outline",
  SUBMITTED: "secondary",
  UNDER_REVIEW: "secondary",
  APPROVED: "default",
  REJECTED: "destructive",
  REVISED: "outline",
};

export default function RequestListPage() {
  const navigate = useNavigate();
  const { data, isLoading } = useRequestList();

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Project Requests</h1>
        <Button onClick={() => navigate("/project-requests/new")}>+ New Request</Button>
      </div>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Title</TableHead>
            <TableHead>Estimated Budget</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Submitted At</TableHead>
            <TableHead className="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {data?.data.map((req) => (
            <TableRow key={req.id}>
              <TableCell>{req.title}</TableCell>
              <TableCell>Rp {req.estimated_budget.toLocaleString("id-ID")}</TableCell>
              <TableCell>
                <Badge variant={statusVariant[req.status]}>{req.status}</Badge>
              </TableCell>
              <TableCell>{req.submitted_at ? new Date(req.submitted_at).toLocaleDateString() : "-"}</TableCell>
              <TableCell className="text-right">
                <Button variant="ghost" size="sm" onClick={() => navigate(`/project-requests/${req.id}`)}>
                  View
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}