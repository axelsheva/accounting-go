package repository

import (
	"context"
	"fmt"
	"time"

	"db/ent"
	"db/ent/transaction"
	"db/ent/user"
)

// TransactionRepository presents a repository for working with transactions
type TransactionRepository struct {
	client      *ent.Client
	balanceRepo *BalanceRepository
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(client *ent.Client, balanceRepo *BalanceRepository) *TransactionRepository {
	return &TransactionRepository{
		client:      client,
		balanceRepo: balanceRepo,
	}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(ctx context.Context, id string, userID int, amount float64,
	currency string, txType transaction.Type) (*ent.Transaction, error) {

	builder := r.client.Transaction.
		Create().
		SetID(id).
		SetUserID(userID).
		SetAmount(amount).
		SetCurrency(currency).
		SetType(txType).
		SetCreatedAt(time.Now())

	tx, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating %s transaction: %w", txType, err)
	}

	var amountWithSign float64
	if txType == transaction.TypeDeposit {
		amountWithSign = amount
	} else {
		amountWithSign = -amount
	}

	err = r.balanceRepo.Upsert(ctx, userID, currency, amountWithSign)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetByID gets a transaction by its ID
func (r *TransactionRepository) GetByID(ctx context.Context, id string) (*ent.Transaction, error) {
	tx, err := r.client.Transaction.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed querying transaction by ID: %w", err)
	}
	return tx, nil
}

// GetAllByUserID gets all transactions of a user
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

// GetAllByUserIDUsingEdge gets all transactions of a user using the edge in the graph
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

// GetByIDAndUserIDAndTypeInTx gets a transaction inside a database transaction
func (r *TransactionRepository) GetByIDAndUserIDAndTypeInTx(ctx context.Context, tx *ent.Tx, txID string, userID int,
	txType transaction.Type) (*ent.Transaction, error) {

	transaction, err := tx.Transaction.
		Query().
		Where(
			transaction.ID(txID),
			transaction.UserID(userID),
			transaction.TypeEQ(txType),
		).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying transaction in tx: %w", err)
	}
	return transaction, nil
}
