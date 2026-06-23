import { useState } from "react";
import { useUsers, useDeactivateUser } from "../hooks/useUsers";
import { Button } from "../../../components/ui/button";
import { Avatar } from "../../../components/ui/avatar";
import { StatusBadge } from "../../../components/ui/status-badge";
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

const roleColor: Record<string, "blue" | "amber" | "gray"> = {
  ADMIN: "blue",
  USER: "amber",
  VIEWER: "gray",
};

export default function UsersPage() {
  const [page] = useState(1);
  const { data, isLoading } = useUsers(page);
  const { mutate: deactivate } = useDeactivateUser();

  if (isLoading) return <div className="text-ink-secondary text-sm">Memuat users...</div>;

  return (
    <div>
      <PageHeader
        title="User Management"
        subtitle={`${data?.meta.total ?? 0} user terdaftar`}
        actions={
          <Button variant="primary">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M12 5v14M5 12h14" />
            </svg>
            Tambah user
          </Button>
        }
      />

      {!data || data.data.length === 0 ? (
        <EmptyState
          icon={
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
              <path d="M16 21v-2a4 4 0 00-4-4H6a4 4 0 00-4 4v2M9 11a4 4 0 100-8 4 4 0 000 8z" />
            </svg>
          }
          title="Belum ada user"
        />
      ) : (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Role</TableHead>
              <TableHead>Status</TableHead>
              <TableHead style={{ width: 100 }}></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data.data.map((u) => (
              <TableRow key={u.id}>
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
                  <StatusBadge color={u.is_active ? "green" : "red"}>
                    {u.is_active ? "Active" : "Inactive"}
                  </StatusBadge>
                </TableCell>
                <TableCell>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => deactivate(u.id)}
                    disabled={!u.is_active}
                  >
                    Nonaktifkan
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