import { useMutation } from "@tanstack/react-query";
import { reportingApi } from "../api/reportingApi";
import { useToast } from "../../../components/ui/toast";

export function useGenerateReport() {
  const toast = useToast();
  return useMutation({
    mutationFn: (payload: { type: string; format: string }) => reportingApi.generate(payload.type, payload.format),
    onError: (error: any) => {
      toast.error("Gagal generate laporan", error?.friendlyMessage || error?.message);
    },
  });
}

export function useGenerateReportForProject() {
  const toast = useToast();
  return useMutation({
    mutationFn: (payload: { projectId: number; type: string; format: string }) =>
      reportingApi.generateForProject(payload.projectId, payload.type, payload.format),
    onError: (error: any) => {
      toast.error("Gagal generate laporan project", error?.friendlyMessage || error?.message);
    },
  });
}

export default useGenerateReport;
