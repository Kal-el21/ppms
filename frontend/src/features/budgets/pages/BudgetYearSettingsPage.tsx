import { useState } from "react";
import { useBudgetYears, useCreateBudgetYear, useUpdateBudgetYear, useDeleteBudgetYear } from "../hooks/useBudgetYears";
import { Card, CardHeader, CardTitle, CardContent } from "../../../components/ui/card";
import { PageHeader } from "../../../components/shared/PageHeader";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../../components/ui/table";
import { BulkActionsDropdown } from "@/components/shared/BulkActionsDropdown";
import { ConfirmDeleteDialog } from "@/components/shared/ConfirmDeleteDialog";
import { EmptyState } from "@/components/shared/EmptyState";
import { TableSkeleton } from "@/components/ui/skeleton";
import { useTableSelection } from "@/components/shared/useTableSelection";

const formatCurrency = (value: number) =>
  `Rp ${Math.round(value || 0).toLocaleString("id-ID")}`;

export default function BudgetYearSettingsPage() {
  const { data, isLoading } = useBudgetYears();
  const createYear = useCreateBudgetYear();
  const updateYear = useUpdateBudgetYear();
  const deleteYear = useDeleteBudgetYear();

  const [year, setYear] = useState("");
  const [capex, setCapex] = useState("");
  const [opex, setOpex] = useState("");

  const [editId, setEditId] = useState<number | null>(null);
  const [editCapex, setEditCapex] = useState("");
  const [editOpex, setEditOpex] = useState("");
  const [editVersion, setEditVersion] = useState(0);

  const [deletingIds, setDeletingIds] = useState<number[]>([]);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);

  const rows = data ?? [];
  const {
    selectedIds,
    toggle,
    toggleAll,
    clear,
    selectedCount,
    isAllSelected,
    isIndeterminate,
  } = useTableSelection<number>();

  const handleCreate = () => {
    const y = Number(year);
    if (!y || y < 2000 || y > 2100) return;
    createYear.mutate({
      year: y,
      capex_ceiling: Number(capex) || 0,
      opex_ceiling: Number(opex) || 0,
    });
    setYear("");
    setCapex("");
    setOpex("");
  };

  const startEdit = (id: number, capexCeiling: number, opexCeiling: number, version: number) => {
    setEditId(id);
    setEditCapex(String(capexCeiling));
    setEditOpex(String(opexCeiling));
    setEditVersion(version);
  };

  const saveEdit = () => {
    if (editId == null) return;
    updateYear.mutate(
      {
        id: editId,
        payload: {
          capex_ceiling: Number(editCapex) || 0,
          opex_ceiling: Number(editOpex) || 0,
          version: editVersion,
        },
      },
      {
        onSettled: () => setEditId(null),
      }
    );
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
      await deleteYear.mutateAsync(id);
    }
    setDeletingIds([]);
    setDeleteDialogOpen(false);
    clear();
  };

  return (
    <div>
      <PageHeader
        title="Pagu Tahunan"
        subtitle="Kelola plafon anggaran CAPEX/OPEX per tahun"
        actions={
          <BulkActionsDropdown
            selectedCount={selectedCount}
            onEdit={
              selectedCount === 1
                ? () => {
                    const id = Array.from(selectedIds)[0];
                    const row = rows.find((r) => r.id === id);
                    if (row)
                      startEdit(row.id, row.capex_ceiling, row.opex_ceiling, row.version);
                  }
                : undefined
            }
            onDelete={handleBulkDelete}
          />
        }
      />

      <Card className="mb-4">
        <CardHeader>
          <CardTitle>Tambah Pagu Tahun</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 sm:grid-cols-4 gap-2.5 items-end">
            <div>
              <p className="text-[11.5px] text-ink-tertiary mb-1">Tahun</p>
              <Input type="number" placeholder="2026" value={year} onChange={(e) => setYear(e.target.value)} />
            </div>
            <div>
              <p className="text-[11.5px] text-ink-tertiary mb-1">CAPEX Ceiling</p>
              <Input type="number" placeholder="0" value={capex} onChange={(e) => setCapex(e.target.value)} />
            </div>
            <div>
              <p className="text-[11.5px] text-ink-tertiary mb-1">OPEX Ceiling</p>
              <Input type="number" placeholder="0" value={opex} onChange={(e) => setOpex(e.target.value)} />
            </div>
            <Button variant="primary" onClick={handleCreate} disabled={createYear.isPending}>
              Tambah
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Daftar Pagu Tahunan</CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <TableSkeleton rows={4} cols={5} />
          ) : rows.length === 0 ? (
            <EmptyState
              icon={<svg className="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8"><path d="M3 3v18h18M7 15l4-4 3 3 5-6" /></svg>}
              title="Belum ada pagu tahunan"
            />
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead style={{ width: 40 }}>
                    <input
                      type="checkbox"
                      checked={rows.length > 0 ? isAllSelected(rows.map((r) => r.id)) : false}
                      ref={(el) => {
                        if (el) el.indeterminate = isIndeterminate(rows.map((r) => r.id));
                      }}
                      onChange={(e) => toggleAll(rows.map((r) => r.id), e.target.checked)}
                    />
                  </TableHead>
                  <TableHead style={{ width: 50 }}>No</TableHead>
                  <TableHead>Tahun</TableHead>
                  <TableHead>CAPEX Ceiling</TableHead>
                  <TableHead>OPEX Ceiling</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {rows.map((row, index) => (
                  <TableRow key={row.id}>
                    <TableCell>
                      <input
                        type="checkbox"
                        checked={selectedIds.has(row.id)}
                        onChange={() => toggle(row.id)}
                      />
                    </TableCell>
                    <TableCell className="text-ink-secondary text-[13px]">{index + 1}</TableCell>
                    {editId === row.id ? (
                      <>
                        <TableCell className="font-medium">
                          <Input type="number" value={row.year} disabled className="w-20" />
                        </TableCell>
                        <TableCell>
                          <Input type="number" value={editCapex} onChange={(e) => setEditCapex(e.target.value)} className="w-40" />
                        </TableCell>
                        <TableCell>
                          <Input type="number" value={editOpex} onChange={(e) => setEditOpex(e.target.value)} className="w-40" />
                        </TableCell>
                        <TableCell>
                          <div className="flex gap-1.5">
                            <Button size="sm" variant="primary" onClick={saveEdit}>
                              Simpan
                            </Button>
                            <Button size="sm" variant="outline" onClick={() => setEditId(null)}>
                              Batal
                            </Button>
                          </div>
                        </TableCell>
                      </>
                    ) : (
                      <>
                        <TableCell className="font-medium">{row.year}</TableCell>
                        <TableCell>{formatCurrency(row.capex_ceiling)}</TableCell>
                        <TableCell>{formatCurrency(row.opex_ceiling)}</TableCell>
                      </>
                    )}
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      <ConfirmDeleteDialog
        open={deleteDialogOpen}
        title="Hapus pagu tahun"
        description={
          deletingIds.length === 1
            ? "Apakah Anda yakin ingin menghapus pagu tahun ini?"
            : `Apakah Anda yakin ingin menghapus ${deletingIds.length} pagu tahun yang dipilih?`
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
