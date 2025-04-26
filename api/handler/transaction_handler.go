package handler

import (
	"net/http"

	"db/ent/transaction"
	"db/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TransactionHandler представляет обработчик для API транзакций
type TransactionHandler struct {
	transactionService *service.TransactionService
}

// NewTransactionHandler создает новый обработчик транзакций
func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// CreateTransactionRequest представляет запрос на создание транзакции
type CreateTransactionRequest struct {
	UserID   int     `json:"user_id" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Currency string  `json:"currency" binding:"required"`
	Type     string  `json:"type" binding:"required,oneof=deposit withdrawal"`
}

// CreateTransaction обрабатывает запрос на создание новой транзакции
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Преобразование строкового типа транзакции в тип для Ent
	var txType transaction.Type
	switch req.Type {
	case "deposit":
		txType = transaction.TypeDeposit
	case "withdrawal":
		txType = transaction.TypeWithdrawal
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid transaction type. Must be 'deposit' or 'withdrawal'",
		})
		return
	}

	// Генерируем уникальный ID для транзакции
	transactionID := uuid.New().String()

	tx, err := h.transactionService.Create(
		c.Request.Context(),
		transactionID,
		req.UserID,
		req.Currency,
		req.Amount,
		txType,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         tx.ID,
		"user_id":    tx.UserID,
		"amount":     tx.Amount,
		"currency":   tx.Currency,
		"type":       tx.Type,
		"created_at": tx.CreatedAt,
	})
}
