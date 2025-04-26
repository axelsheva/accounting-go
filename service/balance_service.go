package service

import (
	"context"
	"db/ent"
	"db/repository"
	"errors"
	"fmt"
)

// BalanceService represents a service for working with balances
type BalanceService struct {
	client      *ent.Client
	balanceRepo *repository.BalanceRepository
	txRepo      *repository.TransactionRepository
}

// NewBalanceService creates a new balance service
func NewBalanceService(client *ent.Client) *BalanceService {
	return &BalanceService{
		client:      client,
		balanceRepo: repository.NewBalanceRepository(client),
		txRepo:      repository.NewTransactionRepository(client),
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

// UpsertBalance creates or updates the balance of a user in a specified currency
func (s *BalanceService) UpsertBalance(ctx context.Context, userID int, currency string, amount float64) error {
	// Call the Upsert method from the repository
	err := s.balanceRepo.Upsert(ctx, userID, currency, amount)
	if err != nil {
		return err
	}

	return nil
}

// DemonstrateBalanceOperations demonstrates balance operations
func (s *BalanceService) DemonstrateBalanceOperations(ctx context.Context, userID int) error {
	fmt.Println("\n--- Demonstrating Balance Operations ---")

	// Get initial EUR balance
	eurBalance, err := s.GetUserBalance(ctx, userID, "EUR")
	if err != nil {
		return fmt.Errorf("failed querying EUR balance: %w", err)
	}

	initialAmount := eurBalance.Amount
	fmt.Printf("Initial EUR balance: %.2f\n", initialAmount)

	// Increment the balance by 100 EUR
	incrementAmount := 100.0
	if err := s.UpsertBalance(ctx, userID, "EUR", incrementAmount); err != nil {
		return fmt.Errorf("failed incrementing balance: %w", err)
	}

	// Get updated balance
	eurBalance, err = s.GetUserBalance(ctx, userID, "EUR")
	if err != nil {
		return fmt.Errorf("failed querying updated EUR balance: %w", err)
	}

	fmt.Printf("EUR balance after increment of %.2f: %.2f\n", incrementAmount, eurBalance.Amount)

	// Decrement the balance
	decrementAmount := -50.0
	if err := s.UpsertBalance(ctx, userID, "EUR", decrementAmount); err != nil {
		return fmt.Errorf("failed decrementing balance: %w", err)
	}

	// Get final balance
	eurBalance, err = s.GetUserBalance(ctx, userID, "EUR")
	if err != nil {
		return fmt.Errorf("failed querying final EUR balance: %w", err)
	}

	fmt.Printf("EUR balance after decrement of %.2f: %.2f\n", decrementAmount, eurBalance.Amount)

	// Try to decrement too much (should fail)
	tooMuchAmount := -eurBalance.Amount - 1000.0
	err = s.UpsertBalance(ctx, userID, "EUR", tooMuchAmount)
	if err != nil {
		fmt.Printf("As expected, decrementing too much (%.2f) failed: %v\n", tooMuchAmount, err)
	} else {
		return errors.New("large withdrawal should have failed but didn't")
	}

	return nil
}
