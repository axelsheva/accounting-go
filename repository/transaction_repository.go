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

// Create creates a new transaction with SQL transaction
func (r *TransactionRepository) Create(ctx context.Context, id string, userID int, amount float64,
	currency string, txType transaction.Type) (*ent.Transaction, error) {

	// Start a transaction
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction: %w", err)
	}

	// Execute the actual logic within the transaction
	transaction, err := r.createWithTx(ctx, tx, id, userID, amount, currency, txType)
	if err != nil {
		// Rollback the transaction in case of error
		if err := tx.Rollback(); err != nil {
			return nil, fmt.Errorf("rolling back transaction: %w (%v)", err, err)
		}
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	return transaction, nil
}

// CreateWithTx creates a new transaction within an existing DB transaction
func (r *TransactionRepository) createWithTx(ctx context.Context, tx *ent.Tx, id string, userID int, amount float64,
	currency string, txType transaction.Type) (*ent.Transaction, error) {

	builder := tx.Transaction.
		Create().
		SetID(id).
		SetUserID(userID).
		SetAmount(amount).
		SetCurrency(currency).
		SetType(txType).
		SetCreatedAt(time.Now())

	transaction, err := builder.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating %s transaction: %w", txType, err)
	}

	var amountWithSign float64
	if txType == "deposit" {
		amountWithSign = amount
	} else {
		amountWithSign = -amount
	}

	// Use the BalanceRepository with the transaction context
	err = r.balanceRepo.UpsertWithTx(ctx, tx, UpsertBalanceParams{
		UserID:   userID,
		Currency: currency,
		Amount:   amountWithSign,
	})
	if err != nil {
		return nil, err
	}

	return transaction, nil
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
