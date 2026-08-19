import { useParams, useNavigate } from "react-router-dom";
import { useState } from "react";
import type { Task, TaskPriority } from "../../tasks/types";
import type { Milestone } from "../../milestones/types";
import type { ProjectRole, ProjectPriority, InitiationType } from "../types";
import {
  useProjectDetail,
  useProjectMembers,
  useChangeProjectStatus,
  useUpdateProject,
  useAddMember,
} from "../hooks/useProjects";
import {
  useMilestones,
  useCreateMilestone,
  useChangeMilestoneStatus,
  useReorderMilestones,
  useUpdateMilestone,
  useDeleteMilestone,
  type CreateMilestonePayload,
} from "../../milestones/hooks/useMilestones";
import {
  useTasks,
  useCreateTask,
  useChangeTaskStatus,
  useUpdateTaskProgress,
  useUpdateTask,
  useDeleteTask,
  useAssignTaskUsers,
  useTaskComments,
  useAddTaskComment,
  type CreateTaskPayload,
} from "../../tasks/hooks/useTasks";
import { useUsers } from "../../users/hooks/useUsers";

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
import { useAuth } from "@/features/auth/context/AuthContext";
import { ConfirmDeleteDialog } from "@/components/shared/ConfirmDeleteDialog";

const icon = (d: string, size = 22) => (
  <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
    <path d={d} />
  </svg>
);

function getPriorityColor(priority: string): string {
  const map: Record<string, string> = {
    URGENT: "var(--danger-600)",
    HIGH: "var(--warning-600)",
    MEDIUM: "var(--primary-600)",
    LOW: "var(--text-tertiary)",
  };
  return map[priority] || "var(--text-tertiary)";
}
const getHealthColor = (health: string) => {
  switch (health) {
    case "GREEN":
      return "green" as const;
    case "YELLOW":
      return "amber" as const;
    case "RED":
      return "red" as const;
    default:
      return "gray" as const;
  }
};

