import { useParams, useNavigate } from "react-router-dom";
import { useState } from "react";
import { useTaskDetail, useUpdateTask, useDeleteTask, useAssignTaskUsers, useAddTaskComment, useTaskComments } from "../hooks/useTasks";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "../../../components/ui/card";

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

  if (isLoading) return <div>Loading...</div>;
  if (!task) return <div>Task not found</div>;

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
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Task #{task.id} — {task.title}</h1>
      </div>

      <Card>
        <CardHeader><CardTitle>Details</CardTitle></CardHeader>
        <CardContent>
          {!editMode ? (
            <div>
              <p><strong>Description:</strong> {task.description}</p>
              <p><strong>Priority:</strong> {task.priority}</p>
              <p><strong>Assignees:</strong> {task.assignee_ids.length ? task.assignee_ids.join(", ") : "-"}</p>
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
              <div key={c.id} className="border-b pb-2">
                <p className="text-sm"><strong>{c.user.full_name}</strong> — {new Date(c.created_at).toLocaleString()}</p>
                <p className="text-sm text-ink-secondary">{c.comment}</p>
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
