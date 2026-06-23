import { useForm, type SubmitHandler } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useNavigate, useParams } from "react-router-dom";
import { useCreateDraft, useRequestDetail, useUpdateDraft } from "../hooks/useRequests";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import { Card, CardHeader, CardTitle, CardContent } from "../../../components/ui/card";
import { PageHeader } from "../../../components/shared/PageHeader";

const schema = z.object({
  title: z.string().min(5, "Minimal 5 karakter"),
  description: z.string().optional(),
  business_goal: z.string().optional(),
  expected_outcome: z.string().optional(),
  estimated_budget: z.number().min(0),
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

  const onSubmit: SubmitHandler<FormValues> = (values) => {
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
    <div className="max-w-2xl">
      <PageHeader
        title={isEdit ? "Edit Draft" : "New Project Request"}
        subtitle="Lengkapi detail pengajuan project sebelum disubmit untuk review."
      />

      <Card>
        <CardHeader>
          <CardTitle>Informasi Request</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <Label htmlFor="title">Judul</Label>
              <Input id="title" placeholder="cth. Migrasi Sistem CRM ke Cloud" {...register("title")} />
              {errors.title && <p className="text-xs text-danger-600 mt-1.5">{errors.title.message}</p>}
            </div>

            <div>
              <Label htmlFor="description">Deskripsi</Label>
              <Input id="description" placeholder="Jelaskan ruang lingkup project secara singkat" {...register("description")} />
            </div>

            <div>
              <Label htmlFor="business_goal">Tujuan Bisnis</Label>
              <Input id="business_goal" placeholder="cth. Meningkatkan efisiensi operasional 20%" {...register("business_goal")} />
            </div>

            <div>
              <Label htmlFor="expected_outcome">Hasil yang Diharapkan</Label>
              <Input id="expected_outcome" placeholder="cth. Sistem CRM baru live dan stabil" {...register("expected_outcome")} />
            </div>

            <div>
              <Label htmlFor="estimated_budget">Estimasi Anggaran (Rp)</Label>
              <Input id="estimated_budget" type="number" placeholder="0" {...register("estimated_budget", { valueAsNumber: true })} />
            </div>

            <div className="flex gap-2 pt-2">
              <Button type="submit" variant="primary" disabled={creating || updating}>
                {isEdit ? "Simpan Draft" : "Buat Draft"}
              </Button>
              <Button type="button" variant="outline" onClick={() => navigate(-1)}>
                Batal
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}