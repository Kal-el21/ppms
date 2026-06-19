import { useParams } from "react-router-dom";
import { useState } from "react";
import {
  useRequestDetail,
  useRevisionHistory,
  useApprovalHistory,
  useSubmitRequest,
  useReviewRequest,
} from "../hooks/useRequests";
import { useAuth } from "../../auth/context/AuthContext";
import { Button } from "../../../components/ui/button";
import { Badge } from "../../../components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "../../../components/ui/card";
import { Input } from "../../../components/ui/input";

export default function RequestDetailPage() {
  const { id } = useParams();
  const requestId = Number(id);
  const { user } = useAuth();

  const { data: request, isLoading } = useRequestDetail(requestId);
  const { data: revisions } = useRevisionHistory(requestId);
  const { data: approvals } = useApprovalHistory(requestId);

  const { mutate: submitRequest, isPending: submitting } = useSubmitRequest();
  const { mutate: review, isPending: reviewing } = useReviewRequest(requestId);

  const [comment, setComment] = useState("");

  if (isLoading || !request) return <div>Loading...</div>;

  const isOwner = user?.id === request.requester_id;
  const isAdmin = user?.system_role === "ADMIN";
  const canSubmit = isOwner && (request.status === "DRAFT" || request.status === "REVISED");
  const canReview = isAdmin && request.status === "UNDER_REVIEW";

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

      {canSubmit && (
        <Button onClick={() => submitRequest(requestId)} disabled={submitting}>
          Submit for Review
        </Button>
      )}

      {canReview && (
        <Card>
          <CardHeader>
            <CardTitle>Review Request</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <Input
              placeholder="Comment (optional)"
              value={comment}
              onChange={(e) => setComment(e.target.value)}
            />
            <div className="flex gap-2">
              <Button
                onClick={() => review({ action: "APPROVED", comment })}
                disabled={reviewing}
              >
                Approve
              </Button>
              <Button
                variant="destructive"
                onClick={() => review({ action: "REJECTED", comment })}
                disabled={reviewing}
              >
                Reject
              </Button>
              <Button
                variant="outline"
                onClick={() => review({ action: "REQUEST_REVISION", comment })}
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