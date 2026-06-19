package handler

import (
	"net/http"
	"strconv"

	auditservice "github.com/Kal-el21/backend/internal/domain/audit/service"
	"github.com/Kal-el21/backend/internal/domain/budget/dto"
	"github.com/Kal-el21/backend/internal/domain/budget/service"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"github.com/Kal-el21/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	service  service.TransactionService
	auditSvc auditservice.AuditService
}

func NewTransactionHandler(service service.TransactionService, auditSvc auditservice.AuditService) *TransactionHandler {
	return &TransactionHandler{service: service, auditSvc: auditSvc}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	budgetID, err := strconv.ParseUint(c.Param("budgetId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid budget id"))
		return
	}

	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, err.Error()))
		return
	}

	createdBy := c.GetUint64("user_id")

	tx, err := h.service.Create(budgetID, createdBy, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Transaksi budget bersifat finansial dan immutable — audit log di sini
	// menjadi SATU-SATUNYA jejak "siapa yang membuat transaksi ini dan kapan"
	// di luar tabel budget_transactions itu sendiri (yang sudah punya created_by/created_at).
	// Tetap di-log untuk konsistensi cross-module audit trail dan pencarian terpusat.
	h.auditSvc.Log(auditservice.LogParams{
		UserID:     &createdBy,
		Module:     "budget",
		Action:     "CREATE_TRANSACTION_" + req.Type,
		EntityType: "BUDGET_TRANSACTION",
		EntityID:   &tx.ID,
		NewData:    tx,
	})

	response.Success(c, http.StatusCreated, tx, "transaction recorded successfully")
}

func (h *TransactionHandler) GetList(c *gin.Context) {
	budgetID, err := strconv.ParseUint(c.Param("budgetId"), 10, 64)
	if err != nil {
		response.Error(c, apperrors.New(apperrors.ErrValidation, "invalid budget id"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	transactions, total, err := h.service.GetByBudgetIDPaginated(budgetID, page, limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    transactions,
		"meta":    gin.H{"page": page, "limit": limit, "total": total},
	})
}
