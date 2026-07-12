import { useNavigate, useParams } from "react-router-dom";
import { useState, type ReactNode } from "react";
import { useUsers } from "../../users/hooks/useUsers";
import {
  useRequestDetail,
  useRevisionHistory,
  useApprovalHistory,
  useSubmitRequest,
  useReviewRequest,
  useReviseRequest,
  useDeleteDraft,
} from "../hooks/useRequests";
import { useAuth } from "../../auth/context/AuthContext";
import { Button } from "../../../components/ui/button";
import { Badge } from "../../../components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "../../../components/ui/card";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import { Select } from "../../../components/ui/select";
import { Textarea } from "../../../components/ui/textarea";
import { PageHeader } from "../../../components/shared/PageHeader";
import type { ProjectRequest, RequestRevision } from "../types";

const initiationOptions = [
  { value: "NEW_INITIATIVE", label: "New Initiative" },
  { value: "RENEWAL", label: "Renewal" },
  { value: "ENHANCEMENT", label: "Enhancement" },
];

const priorityOptions = [
  { value: "LOW", label: "Low" },
  { value: "MEDIUM", label: "Medium" },
  { value: "HIGH", label: "High" },
  { value: "URGENT", label: "Urgent" },
];

const budgetTypeOptions = [
  { value: "CAPEX", label: "CAPEX" },
  { value: "OPEX", label: "OPEX" },
];

type ReviseFormState = {
  title: string;
  description: string;
  business_goal: string;
  expected_outcome: string;
  category: string;
  initiation_type: string;
  priority: string;
  proposed_start_date: string;
  proposed_end_date: string;
  budget_type: string;
  budget_name: string;
  estimated_budget: string;
  notes: string;
};

const emptyReviseForm: ReviseFormState = {
  title: "",
  description: "",
  business_goal: "",
  expected_outcome: "",
  category: "",
  initiation_type: "",
  priority: "MEDIUM",
  proposed_start_date: "",
  proposed_end_date: "",
  budget_type: "",
  budget_name: "",
  estimated_budget: "0",
  notes: "",
};

function formatCurrency(value?: number) {
  return `Rp ${(value ?? 0).toLocaleString("id-ID")}`;
}

function formatDate(value?: string | null) {
  return value ? new Date(value).toLocaleDateString("id-ID") : "-";
}

function toDateInput(value?: string | null) {
  return value ? value.slice(0, 10) : "";
}

function formatEnum(value?: string | null) {
  if (!value) return "-";
  return value
    .split("_")
    .map((part) => part.charAt(0) + part.slice(1).toLowerCase())
    .join(" ");
}

function DetailItem({ label, value }: { label: string; value?: ReactNode }) {
  return (
    <div>
      <p className="mb-1 text-[12px] font-medium uppercase text-ink-tertiary">{label}</p>
      <div className="text-[13.5px] text-ink-primary">{value || "-"}</div>
    </div>
  );
}

function buildReviseForm(request: ProjectRequest): ReviseFormState {
  return {
    title: request.title,
    description: request.description || "",
    business_goal: request.business_goal || "",
    expected_outcome: request.expected_outcome || "",
    category: request.category || "",
    initiation_type: request.initiation_type || "",
    priority: request.priority || "MEDIUM",
    proposed_start_date: toDateInput(request.proposed_start_date),
    proposed_end_date: toDateInput(request.proposed_end_date),
    budget_type: request.budget_type || "",
    budget_name: request.budget_name || "",
    estimated_budget: String(request.estimated_budget ?? 0),
    notes: request.notes || "",
  };
}

function RevisionSnapshot({ revision }: { revision: RequestRevision }) {
  return (
    <div className="border-b border-border pb-3 text-sm last:border-b-0 last:pb-0">
      <div className="flex flex-wrap items-center gap-2">
        <span className="font-medium text-ink-primary">Revision #{revision.revision_number}</span>
        <span className="text-ink-tertiary">{new Date(revision.created_at).toLocaleString("id-ID")}</span>
      </div>
      <p className="mt-1 text-ink-secondary">Reason: {revision.revision_reason}</p>
      <div className="mt-3 grid gap-3 md:grid-cols-4">
        <DetailItem label="Title" value={revision.title} />
        <DetailItem label="Priority" value={formatEnum(revision.priority)} />
        <DetailItem label="Budget" value={formatCurrency(revision.estimated_budget)} />
        <DetailItem label="End Date" value={formatDate(revision.proposed_end_date)} />
      </div>
    </div>
  );
}

