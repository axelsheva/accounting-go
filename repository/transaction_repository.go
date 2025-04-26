package repository

import (
	"context"
	"fmt"
	"time"

	"db/ent"
	"db/ent/transaction"
	"db/ent/user"

	"github.com/google/uuid"
)

// TransactionRepository представляет репозиторий для работы с транзакциями
type TransactionRepository struct {
	client *ent.Client
}

// NewTransactionRepository создаёт новый репозиторий транзакций
func NewTransactionRepository(client *ent.Client) *TransactionRepository {
	return &TransactionRepository{
		client: client,
	}
}

// GenerateTransactionID генерирует уникальный ID транзакции
func (r *TransactionRepository) GenerateTransactionID() string {
	return uuid.New().String()
}

// GeneratePrefixedTransactionID генерирует ID транзакции с префиксом
func (r *TransactionRepository) GeneratePrefixedTransactionID(prefix string) string {
	return prefix + "-" + r.GenerateTransactionID()
}

// Create создаёт новую транзакцию
func (r *TransactionRepository) Create(ctx context.Context, id string, userID int, amount float64,
	currency string, txType transaction.Type, description, status string, completedAt *time.Time) (*ent.Transaction, error) {

	builder := r.client.Transaction.
		Create().
		SetID(id).
		SetUserID(userID).
		SetAmount(amount).
		SetCurrency(currency).
		SetType(txType).
		SetStatus(status).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now())

	if description != "" {
		builder = builder.SetDescription(description)
	}

	if completedAt != nil {
		builder = builder.SetCompletedAt(*completedAt)
	}

	tx, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating %s transaction: %w", txType, err)
	}

	return tx, nil
}

// CreateInTx создаёт новую транзакцию внутри существующей транзакции БД
func (r *TransactionRepository) CreateInTx(ctx context.Context, tx *ent.Tx, id string, userID int, amount float64,
	currency string, txType transaction.Type, description, status string, completedAt *time.Time) (*ent.Transaction, error) {

	builder := tx.Transaction.
		Create().
		SetID(id).
		SetUserID(userID).
		SetAmount(amount).
		SetCurrency(currency).
		SetType(txType).
		SetStatus(status).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now())

	if description != "" {
		builder = builder.SetDescription(description)
	}

	if completedAt != nil {
		builder = builder.SetCompletedAt(*completedAt)
	}

	transaction, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating %s transaction: %w", txType, err)
	}

	return transaction, nil
}

// CreateSampleTransactions создаёт набор тестовых транзакций для пользователя
func (r *TransactionRepository) CreateSampleTransactions(ctx context.Context, user *ent.User) ([]*ent.Transaction, error) {
	var transactions []*ent.Transaction

	// Create a deposit transaction
	now := time.Now()
	depositID := "TX-" + r.GenerateTransactionID()
	deposit, err := r.Create(ctx, depositID, user.ID, 1000.00, "USD", transaction.TypeDeposit, "Initial deposit", "completed", &now)
	if err != nil {
		return nil, err
	}
	transactions = append(transactions, deposit)

	// Create a withdrawal transaction
	withdrawalID := "TX-" + r.GenerateTransactionID()
	withdrawal, err := r.Create(ctx, withdrawalID, user.ID, 150.50, "USD", transaction.TypeWithdrawal, "ATM withdrawal", "completed", &now)
	if err != nil {
		return nil, err
	}
	transactions = append(transactions, withdrawal)

	// Create a pending transfer transaction
	transferID := "TX-" + r.GenerateTransactionID()
	transfer, err := r.Create(ctx, transferID, user.ID, 300.00, "USD", transaction.TypeTransfer, "Transfer to Alice", "pending", nil)
	if err != nil {
		return nil, err
	}
	transactions = append(transactions, transfer)

	return transactions, nil
}

