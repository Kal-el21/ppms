import { useState } from "react";
import { useUsers, useDeactivateUser, useCreateUser, useUpdateUser, useAssignUserRole } from "../hooks/useUsers";
import { useDivisions } from "../../divisions/hooks/useDivisions";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select } from "@/components/ui/select";
import { Avatar } from "@/components/ui/avatar";
import { StatusBadge } from "@/components/ui/status-badge";
import { PageHeader } from "@/components/shared/PageHeader";
import { EmptyState } from "@/components/shared/EmptyState";
import { Pagination } from "@/components/ui/pagination";
import { TableSkeleton } from "@/components/ui/skeleton";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import {
  Table, TableBody, TableCell, TableHead, TableHeader, TableRow,
  TableDensityToggle, type TableDensity,
} from "@/components/ui/table";
import { BulkActionsDropdown } from "@/components/shared/BulkActionsDropdown";
import { ConfirmDeleteDialog } from "@/components/shared/ConfirmDeleteDialog";
import { useTableSelection } from "@/components/shared/useTableSelection";

const roleColor: Record<string, "blue" | "amber" | "gray"> = {
  ADMIN: "blue", USER: "amber", VIEWER: "gray",
};

interface UserFormData {
  full_name: string;
  email: string;
  password: string;
  system_role: string;
  division_id: string;
}