export default function RequestDetailPage() {
  const { id } = useParams();
  const requestId = Number(id);
  const navigate = useNavigate();
  const { user } = useAuth();
  const isAdmin = user?.system_role === "ADMIN";

  const { data: request, isLoading } = useRequestDetail(requestId);
  const { data: revisions } = useRevisionHistory(requestId);
  const { data: approvals } = useApprovalHistory(requestId);
  const { data: usersData } = useUsers(1, 100, isAdmin);

  const { mutate: submitRequest, isPending: submitting } = useSubmitRequest();
  const { mutate: review, isPending: reviewing } = useReviewRequest(requestId);
  const { mutate: reviseRequest, isPending: revising } = useReviseRequest(requestId);
  const { mutate: deleteDraft, isPending: deletingDraft } = useDeleteDraft();

  const [comment, setComment] = useState("");
  const [projectManagerId, setProjectManagerId] = useState("");
  const [showReviseForm, setShowReviseForm] = useState(false);
  const [revisionReason, setRevisionReason] = useState("");
  const [reviseForm, setReviseForm] = useState<ReviseFormState>(emptyReviseForm);

  if (isLoading || !request) return <div>Loading...</div>;

  const isOwner = user?.id === request.requester_id;
  const latestApproval = approvals?.[approvals.length - 1];
  const canSubmit = isOwner && (request.status === "DRAFT" || request.status === "REVISED");
  const canDeleteDraft = isOwner && request.status === "DRAFT";
  const canRevise =
    isOwner &&
    (request.status === "REVISION_REQUESTED" ||
      (request.status === "REJECTED" && latestApproval?.action === "REQUEST_REVISION"));
  const canReview = isAdmin && (request.status === "UNDER_REVIEW" || request.status === "REVISED");
  const projectManagerOptions = (usersData?.data || [])
    .filter((candidate) => candidate.is_active !== false && candidate.system_role !== "VIEWER")
    .map((candidate) => ({
      value: String(candidate.id),
      label: `${candidate.full_name} (${candidate.email})`,
    }));
  const userNameById = new Map((usersData?.data || []).map((candidate) => [candidate.id, candidate.full_name]));
  const revisionDatesInvalid =
    !!reviseForm.proposed_start_date &&
    !!reviseForm.proposed_end_date &&
    reviseForm.proposed_end_date < reviseForm.proposed_start_date;
  const revisionReady =
    reviseForm.title.trim().length >= 5 &&
    reviseForm.initiation_type &&
    reviseForm.proposed_start_date &&
    reviseForm.proposed_end_date &&
    reviseForm.budget_type &&
    reviseForm.budget_name.trim().length >= 2 &&
    revisionReason.trim().length >= 5 &&
    !revisionDatesInvalid;

  const openReviseForm = () => {
    setReviseForm(buildReviseForm(request));
    setRevisionReason("");
    setShowReviseForm(true);
  };

  const handleRevise = () => {
    reviseRequest(
      {
        title: reviseForm.title,
        description: reviseForm.description,
        business_goal: reviseForm.business_goal,
        expected_outcome: reviseForm.expected_outcome,
        category: reviseForm.category,
        initiation_type: reviseForm.initiation_type,
        priority: reviseForm.priority,
        proposed_start_date: reviseForm.proposed_start_date,
        proposed_end_date: reviseForm.proposed_end_date,
        budget_type: reviseForm.budget_type,
        budget_name: reviseForm.budget_name,
        estimated_budget: Number(reviseForm.estimated_budget || 0),
        notes: reviseForm.notes,
        revision_reason: revisionReason,
      },
      {
        onSuccess: () => setShowReviseForm(false),
      }
    );
  };

  const handleDeleteDraft = () => {
    if (!confirm("Hapus draft request ini?")) return;
    deleteDraft(requestId, {
      onSuccess: () => navigate("/project-requests"),
    });
  };

  const handleReview = (action: "APPROVED" | "REJECTED" | "REQUEST_REVISION") => {
    review(
      {
        action,
        comment,
        project_manager_id: action === "APPROVED" ? Number(projectManagerId) : undefined,
      },
      {
        onSuccess: () => {
          setComment("");
          setProjectManagerId("");
        },
      }
    );
  };

  return (
    <div className="space-y-6">
      <PageHeader
        title={request.title}
        subtitle={request.request_number || `REQ-${request.id}`}
        actions={<Badge>{request.status}</Badge>}
      />

      <Card>
        <CardHeader>
          <CardTitle>Info Project</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-4 md:grid-cols-3">
          <DetailItem label="Category" value={request.category} />
          <DetailItem label="Initiation" value={formatEnum(request.initiation_type)} />
          <DetailItem label="Priority" value={formatEnum(request.priority)} />
          <div className="md:col-span-3">
            <DetailItem label="Description" value={request.description} />
          </div>
          <DetailItem label="Business Goal" value={request.business_goal} />
          <DetailItem label="Expected Outcome" value={request.expected_outcome} />
          <DetailItem label="Requester" value={`User #${request.requester_id}`} />
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Anggaran dan Timeline</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-4 md:grid-cols-3">
          <DetailItem label="Budget Type" value={request.budget_type || "-"} />
          <DetailItem label="Budget Name" value={request.budget_name} />
          <DetailItem label="Estimated Budget" value={formatCurrency(request.estimated_budget)} />
          <DetailItem label="Start Date" value={formatDate(request.proposed_start_date)} />
          <DetailItem label="End Date" value={formatDate(request.proposed_end_date)} />
          <DetailItem label="Current Revision" value={String(request.current_revision ?? 0)} />
          <div className="md:col-span-3">
            <DetailItem label="Notes" value={request.notes} />
          </div>
        </CardContent>
      </Card>

      {(canSubmit || canDeleteDraft || canRevise) && (
        <div className="flex flex-wrap gap-2">
          {canSubmit && (
            <Button onClick={() => submitRequest(requestId)} disabled={submitting}>
              {submitting ? "Submitting..." : request.status === "REVISED" ? "Submit Revision" : "Submit for Review"}
            </Button>
          )}
          {canRevise && (
            <Button variant="primary" onClick={openReviseForm} disabled={revising}>
              Revise Request
            </Button>
          )}
          {canDeleteDraft && (
            <Button variant="danger" onClick={handleDeleteDraft} disabled={deletingDraft}>
              {deletingDraft ? "Deleting..." : "Delete Draft"}
            </Button>
          )}
        </div>
      )}

      {showReviseForm && (
        <Card>
          <CardHeader>
            <CardTitle>Revise Request</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 md:grid-cols-2">
              <div className="md:col-span-2">
                <Label>Judul</Label>
                <Input
                  value={reviseForm.title}
                  onChange={(e) => setReviseForm({ ...reviseForm, title: e.target.value })}
                />
              </div>
              <div>
                <Label>Kategori / Tag</Label>
                <Input
                  value={reviseForm.category}
                  onChange={(e) => setReviseForm({ ...reviseForm, category: e.target.value })}
                />
              </div>
              <div>
                <Label>Jenis Inisiasi</Label>
                <Select
                  value={reviseForm.initiation_type}
                  onChange={(e) => setReviseForm({ ...reviseForm, initiation_type: e.target.value })}
                  placeholder="Pilih jenis inisiasi"
                  options={initiationOptions}
                />
              </div>
              <div>
                <Label>Prioritas</Label>
                <Select
                  value={reviseForm.priority}
                  onChange={(e) => setReviseForm({ ...reviseForm, priority: e.target.value })}
                  options={priorityOptions}
                />
              </div>
              <div>
                <Label>Jenis Anggaran</Label>
                <Select
                  value={reviseForm.budget_type}
                  onChange={(e) => setReviseForm({ ...reviseForm, budget_type: e.target.value })}
                  placeholder="Pilih jenis anggaran"
                  options={budgetTypeOptions}
                />
              </div>
              <div>
                <Label>Nama Mata Anggaran</Label>
                <Input
                  value={reviseForm.budget_name}
                  onChange={(e) => setReviseForm({ ...reviseForm, budget_name: e.target.value })}
                />
              </div>
              <div>
                <Label>Estimasi Anggaran</Label>
                <Input
                  type="number"
                  value={reviseForm.estimated_budget}
                  onChange={(e) => setReviseForm({ ...reviseForm, estimated_budget: e.target.value })}
                />
              </div>
              <div>
                <Label>Tanggal Mulai</Label>
                <Input
                  type="date"
                  value={reviseForm.proposed_start_date}
                  onChange={(e) => setReviseForm({ ...reviseForm, proposed_start_date: e.target.value })}
                />
              </div>
              <div>
                <Label>Tanggal Selesai</Label>
                <Input
                  type="date"
                  value={reviseForm.proposed_end_date}
                  onChange={(e) => setReviseForm({ ...reviseForm, proposed_end_date: e.target.value })}
                />
                {revisionDatesInvalid && (
                  <p className="mt-1.5 text-xs text-danger-600">
                    Tanggal selesai tidak boleh sebelum tanggal mulai
                  </p>
                )}
              </div>
              <div className="md:col-span-2">
                <Label>Deskripsi</Label>
                <Textarea
                  value={reviseForm.description}
                  onChange={(e) => setReviseForm({ ...reviseForm, description: e.target.value })}
                />
              </div>
              <div>
                <Label>Tujuan Bisnis</Label>
                <Textarea
                  value={reviseForm.business_goal}
                  onChange={(e) => setReviseForm({ ...reviseForm, business_goal: e.target.value })}
                />
              </div>
              <div>
                <Label>Hasil yang Diharapkan</Label>
                <Textarea
                  value={reviseForm.expected_outcome}
                  onChange={(e) => setReviseForm({ ...reviseForm, expected_outcome: e.target.value })}
                />
              </div>
              <div className="md:col-span-2">
                <Label>Catatan</Label>
                <Textarea
                  value={reviseForm.notes}
                  onChange={(e) => setReviseForm({ ...reviseForm, notes: e.target.value })}
                />
              </div>
              <div className="md:col-span-2">
                <Label>Alasan Revisi</Label>
                <Textarea value={revisionReason} onChange={(e) => setRevisionReason(e.target.value)} />
              </div>
            </div>
            <div className="flex gap-2">
              <Button variant="primary" onClick={handleRevise} disabled={revising || !revisionReady}>
                {revising ? "Saving..." : "Save Revision"}
              </Button>
              <Button variant="outline" onClick={() => setShowReviseForm(false)} disabled={revising}>
                Cancel
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {canReview && (
        <Card>
          <CardHeader>
            <CardTitle>Review Request</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div>
              <Label>Project Manager</Label>
              <Select
                value={projectManagerId}
                onChange={(e) => setProjectManagerId(e.target.value)}
                placeholder="Pilih PM"
                options={projectManagerOptions}
                disabled={reviewing || projectManagerOptions.length === 0}
              />
            </div>
            <Input placeholder="Comment (optional)" value={comment} onChange={(e) => setComment(e.target.value)} />
            <div className="flex flex-wrap gap-2">
              <Button onClick={() => handleReview("APPROVED")} disabled={reviewing || !projectManagerId}>
                Approve
              </Button>
              <Button variant="danger" onClick={() => handleReview("REJECTED")} disabled={reviewing}>
                Reject
              </Button>
              <Button variant="outline" onClick={() => handleReview("REQUEST_REVISION")} disabled={reviewing}>
                Request Revision
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {approvals && approvals.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Approval History</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {approvals.map((approval) => (
              <div key={approval.id} className="border-b border-border pb-3 text-sm last:border-b-0 last:pb-0">
                <div className="flex flex-wrap items-center gap-2">
                  <Badge variant="outline">{approval.action}</Badge>
                  <span className="text-ink-tertiary">
                    {new Date(approval.created_at).toLocaleString("id-ID")}
                  </span>
                </div>
                {approval.project_manager_id && (
                  <p className="mt-1 text-ink-secondary">
                    PM: {userNameById.get(approval.project_manager_id) || `User #${approval.project_manager_id}`}
                  </p>
                )}
                {approval.comment && <p className="mt-1 text-ink-secondary">"{approval.comment}"</p>}
              </div>
            ))}
          </CardContent>
        </Card>
      )}

      {revisions && revisions.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Revision History</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            {revisions.map((revision) => (
              <RevisionSnapshot key={revision.id} revision={revision} />
            ))}
          </CardContent>
        </Card>
      )}
    </div>
  );
}
