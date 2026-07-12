import { useMutation, useQueryClient } from "@tanstack/react-query";
import { importExportApi, type ImportResult } from "../api/importExportApi";
import { useToast } from "../../../components/ui/toast";

// useExport men-download file backup JSON dan memicu unduhan otomatis di browser.
export function useExport() {
  const toast = useToast();
  return useMutation({
    mutationFn: () => importExportApi.exportData(),
    onSuccess: (blob) => {
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      const stamp = new Date().toISOString().slice(0, 19).replace(/[:T]/g, "-");
      link.download = `ppms-backup-${stamp}.json`;
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);
      toast.success("Backup berhasil diunduh");
    },
    onError: (error: any) => {
      toast.error("Gagal mengunduh backup", error?.friendlyMessage || error?.response?.data?.message);
    },
  });
}

// useImport meng-upload file JSON dan me-restore data.
export function useImport() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation<ImportResult, any, File>({
    mutationFn: (file: File) => importExportApi.importData(file),
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: ["projects"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard"] });
      toast.success(
        "Import selesai",
        `${result.imported} project diimpor, ${result.skipped} dilewati`
      );
    },
    onError: (error: any) => {
      toast.error("Gagal mengimpor data", error?.friendlyMessage || error?.response?.data?.message);
    },
  });
}
