import { useState } from "react";
import { useHandovers, useCreateHandover, useMarkReceived, useMarkReturned } from "../hooks/useHandovers";
import { Card, CardHeader, CardTitle, CardContent } from "../../../components/ui/card";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { StatusBadge, getStatusColor } from "../../../components/ui/status-badge";
import { EmptyState } from "../../../components/shared/EmptyState";

interface HandoverSectionProps {
  projectId: number;
}

export default function HandoverSection({ projectId }: HandoverSectionProps) {
  const { data: handovers } = useHandovers(projectId);
  const { mutate: createHandover } = useCreateHandover(projectId);
  const { mutate: markReceived } = useMarkReceived(projectId);
  const { mutate: markReturned, isPending: returning } = useMarkReturned(projectId);

  const [description, setDescription] = useState("");
  const [returningId, setReturningId] = useState<number | null>(null);
  const [returnReason, setReturnReason] = useState("");

  const handleReturn = (handoverId: number, version: number) => {
    if (!returnReason.trim()) return;
    markReturned(
      { handoverId, reason: returnReason, version },
      {
        onSuccess: () => {
          setReturningId(null);
          setReturnReason("");
        },
      }
    );
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Handovers</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex gap-2 mb-5">
          <Input
            placeholder="Deskripsi pengiriman dokumen..."
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className="flex-1"
          />
          <Button
            variant="primary"
            onClick={() => {
              if (!description.trim()) return;
              createHandover({ description });
              setDescription("");
            }}
          >
            Kirim
          </Button>
        </div>

        {!handovers || handovers.length === 0 ? (
          <EmptyState
            icon={
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
                <path d="M21 8a2 2 0 00-1-1.73l-7-4a2 2 0 00-2 0l-7 4A2 2 0 003 8v8a2 2 0 001 1.73l7 4a2 2 0 002 0l7-4A2 2 0 0021 16z" />
              </svg>
            }
            title="Belum ada handover"
          />
        ) : (
          <div className="flex flex-col">
            {handovers.map((h) => (
              <div key={h.id} className="py-3 border-b border-border last:border-b-0">
                <div className="flex items-center justify-between gap-3">
                  <div className="min-w-0">
                    <p className="text-[13px] font-medium m-0 truncate">{h.description}</p>
                    <p className="text-[11.5px] text-ink-tertiary m-0 mt-0.5">
                      {new Date(h.created_at).toLocaleString("id-ID")}
                    </p>
                  </div>
                  <div className="flex items-center gap-2 flex-shrink-0">
                    <StatusBadge color={getStatusColor(h.status)}>{h.status}</StatusBadge>
                    {h.status === "PENDING" && (
                      <>
                        <Button size="sm" variant="ghost" onClick={() => markReceived({ handoverId: h.id, version: h.version })}>
                          Tandai diterima
                        </Button>
                        <Button size="sm" variant="outline" onClick={() => setReturningId(h.id)}>
                          Return
                        </Button>
                      </>
                    )}
                  </div>
                </div>
                {returningId === h.id && (
                  <div className="mt-3 flex gap-2">
                    <Input
                      placeholder="Alasan return..."
                      value={returnReason}
                      onChange={(e) => setReturnReason(e.target.value)}
                    />
                    <Button
                      size="sm"
                      variant="danger"
                      onClick={() => handleReturn(h.id, h.version)}
                      disabled={returning || !returnReason.trim()}
                    >
                      Simpan
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => {
                        setReturningId(null);
                        setReturnReason("");
                      }}
                      disabled={returning}
                    >
                      Batal
                    </Button>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
