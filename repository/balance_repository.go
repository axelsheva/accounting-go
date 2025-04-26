package repository

import (
	"context"
	"fmt"
	"time"

	"db/ent"
	"db/ent/balance"
	"db/ent/user"
	"db/errors"
)

// BalanceRepository представляет репозиторий для работы с балансами
type BalanceRepository struct {
	client *ent.Client
}

// NewBalanceRepository создаёт новый репозиторий балансов
func NewBalanceRepository(client *ent.Client) *BalanceRepository {
	return &BalanceRepository{
		client: client,
	}
}

// GetByUserIDAndCurrency returns the balance of a user in a specified currency
func (r *BalanceRepository) GetByUserIDAndCurrency(ctx context.Context, userID int, currency string) (*ent.Balance, error) {
	balance, err := r.client.Balance.
		Query().
		Where(
			balance.UserID(userID),
			balance.CurrencyEQ(currency),
		).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying %s balance: %w", currency, err)
	}

	return balance, nil
}

// UpsertBalanceParams represents the parameters for the UpsertWithTx method
type UpsertBalanceParams struct {
	UserID   int
	Currency string
	Amount   float64
}

// UpsertWithTx creates or updates a balance within an existing DB transaction
func (r *BalanceRepository) UpsertWithTx(ctx context.Context, tx *ent.Tx, params UpsertBalanceParams) error {
	updated, err := tx.Balance.
		Update().
		Where(
			balance.UserID(params.UserID),
			balance.CurrencyEQ(params.Currency),
		).
		AddAmount(params.Amount).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if updated != 0 {
		return nil
	}
	if err != nil {
		if errors.IsNegativeBalanceConstraintError(err) {
			return errors.ErrInsufficientFunds
		}

		return fmt.Errorf("failed upserting %s balance: %w", params.Currency, err)
	}

	// If the balance is not found, create a new one
	if updated == 0 {
		if params.Amount < 0 {
			return fmt.Errorf("failed creating %s balance: insufficient initial funds", params.Currency)
		}

		_, err = tx.Balance.
			Create().
			SetUserID(params.UserID).
			SetCurrency(params.Currency).
			SetAmount(params.Amount).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed creating %s balance: %w", params.Currency, err)
		}
		return nil
	}

	// If another error occurred
	return fmt.Errorf("failed upserting %s balance: %w", params.Currency, err)
}

// GetAllByUserID returns all balances of a user
func (r *BalanceRepository) GetAllByUserID(ctx context.Context, userID int) ([]*ent.Balance, error) {
	balances, err := r.client.Balance.
		Query().
		Where(balance.UserID(userID)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying balances by user_id: %w", err)
	}

	return balances, nil
}

// GetAllByUserIDUsingEdge returns all balances of a user through the edge in the graph
func (r *BalanceRepository) GetAllByUserIDUsingEdge(ctx context.Context, userID int) ([]*ent.Balance, error) {
	balances, err := r.client.User.
		Query().
		Where(user.ID(userID)).
		QueryBalances().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying user balances: %w", err)
	}

	return balances, nil
}
