import { useParams } from "react-router-dom";
import { useState } from "react";
import type { TaskPriority } from "../../tasks/types";
import { useProjectDetail, useProjectMembers, useChangeProjectStatus } from "../hooks/useProjects";
import { useMilestones, useCreateMilestone, type CreateMilestonePayload } from "../../milestones/hooks/useMilestones";
import { useTasks, useCreateTask, useChangeTaskStatus, useUpdateTaskProgress, type CreateTaskPayload } from "../../tasks/hooks/useTasks";

import { StatusBadge, getStatusColor } from "../../../components/ui/status-badge";
import { HealthBar } from "../../../components/ui/health-bar";
import { ProgressRing } from "../../../components/ui/progress-ring";
import { Avatar, AvatarStack } from "../../../components/ui/avatar";
import { Tabs } from "../../../components/ui/tabs";
import { Card, CardHeader, CardTitle, CardContent } from "../../../components/ui/card";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { EmptyState } from "../../../components/shared/EmptyState";

import BudgetSection from "../../budgets/components/BudgetSection";
import HandoverSection from "../../handovers/components/HandoverSection";
import FileUploadCard from "../../../components/shared/FileUploadCard";
import { Label } from "@/components/ui/label";
import { Select } from "@/components/ui/select";

const icon = (d: string, size = 22) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
    <path d={d} />
  </svg>
);