export default function ProjectDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const projectId = Number(id);
  const { user } = useAuth();

  const { data: project, isLoading } = useProjectDetail(projectId);
  const { data: members } = useProjectMembers(projectId);
  const { data: milestones } = useMilestones(projectId);
  const { data: tasks } = useTasks(projectId);

  const { mutate: changeStatus } = useChangeProjectStatus(projectId);
  const { mutate: updateProject } = useUpdateProject(projectId);
  const { mutate: createMilestone } = useCreateMilestone(projectId);
  const { mutate: changeMilestoneStatus } = useChangeMilestoneStatus(projectId);
  const { mutate: reorderMilestones } = useReorderMilestones(projectId);
  const { mutate: updateMilestone } = useUpdateMilestone(projectId);
  const deleteMilestoneMutation = useDeleteMilestone(projectId);

  const { mutate: createTask } = useCreateTask(projectId);
  const { mutate: changeTaskStatus } = useChangeTaskStatus(projectId);
  const { mutate: updateProgress } = useUpdateTaskProgress(projectId);
  const { mutate: updateTask } = useUpdateTask(projectId);
  const deleteTaskMutation = useDeleteTask(projectId);
  const { mutate: assignTaskUsers } = useAssignTaskUsers(projectId);
  const { mutate: addMember, isPending: addingMember } = useAddMember(projectId);
  const { data: usersData } = useUsers(1, 100, true);

  const users = usersData?.data ?? [];

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
  const [editingMilestoneId, setEditingMilestoneId] = useState<number | null>(null);
  const [milestoneEditForm, setMilestoneEditForm] = useState<CreateMilestonePayload>({
    name: "",
    description: "",
    start_date: null,
    end_date: null,
  });
  const [editingTaskId, setEditingTaskId] = useState<number | null>(null);
  const [taskEditForm, setTaskEditForm] = useState({
    title: "",
    description: "",
    priority: "MEDIUM" as TaskPriority,
    milestone_id: "",
    start_date: "",
    due_date: "",
  });
  const [assigningTaskId, setAssigningTaskId] = useState<number | null>(null);
  const [assigneeInput, setAssigneeInput] = useState("");
  const [commentingTaskId, setCommentingTaskId] = useState<number | null>(null);
  const [commentInput, setCommentInput] = useState("");
  const { data: taskComments } = useTaskComments(projectId, commentingTaskId ?? 0);
  const { mutate: addTaskComment } = useAddTaskComment(projectId, commentingTaskId ?? 0);
  const [editingProject, setEditingProject] = useState(false);
  const [projectEditForm, setProjectEditForm] = useState({
    name: "",
    description: "",
    category: "",
    initiation_type: "" as string,
    priority: "MEDIUM" as ProjectPriority,
    notes: "",
    start_date: "",
    end_date: "",
  });
  const [memberForm, setMemberForm] = useState({
    userId: null as number | null,
    project_role: "MEMBER" as ProjectRole,
  });

  const [milestoneDeleteId, setMilestoneDeleteId] = useState<number | null>(null);
  const [milestoneDeleteOpen, setMilestoneDeleteOpen] = useState(false);
  const [taskDeleteId, setTaskDeleteId] = useState<number | null>(null);
  const [taskDeleteOpen, setTaskDeleteOpen] = useState(false);

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

  const startEditMilestone = (milestone: Milestone) => {
    setEditingMilestoneId(milestone.id);
    setMilestoneEditForm({
      name: milestone.name,
      description: milestone.description || "",
      start_date: milestone.start_date?.slice(0, 10) ?? null,
      end_date: milestone.end_date?.slice(0, 10) ?? null,
    });
  };

  const saveMilestoneEdit = (milestone: Milestone) => {
    if (!milestoneEditForm.name.trim()) return;
    updateMilestone(
      {
        milestoneId: milestone.id,
        payload: {
          name: milestoneEditForm.name,
          description: milestoneEditForm.description || "",
          start_date: milestoneEditForm.start_date || null,
          end_date: milestoneEditForm.end_date || null,
          version: milestone.version,
        },
      },
      { onSuccess: () => setEditingMilestoneId(null) }
    );
  };

  const handleDeleteMilestone = (milestoneId: number) => {
    setMilestoneDeleteId(milestoneId);
    setMilestoneDeleteOpen(true);
  };

  const confirmDeleteMilestone = () => {
    if (milestoneDeleteId == null) return;
    deleteMilestoneMutation.mutate(milestoneDeleteId, { onSuccess: () => setMilestoneDeleteOpen(false) });
    setMilestoneDeleteId(null);
  };

  const moveMilestone = (index: number, direction: -1 | 1) => {
    if (!milestones) return;
    const nextIndex = index + direction;
    if (nextIndex < 0 || nextIndex >= milestones.length) return;

    const ordered = [...milestones];
    const [item] = ordered.splice(index, 1);
    if (!item) return;
    ordered.splice(nextIndex, 0, item);
    reorderMilestones(ordered.map((m) => m.id));
  };

  const startEditTask = (task: Task) => {
    setEditingTaskId(task.id);
    setTaskEditForm({
      title: task.title,
      description: task.description || "",
      priority: task.priority,
      milestone_id: task.milestone_id ? String(task.milestone_id) : "",
      start_date: task.start_date?.slice(0, 10) ?? "",
      due_date: task.due_date?.slice(0, 10) ?? "",
    });
  };

  const saveTaskEdit = (task: Task) => {
    if (!taskEditForm.title.trim()) return;
    updateTask(
      {
        taskId: task.id,
        payload: {
          title: taskEditForm.title,
          description: taskEditForm.description || "",
          priority: taskEditForm.priority,
          milestone_id: taskEditForm.milestone_id ? Number(taskEditForm.milestone_id) : null,
          start_date: taskEditForm.start_date || null,
          due_date: taskEditForm.due_date || null,
          version: task.version,
        },
      },
      { onSuccess: () => setEditingTaskId(null) }
    );
  };

  const handleDeleteTask = (taskId: number) => {
    setTaskDeleteId(taskId);
    setTaskDeleteOpen(true);
  };

  const confirmDeleteTask = () => {
    if (taskDeleteId == null) return;
    deleteTaskMutation.mutate(taskDeleteId, { onSuccess: () => navigate(-1) });
    setTaskDeleteId(null);
    setTaskDeleteOpen(false);
  };

  const startAssignTask = (task: Task) => {
    setAssigningTaskId(task.id);
    setAssigneeInput((task.assignee_ids ?? []).join(", "));
  };

  const saveTaskAssignees = (taskId: number) => {
    const userIds = assigneeInput
      .split(",")
      .map((value) => Number(value.trim()))
      .filter((value) => Number.isInteger(value) && value > 0);
    if (userIds.length === 0) return;
    assignTaskUsers({ taskId, userIds }, { onSuccess: () => setAssigningTaskId(null) });
  };

  const saveTaskComment = (taskId: number) => {
    if (!commentInput.trim()) return;
    addTaskComment(commentInput, { onSuccess: () => setCommentInput("") });
    setCommentingTaskId(taskId);
  };

  const startEditProject = () => {
    if (!project) return;
    setProjectEditForm({
      name: project.name,
      description: project.description || "",
      category: project.category || "",
      initiation_type: project.initiation_type ?? "",
      priority: project.priority,
      notes: project.notes || "",
      start_date: project.start_date?.slice(0, 10) ?? "",
      end_date: project.end_date?.slice(0, 10) ?? "",
    });
    setEditingProject(true);
  };

  const saveProjectEdit = () => {
    if (!project) return;
    if (!projectEditForm.name.trim()) return;
    updateProject(
      {
        name: projectEditForm.name,
        description: projectEditForm.description,
        category: projectEditForm.category,
        initiation_type: (projectEditForm.initiation_type || null) as InitiationType | null,
        priority: projectEditForm.priority,
        notes: projectEditForm.notes,
        start_date: projectEditForm.start_date || null,
        end_date: projectEditForm.end_date || null,
        version: project.version,
      },
      { onSuccess: () => setEditingProject(false) }
    );
  };

  const handleAddMember = () => {
    const userId = memberForm.userId;
    if (!userId || userId <= 0) return;
    addMember(
      { userId, projectRole: memberForm.project_role },
      {
        onSuccess: () =>
          setMemberForm({
            userId: null,
            project_role: "MEMBER",
          }),
      }
    );
  };

  if (isLoading || !project) {
    return <div className="text-ink-secondary text-sm">Memuat project...</div>;
  }

  const atRisk = project.status === "ON_HOLD" || project.progress < 30;
  const pmMembers = members?.filter((m) => m.project_role === "PROJECT_MANAGER") ?? [];
  const currentProjectMember = members?.find((m) => m.user_id === user?.id);
  const canManageMembers = user?.system_role === "ADMIN" || currentProjectMember?.project_role === "PROJECT_MANAGER";

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
            <div className="mt-3 flex flex-wrap items-center gap-2">
              {project.project_code && (
                <StatusBadge color="gray">{project.project_code}</StatusBadge>
              )}
              {project.category && <StatusBadge color="gray">{project.category}</StatusBadge>}
              {project.initiation_type && <StatusBadge color="blue">{project.initiation_type}</StatusBadge>}
              {project.priority && <StatusBadge color={getStatusColor(project.priority)}>{project.priority}</StatusBadge>}
              {project.health && <StatusBadge color={getHealthColor(project.health)}>{project.health}</StatusBadge>}
            </div>
            {project.notes && (
              <p className="text-[12px] text-ink-tertiary mt-2.5 m-0 whitespace-pre-wrap">{project.notes}</p>
            )}
            <div className="mt-3">
              <Button size="sm" variant="outline" onClick={startEditProject}>
                Edit project
              </Button>
            </div>
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

      {editingProject && (
        <Card className="mb-6">
          <CardHeader>
            <CardTitle>Edit project</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 mb-4">
              <div>
                <Label>Nama project</Label>
                <Input value={projectEditForm.name} onChange={(e) => setProjectEditForm({ ...projectEditForm, name: e.target.value })} />
              </div>
              <div>
                <Label>Deskripsi</Label>
                <Input value={projectEditForm.description} onChange={(e) => setProjectEditForm({ ...projectEditForm, description: e.target.value })} />
              </div>
              <div>
                <Label>Kategori</Label>
                <Input value={projectEditForm.category} onChange={(e) => setProjectEditForm({ ...projectEditForm, category: e.target.value })} />
              </div>
              <div>
                <Label>Jenis inisiasi</Label>
                <Select
                  value={projectEditForm.initiation_type}
                  onChange={(e) => setProjectEditForm({ ...projectEditForm, initiation_type: e.target.value })}
                  options={[
                    { value: "", label: "Tidak diatur" },
                    { value: "NEW_INITIATIVE", label: "New Initiative" },
                    { value: "RENEWAL", label: "Renewal" },
                    { value: "ENHANCEMENT", label: "Enhancement" },
                  ]}
                />
              </div>
              <div>
                <Label>Priority</Label>
                <Select
                  value={projectEditForm.priority}
                  onChange={(e) => setProjectEditForm({ ...projectEditForm, priority: e.target.value as ProjectPriority })}
                  options={[
                    { value: "LOW", label: "Low" },
                    { value: "MEDIUM", label: "Medium" },
                    { value: "HIGH", label: "High" },
                    { value: "URGENT", label: "Urgent" },
                  ]}
                />
              </div>
              <div>
                <Label>Tanggal mulai</Label>
                <Input type="date" value={projectEditForm.start_date} onChange={(e) => setProjectEditForm({ ...projectEditForm, start_date: e.target.value })} />
              </div>
              <div>
                <Label>Tanggal selesai</Label>
                <Input type="date" value={projectEditForm.end_date} onChange={(e) => setProjectEditForm({ ...projectEditForm, end_date: e.target.value })} />
              </div>
              <div className="sm:col-span-2">
                <Label>Catatan</Label>
                <Input value={projectEditForm.notes} onChange={(e) => setProjectEditForm({ ...projectEditForm, notes: e.target.value })} />
              </div>
            </div>
            <div className="flex gap-2">
              <Button variant="primary" onClick={saveProjectEdit} disabled={!projectEditForm.name.trim()}>
                Simpan
              </Button>
              <Button variant="outline" onClick={() => setEditingProject(false)}>
                Batal
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

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
                      {milestones.map((m, index) => (
                        <div key={m.id} className="pb-4 border-b border-border last:border-b-0 last:pb-0">
                          {editingMilestoneId === m.id ? (
                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
                              <Input
                                value={milestoneEditForm.name}
                                onChange={(e) => setMilestoneEditForm({ ...milestoneEditForm, name: e.target.value })}
                                placeholder="Nama milestone"
                              />
                              <Input
                                value={milestoneEditForm.description}
                                onChange={(e) => setMilestoneEditForm({ ...milestoneEditForm, description: e.target.value })}
                                placeholder="Deskripsi"
                              />
                              <Input
                                type="date"
                                value={milestoneEditForm.start_date ?? ""}
                                onChange={(e) => setMilestoneEditForm({ ...milestoneEditForm, start_date: e.target.value })}
                              />
                              <Input
                                type="date"
                                value={milestoneEditForm.end_date ?? ""}
                                onChange={(e) => setMilestoneEditForm({ ...milestoneEditForm, end_date: e.target.value })}
                              />
                              <div className="flex gap-2 sm:col-span-2">
                                <Button size="sm" variant="primary" onClick={() => saveMilestoneEdit(m)}>
                                  Simpan
                                </Button>
                                <Button size="sm" variant="outline" onClick={() => setEditingMilestoneId(null)}>
                                  Batal
                                </Button>
                              </div>
                            </div>
                          ) : (
                            <div className="flex items-center gap-4">
                              <div className="flex-1 min-w-0">
                                <div className="flex items-center justify-between gap-2 mb-1.5">
                                  <span className="text-[13.5px] font-medium truncate">{m.name}</span>
                                  <StatusBadge color={getStatusColor(m.status)}>{m.status}</StatusBadge>
                                </div>
                                <HealthBar progress={m.progress} atRisk={m.status === "CANCELLED"} />
                              </div>
                              <span className="text-[13px] font-semibold text-ink-secondary w-10 text-right flex-shrink-0">
                                {m.progress.toFixed(0)}%
                              </span>
                              <select
                                value={m.status}
                                onChange={(e) => changeMilestoneStatus({ milestoneId: m.id, status: e.target.value, version: m.version })}
                                className="h-7 px-2 rounded-md border border-border-strong bg-surface text-[11.5px] font-medium flex-shrink-0 cursor-pointer"
                              >
                                <option value="PLANNED">Planned</option>
                                <option value="ACTIVE">Active</option>
                                <option value="COMPLETED">Completed</option>
                                <option value="CANCELLED">Cancelled</option>
                              </select>
                              <div className="flex gap-1 flex-shrink-0">
                                <Button size="sm" variant="ghost" onClick={() => moveMilestone(index, -1)} disabled={index === 0}>
                                  Naik
                                </Button>
                                <Button size="sm" variant="ghost" onClick={() => moveMilestone(index, 1)} disabled={index === milestones.length - 1}>
                                  Turun
                                </Button>
                                <Button size="sm" variant="outline" onClick={() => startEditMilestone(m)}>
                                  Edit
                                </Button>
                                <Button size="sm" variant="danger" onClick={() => handleDeleteMilestone(m.id)}>
                                  Hapus
                                </Button>
                              </div>
                            </div>
                          )}
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
                        <div key={t.id} className="py-3 border-b border-border last:border-b-0">
                          {editingTaskId === t.id ? (
                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
                              <Input
                                value={taskEditForm.title}
                                onChange={(e) => setTaskEditForm({ ...taskEditForm, title: e.target.value })}
                                placeholder="Judul task"
                              />
                              <Select
                                value={taskEditForm.priority}
                                onChange={(e) => setTaskEditForm({ ...taskEditForm, priority: e.target.value as TaskPriority })}
                                options={[
                                  { value: "LOW", label: "Low" },
                                  { value: "MEDIUM", label: "Medium" },
                                  { value: "HIGH", label: "High" },
                                  { value: "URGENT", label: "Urgent" },
                                ]}
                              />
                              <Input
                                value={taskEditForm.description}
                                onChange={(e) => setTaskEditForm({ ...taskEditForm, description: e.target.value })}
                                placeholder="Deskripsi"
                              />
                              <Select
                                value={taskEditForm.milestone_id}
                                onChange={(e) => setTaskEditForm({ ...taskEditForm, milestone_id: e.target.value })}
                                options={[
                                  { value: "", label: "Tanpa milestone" },
                                  ...(milestones || []).map((m) => ({ value: String(m.id), label: m.name })),
                                ]}
                              />
                              <Input
                                type="date"
                                value={taskEditForm.start_date}
                                onChange={(e) => setTaskEditForm({ ...taskEditForm, start_date: e.target.value })}
                              />
                              <Input
                                type="date"
                                value={taskEditForm.due_date}
                                onChange={(e) => setTaskEditForm({ ...taskEditForm, due_date: e.target.value })}
                              />
                              <div className="flex gap-2 sm:col-span-2">
                                <Button size="sm" variant="primary" onClick={() => saveTaskEdit(t)}>
                                  Simpan
                                </Button>
                                <Button size="sm" variant="outline" onClick={() => setEditingTaskId(null)}>
                                  Batal
                                </Button>
                              </div>
                            </div>
                          ) : (
                            <div className="flex items-center gap-3">
                              <span
                                className="h-2 w-2 rounded-sm flex-shrink-0"
                                style={{ background: getPriorityColor(t.priority) }}
                              />
                              <div className="flex-1 min-w-0">
                                <p className="text-[13px] font-medium m-0 truncate">{t.title}</p>
                                <p className="text-[11.5px] text-ink-tertiary m-0">
                                  {t.assignee_ids?.length ? `Assignee: ${t.assignee_ids.join(", ")}` : "Belum ada assignee"}
                                </p>
                              </div>

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

                              <div className="flex gap-1 flex-shrink-0">
                                <Button size="sm" variant="ghost" onClick={() => setCommentingTaskId(commentingTaskId === t.id ? null : t.id)}>
                                  Komentar
                                </Button>
                                <Button size="sm" variant="outline" onClick={() => startAssignTask(t)}>
                                  Assign
                                </Button>
                                <Button size="sm" variant="outline" onClick={() => startEditTask(t)}>
                                  Edit
                                </Button>
                                <Button size="sm" variant="danger" onClick={() => handleDeleteTask(t.id)}>
                                  Hapus
                                </Button>
                              </div>
                            </div>
                          )}

                          {assigningTaskId === t.id && (
                            <div className="mt-3 flex gap-2">
                              <Input
                                value={assigneeInput}
                                onChange={(e) => setAssigneeInput(e.target.value)}
                                placeholder="User ID, pisahkan dengan koma"
                              />
                              <Button size="sm" variant="primary" onClick={() => saveTaskAssignees(t.id)}>
                                Simpan
                              </Button>
                              <Button size="sm" variant="outline" onClick={() => setAssigningTaskId(null)}>
                                Batal
                              </Button>
                            </div>
                          )}

                          {commentingTaskId === t.id && (
                            <div className="mt-3 rounded-md border border-border bg-surface-secondary p-3">
                              <div className="space-y-2 mb-3">
                                {!taskComments || taskComments.data.length === 0 ? (
                                  <p className="text-[12px] text-ink-tertiary m-0">Belum ada komentar.</p>
                                ) : (
                                  taskComments.data.map((comment) => (
                                    <div key={comment.id} className="text-[12.5px]">
                                      <p className="m-0">{comment.comment}</p>
                                      <p className="m-0 text-[11px] text-ink-tertiary">
                                        User #{comment.user_id} - {new Date(comment.created_at).toLocaleString("id-ID")}
                                      </p>
                                    </div>
                                  ))
                                )}
                              </div>
                              <div className="flex gap-2">
                                <Input
                                  value={commentInput}
                                  onChange={(e) => setCommentInput(e.target.value)}
                                  placeholder="Tulis komentar..."
                                />
                                <Button size="sm" variant="primary" onClick={() => saveTaskComment(t.id)} disabled={!commentInput.trim()}>
                                  Kirim
                                </Button>
                              </div>
                            </div>
                          )}
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
                  {canManageMembers && (
                    <div className="grid grid-cols-1 sm:grid-cols-[1fr_180px_auto] gap-2.5 mb-5">
                      <div>
                        <Label>User</Label>
                        <Select
                          value={memberForm.userId ? String(memberForm.userId) : ""}
                          onChange={(e) =>
                            setMemberForm({
                              ...memberForm,
                              userId: e.target.value ? Number(e.target.value) : null,
                            })
                          }
                          placeholder="Pilih user"
                          options={[
                            { value: "", label: "Pilih user" },
                            ...users
                              .filter((u) => u.is_active !== false)
                              .map((u) => ({ value: String(u.id), label: u.full_name })),
                          ]}
                        />
                      </div>
                      <div>
                        <Label>Project role</Label>
                        <Select
                          value={memberForm.project_role}
                          onChange={(e) =>
                            setMemberForm({ ...memberForm, project_role: e.target.value as ProjectRole })
                          }
                          options={[
                            { value: "PROJECT_MANAGER", label: "Project Manager" },
                            { value: "MEMBER", label: "Member" },
                            { value: "OBSERVER", label: "Observer" },
                          ]}
                        />
                      </div>
                      <div className="flex items-end">
                        <Button
                          variant="primary"
                          onClick={handleAddMember}
                          disabled={addingMember || !memberForm.userId}
                        >
                          Tambah member
                        </Button>
                      </div>
                    </div>
                  )}

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

      <ConfirmDeleteDialog
        open={milestoneDeleteOpen}
        title="Hapus milestone"
        description="Apakah Anda yakin ingin menghapus milestone ini?"
        onConfirm={confirmDeleteMilestone}
        onCancel={() => {
          setMilestoneDeleteOpen(false);
          setMilestoneDeleteId(null);
        }}
        isDeleting={deleteMilestoneMutation.isPending}
      />
      <ConfirmDeleteDialog
        open={taskDeleteOpen}
        title="Hapus task"
        description="Apakah Anda yakin ingin menghapus task ini?"
        onConfirm={confirmDeleteTask}
        onCancel={() => {
          setTaskDeleteOpen(false);
          setTaskDeleteId(null);
        }}
        isDeleting={deleteTaskMutation.isPending}
      />
    </div>
  );
}
