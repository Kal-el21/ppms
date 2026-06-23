import { useState } from "react";
import { useDivisions, useCreateDivision } from "../hooks/useDivisions";
import { useAuth } from "../../auth/context/AuthContext";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { PageHeader } from "../../../components/shared/PageHeader";
import { EmptyState } from "../../../components/shared/EmptyState";
import { Card, CardContent } from "../../../components/ui/card";

export default function DivisionsPage() {
  const { user } = useAuth();
  const { data, isLoading } = useDivisions();
  const { mutate: createDivision, isPending } = useCreateDivision();

  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");

  const isAdmin = user?.system_role === "ADMIN";

  if (isLoading) return <div className="text-ink-secondary text-sm">Memuat divisions...</div>;

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
        subtitle={`${data?.length ?? 0} divisi terdaftar`}
        actions={
          isAdmin ? (
            <Button variant="primary" onClick={() => setShowForm((v) => !v)}>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M12 5v14M5 12h14" />
              </svg>
              Tambah divisi
            </Button>
          ) : undefined
        }
      />

      {showForm && (
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

      {!data || data.length === 0 ? (
        <EmptyState
          icon={
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
              <path d="M3 21h18M5 21V7l8-4v18M19 21V11l-6-4" />
            </svg>
          }
          title="Belum ada divisi"
        />
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3.5">
          {data.map((division) => (
            <Card key={division.id}>
              <CardContent className="pt-5">
                <div className="h-9 w-9 rounded-md bg-primary-50 dark:bg-primary-900/30 text-primary-700 dark:text-primary-400 flex items-center justify-center font-semibold text-[13px] mb-3">
                  {division.name.slice(0, 2).toUpperCase()}
                </div>
                <p className="text-[13.5px] font-semibold m-0 mb-1">{division.name}</p>
                <p className="text-[12px] text-ink-tertiary m-0 leading-relaxed">
                  {division.description || "Tidak ada deskripsi"}
                </p>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}