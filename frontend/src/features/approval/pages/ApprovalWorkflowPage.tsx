import { useEffect, useState } from "react";
import {
  useApprovalLevels,
  useApprovalWorkflows,
  useCreateApprovalLevel,
  useCreateApprovalWorkflow,
} from "../hooks/useApprovalWorkflows";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select } from "@/components/ui/select";
import { StatusBadge } from "@/components/ui/status-badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { PageHeader } from "@/components/shared/PageHeader";
import { BulkActionsDropdown } from "@/components/shared/BulkActionsDropdown";
import { EmptyState } from "@/components/shared/EmptyState";
import { TableSkeleton } from "@/components/ui/skeleton";
import { useTableSelection } from "@/components/shared/useTableSelection";

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

  const { data: levels, isLoading: loadingLevels } = useApprovalLevels(selectedWorkflowId);
  const { mutate: createLevel, isPending: creatingLevel } = useCreateApprovalLevel(selectedWorkflowId);

  const workflowsList = workflows ?? [];
  const levelsList = levels ?? [];

  const workflowSelection = useTableSelection<number>();
  const levelSelection = useTableSelection<number>();

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

  const [levelOrder, setLevelOrder] = useState("");
  const [roleRequired, setRoleRequired] = useState("ADMIN");

  return (
    <div className="space-y-5">
      <PageHeader
        title="Approval Workflows"
        subtitle="Kelola workflow dan level approval untuk pengembangan multi-level approval."
      />

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
          <CardHeader className="flex items-center justify-between">
            <CardTitle>Workflow</CardTitle>
            <BulkActionsDropdown selectedCount={workflowSelection.selectedCount} />
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <TableSkeleton rows={3} cols={4} />
            ) : !workflows || workflows.length === 0 ? (
              <EmptyState icon={emptyIcon} title="Belum ada workflow" />
            ) : (
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead style={{ width: 40 }}>
                        <input
                          type="checkbox"
                          checked={workflowSelection.isAllSelected(workflowsList.map((w) => w.id))}
                          ref={(el) => {
                            if (el) el.indeterminate = workflowSelection.isIndeterminate(workflowsList.map((w) => w.id));
                          }}
                          onChange={(e) => workflowSelection.toggleAll(workflowsList.map((w) => w.id), e.target.checked)}
                        />
                      </TableHead>
                      <TableHead style={{ width: 50 }}>No</TableHead>
                      <TableHead>Name</TableHead>
                      <TableHead>Status</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {workflows.map((workflow, index) => (
                      <TableRow
                        key={workflow.id}
                        className="cursor-pointer"
                        onClick={() => setSelectedWorkflowId(workflow.id)}
                      >
                        <TableCell>
                          <input
                            type="checkbox"
                            checked={workflowSelection.selectedIds.has(workflow.id)}
                            onChange={() => workflowSelection.toggle(workflow.id)}
                          />
                        </TableCell>
                        <TableCell className="text-ink-secondary text-[13px]">{index + 1}</TableCell>
                        <TableCell className="font-medium text-[13px]">{workflow.name}</TableCell>
                        <TableCell>
                          <StatusBadge color={workflow.is_active ? "green" : "gray"}>
                            {workflow.is_active ? "ACTIVE" : "INACTIVE"}
                          </StatusBadge>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
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
              <TableSkeleton rows={3} cols={3} />
            ) : !selectedWorkflowId ? (
              <EmptyState icon={emptyIcon} title="Pilih workflow" />
            ) : !levels || levels.length === 0 ? (
              <EmptyState icon={emptyIcon} title="Belum ada level" />
            ) : (
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead style={{ width: 40 }}>
                        <input
                          type="checkbox"
                          checked={levelSelection.isAllSelected(levelsList.map((l) => l.id))}
                          ref={(el) => {
                            if (el) el.indeterminate = levelSelection.isIndeterminate(levelsList.map((l) => l.id));
                          }}
                          onChange={(e) => levelSelection.toggleAll(levelsList.map((l) => l.id), e.target.checked)}
                        />
                      </TableHead>
                      <TableHead style={{ width: 50 }}>No</TableHead>
                      <TableHead>Level Order</TableHead>
                      <TableHead>Role Required</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {levels.map((level, index) => (
                      <TableRow key={level.id}>
                        <TableCell>
                          <input
                            type="checkbox"
                            checked={levelSelection.selectedIds.has(level.id)}
                            onChange={() => levelSelection.toggle(level.id)}
                          />
                        </TableCell>
                        <TableCell className="text-ink-secondary text-[13px]">{index + 1}</TableCell>
                        <TableCell className="font-medium text-[13px]">Level {level.level_order}</TableCell>
                        <TableCell>
                          <StatusBadge color="blue">{level.role_required}</StatusBadge>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