export default function ProjectDetailPage() {
  const { id } = useParams();
  const projectId = Number(id);

  const { data: project, isLoading } = useProjectDetail(projectId);
  const { data: members } = useProjectMembers(projectId);
  const { data: milestones } = useMilestones(projectId);
  const { data: tasks } = useTasks(projectId);

  const { mutate: changeStatus } = useChangeProjectStatus(projectId);
  const { mutate: createMilestone } = useCreateMilestone(projectId);
  const { mutate: createTask } = useCreateTask(projectId);
  const { mutate: changeTaskStatus } = useChangeTaskStatus(projectId);
  const { mutate: updateProgress } = useUpdateTaskProgress(projectId);

  // Update state form milestone — semua string agar konsisten dengan API & DateInput
  const [milestoneForm, setMilestoneForm] = useState<CreateMilestonePayload>({
    name: "",
    description: "",
    start_date: null,
    end_date: null,
  });

  const [taskForm, setTaskForm] = useState<CreateTaskPayload>({
    title: "",
    description: "",
    priority: "MEDIUM",
    due_date: null,
  });

  const handleCreateMilestone = () => {
    if (!milestoneForm.name.trim()) return;
    createMilestone({
      name: milestoneForm.name,
      description: milestoneForm.description,
      start_date: milestoneForm.start_date || undefined,
      end_date: milestoneForm.end_date || undefined,
    });
    setMilestoneForm({ name: "", description: "", start_date: null, end_date: null });
  };

  const handleCreateTask = () => {
    if (!taskForm.title.trim()) return;
    createTask({
      title: taskForm.title,
      description: taskForm.description,
      priority: taskForm.priority,
      due_date: taskForm.due_date || undefined,
    });
    setTaskForm({ title: "", description: "", priority: "MEDIUM", due_date: null });
  };

  if (isLoading || !project) {
    return <div className="text-ink-secondary text-sm">Memuat project...</div>;
  }

  const atRisk = project.status === "ON_HOLD" || project.progress < 30;
  const pmMembers = members?.filter((m) => m.project_role === "PROJECT_MANAGER") ?? [];

  return (
    <div>
      {/* Header */}
      <div className="flex items-start justify-between mb-6">
        <div className="flex items-start gap-4">
          <ProgressRing progress={project.progress} atRisk={atRisk} />
          <div>
            <div className="flex items-center gap-2.5 mb-1.5">
              <h1 className="text-[20px] font-semibold tracking-tight m-0">{project.name}</h1>
              <StatusBadge color={getStatusColor(project.status)}>{project.status}</StatusBadge>
            </div>
            <p className="text-[13px] text-ink-secondary m-0 max-w-xl">{project.description || "Tidak ada deskripsi."}</p>
            {pmMembers.length > 0 && (
              <div className="flex items-center gap-2 mt-2.5">
                <span className="text-[11.5px] text-ink-tertiary">PM:</span>
                <AvatarStack names={pmMembers.map((m) => `User ${m.user_id}`)} />
              </div>
            )}
          </div>
        </div>

        <select
          value={project.status}
          onChange={(e) => changeStatus({ status: e.target.value, version: project.version })}
          className="h-9 px-3 rounded-md border border-border-strong bg-surface text-[13px] font-medium cursor-pointer focus:outline-none focus:ring-2 focus:ring-primary-500"
        >
          <option value="PLANNED">Planned</option>
          <option value="ACTIVE">Active</option>
          <option value="ON_HOLD">On hold</option>
          <option value="COMPLETED">Completed</option>
          <option value="CANCELLED">Cancelled</option>
        </select>
      </div>

      {/* Tabs */}
      <Tabs
        tabs={[
          { key: "overview", label: "Overview" },
          { key: "milestones", label: "Milestones", count: milestones?.length },
          { key: "tasks", label: "Tasks", count: tasks?.length },
          { key: "budget", label: "Budget" },
          { key: "handovers", label: "Handovers" },
          { key: "files", label: "Files" },
          { key: "members", label: "Members", count: members?.length },
        ]}
        defaultTab="overview"
      >
        {(activeTab) => (
          <>
            {activeTab === "overview" && (
              <div className="grid grid-cols-1 lg:grid-cols-[1.3fr_1fr] gap-4">
                <Card>
                  <CardHeader>
                    <CardTitle>Milestone progress</CardTitle>
                  </CardHeader>
                  <CardContent>
                    {!milestones || milestones.length === 0 ? (
                      <EmptyState
                        icon={icon("M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5", 20)}
                        title="Belum ada milestone"
                        description="Tambahkan milestone untuk membagi project menjadi target-target utama."
                      />
                    ) : (
                      <div className="flex flex-col gap-3.5">
                        {milestones.map((m) => (
                          <div key={m.id}>
                            <div className="flex items-center justify-between mb-1.5">
                              <span className="text-[13px] font-medium">{m.name}</span>
                              <span className="text-xs font-semibold text-ink-secondary">{m.progress.toFixed(0)}%</span>
                            </div>
                            <HealthBar progress={m.progress} atRisk={m.status === "CANCELLED"} />
                          </div>
                        ))}
                      </div>
                    )}
                  </CardContent>
                </Card>

                <Card>
                  <CardHeader>
                    <CardTitle>Ringkasan task</CardTitle>
                  </CardHeader>
                  <CardContent>
                    {!tasks || tasks.length === 0 ? (
                      <p className="text-[12.5px] text-ink-tertiary">Belum ada task.</p>
                    ) : (
                      <div className="flex flex-col gap-2.5">
                        {(["TODO", "IN_PROGRESS", "DONE", "CANCELLED"] as const).map((status) => {
                          const count = tasks.filter((t) => t.status === status).length;
                          return (
                            <div key={status} className="flex items-center justify-between">
                              <StatusBadge color={getStatusColor(status)}>{status}</StatusBadge>
                              <span className="text-[13px] font-semibold">{count}</span>
                            </div>
                          );
                        })}
                      </div>
                    )}
                  </CardContent>
                </Card>
              </div>
            )}

            {activeTab === "milestones" && (
              <Card>
                <CardHeader>
                  <CardTitle>Milestones</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5 mb-5">
                    <div>
                      <Label>Nama milestone</Label>
                      <Input placeholder="cth. Phase 1 — Discovery" value={milestoneForm.name} onChange={(e) => setMilestoneForm({ ...milestoneForm, name: e.target.value })} />
                    </div>
                    <div>
                      <Label>Deskripsi</Label>
                      <Input placeholder="Opsional" value={milestoneForm.description} onChange={(e) => setMilestoneForm({ ...milestoneForm, description: e.target.value })} />
                    </div>
                    <div>
                      <Label>Tanggal mulai</Label>
                      <Input type="date" value={milestoneForm.start_date ?? ""} onChange={(e) => setMilestoneForm({ ...milestoneForm, start_date: e.target.value })} />
                    </div>
                    <div>
                      <Label>Tanggal selesai</Label>
                      <Input type="date" value={milestoneForm.end_date ?? ""} onChange={(e) => setMilestoneForm({ ...milestoneForm, end_date: e.target.value })} />
                    </div>
                  </div>
                  <Button variant="primary" onClick={handleCreateMilestone}>
                    Tambah milestone
                  </Button>

                  {!milestones || milestones.length === 0 ? (
                    <EmptyState
                      icon={icon("M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5", 20)}
                      title="Belum ada milestone"
                    />
                  ) : (
                    <div className="flex flex-col gap-4">
                      {milestones.map((m) => (
                        <div key={m.id} className="flex items-center gap-4 pb-4 border-b border-border last:border-b-0 last:pb-0">
                          <div className="flex-1 min-w-0">
                            <div className="flex items-center justify-between mb-1.5">
                              <span className="text-[13.5px] font-medium">{m.name}</span>
                              <StatusBadge color={getStatusColor(m.status)}>{m.status}</StatusBadge>
                            </div>
                            <HealthBar progress={m.progress} atRisk={m.status === "CANCELLED"} />
                          </div>
                          <span className="text-[13px] font-semibold text-ink-secondary w-10 text-right flex-shrink-0">
                            {m.progress.toFixed(0)}%
                          </span>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            )}

            {activeTab === "tasks" && (
              <Card>
                <CardHeader>
                  <CardTitle>Tasks</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5 mb-5">
                    <div>
                      <Label>Judul task</Label>
                      <Input placeholder="cth. Setup CI/CD pipeline" value={taskForm.title} onChange={(e) => setTaskForm({ ...taskForm, title: e.target.value })} />
                    </div>
                    <div>
                      <Label>Priority</Label>
                      <Select
                        value={taskForm.priority}
                        onChange={(e) => setTaskForm({ ...taskForm, priority: e.target.value as TaskPriority })}
                        options={[
                          { value: "LOW", label: "Low" },
                          { value: "MEDIUM", label: "Medium" },
                          { value: "HIGH", label: "High" },
                          { value: "URGENT", label: "Urgent" },
                        ]}
                      />
                    </div>
                    <div>
                      <Label>Deskripsi</Label>
                      <Input placeholder="Opsional" value={taskForm.description} onChange={(e) => setTaskForm({ ...taskForm, description: e.target.value })} />
                    </div>
                    <div>
                      <Label>Due date</Label>
                      <Input type="date" value={taskForm.due_date ?? ""} onChange={(e) => setTaskForm({ ...taskForm, due_date: e.target.value })} />
                    </div>
                  </div>
                  <Button variant="primary" onClick={handleCreateTask}>
                    Tambah task
                  </Button>

                  {!tasks || tasks.length === 0 ? (
                    <EmptyState
                      icon={icon("M9 11l3 3L22 4M21 12v7a2 2 0 01-2 2H5a2 2 0 01-2-2V5a2 2 0 012-2h11", 20)}
                      title="Belum ada task"
                    />
                  ) : (
                    <div className="flex flex-col">
                      {tasks.map((t) => (
                        <div
                          key={t.id}
                          className="flex items-center gap-3 py-3 border-b border-border last:border-b-0"
                        >
                          <span
                            className="h-2 w-2 rounded-sm flex-shrink-0"
                            style={{
                              background:
                                t.priority === "URGENT"
                                  ? "#DC2626"
                                  : t.priority === "HIGH"
                                  ? "#D97706"
                                  : t.priority === "MEDIUM"
                                  ? "#2563EB"
                                  : "#94A3B8",
                            }}
                          />
                          <span className="text-[13px] font-medium flex-1 min-w-0 truncate">{t.title}</span>

                          <div className="w-24 flex-shrink-0">
                            <input
                              type="range"
                              min={0}
                              max={100}
                              value={t.progress}
                              onChange={(e) =>
                                updateProgress({ taskId: t.id, progress: Number(e.target.value), version: t.version })
                              }
                              className="w-full"
                            />
                          </div>
                          <span className="text-xs font-semibold text-ink-secondary w-8 text-right flex-shrink-0">
                            {t.progress}%
                          </span>

                          <select
                            value={t.status}
                            onChange={(e) => changeTaskStatus({ taskId: t.id, status: e.target.value, version: t.version })}
                            className="h-7 px-2 rounded-md border border-border-strong bg-surface text-[11.5px] font-medium flex-shrink-0 cursor-pointer"
                          >
                            <option value="TODO">Todo</option>
                            <option value="IN_PROGRESS">In progress</option>
                            <option value="DONE">Done</option>
                            <option value="CANCELLED">Cancelled</option>
                          </select>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            )}

            {activeTab === "budget" && <BudgetSection projectId={projectId} />}

            {activeTab === "handovers" && <HandoverSection projectId={projectId} />}

            {activeTab === "files" && <FileUploadCard entityType="PROJECT" entityId={projectId} />}

            {activeTab === "members" && (
              <Card>
                <CardHeader>
                  <CardTitle>Members</CardTitle>
                </CardHeader>
                <CardContent>
                  {!members || members.length === 0 ? (
                    <EmptyState
                      icon={icon("M16 21v-2a4 4 0 00-4-4H6a4 4 0 00-4 4v2M9 11a4 4 0 100-8 4 4 0 000 8z", 20)}
                      title="Belum ada member"
                    />
                  ) : (
                    <div className="flex flex-col">
                      {members.map((m) => (
                        <div key={m.id} className="flex items-center gap-3 py-3 border-b border-border last:border-b-0">
                          <Avatar name={`User ${m.user_id}`} />
                          <div className="flex-1 min-w-0">
                            <p className="text-[13px] font-medium m-0">User #{m.user_id}</p>
                            <p className="text-[11.5px] text-ink-tertiary m-0">
                              Bergabung {new Date(m.joined_at).toLocaleDateString("id-ID")}
                            </p>
                          </div>
                          <StatusBadge color={m.project_role === "PROJECT_MANAGER" ? "blue" : "gray"}>
                            {m.project_role}
                          </StatusBadge>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            )}
          </>
        )}
      </Tabs>
    </div>
  );
}