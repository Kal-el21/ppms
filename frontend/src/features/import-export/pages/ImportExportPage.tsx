import { useRef, useState } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "../../../components/ui/card";
import { PageHeader } from "../../../components/shared/PageHeader";
import { Button } from "../../../components/ui/button";
import { MetricCard } from "../../../components/ui/metric-card";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "../../../components/ui/dialog";
import { useExport, useImport } from "../hooks/useImportExport";
import type { ImportResult } from "../api/importExportApi";
import { useToast } from "../../../components/ui/toast";

interface Preview {
  fileName: string;
  size: number;
  version?: string;
  projectCount: number;
}

const formatSize = (bytes: number) => {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
};

export default function ImportExportPage() {
  const toast = useToast();
  const exportMutation = useExport();
  const importMutation = useImport();

  const inputRef = useRef<HTMLInputElement>(null);
  const [file, setFile] = useState<File | null>(null);
  const [preview, setPreview] = useState<Preview | null>(null);
  const [dragOver, setDragOver] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [result, setResult] = useState<ImportResult | null>(null);

  const handleFile = async (picked: File) => {
    if (!picked.name.toLowerCase().endsWith(".json")) {
      toast.error("Format tidak didukung", "Hanya file .json yang diperbolehkan");
      return;
    }
    setResult(null);
    try {
      const text = await picked.text();
      const parsed = JSON.parse(text);
      const projects = Array.isArray(parsed?.projects) ? parsed.projects : null;
      if (!projects) {
        toast.error("File tidak valid", "Struktur JSON tidak memiliki daftar 'projects'");
        return;
      }
      setFile(picked);
      setPreview({
        fileName: picked.name,
        size: picked.size,
        version: parsed?.version,
        projectCount: projects.length,
      });
    } catch {
      toast.error("File tidak valid", "Gagal membaca JSON. Pastikan file tidak rusak.");
    }
  };

  const onDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setDragOver(false);
    const dropped = e.dataTransfer.files?.[0];
    if (dropped) handleFile(dropped);
  };

  const resetImport = () => {
    setFile(null);
    setPreview(null);
    if (inputRef.current) inputRef.current.value = "";
  };

  const confirmImport = () => {
    if (!file) return;
    importMutation.mutate(file, {
      onSuccess: (res) => {
        setResult(res);
        setConfirmOpen(false);
        resetImport();
      },
      onError: () => setConfirmOpen(false),
    });
  };

  return (
    <div className="max-w-4xl">
      <PageHeader
        title="Import / Export Data"
        subtitle="Backup seluruh data project ke file JSON atau restore dari file backup."
      />

      {/* EXPORT */}
      <Card className="mb-5">
        <CardHeader>
          <CardTitle>Export Backup</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-[13px] text-ink-secondary mb-4">
            Unduh seluruh project yang dapat Anda akses beserta members, milestones, tasks,
            budget, dan transaksi dalam satu file JSON.
          </p>
          <Button
            variant="primary"
            onClick={() => exportMutation.mutate()}
            disabled={exportMutation.isPending}
          >
            {exportMutation.isPending ? "Menyiapkan..." : "Download Backup JSON"}
          </Button>
        </CardContent>
      </Card>

      {/* IMPORT */}
      <Card className="mb-5">
        <CardHeader>
          <CardTitle>Import Backup</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-[13px] text-ink-secondary mb-4">
            Upload file backup JSON untuk membuat ulang project. Setiap project akan dibuat
            sebagai project baru dengan kode project baru — data yang ada tidak akan ditimpa.
          </p>

          <div
            onDragOver={(e) => {
              e.preventDefault();
              setDragOver(true);
            }}
            onDragLeave={() => setDragOver(false)}
            onDrop={onDrop}
            onClick={() => inputRef.current?.click()}
            className={`flex flex-col items-center justify-center gap-2 rounded-lg border-2 border-dashed px-6 py-10 cursor-pointer transition-colors ${
              dragOver
                ? "border-primary-500 bg-primary-50 dark:bg-primary-900/20"
                : "border-border hover:border-primary-400 hover:bg-surface-secondary"
            }`}
          >
            <svg className="h-8 w-8 text-ink-tertiary" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
              <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M17 8l-5-5-5 5M12 3v12" />
            </svg>
            <p className="text-[13px] font-medium text-ink-primary">
              Tarik file ke sini atau klik untuk memilih
            </p>
            <p className="text-[11.5px] text-ink-tertiary">Hanya file .json (maks. 20MB)</p>
            <input
              ref={inputRef}
              type="file"
              accept="application/json,.json"
              className="hidden"
              onChange={(e) => {
                const picked = e.target.files?.[0];
                if (picked) handleFile(picked);
              }}
            />
          </div>

          {preview && (
            <div className="mt-4 flex items-center justify-between rounded-lg border border-border bg-surface-secondary px-4 py-3">
              <div>
                <p className="text-[13px] font-medium text-ink-primary">{preview.fileName}</p>
                <p className="text-[11.5px] text-ink-tertiary">
                  {formatSize(preview.size)} · {preview.projectCount} project
                  {preview.version ? ` · versi ${preview.version}` : ""}
                </p>
              </div>
              <div className="flex gap-2">
                <Button variant="outline" size="sm" onClick={resetImport}>
                  Batal
                </Button>
                <Button
                  variant="primary"
                  size="sm"
                  onClick={() => setConfirmOpen(true)}
                  disabled={importMutation.isPending || preview.projectCount === 0}
                >
                  Import
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* RESULT */}
      {result && (
        <Card>
          <CardHeader>
            <CardTitle>Hasil Import</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 mb-4">
              <MetricCard
                label="Total Project"
                value={result.total_projects}
                iconColor="blue"
                icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M3 7l4-4h6l4 4h4v13H3z" /></svg>}
              />
              <MetricCard
                label="Berhasil Diimpor"
                value={result.imported}
                iconColor="green"
                icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M20 6L9 17l-5-5" /></svg>}
              />
              <MetricCard
                label="Dilewati"
                value={result.skipped}
                iconColor="amber"
                icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M12 9v4M12 16h.01M10.3 3.6L1.6 18a2 2 0 001.7 3h17.4a2 2 0 001.7-3L13.7 3.6a2 2 0 00-3.4 0z" /></svg>}
              />
            </div>

            {result.errors.length > 0 && (
              <div className="rounded-lg border border-danger-200 bg-danger-50 dark:bg-danger-900/20 dark:border-danger-900/40 px-4 py-3">
                <p className="text-[12.5px] font-semibold text-danger-700 dark:text-danger-400 mb-2">
                  {result.errors.length} project gagal diimpor
                </p>
                <ul className="list-disc list-inside space-y-1">
                  {result.errors.map((err, idx) => (
                    <li key={idx} className="text-[12px] text-danger-700 dark:text-danger-400">
                      {err}
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* CONFIRM DIALOG */}
      <Dialog open={confirmOpen} onOpenChange={setConfirmOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Konfirmasi Import</DialogTitle>
            <DialogDescription>
              {preview
                ? `${preview.projectCount} project akan dibuat sebagai project baru dari file "${preview.fileName}". Lanjutkan?`
                : "Lanjutkan import?"}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setConfirmOpen(false)} disabled={importMutation.isPending}>
              Batal
            </Button>
            <Button variant="primary" onClick={confirmImport} disabled={importMutation.isPending}>
              {importMutation.isPending ? "Mengimpor..." : "Ya, Import"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
