import { useState } from "react";
import { useGenerateReport, useGenerateReportForProject } from "../hooks/useReporting";
import { useAuth } from "../../auth/context/AuthContext";
import { Card, CardHeader, CardTitle, CardContent } from "../../../components/ui/card";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import { PageHeader } from "../../../components/shared/PageHeader";

const reportTypes = [
  { value: "PROJECT", label: "Project Report", desc: "Status & progress seluruh project" },
  { value: "MILESTONE", label: "Milestone Report", desc: "Pencapaian milestone per project" },
  { value: "TASK", label: "Task Report", desc: "Detail task & prioritas" },
  { value: "BUDGET", label: "Budget Report", desc: "Alokasi & penggunaan anggaran" },
  { value: "HANDOVER", label: "Handover Report", desc: "Riwayat serah terima dokumen" },
];

const formatIcon = (format: "PDF" | "EXCEL") =>
  format === "PDF" ? (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
      <path d="M14 2v6h6" />
    </svg>
  ) : (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <rect x="3" y="3" width="18" height="18" rx="2" />
      <path d="M3 9h18M9 21V9" />
    </svg>
  );

export default function ReportingPage() {
  const { user } = useAuth();
  const isAdmin = user?.system_role === "ADMIN";

  const [type, setType] = useState("PROJECT");
  const [format, setFormat] = useState<"PDF" | "EXCEL">("PDF");
  const [projectId, setProjectId] = useState("");
  const [isGenerating, setIsGenerating] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const downloadBlob = (blob: Blob, fileName: string) => {
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = fileName;
    link.click();
    window.URL.revokeObjectURL(url);
  };

  const { mutateAsync: generate } = useGenerateReport();
  const { mutateAsync: generateForProject } = useGenerateReportForProject();

  const handleGenerate = async () => {
    setErrorMsg(null);
    setIsGenerating(true);
    try {
      const extension = format === "PDF" ? "pdf" : "xlsx";
      const fileName = `${type.toLowerCase()}_report.${extension}`;

      let blob: Blob;
      if (projectId) {
        blob = await generateForProject({ projectId: Number(projectId), type, format });
      } else {
        if (!isAdmin) {
          setErrorMsg("Hanya ADMIN yang bisa generate laporan sistem-wide. Isi Project ID di mana Anda menjadi PM.");
          setIsGenerating(false);
          return;
        }
        blob = await generate({ type, format });
      }

      downloadBlob(blob, fileName);
    } catch (err: any) {
      setErrorMsg(err?.response?.data?.message || "Gagal generate laporan");
    } finally {
      setIsGenerating(false);
    }
  };

  return (
    <div>
      <PageHeader title="Reporting" subtitle="Generate laporan project dalam format PDF atau Excel" />

      <Card>
        <CardHeader>
          <CardTitle>Pilih jenis laporan</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5 mb-6">
            {reportTypes.map((rt) => (
              <button
                key={rt.value}
                onClick={() => setType(rt.value)}
                className={`text-left rounded-lg border p-3.5 transition-colors ${
                  type === rt.value
                    ? "border-primary-600 bg-primary-50 dark:bg-primary-900/20"
                    : "border-border hover:bg-surface-secondary"
                }`}
              >
                <p
                  className={`text-[13px] font-semibold m-0 mb-0.5 ${
                    type === rt.value ? "text-primary-700 dark:text-primary-400" : "text-ink-primary"
                  }`}
                >
                  {rt.label}
                </p>
                <p className="text-[11.5px] text-ink-tertiary m-0">{rt.desc}</p>
              </button>
            ))}
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-5">
            <div>
              <Label htmlFor="projectId">
                Project ID {isAdmin ? "(opsional)" : "(wajib — Anda harus PM project ini)"}
              </Label>
              <Input
                id="projectId"
                type="number"
                placeholder="cth. 1"
                value={projectId}
                onChange={(e) => setProjectId(e.target.value)}
              />
              <p className="text-[11.5px] text-ink-tertiary mt-1.5">
                {isAdmin
                  ? "Kosongkan untuk laporan gabungan seluruh project."
                  : "Laporan project-scoped membutuhkan role Project Manager."}
              </p>
            </div>

            <div>
              <Label>Format file</Label>
              <div className="flex gap-2">
                {(["PDF", "EXCEL"] as const).map((f) => (
                  <button
                    key={f}
                    onClick={() => setFormat(f)}
                    className={`flex-1 flex items-center justify-center gap-2 h-9 rounded-md border text-[13px] font-medium transition-colors ${
                      format === f
                        ? "border-primary-600 bg-primary-50 text-primary-700 dark:bg-primary-900/20 dark:text-primary-400"
                        : "border-border-strong text-ink-secondary hover:bg-surface-secondary"
                    }`}
                  >
                    {formatIcon(f)}
                    {f}
                  </button>
                ))}
              </div>
            </div>
          </div>

          {errorMsg && (
            <div className="flex items-start gap-2 rounded-md bg-danger-50 dark:bg-danger-900/20 border border-danger-200 dark:border-danger-900/40 px-3 py-2.5 mb-4">
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-danger-600 flex-shrink-0 mt-0.5">
                <circle cx="12" cy="12" r="9" />
                <path d="M12 8v5M12 16h.01" />
              </svg>
              <p className="text-[12.5px] text-danger-700 dark:text-danger-400 m-0">{errorMsg}</p>
            </div>
          )}

          <Button variant="primary" onClick={handleGenerate} disabled={isGenerating} className="w-full sm:w-auto">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" />
              <path d="M7 10l5 5 5-5M12 15V3" />
            </svg>
            {isGenerating ? "Membuat laporan..." : "Generate & Download"}
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}