import { useParams, useNavigate } from "react-router-dom";
import { useState, useMemo } from "react";
import { useMilestones, useUpdateMilestone, useDeleteMilestone } from "../hooks/useMilestones";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "../../../components/ui/card";
import { StatusBadge, getStatusColor } from "../../../components/ui/status-badge";
import { EmptyState } from "../../../components/shared/EmptyState";
import { CardSkeleton } from "../../../components/ui/skeleton";

const notFoundIcon = (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
    <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
  </svg>
);

export default function MilestoneDetailPage() {
  const { projectId, milestoneId } = useParams();
  const pId = Number(projectId);
  const mId = Number(milestoneId);
  const navigate = useNavigate();

  const { data: milestones, isLoading } = useMilestones(pId);
  const milestone = useMemo(() => milestones?.find((m: any) => m.id === mId), [milestones, mId]);

  const { mutate: updateMilestone } = useUpdateMilestone(pId);
  const { mutate: deleteMilestone } = useDeleteMilestone(pId);

  const [editMode, setEditMode] = useState(false);
  const [name, setName] = useState(milestone?.name || "");
  const [description, setDescription] = useState(milestone?.description || "");

  if (isLoading) return <CardSkeleton />;
  if (!milestone) {
    return (
      <EmptyState
        icon={notFoundIcon}
        title="Milestone tidak ditemukan"
        description="Milestone yang Anda cari mungkin telah dihapus atau tidak tersedia."
      />
    );
  }

  const handleSave = () => {
    updateMilestone({ milestoneId: mId, payload: { name, description, version: milestone.version } });
    setEditMode(false);
  };

  const handleDelete = () => {
    if (!confirm("Hapus milestone ini?")) return;
    deleteMilestone(mId, { onSuccess: () => navigate(-1) });
  };

  return (
    <div className="space-y-6">
      <div className="flex items-start justify-between">
        <h1 className="text-[20px] font-semibold tracking-tight m-0">Milestone #{milestone.id} — {milestone.name}</h1>
        <StatusBadge color={getStatusColor(milestone.status)}>{milestone.status}</StatusBadge>
      </div>

      <Card>
        <CardHeader><CardTitle>Details</CardTitle></CardHeader>
        <CardContent>
          {!editMode ? (
            <div className="space-y-4">
              <div>
                <p className="text-[12.5px] text-ink-tertiary mb-1">Description</p>
                <p className="text-[13.5px] text-ink-primary m-0">{milestone.description}</p>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <p className="text-[12.5px] text-ink-tertiary mb-1">Start</p>
                  <p className="text-[13.5px] text-ink-primary m-0">
                    {milestone.start_date ? new Date(milestone.start_date).toLocaleDateString() : "-"}
                  </p>
                </div>
                <div>
                  <p className="text-[12.5px] text-ink-tertiary mb-1">End</p>
                  <p className="text-[13.5px] text-ink-primary m-0">
                    {milestone.end_date ? new Date(milestone.end_date).toLocaleDateString() : "-"}
                  </p>
                </div>
              </div>
            </div>
          ) : (
            <div className="space-y-3">
              <div>
                <Label>Name</Label>
                <Input value={name} onChange={(e) => setName(e.target.value)} />
              </div>
              <div>
                <Label>Description</Label>
                <Input value={description} onChange={(e) => setDescription(e.target.value)} />
              </div>
              <div className="flex gap-2">
                <Button onClick={handleSave}>Save</Button>
                <Button variant="outline" onClick={() => setEditMode(false)}>Cancel</Button>
              </div>
            </div>
          )}

          <div className="mt-4 flex gap-2">
            <Button onClick={() => setEditMode((s) => !s)}>{editMode ? "Close" : "Edit Milestone"}</Button>
            <Button variant="danger" onClick={handleDelete}>Delete</Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
