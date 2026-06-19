import { useNavigate } from "react-router-dom";
import { useProjectList } from "../hooks/useProjects";
import { Badge } from "../../../components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../../../components/ui/table";

export default function ProjectListPage() {
  const navigate = useNavigate();
  const { data, isLoading } = useProjectList();

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-semibold">Projects</h1>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Progress</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {data?.data.map((project) => (
            <TableRow key={project.id} className="cursor-pointer" onClick={() => navigate(`/projects/${project.id}`)}>
              <TableCell>{project.name}</TableCell>
              <TableCell>
                <Badge>{project.status}</Badge>
              </TableCell>
              <TableCell>{project.progress.toFixed(0)}%</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}