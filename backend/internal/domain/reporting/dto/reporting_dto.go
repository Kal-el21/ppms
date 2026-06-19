package dto

type ReportType string

const (
	ReportProject   ReportType = "PROJECT"
	ReportMilestone ReportType = "MILESTONE"
	ReportTask      ReportType = "TASK"
	ReportBudget    ReportType = "BUDGET"
	ReportHandover  ReportType = "HANDOVER"
)

type GenerateReportRequest struct {
	Type      string  `json:"type" validate:"required,oneof=PROJECT MILESTONE TASK BUDGET HANDOVER"`
	ProjectID *uint64 `json:"project_id"`
	Format    string  `json:"format" validate:"required,oneof=PDF EXCEL"`
}

type ProjectReportRow struct {
	ID        uint64
	Name      string
	Status    string
	Progress  float64
	StartDate string
	EndDate   string
}

type TaskReportRow struct {
	ID       uint64
	Title    string
	Status   string
	Priority string
	Progress int
	DueDate  string
}

type BudgetReportRow struct {
	ProjectName     string
	AllocatedBudget float64
	UsedBudget      float64
	RemainingBudget float64
}

type MilestoneReportRow struct {
	ProjectName string
	Name        string
	Status      string
	Progress    float64
	StartDate   string
	EndDate     string
}

type HandoverReportRow struct {
	ProjectName  string
	SenderName   string
	ReceiverName string
	Status       string
	DeliveryDate string
	ReceivedAt   string
}
