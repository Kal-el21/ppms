import { useState } from "react";
import { reportingApi } from "../api/reportingApi";
import { useAuth } from "../../auth/context/AuthContext";
import { Card, CardContent, CardHeader, CardTitle } from "../../../components/ui/card";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";

export default function ReportingPage() {
  const { user } = useAuth();
  const isAdmin = user?.system_role === "ADMIN";

  const [type, setType] = useState("PROJECT");
  const [format, setFormat] = useState("PDF");
  const [projectId, setProjectId] = useState("");
  const [isGenerating, setIsGenerating] = useState(false);

  const downloadBlob = (blob: Blob, fileName: string) => {
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = fileName;
    link.click();
    window.URL.revokeObjectURL(url);
  };

  const handleGenerate = async () => {
    setIsGenerating(true);
    try {
      const extension = format === "PDF" ? "pdf" : "xlsx";
      const fileName = `${type.toLowerCase()}_report.${extension}`;

      let blob: Blob;
      if (projectId) {
        // Project-scoped: butuh role PROJECT_MANAGER di project tersebut (atau ADMIN)
        blob = await reportingApi.generateForProject(Number(projectId), type, format);
      } else {
        // System-wide: ADMIN only
        if (!isAdmin) {
          alert("Only ADMIN can generate system-wide reports. Please specify a Project ID where you are PROJECT_MANAGER.");
          setIsGenerating(false);
          return;
        }
        blob = await reportingApi.generate(type, format);
      }

      downloadBlob(blob, fileName);
    } catch (err: any) {
      alert(err?.response?.data?.message || "Failed to generate report");
    } finally {
      setIsGenerating(false);
    }
  };

  return (
    <Card className="max-w-md">
      <CardHeader>
        <CardTitle>Generate Report</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div>
          <label className="text-sm font-medium">Report Type</label>
          <select value={type} onChange={(e) => setType(e.target.value)} className="w-full rounded border px-2 py-2 text-sm">
            <option value="PROJECT">Project Report</option>
            <option value="MILESTONE">Milestone Report</option>
            <option value="TASK">Task Report</option>
            <option value="BUDGET">Budget Report</option>
            <option value="HANDOVER">Handover Report</option>
          </select>
        </div>

        <div>
          <label className="text-sm font-medium">
            Project ID {isAdmin ? "(optional — leave empty for system-wide report)" : "(required — you must be PM of this project)"}
          </label>
          <Input
            type="number"
            placeholder="e.g. 1"
            value={projectId}
            onChange={(e) => setProjectId(e.target.value)}
          />
        </div>

        <div>
          <label className="text-sm font-medium">Format</label>
          <select value={format} onChange={(e) => setFormat(e.target.value)} className="w-full rounded border px-2 py-2 text-sm">
            <option value="PDF">PDF</option>
            <option value="EXCEL">Excel</option>
          </select>
        </div>

        <Button onClick={handleGenerate} disabled={isGenerating} className="w-full">
          {isGenerating ? "Generating..." : "Generate & Download"}
        </Button>
      </CardContent>
    </Card>
  );
}