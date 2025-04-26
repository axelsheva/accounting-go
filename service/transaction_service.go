package service

import (
	"context"
	"accounting/ent"
	"accounting/ent/transaction"
	"accounting/repository"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// TransactionService presents a service for working with transactions
type TransactionService struct {
	txRepo      *repository.TransactionRepository
	balanceRepo *repository.BalanceRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(client *ent.Client) *TransactionService {
	balanceRepo := repository.NewBalanceRepository(client)

	return &TransactionService{
		balanceRepo: balanceRepo,
		txRepo:      repository.NewTransactionRepository(client, balanceRepo),
	}
}

func (s *TransactionService) Create(ctx context.Context, id string, userID int, currency string, amount float64, txType transaction.Type) (*ent.Transaction, error) {
	tx, err := s.txRepo.Create(ctx, id, userID, amount, currency, txType)
	if err != nil {
		return nil, fmt.Errorf("transaction service - create transaction: %w", err)
	}
	return tx, nil
}

// GetTransactionByID gets a transaction by its ID
func (s *TransactionService) GetTransactionByID(ctx context.Context, id string) (*ent.Transaction, error) {
	tx, err := s.txRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("transaction service - get transaction by id: %w", err)
	}
	return tx, nil
}

// GetAllTransactionsByUserID gets all transactions of a user
func (s *TransactionService) GetAllTransactionsByUserID(ctx context.Context, userID int) ([]*ent.Transaction, error) {
	txs, err := s.txRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("transaction service - get all transactions by user id: %w", err)
	}
	return txs, nil
}

// TestIdempotency demonstrates how the transaction ID prevents duplicate transactions
func (s *TransactionService) TestIdempotency(ctx context.Context, user *ent.User) error {
	// Create a fixed ID for demonstration
	fixedID := uuid.New().String()

	// First attempt to create a transaction - should succeed
	fmt.Println("\n--- Testing Idempotency ---")
	fmt.Println("First attempt with fixed transaction ID:", fixedID)

	tx1, err := s.txRepo.Create(ctx, fixedID, user.ID, 500.00, "USD", transaction.TypeDeposit)
	if err != nil {
		return fmt.Errorf("failed first attempt: %w", err)
	}
	fmt.Printf("First transaction created successfully: %v\n", tx1)

	// Second attempt with the same ID - should fail due to the constraint
	fmt.Println("\nSecond attempt with same transaction ID:", fixedID)
	_, err = s.txRepo.Create(ctx, fixedID, user.ID, 500.00, "USD", transaction.TypeDeposit)

	if err != nil {
		fmt.Printf("Second attempt failed as expected: %v\n", err)
		fmt.Println("This demonstrates idempotency - can't insert the same transaction ID twice")
		return nil
	}

	// If we got here, something went wrong - the second attempt should have failed
	return fmt.Errorf("idempotency test failed: duplicate transaction with same ID was allowed")
}

// QueryTransactions demonstrates how to query transactions from the database
func (s *TransactionService) QueryTransactions(ctx context.Context, userID int) error {
	// Find all transactions
	transactions, err := s.txRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed querying transactions: %w", err)
	}
	fmt.Printf("Found %d transactions\n", len(transactions))

	// Get transactions for a specific user through the repository
	userTxs, err := s.txRepo.GetAllByUserIDUsingEdge(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed querying user with transactions: %w", err)
	}
	fmt.Printf("User '%d' has %d transactions\n", userID, len(userTxs))

	return nil
}

// DemonstrateBalanceOperations demonstrates balance operations
func (s *TransactionService) DemonstrateBalanceOperations(ctx context.Context, userID int) error {
	fmt.Println("\n--- Demonstrating Balance Operations ---")

	// Get initial EUR balance
	eurBalance, err := s.balanceRepo.GetByUserIDAndCurrency(ctx, userID, "EUR")
	if err != nil {
		return fmt.Errorf("failed querying EUR balance: %w", err)
	}

	initialAmount := eurBalance.Amount
	fmt.Printf("Initial EUR balance: %.2f\n", initialAmount)

	// Increment the balance by 100 EUR
	incrementAmount := 100.0
	_, err = s.txRepo.Create(ctx, uuid.New().String(), userID, incrementAmount, "EUR", transaction.TypeDeposit)
	if err != nil {
		return fmt.Errorf("failed incrementing balance: %w", err)
	}

	// Get updated balance
	eurBalance, err = s.balanceRepo.GetByUserIDAndCurrency(ctx, userID, "EUR")
	if err != nil {
		return fmt.Errorf("failed querying updated EUR balance: %w", err)
	}

	fmt.Printf("EUR balance after increment of %.2f: %.2f\n", incrementAmount, eurBalance.Amount)

	// Decrement the balance
	decrementAmount := 50.0
	_, err = s.txRepo.Create(ctx, uuid.New().String(), userID, decrementAmount, "EUR", transaction.TypeWithdrawal)
	if err != nil {
		return fmt.Errorf("failed decrementing balance: %w", err)
	}

	// Get final balance
	eurBalance, err = s.balanceRepo.GetByUserIDAndCurrency(ctx, userID, "EUR")
	if err != nil {
		return fmt.Errorf("failed querying final EUR balance: %w", err)
	}

	fmt.Printf("EUR balance after decrement of %.2f: %.2f\n", decrementAmount, eurBalance.Amount)

	// Try to decrement too much (should fail)
	tooMuchAmount := eurBalance.Amount + 1000.0
	_, err = s.txRepo.Create(ctx, uuid.New().String(), userID, tooMuchAmount, "EUR", transaction.TypeWithdrawal)
	if err != nil {
		fmt.Printf("As expected, decrementing too much (%.2f) failed: %v\n", tooMuchAmount, err)
	} else {
		return errors.New("large withdrawal should have failed but didn't")
	}

	return nil
}
