package service

import (
	"context"
	"db/ent"
	"db/repository"
	"fmt"
)

// BalanceService represents a service for working with balances
type BalanceService struct {
	client      *ent.Client
	balanceRepo *repository.BalanceRepository
}

// NewBalanceService creates a new balance service
func NewBalanceService(client *ent.Client) *BalanceService {
	return &BalanceService{
		client:      client,
		balanceRepo: repository.NewBalanceRepository(client),
	}
}

// GetUserBalances gets all balances of a user
func (s *BalanceService) GetUserBalances(ctx context.Context, userID int) ([]*ent.Balance, error) {
	balances, err := s.balanceRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("balance service - get user balances: %w", err)
	}
	return balances, nil
}

// GetUserBalance gets the balance of a user in a specified currency
func (s *BalanceService) GetUserBalance(ctx context.Context, userID int, currency string) (*ent.Balance, error) {
	balance, err := s.balanceRepo.GetByUserIDAndCurrency(ctx, userID, currency)
	if err != nil {
		return nil, fmt.Errorf("balance service - get user balance: %w", err)
	}
	return balance, nil
}
