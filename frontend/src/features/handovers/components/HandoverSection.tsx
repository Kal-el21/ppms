import { useState } from "react";
import { useHandovers, useCreateHandover, useMarkReceived } from "../hooks/useHandovers";
import { Card, CardContent, CardHeader, CardTitle } from "../../../components/ui/card";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { Badge } from "../../../components/ui/badge";

interface HandoverSectionProps {
  projectId: number;
}

export default function HandoverSection({ projectId }: HandoverSectionProps) {
  const { data: handovers } = useHandovers(projectId);
  const { mutate: createHandover } = useCreateHandover(projectId);
  const { mutate: markReceived } = useMarkReceived(projectId);

  const [description, setDescription] = useState("");

  return (
    <Card>
      <CardHeader>
        <CardTitle>Handovers</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="flex gap-2">
          <Input
            placeholder="Shipment description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
          <Button
            onClick={() => {
              createHandover({ description });
              setDescription("");
            }}
          >
            Record Shipment
          </Button>
        </div>

        {handovers?.map((h) => (
          <div key={h.id} className="flex items-center justify-between border-b pb-2 text-sm">
            <div>
              <p>{h.description}</p>
              <p className="text-xs text-slate-500">{new Date(h.created_at).toLocaleString()}</p>
            </div>
            <div className="flex items-center gap-2">
              <Badge variant={h.status === "RECEIVED" ? "default" : h.status === "RETURNED" ? "destructive" : "outline"}>
                {h.status}
              </Badge>
              {h.status === "PENDING" && (
                <Button size="sm" variant="ghost" onClick={() => markReceived({ handoverId: h.id, version: h.version })}>
                  Mark Received
                </Button>
              )}
            </div>
          </div>
        ))}
      </CardContent>
    </Card>
  );
}