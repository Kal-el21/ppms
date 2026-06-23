import { useNavigate } from "react-router-dom";
import { useRequestList } from "../hooks/useRequests";
import { Button } from "../../../components/ui/button";
import { StatusBadge, getStatusColor } from "../../../components/ui/status-badge";
import { PageHeader } from "../../../components/shared/PageHeader";
import { EmptyState } from "../../../components/shared/EmptyState";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../../components/ui/table";

export default function RequestListPage() {
  const navigate = useNavigate();
  const { data, isLoading } = useRequestList();

  if (isLoading) return <div className="text-ink-secondary text-sm">Memuat requests...</div>;

  return (
    <div>
      <PageHeader
        title="Project Requests"
        subtitle={`${data?.meta.total ?? 0} request tercatat`}
        actions={
          <Button variant="primary" onClick={() => navigate("/project-requests/new")}>
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M12 5v14M5 12h14" />
            </svg>
            Request baru
          </Button>
        }
      />

      {!data || data.data.length === 0 ? (
        <EmptyState
          icon={
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
              <rect x="3" y="4" width="18" height="18" rx="2" />
              <path d="M3 10h18M9 3v4" />
            </svg>
          }
          title="Belum ada project request"
          description="Mulai dengan membuat draft pengajuan project baru."
          action={
            <Button variant="primary" size="sm" onClick={() => navigate("/project-requests/new")}>
              Buat request
            </Button>
          }
        />
      ) : (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Title</TableHead>
              <TableHead>Estimated Budget</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Submitted</TableHead>
              <TableHead style={{ width: 80 }}></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.data.map((req) => (
              <TableRow key={req.id} className="cursor-pointer" onClick={() => navigate(`/project-requests/${req.id}`)}>
                <TableCell>
                  <div className="font-medium text-ink-primary">{req.title}</div>
                  <div className="text-[11.5px] text-ink-tertiary mt-0.5">REQ-{req.id}</div>
                </TableCell>
                <TableCell className="font-medium">Rp {req.estimated_budget.toLocaleString("id-ID")}</TableCell>
                <TableCell>
                  <StatusBadge color={getStatusColor(req.status)}>{req.status}</StatusBadge>
                </TableCell>
                <TableCell className="text-ink-secondary">
                  {req.submitted_at ? new Date(req.submitted_at).toLocaleDateString("id-ID") : "—"}
                </TableCell>
                <TableCell>
                  <Button variant="ghost" size="sm">
                    Lihat
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      )}
    </div>
  );
}