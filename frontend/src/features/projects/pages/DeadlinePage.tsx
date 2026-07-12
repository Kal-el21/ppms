import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useProjectDeadlines } from "../hooks/useProjects";
import { useDeleteProject } from "../hooks/useProjects";
import { type ProjectDeadline } from "../types";
import { Card, CardHeader, CardTitle, CardContent } from "../../../components/ui/card";
import { PageHeader } from "../../../components/shared/PageHeader";
import { StatusBadge, getStatusColor } from "../../../components/ui/status-badge";
import { HealthBar } from "../../../components/ui/health-bar";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "../../../components/ui/table";
import { Tabs } from "../../../components/ui/tabs";
import { BulkActionsDropdown } from "@/components/shared/BulkActionsDropdown";
import { ConfirmDeleteDialog } from "@/components/shared/ConfirmDeleteDialog";
import { useTableSelection } from "@/components/shared/useTableSelection";

const WINDOWS = [
  { key: "overdue", label: "Overdue" },
  { key: "30", label: "Due 30 Days" },
  { key: "60", label: "Due 60 Days" },
  { key: "90", label: "Due 90 Days" },
] as const;

const formatDate = (value?: string | null) =>
  value ? new Date(value).toLocaleDateString("id-ID", { day: "numeric", month: "short", year: "numeric" }) : "—";

function DeadlineTabContent({ list, onNavigate }: { list: ProjectDeadline[]; onNavigate: (id: number) => void }) {
  const deleteProject = useDeleteProject();
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
      await deleteProject.mutateAsync(id);
    }
    setDeletingIds([]);
    setDeleteDialogOpen(false);
    clear();
  };

  return (
    <>
      <div className="mb-4">
        <BulkActionsDropdown selectedCount={selectedCount} onDelete={handleBulkDelete} />
      </div>
      {list.length === 0 ? (
        <p className="text-[13px] text-ink-tertiary">Tidak ada project pada window ini.</p>
      ) : (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead style={{ width: 40 }}>
                <input
                  type="checkbox"
                  checked={isAllSelected(list.map((d) => d.id))}
                  ref={(el) => {
                    if (el) el.indeterminate = isIndeterminate(list.map((d) => d.id));
                  }}
                  onChange={(e) => toggleAll(list.map((d) => d.id), e.target.checked)}
                />
              </TableHead>
              <TableHead style={{ width: 50 }}>No</TableHead>
              <TableHead>Kode</TableHead>
              <TableHead>Project</TableHead>
              <TableHead>End Date</TableHead>
              <TableHead>Status</TableHead>
              <TableHead style={{ width: 160 }}>Progress</TableHead>
              <TableHead>Sisa Waktu</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {list.map((d, index) => {
              const atRisk = d.status === "ON_HOLD" || d.progress < 30;
              return (
                <TableRow
                  key={d.id}
                  className="cursor-pointer"
                  onClick={() => onNavigate(d.id)}
                >
                  <TableCell>
                    <input type="checkbox" checked={selectedIds.has(d.id)} onChange={() => toggle(d.id)} />
                  </TableCell>
                  <TableCell className="text-ink-secondary text-[13px]">{index + 1}</TableCell>
                  <TableCell className="text-[12.5px] text-ink-secondary">{d.project_code}</TableCell>
                  <TableCell className="font-medium text-ink-primary">{d.name}</TableCell>
                  <TableCell className="text-[12.5px] text-ink-secondary">{formatDate(d.end_date)}</TableCell>
                  <TableCell>
                    <StatusBadge color={getStatusColor(d.status)}>{d.status}</StatusBadge>
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <HealthBar progress={d.progress} atRisk={atRisk} className="flex-1" />
                      <span className="text-xs font-semibold text-ink-secondary w-9 text-right">
                        {d.progress.toFixed(0)}%
                      </span>
                    </div>
                  </TableCell>
                  <TableCell
                    className={
                      d.days_remaining < 0
                        ? "text-danger-600 font-semibold text-[12.5px]"
                        : d.days_remaining <= 30
                        ? "text-warning-600 font-semibold text-[12.5px]"
                        : "text-ink-secondary text-[12.5px]"
                    }
                  >
                    {d.days_remaining < 0 ? `${Math.abs(d.days_remaining)} hari lalu` : `${d.days_remaining} hari`}
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      )}
      <ConfirmDeleteDialog
        open={deleteDialogOpen}
        title="Hapus project"
        description={
          deletingIds.length === 1
            ? "Apakah Anda yakin ingin menghapus project ini?"
            : `Apakah Anda yakin ingin menghapus ${deletingIds.length} project yang dipilih?`
        }
        onConfirm={confirmDelete}
        onCancel={() => {
          setDeleteDialogOpen(false);
          setDeletingIds([]);
        }}
        isDeleting={deletingIds.length > 0}
      />
    </>
  );
}

export default function DeadlinePage() {
  const navigate = useNavigate();
  const overdue = useProjectDeadlines("overdue");
  const d30 = useProjectDeadlines("30");
  const d60 = useProjectDeadlines("60");
  const d90 = useProjectDeadlines("90");

  const map: Record<string, { data?: ProjectDeadline[]; isLoading: boolean }> = {
    overdue: { data: overdue.data, isLoading: overdue.isLoading },
    "30": { data: d30.data, isLoading: d30.isLoading },
    "60": { data: d60.data, isLoading: d60.isLoading },
    "90": { data: d90.data, isLoading: d90.isLoading },
  };

  return (
    <div>
      <PageHeader title="Deadline Projects" subtitle="Project berdasarkan tenggat waktu" />

      <Tabs
        tabs={WINDOWS.map((w) => ({ key: w.key, label: w.label }))}
        defaultTab="overdue"
      >
        {(activeTab) => {
          const current = map[activeTab];
          const list = current?.data ?? [];
          return (
            <Card className="mt-4">
              <CardHeader>
                <CardTitle>
                  {WINDOWS.find((w) => w.key === activeTab)?.label} ({list.length})
                </CardTitle>
              </CardHeader>
              <CardContent>
                {current?.isLoading ? (
                  <div className="h-24 bg-surface-tertiary rounded-md animate-pulse" />
                ) : (
                  <DeadlineTabContent list={list as ProjectDeadline[]} onNavigate={(id) => navigate(`/projects/${id}`)} />
                )}
              </CardContent>
            </Card>
          );
        }}
      </Tabs>
    </div>
  );
}
