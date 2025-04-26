package main

import (
	"context"
	"fmt"
	"log"

	"db/ent"
	"db/service"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	// PostgreSQL connection string
	// Replace with your actual database credentials
	dsn := "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable"

	// Create an ent client
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	// Run the auto migration tool to create all schema resources
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Working with the client
	ctx := context.Background()

	// Create services
	userService := service.NewUserService(client)
	balanceService := service.NewBalanceService(client)
	transactionService := service.NewTransactionService(client)

	// Create a user
	user, err := userService.CreateRandomUser(ctx)
	if err != nil {
		log.Fatalf("failed creating user: %v", err)
	}
	fmt.Printf("Created user with ID: %d\n", user.ID)

	// Create balances for the user
	_, err = transactionService.Create(ctx, uuid.New().String(), user.ID, "USD", 1000.00, "deposit")
	if err != nil {
		log.Fatalf("failed creating balances: %v", err)
	}
	_, err = transactionService.Create(ctx, uuid.New().String(), user.ID, "EUR", 500.00, "deposit")
	if err != nil {
		log.Fatalf("failed creating balances: %v", err)
	}
	_, err = transactionService.Create(ctx, uuid.New().String(), user.ID, "RUB", 50000.00, "deposit")
	if err != nil {
		log.Fatalf("failed creating balances: %v", err)
	}

	// Query transactions
	if err := transactionService.QueryTransactions(ctx, user.ID); err != nil {
		log.Fatalf("failed querying transactions: %v", err)
	}

	// Demonstrate idempotency with same transaction ID
	if err := transactionService.TestIdempotency(ctx, user); err != nil {
		log.Fatalf("failed testing idempotency: %v", err)
	}

	// Query balances
	balances, err := balanceService.GetUserBalances(ctx, user.ID)
	if err != nil {
		log.Fatalf("failed querying balances: %v", err)
	}
	fmt.Printf("\n--- Querying Balances ---\n")
	fmt.Printf("User with ID %d has %d balances:\n", user.ID, len(balances))
	for _, b := range balances {
		fmt.Printf("- %s: %.2f\n", b.Currency, b.Amount)
	}

	// Demonstrate balance update with a transaction
	if err := UpdateBalanceWithTransaction(ctx, client, balanceService, transactionService, user.ID); err != nil {
		log.Fatalf("failed updating balance: %v", err)
	}

	// Demonstrate queries using direct user_id field
	if err := QueryByUserID(ctx, client, balanceService, transactionService, user.ID); err != nil {
		log.Fatalf("failed querying by user_id: %v", err)
	}

	// Demonstrate balance increment and decrement
	if err := transactionService.DemonstrateBalanceOperations(ctx, user.ID); err != nil {
		log.Fatalf("failed demonstrating balance operations: %v", err)
	}

	fmt.Println("Ent with PostgreSQL setup completed successfully!")
}

// QueryByID demonstrates how to query entities by their IDs
func QueryByID(ctx context.Context, userService *service.UserService, txService *service.TransactionService, userID int, transactionID string) error {
	// Get user by ID
	u, err := userService.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed querying user by ID: %w", err)
	}
	fmt.Printf("Found user by ID %d: %s (%s)\n", u.ID, u.Name, u.Email)

	// Get transaction by ID
	tx, err := txService.GetTransactionByID(ctx, transactionID)
	if err != nil {
		return fmt.Errorf("failed querying transaction by ID: %w", err)
	}
	fmt.Printf("Found transaction by ID %s: amount: %.2f %s\n",
		tx.ID, tx.Amount, tx.Currency)

	// Get all transactions for a user by user ID
	userWithTx, err := userService.GetUserWithTransactions(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed querying user's transactions: %w", err)
	}
	fmt.Printf("User with ID %d has %d transactions\n", userID, len(userWithTx.Edges.Transactions))

	// Get user for a transaction by transaction ID
	txUser, err := tx.QueryUser().Only(ctx)
	if err != nil {
		return fmt.Errorf("failed querying transaction's user: %w", err)
	}
	fmt.Printf("Transaction with ID %s belongs to user: %s (ID: %d)\n",
		transactionID, txUser.Name, txUser.ID)

	return nil
}

// UpdateBalanceWithTransaction demonstrates updating a balance when a transaction occurs
func UpdateBalanceWithTransaction(ctx context.Context, client *ent.Client, balanceService *service.BalanceService, transactionService *service.TransactionService, userID int) error {
	fmt.Println("\n--- Updating Balance with Transaction ---")

	// Get the user's USD balance
	usdBalance, err := balanceService.GetUserBalance(ctx, userID, "USD")
	if err != nil {
		return fmt.Errorf("failed querying user's USD balance: %w", err)
	}

	fmt.Printf("Initial USD balance: %.2f\n", usdBalance.Amount)

	// Deposit amount - use IncrementBalance from the service
	depositAmount := 250.0
	_, err = transactionService.Create(ctx, uuid.New().String(), userID, "USD", depositAmount, "deposit")
	if err != nil {
		return fmt.Errorf("failed incrementing balance: %w", err)
	}

	// Get updated balance
	updatedBalance, err := balanceService.GetUserBalance(ctx, userID, "USD")
	if err != nil {
		return fmt.Errorf("failed querying updated balance: %w", err)
	}

	fmt.Printf("Updated USD balance after deposit of %.2f: %.2f\n", depositAmount, updatedBalance.Amount)

	return nil
}

// QueryByUserID demonstrates how to query entities using the direct user_id field
func QueryByUserID(ctx context.Context, client *ent.Client, balanceService *service.BalanceService, txService *service.TransactionService, userID int) error {
	fmt.Println("\n--- Querying By User ID Field ---")

	// Get all balances for user using user_id field directly
	balances, err := balanceService.GetUserBalances(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed querying balances by user_id: %w", err)
	}
	fmt.Printf("Found %d balances using user_id field directly\n", len(balances))

	// Get all transactions for user using user_id field directly
	transactions, err := txService.GetAllTransactionsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed querying transactions by user_id: %w", err)
	}
	fmt.Printf("Found %d transactions using user_id field directly\n", len(transactions))

	// Get USD balance using user_id and currency directly
	usdBalance, err := balanceService.GetUserBalance(ctx, userID, "USD")
	if err != nil {
		return fmt.Errorf("failed querying USD balance by user_id: %w", err)
	}
	fmt.Printf("Found USD balance with amount %.2f using user_id field directly\n", usdBalance.Amount)

	// Get all transactions for user using user_id field directly
	transactions, err = txService.GetAllTransactionsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed querying transactions: %w", err)
	}
	fmt.Printf("Found %d transactions\n", len(transactions))

	return nil
}
