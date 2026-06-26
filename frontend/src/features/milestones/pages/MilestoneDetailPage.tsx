import { useParams, useNavigate } from "react-router-dom";
import { useState, useMemo } from "react";
import { useMilestones, useUpdateMilestone, useDeleteMilestone } from "../hooks/useMilestones";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "../../../components/ui/card";

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

  if (isLoading) return <div>Loading...</div>;
  if (!milestone) return <div>Milestone not found</div>;

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
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Milestone #{milestone.id} — {milestone.name}</h1>
      </div>

      <Card>
        <CardHeader><CardTitle>Details</CardTitle></CardHeader>
        <CardContent>
          {!editMode ? (
            <div>
              <p><strong>Description:</strong> {milestone.description}</p>
              <p><strong>Start:</strong> {milestone.start_date ? new Date(milestone.start_date).toLocaleDateString() : "-"}</p>
              <p><strong>End:</strong> {milestone.end_date ? new Date(milestone.end_date).toLocaleDateString() : "-"}</p>
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
