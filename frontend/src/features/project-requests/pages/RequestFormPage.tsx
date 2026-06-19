import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useNavigate, useParams } from "react-router-dom";
import { useCreateDraft, useRequestDetail, useUpdateDraft } from "../hooks/useRequests";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "../../../components/ui/card";

const schema = z.object({
  title: z.string().min(5, "Minimal 5 karakter"),
  description: z.string().optional(),
  business_goal: z.string().optional(),
  expected_outcome: z.string().optional(),
  estimated_budget: z.coerce.number().min(0),
});

type FormValues = z.infer<typeof schema>;

export default function RequestFormPage() {
  const { id } = useParams();
  const isEdit = !!id;
  const navigate = useNavigate();

  const { data: existing } = useRequestDetail(Number(id));
  const { mutate: create, isPending: creating } = useCreateDraft();
  const { mutate: update, isPending: updating } = useUpdateDraft(Number(id));

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    values: existing
      ? {
          title: existing.title,
          description: existing.description,
          business_goal: existing.business_goal,
          expected_outcome: existing.expected_outcome,
          estimated_budget: existing.estimated_budget,
        }
      : undefined,
  });

  const onSubmit = (values: FormValues) => {
    if (isEdit && existing) {
      update(
        { ...values, version: existing.version },
        { onSuccess: () => navigate(`/project-requests/${id}`) }
      );
    } else {
      create(values, {
        onSuccess: (data) => navigate(`/project-requests/${data.id}`),
      });
    }
  };

  return (
    <Card className="max-w-2xl">
      <CardHeader>
        <CardTitle>{isEdit ? "Edit Draft" : "New Project Request"}</CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="title">Title</Label>
            <Input id="title" {...register("title")} />
            {errors.title && <p className="text-sm text-red-500">{errors.title.message}</p>}
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Input id="description" {...register("description")} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="business_goal">Business Goal</Label>
            <Input id="business_goal" {...register("business_goal")} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="expected_outcome">Expected Outcome</Label>
            <Input id="expected_outcome" {...register("expected_outcome")} />
          </div>

          <div className="space-y-2">
            <Label htmlFor="estimated_budget">Estimated Budget (Rp)</Label>
            <Input id="estimated_budget" type="number" {...register("estimated_budget")} />
          </div>

          <Button type="submit" disabled={creating || updating}>
            {isEdit ? "Save Draft" : "Create Draft"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}