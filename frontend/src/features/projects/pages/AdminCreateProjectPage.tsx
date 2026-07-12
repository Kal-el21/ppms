import { useForm, type SubmitHandler } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useNavigate } from "react-router-dom";
import { useCreateDirectProject } from "../hooks/useProjects";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { Textarea } from "../../../components/ui/textarea";
import { Label } from "../../../components/ui/label";
import { Select } from "../../../components/ui/select";
import { Card, CardHeader, CardTitle, CardContent } from "../../../components/ui/card";
import { PageHeader } from "../../../components/shared/PageHeader";
import type { CreateProjectDirectRequest } from "../types";
import type { ProjectPriority, BudgetType } from "../types";

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
  name: z.string().min(3, "Minimal 3 karakter"),
  description: z.string().max(2000, "Maksimal 2000 karakter").optional(),
  category: z.string().max(100, "Maksimal 100 karakter").optional(),
  initiation_type: z
    .string()
    .min(1, "Pilih jenis inisiasi")
    .refine((value) => initiationValues.includes(value as (typeof initiationValues)[number]), "Jenis inisiasi tidak valid"),
  priority: z
    .string()
    .min(1, "Pilih prioritas")
    .refine((value) => priorityValues.includes(value as (typeof priorityValues)[number]), "Prioritas tidak valid"),
  notes: z.string().max(2000, "Maksimal 2000 karakter").optional(),
  start_date: z.string().optional(),
  end_date: z.string().optional(),
  budget_type: z
    .string()
    .min(1, "Pilih jenis anggaran")
    .refine((value) => budgetTypeValues.includes(value as (typeof budgetTypeValues)[number]), "Jenis anggaran tidak valid"),
  budget_name: z.string().min(2, "Nama mata anggaran wajib diisi").max(200, "Maksimal 200 karakter"),
  allocated_budget: z.number().min(0, "Anggaran tidak boleh negatif"),
}).refine(
  (value) => !value.start_date || !value.end_date || value.end_date >= value.start_date,
  {
    path: ["end_date"],
    message: "Tanggal selesai tidak boleh sebelum tanggal mulai",
  }
);

type FormValues = z.infer<typeof schema>;

const defaultValues: FormValues = {
  name: "",
  description: "",
  category: "",
  initiation_type: "",
  priority: "MEDIUM",
  notes: "",
  start_date: "",
  end_date: "",
  budget_type: "",
  budget_name: "",
  allocated_budget: 0,
};

function FieldError({ message }: { message?: string }) {
  if (!message) return null;
  return <p className="mt-1.5 text-xs text-danger-600">{message}</p>;
}

export default function AdminCreateProjectPage() {
  const navigate = useNavigate();
  const { mutate: createDirect, isPending } = useCreateDirectProject();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues,
  });

  const onSubmit: SubmitHandler<FormValues> = (values) => {
    const payload: CreateProjectDirectRequest = {
      name: values.name,
      description: values.description || "",
      category: values.category || "",
      initiation_type: (values.initiation_type || null) as CreateProjectDirectRequest["initiation_type"],
      priority: values.priority as ProjectPriority,
      notes: values.notes || "",
      start_date: values.start_date || null,
      end_date: values.end_date || null,
      budget_type: values.budget_type as BudgetType,
      budget_name: values.budget_name,
      allocated_budget: values.allocated_budget,
    };
    createDirect(payload, {
      onSuccess: (data) => navigate(`/projects/${data.id}`),
    });
  };

  return (
    <div className="max-w-5xl">
      <PageHeader
        title="Tambah Project"
        subtitle="Buat project langsung tanpa melalui request workflow."
      />

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
        <Card>
          <CardHeader>
            <CardTitle>Info Project</CardTitle>
          </CardHeader>
          <CardContent className="grid gap-4 md:grid-cols-2">
            <div className="md:col-span-2">
              <Label htmlFor="name">Nama Project</Label>
              <Input id="name" placeholder="cth. Migrasi Sistem CRM ke Cloud" {...register("name")} />
              <FieldError message={errors.name?.message} />
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

            <div className="md:col-span-2">
              <Label htmlFor="notes">Catatan</Label>
              <Textarea id="notes" placeholder="Catatan portfolio atau risiko awal" {...register("notes")} />
              <FieldError message={errors.notes?.message} />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Timeline</CardTitle>
          </CardHeader>
          <CardContent className="grid gap-4 md:grid-cols-2">
            <div>
              <Label htmlFor="start_date">Tanggal Mulai</Label>
              <Input id="start_date" type="date" {...register("start_date")} />
              <FieldError message={errors.start_date?.message} />
            </div>

            <div>
              <Label htmlFor="end_date">Tanggal Selesai</Label>
              <Input id="end_date" type="date" {...register("end_date")} />
              <FieldError message={errors.end_date?.message} />
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
              <Label htmlFor="allocated_budget">Alokasi Anggaran (Rp)</Label>
              <Input id="allocated_budget" type="number" placeholder="0" {...register("allocated_budget", { valueAsNumber: true })} />
              <FieldError message={errors.allocated_budget?.message} />
            </div>
          </CardContent>
        </Card>

        <div className="flex gap-2">
          <Button type="submit" variant="primary" disabled={isPending}>
            Buat Project
          </Button>
          <Button type="button" variant="outline" onClick={() => navigate("/projects")}>
            Batal
          </Button>
        </div>
      </form>
    </div>
  );
}
