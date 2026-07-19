import { useParams, useNavigate } from "react-router-dom";
import { useState } from "react";
import { useTaskDetail, useUpdateTask, useDeleteTask, useAssignTaskUsers, useAddTaskComment, useTaskComments } from "../hooks/useTasks";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "../../../components/ui/card";
import { StatusBadge, getStatusColor } from "../../../components/ui/status-badge";
import { EmptyState } from "../../../components/shared/EmptyState";
import { CardSkeleton } from "../../../components/ui/skeleton";

const notFoundIcon = (
  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
    <path d="M9 11l3 3L22 4M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11" />
  </svg>
);

export default function TaskDetailPage() {
  const { projectId, taskId } = useParams();
  const pId = Number(projectId);
  const tId = Number(taskId);
  const navigate = useNavigate();

  const { data: task, isLoading } = useTaskDetail(pId, tId);
  const { mutate: updateTask, isPending: updating } = useUpdateTask(pId);
  const { mutate: deleteTask, isPending: deleting } = useDeleteTask(pId);
  const { mutate: assignUsers } = useAssignTaskUsers(pId);
  const { mutate: addComment } = useAddTaskComment(pId, tId);
  const { data: comments } = useTaskComments(pId, tId);

  const [editMode, setEditMode] = useState(false);
  const [title, setTitle] = useState(task?.title || "");
  const [description, setDescription] = useState(task?.description || "");
  const [assigneeInput, setAssigneeInput] = useState("");
  const [newComment, setNewComment] = useState("");

  if (isLoading) return <CardSkeleton />;
  if (!task) {
    return (
      <EmptyState
        icon={notFoundIcon}
        title="Task tidak ditemukan"
        description="Task yang Anda cari mungkin telah dihapus atau tidak tersedia."
      />
    );
  }

  const handleSave = () => {
    updateTask({ taskId: tId, payload: { title, description, version: task.version } });
    setEditMode(false);
  };

  const handleDelete = () => {
    if (!confirm("Hapus task ini?")) return;
    deleteTask(tId, { onSuccess: () => navigate(-1) });
  };

  const handleAssign = () => {
    const ids = assigneeInput.split(",").map((s) => Number(s.trim())).filter(Boolean);
    if (ids.length === 0) return;
    assignUsers({ taskId: tId, userIds: ids });
    setAssigneeInput("");
  };

  const handleAddComment = () => {
    if (!newComment.trim()) return;
    addComment(newComment);
    setNewComment("");
  };

  return (
    <div className="space-y-6">
      <div className="flex items-start justify-between">
        <h1 className="text-[20px] font-semibold tracking-tight m-0">Task #{task.id} — {task.title}</h1>
      </div>

      <Card>
        <CardHeader><CardTitle>Details</CardTitle></CardHeader>
        <CardContent>
          {!editMode ? (
            <div className="space-y-4">
              <div className="flex flex-wrap items-center gap-2">
                <StatusBadge color={getStatusColor(task.status)}>{task.status}</StatusBadge>
                <StatusBadge color={getStatusColor(task.priority)}>{task.priority}</StatusBadge>
              </div>
              <div>
                <p className="text-[12.5px] text-ink-tertiary mb-1">Description</p>
                <p className="text-[13.5px] text-ink-primary m-0">{task.description}</p>
              </div>
              <div>
                <p className="text-[12.5px] text-ink-tertiary mb-1">Assignees</p>
                <p className="text-[13.5px] text-ink-primary m-0">
                  {task.assignee_ids.length ? task.assignee_ids.join(", ") : "-"}
                </p>
              </div>
            </div>
          ) : (
            <div className="space-y-3">
              <div>
                <Label>Title</Label>
                <Input value={title} onChange={(e) => setTitle(e.target.value)} />
              </div>
              <div>
                <Label>Description</Label>
                <Input value={description} onChange={(e) => setDescription(e.target.value)} />
              </div>
              <div className="flex gap-2">
                <Button onClick={handleSave} disabled={updating}>Save</Button>
                <Button variant="outline" onClick={() => setEditMode(false)}>Cancel</Button>
              </div>
            </div>
          )}

          <div className="mt-4 flex gap-2">
            <Button onClick={() => setEditMode((s) => !s)}>{editMode ? "Close" : "Edit Task"}</Button>
            <Button variant="danger" onClick={handleDelete} disabled={deleting}>Delete</Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader><CardTitle>Assignees</CardTitle></CardHeader>
        <CardContent>
          <div className="flex gap-2">
            <Input placeholder="Comma-separated user IDs (e.g. 2,3)" value={assigneeInput} onChange={(e) => setAssigneeInput(e.target.value)} />
            <Button onClick={handleAssign}>Assign</Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader><CardTitle>Comments</CardTitle></CardHeader>
        <CardContent>
          <div className="space-y-3">
            {(comments?.data || []).map((c: any) => (
              <div key={c.id} className="border-b border-border pb-2">
                <p className="text-[12.5px] text-ink-primary m-0">
                  <strong>{c.user.full_name}</strong> — {new Date(c.created_at).toLocaleString()}
                </p>
                <p className="text-[13px] text-ink-secondary m-0 mt-0.5">{c.comment}</p>
              </div>
            ))}

            <div className="mt-3">
              <Input placeholder="Add comment..." value={newComment} onChange={(e) => setNewComment(e.target.value)} />
              <div className="mt-2"><Button onClick={handleAddComment}>Add Comment</Button></div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
