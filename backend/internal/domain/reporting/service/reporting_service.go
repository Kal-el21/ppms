package service

import (
	"bytes"
	"fmt"

	"github.com/Kal-el21/backend/internal/domain/reporting/dto"
	"github.com/Kal-el21/backend/internal/domain/reporting/repository"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type ReportingService interface {
	Generate(req dto.GenerateReportRequest) ([]byte, string, string, error) // returns: fileBytes, fileName, contentType, error
}

type reportingService struct {
	repo repository.ReportingRepository
}

func NewReportingService(repo repository.ReportingRepository) ReportingService {
	return &reportingService{repo: repo}
}

func (s *reportingService) Generate(req dto.GenerateReportRequest) ([]byte, string, string, error) {
	switch dto.ReportType(req.Type) {
	case dto.ReportProject:
		return s.generateProjectReport(req.ProjectID, req.Format)
	case dto.ReportTask:
		return s.generateTaskReport(req.ProjectID, req.Format)
	case dto.ReportBudget:
		return s.generateBudgetReport(req.ProjectID, req.Format)
	case dto.ReportMilestone:
		return s.generateMilestoneReport(req.ProjectID, req.Format)
	case dto.ReportHandover:
		return s.generateHandoverReport(req.ProjectID, req.Format)
	default:
		return nil, "", "", fmt.Errorf("report type %s not yet implemented", req.Type)
	}
}

func (s *reportingService) generateMilestoneReport(projectID *uint64, format string) ([]byte, string, string, error) {
	rows, err := s.repo.GetMilestoneReport(projectID)
	if err != nil {
		return nil, "", "", err
	}

	if format == "PDF" {
		return generateMilestonePDF(rows)
	}
	return generateMilestoneExcel(rows)
}

func (s *reportingService) generateHandoverReport(projectID *uint64, format string) ([]byte, string, string, error) {
	rows, err := s.repo.GetHandoverReport(projectID)
	if err != nil {
		return nil, "", "", err
	}

	if format == "PDF" {
		return generateHandoverPDF(rows)
	}
	return generateHandoverExcel(rows)
}

func (s *reportingService) generateProjectReport(projectID *uint64, format string) ([]byte, string, string, error) {
	rows, err := s.repo.GetProjectReport(projectID)
	if err != nil {
		return nil, "", "", err
	}

	if format == "PDF" {
		return generateProjectPDF(rows)
	}
	return generateProjectExcel(rows)
}

func (s *reportingService) generateTaskReport(projectID *uint64, format string) ([]byte, string, string, error) {
	rows, err := s.repo.GetTaskReport(projectID)
	if err != nil {
		return nil, "", "", err
	}

	if format == "PDF" {
		return generateTaskPDF(rows)
	}
	return generateTaskExcel(rows)
}

func (s *reportingService) generateBudgetReport(projectID *uint64, format string) ([]byte, string, string, error) {
	rows, err := s.repo.GetBudgetReport(projectID)
	if err != nil {
		return nil, "", "", err
	}

	if format == "PDF" {
		return generateBudgetPDF(rows)
	}
	return generateBudgetExcel(rows)
}

// ===== PDF Generators =====

func generateProjectPDF(rows []dto.ProjectReportRow) ([]byte, string, string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Project Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(60, 8, "Name", "1", 0, "", false, 0, "")
	pdf.CellFormat(30, 8, "Status", "1", 0, "", false, 0, "")
	pdf.CellFormat(25, 8, "Progress", "1", 0, "", false, 0, "")
	pdf.CellFormat(35, 8, "Start Date", "1", 0, "", false, 0, "")
	pdf.CellFormat(35, 8, "End Date", "1", 1, "", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	for _, row := range rows {
		pdf.CellFormat(60, 8, row.Name, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 8, row.Status, "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 8, fmt.Sprintf("%.1f%%", row.Progress), "1", 0, "", false, 0, "")
		pdf.CellFormat(35, 8, row.StartDate, "1", 0, "", false, 0, "")
		pdf.CellFormat(35, 8, row.EndDate, "1", 1, "", false, 0, "")
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", "", err
	}

	return buf.Bytes(), "project_report.pdf", "application/pdf", nil
}

func generateTaskPDF(rows []dto.TaskReportRow) ([]byte, string, string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Task Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(70, 8, "Title", "1", 0, "", false, 0, "")
	pdf.CellFormat(30, 8, "Status", "1", 0, "", false, 0, "")
	pdf.CellFormat(25, 8, "Priority", "1", 0, "", false, 0, "")
	pdf.CellFormat(25, 8, "Progress", "1", 0, "", false, 0, "")
	pdf.CellFormat(35, 8, "Due Date", "1", 1, "", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	for _, row := range rows {
		pdf.CellFormat(70, 8, row.Title, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 8, row.Status, "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 8, row.Priority, "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 8, fmt.Sprintf("%d%%", row.Progress), "1", 0, "", false, 0, "")
		pdf.CellFormat(35, 8, row.DueDate, "1", 1, "", false, 0, "")
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", "", err
	}

	return buf.Bytes(), "task_report.pdf", "application/pdf", nil
}

func generateBudgetPDF(rows []dto.BudgetReportRow) ([]byte, string, string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Budget Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(60, 8, "Project", "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 8, "Allocated", "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 8, "Used", "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 8, "Remaining", "1", 1, "", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	for _, row := range rows {
		pdf.CellFormat(60, 8, row.ProjectName, "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 8, fmt.Sprintf("%.2f", row.AllocatedBudget), "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 8, fmt.Sprintf("%.2f", row.UsedBudget), "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 8, fmt.Sprintf("%.2f", row.RemainingBudget), "1", 1, "", false, 0, "")
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", "", err
	}

	return buf.Bytes(), "budget_report.pdf", "application/pdf", nil
}

// ===== Excel Generators =====

func generateProjectExcel(rows []dto.ProjectReportRow) ([]byte, string, string, error) {
	f := excelize.NewFile()
	sheet := "Project Report"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Name", "Status", "Progress (%)", "Start Date", "End Date"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, row := range rows {
		r := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", r), row.Name)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", r), row.Status)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", r), row.Progress)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", r), row.StartDate)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", r), row.EndDate)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", "", err
	}

	return buf.Bytes(), "project_report.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
}

