import { useRef, useState } from "react";
import { Button } from "../ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "../ui/card";
import { EmptyState } from "./EmptyState";
import {
  useAttachments,
  useUploadAttachment,
  useDeleteAttachment,
  useAttachmentVersions,
  handleDownload,
} from "../../features/attachments/hooks/useAttachments";
import type { EntityType } from "../../features/attachments/types";

interface FileUploadCardProps {
  entityType: EntityType;
  entityId: number;
}

export default function FileUploadCard({ entityType, entityId }: FileUploadCardProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [versionAttachmentId, setVersionAttachmentId] = useState<number | null>(null);
  const { data: attachments, isLoading } = useAttachments(entityType, entityId);
  const { data: versions, isLoading: versionsLoading } = useAttachmentVersions(versionAttachmentId ?? 0);
  const { mutate: upload, isPending } = useUploadAttachment(entityType, entityId);
  const { mutate: deleteFile } = useDeleteAttachment(entityType, entityId);

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      upload(file);
      e.target.value = "";
    }
  };

  const formatSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Attachments</CardTitle>
      </CardHeader>
      <CardContent>
        <input
          ref={fileInputRef}
          type="file"
          className="hidden"
          onChange={handleFileSelect}
          accept=".pdf,.jpg,.jpeg,.png,.webp,.doc,.docx,.xls,.xlsx,.csv,.txt"
        />
        <Button variant="outline" size="sm" onClick={() => fileInputRef.current?.click()} disabled={isPending} className="mb-4">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M17 8l-5-5-5 5M12 3v12" />
          </svg>
          {isPending ? "Mengupload..." : "Upload file"}
        </Button>

        {isLoading ? (
          <p className="text-[12.5px] text-ink-tertiary">Memuat attachments...</p>
        ) : attachments?.length === 0 ? (
          <EmptyState
            icon={
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
                <path d="M13 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V9z" />
                <path d="M13 2v7h7" />
              </svg>
            }
            title="Belum ada file"
          />
        ) : (
          <div className="flex flex-col">
            {attachments?.map((a) => (
              <div key={a.id} className="flex items-center justify-between gap-3 py-2.5 border-b border-border last:border-b-0">
                <div className="flex items-center gap-2.5 min-w-0">
                  <div className="h-8 w-8 rounded-md bg-primary-50 dark:bg-primary-900/30 text-primary-700 dark:text-primary-400 flex items-center justify-center flex-shrink-0">
                    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M13 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V9z" />
                      <path d="M13 2v7h7" />
                    </svg>
                  </div>
                  <div className="min-w-0">
                    <p className="text-[13px] font-medium m-0 truncate">{a.original_name}</p>
                    <p className="text-[11px] text-ink-tertiary m-0">
                      {formatSize(a.file_size)} &middot; v{a.version}
                    </p>
                  </div>
                </div>
                <div className="flex gap-1 flex-shrink-0">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setVersionAttachmentId(versionAttachmentId === a.id ? null : a.id)}
                  >
                    Versions
                  </Button>
                  <Button variant="ghost" size="sm" onClick={() => handleDownload(a.id)}>
                    Download
                  </Button>
                  <Button variant="ghost" size="sm" onClick={() => deleteFile(a.id)}>
                    Hapus
                  </Button>
                </div>
              </div>
            ))}
          </div>
        )}

        {versionAttachmentId && (
          <div className="mt-4 rounded-md border border-border bg-surface-secondary p-3">
            <p className="text-[12px] font-semibold uppercase tracking-wide text-ink-tertiary mb-2">Version history</p>
            {versionsLoading ? (
              <p className="text-[12.5px] text-ink-tertiary m-0">Memuat versions...</p>
            ) : !versions || versions.length === 0 ? (
              <p className="text-[12.5px] text-ink-tertiary m-0">Belum ada version history.</p>
            ) : (
              <div className="flex flex-col">
                {versions.map((version) => (
                  <div key={version.id} className="flex items-center justify-between gap-3 py-2 border-b border-border last:border-b-0">
                    <div className="min-w-0">
                      <p className="text-[12.5px] font-medium m-0 truncate">
                        v{version.version} - {version.original_name}
                      </p>
                      <p className="text-[11px] text-ink-tertiary m-0">
                        {formatSize(version.file_size)} &middot; {new Date(version.created_at).toLocaleString("id-ID")}
                      </p>
                    </div>
                    <Button variant="ghost" size="sm" onClick={() => handleDownload(version.id)}>
                      Download
                    </Button>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
