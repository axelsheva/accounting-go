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

// Upsert creates a new balance or updates an existing one for a specified user and currency
func (r *BalanceRepository) Upsert(ctx context.Context, userID int, currency string, amount float64) error {
	updated, err := r.client.Balance.
		Update().
		Where(
			balance.UserID(userID),
			balance.CurrencyEQ(currency),
		).
		AddAmount(amount).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		if errors.IsNegativeBalanceConstraintError(err) {
			return errors.ErrInsufficientFunds
		}

		return fmt.Errorf("failed upserting %s balance: %w", currency, err)
	}
	if updated != 0 {
		return nil
	}

	// If the balance is not found, create a new one
	if updated == 0 {
		if amount < 0 {
			return fmt.Errorf("failed creating %s balance: %w", currency, err)
		}

		_, err = r.client.Balance.
			Create().
			SetUserID(userID).
			SetCurrency(currency).
			SetAmount(amount).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed creating %s balance: %w", currency, err)
		}
		return nil
	}

	// If another error occurred
	return fmt.Errorf("failed upserting %s balance: %w", currency, err)
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