func generateTaskExcel(rows []dto.TaskReportRow) ([]byte, string, string, error) {
	f := excelize.NewFile()
	sheet := "Task Report"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Title", "Status", "Priority", "Progress (%)", "Due Date"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, row := range rows {
		r := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", r), row.Title)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", r), row.Status)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", r), row.Priority)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", r), row.Progress)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", r), row.DueDate)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", "", err
	}

	return buf.Bytes(), "task_report.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
}

func generateBudgetExcel(rows []dto.BudgetReportRow) ([]byte, string, string, error) {
	f := excelize.NewFile()
	sheet := "Budget Report"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Project", "Allocated", "Used", "Remaining"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, row := range rows {
		r := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", r), row.ProjectName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", r), row.AllocatedBudget)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", r), row.UsedBudget)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", r), row.RemainingBudget)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", "", err
	}

	return buf.Bytes(), "budget_report.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
}

// ===== Milestone Report Generators =====

func generateMilestonePDF(rows []dto.MilestoneReportRow) ([]byte, string, string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Milestone Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(45, 8, "Project", "1", 0, "", false, 0, "")
	pdf.CellFormat(50, 8, "Milestone", "1", 0, "", false, 0, "")
	pdf.CellFormat(25, 8, "Status", "1", 0, "", false, 0, "")
	pdf.CellFormat(20, 8, "Progress", "1", 0, "", false, 0, "")
	pdf.CellFormat(25, 8, "Start", "1", 0, "", false, 0, "")
	pdf.CellFormat(25, 8, "End", "1", 1, "", false, 0, "")

	pdf.SetFont("Arial", "", 8)
	for _, row := range rows {
		pdf.CellFormat(45, 8, row.ProjectName, "1", 0, "", false, 0, "")
		pdf.CellFormat(50, 8, row.Name, "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 8, row.Status, "1", 0, "", false, 0, "")
		pdf.CellFormat(20, 8, fmt.Sprintf("%.1f%%", row.Progress), "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 8, row.StartDate, "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 8, row.EndDate, "1", 1, "", false, 0, "")
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", "", err
	}

	return buf.Bytes(), "milestone_report.pdf", "application/pdf", nil
}

func generateMilestoneExcel(rows []dto.MilestoneReportRow) ([]byte, string, string, error) {
	f := excelize.NewFile()
	sheet := "Milestone Report"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Project", "Milestone", "Status", "Progress (%)", "Start Date", "End Date"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, row := range rows {
		r := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", r), row.ProjectName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", r), row.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", r), row.Status)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", r), row.Progress)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", r), row.StartDate)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", r), row.EndDate)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", "", err
	}

	return buf.Bytes(), "milestone_report.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
}

// ===== Handover Report Generators =====

func generateHandoverPDF(rows []dto.HandoverReportRow) ([]byte, string, string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Handover Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(40, 8, "Project", "1", 0, "", false, 0, "")
	pdf.CellFormat(35, 8, "Sender", "1", 0, "", false, 0, "")
	pdf.CellFormat(35, 8, "Receiver", "1", 0, "", false, 0, "")
	pdf.CellFormat(25, 8, "Status", "1", 0, "", false, 0, "")
	pdf.CellFormat(25, 8, "Delivery", "1", 0, "", false, 0, "")
	pdf.CellFormat(30, 8, "Received At", "1", 1, "", false, 0, "")

	pdf.SetFont("Arial", "", 8)
	for _, row := range rows {
		pdf.CellFormat(40, 8, row.ProjectName, "1", 0, "", false, 0, "")
		pdf.CellFormat(35, 8, row.SenderName, "1", 0, "", false, 0, "")
		pdf.CellFormat(35, 8, row.ReceiverName, "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 8, row.Status, "1", 0, "", false, 0, "")
		pdf.CellFormat(25, 8, row.DeliveryDate, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 8, row.ReceivedAt, "1", 1, "", false, 0, "")
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", "", err
	}

	return buf.Bytes(), "handover_report.pdf", "application/pdf", nil
}

func generateHandoverExcel(rows []dto.HandoverReportRow) ([]byte, string, string, error) {
	f := excelize.NewFile()
	sheet := "Handover Report"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"Project", "Sender", "Receiver", "Status", "Delivery Date", "Received At"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	for i, row := range rows {
		r := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", r), row.ProjectName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", r), row.SenderName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", r), row.ReceiverName)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", r), row.Status)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", r), row.DeliveryDate)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", r), row.ReceivedAt)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", "", err
	}

	return buf.Bytes(), "handover_report.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
}
