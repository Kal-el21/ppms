import { useNavigate, useParams } from "react-router-dom";
import { useState } from "react";
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
  const [reviseForm, setReviseForm] = useState({
    title: "",
    description: "",
    business_goal: "",
    expected_outcome: "",
    estimated_budget: "",
  });

  if (isLoading || !request) return <div>Loading...</div>;

  const isOwner = user?.id === request.requester_id;
  const latestApproval = approvals?.[approvals.length - 1];
  const canSubmit = isOwner && request.status === "DRAFT";
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

  const openReviseForm = () => {
    setReviseForm({
      title: request.title,
      description: request.description || "",
      business_goal: request.business_goal || "",
      expected_outcome: request.expected_outcome || "",
      estimated_budget: String(request.estimated_budget ?? 0),
    });
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
        estimated_budget: Number(reviseForm.estimated_budget || 0),
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
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">{request.title}</h1>
        <Badge>{request.status}</Badge>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Details</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2 text-sm">
          <p><strong>Description:</strong> {request.description}</p>
          <p><strong>Business Goal:</strong> {request.business_goal}</p>
          <p><strong>Expected Outcome:</strong> {request.expected_outcome}</p>
          <p><strong>Estimated Budget:</strong> Rp {request.estimated_budget.toLocaleString("id-ID")}</p>
        </CardContent>
      </Card>

      {(canSubmit || canDeleteDraft || canRevise) && (
        <div className="flex flex-wrap gap-2">
          {canSubmit && (
            <Button onClick={() => submitRequest(requestId)} disabled={submitting}>
              {submitting ? "Submitting..." : "Submit for Review"}
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
            <div>
              <Label>Judul</Label>
              <Input
                value={reviseForm.title}
                onChange={(e) => setReviseForm({ ...reviseForm, title: e.target.value })}
              />
            </div>
            <div>
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
            <div>
              <Label>Estimasi Anggaran</Label>
              <Input
                type="number"
                value={reviseForm.estimated_budget}
                onChange={(e) => setReviseForm({ ...reviseForm, estimated_budget: e.target.value })}
              />
            </div>
            <div>
              <Label>Alasan Revisi</Label>
              <Textarea
                value={revisionReason}
                onChange={(e) => setRevisionReason(e.target.value)}
              />
            </div>
            <div className="flex gap-2">
              <Button
                variant="primary"
                onClick={handleRevise}
                disabled={revising || !reviseForm.title.trim() || revisionReason.trim().length < 5}
              >
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
            <Input
              placeholder="Comment (optional)"
              value={comment}
              onChange={(e) => setComment(e.target.value)}
            />
            <div className="flex gap-2">
              <Button
                onClick={() => handleReview("APPROVED")}
                disabled={reviewing || !projectManagerId}
              >
                Approve
              </Button>
              <Button
                variant="danger"
                onClick={() => handleReview("REJECTED")}
                disabled={reviewing}
              >
                Reject
              </Button>
              <Button
                variant="outline"
                onClick={() => handleReview("REQUEST_REVISION")}
                disabled={reviewing}
              >
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
          <CardContent className="space-y-2">
            {approvals.map((a) => (
              <div key={a.id} className="border-b pb-2 text-sm">
                <p><Badge variant="outline">{a.action}</Badge> — {new Date(a.created_at).toLocaleString()}</p>
                {a.project_manager_id && (
                  <p className="text-slate-500">PM: {userNameById.get(a.project_manager_id) || `User #${a.project_manager_id}`}</p>
                )}
                {a.comment && <p className="text-slate-500">"{a.comment}"</p>}
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
          <CardContent className="space-y-2">
            {revisions.map((r) => (
              <div key={r.id} className="border-b pb-2 text-sm">
                <p>Revision #{r.revision_number} — {new Date(r.created_at).toLocaleString()}</p>
                <p className="text-slate-500">Reason: {r.revision_reason}</p>
              </div>
            ))}
          </CardContent>
        </Card>
      )}
    </div>
  );
}
