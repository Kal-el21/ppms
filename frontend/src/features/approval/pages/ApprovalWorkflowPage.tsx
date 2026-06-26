import { useEffect, useState } from "react";
import {
  useApprovalLevels,
  useApprovalWorkflows,
  useCreateApprovalLevel,
  useCreateApprovalWorkflow,
} from "../hooks/useApprovalWorkflows";
import { Button } from "../../../components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "../../../components/ui/card";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import { Select } from "../../../components/ui/select";
import { EmptyState } from "../../../components/shared/EmptyState";
import { StatusBadge } from "../../../components/ui/status-badge";

const emptyIcon = (
  <svg className="h-5 w-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
    <path d="M9 12l2 2 4-4" />
    <path d="M20 6L9 17l-5-5" />
  </svg>
);

export default function ApprovalWorkflowPage() {
  const { data: workflows, isLoading } = useApprovalWorkflows();
  const { mutate: createWorkflow, isPending: creatingWorkflow } = useCreateApprovalWorkflow();

  const [selectedWorkflowId, setSelectedWorkflowId] = useState(0);
  const [workflowName, setWorkflowName] = useState("");
  const [levelOrder, setLevelOrder] = useState("");
  const [roleRequired, setRoleRequired] = useState("ADMIN");

  const { data: levels, isLoading: loadingLevels } = useApprovalLevels(selectedWorkflowId);
  const { mutate: createLevel, isPending: creatingLevel } = useCreateApprovalLevel(selectedWorkflowId);

  useEffect(() => {
    if (!selectedWorkflowId && workflows && workflows.length > 0) {
      setSelectedWorkflowId(workflows[0].id);
    }
  }, [selectedWorkflowId, workflows]);

  const handleCreateWorkflow = () => {
    if (!workflowName.trim()) return;
    createWorkflow(workflowName, {
      onSuccess: (workflow) => {
        setWorkflowName("");
        setSelectedWorkflowId(workflow.id);
      },
    });
  };

  const handleCreateLevel = () => {
    if (!selectedWorkflowId || !levelOrder.trim() || !roleRequired.trim()) return;
    createLevel(
      {
        level_order: Number(levelOrder),
        role_required: roleRequired,
      },
      {
        onSuccess: () => {
          setLevelOrder("");
          setRoleRequired("ADMIN");
        },
      }
    );
  };

  return (
    <div className="space-y-5">
      <div>
        <h1 className="text-[20px] font-semibold tracking-tight m-0">Approval Workflows</h1>
        <p className="text-[13px] text-ink-secondary mt-1 m-0">
          Kelola workflow dan level approval untuk pengembangan multi-level approval.
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Buat workflow</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex gap-2">
            <Input
              placeholder="cth. Default Project Request Approval"
              value={workflowName}
              onChange={(e) => setWorkflowName(e.target.value)}
            />
            <Button variant="primary" onClick={handleCreateWorkflow} disabled={creatingWorkflow || !workflowName.trim()}>
              Simpan
            </Button>
          </div>
        </CardContent>
      </Card>

      <div className="grid grid-cols-1 lg:grid-cols-[0.9fr_1.3fr] gap-4">
        <Card>
          <CardHeader>
            <CardTitle>Workflow</CardTitle>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <p className="text-[12.5px] text-ink-tertiary">Memuat workflow...</p>
            ) : !workflows || workflows.length === 0 ? (
              <EmptyState icon={emptyIcon} title="Belum ada workflow" />
            ) : (
              <div className="flex flex-col">
                {workflows.map((workflow) => (
                  <button
                    key={workflow.id}
                    type="button"
                    onClick={() => setSelectedWorkflowId(workflow.id)}
                    className={`flex items-center justify-between gap-3 py-3 px-2 rounded-md text-left ${
                      selectedWorkflowId === workflow.id ? "bg-primary-50 text-primary-700" : "hover:bg-surface-secondary"
                    }`}
                  >
                    <span className="text-[13px] font-medium">{workflow.name}</span>
                    <StatusBadge color={workflow.is_active ? "green" : "gray"}>
                      {workflow.is_active ? "ACTIVE" : "INACTIVE"}
                    </StatusBadge>
                  </button>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Approval Levels</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 sm:grid-cols-[120px_1fr_auto] gap-2.5 mb-5">
              <div>
                <Label>Urutan</Label>
                <Input type="number" min={1} value={levelOrder} onChange={(e) => setLevelOrder(e.target.value)} />
              </div>
              <div>
                <Label>Role</Label>
                <Select
                  value={roleRequired}
                  onChange={(e) => setRoleRequired(e.target.value)}
                  options={[
                    { value: "ADMIN", label: "Admin" },
                    { value: "USER", label: "User" },
                    { value: "VIEWER", label: "Viewer" },
                  ]}
                />
              </div>
              <div className="flex items-end">
                <Button
                  variant="primary"
                  onClick={handleCreateLevel}
                  disabled={creatingLevel || !selectedWorkflowId || !levelOrder.trim()}
                >
                  Tambah
                </Button>
              </div>
            </div>

            {loadingLevels ? (
              <p className="text-[12.5px] text-ink-tertiary">Memuat levels...</p>
            ) : !selectedWorkflowId ? (
              <EmptyState icon={emptyIcon} title="Pilih workflow" />
            ) : !levels || levels.length === 0 ? (
              <EmptyState icon={emptyIcon} title="Belum ada level" />
            ) : (
              <div className="flex flex-col">
                {levels.map((level) => (
                  <div key={level.id} className="flex items-center justify-between gap-3 py-3 border-b border-border last:border-b-0">
                    <div>
                      <p className="text-[13px] font-medium m-0">Level {level.level_order}</p>
                      <p className="text-[11.5px] text-ink-tertiary m-0">Workflow #{level.workflow_id}</p>
                    </div>
                    <StatusBadge color="blue">{level.role_required}</StatusBadge>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
