package handler

import (
	"net/http"

	"db/ent/transaction"
	"db/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

// TransactionHandler represents the handler for transaction API
type TransactionHandler struct {
	transactionService *service.TransactionService
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(transactionService *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// CreateTransactionRequest represents a request to create a transaction
type CreateTransactionRequest struct {
	UserID   int     `json:"user_id" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Currency string  `json:"currency" binding:"required"`
	Type     string  `json:"type" binding:"required,oneof=deposit withdrawal"`
}

// CreateTransaction handles the request to create a new transaction
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Convert string transaction type to Ent type
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

	// Generate a unique ID for the transaction
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

	// Use jsoniter if available via middleware, otherwise fallback to standard c.JSON
	if jsonValue, exists := c.Get("json"); exists {
		if json, ok := jsonValue.(jsoniter.API); ok {
			data, err := json.Marshal(gin.H{
				"id":         tx.ID,
				"user_id":    tx.UserID,
				"amount":     tx.Amount,
				"currency":   tx.Currency,
				"type":       tx.Type,
				"created_at": tx.CreatedAt,
			})
			if err == nil {
				c.Data(http.StatusCreated, "application/json", data)
				return
			}
		}
	}

	// Fallback to standard Gin JSON marshaling
	c.JSON(http.StatusCreated, gin.H{
		"id":         tx.ID,
		"user_id":    tx.UserID,
		"amount":     tx.Amount,
		"currency":   tx.Currency,
		"type":       tx.Type,
		"created_at": tx.CreatedAt,
	})
}
