package service

import (
	"testing"
	"time"

	"github.com/Kal-el21/backend/internal/domain/project/entity"
)

func timePtr(t time.Time) *time.Time { return &t }

func TestCalculateProjectHealth(t *testing.T) {
	now := time.Date(2026, 7, 9, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name     string
		project  *entity.Project
		progress float64
		want     string
	}{
		{
			name:     "completed is green",
			project:  &entity.Project{Status: entity.ProjectCompleted},
			progress: 40,
			want:     "GREEN",
		},
		{
			name:     "cancelled is red",
			project:  &entity.Project{Status: entity.ProjectCancelled},
			progress: 10,
			want:     "RED",
		},
		{
			name:     "on hold is yellow",
			project:  &entity.Project{Status: entity.ProjectOnHold},
			progress: 10,
			want:     "YELLOW",
		},
		{
			name:     "overdue and not finished is red",
			project:  &entity.Project{Status: entity.ProjectActive, EndDate: timePtr(time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC))},
			progress: 50,
			want:     "RED",
		},
		{
			name:     "due within 7 days below 80 is yellow",
			project:  &entity.Project{Status: entity.ProjectActive, EndDate: timePtr(time.Date(2026, 7, 12, 0, 0, 0, 0, time.UTC))},
			progress: 50,
			want:     "YELLOW",
		},
		{
			name:     "active low progress is yellow",
			project:  &entity.Project{Status: entity.ProjectActive, EndDate: timePtr(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC))},
			progress: 10,
			want:     "YELLOW",
		},
		{
			name:     "active healthy is green",
			project:  &entity.Project{Status: entity.ProjectActive, EndDate: timePtr(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC))},
			progress: 60,
			want:     "GREEN",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := calculateProjectHealth(c.project, c.progress, now)
			if got != c.want {
				t.Errorf("calculateProjectHealth() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestNormalizeProjectPriority(t *testing.T) {
	got, err := normalizeProjectPriority("urgent")
	if err != nil || got != "URGENT" {
		t.Errorf("normalizeProjectPriority(urgent) = %q, %v", got, err)
	}
	if _, err := normalizeProjectPriority("bogus"); err == nil {
		t.Errorf("expected error for invalid priority")
	}
}
