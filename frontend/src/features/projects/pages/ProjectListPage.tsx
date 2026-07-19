import { useState, useCallback } from "react";
import { Component, type ErrorInfo, type ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/features/auth/context/AuthContext";
import { useProjectList, type ProjectListFilters } from "../hooks/useProjects";
import { useDeleteProject } from "../hooks/useProjects";
import { StatusBadge, getStatusColor } from "../../../components/ui/status-badge";
import { HealthBar } from "../../../components/ui/health-bar";
import { PageHeader } from "../../../components/shared/PageHeader";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { Select } from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
  TableDensityToggle,
  type TableDensity,
} from "../../../components/ui/table";
import { BulkActionsDropdown } from "@/components/shared/BulkActionsDropdown";
import { useTableSelection } from "@/components/shared/useTableSelection";

interface ProjectListErrorBoundaryProps {
  children: ReactNode;
}

interface ProjectListErrorBoundaryState {
  hasError: boolean;
  errorMessage: string | null;
}

class ProjectListErrorBoundary extends Component<
  ProjectListErrorBoundaryProps,
  ProjectListErrorBoundaryState
> {
  state: ProjectListErrorBoundaryState = {
    hasError: false,
    errorMessage: null,
  };

  static getDerivedStateFromError(error: Error): ProjectListErrorBoundaryState {
    return {
      hasError: true,
      errorMessage: error.message,
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error("ProjectListPage render error", error, errorInfo);
  }

  reset = () => {
    this.setState({
      hasError: false,
      errorMessage: null,
    });
  };

  render() {
    if (this.state.hasError) {
      return (
        <div className="max-w-xl rounded-lg border border-danger-200 bg-danger-50 px-4 py-3 dark:border-danger-900/40 dark:bg-danger-900/20">
          <p className="text-[13.5px] font-semibold text-danger-700 dark:text-danger-400">
            Gagal menampilkan daftar project
          </p>
          <p className="mt-1 text-[12.5px] text-danger-700/80 dark:text-danger-400/80">
            {this.state.errorMessage || "Terjadi kesalahan saat memuat tampilan project."}
          </p>
          <Button type="button" variant="outline" size="sm" className="mt-3" onClick={this.reset}>
            Coba lagi
          </Button>
        </div>
      );
    }

    return this.props.children;
  }
}

const getHealthColor = (health?: string) => {
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

const formatDate = (value?: string | null) =>
  value ? new Date(value).toLocaleDateString("id-ID", { year: "numeric", month: "short", day: "numeric" }) : "—";

const formatCurrency = (value?: number) =>
  typeof value === "number" ? `Rp ${Math.round(value).toLocaleString("id-ID")}` : "—";

function ProjectListContent() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const isAdmin = user?.system_role === "ADMIN";

  const [filters, setFilters] = useState<ProjectListFilters>({
    page: 1,
    limit: 20,
    sort: "newest",
  });

  const setFilter = useCallback(
    (patch: Partial<ProjectListFilters>) => {
      setFilters((prev) => ({ ...prev, ...patch, page: 1 }));
    },
    []
  );

  const { data, isLoading } = useProjectList(filters);
  const deleteProject = useDeleteProject();
  const [density, setDensity] = useState<TableDensity>("comfortable");

  const list = data?.data ?? [];
  const {
    selectedIds,
    toggle,
    toggleAll,
    selectedCount,
    isAllSelected,
    isIndeterminate,
  } = useTableSelection<number>();

  const handleBulkDelete = () => {
    const ids = Array.from(selectedIds);
    if (ids.length === 0) return;
    deleteProject.mutate(ids[0]);
  };

  return (
    <div>
      <PageHeader
        title="Projects"
        subtitle={`${data?.meta.total ?? 0} project terdaftar`}
        actions={
          <div className="flex gap-2">
            <TableDensityToggle value={density} onChange={setDensity} />
            <BulkActionsDropdown
              selectedCount={selectedCount}
              onDelete={handleBulkDelete}
            />
            {isAdmin && (
              <Button variant="outline" size="sm" onClick={() => navigate("/projects/new")}>
                Tambah Project
              </Button>
            )}
            <Button variant="primary" size="sm" onClick={() => navigate("/project-requests/new")}>
              Request Baru
            </Button>
          </div>
        }
      />

      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-2.5 mb-4">
        <Input
          placeholder="Cari project..."
          value={filters.search ?? ""}
          onChange={(e) => setFilter({ search: e.target.value })}
          className="col-span-2"
        />
        <Select
          value={filters.status ?? ""}
          onChange={(e) => setFilter({ status: e.target.value })}
          options={[
            { value: "", label: "Semua status" },
            { value: "PLANNED", label: "Planned" },
            { value: "ACTIVE", label: "Active" },
            { value: "ON_HOLD", label: "On Hold" },
            { value: "COMPLETED", label: "Completed" },
            { value: "CANCELLED", label: "Cancelled" },
          ]}
        />
        <Select
          value={filters.budget_type ?? ""}
          onChange={(e) => setFilter({ budget_type: e.target.value })}
          options={[
            { value: "", label: "Semua budget" },
            { value: "CAPEX", label: "CAPEX" },
            { value: "OPEX", label: "OPEX" },
          ]}
        />
        <Select
          value={filters.initiation_type ?? ""}
          onChange={(e) => setFilter({ initiation_type: e.target.value })}
          options={[
            { value: "", label: "Semua inisiasi" },
            { value: "NEW_INITIATIVE", label: "New Initiative" },
            { value: "RENEWAL", label: "Renewal" },
            { value: "ENHANCEMENT", label: "Enhancement" },
          ]}
        />
        <Select
          value={filters.priority ?? ""}
          onChange={(e) => setFilter({ priority: e.target.value })}
          options={[
            { value: "", label: "Semua priority" },
            { value: "LOW", label: "Low" },
            { value: "MEDIUM", label: "Medium" },
            { value: "HIGH", label: "High" },
            { value: "URGENT", label: "Urgent" },
          ]}
        />
        <Select
          value={filters.sort ?? "newest"}
          onChange={(e) => setFilter({ sort: e.target.value })}
          options={[
            { value: "newest", label: "Terbaru" },
            { value: "code", label: "Kode Project" },
            { value: "name", label: "Nama" },
            { value: "end_date", label: "Deadline" },
            { value: "progress", label: "Progress" },
            { value: "budget", label: "Budget" },
          ]}
        />
      </div>

      <Table density={density}>
        <TableHeader>
          <TableRow>
            <TableHead style={{ width: 40 }}>
              <input
                type="checkbox"
                checked={isAllSelected(list.map((p) => p.id))}
                ref={(el) => {
                  if (el) el.indeterminate = isIndeterminate(list.map((p) => p.id));
                }}
                onChange={(e) => toggleAll(list.map((p) => p.id), e.target.checked)}
              />
            </TableHead>
            <TableHead style={{ width: 50 }}>No</TableHead>
            <TableHead>Project</TableHead>
            <TableHead>Kategori</TableHead>
            <TableHead>Inisiasi</TableHead>
            <TableHead>Priority</TableHead>
            <TableHead>Budget</TableHead>
            <TableHead>Mulai</TableHead>
            <TableHead>Selesai</TableHead>
            <TableHead>Health</TableHead>
            <TableHead style={{ width: 200 }}>Progress</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {list.map((project, index) => {
            const atRisk = project.status === "ON_HOLD" || project.progress < 30;
            return (
              <TableRow
                key={project.id}
                className="cursor-pointer"
                onClick={() => navigate(`/projects/${project.id}`)}
              >
                <TableCell onClick={(e) => e.stopPropagation()}>
                  <input
                    type="checkbox"
                    checked={selectedIds.has(project.id)}
                    onChange={() => toggle(project.id)}
                  />
                </TableCell>
                <TableCell className="text-ink-secondary text-[13px]" onClick={(e) => e.stopPropagation()}>
                  {(filters.page ? (filters.page - 1) * (filters.limit || 20) : 0) + index + 1}
                </TableCell>
                <TableCell>
                  <div className="font-medium text-ink-primary">{project.name}</div>
                  <div className="text-[11.5px] text-ink-tertiary">{project.project_code || "—"}</div>
                  {(project.budget_allocated ?? 0) > 0 && (
                    <div className="mt-1 text-[11px] text-ink-tertiary">
                      {formatCurrency(project.budget_allocated)} / {formatCurrency(project.budget_used)} used
                    </div>
                  )}
                </TableCell>
                <TableCell>
                  {project.category ? <StatusBadge color="gray">{project.category}</StatusBadge> : "—"}
                </TableCell>
                <TableCell>
                  {project.initiation_type ? (
                    <StatusBadge color="blue">{project.initiation_type}</StatusBadge>
                  ) : (
                    "—"
                  )}
                </TableCell>
                <TableCell>
                  <StatusBadge color={getStatusColor(project.priority)}>{project.priority}</StatusBadge>
                </TableCell>
                <TableCell>
                  {project.budget_type ? (
                    <StatusBadge color={project.budget_type === "CAPEX" ? "indigo" : "teal"}>
                      {project.budget_type}
                    </StatusBadge>
                  ) : (
                    "—"
                  )}
                </TableCell>
                <TableCell className="text-[12.5px] text-ink-secondary">{formatDate(project.start_date)}</TableCell>
                <TableCell className="text-[12.5px] text-ink-secondary">{formatDate(project.end_date)}</TableCell>
                <TableCell>
                  <StatusBadge color={getHealthColor(project.health)}>{project.health}</StatusBadge>
                </TableCell>
                <TableCell onClick={(e) => e.stopPropagation()}>
                  <div className="flex items-center gap-2.5">
                    <HealthBar progress={project.progress} atRisk={atRisk} className="flex-1" />
                    <span className="text-xs font-semibold text-ink-secondary w-9 text-right">
                      {(project.progress ?? 0).toFixed(0)}%
                    </span>
                  </div>
                </TableCell>
              </TableRow>
            );
          })}
        </TableBody>
      </Table>

      {!isLoading && list.length === 0 && (
        <p className="text-[13px] text-ink-tertiary mt-6 text-center">
          Tidak ada project yang cocok dengan filter.
        </p>
      )}
    </div>
  );
}

export default function ProjectListPage() {
  return (
    <ProjectListErrorBoundary>
      <ProjectListContent />
    </ProjectListErrorBoundary>
  );
}