export default function UsersPage() {
  const [page, setPage] = useState(1);
  const limit = 20;
  const { data, isLoading, refetch } = useUsers(page, limit);
  const { data: divisions } = useDivisions();
  const { mutate: deactivate } = useDeactivateUser();
  const { mutate: assignRole } = useAssignUserRole();
  const { mutate: updateUser } = useUpdateUser();
  const { mutate: createUser, isPending: creating } = useCreateUser();

  const [showForm, setShowForm] = useState(false);
  const [editingUserId, setEditingUserId] = useState<number | null>(null);
  const [form, setForm] = useState<UserFormData>({
    full_name: "", email: "", password: "", system_role: "USER", division_id: "",
  });

  const [deletingIds, setDeletingIds] = useState<number[]>([]);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [density, setDensity] = useState<TableDensity>("comfortable");

  const divisionMap = Object.fromEntries((divisions || []).map((d) => [String(d.id), d.name]));

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

  const openCreateForm = () => {
    setEditingUserId(null);
    setForm({ full_name: "", email: "", password: "", system_role: "USER", division_id: "" });
    setShowForm(true);
  };

  const openEditForm = (user: { id: number; full_name: string; system_role: string; division_id: number | null }) => {
    setEditingUserId(user.id);
    setShowForm(true);
    setForm({
      full_name: user.full_name,
      email: "",
      password: "",
      system_role: user.system_role,
      division_id: user.division_id ? String(user.division_id) : "",
    });
  };

  const closeForm = () => {
    setShowForm(false);
    setEditingUserId(null);
    setForm({ full_name: "", email: "", password: "", system_role: "USER", division_id: "" });
  };

  const handleCreate = () => {
    createUser(
      {
        full_name: form.full_name,
        email: form.email,
        password: form.password,
        system_role: form.system_role,
        division_id: form.division_id ? Number(form.division_id) : null,
      },
      { onSuccess: closeForm }
    );
  };

  const handleUpdate = () => {
    if (!editingUserId) return;
    updateUser(
      { id: editingUserId, payload: { full_name: form.full_name, division_id: form.division_id ? Number(form.division_id) : null } },
      { onSuccess: closeForm }
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
      await deactivate(id);
    }
    setDeletingIds([]);
    setDeleteDialogOpen(false);
    clear();
    refetch();
  };

  const handleEdit = () => {
    if (selectedCount !== 1) return;
    const id = Array.from(selectedIds)[0];
    const user = rows.find((u) => u.id === id);
    if (user) {
      openEditForm(user);
    }
  };

  return (
    <div>
      <PageHeader
        title="User Management"
        subtitle={`${data?.meta.total ?? 0} user terdaftar`}
        actions={
          <>
            <TableDensityToggle value={density} onChange={setDensity} />
            <Button variant="primary" onClick={openCreateForm}>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M12 5v14M5 12h14" />
              </svg>
              Tambah user
            </Button>
          </>
        }
      />

      {showForm && (
        <Card className="mb-5 border-primary-200 dark:border-primary-900/50">
          <CardHeader><CardTitle>{editingUserId ? "Edit User" : "Tambah User Baru"}</CardTitle></CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 mb-4">
              <div>
                <Label>Nama lengkap</Label>
                <Input placeholder="cth. Budi Santoso" value={form.full_name} onChange={(e) => setForm({ ...form, full_name: e.target.value })} />
              </div>
              {!editingUserId && (
                <div>
                  <Label>Email</Label>
                  <Input type="email" placeholder="budi@perusahaan.com" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
                </div>
              )}
              {!editingUserId && (
                <div>
                  <Label>Password awal</Label>
                  <Input type="password" placeholder="Min. 8 karakter" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} />
                </div>
              )}
              <div>
                <Label>System Role</Label>
                <Select
                  value={form.system_role}
                  onChange={(e) => setForm({ ...form, system_role: e.target.value })}
                  options={[
                    { value: "ADMIN", label: "Admin" },
                    { value: "USER", label: "User" },
                    { value: "VIEWER", label: "Viewer" },
                  ]}
                />
              </div>
              <div>
                <Label>Division</Label>
                <Select
                  value={form.division_id}
                  onChange={(e) => setForm({ ...form, division_id: e.target.value })}
                  placeholder="Pilih divisi (opsional)"
                  options={(divisions || []).map((d) => ({ value: String(d.id), label: d.name }))}
                />
              </div>
            </div>
            <div className="flex gap-2">
              <Button variant="primary" onClick={editingUserId ? handleUpdate : handleCreate} disabled={creating}>
                {creating ? "Menyimpan..." : "Simpan"}
              </Button>
              <Button variant="outline" onClick={() => setShowForm(false)}>Batal</Button>
            </div>
          </CardContent>
        </Card>
      )}

      {isLoading ? (
        <TableSkeleton rows={8} cols={6} />
      ) : !data || data.data.length === 0 ? (
        <EmptyState
          icon={<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8"><path d="M16 21v-2a4 4 0 00-4-4H6a4 4 0 00-4 4v2M9 11a4 4 0 100-8 4 4 0 000 8z" /></svg>}
          title="Belum ada user"
        />
      ) : (
        <>
          <Table density={density}>
            <TableHeader>
              <TableRow>
                <TableHead style={{ width: 40 }}>
                  <input
                    type="checkbox"
                    checked={isAllSelected(rows.map((u) => u.id))}
                    ref={(el) => {
                      if (el) el.indeterminate = isIndeterminate(rows.map((u) => u.id));
                    }}
                    onChange={(e) => toggleAll(rows.map((u) => u.id), e.target.checked)}
                  />
                </TableHead>
                <TableHead style={{ width: 50 }}>No</TableHead>
                <TableHead>Name</TableHead>
                <TableHead>Email</TableHead>
                <TableHead>Role</TableHead>
                <TableHead>Division</TableHead>
                <TableHead>Status</TableHead>
                <TableHead style={{ width: 120 }}></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {rows.map((u, index) => (
                <TableRow key={u.id}>
                  <TableCell>
                    <input
                      type="checkbox"
                      checked={selectedIds.has(u.id)}
                      onChange={() => toggle(u.id)}
                    />
                  </TableCell>
                  <TableCell className="text-ink-secondary text-[13px]">{(page - 1) * limit + index + 1}</TableCell>
                  <TableCell>
                    <div className="flex items-center gap-2.5">
                      <Avatar name={u.full_name} size="sm" />
                      <span className="font-medium text-ink-primary">{u.full_name}</span>
                    </div>
                  </TableCell>
                  <TableCell className="text-ink-secondary">{u.email}</TableCell>
                  <TableCell>
                    <StatusBadge color={roleColor[u.system_role] || "gray"}>{u.system_role}</StatusBadge>
                  </TableCell>
                  <TableCell>
                    {u.division_id ? (
                      <span className="text-ink-secondary">{divisionMap[String(u.division_id)] || `Div #${u.division_id}`}</span>
                    ) : (
                      <span className="text-ink-tertiary">—</span>
                    )}
                  </TableCell>
                  <TableCell>
                    <StatusBadge color={u.is_active ? "green" : "red"}>{u.is_active ? "Active" : "Inactive"}</StatusBadge>
                  </TableCell>
                  <TableCell className="flex flex-wrap gap-2">
                    <Select
                      value={u.system_role}
                      onChange={(e) => assignRole({ id: u.id, systemRole: e.target.value })}
                      className="w-32"
                      options={[
                        { value: "ADMIN", label: "Admin" },
                        { value: "USER", label: "User" },
                        { value: "VIEWER", label: "Viewer" },
                      ]}
                    />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
          <div className="mt-4">
            <BulkActionsDropdown
              selectedCount={selectedCount}
              onEdit={handleEdit}
              onDelete={handleBulkDelete}
            />
          </div>
          <Pagination page={page} total={data.meta.total} limit={limit} onPageChange={setPage} />
        </>
      )}

      <ConfirmDeleteDialog
        open={deleteDialogOpen}
        title={`Nonaktifkan ${deletingIds.length} user`}
        description={
          deletingIds.length === 1
            ? "User yang dinonaktifkan tidak akan bisa login kembali."
            : `${deletingIds.length} user yang dipilih akan dinonaktifkan dan tidak bisa login kembali.`
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
