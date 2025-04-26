package repository

import (
	"context"
	"fmt"
	"time"

	"db/ent"
	"db/ent/user"

	"github.com/google/uuid"
)

// UserRepository представляет репозиторий для работы с пользователями
type UserRepository struct {
	client *ent.Client
}

// NewUserRepository создаёт новый репозиторий пользователей
func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

// Create создаёт нового пользователя с использованием SQL транзакции
func (r *UserRepository) Create(ctx context.Context, name string, email string, age int) (*ent.User, error) {
	// Start a transaction
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed starting transaction: %w", err)
	}

	// Execute the actual logic within the transaction
	user, err := r.CreateWithTx(ctx, tx, name, email, age)
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

	return user, nil
}

// CreateWithTx создаёт нового пользователя внутри существующей SQL транзакции
func (r *UserRepository) CreateWithTx(ctx context.Context, tx *ent.Tx, name string, email string, age int) (*ent.User, error) {
	u, err := tx.User.
		Create().
		SetName(name).
		SetEmail(email).
		SetAge(age).
		SetCreatedAt(time.Now()).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	return u, nil
}

// CreateRandom creates a user with a random email
func (r *UserRepository) CreateRandom(ctx context.Context) (*ent.User, error) {
	return r.Create(ctx, "John Doe", uuid.New().String(), 30)
}

// GetByID gets a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int) (*ent.User, error) {
	u, err := r.client.User.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed querying user by ID: %w", err)
	}
	return u, nil
}

// GetWithTransactions gets a user together with its transactions
func (r *UserRepository) GetWithTransactions(ctx context.Context, id int) (*ent.User, error) {
	userWithTx, err := r.client.User.
		Query().
		Where(user.ID(id)).
		WithTransactions().
		Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying user with transactions: %w", err)
	}

	return userWithTx, nil
}
