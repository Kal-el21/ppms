import { Component, type ErrorInfo, type ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import { useProjectList } from "../hooks/useProjects";
import { StatusBadge, getStatusColor } from "../../../components/ui/status-badge";
import { HealthBar } from "../../../components/ui/health-bar";
import { PageHeader } from "../../../components/shared/PageHeader";
import { Button } from "../../../components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../../components/ui/table";
import { TableSkeleton } from "@/components/ui/skeleton";

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

function ProjectListContent() {
  const navigate = useNavigate();
  const { data, isLoading } = useProjectList();

  if (isLoading) {
    return (
      <div>
        <div className="flex justify-between mb-6">
          <div className="h-7 w-32 bg-surface-tertiary rounded-md animate-pulse" />
        </div>
        <TableSkeleton rows={6} cols={3} />
      </div>
    );
  }

  return (
    <div>
      <PageHeader title="Projects" subtitle={`${data?.meta.total ?? 0} project terdaftar`} />

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Project</TableHead>
            <TableHead>Status</TableHead>
            <TableHead style={{ width: 220 }}>Progress</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {data?.data.map((project) => {
            const atRisk = project.status === "ON_HOLD" || project.progress < 30;
            return (
              <TableRow
                key={project.id}
                className="cursor-pointer"
                onClick={() => navigate(`/projects/${project.id}`)}
              >
                <TableCell>
                  <div className="font-medium text-ink-primary">{project.name}</div>
                </TableCell>
                <TableCell>
                  <StatusBadge color={getStatusColor(project.status)}>{project.status}</StatusBadge>
                </TableCell>
                <TableCell>
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
