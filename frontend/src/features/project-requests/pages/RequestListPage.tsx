import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useRequestList } from "../hooks/useRequests";
import { useDeleteDraft } from "../hooks/useRequests";
import { Button } from "@/components/ui/button";
import { StatusBadge, getStatusColor } from "@/components/ui/status-badge";
import { PageHeader } from "@/components/shared/PageHeader";
import { EmptyState } from "@/components/shared/EmptyState";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { TableSkeleton } from "@/components/ui/skeleton";
import { BulkActionsDropdown } from "@/components/shared/BulkActionsDropdown";
import { ConfirmDeleteDialog } from "@/components/shared/ConfirmDeleteDialog";
import { useTableSelection } from "@/components/shared/useTableSelection";

export default function RequestListPage() {
  const navigate = useNavigate();
  const { data, isLoading, refetch } = useRequestList();
  const { mutateAsync: deleteDraft } = useDeleteDraft();

  const rows = data?.data ?? [];
  const {
    selectedIds,
    toggle,
    toggleAll,
    clear,
    selectedCount,
    isAllSelected,
    isIndeterminate,
  } = useTableSelection<number>();

  const [deletingIds, setDeletingIds] = useState<number[]>([]);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  const handleBulkDelete = () => {
    const ids = Array.from(selectedIds);
    if (ids.length === 0) return;
    setDeletingIds(ids);
    setDeleteDialogOpen(true);
  };

  const confirmDelete = async () => {
    if (deletingIds.length === 0) return;
    for (const id of deletingIds) {
      await deleteDraft(id);
    }
    setDeletingIds([]);
    setDeleteDialogOpen(false);
    clear();
    refetch();
  };

  const handleRowClick = (reqId: number) => {
    navigate(`/project-requests/${reqId}`);
  };

  if (isLoading) return <TableSkeleton rows={5} cols={4} />;

  return (
    <div>
      <PageHeader
        title="Project Requests"
        subtitle={`${data?.meta.total ?? 0} request tercatat`}
        actions={
          <div className="flex gap-2">
            <BulkActionsDropdown
              selectedCount={selectedCount}
              onDelete={handleBulkDelete}
            />
            <Button variant="primary" onClick={() => navigate("/project-requests/new")}>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M12 5v14M5 12h14" />
              </svg>
              Request baru
            </Button>
          </div>
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
        <>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead style={{ width: 40 }}>
                  <input
                    type="checkbox"
                    checked={isAllSelected(rows.map((r) => r.id))}
                    ref={(el) => {
                      if (el) el.indeterminate = isIndeterminate(rows.map((r) => r.id));
                    }}
                    onChange={(e) => toggleAll(rows.map((r) => r.id), e.target.checked)}
                  />
                </TableHead>
                <TableHead style={{ width: 50 }}>No</TableHead>
                <TableHead>Title</TableHead>
                <TableHead>Estimated Budget</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Submitted</TableHead>
                <TableHead style={{ width: 80 }}></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {data.data.map((req, index) => (
                <TableRow
                  key={req.id}
                  className="cursor-pointer"
                  onClick={() => handleRowClick(req.id)}
                >
                  <TableCell onClick={(e) => e.stopPropagation()}>
                    <input
                      type="checkbox"
                      checked={selectedIds.has(req.id)}
                      onChange={() => toggle(req.id)}
                    />
                  </TableCell>
                  <TableCell className="text-ink-secondary text-[13px]" onClick={(e) => e.stopPropagation()}>
                    {(data.meta.page - 1) * (data.meta.limit || 20) + index + 1}
                  </TableCell>
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
                  <TableCell onClick={(e) => e.stopPropagation()}>
                    <Button variant="ghost" size="sm">
                      Lihat
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
          <ConfirmDeleteDialog
            open={deleteDialogOpen}
            title="Hapus request"
            description={
              deletingIds.length === 1
                ? "Apakah Anda yakin ingin menghapus request ini?"
                : `Apakah Anda yakin ingin menghapus ${deletingIds.length} request yang dipilih?`
            }
            onConfirm={confirmDelete}
            onCancel={() => {
              setDeleteDialogOpen(false);
              setDeletingIds([]);
            }}
            isDeleting={deletingIds.length > 0}
          />
        </>
      )}
    </div>
  );
}