// GetByID получает транзакцию по её ID
func (r *TransactionRepository) GetByID(ctx context.Context, id string) (*ent.Transaction, error) {
	tx, err := r.client.Transaction.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed querying transaction by ID: %w", err)
	}
	return tx, nil
}

// UpdateStatus обновляет статус транзакции
func (r *TransactionRepository) UpdateStatus(ctx context.Context, id string, status string) (*ent.Transaction, error) {
	tx, err := r.client.Transaction.
		UpdateOneID(id).
		SetStatus(status).
		SetUpdatedAt(time.Now()).
		SetCompletedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed updating transaction status: %w", err)
	}
	return tx, nil
}

// UpdateStatusInTx обновляет статус транзакции внутри существующей транзакции БД
func (r *TransactionRepository) UpdateStatusInTx(ctx context.Context, tx *ent.Tx, transaction *ent.Transaction,
	status string) (*ent.Transaction, error) {

	updatedTx, err := tx.Transaction.
		UpdateOne(transaction).
		SetStatus(status).
		SetUpdatedAt(time.Now()).
		SetCompletedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed updating transaction status: %w", err)
	}
	return updatedTx, nil
}

// GetAllByUserID получает все транзакции пользователя
func (r *TransactionRepository) GetAllByUserID(ctx context.Context, userID int) ([]*ent.Transaction, error) {
	txs, err := r.client.Transaction.
		Query().
		Where(transaction.UserID(userID)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying transactions by user_id: %w", err)
	}
	return txs, nil
}

// GetAllByUserIDUsingEdge получает все транзакции пользователя через связь в графе
func (r *TransactionRepository) GetAllByUserIDUsingEdge(ctx context.Context, userID int) ([]*ent.Transaction, error) {
	txs, err := r.client.User.
		Query().
		Where(user.ID(userID)).
		QueryTransactions().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user's transactions: %w", err)
	}
	return txs, nil
}

// GetAllByStatus получает все транзакции с указанным статусом
func (r *TransactionRepository) GetAllByStatus(ctx context.Context, status string) ([]*ent.Transaction, error) {
	txs, err := r.client.Transaction.
		Query().
		Where(transaction.StatusEQ(status)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying %s transactions: %w", status, err)
	}
	return txs, nil
}

// GetByUserIDAndTypeAndStatus получает транзакцию по userID, типу и статусу
func (r *TransactionRepository) GetByUserIDAndTypeAndStatus(ctx context.Context, userID int, txType transaction.Type, status string) (*ent.Transaction, error) {
	tx, err := r.client.Transaction.
		Query().
		Where(
			transaction.UserID(userID),
			transaction.TypeEQ(txType),
			transaction.StatusEQ(status),
		).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying transaction by user_id, type and status: %w", err)
	}
	return tx, nil
}

// GetByIDAndUserIDAndTypeAndStatus получает транзакцию по её ID, userID, типу и статусу
func (r *TransactionRepository) GetByIDAndUserIDAndTypeAndStatus(ctx context.Context, txID string, userID int,
	txType transaction.Type, status string) (*ent.Transaction, error) {

	tx, err := r.client.Transaction.
		Query().
		Where(
			transaction.ID(txID),
			transaction.UserID(userID),
			transaction.TypeEQ(txType),
			transaction.StatusEQ(status),
		).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying transaction by id, user_id, type and status: %w", err)
	}
	return tx, nil
}

// GetByIDAndUserIDAndTypeAndStatusInTx получает транзакцию внутри транзакции БД
func (r *TransactionRepository) GetByIDAndUserIDAndTypeAndStatusInTx(ctx context.Context, tx *ent.Tx, txID string, userID int,
	txType transaction.Type, status string) (*ent.Transaction, error) {

	transaction, err := tx.Transaction.
		Query().
		Where(
			transaction.ID(txID),
			transaction.UserID(userID),
			transaction.TypeEQ(txType),
			transaction.StatusEQ(status),
		).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying transaction in tx: %w", err)
	}
	return transaction, nil
}
