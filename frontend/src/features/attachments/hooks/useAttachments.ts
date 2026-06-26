import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { attachmentApi } from "../api/attachmentApi";
import type { EntityType } from "../types";

export function useAttachments(entityType: EntityType, entityId: number) {
  return useQuery({
    queryKey: ["attachments", entityType, entityId],
    queryFn: () => attachmentApi.getByEntity(entityType, entityId),
    enabled: !!entityId,
  });
}

export function useUploadAttachment(entityType: EntityType, entityId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (file: File) => attachmentApi.upload(file, entityType, entityId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["attachments", entityType, entityId] }),
  });
}

export function useDeleteAttachment(entityType: EntityType, entityId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => attachmentApi.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["attachments", entityType, entityId] }),
  });
}

export function useAttachmentVersions(id: number) {
  return useQuery({
    queryKey: ["attachments", id, "versions"],
    queryFn: () => attachmentApi.getVersions(id),
    enabled: !!id,
  });
}

export async function handleDownload(id: number) {
  const { download_url, file_name } = await attachmentApi.getDownloadUrl(id);
  const link = document.createElement("a");
  link.href = download_url;
  link.download = file_name;
  link.click();
}
