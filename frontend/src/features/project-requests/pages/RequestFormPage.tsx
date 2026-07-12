import { useForm, type SubmitHandler } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useNavigate, useParams } from "react-router-dom";
import { useCreateDraft, useRequestDetail, useUpdateDraft } from "../hooks/useRequests";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { Textarea } from "../../../components/ui/textarea";
import { Label } from "../../../components/ui/label";
import { Select } from "../../../components/ui/select";
import { Card, CardHeader, CardTitle, CardContent } from "../../../components/ui/card";
import { PageHeader } from "../../../components/shared/PageHeader";
import type { ProjectRequest } from "../types";

const initiationValues = ["NEW_INITIATIVE", "RENEWAL", "ENHANCEMENT"] as const;
const priorityValues = ["LOW", "MEDIUM", "HIGH", "URGENT"] as const;
const budgetTypeValues = ["CAPEX", "OPEX"] as const;

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

const schema = z.object({
  title: z.string().min(5, "Minimal 5 karakter"),
  description: z.string().max(2000, "Maksimal 2000 karakter").optional(),
  business_goal: z.string().max(1000, "Maksimal 1000 karakter").optional(),
  expected_outcome: z.string().max(1000, "Maksimal 1000 karakter").optional(),
  category: z.string().max(100, "Maksimal 100 karakter").optional(),
  initiation_type: z
    .string()
    .min(1, "Pilih jenis inisiasi")
    .refine((value) => initiationValues.includes(value as (typeof initiationValues)[number]), "Jenis inisiasi tidak valid"),
  priority: z
    .string()
    .min(1, "Pilih prioritas")
    .refine((value) => priorityValues.includes(value as (typeof priorityValues)[number]), "Prioritas tidak valid"),
  proposed_start_date: z.string().min(1, "Tanggal mulai wajib diisi"),
  proposed_end_date: z.string().min(1, "Tanggal selesai wajib diisi"),
  budget_type: z
    .string()
    .min(1, "Pilih jenis anggaran")
    .refine((value) => budgetTypeValues.includes(value as (typeof budgetTypeValues)[number]), "Jenis anggaran tidak valid"),
  budget_name: z.string().min(2, "Nama mata anggaran wajib diisi").max(200, "Maksimal 200 karakter"),
  estimated_budget: z.number().min(0, "Anggaran tidak boleh negatif"),
  notes: z.string().max(2000, "Maksimal 2000 karakter").optional(),
}).refine(
  (value) => !value.proposed_start_date || !value.proposed_end_date || value.proposed_end_date >= value.proposed_start_date,
  {
    path: ["proposed_end_date"],
    message: "Tanggal selesai tidak boleh sebelum tanggal mulai",
  }
);

type FormValues = z.infer<typeof schema>;

const defaultValues: FormValues = {
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
  estimated_budget: 0,
  notes: "",
};

function toDateInput(value?: string | null) {
  return value ? value.slice(0, 10) : "";
}

function toFormValues(request?: ProjectRequest): FormValues {
  if (!request) return defaultValues;

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
    estimated_budget: request.estimated_budget ?? 0,
    notes: request.notes || "",
  };
}

