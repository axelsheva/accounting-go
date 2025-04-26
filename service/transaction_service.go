package service

import (
	"context"
	"db/ent"
	"db/ent/transaction"
	"db/repository"
	"fmt"
	"time"
)

// TransactionService представляет сервис для работы с транзакциями
type TransactionService struct {
	txRepo *repository.TransactionRepository
}

// NewTransactionService создает новый сервис транзакций
func NewTransactionService(client *ent.Client) *TransactionService {
	return &TransactionService{
		txRepo: repository.NewTransactionRepository(client),
	}
}

// CreateSampleTransactions создает набор тестовых транзакций для пользователя
func (s *TransactionService) CreateSampleTransactions(ctx context.Context, user *ent.User) ([]*ent.Transaction, error) {
	transactions, err := s.txRepo.CreateSampleTransactions(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("transaction service - create sample transactions: %w", err)
	}
	return transactions, nil
}

// GetTransactionByID получает транзакцию по её ID
func (s *TransactionService) GetTransactionByID(ctx context.Context, id string) (*ent.Transaction, error) {
	tx, err := s.txRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("transaction service - get transaction by id: %w", err)
	}
	return tx, nil
}

// GetAllTransactionsByUserID получает все транзакции пользователя
func (s *TransactionService) GetAllTransactionsByUserID(ctx context.Context, userID int) ([]*ent.Transaction, error) {
	txs, err := s.txRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("transaction service - get all transactions by user id: %w", err)
	}
	return txs, nil
}

// GetTransactionsByStatus получает все транзакции с указанным статусом
func (s *TransactionService) GetTransactionsByStatus(ctx context.Context, status string) ([]*ent.Transaction, error) {
	txs, err := s.txRepo.GetAllByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("transaction service - get transactions by status: %w", err)
	}
	return txs, nil
}

// TestIdempotency демонстрирует, как ID транзакции предотвращает дублирование транзакций
func (s *TransactionService) TestIdempotency(ctx context.Context, user *ent.User) error {
	// Создаем фиксированный ID для демонстрации
	fixedID := s.txRepo.GenerateTransactionID()

	// Первая попытка создания транзакции - должна успешно пройти
	fmt.Println("\n--- Testing Idempotency ---")
	fmt.Println("First attempt with fixed transaction ID:", fixedID)

	now := time.Now()
	tx1, err := s.txRepo.Create(ctx, fixedID, user.ID, 500.00, "USD", transaction.TypeDeposit, "Idempotency test", "completed", &now)
	if err != nil {
		return fmt.Errorf("failed first attempt: %w", err)
	}
	fmt.Printf("First transaction created successfully: %v\n", tx1)

	// Вторая попытка с тем же ID - должна завершиться ошибкой нарушения ограничения
	fmt.Println("\nSecond attempt with same transaction ID:", fixedID)
	_, err = s.txRepo.Create(ctx, fixedID, user.ID, 500.00, "USD", transaction.TypeDeposit, "Idempotency test - duplicate", "completed", &now)

	if err != nil {
		fmt.Printf("Second attempt failed as expected: %v\n", err)
		fmt.Println("This demonstrates idempotency - can't insert the same transaction ID twice")
		return nil
	}

	// Если мы попали сюда, что-то пошло не так - вторая попытка должна была завершиться неудачей
	return fmt.Errorf("idempotency test failed: duplicate transaction with same ID was allowed")
}

// QueryTransactions демонстрирует, как запрашивать транзакции из базы данных
func (s *TransactionService) QueryTransactions(ctx context.Context, userID int) error {
	// Находим все транзакции
	transactions, err := s.txRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed querying transactions: %w", err)
	}
	fmt.Printf("Found %d transactions\n", len(transactions))

	// Находим завершенные транзакции
	completed, err := s.GetTransactionsByStatus(ctx, "completed")
	if err != nil {
		return fmt.Errorf("failed querying completed transactions: %w", err)
	}
	fmt.Printf("Found %d completed transactions\n", len(completed))

	// Находим ожидающие транзакции
	pending, err := s.GetTransactionsByStatus(ctx, "pending")
	if err != nil {
		return fmt.Errorf("failed querying pending transactions: %w", err)
	}
	fmt.Printf("Found %d pending transactions\n", len(pending))

	// Получаем транзакции для конкретного пользователя через репозиторий
	userTxs, err := s.txRepo.GetAllByUserIDUsingEdge(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed querying user with transactions: %w", err)
	}
	fmt.Printf("User '%d' has %d transactions\n", userID, len(userTxs))

	return nil
}
