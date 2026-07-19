import { useState } from "react";
import { useDivisions, useCreateDivision, useUpdateDivision, useDeleteDivision } from "../hooks/useDivisions";
import { useAuth } from "../../auth/context/AuthContext";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { PageHeader } from "@/components/shared/PageHeader";
import { EmptyState } from "@/components/shared/EmptyState";
import { TableSkeleton } from "@/components/ui/skeleton";
import { Card, CardContent } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { BulkActionsDropdown } from "@/components/shared/BulkActionsDropdown";
import { ConfirmDeleteDialog } from "@/components/shared/ConfirmDeleteDialog";
import { useTableSelection } from "@/components/shared/useTableSelection";

export default function DivisionsPage() {
  const { user } = useAuth();
  const { data, isLoading } = useDivisions();
  const { mutate: createDivision, isPending } = useCreateDivision();
  const { mutateAsync: updateDivision } = useUpdateDivision();
  const { mutateAsync: deleteDivision } = useDeleteDivision();

  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editName, setEditName] = useState("");
  const [editDesc, setEditDesc] = useState("");
  const [deletingIds, setDeletingIds] = useState<number[]>([]);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  void showForm;
  void deleteDialogOpen;

  const isAdmin = user?.system_role === "ADMIN";

  const {
    selectedIds,
    toggle,
    toggleAll,
    clear,
    selectedCount,
    isAllSelected,
    isIndeterminate,
  } = useTableSelection<number>();

  if (isLoading) return <TableSkeleton rows={6} cols={4} />;

  const divisions = data ?? [];

  const handleEdit = (id: number, name: string, description: string) => {
    setEditingId(id);
    setEditName(name);
    setEditDesc(description);
  };

  const handleSaveEdit = async () => {
    if (!editingId || !editName.trim()) return;
    await updateDivision({ id: editingId, payload: { name: editName, description: editDesc } });
    setEditingId(null);
  };

  const handleBulkDelete = () => {
    const ids = Array.from(selectedIds);
    if (ids.length === 0) return;
    setDeletingIds(ids);
    setDeleteDialogOpen(true);
  };

  const confirmDelete = async () => {
    if (deletingIds.length === 0) return;
    for (const id of deletingIds) {
      await deleteDivision(id);
    }
    setDeletingIds([]);
    setDeleteDialogOpen(false);
    clear();
  };

  const handleCreate = () => {
    if (!name.trim()) return;
    createDivision(
      { name, description },
      {
        onSuccess: () => {
          setName("");
          setDescription("");
          setShowForm(false);
        },
      }
    );
  };

  return (
    <div>
      <PageHeader
        title="Divisions"
        subtitle={`${divisions.length} divisi terdaftar`}
        actions={
          isAdmin ? (
            <BulkActionsDropdown
              selectedCount={selectedCount}
              onEdit={
                selectedCount === 1
                  ? () => {
                      const id = Array.from(selectedIds)[0];
                      const division = divisions.find((d) => d.id === id);
                      if (division) handleEdit(division.id, division.name, division.description || "");
                    }
                  : undefined
              }
              onDelete={handleBulkDelete}
            />
          ) : undefined
        }
      />

      {showForm && isAdmin && (
        <Card className="mb-5">
          <CardContent className="pt-5">
            <div className="flex gap-2">
              <Input placeholder="Nama divisi" value={name} onChange={(e) => setName(e.target.value)} />
              <Input placeholder="Deskripsi" value={description} onChange={(e) => setDescription(e.target.value)} />
              <Button variant="primary" onClick={handleCreate} disabled={isPending}>
                Simpan
              </Button>
              <Button variant="outline" onClick={() => setShowForm(false)}>
                Batal
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {!divisions || divisions.length === 0 ? (
        <EmptyState
          icon={
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
              <path d="M3 21h18M5 21V7l8-4v18M19 21V11l-6-4" />
            </svg>
          }
          title="Belum ada divisi"
        />
      ) : (
        <Table>
          <TableHeader>
            <TableRow>
              {isAdmin && (
                <TableHead style={{ width: 40 }}>
                  <input
                    type="checkbox"
                    checked={isAllSelected(divisions.map((d) => d.id))}
                    ref={(el) => {
                      if (el) el.indeterminate = isIndeterminate(divisions.map((d) => d.id));
                    }}
                    onChange={(e) => toggleAll(divisions.map((d) => d.id), e.target.checked)}
                  />
                </TableHead>
              )}
              <TableHead style={{ width: 50 }}>No</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Description</TableHead>
              {isAdmin && editingId !== null && <TableHead style={{ width: 100 }}>Aksi</TableHead>}
            </TableRow>
          </TableHeader>
          <TableBody>
            {divisions.map((division, index) => (
              <TableRow key={division.id}>
                {isAdmin && (
                  <TableCell>
                    <input
                      type="checkbox"
                      checked={selectedIds.has(division.id)}
                      onChange={() => toggle(division.id)}
                    />
                  </TableCell>
                )}
                <TableCell className="text-ink-secondary text-[13px]">{index + 1}</TableCell>
                <TableCell>
                  {editingId === division.id ? (
                    <Input
                      value={editName}
                      onChange={(e) => setEditName(e.target.value)}
                      className="w-48"
                    />
                  ) : (
                    <span className="font-medium text-ink-primary">{division.name}</span>
                  )}
                </TableCell>
                <TableCell>
                  {editingId === division.id ? (
                    <Input
                      value={editDesc}
                      onChange={(e) => setEditDesc(e.target.value)}
                      className="w-64"
                    />
                  ) : (
                    <span className="text-ink-tertiary text-[13px]">
                      {division.description || "Tidak ada deskripsi"}
                    </span>
                  )}
                </TableCell>
                {isAdmin && editingId === division.id && (
                  <TableCell>
                    <div className="flex gap-1.5">
                      <Button size="sm" variant="primary" onClick={handleSaveEdit}>
                        Simpan
                      </Button>
                      <Button size="sm" variant="outline" onClick={() => setEditingId(null)}>
                        Batal
                      </Button>
                    </div>
                  </TableCell>
                )}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      )}

      <ConfirmDeleteDialog
        open={deleteDialogOpen}
        title="Hapus divisi"
        description={
          deletingIds.length === 1
            ? "Apakah Anda yakin ingin menghapus divisi ini?"
            : `Apakah Anda yakin ingin menghapus ${deletingIds.length} divisi yang dipilih?`
        }
        onConfirm={confirmDelete}
        onCancel={() => {
          setDeleteDialogOpen(false);
          setDeletingIds([]);
        }}
        isDeleting={deletingIds.length > 0}
      />
    </div>
  );
}