function FieldError({ message }: { message?: string }) {
  if (!message) return null;
  return <p className="mt-1.5 text-xs text-danger-600">{message}</p>;
}

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
    values: toFormValues(existing),
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
    <div className="max-w-5xl">
      <PageHeader
        title={isEdit ? "Edit Draft" : "New Project Request"}
        subtitle="Lengkapi data request sebelum dikirim untuk review Admin."
      />

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
        <Card>
          <CardHeader>
            <CardTitle>Info Project</CardTitle>
          </CardHeader>
          <CardContent className="grid gap-4 md:grid-cols-2">
            <div className="md:col-span-2">
              <Label htmlFor="title">Judul</Label>
              <Input id="title" placeholder="cth. Migrasi Sistem CRM ke Cloud" {...register("title")} />
              <FieldError message={errors.title?.message} />
            </div>

            <div>
              <Label htmlFor="category">Kategori / Tag</Label>
              <Input id="category" placeholder="cth. Infrastruktur, Aplikasi, Security" {...register("category")} />
              <FieldError message={errors.category?.message} />
            </div>

            <div>
              <Label htmlFor="initiation_type">Jenis Inisiasi</Label>
              <Select
                id="initiation_type"
                placeholder="Pilih jenis inisiasi"
                options={initiationOptions}
                {...register("initiation_type")}
              />
              <FieldError message={errors.initiation_type?.message} />
            </div>

            <div>
              <Label htmlFor="priority">Prioritas</Label>
              <Select id="priority" options={priorityOptions} {...register("priority")} />
              <FieldError message={errors.priority?.message} />
            </div>

            <div className="md:col-span-2">
              <Label htmlFor="description">Deskripsi</Label>
              <Textarea
                id="description"
                placeholder="Jelaskan ruang lingkup project secara singkat"
                {...register("description")}
              />
              <FieldError message={errors.description?.message} />
            </div>

            <div>
              <Label htmlFor="business_goal">Tujuan Bisnis</Label>
              <Textarea
                id="business_goal"
                placeholder="cth. Meningkatkan efisiensi operasional 20%"
                {...register("business_goal")}
              />
              <FieldError message={errors.business_goal?.message} />
            </div>

            <div>
              <Label htmlFor="expected_outcome">Hasil yang Diharapkan</Label>
              <Textarea
                id="expected_outcome"
                placeholder="cth. Sistem CRM baru live dan stabil"
                {...register("expected_outcome")}
              />
              <FieldError message={errors.expected_outcome?.message} />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Anggaran</CardTitle>
          </CardHeader>
          <CardContent className="grid gap-4 md:grid-cols-3">
            <div>
              <Label htmlFor="budget_type">Jenis Anggaran</Label>
              <Select
                id="budget_type"
                placeholder="Pilih jenis anggaran"
                options={budgetTypeOptions}
                {...register("budget_type")}
              />
              <FieldError message={errors.budget_type?.message} />
            </div>

            <div>
              <Label htmlFor="budget_name">Nama Mata Anggaran</Label>
              <Input id="budget_name" placeholder="cth. IT Capital Budget 2026" {...register("budget_name")} />
              <FieldError message={errors.budget_name?.message} />
            </div>

            <div>
              <Label htmlFor="estimated_budget">Estimasi Anggaran (Rp)</Label>
              <Input id="estimated_budget" type="number" placeholder="0" {...register("estimated_budget", { valueAsNumber: true })} />
              <FieldError message={errors.estimated_budget?.message} />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Timeline dan Catatan</CardTitle>
          </CardHeader>
          <CardContent className="grid gap-4 md:grid-cols-2">
            <div>
              <Label htmlFor="proposed_start_date">Tanggal Mulai</Label>
              <Input id="proposed_start_date" type="date" {...register("proposed_start_date")} />
              <FieldError message={errors.proposed_start_date?.message} />
            </div>

            <div>
              <Label htmlFor="proposed_end_date">Tanggal Selesai</Label>
              <Input id="proposed_end_date" type="date" {...register("proposed_end_date")} />
              <FieldError message={errors.proposed_end_date?.message} />
            </div>

            <div className="md:col-span-2">
              <Label htmlFor="notes">Catatan</Label>
              <Textarea id="notes" placeholder="Catatan portfolio atau risiko awal" {...register("notes")} />
              <FieldError message={errors.notes?.message} />
            </div>
          </CardContent>
        </Card>

        <div className="flex gap-2">
          <Button type="submit" variant="primary" disabled={creating || updating}>
            {isEdit ? "Simpan Draft" : "Buat Draft"}
          </Button>
          <Button type="button" variant="outline" onClick={() => navigate(-1)}>
            Batal
          </Button>
        </div>
      </form>
    </div>
  );
}
