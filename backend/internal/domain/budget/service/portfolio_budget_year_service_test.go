package service

import (
	"errors"
	"testing"

	"github.com/Kal-el21/backend/internal/domain/budget/dto"
	"github.com/Kal-el21/backend/internal/domain/budget/entity"
	"github.com/Kal-el21/backend/internal/domain/budget/repository"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"gorm.io/gorm"
)

type fakePortfolioBudgetYearRepo struct {
	years []entity.PortfolioBudgetYear
}

func (f *fakePortfolioBudgetYearRepo) Create(y *entity.PortfolioBudgetYear) error {
	stored := *y
	if stored.Version == 0 {
		stored.Version = 1
	}
	f.years = append(f.years, stored)
	return nil
}

func (f *fakePortfolioBudgetYearRepo) FindByID(id uint64) (*entity.PortfolioBudgetYear, error) {
	for i := range f.years {
		if f.years[i].ID == id {
			return &f.years[i], nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakePortfolioBudgetYearRepo) FindByYear(year int) (*entity.PortfolioBudgetYear, error) {
	for i := range f.years {
		if f.years[i].Year == year {
			return &f.years[i], nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakePortfolioBudgetYearRepo) FindAll() ([]entity.PortfolioBudgetYear, error) {
	return f.years, nil
}

func (f *fakePortfolioBudgetYearRepo) UpdateWithVersionCheck(y *entity.PortfolioBudgetYear, expectedVersion int) (int64, error) {
	for i := range f.years {
		if f.years[i].ID == y.ID && f.years[i].Version == expectedVersion {
			f.years[i].CapexCeiling = y.CapexCeiling
			f.years[i].OpexCeiling = y.OpexCeiling
			f.years[i].Version++
			return 1, nil
		}
	}
	return 0, nil
}

func (f *fakePortfolioBudgetYearRepo) Delete(id uint64) error {
	for i := range f.years {
		if f.years[i].ID == id {
			f.years = append(f.years[:i], f.years[i+1:]...)
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func newFakeYearService() (*portfolioBudgetYearService, *fakePortfolioBudgetYearRepo) {
	repo := &fakePortfolioBudgetYearRepo{}
	return &portfolioBudgetYearService{repo: repo}, repo
}

func isValidation(err error) bool {
	var ae *apperrors.AppError
	return errors.As(err, &ae) && ae.Code == apperrors.ErrValidation
}

func TestPortfolioBudgetYearCreateValidation(t *testing.T) {
	svc, _ := newFakeYearService()

	if _, err := svc.Create(1, dto.CreatePortfolioBudgetYearRequest{Year: 1999, CapexCeiling: 0, OpexCeiling: 0}); !isValidation(err) {
		t.Errorf("expected validation error for out-of-range year, got %v", err)
	}

	if _, err := svc.Create(1, dto.CreatePortfolioBudgetYearRequest{Year: 2026, CapexCeiling: -1, OpexCeiling: 0}); !isValidation(err) {
		t.Errorf("expected validation error for negative ceiling, got %v", err)
	}

	if _, err := svc.Create(1, dto.CreatePortfolioBudgetYearRequest{Year: 2026, CapexCeiling: 100, OpexCeiling: 50}); err != nil {
		t.Errorf("expected success, got %v", err)
	}

	if _, err := svc.Create(1, dto.CreatePortfolioBudgetYearRequest{Year: 2026, CapexCeiling: 100, OpexCeiling: 50}); err == nil {
		t.Errorf("expected duplicate-year error")
	}
}

func TestPortfolioBudgetYearUpdateVersionMismatch(t *testing.T) {
	svc, _ := newFakeYearService()
	created, err := svc.Create(1, dto.CreatePortfolioBudgetYearRequest{Year: 2027, CapexCeiling: 100, OpexCeiling: 50})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Wrong version → no rows updated → version mismatch
	if _, err := svc.Update(created.ID, 1, dto.UpdatePortfolioBudgetYearRequest{CapexCeiling: 200, OpexCeiling: 100, Version: 999}); err == nil {
		t.Errorf("expected version mismatch error")
	}

	// Correct version → success
	if _, err := svc.Update(created.ID, 1, dto.UpdatePortfolioBudgetYearRequest{CapexCeiling: 200, OpexCeiling: 100, Version: 1}); err != nil {
		t.Errorf("expected success on correct version, got %v", err)
	}
}

var _ repository.PortfolioBudgetYearRepository = (*fakePortfolioBudgetYearRepo)(nil)
