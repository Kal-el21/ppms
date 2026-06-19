import { useParams } from "react-router-dom";
import { useProjectDetail, useProjectMembers } from "../hooks/useProjects";
import { useMilestones, useCreateMilestone } from "../../milestones/hooks/useMilestones";
import { useTasks, useCreateTask, useChangeTaskStatus, useUpdateTaskProgress } from "../../tasks/hooks/useTasks";
import { Badge } from "../../../components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "../../../components/ui/card";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { useState } from "react";
import HandoverSection from "@/features/handovers/components/HandoverSection";

export default function ProjectDetailPage() {
  const { id } = useParams();
  const projectId = Number(id);

  const { data: project, isLoading } = useProjectDetail(projectId);
  const { data: members } = useProjectMembers(projectId);
  const { data: milestones } = useMilestones(projectId);
  const { data: tasks } = useTasks(projectId);

  const { mutate: createMilestone } = useCreateMilestone(projectId);
  const { mutate: createTask } = useCreateTask(projectId);
  const { mutate: changeTaskStatus } = useChangeTaskStatus(projectId);
  const { mutate: updateProgress } = useUpdateTaskProgress(projectId);

  const [milestoneName, setMilestoneName] = useState("");
  const [taskTitle, setTaskTitle] = useState("");

  if (isLoading || !project) return <div>Loading...</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">{project.name}</h1>
        <Badge>{project.status}</Badge>
      </div>

      <p className="text-slate-500">Progress: {project.progress.toFixed(0)}%</p>

      <Card>
        <CardHeader>
          <CardTitle>Members</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          {members?.map((m) => (
            <div key={m.id} className="flex justify-between text-sm">
              <span>User #{m.user_id}</span>
              <Badge variant="outline">{m.project_role}</Badge>
            </div>
          ))}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Milestones</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="flex gap-2">
            <Input
              placeholder="Milestone name"
              value={milestoneName}
              onChange={(e) => setMilestoneName(e.target.value)}
            />
            <Button
              onClick={() => {
                createMilestone({ name: milestoneName, description: "" });
                setMilestoneName("");
              }}
            >
              Add
            </Button>
          </div>

          {milestones?.map((m) => (
            <div key={m.id} className="flex justify-between border-b pb-2 text-sm">
              <span>{m.name}</span>
              <div className="flex items-center gap-2">
                <span>{m.progress.toFixed(0)}%</span>
                <Badge variant="outline">{m.status}</Badge>
              </div>
            </div>
          ))}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Tasks</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="flex gap-2">
            <Input
              placeholder="Task title"
              value={taskTitle}
              onChange={(e) => setTaskTitle(e.target.value)}
            />
            <Button
              onClick={() => {
                createTask({ title: taskTitle, priority: "MEDIUM" });
                setTaskTitle("");
              }}
            >
              Add
            </Button>
          </div>

          {tasks?.map((t) => (
            <div key={t.id} className="flex items-center justify-between border-b pb-2 text-sm">
              <span>{t.title}</span>
              <div className="flex items-center gap-2">
                <input
                  type="range"
                  min={0}
                  max={100}
                  value={t.progress}
                  onChange={(e) =>
                    updateProgress({ taskId: t.id, progress: Number(e.target.value), version: t.version })
                  }
                  className="w-24"
                />
                <span className="w-10 text-right">{t.progress}%</span>
                <select
                  value={t.status}
                  onChange={(e) => changeTaskStatus({ taskId: t.id, status: e.target.value, version: t.version })}
                  className="rounded border px-2 py-1 text-xs"
                >
                  <option value="TODO">TODO</option>
                  <option value="IN_PROGRESS">IN_PROGRESS</option>
                  <option value="DONE">DONE</option>
                  <option value="CANCELLED">CANCELLED</option>
                </select>
              </div>
            </div>
          ))}
        </CardContent>
        <HandoverSection projectId={projectId} />
      </Card>
    </div>
  );
}